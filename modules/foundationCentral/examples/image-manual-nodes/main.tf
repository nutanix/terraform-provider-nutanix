terraform {
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

// use manual-based-node-imaging 
module "batch" {
  // source where module code is present in local machine
  source       = "../../modules/foundationCentral/manual_node_imaging"
  cluster_name = "test-fc"

  // Required setting to provide network settings
  common_network_settings = {
    cvm_dns_servers : [
      "10.x.x.xx"
    ]
    hypervisor_dns_servers : [
      "10.x.x.xx"
    ]
    cvm_ntp_servers : [
      "<>"
    ]
    hypervisor_ntp_servers : [
      "<>"
    ]
  }


  /*[Required] All the required info about the nodes such as cvm_ip, hypervisor_ip, ipmi_ip, etc are mandantory to be
given . This is manual node imaging ,so no information about the nodes will be fetched by datasources . 
*/
  node_list = [{
    cvm_gateway                   = "10.x.xx.xx"
    cvm_netmask                   = "10.x.xx.xx"
    cvm_ip                        = "10.x.xx.xx"
    hypervisor_gateway            = "10.x.xx.xx"
    hypervisor_netmask            = "10.x.xx.xx"
    hypervisor_ip                 = "10.x.xx.xx"
    hypervisor_hostname           = "HOST-1"
    imaged_node_uuid              = "<NODE-UUID>"
    use_existing_network_settings = false
    ipmi_gateway                  = "10.x.xx.xx"
    ipmi_netmask                  = "10.x.xx.xx"
    ipmi_ip                       = "10.x.xx.xx"
    image_now                     = true
    hypervisor_type               = "kvm"
    },
    {
      cvm_gateway                   = "10.x.xx.xx"
      cvm_netmask                   = "10.x.xx.xx"
      cvm_ip                        = "10.x.xx.xx"
      hypervisor_gateway            = "10.x.xx.xx"
      hypervisor_netmask            = "10.x.xx.xx"
      hypervisor_ip                 = "10.x.xx.xx"
      hypervisor_hostname           = "HOST-2"
      imaged_node_uuid              = "<NODE-UUID>"
      use_existing_network_settings = false
      ipmi_gateway                  = "10.x.xx.xx"
      ipmi_netmask                  = "10.x.xx.xx"
      ipmi_ip                       = "10.x.xx.xx"
      image_now                     = true
      hypervisor_type               = "kvm"
  }]
  redundancy_factor = 2

  //[Required] provide AOS Package URL
  aos_package_url = "<AOS-PACKAGE-URL>"
}

output "k1" {
  value = module.batch
}