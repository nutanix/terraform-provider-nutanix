terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.4.0"
    }
  }
}

# Defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}
# Example: Creating a new cluster profile resource with all possible attributes

resource "nutanix_cluster_profile_v2" "example" {
  name              = "tf-cluster-profile"
  description       = "Example Cluster Profile created via Terraform"
  allowed_overrides = ["NTP_SERVER_CONFIG", "SNMP_SERVER_CONFIG"]

  name_server_ip_list {
    ipv4 {
      value = "240.29.254.180"
    }
    ipv6 {
      value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"
    }
  }
  ntp_server_ip_list {
    ipv4 {
      value = "240.29.254.180"
    }
    ipv6 {
      value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"
    }
    fqdn {
      value = "ntp.example.com"
    }
  }
  smtp_server {
    email_address = "email@example.com"
    type          = "SSL"
    server {
      ip_address {
        ipv4 {
          value = "240.29.254.180"
        }
        ipv6 {
          value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"
        }
        fqdn {
          value = "smtp.example.com"
        }
      }
    }
  }
  nfs_subnet_white_list = ["10.110.106.45/255.255.255.255"]
  snmp_config {
    is_enabled = true
    users {
      username  = "snmpuser1"
      auth_type = "MD5"
      auth_key  = "Example_SNMP_user_authentication_key"
      priv_type = "DES"
      priv_key  = "Example_SNMP_user_encryption_key"
    }
    transports {
      protocol = "UDP"
      port     = 21
    }
    traps {
      address {
        ipv4 {
          value         = "240.29.254.180"
          prefix_length = 24
        }
        ipv6 {
          value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"
        }
      }
      username         = "trapuser"
      protocol         = "UDP"
      port             = 59
      should_inform    = false
      engine_id        = "0x1234567890abcdef12"
      version          = "V2"
      receiver_name    = "trap-receiver"
      community_string = "snmp-server community public RO 192.168.1.0 255.255.255.0"
    }
  }
  rsyslog_server_list {
    server_name      = "exampleServer1"
    port             = 29
    network_protocol = "UDP"
    ip_address {
      ipv4 {
        value = "240.29.254.180"
      }
      ipv6 {
        value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"
      }
    }
    modules {
      name                     = "CASSANDRA"
      log_severity_level       = "EMERGENCY"
      should_log_monitor_files = true
    }
    modules {
      name                     = "CURATOR"
      log_severity_level       = "ERROR"
      should_log_monitor_files = false
    }
  }
  pulse_status {
    is_enabled          = false
    pii_scrubbing_level = "DEFAULT"
  }
}

output "cluster_profile_all_attrs" {
  value = nutanix_cluster_profile_v2.example
}


# List all cluster profiles
data "nutanix_cluster_profiles_v2" "list-cluster-profiles" {
}

# Filter cluster profiles by cluster count
data "nutanix_cluster_profiles_v2" "filtered-cluster-profiles" {
  filter = "clusterCount eq 62"
}

# Get paginated cluster profiles
data "nutanix_cluster_profiles_v2" "paged-cluster-profiles" {
  page  = 1
  limit = 10
}

# Get ordered cluster profiles by name
data "nutanix_cluster_profiles_v2" "ordered-cluster-profiles" {
  order_by = "name"
}

# Get cluster profiles with selected fields only
data "nutanix_cluster_profiles_v2" "selected-cluster-profiles" {
  select = "name,description"
}

# Filter cluster profiles by name
data "nutanix_cluster_profiles_v2" "filtered-by-name" {
  filter = "name eq 'Test Cluster Profile'"
}

# Filter cluster profiles by create time
data "nutanix_cluster_profiles_v2" "filtered-by-create-time" {
  filter = "createTime eq '2009-09-23T14:30:00-07:00'"
}

# Get cluster profiles ordered by create time descending
data "nutanix_cluster_profiles_v2" "ordered-by-create-time" {
  order_by = "createTime desc"
}

# Output examples
output "all_cluster_profiles" {
  value = data.nutanix_cluster_profiles_v2.list-cluster-profiles.cluster_profiles
}

output "filtered_cluster_profiles" {
  value = data.nutanix_cluster_profiles_v2.filtered-cluster-profiles.cluster_profiles
}

output "paged_cluster_profiles" {
  value = data.nutanix_cluster_profiles_v2.paged-cluster-profiles.cluster_profiles
}

