package nutanix

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNutanixVolumeGroup() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getVGSchema(),
	}
}

func getVGSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"creation_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_hash": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"categories": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"project_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"owner_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"api_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"availability_zone_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"cluster_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"host_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"hypervisor_type": {
			Type:     schema.TypeString,
			Computed: true,
		},

		// RESOURCES ARGUMENTS

		"num_vnuma_nodes": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"nic_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nic_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"floating_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"model": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"network_function_nic_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"mac_address": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ip_endpoint_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"type": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"network_function_chain_reference": {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Required: true,
								},
								"name": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"subnet_reference": {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Required: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"subnet_reference_name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"guest_os_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"power_state": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"nutanix_guest_tools": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"available_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"iso_mount_state": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"state": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"guest_os_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"enabled_capability_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"vss_snapshot_capable": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"is_reachable": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"vm_mobility_drivers_installed": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
		"num_vcpus_per_socket": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"num_sockets": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"gpu_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"frame_buffer_size_mib": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vendor": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"pci_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"fraction": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"mode": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"num_virtual_display_heads": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"guest_driver_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"device_id": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"parent_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"memory_size_mib": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"boot_device_order_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"boot_device_disk_address": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_index": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"adapter_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"boot_device_mac_address": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"hardware_clock_timezone": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"guest_customization_cloud_init": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"meta_data": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"user_data": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"custom_key_values": {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"guest_customization_is_overridable": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"guest_customization_sysprep": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"install_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"unattend_xml": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"custom_key_values": {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"should_fail_on_script_failure": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"enable_script_exec": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"power_state_mechanism": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vga_console_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"disk_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"disk_size_bytes": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"disk_size_mib": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"device_properties": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"device_type": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"disk_address": {
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"device_index": {
												Type:     schema.TypeInt,
												Required: true,
											},
											"adapter_type": {
												Type:     schema.TypeString,
												Required: true,
											},
										},
									},
								},
							},
						},
					},
					"data_source_reference": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Required: true,
								},
								"name": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},

					"volume_group_reference": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Required: true,
								},
								"name": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
	}
}
