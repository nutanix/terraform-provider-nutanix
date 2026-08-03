
resource "nutanix_images_v2" "windows-image" {
  name = local.images.windows_image
  type = "DISK_IMAGE"
  source {
    url_source {
      url = local.images.windows_image_url
    }
  }
}

# Ubuntu + Rocky8/NGT gold images: reuse-if-exists. prepare_pc.py may already
# have registered these by name (its idempotent _ensure_image), so look them up
# first and only create when absent -- this stops the two provisioners from
# producing duplicate same-named images regardless of run order. VMs reference
# these via local.ubuntu_image_ext_id / local.ngt_image_ext_id below (NOT the
# resource .id, which does not exist when count = 0).
data "nutanix_images_v2" "ubuntu-image" {
  filter = "name eq '${local.images.ubuntu_image}'"
}

resource "nutanix_images_v2" "ubuntu-image" {
  count = length(data.nutanix_images_v2.ubuntu-image.images) == 0 ? 1 : 0
  name  = local.images.ubuntu_image
  type  = "DISK_IMAGE"
  source {
    url_source {
      url = local.images.ubuntu_image_url
    }
  }
}

data "nutanix_images_v2" "ngt-image" {
  filter = "name eq '${local.images.ngt_image}'"
}

resource "nutanix_images_v2" "ngt-image" {
  count = length(data.nutanix_images_v2.ngt-image.images) == 0 ? 1 : 0
  name  = local.images.ngt_image
  type  = "DISK_IMAGE"
  source {
    url_source {
      url = local.images.ngt_image_url
    }
  }
}

locals {
  # Prefer the pre-existing (data-source) image; fall back to the one created above.
  ubuntu_image_ext_id = length(data.nutanix_images_v2.ubuntu-image.images) > 0 ? data.nutanix_images_v2.ubuntu-image.images[0].ext_id : nutanix_images_v2.ubuntu-image[0].id
  ngt_image_ext_id    = length(data.nutanix_images_v2.ngt-image.images) > 0 ? data.nutanix_images_v2.ngt-image.images[0].ext_id : nutanix_images_v2.ngt-image[0].id
}

# NKP deploy/bastion VM.
# An Ubuntu VM cloned from the ubuntu gold image
# (local.ubuntu_image_ext_id) used to run the `nkp create cluster nutanix`
# workflow. NKP bootstraps a local KIND (docker) cluster, so this VM is sized
# with extra vCPU/RAM/disk. cloud-init sets the 'ubuntu' user password and
# enables SSH password auth so the automation (e2e_nkp_deploy.py) can log in.
resource "nutanix_virtual_machine_v2" "nkp-bastion" {
  name                 = "nkp-bastion"
  description          = "NKP deploy bastion (runs nkp create cluster nutanix)"
  num_sockets          = 4
  num_cores_per_socket = 1
  memory_size_bytes    = 8 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster_ext_id
  }
  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        data_source {
          reference {
            image_reference {
              image_ext_id = local.ubuntu_image_ext_id
            }
          }
        }
        disk_size_bytes = 60 * 1024 * 1024 * 1024
      }
    }
  }
  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = data.nutanix_subnets_v2.vmm-subnet.subnets[0].ext_id
        }
        vlan_mode = "ACCESS"
      }
    }
  }
  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
  guest_customization {
    config {
      cloud_init {
        datasource_type = "CONFIG_DRIVE_V2"
        cloud_init_script {
          user_data {
            value = base64encode(<<-EOT
              #cloud-config
              ssh_pwauth: true
              chpasswd:
                expire: false
                list: |
                  ubuntu:Nutanix.123
            EOT
            )
          }
        }
      }
    }
  }
  power_state = "ON"
  lifecycle {
    ignore_changes = [guest_customization, guest_tools, cd_roms]
  }
  depends_on = [nutanix_images_v2.ubuntu-image, data.nutanix_subnets_v2.vmm-subnet]
}

# VM integration test VM with overlay subnet
resource "nutanix_virtual_machine_v2" "vm-overlay-subnet" {
  name              = local.integration_vm.name
  is_agent_vm       = false
  num_sockets       = 1
  memory_size_bytes = 4 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster_ext_id
  }
  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
  nics {
    nic_backing_info {
      virtual_ethernet_nic {
        is_connected = true
      }
    }
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        ipv4_config {
          ip_address {
            value = local.integration_vm.ip_address
          }
          should_assign_ip = true
        }
        subnet {
          ext_id = nutanix_subnet_v2.overlay-subnet.id
        }
        vlan_mode = "ACCESS"
      }
    }
  }
  power_state = "ON"
  lifecycle {
    ignore_changes = [nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.ipv4_config.0.should_assign_ip]
  }
}

