package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixNodeV2 reads details of a single node registered with a
// specific claim token. It requires the parent claim token id and the node id.
func DatasourceNutanixNodeV2() *schema.Resource {
	nodeSchema := datasourceClaimTokenNodeSchema()
	nodeSchema["claim_token_ext_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "External ID of the claim token",
	}
	// ext_id is a required input for the singular node datasource.
	nodeSchema["ext_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "External ID of the node",
	}
	return &schema.Resource{
		ReadContext: DatasourceNutanixNodeV2Read,
		Schema:      nodeSchema,
	}
}

func DatasourceNutanixNodeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	claimTokenExtID := d.Get("claim_token_ext_id").(string)
	extID := d.Get("ext_id").(string)

	resp, err := conn.ClaimTokensAPIInstance.GetNodeByClaimTokenId(utils.StringPtr(claimTokenExtID), utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching node : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.ClaimTokenNode)

	node := flattenClaimTokenNode(&getResp)
	for k, v := range node {
		if k == "ext_id" {
			continue
		}
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
