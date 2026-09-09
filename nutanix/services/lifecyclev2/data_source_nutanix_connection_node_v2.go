package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixConnectionNodeV2 returns details of a node discovered through a
// specific hardware provider connection.
func DatasourceNutanixConnectionNodeV2() *schema.Resource {
	dsSchema := discoveredNodeDatasourceSchema()
	dsSchema["hardware_provider_ext_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "External ID of the hardware provider",
	}
	dsSchema["connection_ext_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "External ID of the hardware provider connection",
	}
	dsSchema["ext_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "External ID of the discovered node",
	}
	return &schema.Resource{
		ReadContext: DatasourceNutanixConnectionNodeV2Read,
		Schema:      dsSchema,
	}
}

func DatasourceNutanixConnectionNodeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	hpExtID := d.Get("hardware_provider_ext_id").(string)
	connExtID := d.Get("connection_ext_id").(string)
	extID := d.Get("ext_id").(string)

	resp, err := conn.HardwareProvidersAPIInstance.GetConnectionNodeById(utils.StringPtr(hpExtID), utils.StringPtr(connExtID), utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching connection node : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.DiscoveredNode)

	for key, value := range flattenDiscoveredNode(&getResp) {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(extID)
	return nil
}
