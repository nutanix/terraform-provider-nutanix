# Note - Replace appropriate values of variables in terraform.tfvars file as per setup

terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
    }
  }
}

# Defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

# Get cluster information
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

# Local variable for cluster UUID
locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

# Example 1: Update SSL certificate with provided certificate data
# Uncomment this block if you want to import/update a certificate
resource "nutanix_ssl_certificate_v2" "example_with_cert" {
  cluster_ext_id        = local.cluster_ext_id
  passphrase            = var.passphrase != "" ? var.passphrase : null
  private_key           = var.private_key != "" ? base64encode(var.private_key) : null
  public_certificate    = var.public_certificate != "" ? base64encode(var.public_certificate) : null
  ca_chain             = var.ca_chain != "" ? base64encode(var.ca_chain) : null
  private_key_algorithm = "RSA_2048"
}

# Example 2: Regenerate self-signed certificate (only RSA_2048 supported)
# Uncomment this block if you want to regenerate a self-signed certificate
resource "nutanix_ssl_certificate_v2" "example_regenerate" {
  cluster_ext_id        = local.cluster_ext_id
  private_key_algorithm = "RSA_2048"
}

# Read SSL certificate details
data "nutanix_ssl_certificate_v2" "example" {
  cluster_ext_id = local.cluster_ext_id
  depends_on     = [nutanix_ssl_certificate_v2.example_regenerate]
}

