# config.yaml is the single source of truth for every input below (creds,
# fixture names, provisioning specs). Read it once, then derive locals.
locals {
  conf = yamldecode(file("${path.module}/../config.yaml"))

  pc_prep     = local.conf.pc_prep
  images      = local.conf.images
  gc_subnet   = local.conf.networking.gc_subnet
  dns_servers = local.conf.dns.servers
  ntp_servers = local.conf.ntp.servers
  # Base Windows VM name the guest-customization template-deploy tests look up
  # (vmm.gc_profile.vm_name). Keep this VM's name in sync with that config key.
  gc_profile = local.conf.vmm.gc_profile

  # Provider connection (was terraform.tfvars -> nutanix_*).
  pc = local.conf.pc

  # Bare PE the deploy-PC test targets (prism.deploy_pc). Used to pre-create the
  # deploy network on that PE via the aliased nutanix.deploy_pe provider.
  deploy_pc = local.conf.prism.deploy_pc

  # Provisioning inputs (was terraform.tfvars). Names come from the domain
  # blocks; the numeric/IP specs come from the terraform.* block.
  external_nat_subnet = merge(local.conf.terraform.external_nat_subnet, {
    name = local.conf.networking.external_nat_subnet
  })
  overlay_subnet = merge(local.conf.terraform.overlay_subnet, {
    name = local.conf.networking.overlay_subnet
  })
  integration_vm = {
    name       = local.conf.vmm.integration_vm
    ip_address = local.conf.terraform.integration_vm.ip_address
  }
  # Rocky8 VM with NGT installed (vmm.ngt_vm) that the ngt_configuration /
  # ngt_installation data-source tests look up by name.
  ngt_vm = local.conf.vmm.ngt_vm
  object_store = merge(local.conf.object_store, {
    subnet = merge(local.conf.object_store.subnet, {
      subnet_name = local.conf.object_store.subnet_name
    })
  })
}

# define common variables
locals {
  pcFilter  = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
  aosFilter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}


# pull cluster data
data "nutanix_clusters_v2" "clusters" {
  filter = local.aosFilter
}

# define local variables
locals {
  cluster_ext_id   = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
  gc_subnet_ext_id = nutanix_subnet_v2.vlan-255.id
}


data "nutanix_storage_containers_v2" "default_container" {
  filter = "startswith(name,'default-container-')"
}

locals {
  default_container_uuid = data.nutanix_storage_containers_v2.default_container.storage_containers[0].ext_id
}