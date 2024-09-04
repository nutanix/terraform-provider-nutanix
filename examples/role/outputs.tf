output "k8s_infra_provision_role_id" {
  value = data.nutanix_role.kubernetes_infrastructure_provision.id
}

output "k8s_data_services_system_role_id" {
  value = data.nutanix_role.k8s_data_services_system.id
}

output "csi_system_role_id" {
  value = data.nutanix_role.csi_system.id
}