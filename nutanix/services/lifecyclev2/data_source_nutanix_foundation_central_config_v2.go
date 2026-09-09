package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/mgmt"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

// DatasourceNutanixFoundationCentralConfigV2 reads the current configuration settings
// of the Foundation Central service. Foundation Central exposes a single, cluster-wide
// configuration object, so this datasource takes no id input and returns one object.
func DatasourceNutanixFoundationCentralConfigV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixFoundationCentralConfigV2Read,
		Schema: map[string]*schema.Schema{
			"ahv_installation_timeout_minutes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timeout in minutes for AHV installation",
			},
			"aos_download_timeout_minutes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timeout in minutes for AOS download",
			},
			"commit_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Foundation Central Commit ID",
			},
			"tls_certificate_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA-256 fingerprint of the TLS certificate used by Foundation Central service",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Foundation Central Version",
			},
		},
	}
}

func DatasourceNutanixFoundationCentralConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.FoundationCentralConfigAPIInstance.GetFoundationCentralConfig()
	if err != nil {
		return diag.Errorf("error while fetching foundation central config : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.FoundationCentralConfig)

	if err := d.Set("ahv_installation_timeout_minutes", getResp.AhvInstallationTimeoutMinutes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("aos_download_timeout_minutes", getResp.AosDownloadTimeoutMinutes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("commit_id", getResp.CommitId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tls_certificate_fingerprint", getResp.TlsCertificateFingerprint); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", getResp.Version); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
