package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamAuthn "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

// DatasourceNutanixWelcomeBannerV2 fetches the configured welcome banner.
func DatasourceNutanixWelcomeBannerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixWelcomeBannerV2Read,
		Schema: map[string]*schema.Schema{
			"content": {
				Description: "Content of the welcome banner.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_time": {
				Description: "Creation time of the welcome banner.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_enabled": {
				Description: "Flag to denote whether the welcome banner is enabled or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"last_updated_time": {
				Description: "Last updated time of the welcome banner.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func DatasourceNutanixWelcomeBannerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	resp, err := conn.WelcomeBannerAPIInstance.GetWelcomeBanner()
	if err != nil {
		return diag.Errorf("error while fetching welcome banner: %v", err)
	}

	if resp == nil || resp.GetData() == nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty welcome banner.",
		}}
	}

	getResp := resp.GetData().(iamAuthn.WelcomeBanner)

	if err := d.Set("content", getResp.Content); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("is_enabled", getResp.IsEnabled); err != nil {
		return diag.FromErr(err)
	}
	if getResp.LastUpdatedTime != nil {
		if err := d.Set("last_updated_time", getResp.LastUpdatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(resource.UniqueId())
	return nil
}
