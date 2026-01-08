package vmmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixOvaDownloadV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixOvaDownloadV2Create,
		ReadContext:   ResourceNutanixOvaDownloadV2Read,
		UpdateContext: ResourceNutanixOvaDownloadV2Update,
		DeleteContext: ResourceNutanixOvaDownloadV2Delete,
		Schema: map[string]*schema.Schema{
			"ova_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ova_file_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixOvaDownloadV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	ovaExtID := d.Get("ova_ext_id")
	resp, err := conn.OvasAPIInstance.GetFileByOvaId(utils.StringPtr(ovaExtID.(string)))
	if err != nil {
		return diag.Errorf("error Downloading Ova file: %v", err)
	}

	aJSON, _ := json.MarshalIndent(resp, "", "  ")
	log.Printf("[DEBUG] Downloaded OVA file response: %s", aJSON)

	respData := resp.Data
	if respData == nil {
		return diag.Errorf("error Downloading Ova file: %v", resp)
	}

	filePath := respData.GetValue().(content.FileDetail).Path

	log.Printf("[DEBUG] OVA file path: %s", utils.StringValue(filePath))

	if err := d.Set("ova_file_path", utils.StringValue(filePath)); err != nil {
		return diag.Errorf("error setting ova_file_path: %v", err)
	}

	// This is an action resource that does not maintain state and has no associated task.
	// The resource ID is set to the OVA ext_id for traceability.
	d.SetId(ovaExtID.(string))
	return ResourceNutanixOvaDownloadV2Read(ctx, d, meta)
}

func ResourceNutanixOvaDownloadV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixOvaDownloadV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixOvaDownloadV2Create(ctx, d, meta)
}

func ResourceNutanixOvaDownloadV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
