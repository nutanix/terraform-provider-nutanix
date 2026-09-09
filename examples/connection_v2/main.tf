terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
    }
  }
}

# defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Look up the available hardware providers on the cluster. Hardware providers are
# not created through Terraform; a connection is always scoped to one.
data "nutanix_hardware_providers_v2" "all" {}

data "nutanix_hardware_provider_v2" "example" {
  ext_id = var.hardware_provider_ext_id
}

# --- Connection resource: url_endpoint + basic_auth variant ---
resource "nutanix_connection_v2" "url_basic" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  name                     = "connection-url-basic"
  region                   = "us-west"

  access_details {
    auth {
      basic_auth {
        username = var.connection_username
        password = var.connection_password
      }
    }
    endpoint {
      url_endpoint {
        url = "https://hardware-provider.example.com"
      }
    }
  }
}

# --- Connection resource: ip_range_endpoint variant ---
resource "nutanix_connection_v2" "ip_range" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  name                     = "connection-ip-range"

  access_details {
    auth {
      basic_auth {
        username = var.connection_username
        password = var.connection_password
      }
    }
    endpoint {
      ip_range_endpoint {
        ip_ranges {
          start_ip {
            ipv4 {
              value = "10.0.0.1"
            }
          }
          end_ip {
            ipv4 {
              value = "10.0.0.10"
            }
          }
        }
      }
    }
  }
}

# --- Connection resource: ip_address_endpoint + api_key_auth variant ---
resource "nutanix_connection_v2" "ip_address" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  name                     = "connection-ip-address"

  access_details {
    auth {
      api_key_auth {
        api_key_id     = "api-key-id"
        api_key_secret = "api-key-secret"
      }
    }
    endpoint {
      ip_address_endpoint {
        ip_addresses {
          ipv4 {
            value = "10.0.0.5"
          }
        }
      }
    }
  }
}

# --- Action resource: create a connection via the standalone create-action surface ---
resource "nutanix_create_connection_by_hardware_provider_id_v2" "example" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  name                     = "connection-action-create"
  region                   = "us-east"

  access_details {
    auth {
      basic_auth {
        username = var.connection_username
        password = var.connection_password
      }
    }
    endpoint {
      url_endpoint {
        url = "https://hardware-provider.example.com"
      }
    }
  }
}

# --- Action resource: refresh nodes / resource pools of a connection ---
resource "nutanix_refresh_connection_v2" "example" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  ext_id                   = nutanix_connection_v2.url_basic.ext_id

  refresh_resources_spec {
    should_refresh_ip_pools              = true
    should_refresh_mac_pools             = true
    should_refresh_server_identity_pools = true
  }
}

# --- Datasources: singular + plural connection ---
data "nutanix_connection_v2" "get" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  ext_id                   = nutanix_connection_v2.url_basic.ext_id
}

data "nutanix_connections_v2" "list" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  depends_on               = [nutanix_connection_v2.url_basic]
}

# --- Datasources: nodes discovered through a connection ---
data "nutanix_connection_node_v2" "get_node" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  connection_ext_id        = nutanix_connection_v2.url_basic.ext_id
  ext_id                   = "00000000-0000-0000-0000-000000000000"
}

data "nutanix_nodes_v2" "list_nodes" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  connection_ext_id        = nutanix_connection_v2.url_basic.ext_id
}

# --- Datasources: IP / MAC / server-identity pools ---
data "nutanix_ip_pool_v2" "get_ip_pool" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  connection_ext_id        = nutanix_connection_v2.url_basic.ext_id
  ext_id                   = "00000000-0000-0000-0000-000000000000"
}

data "nutanix_ip_pools_v2" "list_ip_pools" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  connection_ext_id        = nutanix_connection_v2.url_basic.ext_id
}

data "nutanix_mac_pool_v2" "get_mac_pool" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  connection_ext_id        = nutanix_connection_v2.url_basic.ext_id
  ext_id                   = "00000000-0000-0000-0000-000000000000"
}

data "nutanix_mac_pools_v2" "list_mac_pools" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  connection_ext_id        = nutanix_connection_v2.url_basic.ext_id
}

data "nutanix_server_identity_pool_v2" "get_sip" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  connection_ext_id        = nutanix_connection_v2.url_basic.ext_id
  ext_id                   = "00000000-0000-0000-0000-000000000000"
}

data "nutanix_server_identity_pools_v2" "list_sips" {
  hardware_provider_ext_id = var.hardware_provider_ext_id
  connection_ext_id        = nutanix_connection_v2.url_basic.ext_id
}
