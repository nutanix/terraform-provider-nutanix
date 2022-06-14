terraform{
  required_providers{
    nutanix = {
      source = "nutanix/nutanix"
      version = "1.6.0"
    }
  }
}
provider "nutanix" {
  username  = "admin"
  password  = "Nutanix/123456"
  endpoint  = "10.xx.xx.xx"
  insecure  = true
  port      = 9440
}


#pull all clusters data
data "nutanix_clusters" "clusters"{}

#create local variable pointing to desired cluster
locals {
	cluster1 = [
	  for cluster in data.nutanix_clusters.clusters.entities :
	  cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

#creating subnet
resource "nutanix_subnet" "acc" {
    name        = "vpc-sub"
    description = "Description of my unit test VLAN updated"
    subnet_type = "VLAN"
    cluster_uuid = local.cluster1
    vlan_id = 121
    subnet_ip          = "10.xx.xx.xx"
    default_gateway_ip = "10.xx.xx.xx"
    prefix_length = 24

    is_external = true
    enable_nat = false
    ip_config_pool_list_ranges = ["10.xx.xx.xx 10.xx.xx.xx"]
            
}

// creating VPC

resource "nutanix_vpc" "test1" {
    name = "testtNew-1"

    // Ext Subnet Reference
    external_subnet_reference_uuid = [
        resource.nutanix_subnet.acc.id
    ]
    
    common_domain_name_server_ip_list{
            ip = "x.x.x.x"
    }
    
    externally_routable_prefix_list{
        ip=  "192.xx.x.xx"
        prefix_length= 24
    }
    
}


// dataSources for VPC


//dataSource to get details for an entity

data "nutanix_vpc" "vpc"{
    // vpc uuid required to get VPC entity
    vpc_uuid = ""
}

output "vpcOut1" {
   value =  data.nutanix_vpc.vpc
}


// dataSource to all vpc present

data "nutanix_vpc_list" "vpclist"{
    // Optional paramters are length and offset
    length = 10
}

output "vpcOut2" {
   value =  data.nutanix_vpc_list.vpclist
}