package licensingv2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/models/licensing/v4/agreements"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixEULAV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixEULAReadV2,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"acceptances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accepted_by": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"acceptance_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixEULAReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LicensingAPI

	resp, err := conn.LicensingEULAAPIInstance.GetEula()
	if err != nil {
		return diag.Errorf("error while retrieving EULA: %v", err)
	}

	if resp == nil {
		return diag.Errorf("error while retrieving EULA: empty response")
	}

	getResp := resp.Data.GetValue().(import1.Eula)
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("content", getResp.Content); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_time", flattenTime(getResp.UpdatedTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", getResp.Version); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_enabled", getResp.IsEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("acceptances", flattenAcceptances(getResp.Acceptances)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

// flatten funcs
func flattenLinks(links []response.ApiLink) []map[string]interface{} {
	if len(links) > 0 {
		linkList := make([]map[string]interface{}, 0)
		for _, link := range links {
			linkMap := make(map[string]interface{})
			if link.Rel != nil {
				linkMap["rel"] = utils.StringValue(link.Rel)
			}
			if link.Href != nil {
				linkMap["href"] = utils.StringValue(link.Href)
			}

			linkList = append(linkList, linkMap)
		}
		return linkList
	}
	return nil
}

func flattenAcceptances(acceptances []import1.Acceptance) []map[string]interface{} {
	if len(acceptances) > 0 {
		acceptanceList := make([]map[string]interface{}, 0)
		for _, acceptance := range acceptances {
			acceptanceMap := make(map[string]interface{})
			if acceptance.AcceptedBy != nil {
				acceptanceMap["accepted_by"] = flattenEndUser(acceptance.AcceptedBy)
			}
			if acceptance.AcceptanceTime != nil {
				acceptanceMap["acceptance_time"] = flattenTime(acceptance.AcceptanceTime)
			}

			acceptanceList = append(acceptanceList, acceptanceMap)
		}
		return acceptanceList
	}
	return nil
}

func flattenTime(inTime *time.Time) string {
	if inTime != nil {
		return inTime.UTC().Format(time.RFC3339)
	}
	return ""
}

func flattenEndUser(endUser *import1.EndUser) map[string]interface{} {
	if endUser == nil {
		return nil
	}
	m := make(map[string]interface{})
	m["company_name"] = endUser.CompanyName
	m["job_title"] = endUser.JobTitle
	m["login_id"] = endUser.LoginId
	m["user_name"] = endUser.UserName
	return m
}
