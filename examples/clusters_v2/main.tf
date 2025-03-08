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


# Add Cluster with minimum configuration.

resource "nutanix_cluster_v2" "example" {
  name = "terraform-example-cluster"
  nodes {
    node_list {
      controller_vm_ip {
        ipv4 {
          value = var.cvm_ip
        }
      }
    }
  }
  config {
    cluster_function  = ["AOS"]
    redundancy_factor = 1
    cluster_arch      = "X86_64"
    fault_tolerance_state {
      domain_awareness_level = "DISK"
    }
  }
  # after create a cluster you need to reset the pe ui password
  provisioner "local-exec" {
    command = "sshpass -p '${var.pe_password}' ssh ${var.pe_username}@${var.cvm_ip} '/home/nutanix/prism/cli/ncli user reset-password user-name=${var.username} password=${var.password}'"
  }
}

# after create a cluster you need to register the cluster with prism central
# to be able to do read, update, delete operations on the cluster and use it
resource "nutanix_pc_registration_v2" "pc1" {
  pc_ext_id = var.pc_ext_id
  remote_cluster {
    aos_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = var.cvm_ip
          }
        }
        credentials {
          authentication {
            username = var.username
            password = var.password
          }
        }
      }
    }
  }
}

# List all  Directory Services.
data "nutanix_clusters_v2" "example" {
  depends_on = [nutanix_pc_registration_v2.pc1]
}

# List Clusters with filter.
data "nutanix_clusters_v2" "example" {
  filter     = "name eq '${nutanix_cluster_v2.example.name}'"
  depends_on = [nutanix_pc_registration_v2.pc1]
}


# Get a Directory Service.
data "nutanix_cluster_v2" "example" {
  ext_id     = nutanix_cluster_v2.example.id
  depends_on = [nutanix_pc_registration_v2.pc1]
}
