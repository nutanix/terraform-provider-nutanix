package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	networkingConfig "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	networkingVsReq "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/request/virtualswitches"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVirtualSwitchV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixVirtualSwitchV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of Virtual Switch",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User-visible Virtual Switch name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Input body to configure a Virtual Switch",
			},
			"bond_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cluster configuration list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_ip_address": {
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
						"hosts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_nics": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"internal_bridge_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_address": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
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
												"prefix_length": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"active_uplink": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"route_table": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"vlan_identifier": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"mtu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"igmp_spec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_snooping_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"snooping_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"querier_spec": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_querier_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"vlan_id_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
					},
				},
			},
			"is_quick_mode": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_deployment_error": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_update_in_progress": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_delete_in_progress": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"owner_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared_with_projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DatasourceMetadataSchemaV2(),
				},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceNutanixVirtualSwitchV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	extID := d.Get("ext_id").(string)
	getReq := networkingVsReq.GetVirtualSwitchByIdRequest{
		ExtId: utils.StringPtr(extID),
	}

	resp, err := conn.VirtualSwitchAPI.GetVirtualSwitchById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error while fetching Virtual Switch: %v", err)
	}

	vs := resp.Data.GetValue().(networkingConfig.VirtualSwitch)

	d.SetId(utils.StringValue(vs.ExtId))

	if err := d.Set("name", vs.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", vs.Description); err != nil {
		return diag.FromErr(err)
	}
	if vs.BondMode != nil {
		if err := d.Set("bond_mode", vs.BondMode.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("clusters", flattenVirtualSwitchClusters(vs.Clusters)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mtu", vs.Mtu); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("igmp_spec", flattenIgmpSpec(vs.IgmpSpec)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_quick_mode", vs.IsQuickMode); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_default", vs.IsDefault); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_deployment_error", vs.HasDeploymentError); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_update_in_progress", vs.HasUpdateInProgress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_delete_in_progress", vs.HasDeleteInProgress); err != nil {
		return diag.FromErr(err)
	}
	if vs.OwnerType != nil {
		if err := d.Set("owner_type", vs.OwnerType.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("project_ext_id", vs.ProjectExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("shared_with_projects", vs.SharedWithProjects); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(vs.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(vs.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", vs.TenantId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
