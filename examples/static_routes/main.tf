terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.6.0"
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


#create one static route for vpc with external subnet
resource "nutanix_static_routes" "scn" {
  vpc_uuid = "{{vpc_uuid}}"

  static_routes_list{
    # destination prefix format 10.x.x.x/x
    destination= "10.x.x.x/x"
    # required ext subnet uuid for next hop
    external_subnet_reference_uuid = "{{ext_subnet_uuid}}" 
  }
}


#create multiple static routes for vpc with external subnet
resource "nutanix_static_routes" "scn" {
  vpc_uuid = "{{vpc_uuid}}"

  static_routes_list{
    # destination prefix format 10.x.x.x/x
    destination= "10.x.x.x/x"
    # required ext subnet uuid for next hop
    external_subnet_reference_uuid = "{{ext_subnet_uuid}}" 
  }

  static_routes_list{
    destination= "10.x.x.x/x"
    external_subnet_reference_uuid = "{{ext_subnet_uuid}}"
  }
}