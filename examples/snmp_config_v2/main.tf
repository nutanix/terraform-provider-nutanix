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
# Status mode: only `is_enabled` is set, so the resource manages the
# cluster-wide SNMP enabled flag (UpdateSnmpStatus). Omit `port` and
# `protocol` to opt into this mode.
# ---------------------------------------------------------------------------
resource "nutanix_snmp_config_v2" "status" {
  cluster_ext_id = var.cluster_ext_id
  is_enabled     = var.snmp_is_enabled
}

# ---------------------------------------------------------------------------
# Transport mode: `port` and `protocol` are set, so the resource manages a
# single SNMP transport on the cluster (AddSnmpTransport on create,
# RemoveSnmpTransport on delete). `port`/`protocol` are immutable; changing
# either forces resource replacement.
# ---------------------------------------------------------------------------
resource "nutanix_snmp_config_v2" "transport" {
  cluster_ext_id = var.cluster_ext_id
  port           = var.snmp_port
  protocol       = var.snmp_protocol
}

# Singular datasource: returns the full SNMP config of the cluster
# (is_enabled, transports, traps, users) in a single read.
data "nutanix_snmp_config_v2" "current" {
  cluster_ext_id = var.cluster_ext_id

  depends_on = [
    nutanix_snmp_config_v2.status,
    nutanix_snmp_config_v2.transport,
  ]
}

output "snmp_is_enabled" {
  value = data.nutanix_snmp_config_v2.current.is_enabled
}
output "snmp_transports" {
  value = data.nutanix_snmp_config_v2.current.transports
}
