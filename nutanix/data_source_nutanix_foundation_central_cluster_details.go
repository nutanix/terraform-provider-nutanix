package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNutanixFCClusterDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixFCClusterDetailsRead,
		Schema: map[string]*schema.Schema{
			"imaged_cluster_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"current_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"archived": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cluster_external_ip": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"imaged_node_uuid_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"common_network_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cvm_dns_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"hypervisor_dns_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cvm_ntp_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"hypervisor_ntp_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"storage_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"redundancy_factor": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"foundation_init_node_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workflow_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"foundation_init_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blocks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"block_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nodes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cvm_ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"fc_imaged_node_uuid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hypervisor": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hypervisor_hostname": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hypervisor_ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"image_now": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"ipmi_ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ipv6_address": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"node_position": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"node_serial": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"clusters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_external_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_init_now": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"cluster_init_successful": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"cluster_members": {
										Type:     schema.TypeString,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"cluster_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cvm_dns_servers": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cvm_ntp_servers": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"redundancy_factor": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"timezone": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cvm_gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cvm_netmask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_servers": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hyperv_product_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hyperv_sku": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_iso_url": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hypervisor_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sha256sum": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"hypervisor_isos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hypervisor_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sha256sum": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"hypervisor_netmask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipmi_gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipmi_netmask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nos_package_url": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sha256sum": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"cluster_status": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_creation_started": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"intent_picked_up": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"imaging_stopped": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"node_progress_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"imaged_node_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"imaging_stopped": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"intent_picked_up": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"percent_complete": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"message_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"aggregate_percent_complete": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"current_foundation_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_progress_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"percent_complete": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"message_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"foundation_session_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cluster_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"destroyed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixFCClusterDetailsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
