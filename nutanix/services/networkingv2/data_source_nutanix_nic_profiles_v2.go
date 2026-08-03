package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/request/nicprofiles"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func DataSourceNutanixNicProfilesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceNutanixNicProfilesV2Read,
		Schema: map[string]*schema.Schema{
			"nic_profiles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DataSourceNutanixNicProfileV2(),
			},
		},
	}
}

func DataSourceNutanixNicProfilesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	listRequest := import2.ListNicProfilesRequest{}
	resp, err := conn.NicProfilesAPI.ListNicProfiles(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while fetching NIC profiles: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("nic_profiles", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return nil
	}

	raw := resp.Data.GetValue()
	out := make([]interface{}, 0)
	switch v := raw.(type) {
	case []import1.NicProfile:
		out = make([]interface{}, len(v))
		for i, item := range v {
			out[i] = flattenNicProfileEntity(make(map[string]interface{}), item)
		}
	case []*import1.NicProfile:
		for _, item := range v {
			if item != nil {
				out = append(out, flattenNicProfileEntity(make(map[string]interface{}), *item))
			}
		}
	case []import1.NicProfileProjection:
		out = make([]interface{}, len(v))
		for i, item := range v {
			out[i] = flattenNicProfileProjectionEntity(make(map[string]interface{}), item)
		}
	case []*import1.NicProfileProjection:
		for _, item := range v {
			if item != nil {
				out = append(out, flattenNicProfileProjectionEntity(make(map[string]interface{}), *item))
			}
		}
	default:
		return diag.Errorf("unexpected NIC profiles response type: %T", raw)
	}

	if err := d.Set("nic_profiles", out); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
