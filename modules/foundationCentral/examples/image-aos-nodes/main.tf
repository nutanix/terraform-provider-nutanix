tterraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">1.5.0-beta.2"
    }
  }
}

provider "nutanix" {
  username = "user"
  password = "pass"
  endpoint = "10.x.xx.xx"
  insecure = true
  port     = 9440
}

// use aos-based-node-imaging based on given node-serials-filter
module "batch1" {

  // source where module code is present in local machine
  source       = "../../modules/foundationCentral/aos_based_imaging"
  cluster_name = "FC-AOS-Mod-2"

  /* 
    Mention Node-Serial => node serials of the node which have to imaged . 
    All the other info about the nodes such as cvm_ip, hypervisor_ip, netmask, etc will be fetched internally.
    Input given here will be of more priority than fetched details of nodes. 
    */
  node_info = {
    "<Node-Serial-1>" : {}
    "<Node-Serial-2>" : {}
    "<Node-Serial-3>" : {}
  }
  // Required setting to provide network settings
  common_network_settings = {
    cvm_dns_servers : [
      "10.4.8.15"
    ]
    hypervisor_dns_servers : [
      "10.4.8.15"
    ]
    cvm_ntp_servers : [
      "0.pool.ntp.org"
    ]
    hypervisor_ntp_servers : [
      "0.pool.ntp.org"
    ]
  }
  // provide hypervisor iso details
  hypervisor_iso_details = {
    url = "<HYPERVISOR-ISO-URL>"
  }

  // give aos package url to download
  aos_package_url = "<AOS-PACKAGE-URL>"

  //Optional field to skip cluster 
  skip_cluster_creation = true
}

output "batch" {
  value = module.batch1
}