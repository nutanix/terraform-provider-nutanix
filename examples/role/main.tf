provider "nutanix" {
  username     = var.user
  password     = var.password
  endpoint     = var.endpoint
  insecure     = var.insecure
  port         = var.port
  wait_timeout = 60
}

# Create Kubernetes Infrastructure Provision role
# ---------------
data "nutanix_permission" "k8s_infra_provision_permissions" {
  for_each        = toset(var.k8s_infra_provision_permissions)
  permission_name = each.key
}

resource "nutanix_role" "kubernetes_infrastructure_provision" {
  name        = "Kubernetes Infrastructure Provision"
  description = "Access for Kubernetes cluster infrastructure VMs resources"
  dynamic "permission_reference_list" {
    for_each = data.nutanix_permission.k8s_infra_provision_permissions
    content {
      kind = "permission"
      uuid = permission_reference_list.value.id
    }
  }
}

data "nutanix_role" "kubernetes_infrastructure_provision" {
  role_id = nutanix_role.kubernetes_infrastructure_provision.id
}

# Create CSI System role
# ---------------
data "nutanix_permission" "csi_system_role_permissions" {
  for_each        = toset(var.csi_system_role_permissions)
  permission_name = each.key
}

resource "nutanix_role" "csi_system" {
  name        = "CSI System"
  description = "Full access for Kubernetes cluster infrastructure resources for CSI"
  dynamic "permission_reference_list" {
    for_each = data.nutanix_permission.csi_system_role_permissions
    content {
      kind = "permission"
      uuid = permission_reference_list.value.id
    }
  }
}

data "nutanix_role" "csi_system" {
  role_id = nutanix_role.csi_system.id
}

# Create Kubernetes Data Services System role
# ---------------
data "nutanix_permission" "k8s_data_services_system_role_permissions" {
  for_each        = toset(var.k8s_data_services_system_role_permissions)
  permission_name = each.key
}

resource "nutanix_role" "k8s_data_services_system" {
  name        = "Kubernetes Data Services System"
  description = "Full access for Kubernetes cluster infrastructure resources for Kubernetes Data Services"
  dynamic "permission_reference_list" {
    for_each = data.nutanix_permission.k8s_data_services_system_role_permissions
    content {
      kind = "permission"
      uuid = permission_reference_list.value.id
    }
  }
}

data "nutanix_role" "k8s_data_services_system" {
  role_id = nutanix_role.k8s_data_services_system.id
}
