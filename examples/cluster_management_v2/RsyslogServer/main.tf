terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Create an RSYSLOG server configuration
resource "nutanix_rsyslog_server_v2" "example" {
  cluster_ext_id   = "00000000-0000-0000-0000-000000000000"
  server_name      = "example-rsyslog-server"
  port             = 514
  network_protocol = "UDP"

  ip_address {
    ipv4 {
      value = "10.0.0.1"
    }
  }

  modules {
    name                     = "CASSANDRA"
    log_severity_level       = "INFO"
    should_log_monitor_files = true
  }

  modules {
    name                     = "STARGATE"
    log_severity_level       = "WARNING"
    should_log_monitor_files = false
  }
}

# Fetch a specific RSYSLOG server by ext_id
data "nutanix_rsyslog_server_v2" "example" {
  cluster_ext_id = nutanix_rsyslog_server_v2.example.cluster_ext_id
  ext_id         = nutanix_rsyslog_server_v2.example.ext_id
}

# List all RSYSLOG servers for a cluster
data "nutanix_rsyslog_servers_v2" "example" {
  cluster_ext_id = nutanix_rsyslog_server_v2.example.cluster_ext_id
}
