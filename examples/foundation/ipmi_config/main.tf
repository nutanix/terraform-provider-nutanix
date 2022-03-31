// nutanix_foundation_ipmi_config resource was introduced in nutanix/nutanix version >=1.4.2
terraform{
    required_providers{
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.4.2"
        }
    }
}

// default foundation_port is 8000 so can be ignored
provider "nutanix" {
    // foundation_port = 8000
    foundation_endpoint = "10.xx.xx.xx"
}

/*
Description:
Here we are configuring ipmi in two nodes of two blocks each with
username = abc & password = xyz
*/
resource "nutanix_foundation_ipmi_config" "impi1" {
  ipmi_user = "abc"
  ipmi_password = "xyz"
  ipmi_gateway = "xx.xx.xx.xx"
  ipmi_netmask = "xx.xx.xx.xx"
  blocks{
    nodes {
          ipmi_mac = "xx:xx:xx:xx:xx:xx"
          ipmi_configure_now =  true
          ipmi_ip = "10.xx.xx.xx"
    }
    nodes {
          ipmi_mac = "xx:xx:xx:xx:xx:xx"
          ipmi_configure_now =  true
          ipmi_ip = "10.xx.xx.xx"
    }
    block_id = "xxxx"
  }
  blocks{
    nodes {
          ipmi_mac = "xx:xx:xx:xx:xx:xx"
          ipmi_configure_now =  true
          ipmi_ip = "10.xx.xx.xx"
    }
    nodes {
          ipmi_mac = "xx:xx:xx:xx:xx:xx"
          ipmi_configure_now =  true
          ipmi_ip = "10.xx.xx.xx"
    }
    block_id = "xxxx"
  }
}