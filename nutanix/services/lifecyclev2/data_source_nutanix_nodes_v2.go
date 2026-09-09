package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixNodesV2 reads a list of nodes registered with a specific
// claim token. It requires the parent claim token id.
func DatasourceNutanixNodesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixNodesV2Read,
		Schema: map[string]*schema.Schema{
			"claim_token_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "External ID of the claim token",
			},
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: datasourceClaimTokenNodeSchema(),
				},
			},
		},
	}
}

func DatasourceNutanixNodesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	claimTokenExtID := d.Get("claim_token_ext_id").(string)

	var page, limit *int
	var filter, orderBy, selectQ *string

	if v, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		selectQ = utils.StringPtr(v.(string))
	}

	resp, err := conn.ClaimTokensAPIInstance.ListNodesByClaimTokenId(utils.StringPtr(claimTokenExtID), page, limit, filter, orderBy, selectQ)
	if err != nil {
		return diag.Errorf("error while listing nodes : %v", err)
	}

	d.SetId(resource.UniqueId())

	if resp.Data == nil {
		if err := d.Set("nodes", make([]map[string]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of nodes.",
		}}
	}

	nodes := resp.Data.GetValue().([]import1.ClaimTokenNode)

	if err := d.Set("nodes", flattenClaimTokenNodes(nodes)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
