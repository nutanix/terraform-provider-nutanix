terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.8.0"
        }
    }
}

#defining nutanix configuration
provider "nutanix"{
  ndb_username = var.ndb_username
  ndb_password = var.ndb_password
  ndb_endpoint = var.ndb_endpoint
  insecure = true
}


#resource to create network
resource "nutanix_ndb_network" "name" {
    name= "test-sub"
    type="Static"
    cluster_id = "{{ cluster_id }}"
    gateway= "{{ gatway for the vlan }}"
    subnet_mask = "{{ subnet mask for the vlan}}"
    primary_dns = " {{ primary dns for the vlan }}"
    secondary_dns= "{{secondary dns for the vlan }}"
    ip_pools{
        start_ip = "{{ starting address range}}"
        end_ip = "{{ ending address range }}"
    }
}

#data source to get network
data "nutanix_ndb_network" "net"{
    name = "{{ name of network }}"
}

data "nutanix_ndb_network" "net"{
    id = "{{ id of network }}"
}

#data source to get List of networks
data "nutanix_ndb_networks" "nets"{ }