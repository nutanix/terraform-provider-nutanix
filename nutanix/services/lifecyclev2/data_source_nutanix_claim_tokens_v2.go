package lifecyclev2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixClaimTokensV2 reads a paginated list of all claim tokens.
func DatasourceNutanixClaimTokensV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixClaimTokensV2Read,
		Schema: map[string]*schema.Schema{
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
			"claim_tokens": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixClaimTokenV2(),
			},
		},
	}
}

func DatasourceNutanixClaimTokensV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

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

	resp, err := conn.ClaimTokensAPIInstance.ListClaimTokens(page, limit, filter, orderBy, selectQ)
	if err != nil {
		return diag.Errorf("error while listing claim tokens : %v", err)
	}

	d.SetId(resource.UniqueId())

	if resp.Data == nil {
		if err := d.Set("claim_tokens", make([]map[string]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of claim tokens.",
		}}
	}

	tokens := resp.Data.GetValue().([]import1.ClaimToken)

	if err := d.Set("claim_tokens", flattenClaimTokens(tokens)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenClaimTokens(tokens []import1.ClaimToken) []map[string]interface{} {
	if len(tokens) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(tokens))
	for i := range tokens {
		token := tokens[i]
		m := map[string]interface{}{
			"ext_id":              utils.StringValue(token.ExtId),
			"name":                utils.StringValue(token.Name),
			"max_usage_count":     utils.IntValue(token.MaxUsageCount),
			"current_usage_count": utils.IntValue(token.CurrentUsageCount),
			"owner_ext_id":        utils.StringValue(token.OwnerExtId),
			"tenant_id":           utils.StringValue(token.TenantId),
			"links":               flattenLifecycleLinks(token.Links),
		}
		if token.ExpiryTime != nil {
			m["expiry_time"] = token.ExpiryTime.Format(time.RFC3339)
		}
		if token.CreatedTime != nil {
			m["created_time"] = token.CreatedTime.Format(time.RFC3339)
		}
		out = append(out, m)
	}
	return out
}
