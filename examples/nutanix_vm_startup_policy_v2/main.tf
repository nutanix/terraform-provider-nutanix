terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.0"
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

# Create categories for VM startup policy groups
resource "nutanix_category_v2" "group1_category" {
  key   = "vm-startup-group1"
  value = "vm-startup-group1-value"
}

resource "nutanix_category_v2" "group2_category" {
  key   = "vm-startup-group2"
  value = "vm-startup-group2-value"
}

# Create a VM startup policy with PowerOn start condition
resource "nutanix_vm_startup_policy_v2" "example_power_on" {
  name        = "example-vm-startup-policy-power-on"
  description = "Example VM startup policy with power on start condition"

  groups {
    categories {
      ext_id = nutanix_category_v2.group1_category.id
    }
  }

  groups {
    categories {
      ext_id = nutanix_category_v2.group2_category.id
    }
  }

  start_conditions {
    delay_duration_secs = 30
    power_state_criteria {
      power_on {}
    }
  }
}

# Create a VM startup policy associated with a project
resource "nutanix_vm_startup_policy_v2" "example_with_project" {
  name           = "example-vm-startup-policy-with-project"
  description    = "Example VM startup policy associated with a project"
  project_ext_id = "<project-uuid>"

  groups {
    categories {
      ext_id = nutanix_category_v2.group1_category.id
    }
  }

  groups {
    categories {
      ext_id = nutanix_category_v2.group2_category.id
    }
  }

  start_conditions {
    delay_duration_secs = 30
    power_state_criteria {
      power_on {}
    }
  }
}

# Create a VM startup policy with GuestBootup start condition
resource "nutanix_vm_startup_policy_v2" "example_guest_bootup" {
  name        = "example-vm-startup-policy-guest-bootup"
  description = "Example VM startup policy with guest bootup start condition"

  groups {
    categories {
      ext_id = nutanix_category_v2.group1_category.id
    }
  }

  groups {
    categories {
      ext_id = nutanix_category_v2.group2_category.id
    }
  }

  start_conditions {
    delay_duration_secs = 60
    power_state_criteria {
      guest_bootup {
        timeout_duration_secs = 120
      }
    }
  }
}

# Singular datasource - fetch by ID
data "nutanix_vm_startup_policy_v2" "get_policy" {
  ext_id = nutanix_vm_startup_policy_v2.example_power_on.id
}

# Plural datasource - list all policies
data "nutanix_vm_startup_policies_v2" "list_policies" {
  depends_on = [nutanix_vm_startup_policy_v2.example_power_on]
}

# List VM compliance states for a policy
data "nutanix_vm_startup_policy_vm_compliance_states_v2" "compliance" {
  vm_startup_policy_ext_id = nutanix_vm_startup_policy_v2.example_power_on.id
}

# List dependency conflicts for a policy
data "nutanix_vm_startup_policy_dependency_conflicts_v2" "dep_conflicts" {
  vm_startup_policy_ext_id = nutanix_vm_startup_policy_v2.example_power_on.id
}

# List start condition conflicts for a policy
data "nutanix_vm_startup_policy_start_condition_conflicts_v2" "sc_conflicts" {
  vm_startup_policy_ext_id = nutanix_vm_startup_policy_v2.example_power_on.id
}
