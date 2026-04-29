variable "nutanix_endpoint" {}
variable "nutanix_username" {}
variable "nutanix_password" {}
variable "nutanix_port" {}
variable "nutanix_insecure" {}

provider "nutanix" {
  endpoint = var.nutanix_endpoint
  username = var.nutanix_username
  password = var.nutanix_password
  port     = var.nutanix_port
  insecure = var.nutanix_insecure
}

# List AOS clusters
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

# Run System-Defined Checks on the cluster
resource "nutanix_run_system_defined_checks_v2" "example" {
  cluster_ext_id        = local.clusterExtID
  should_run_all_checks = true
  should_anonymize      = false
  should_send_report_to_configured_recipients = true
}
