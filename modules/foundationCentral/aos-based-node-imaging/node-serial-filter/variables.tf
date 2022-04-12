terraform {
  experiments = [module_variable_optional_attrs]
}

//[Optional]: Cluster external IP
variable "cluster_external_ip" {
  description = "External management ip of the cluster"
  type        = string
  default     = ""
}

//[Required]: Redundancy factor of the cluster.
variable "redundancy_factor" {
  type    = number
  default = 2
}

//[Required]: URL to download AOS package. Required only if imaging is needed.
variable "aos_package_url" {
  description = "aos package file"
  type        = string
}


//[Optional]: Number of storage only nodes in the cluster. AHV iso for storage node will be taken from aos package.
variable "storage_node_count" {
  description = "Number of storage only nodes in the cluster. AHV iso for storage node will be taken from aos package"
  type        = number
  default     = null
}

//[Required]: Name of the cluster.
variable "cluster_name" {
  description = "Name of the cluster"
  type        = string
}

//[Optional] : Number of nodes in the cluster. 
variable "cluster_size" {
  description = "Number of nodes in the cluster"
  type        = number
  default     = 0
}

//[Optional]: 
variable "skip_cluster_creation" {
  description = "Skip cluster creation"
  type        = bool
  default     = false
}
//[Optional]: Sha256sum of AOS package.
variable "aos_package_sha256sum" {
  description = "Sha256sum of AOS package"
  type        = string
  default     = ""
}

//[Optional]: Timezone to be set on the cluster.
variable "timezone" {
  description = "Timezone to be set on the cluster"
  type        = string
  default     = ""
}

//[Required]: Common network settings across the nodes in the cluster.
variable "common_network_settings" {
  description = "Common network settings across the nodes in the cluster"
  type = object({
    cvm_dns_servers        = list(string)
    hypervisor_dns_servers = list(string)
    cvm_ntp_servers        = list(string)
    hypervisor_ntp_servers = list(string)
  })
}

//[Optional]: Details of the hypervisor iso.
variable "hypervisor_iso_details" {
  description = "Details of the hypervisor iso"
  type = object({
    hyperv_sku         = optional(string)
    url                = string
    hyperv_product_key = optional(string)
    sha256sum          = optional(string)
  })
  default = null
}

// List of details of nodes out of which the cluster needs to be created.
variable "node_info" {
  description = "List of details of nodes out of which the cluster needs to be created."
  type = map(object({
    cvm_gateway                   = optional(string)
    cvm_ip                        = optional(string)
    cvm_netmask                   = optional(string)
    cvm_ram_gb                    = optional(number)
    cvm_vlan_id                   = optional(number)
    hypervisor_gateway            = optional(string)
    hypervisor_hostname           = optional(string)
    hypervisor_ip                 = optional(string)
    hypervisor_netmask            = optional(string)
    hypervisor_type               = optional(string)
    image_now                     = optional(bool)
    imaged_node_uuid              = optional(string)
    ipmi_gateway                  = optional(string)
    ipmi_ip                       = optional(string)
    ipmi_netmask                  = optional(string)
    rdma_passthrough              = optional(bool)
    use_existing_network_settings = optional(string)
    node_serial                   = optional(string)
    hardware_attributes_override  = optional(map(any))
  }))
}