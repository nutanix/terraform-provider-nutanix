package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixAddressGroupsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixAddressGroupsV2Read,
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
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_addresses": SchemaForValuePrefixLength(),
						"ip_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"policy_references": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixAddressGroupsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	// initialize query params
	var filter, orderBy, selects *string
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
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	resp, err := conn.AddressGroupAPIInstance.ListAddressGroups(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching address groups : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("address_groups", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of address groups.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.AddressGroup)
	if err := d.Set("address_groups", flattenAddressGroupsEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenAddressGroupsEntities(pr []import1.AddressGroup) []interface{} {
	if len(pr) > 0 {
		addGroups := make([]interface{}, len(pr))

		for k, v := range pr {
			add := make(map[string]interface{})

			add["ext_id"] = v.ExtId
			add["name"] = v.Name
			add["description"] = v.Description

			if v.Ipv4Addresses != nil {
				add["ipv4_addresses"] = flattenIPv4AddressMicroSeg(v.Ipv4Addresses)
			}
			if v.IpRanges != nil {
				add["ip_ranges"] = flattenIPv4Range(v.IpRanges)
			}
			if v.PolicyReferences != nil {
				add["policy_references"] = flattenListofString(v.PolicyReferences)
			}
			if v.CreatedBy != nil {
				add["created_by"] = v.CreatedBy
			}
			if v.Links != nil {
				add["links"] = flattenLinksMicroSeg(v.Links)
			}
			if v.TenantId != nil {
				add["tenant_id"] = v.TenantId
			}

			addGroups[k] = add
		}
		return addGroups
	}
	return nil
}
