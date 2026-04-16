



# This Terraform script will do:
# 1. Deploy an object store with one worker node
# 2. Create Certificate for an object store
# 3. List all certificates for an object store
# 4. Fetch certificate details for an object store


# NOTE:
# 1. Before Deleting object store, make sure to delete buckets inside it
#    Currently, we are not supporting delete bucket API in terraform
# 2. Object store Update is used only to resume deployment of object store when it fails,
#    the state will be OBJECT_STORE_DEPLOYMENT_FAILED, update will resume the deployment


terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.3.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

# subnet name to be used for object store
locals {
  subnetName = "objects.800"
}

# Fetching cluster and subnet details
data "nutanix_clusters_v2" "clusters" {}

data "nutanix_subnets_v2" "subnets" {
  filter = "name eq '${local.subnetName}'"
}

# Fetching cluster and subnet ext_id
locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
  subnetExtId = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
}

# Deploying an object store
resource "nutanix_object_store_v2" "example" {
  name                     = "tf-example-os"
  description              = "terraform create object store example"
  deployment_version       = "5.1.1"
  domain                   = "msp.pc-idbc.nutanix.com"
  num_worker_nodes         = 1
  cluster_ext_id           = local.clusterExtId
  total_capacity_gib       = 20 * pow(1024, 3) # 20 GB
  public_network_reference = local.subnetExtId
  public_network_ips {
    ipv4 {
      value = "10.44.77.123"
    }
  }
  storage_network_reference = local.subnetExtId
  storage_network_dns_ip {
    ipv4 {
      value = "10.44.77.124"
    }
  }
  storage_network_vip {
    ipv4 {
      value = "10.44.77.125"
    }
  }
}

# This is example of creating certificate for object store
# check API Ref for more details
# create object_store_cert.json file, file content :]
# {
#   "alternateIps": [
#     {
#       "ipv4": {
#         "value": "10.44.77.123"
#       }
#     }
#   ]
# }

# Creating certificate for object store
resource "nutanix_object_store_certificate_v2" "example" {
  object_store_ext_id = nutanix_object_store_v2.example.id
  # path to certificate json file
  path = "./object_store_cert.json"
}


#List all certificates for object store
data "nutanix_certificates_v2" "list" {
  object_store_ext_id = nutanix_object_store_v2.example.id
  depends_on          = [nutanix_object_store_certificate_v2.example]
}

# fetching certificate details for object store
data "nutanix_certificate_v2" "fetch" {
  object_store_ext_id = nutanix_object_store_v2.example.id
  ext_id              = nutanix_object_store_certificate_v2.example.id
  depends_on          = [nutanix_object_store_certificate_v2.example]
}
