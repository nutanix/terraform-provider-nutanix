terraform{
    required_providers {
        nutanix = {
            source  = "nutanixtemp/nutanix"
            version = "1.99.99"
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

# create PBR with vpc uuid with any source or destination or protocol with permit action

resource "nutanix_pbr" "pbr1" {
    name = "test1"
    priority = 123
 
 
    vpc_reference_uuid = var.vpc_reference_uuid
    source{
        address_type = "ALL"
    }
    destination{
        address_type = "ALL"
    }

    protocol_type = "ALL"
   action = "PERMIT"
}


# create PBR with vpc uuid with source external and destination network with reroute action and  tcp port rangelist

resource "nutanix_pbr" "pbr2" {
    name = "test2"
    priority = 132
 
 
    vpc_reference_uuid = var.vpc_reference_uuid
    source{
        address_type = "INTERNET"
    }
    destination{
        subnet_ip=  "1.2.2.0"
        prefix_length=  24
    }

    protocol_type = "TCP"
    protocol_parameters{
        tcp{
            source_port_range_list{
                end_port  = 50
                start_port = 50
            }
            destination_port_range_list{
                end_port  = 40
                start_port = 40
            }
        }
    }

    action = "REROUTE"
    service_ip_list = ["10.2.2.34"]
}

#create PBR with vpc uuid with source network and destination external with reroute action and  udp port rangelist

resource "nutanix_pbr" "pbr3" {
    name = "test3"
    priority = 212

    vpc_reference_uuid = var.vpc_reference_uuid
    source{
        subnet_ip=  "1.2.2.0"
        prefix_length=  24
    }
    destination{
        address_type = "INTERNET"
    }

    protocol_type = "UDP"
    protocol_parameters{
        udp{
            source_port_range_list{
                end_port  = 50
                start_port = 50
            }
            destination_port_range_list{
                end_port  = 40
                start_port = 40
            }
        }
    }

    action = "REROUTE"
    service_ip_list = ["10.2.2.34"]
}

#create PBR with vpc name with source external and destination network with reroute action and icmp

resource "nutanix_pbr" "pbr4" {
    name = "test4"
    priority = 222

    vpc_reference_uuid = var.vpc_reference_uuid
    source{
        address_type = "INTERNET"
    }
    destination{
        subnet_ip=  "1.2.2.0"
        prefix_length=  24
    }

    protocol_type = "ICMP"
    protocol_parameters{
        icmp {
        icmp_type = 2
        icmp_code = 20
        }
    }

    action = "REROUTE"
    service_ip_list = ["10.2.2.34"]  
}

#create PBR with vpc uuid with source Any and destination external and deny action with protocol number

resource "nutanix_pbr" "pbr5" {
    name = "test5"
    priority = 252

    vpc_reference_uuid = var.vpc_reference_uuid
    source{
        address_type = "ALL"
    }
    destination{
        address_type = "INTERNET"
    }

    protocol_type = "PROTOCOL_NUMBER"
    protocol_parameters{
        protocol_number = "2"
    }
    action = "DENY"
}