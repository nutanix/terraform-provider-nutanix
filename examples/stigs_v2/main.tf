terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.4.0"
    }
  }
}



data "nutanix_stigs_v2" "all" {}

data "nutanix_stigs_v2" "filtered-status" {
  filter = "status eq Security.Report.StigStatus'APPLICABLE'"
}

data "nutanix_stigs_v2" "filtered-severity" {
  filter = "severity eq Security.Report.Severity'HIGH'"
}

data "nutanix_stigs_v2" "limited" {
  limit  = 2
  select = "stigVersion,status"
}
