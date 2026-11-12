# ==============================================================================
# Tenant Admin: Infrastructure Management
# ==============================================================================
#
# What this file does:
#   Shows how a Tenant Admin creates infrastructure on a domain that the Service
#   Provider Admin has already allocated to the tenant. This example covers a
#   simple VM workflow (image + virtual machine).
#
# Flow:
#   1. Authenticate as Tenant Admin via nutanix_central_endpoint + tenant_name
#      + domain_name (see provider comments below for how to map the Tenant URL).
#   2. Discover the tenant project and its resource-group placement targets
#      (cluster and storage container from the domain allocation).
#   3. Create a disk image in the tenant project.
#   4. Create and power on a VM using that image, cluster, and storage container.
#
# Prerequisites:
#   1. SP Admin has created the tenant
#   2. Domain allocation with quota / placement targets already exists
#   3. Tenant Admin credentials are available
#
# ==============================================================================

# ==============================================================================
# TERRAFORM CONFIGURATION
# ==============================================================================

terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.1"
    }
  }
}


# ==============================================================================
# PROVIDER CONFIGURATION FOR TENANT ADMIN
# ==============================================================================
#
# How to set these values from the Tenant Admin URL you receive from the SP Admin:
#
#   https://<tenant_name>.tenants.<service_provider_fqdn>
#
# Example:
#   https://my-tenant.tenants.nconprem-example.nxncm.com
#
#   - tenant_name              = "my-tenant"
#   - nutanix_central_endpoint = "nconprem-example.nxncm.com"   # the service provider FQDN only
#   - domain_name              = allocated PC domain name (from the SP Admin domain allocation)
#
# Do not put the full tenant URL into nutanix_central_endpoint — use only the
# service provider FQDN (the part after ".tenants.").
#
# Note For Service Provider Admin:
#   Calls to Nutanix Central may return 404 when the SP FQDN Host header includes
#   an explicit :443 port. This can affect Terraform and other API clients. Until
#   a fix ships, add both <host> and <host>:443 to ALLOWED_CALLBACK_HOSTS on the
#   Dion service and restart it. Details:
#   https://portal.nutanix.com/page/documents/kbs/details?targetId=kA0VO000000EVaD0AW

provider "nutanix" {
  username                 = var.tenant_username
  password                 = var.tenant_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  domain_name              = var.domain_name
  tenant_name              = var.tenant_name
  port                     = 443
}



# Fetch the Project
data "nutanix_projects_v2" "projects" {}

# Get clusters (for VM creation)
data "nutanix_resource_groups_v2" "resource_groups" {
  filter = "projectExtId eq '${data.nutanix_projects_v2.projects.projects[0].ext_id}'"
}


# ==============================================================================
# LOCAL VALUES
# ==============================================================================

locals {
  cluster_ext_id           = data.nutanix_resource_groups_v2.resource_groups.resource_groups[0].placement_targets[0].cluster_ext_id
  storage_container_ext_id = data.nutanix_resource_groups_v2.resource_groups.resource_groups[0].placement_targets[0].storage_containers[0].ext_id
}




# VM Workflow

# Image
resource "nutanix_images_v2" "image" {
  name                     = "tenant image"
  description              = "tenant admin workflow image"
  type                     = "DISK_IMAGE"
  project_ext_id           = data.nutanix_projects_v2.projects.projects[0].ext_id
  cluster_location_ext_ids = [local.cluster_ext_id]
  source {
    url_source {
      url = var.image_url
    }
  }
}

resource "nutanix_virtual_machine_v2" "database_vm" {
  name                 = "tenant admin workflow vm"
  description          = "tenant admin workflow vm"
  num_cores_per_socket = 1
  num_sockets          = 1
  memory_size_bytes    = 2 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster_ext_id
  }
  project {
    ext_id = data.nutanix_projects_v2.projects.projects[0].ext_id
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
              image_ext_id = nutanix_images_v2.image.id
            }
          }
        }
        storage_container {
          ext_id = local.storage_container_ext_id
        }
      }
    }
  }
  cd_roms {
    disk_address {
      bus_type = "IDE"
      index    = 1
    }
  }
  boot_config {
    legacy_boot {
      boot_order = ["DISK"]
    }
  }
  power_state = "ON"
  depends_on  = [nutanix_images_v2.image]
}