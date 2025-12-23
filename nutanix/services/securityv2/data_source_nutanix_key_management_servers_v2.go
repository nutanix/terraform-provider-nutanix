// Package securityv2 provides resources for managing security-related configurations in Nutanix.
package securityv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/security/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixKeyManagementServersV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixKeyManagementServersV2Read,
		Schema: map[string]*schema.Schema{
			"kms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixKeyManagementServerV2(),
			},
		},
	}
}

func DatasourceNutanixKeyManagementServersV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	resp, err := conn.KeyManagementServersAPIInstance.ListKeyManagementServers()
	if err != nil {
		return diag.Errorf("error while listing key management server : %v", err)
	}

	if resp.Data == nil {
		if setErr := d.Set("kms", []map[string]interface{}{}); setErr != nil {
			return diag.FromErr(setErr)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No Key Management Servers found",
			Detail:   "The API returned an empty list of key management servers.",
		}}
	}

	listRespValue, ok := resp.Data.GetValue().([]config.KeyManagementServer)
	if !ok {
		return diag.Errorf("error: unexpected response type from list API, expected []KeyManagementServer")
	}
	listResp := listRespValue

	if len(listResp) == 0 {
		if setErr := d.Set("kms", []map[string]interface{}{}); setErr != nil {
			return diag.FromErr(setErr)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No Key Management Servers found",
			Detail:   "The API returned an empty list of key management servers.",
		}}
	}

	kmsList, err := flattenKeyManagementServer(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("kms", kmsList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())

	return nil
}

func flattenKeyManagementServer(kmsList []config.KeyManagementServer) ([]map[string]interface{}, error) {
	kmsFlattened := make([]map[string]interface{}, 0, len(kmsList))
	for _, kms := range kmsList {
		accessInformation, flattenErr := flattenAccessInformation(kms.GetAccessInformation())
		if flattenErr != nil {
			return nil, flattenErr
		}
		kmsMap := map[string]interface{}{
			"name":               utils.StringValue(kms.Name),
			"ext_id":             utils.StringValue(kms.ExtId),
			"tenant_id":          utils.StringValue(kms.TenantId),
			"access_information": accessInformation,
			"links":              flattenLinks(kms.Links),
			"creation_timestamp": utils.TimeStringValue(kms.CreationTimestamp),
		}
		kmsFlattened = append(kmsFlattened, kmsMap)
	}
	return kmsFlattened, nil
}
