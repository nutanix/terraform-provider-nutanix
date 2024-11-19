package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import7 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/vmm/v4/ahv/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVMNetworkDevicesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVMNetworkDevicesV2Read,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backing_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"model": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mac_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_connected": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"num_queues": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"network_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nic_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_function_chain": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"network_function_nic_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"vlan_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"trunked_vlans": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"should_allow_unknown_macs": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"ipv4_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"should_assign_ip": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"ip_address": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"secondary_ip_address_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									// not visible in API reference
									"ipv4_info": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"learned_ip_addresses": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Required: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Default:  32,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixVMNetworkDevicesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	// initialize query params
	var filter *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}

	resp, err := conn.VMAPIInstance.ListNicsByVmId(utils.StringPtr(vmExtID.(string)), page, limit, filter)
	if err != nil {

		return diag.Errorf("error while fetching network devices : %v", err)
	}
	getResp := resp.Data

	if getResp != nil {
		nics := getResp.GetValue().([]import7.Nic)
		if err := d.Set("network_devices", flattenNetworkDeviceEntities(nics)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenNetworkDeviceEntities(pr []import7.Nic) []interface{} {
	if len(pr) > 0 {
		vspecs := make([]interface{}, len(pr))

		for k, v := range pr {
			spec := make(map[string]interface{})

			if v.ExtId != nil {
				spec["ext_id"] = v.ExtId
			}
			if v.BackingInfo != nil {
				spec["backing_info"] = flattenEmulatedNic(v.BackingInfo)
			}
			if v.NetworkInfo != nil {
				spec["network_info"] = flattenNicNetworkInfo(v.NetworkInfo)
			}
			vspecs[k] = spec
		}
		return vspecs
	}
	return nil
}
