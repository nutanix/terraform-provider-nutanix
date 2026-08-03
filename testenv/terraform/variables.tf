# No input variables: every input is sourced from testenv/config.yaml via the
# yamldecode locals in main.tf (single source of truth). The former
# terraform.tfvars variables (nutanix_*, external_nat_subnet, overlay_subnet,
# integrationVM, object_store) now live under config.yaml's pc / networking /
# vmm / object_store / terraform blocks.
