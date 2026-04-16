package ndb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixNDBProfileAvailableIPs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBProfileAvailableIPsRead,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"available_ips": {
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
						"property_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"managed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixNDBProfileAvailableIPsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	profileID := ""
	if ID, ok := d.GetOk("profile_id"); ok {
		profileID = ID.(string)
	}
	resp, err := conn.Service.GetAvailableIPs(ctx, profileID)
	if err != nil {
		return diag.FromErr(err)
	}

	if e := d.Set("available_ips", flattenAvailableIPsResponse(resp)); err != nil {
		return diag.FromErr(e)
	}

	d.SetId(profileID)
	return nil
}

func flattenAvailableIPsResponse(ips *era.GetNetworkAvailableIPs) []map[string]interface{} {
	if ips != nil {
		lst := []map[string]interface{}{}
		for _, v := range *ips {
			d := map[string]interface{}{}

			d["id"] = v.ID
			d["name"] = v.Name
			d["property_name"] = v.PropertyName
			d["type"] = v.Type
			d["managed"] = v.Managed
			d["cluster_id"] = v.ClusterID
			d["cluster_name"] = v.ClusterName
			d["ip_addresses"] = utils.StringValueSlice(v.IPAddresses)

			lst = append(lst, d)
		}
		return lst
	}
	return nil
}
