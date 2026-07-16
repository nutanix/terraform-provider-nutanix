package iamv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	authn "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// welcomeBannerSingletonID is a fixed identifier used for the welcome banner
// resource. The welcome banner is a cluster-wide singleton configuration and
// the API does not expose an ext_id for it, so a constant ID is used.
const welcomeBannerSingletonID = "welcome_banner"

// ResourceNutanixWelcomeBannerV2 defines the schema and CRUD handlers for the
// welcome banner singleton resource. The welcome banner is updated via the
// UpdateWelcomeBanner API and read via the GetWelcomeBanner API.
func ResourceNutanixWelcomeBannerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixWelcomeBannerV2Create,
		ReadContext:   ResourceNutanixWelcomeBannerV2Read,
		UpdateContext: ResourceNutanixWelcomeBannerV2Update,
		DeleteContext: ResourceNutanixWelcomeBannerV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"content": {
				Description: "Content of the welcome banner.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"is_enabled": {
				Description: "Flag to denote whether the welcome banner is enabled or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"created_time": {
				Description: "Creation time of the welcome banner.",
				Type:        schema.TypeString,
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

func ResourceNutanixWelcomeBannerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	body := &authn.WelcomeBanner{}
	if content, ok := d.GetOk("content"); ok {
		body.Content = utils.StringPtr(content.(string))
	}
	if isEnabled, ok := d.GetOkExists("is_enabled"); ok {
		body.IsEnabled = utils.BoolPtr(isEnabled.(bool))
	}

	resp, err := conn.WelcomeBannerAPIInstance.UpdateWelcomeBanner(body)
	if err != nil {
		return diag.Errorf("error while creating welcome banner: %v", err)
	}
	log.Printf("[DEBUG] Welcome banner created. Response: %v", resp)

	d.SetId(welcomeBannerSingletonID)
	return ResourceNutanixWelcomeBannerV2Read(ctx, d, meta)
}

func ResourceNutanixWelcomeBannerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	resp, err := conn.WelcomeBannerAPIInstance.GetWelcomeBanner()
	if err != nil {
		return diag.Errorf("error while reading welcome banner: %v", err)
	}

	if resp.Data == nil {
		d.SetId("")
		return nil
	}

	banner := resp.Data.GetValue().(authn.WelcomeBanner)

	if err := d.Set("content", banner.Content); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_enabled", banner.IsEnabled); err != nil {
		return diag.FromErr(err)
	}
	if banner.CreatedTime != nil {
		if err := d.Set("created_time", banner.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if banner.LastUpdatedTime != nil {
		if err := d.Set("last_updated_time", banner.LastUpdatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Id() == "" {
		d.SetId(welcomeBannerSingletonID)
	}
	return nil
}

func ResourceNutanixWelcomeBannerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	body := &authn.WelcomeBanner{}
	if content, ok := d.GetOk("content"); ok {
		body.Content = utils.StringPtr(content.(string))
	}
	if isEnabled, ok := d.GetOkExists("is_enabled"); ok {
		body.IsEnabled = utils.BoolPtr(isEnabled.(bool))
	}

	resp, err := conn.WelcomeBannerAPIInstance.UpdateWelcomeBanner(body)
	if err != nil {
		return diag.Errorf("error while updating welcome banner: %v", err)
	}
	log.Printf("[DEBUG] Welcome banner updated. Response: %v", resp)

	return ResourceNutanixWelcomeBannerV2Read(ctx, d, meta)
}

func ResourceNutanixWelcomeBannerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	// The welcome banner is a singleton config that cannot be deleted. Deleting
	// the Terraform resource resets it to a disabled state with empty content.
	body := &authn.WelcomeBanner{
		Content:   utils.StringPtr(""),
		IsEnabled: utils.BoolPtr(false),
	}

	resp, err := conn.WelcomeBannerAPIInstance.UpdateWelcomeBanner(body)
	if err != nil {
		return diag.Errorf("error while deleting (resetting) welcome banner: %v", err)
	}
	log.Printf("[DEBUG] Welcome banner reset on delete. Response: %v", resp)

	d.SetId("")
	return nil
}
