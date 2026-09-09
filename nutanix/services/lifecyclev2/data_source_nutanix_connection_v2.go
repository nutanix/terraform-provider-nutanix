package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixConnectionV2 returns details of a hardware provider connection
// identified by its external ID.
func DatasourceNutanixConnectionV2() *schema.Resource {
	dsSchema := connectionDatasourceSchema()
	dsSchema["hardware_provider_ext_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "External ID of the hardware provider",
	}
	// ext_id is the input identifier for the singular datasource, so it is Required.
	dsSchema["ext_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "External ID of the hardware provider connection",
	}
	return &schema.Resource{
		ReadContext: DatasourceNutanixConnectionV2Read,
		Schema:      dsSchema,
	}
}

func DatasourceNutanixConnectionV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	hpExtID := d.Get("hardware_provider_ext_id").(string)
	extID := d.Get("ext_id").(string)

	resp, err := conn.HardwareProvidersAPIInstance.GetConnectionById(utils.StringPtr(hpExtID), utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching connection : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Connection)

	for key, value := range flattenConnection(&getResp) {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(extID)
	return nil
}
