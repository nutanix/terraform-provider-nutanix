// resources/datasources used in this file were introduced in nutanix/nutanix version 1.5.0-beta.2
terraform {
    required_providers {
      nutanix = {
          source = "nutanix/nutanix"
          version = ">=1.5.0-beta.2"
      }
    }
}

provider "nutanix" {
    username  = "user"
    password  = "pass"
    endpoint  = "10.x.xx.xx"
    insecure  = true
    port      = 9440
}


// resource to image and create cluster

resource "nutanix_foundation_central_image_cluster" "res"{
  cluster_name = "test-FC"
  cluster_external_ip = "<CLUSTER-IP>"
  common_network_settings{
    cvm_dns_servers=[
        "xx.x.xx.xx"
    ]
    hypervisor_dns_servers=[
        "xx.x.xx.xx"
    ]
    cvm_ntp_servers=[
        "<cvm-ntp>"
    ]
    hypervisor_ntp_servers=[
        "<hypervisor-ntp>"
    ]
  }
    redundancy_factor = 2
    node_list{
      cvm_gateway="10.xx.xx.xx"
      cvm_netmask="xx.xx.xx.xx"
      cvm_ip="10.x.xx.xx"
      hypervisor_gateway="10.x.x.xx"
      hypervisor_netmask="xx.xx.xx.xx"
      hypervisor_ip="10.x.xx.xx"
      hypervisor_hostname="HOST-1"
      imaged_node_uuid="<NODE-UUID>"
      use_existing_network_settings=false
      ipmi_gateway="10.x.xx.xx"
      ipmi_netmask="10.x.xx.xx"
      ipmi_ip="10.x.xx.xx"
      image_now=true
      hypervisor_type="kvm"
      hardware_attributes_override = {
        default_workload="vdi"
        lcm_family= "smc_gen_10"
        maybe_1GbE_only= true
        robo_mixed_hypervisor= true
      }
    }
    node_list{
        cvm_gateway="10.xx.xx.xx"
        cvm_netmask="xx.xx.xx.xx"
        cvm_ip="10.x.xx.xx"
        hypervisor_gateway="10.x.x.xx"
        hypervisor_netmask="xx.xx.xx.xx"
        hypervisor_ip="10.x.xx.xx"
        hypervisor_hostname="HOST-2"
        imaged_node_uuid="<NODE-UUID>"
        use_existing_network_settings=false
        ipmi_gateway="10.x.xx.xx"
        ipmi_netmask="10.x.xx.xx"
        ipmi_ip="10.x.xx.xx"
        image_now=true
        hypervisor_type="kvm"
    }
    node_list{
        cvm_gateway="10.xx.xx.xx"
        cvm_netmask="xx.xx.xx.xx"
        cvm_ip="10.x.xx.xx"
        hypervisor_gateway="10.x.x.xx"
        hypervisor_netmask="xx.xx.xx.xx"
        hypervisor_ip="10.x.xx.xx"
        hypervisor_hostname="HOST-3"
        imaged_node_uuid="<NODE-UUID>"
        use_existing_network_settings=false
        ipmi_gateway="10.x.xx.xx"
        ipmi_netmask="10.x.xx.xx"
        ipmi_ip="10.x.xx.xx"
        image_now=true
        hypervisor_type="kvm"
    }
    aos_package_url="<URL>"

    //pass true to skip cluster creation
    skip_cluster_creation = true


}

output "res1"{
    value = resource.nutanix_foundation_central_image_cluster.res
}


//datasource to List all the clusters created using Foundation Central.

data "nutanix_foundation_central_imaged_clusters_list" "cls" {}

output "cls1" {
  value = data.nutanix_foundation_central_imaged_clusters_list.cls
}

// datasource to Get a cluster created using Foundation Central.
data "nutanix_foundation_central_cluster_details" "clsDet" {
  imaged_cluster_uuid = "<imaged_cluster_uuid>"
}

output "clsDet1" {
  value = data.nutanix_foundation_central_cluster_details.clsDet
}