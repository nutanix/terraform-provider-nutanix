package nutanix

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixNDBNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBNetworkCreate,
		ReadContext:   resourceNutanixNDBNetworkRead,
		UpdateContext: resourceNutanixNDBNetworkUpdate,
		DeleteContext: resourceNutanixNDBNetworkDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_pools": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"end_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"vlan_gateway": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan_subnet_mask": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan_primary_dns": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan_secondary_dns": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan_dns_doamin": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// computed

			"managed": {
				Type:     schema.TypeBool,
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

func resourceNutanixNDBNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	req := &era.NetworkIntentInput{}

	if name, ok := d.GetOk("name"); ok {
		req.Name = utils.StringPtr(name.(string))
	}
	if clsID, ok := d.GetOk("cluster_id"); ok {
		req.ClusterID = utils.StringPtr(clsID.(string))
	}
	if netType, ok := d.GetOk("type"); ok {
		req.Type = utils.StringPtr(netType.(string))
	}
	if ipPools, ok := d.GetOk("ip_pools"); ok {
		ipPoolList := ipPools.([]interface{})

		poolList := make([]*era.IPPools, 0)
		for _, v := range ipPoolList {
			pool := &era.IPPools{}
			val := v.(map[string]interface{})

			if start, ok := val["start_ip"]; ok {
				pool.StartIP = utils.StringPtr(start.(string))
			}
			if end, ok := val["end_ip"]; ok {
				pool.EndIP = utils.StringPtr(end.(string))
			}
			poolList = append(poolList, pool)
		}
		req.IPPools = poolList
	}

	props := make([]*era.Properties, 0)
	if vlanGateway, ok := d.GetOk("vlan_gateway"); ok {
		props = append(props, &era.Properties{
			Name:  utils.StringPtr("VLAN_GATEWAY"),
			Value: utils.StringPtr(vlanGateway.(string)),
		})
	}

	if vlanSubnetMask, ok := d.GetOk("vlan_subnet_mask"); ok {
		props = append(props, &era.Properties{
			Name:  utils.StringPtr("VLAN_SUBNET_MASK"),
			Value: utils.StringPtr(vlanSubnetMask.(string)),
		})
	}

	if vlanPrimaryDns, ok := d.GetOk("vlan_primary_dns"); ok {
		props = append(props, &era.Properties{
			Name:  utils.StringPtr("VLAN_PRIMARY_DNS"),
			Value: utils.StringPtr(vlanPrimaryDns.(string)),
		})
	}

	if vlanSecDns, ok := d.GetOk("vlan_secondary_dns"); ok {
		props = append(props, &era.Properties{
			Name:  utils.StringPtr("VLAN_SECONDARY_DNS"),
			Value: utils.StringPtr(vlanSecDns.(string)),
		})
	}

	if vlanDnsDomain, ok := d.GetOk("vlan_dns_domain"); ok {
		props = append(props, &era.Properties{
			Name:  utils.StringPtr("VLAN_DNS_DOMAIN"),
			Value: utils.StringPtr(vlanDnsDomain.(string)),
		})
	}

	req.Properties = props
	// api to create network in ndb
	resp, err := conn.Service.CreateNetwork(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.ID)
	return resourceNutanixNDBNetworkRead(ctx, d, meta)
}

func resourceNutanixNDBNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	resp, err := conn.Service.GetNetwork(ctx, d.Id(), "")
	if err != nil {
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
	return nil
}
func resourceNutanixNDBNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNutanixNDBNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	resp, err := conn.Service.DeleteNetwork(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == utils.StringPtr("vLAN Successfully Removed.") {
		log.Printf("NDB Network with %s id is deleted successfully", d.Id())
		d.SetId("")
	}
	return nil
}

func flattenPropertiesMap(pm *era.NetworkPropertiesmap) []interface{} {
	if pm != nil {
		propMap := []interface{}{}
		prop := map[string]interface{}{}

		prop["vlan_gateway"] = pm.VLANGateway
		prop["vlan_primary_dns"] = pm.VLANPrimaryDNS
		prop["vlan_secondary_dns"] = pm.VLANSecondaryDNS
		prop["vlan_subnet_mask"] = pm.VLANSubnetMask

		propMap = append(propMap, prop)
		return propMap
	}
	return nil
}
