package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixHardwareProviderV2 returns details of a hardware provider
// identified by its external ID.
func DatasourceNutanixHardwareProviderV2() *schema.Resource {
	dsSchema := hardwareProviderDatasourceSchema()
	dsSchema["ext_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "External ID of the hardware provider",
	}
	return &schema.Resource{
		ReadContext: DatasourceNutanixHardwareProviderV2Read,
		Schema:      dsSchema,
	}
}

func DatasourceNutanixHardwareProviderV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.HardwareProvidersAPIInstance.GetHardwareProviderById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching hardware provider : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.HardwareProvider)

	for key, value := range flattenHardwareProvider(&getResp) {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(extID)
	return nil
}
