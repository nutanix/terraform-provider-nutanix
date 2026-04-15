output "management_subnet_id" {
  value = nutanix_subnet_v2.management.id
}

output "network_function_id" {
  value = nutanix_network_function_v2.nf.id
}

output "network_function_nic_pairs" {
  value = data.nutanix_network_function_v2.nf.nic_pairs
}

output "nf_vm_ids" {
  value = {
    for name, details in local.nf_vm_details :
    name => details.vm_ext_id
  }
}

output "nf_vm_management_ips" {
  value = {
    for name, details in local.nf_vm_details :
    name => details.management_ips
  }
}
