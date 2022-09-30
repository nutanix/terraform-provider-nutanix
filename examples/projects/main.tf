terraform{
    required_providers {
        nutanix = {
            source  = "nutanix/nutanix"
            version = "1.7.1"
        }
    }
}

#defining nutanix configuration
provider "nutanix"{
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port = 9440
  insecure = true
}

# set use_project_internal flag to use user-role mapping 

data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

# ### Define Terraform Managed Subnets
resource "nutanix_subnet" "infra-managed-network-140" {
  # What cluster will this VLAN live on?
  cluster_uuid = local.cluster1

  # General Information
  name        = "infra-managed-network-140"
  vlan_id     = 140
  subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
  subnet_ip = "10.xx.xx.xx"

  default_gateway_ip = "10.xx.xx.xx"
  prefix_length      = 24

  dhcp_options = {
    boot_file_name   = "bootfile"
    domain_name      = "lab"
    tftp_server_name = "10.xx.xx.xx"
  }

  dhcp_server_address = {
    ip = "10.xx.xx.xx"
  }

  dhcp_domain_name_server_list = ["10.xx.xx.xx"]
  dhcp_domain_search_list      = ["ntnxlab.local"]
  #ip_config_pool_list_ranges   = ["10.xx.xx.xx 10.xx.xx.xx"] 
}

# Note: user reference and acp->user_reference should be same for mapping the role. Also whenever acp is given
# it's mandate to provide cluster_uuid to get the filter context list and scope of each defined user.

resource "nutanix_project" "testp1" {
    name        = "testProj"
    description = "test project description"

    # cluster uuid is required to map acp in projects
    cluster_uuid = "${local.cluster1}" 

    # set this use_project_internal flag for using projects_internal API 
    use_project_internal=true

    # set project collaboration, default it is true
    enable_collab = true
    default_subnet_reference{
        kind="subnet"
        uuid=resource.nutanix_subnet.sub.id
    }
    user_reference_list{
      name= "{{user_name}}"
      kind= "user"
      uuid= "{{user_uuid}}"
    }
    subnet_reference_list{
          uuid=resource.nutanix_subnet.sub.id
    }
    acp{
        # acp name consists name_uuid string, it should be different for each acp. 
        name="{{acp_name}}"
        role_reference{
            kind= "role"
            uuid= "{{role_uuid}}"
            name="Developer"
        }
        user_reference_list{
            name= "{{user_name}}"
            kind= "user"
            uuid= "{{user_uuid}}"
        }
        description= "descripton"
    }
    api_version = "3.1"

  # to enable project quotas
    project_quota{
      vcpu = 1
      disk = 2147483648
      memory = 2147483648
    }
}
