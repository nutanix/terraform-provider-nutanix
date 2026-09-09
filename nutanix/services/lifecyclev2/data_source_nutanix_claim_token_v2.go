package lifecyclev2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixClaimTokenV2 reads details of a single claim token identified
// by its external ID.
func DatasourceNutanixClaimTokenV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixClaimTokenV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A globally unique identifier of an instance that is suitable for external consumption.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the claim token",
			},
			"expiry_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiry time of the claim token",
			},
			"max_usage_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of times the claim token can be used for registering nodes",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the claim token was created",
			},
			"current_usage_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of times the claim token has been used for registering nodes",
			},
			"owner_ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External ID of the owner",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
			},
			"links": datasourceLifecycleLinksSchema(),
		},
	}
}

func DatasourceNutanixClaimTokenV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.ClaimTokensAPIInstance.GetClaimTokenById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching claim token : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.ClaimToken)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if getResp.ExpiryTime != nil {
		if err := d.Set("expiry_time", getResp.ExpiryTime.Format(time.RFC3339)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("max_usage_count", getResp.MaxUsageCount); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.Format(time.RFC3339)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("current_usage_count", getResp.CurrentUsageCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLifecycleLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