# Subnet the NGT integration VM attaches to (looked up by name, e.g. vlan.800).
data "nutanix_subnets_v2" "vmm-subnet" {
  filter = "name eq '${local.conf.vmm.subnet_name}'"
}

data "nutanix_storage_containers_v2" "ngt-sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

# Rocky8 integration VM with NGT installed (SELF_SERVICE_RESTORE + VSS_SNAPSHOT),
# looked up by name (vmm.ngt_vm.name) by the ngt_configuration / ngt_installation
# data-source tests. Distinct from vmm.ngt.ngt_upgrade_vm_name (that VM carries an
# OLD NGT for the upgrade test and is created by prepare_pc.py).
resource "nutanix_virtual_machine_v2" "ngt-vm" {
  name                 = local.ngt_vm.name
  description          = "integration test vm with ngt image"
  num_cores_per_socket = 1
  num_sockets          = 1
  memory_size_bytes    = 4 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster_ext_id
  }

  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        data_source {
          reference {
            image_reference {
              image_ext_id = local.ngt_image_ext_id
            }
          }
        }
        disk_size_bytes = 20 * 1024 * 1024 * 1024
      }
    }
  }

  cd_roms {
    disk_address {
      bus_type = "IDE"
      index    = 0
    }
  }

  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = data.nutanix_subnets_v2.vmm-subnet.subnets[0].ext_id
        }
        vlan_mode = "ACCESS"
      }
    }
  }

  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }

  power_state = "ON"

  lifecycle {
    ignore_changes = [guest_tools]
  }

  depends_on = [data.nutanix_clusters_v2.clusters, nutanix_images_v2.ngt-image, data.nutanix_storage_containers_v2.ngt-sc]
}

resource "nutanix_ngt_installation_v2" "install-ngt-integration" {
  ext_id = nutanix_virtual_machine_v2.ngt-vm.id
  credential {
    username = local.ngt_vm.username
    password = local.ngt_vm.password
  }
  reboot_preference {
    schedule_type = "IMMEDIATE"
  }
  capablities = ["SELF_SERVICE_RESTORE", "VSS_SNAPSHOT"]
  depends_on  = [nutanix_virtual_machine_v2.ngt-vm]
  lifecycle {
    ignore_changes = [capablities]
  }
}



# gc_profile vm 
# VM Create
resource "nutanix_virtual_machine_v2" "gc_profile_vm" {
  name                 = local.gc_profile.vm_name
  num_sockets          = 2
  num_cores_per_socket = 1
  memory_size_bytes    = 3 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster_ext_id
  }
  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        data_source {
          reference {
            image_reference {
              image_ext_id = nutanix_images_v2.windows-image.id
            }
          }
        }
        disk_size_bytes = 40 * 1024 * 1024 * 1024
      }
    }
  }
  cd_roms {
    disk_address {
      bus_type = "SATA"
      index    = 0
    }
  }
  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = local.gc_subnet_ext_id
        }
        vlan_mode                 = "ACCESS"
        should_allow_unknown_macs = false
        ipv4_config {
          should_assign_ip = true
        }
      }
    }
    nic_backing_info {
      virtual_ethernet_nic {
        is_connected = true
      }
    }
  }
  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
  power_state = "ON"
  lifecycle {
    ignore_changes = [
      guest_tools,
      cd_roms
    ]
  }
}


# Give the Windows guest time to finish sysprep/OOBE (which applies the
# Administrator password and runs the WinRM FirstLogonCommands) before NGT tries
# to authenticate -- otherwise the install fails with "credentials are invalid".
resource "time_sleep" "wait-for-gc-profile-sysprep" {
  depends_on      = [nutanix_virtual_machine_v2.gc_profile_vm]
  create_duration = "2m"
}

resource "nutanix_ngt_insert_iso_v2" "gc_profile_vm_insert-iso" {
  ext_id         = nutanix_virtual_machine_v2.gc_profile_vm.id
  capablities    = ["VSS_SNAPSHOT", "SELF_SERVICE_RESTORE"]
  is_config_only = false
  depends_on     = [time_sleep.wait-for-gc-profile-sysprep]
  lifecycle {
    ignore_changes = [capablities]
  }
}

resource "nutanix_ngt_installation_v2" "install-ngt-gc-profile" {
  ext_id      = nutanix_virtual_machine_v2.gc_profile_vm.id
  capablities = ["VSS_SNAPSHOT", "SELF_SERVICE_RESTORE"]
  credential {
    username = local.gc_profile.default_image_username
    password = local.gc_profile.default_image_password
  }
  reboot_preference {
    schedule_type = "IMMEDIATE"
  }
  depends_on = [
  nutanix_ngt_insert_iso_v2.gc_profile_vm_insert-iso]
  lifecycle {
    ignore_changes = [capablities]
  }
}
