/*
Description:
- Here we will image 3 nodes (can be called superman19 nodes), where the ipmi related info
  will be pulled from a hashicorp vault running on a external machine.
- In vault, the creds are saved under nodes/superman19 path as k-v defined as:
    "ipmi_info" : {
        "node_serial_1" : {
            "ipmi_ip" : "<ipmi-ip>",
            "ipmi_user" : "<ipmi-user>",
            "ipmi_password" : "<ipmi-pass>",
        },
        "node_serial_2" : {
            "ipmi_ip" : "<ipmi-ip>",
            "ipmi_user" : "<ipmi-user>",
            "ipmi_password" : "<ipmi-pass>",
        },
        .
        .
        .
        .
    }
*/

/*
[IMPORTANT]
Please note this from hashicorp/vault documentation:
Interacting with Vault from Terraform causes any secrets that you read and write to be persisted in both Terraform's state file and in any generated plan files. 
For any Terraform module that reads or writes Vault secrets, these files should be treated as sensitive and protected accordingly.
Docs: https://registry.terraform.io/providers/hashicorp/vault/latest/docs
*/

// pull hashicorp vault and nutanix provider
terraform {
    required_providers{
        vault = {
            source = "hashicorp/vault"
            version = "3.5.0"
        }
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.5.0-beta"
        }
    }
}

// initialize vault. This internally uses VAULT_ADDR & VAULT_TOKEN environment variables for authentication
provider "vault" {
    address = "{{ address }}"
}

// initialize nutanix provider
provider "nutanix" {
    foundation_endpoint = "10.xx.xx.xx"
}

// pull nos packages info
data "nutanix_foundation_nos_packages" "nos" {}

// get ipmi secrets from vault from secret path = nodes/superman19, where superman19 secrets 
// is having ipmi info (creds & ip) of group of nodes
data "vault_generic_secret" "superman19_ipmi_info" {
    path = "nodes/superman19"
}

// since the data source vault_generic_secret schema defines the data as sensitive, do conversion and decodig first
locals {
    ipmi_info = jsondecode(nonsensitive(data.vault_generic_secret.superman19_ipmi_info.data["ipmi_info"]))
}

// image nodes with using IPMI way where the ipmi info comes from vault with key as node-serial of nodes
resource "nutanix_foundation_image_nodes" "batch1" {
    timeouts {
        create = "70m"
    }

    nos_package = data.nutanix_foundation_nos_packages.nos.entities[0]
    cvm_netmask = "xx.xx.xx.xx"
    cvm_gateway = "10.xx.xx.xx"
    hypervisor_gateway = "10.xx.xx.xx"
    hypervisor_netmask = "xx.xx.xx.xx"
    ipmi_gateway = "10.xx.xx.xx"
    ipmi_netmask = "xx.xx.xx.xx"

    blocks{
        nodes{
            hypervisor_hostname="superman19-1"
            hypervisor_ip= "10.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            cvm_ip= "10.xx.xx.xx"
            node_position= "A"
            ipmi_ip = local.ipmi_info["<node-serial-1>"]["ipmi_ip"]
            ipmi_user = local.ipmi_info["<node-serial-1>"]["ipmi_user"]
            ipmi_password = local.ipmi_info["<node-serial-1>"]["ipmi_password"]
        }
        nodes{
            hypervisor_hostname="superman19-2"
            hypervisor_ip= "10.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            cvm_ip= "10.xx.xx.xx"
            node_position= "B"
            ipmi_ip = local.ipmi_info["<node-serial-2>"]["ipmi_ip"]
            ipmi_user = local.ipmi_info["<node-serial-2>"]["ipmi_user"]
            ipmi_password = local.ipmi_info["<node-serial-2>"]["ipmi_password"]
        }
        nodes{
            hypervisor_hostname="superman19-3"
            hypervisor_ip= "10.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            cvm_ip= "10.xx.xx.xx"
            node_position= "C"
            ipmi_ip = local.ipmi_info["<node-serial-3>"]["ipmi_ip"]
            ipmi_user = local.ipmi_info["<node-serial-3>"]["ipmi_user"]
            ipmi_password = local.ipmi_info["<node-serial-3>"]["ipmi_password"]
        }
    }

    clusters {
        redundancy_factor = 2
        cluster_name = "superman19"
        cluster_init_now = true
        cluster_members = ["10.xx.xx.xx","10.xx.xx.xx","10.xx.xx.xx"]
    }
}

// output the imaging session details
output "session" {
    value = resource.nutanix_foundation_image_nodes.batch1
}