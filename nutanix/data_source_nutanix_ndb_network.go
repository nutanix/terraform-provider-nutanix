package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNutanixEraNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraNetworkRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}
func dataSourceNutanixEraNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	name, nok := d.GetOk("name")
	networkID, iok := d.GetOk("id")

	if !nok && !iok {
		return diag.Errorf("either name or id is required to get the network details")
	}

	resp, err := conn.Service.GetNetwork(ctx, name.(string), networkID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", resp.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", resp.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("managed", resp.Managed); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_id", resp.ClusterID); err != nil {
		return diag.FromErr(err)
	}
	props := []interface{}{}
	for _, prop := range resp.Properties {
		props = append(props, map[string]interface{}{
			"name":   prop.Name,
			"value":  prop.Value,
			"secure": prop.Secure,
		})
	}
	if err := d.Set("properties", props); err != nil {
		return diag.FromErr(err)
	}

	if resp.PropertiesMap != nil {
		d.Set("properties_map", flattenPropertiesMap(resp.PropertiesMap))
	}

	if resp.StretchedVlanID != nil {
		d.Set("stretched_vlan_id", resp.StretchedVlanID)
	}

	d.SetId(*resp.ID)
	return nil
}
