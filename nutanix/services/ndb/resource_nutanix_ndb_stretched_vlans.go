package ndb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBStretchedVlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBStretchedVlanCreate,
		ReadContext:   resourceNutanixNDBStretchedVlanRead,
		UpdateContext: resourceNutanixNDBStretchedVlanUpdate,
		DeleteContext: resourceNutanixNDBStretchedVlanDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vlan_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"subnet_mask": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			//computed field
			"vlans_list": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
				},
			},
		},
	}
}

func resourceNutanixNDBStretchedVlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.StretchedVlansInput{}

	if name, ok := d.GetOk("name"); ok {
		req.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		req.Description = utils.StringPtr(desc.(string))
	}
	if netType, ok := d.GetOk("type"); ok {
		req.Type = utils.StringPtr(netType.(string))
	}
	if vlanIDs, ok := d.GetOk("vlan_ids"); ok {
		res := make([]*string, 0)
		vlanList := vlanIDs.([]interface{})
		for _, v := range vlanList {
			res = append(res, utils.StringPtr(v.(string)))
		}

		req.VlanIDs = res
	}

	// api to stretched vlan

	resp, err := conn.Service.CreateStretchedVlan(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.ID)
	log.Printf("NDB Stretched Vlan with %s id is created successfully", d.Id())
	return resourceNutanixNDBStretchedVlanRead(ctx, d, meta)
}

func resourceNutanixNDBStretchedVlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	// check if d.Id() is nil
	if d.Id() == "" {
		return diag.Errorf("stretched vlan id is required for read operation")
	}
	resp, err := conn.Service.GetStretchedVlan(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", resp.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenStretchedVlanMetadata(resp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vlans_list", flattenStretchedVlans(resp.Vlans)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNutanixNDBStretchedVlanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	updateReq := &era.StretchedVlansInput{}
	metadata := &era.StretchedVlanMetadata{}
	// api to network api

	resp, err := conn.Service.GetStretchedVlan(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp != nil {
		updateReq.Name = resp.Name
		updateReq.Description = resp.Name
		// updateReq.Metadata = resp.Metadata
		updateReq.Type = resp.Type

		// get the vlans ids
		if resp.Vlans != nil {
			vlans := make([]*string, 0)
			for _, v := range resp.Vlans {
				vlans = append(vlans, v.ID)
			}
			updateReq.VlanIDs = vlans
		}

		if resp.Metadata != nil {
			metadata = resp.Metadata
			updateReq.Metadata = metadata
		}
	}

	if d.HasChange("name") {
		updateReq.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateReq.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("type") {
		updateReq.Type = utils.StringPtr(d.Get("type").(string))
	}
	if d.HasChange("vlan_ids") {
		res := make([]*string, 0)
		vlanList := d.Get("vlan_ids").([]interface{})
		for _, v := range vlanList {
			res = append(res, utils.StringPtr(v.(string)))
		}

		updateReq.VlanIDs = res
	}
	if d.HasChange("metadata") {
		metadataList := d.Get("metadata").([]interface{})
		for _, v := range metadataList {
			val := v.(map[string]interface{})

			if gateway, ok := val["gateway"]; ok && len(gateway.(string)) > 0 {
				metadata.Gateway = utils.StringPtr(gateway.(string))
			}
			if subnetMask, ok := val["subnet_mask"]; ok && len(subnetMask.(string)) > 0 {
				metadata.SubnetMask = utils.StringPtr(subnetMask.(string))
			}
		}
		updateReq.Metadata = metadata
	}

	updateResp, er := conn.Service.UpdateStretchedVlan(ctx, d.Id(), updateReq)
	if er != nil {
		return diag.FromErr(er)
	}

	if updateResp != nil {
		log.Printf("NDB Stretched Vlan with %s id is updated successfully", d.Id())
	}
	return resourceNutanixNDBStretchedVlanRead(ctx, d, meta)
}

func resourceNutanixNDBStretchedVlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	resp, err := conn.Service.DeleteStretchedVlan(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == utils.StringPtr("vLAN Successfully Removed.") {
		log.Printf("NDB Stretched Vlan with %s id is deleted successfully", d.Id())
		d.SetId("")
	}
	return nil
}

func flattenStretchedVlans(net []*era.NetworkIntentResponse) []interface{} {
	if len(net) > 0 {
		netList := make([]interface{}, len(net))

		for _, v := range net {
			nwt := map[string]interface{}{}
			nwt["id"] = v.ID
			nwt["name"] = v.Name
			nwt["type"] = v.Type
			nwt["cluster_id"] = v.ClusterID
			nwt["managed"] = v.Managed
			if v.Properties != nil {
				props := []interface{}{}
				for _, prop := range v.Properties {
					props = append(props, map[string]interface{}{
						"name":   prop.Name,
						"value":  prop.Value,
						"secure": prop.Secure,
					})
				}
				nwt["properties"] = props
			}
			if v.PropertiesMap != nil {
				nwt["properties_map"] = flattenPropertiesMap(v.PropertiesMap)
			}
			if v.StretchedVlanID != nil {
				nwt["stretched_vlan_id"] = v.StretchedVlanID
			}
			netList = append(netList, nwt)
		}
		return netList
	}
	return nil
}

func flattenStretchedVlanMetadata(pr *era.StretchedVlanMetadata) []interface{} {
	if pr != nil {
		metaList := make([]interface{}, 0)

		meta := map[string]interface{}{}

		meta["gateway"] = pr.Gateway
		meta["subnet_mask"] = pr.SubnetMask

		metaList = append(metaList, meta)
		return metaList
	}
	return nil
}
