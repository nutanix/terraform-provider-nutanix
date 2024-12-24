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

# initiate the template
resource "nutanix_template_guest_os_actions_v2" "example-initiate" {
  ext_id     = resource.nutanix_template_v2.example-1.id
  action     = "initiate"
  version_id = resource.nutanix_template_v2.example-1.template_version_spec.0.ext_id
  depends_on = [nutanix_template_v2.example-1]
}

resource "nutanix_template_guest_os_actions_v2" "example-cancel" {
  ext_id     = resource.nutanix_template_v2.example-1.id
  action     = "cancel"
  depends_on = [nutanix_template_guest_os_actions_v2.example-initiate]
}
