package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixSecretV2 reads the secret value of a claim token identified by
// its external ID.
func DatasourceNutanixSecretV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixSecretV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "External ID of the claim token",
			},
			"secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Claim token secret data",
			},
		},
	}
}

func DatasourceNutanixSecretV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.ClaimTokensAPIInstance.GetSecretByClaimTokenId(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching claim token secret : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.ClaimTokenSecret)

	if err := d.Set("secret", getResp.Secret); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(extID)
	return nil
}
