// Package securityv2 provides resources for managing security-related configurations in Nutanix.
package securityv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/security/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixKeyManagementServerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixKeyManagementServerV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_information": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_secret": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"credential_expiry_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"truncated_client_secret": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixKeyManagementServerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.KeyManagementServersAPIInstance.GetKeyManagementServerById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching key management server : %v", err)
	}

	getRespValue, ok := resp.Data.GetValue().(config.KeyManagementServer)
	if !ok {
		return diag.Errorf("error: unexpected response type from get API, expected KeyManagementServer")
	}
	getResp := getRespValue

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	accessInfo, flattenErr := flattenAccessInformation(getResp.AccessInformation)
	if flattenErr != nil {
		return diag.FromErr(flattenErr)
	}
	if err := d.Set("access_information", accessInfo); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenAccessInformation(azureAccessInformation *config.AzureAccessInformation) ([]map[string]interface{}, error) {
	if azureAccessInformation == nil {
		return nil, fmt.Errorf("access information is nil")
	}

	flattenedAccessInfo := make([]map[string]interface{}, 1)
	flattenedAccessInfo[0] = map[string]interface{}{
		"endpoint_url":            utils.StringValue(azureAccessInformation.EndpointUrl),
		"key_id":                  utils.StringValue(azureAccessInformation.KeyId),
		"tenant_id":               utils.StringValue(azureAccessInformation.TenantId),
		"client_id":               utils.StringValue(azureAccessInformation.ClientId),
		"client_secret":           utils.StringValue(azureAccessInformation.ClientSecret),
		"truncated_client_secret": utils.StringValue(azureAccessInformation.TruncatedClientSecret),
		"credential_expiry_date":  utils.TimeValue(azureAccessInformation.CredentialExpiryDate).Format("2006-01-02"),
	}
	return flattenedAccessInfo, nil
}

func flattenLinks(links []response.ApiLink) []interface{} {
	if len(links) > 0 {
		flattenedLinks := make([]interface{}, len(links))

		for k, v := range links {
			link := make(map[string]interface{})

			if v.Href != nil {
				link["href"] = v.Href
			}
			if v.Rel != nil {
				link["rel"] = v.Rel
			}
			flattenedLinks[k] = link
		}
		return flattenedLinks
	}
	return nil
}
