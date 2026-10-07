terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# ---------------------------------------------------------------------------
# SNMP V2 trap: identified by community_string (no SNMP user required).
# ---------------------------------------------------------------------------
resource "nutanix_snmp_trap_v2" "v2" {
  cluster_ext_id   = var.cluster_ext_id
  version          = "V2"
  community_string = "public"
  port             = var.snmp_trap_v2_port
  protocol         = "UDP"

  address {
    ipv4 {
      value = var.snmp_trap_v2_ipv4
    }
  }
}

# ---------------------------------------------------------------------------
# SNMP V3 trap: must reference an existing SNMP user on the same cluster
# by username. Below we provision the user inline so the example is
# self-contained; reference an existing user instead by setting `username`
# to a literal string and removing the dependent resource.
# ---------------------------------------------------------------------------
resource "nutanix_snmp_user_v2" "v3_user" {
  cluster_ext_id = var.cluster_ext_id
  username       = var.snmp_v3_username
  auth_type      = "SHA"
  auth_key       = "auth-key-12345678"
  priv_type      = "AES"
  priv_key       = "priv-key-12345678"
}

resource "nutanix_snmp_trap_v2" "v3" {
  cluster_ext_id = var.cluster_ext_id
  version        = "V3"
  username       = nutanix_snmp_user_v2.v3_user.username
  port           = var.snmp_trap_v3_port
  protocol       = "UDP"

  address {
    ipv4 {
      value = var.snmp_trap_v3_ipv4
    }
  }
}

# Read a single trap back via the singular datasource (uses the V2 trap
# created above; swap to nutanix_snmp_trap_v2.v3.ext_id for the V3 trap).
data "nutanix_snmp_trap_v2" "by_id" {
  cluster_ext_id = var.cluster_ext_id
  ext_id         = nutanix_snmp_trap_v2.v2.ext_id
}

output "snmp_trap_v2" {
  value = data.nutanix_snmp_trap_v2.by_id
}
