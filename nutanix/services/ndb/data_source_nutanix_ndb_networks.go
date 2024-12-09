package ndb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixEraNetworks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraNetworksRead,
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"managed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stretched_vlan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"secure": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"properties_map": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vlan_subnet_mask": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vlan_primary_dns": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vlan_secondary_dns": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vlan_gateway": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dbserver_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dbserver_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ip_pools": {
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
									"addresses": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"status": {
													Type:     schema.TypeString,
													Computed: true,
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

func dataSourceNutanixEraNetworksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	resp, err := conn.Service.ListNetwork(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if e := d.Set("networks", flattenNetworkListResponse(resp)); e != nil {
		return diag.FromErr(e)
	}

	uuid, er := uuid.GenerateUUID()

	if er != nil {
		return diag.Errorf("Error generating UUID for era networks: %+v", er)
	}
	d.SetId(uuid)
	return nil
}

func flattenNetworkListResponse(ntw *era.ListNetworkResponse) []interface{} {
	if ntw != nil {
		networkList := make([]interface{}, 0)

		for _, v := range *ntw {
			val := map[string]interface{}{}
			val["name"] = v.Name
			val["id"] = v.ID
			val["type"] = v.Type
			val["cluster_id"] = v.ClusterID
			val["properties"] = flattenNetworkProperties(v.Properties)
			if v.PropertiesMap != nil {
				val["properties_map"] = flattenPropertiesMap(v.PropertiesMap)
			}
			if v.Managed != nil {
				val["managed"] = v.Managed
			}
			if v.StretchedVlanID != nil {
				val["stretched_vlan_id"] = v.StretchedVlanID
			}
			if v.IPPools != nil {
				val["ip_pools"] = flattenIPPools(v.IPPools)
			}
			if v.IPAddresses != nil {
				val["ip_addresses"] = flattenIPAddress(v.IPAddresses)
			}

			networkList = append(networkList, val)
		}
		return networkList
	}
	return nil
}

func flattenNetworkProperties(erp []*era.Properties) []map[string]interface{} {
	if len(erp) > 0 {
		res := make([]map[string]interface{}, len(erp))

		for k, v := range erp {
			ents := make(map[string]interface{})
			ents["name"] = v.Name
			ents["value"] = v.Value
			ents["secure"] = v.Secure
			res[k] = ents
		}
		return res
	}
	return nil
}
