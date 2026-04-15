package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixFloatingIPsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixFloatingIPsV2Read,
		Schema: map[string]*schema.Schema{
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
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"floating_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"association": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_nic_association": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vm_nic_reference": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"vpc_reference": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"private_ip_association": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_reference": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"private_ip": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": SchemaForValuePrefixLength(),
															"ipv6": SchemaForValuePrefixLength(),
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"floating_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLength(),
									"ipv6": SchemaForValuePrefixLength(),
								},
							},
						},
						"external_subnet_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_subnet": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     DataSourceNutanixSubnetV2(),
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"floating_ip_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"association_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_nic_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: DataSourceVPCSchemaV2(),
							},
						},
						"vm_nic": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"private_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
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
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: DatasourceMetadataSchemaV2(),
							},
						},
					},
				},
			},
		},
	}
}

func datasourceNutanixFloatingIPsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	// initialize query params
	var filter, orderBy, expand *string
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
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if expandf, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandf.(string))
	} else {
		expand = nil
	}

	resp, err := conn.FloatingIPAPIInstance.ListFloatingIps(page, limit, filter, orderBy, expand)
	if err != nil {
		return diag.Errorf("error while fetching floating_ips : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("floating_ips", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of floating IPs.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.FloatingIp)
	if err := d.Set("floating_ips", flattenFloatingIPsEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenFloatingIPsEntities(pr []import1.FloatingIp) []map[string]interface{} {
	if len(pr) > 0 {
		fips := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			fip := make(map[string]interface{})

			fip["ext_id"] = v.ExtId
			fip["name"] = v.Name
			fip["description"] = v.Description
			fip["association"] = flattenAssociation(v.Association)
			fip["floating_ip"] = flattenFloatingIP(v.FloatingIp)
			fip["external_subnet_reference"] = v.ExternalSubnetReference
			fip["external_subnet"] = flattenExternalSubnet(v.ExternalSubnet)
			fip["private_ip"] = v.PrivateIp
			fip["floating_ip_value"] = v.FloatingIpValue
			fip["association_status"] = v.AssociationStatus.GetName()
			fip["vpc_reference"] = v.VpcReference
			fip["vm_nic_reference"] = v.VmNicReference
			fip["vpc"] = flattenVpc(v.Vpc)
			fip["vm_nic"] = flattenVMNic(v.VmNic)
			fip["links"] = flattenLinks(v.Links)
			fip["tenant_id"] = v.TenantId
			fip["metadata"] = flattenMetadata(v.Metadata)

			fips[k] = fip
		}
		return fips
	}
	return nil
}
