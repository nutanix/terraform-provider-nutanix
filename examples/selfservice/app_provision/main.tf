terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.2.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  insecure = true
  port     = var.nutanix_port
}

// Example 1: Provision an application (launch blueprint)
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = var.blueprint_name
    app_name        = var.app_name
    app_description = var.app_description
}


// Example 2: Provision with runtime editable
// Read runtime editable fields in a blueprint
data "nutanix_blueprint_runtime_editables" "example" {
    bp_name = var.blueprint_name
}

// dumps fetched runtime editables into a readable json file
resource "local_file" "dump_runtime_value" {
    content  = jsonencode(data.nutanix_blueprint_runtime_editables.example.runtime_editables)
    filename = var.file_name
}

// In above dumped file. Inside substrate_list look for name of substrate whose value you want to change.
// Copy complete value from "value" key and use jq to format string.
// Use command (in terminal) to format string: echo '<copied-value-with-string-quotes>' | jq -r | jq
// Copy the formatted string and paste it in your variables.tf in place of <jq-formatted-value> (example present in variables.tf)

// Launch blueprint and provision your application
resource "nutanix_self_service_app_provision" "test" {
   bp_name         = var.blueprint_name
   app_name        = var.app_name
   app_description = var.app_description

   runtime_editables {
    substrate_list {
       name = var.substrate_name
       value = var.substrate_value
     }
   }
}

// Example 3: Run system action

resource "nutanix_self_service_app_provision" "test" {
    app_uuid        = nutanix_self_service_app_provision.test.id
    action = var.system_action_name
}

# Alternatively you can also run system action by using app name
resource "nutanix_self_service_app_provision" "test" {
    app_name        = var.app_name
    action = var.system_action_name
}

// Example 4: Soft delete Application

# Step 1: Provision application 
# Step 3: set soft_delete attribute as true and run terraform apply
# Step 4: Run terraform destroy to soft delete application.

resource "nutanix_self_service_app_provision" "test" {
    app_uuid        = nutanix_self_service_app_provision.test.id
    soft_delete     = true
}
