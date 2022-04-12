terraform {
  experiments = [module_variable_optional_attrs]
}

//[Required]: URL to download AOS package. Required only if imaging is needed.
variable "aos_package_url" {
  description = "aos package file"
  type        = string
}

//[Required]: Node details to image and create cluster
variable "node_list" {
  description = "List of details of nodes out of which the cluster needs to be created."
  type = list(object({
    cvm_gateway                   = string
    ipmi_netmask                  = string
    rdma_passthrough              = optional(bool)
    imaged_node_uuid              = string
    cvm_vlan_id                   = optional(string)
    hypervisor_type               = string
    image_now                     = bool
    hypervisor_hostname           = string
    hypervisor_netmask            = string
    cvm_netmask                   = string
    ipmi_ip                       = string
    hypervisor_gateway            = string
    cvm_ram_gb                    = optional(number)
    cvm_ip                        = string
    hypervisor_ip                 = string
    use_existing_network_settings = string
    ipmi_gateway                  = string
  }))
}


//[Optional]
variable "cluster_external_ip" {
  description = "External management ip of the cluster."
  type        = string
  default     = ""
}


//[Required]: Redundancy factor of the cluster.
variable "redundancy_factor" {
  type    = number
  default = 2
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

//[Optional]: skip cluster creation
variable "skip_cluster_creation" {
  description = "skip cluster creation"
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
    hyperv_sku         = string
    url                = optional(string)
    hyperv_product_key = string
    sha256sum          = string
  })
  default = null
}