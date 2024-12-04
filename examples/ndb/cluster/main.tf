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

## resource to register cluster in ndb

resource "nutanix_ndb_cluster" "clsname" {
  name= "cls-name"
  description = "cluster description"
  cluster_ip = "{{ clusterIP }}"
  username= "{{ cluster username }}"
  password = "{{ cluster password }}"
  storage_container = "{{ storage container}}"
  agent_network_info{
	dns = "{{ dns }}"
	ntp = "{{ ntp }}"
  }
  networks_info{
	type = "DHCP"
	network_info{
		vlan_name = "vlan_static"
		static_ip = "{{ staticIP }}"
		gateway = "{{ Gateway }}"
		subnet_mask="{{ subnetMask }}"
	}
	access_type = [
        "PRISM",
        "DSIP",
        "DBSERVER"
      ]
  }
}