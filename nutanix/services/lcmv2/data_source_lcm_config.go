package lcmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lcmconfigimport1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/resources"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixLcmConfigV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixLcmConfigV2Create,
		Schema: map[string]*schema.Schema{
			"x_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_auto_inventory_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"auto_inventory_schedule": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connectivity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_https_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"supported_software_entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Description: "List of String",
					Type:        schema.TypeString,
				},
			},
			"deprecated_software_entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Description: "List of String",
					Type:        schema.TypeString,
				},
			},
			"is_framework_bundle_uploaded": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_module_auto_upgrade_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixLcmConfigV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterExtID := d.Get("x_cluster_id").(string)
	var clusterID *string
	if clusterExtID != "" {
		clusterID = utils.StringPtr(clusterExtID)
	} else {
		clusterID = nil
	}
	resp, err := conn.LcmConfigAPIInstance.GetConfig(clusterID)
	if err != nil {
		return diag.Errorf("error while fetching the Lcm config : %v", err)
	}

	lcmConfig := resp.Data.GetValue().(lcmconfigimport1.Config)
	if err := d.Set("tenant_id", lcmConfig.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(lcmConfig.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("url", lcmConfig.Url); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_auto_inventory_enabled", lcmConfig.IsAutoInventoryEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("auto_inventory_schedule", lcmConfig.AutoInventorySchedule); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", lcmConfig.Version); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_version", lcmConfig.DisplayVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("connectivity_type", lcmConfig.ConnectivityType.GetName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_https_enabled", lcmConfig.IsHttpsEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("supported_software_entities", lcmConfig.SupportedSoftwareEntities); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deprecated_software_entities", lcmConfig.DeprecatedSoftwareEntities); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_framework_bundle_uploaded", lcmConfig.IsFrameworkBundleUploaded); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_module_auto_upgrade_enabled", lcmConfig.HasModuleAutoUpgradeEnabled); err != nil {
		return diag.FromErr(err)
	}

	// Set Id to random UUID
	d.SetId(utils.GenUUID())
	return nil
}
