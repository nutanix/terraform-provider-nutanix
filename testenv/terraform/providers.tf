
terraform {
  required_providers {
    nutanix = {
      source  = "nutanixtemp/nutanix"
      version = "1.99.99"
      # source  = "nutanix/nutanix"
      # version = "2.4.3"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.9"
    }
    # local = {
    #   source  = "hashicorp/local"
    #   version = "~> 2.0"
    # }
    # http = {
    #   source  = "hashicorp/http"
    #   version = "~> 3.0"
    # }
    # null = {
    #   source  = "hashicorp/null"
    #   version = "~> 3.1"
    # }
  }
}

provider "nutanix" {
  username = local.pc.username
  password = local.pc.password
  endpoint = local.pc.endpoint
  insecure = true
  port     = local.pc.port
}

# Bare Prism Element that the prismv2 deploy-PC test deploys a fresh PC onto.
# It is a different cluster than the local PC above, so it needs its own
# provider. Used only to pre-create the deploy network (networking.tf:
# nutanix_subnet_v2.deploy-pc-external-subnet).
provider "nutanix" {
  alias    = "deploy_pe"
  username = local.pc.username
  password = local.pc.password
  endpoint = local.deploy_pc.pe_ip
  insecure = true
  port     = local.pc.port
}
