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


# Create a new Template from done only using vm reference
resource "nutanix_template_v2" "example-1" {
  template_name        = "tf-example-template"
  template_description = "test create template from vm using terraform"
  template_version_spec {
    version_source {
      template_vm_reference {
        ext_id = "<VM_UUID>"
      }
    }
  }
}

# Create a new Template from an existing VM. with Guest Customization
resource "nutanix_template_v2" "example-2" {
  template_name        = "tf-example-template"
  template_description = "test create template from vm using terraform"
  template_version_spec {
    version_source {
      template_vm_reference {
        ext_id = "<VM_UUID>"
      }
    }
  }
  guest_customization {
    config {
      sysprep {
        sysprep_script {
          custom_key_values {
            key_value_pairs {
              name = "locale"
              value {
                string = "en-PS"
              }
            }
          }
        }
      }
    }
  }
}

# for updating the existing template, we can use template_version_reference or template_vm_reference, only one of them
# version name and version description are mandatory fields on update operation
# to update template and override the existing configuration, we will use template_version_reference
resource "nutanix_template_v2" "example-1" {
  template_name        = "tf-example-template"
  template_description = "test create template from vm using terraform"
  template_version_spec {
    version_name        = "2.0.0"
    version_description = "updating version from initial to 2.0.0"
    is_active_version   = true
    version_source {
      template_vm_reference {
        ext_id = "<VM_UUID>"
      }
      template_version_reference {
        # if version id is not provided, it will use the latest version of the template by default
        version_id = "<TEMPLATE_VERSION_UUID>"
        override_vm_config {
          name                 = "tf-test-vm-2.0.0"
          memory_size_bytes    = 3 * 1024 * 1024 * 1024 # 3 GB
          num_cores_per_socket = 2
          num_sockets          = 2
          num_threads_per_core = 2
          guest_customization {
            config {
              cloud_init {
                cloud_init_script {
                  user_data {
                    value = base64encode("#cloud-config\nusers:\n  - name: ubuntu\n    ssh-authorized-keys:\n      - ssh-rsa DUMMYSSH mypass\n    sudo: ['ALL=(ALL) NOPASSWD:ALL']")
                  }
                  custom_key_values {
                    key_value_pairs {
                      name = "locale"
                      value {
                        string = "en-PS"
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}


# to update template and override the existing configuration, we will use template_vm_reference

resource "nutanix_template_v2" "example-1" {
  template_name        = "tf-example-template"
  template_description = "test create template from vm using terraform"
  template_version_spec {
    version_name        = "2.0.0"
    version_description = "updating version from initial to 2.0.0"
    is_active_version   = true
    version_source {
      template_vm_reference {
        ext_id = "<New_VM_UUID>"
      }
    }
  }
}


# List all the Templates in the system.
data "nutanix_templates_v2" "templates-1" {}

# List all the Templates in the system with a filter.
data "nutanix_templates_v2" "templates-2" {
  filter = "templateName eq '${nutanix_template_v2.example-1.template_name}'"
}

# List all the Templates in the system with a limit.
data "nutanix_templates_v2" "templates-3" {
  limit = 3
}

# Get a specific Template by UUID.
data "nutanix_template_v2" "template-1" {
  ext_id = nutanix_template_v2.example-1.ext_id
  # or -> ext_id = nutanix_template_v2.example-1.id
}
