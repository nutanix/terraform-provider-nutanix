package passwordmanagerv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clusterConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/response"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixPasswordManagersV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixPasswordManagerV2Read,
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
			"passwords": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": schemaForLinks(),
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"prefix_length": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"cluster_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiry_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_hsp_in_use": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixPasswordManagerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	// initialize query params
	var filter, orderBy, selects *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	}

	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	}

	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	}

	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	}

	var extraParam *string = nil
	resp, err := conn.PasswordManagerAPI.ListSystemUserPasswords(page, limit, filter, orderBy, selects, extraParam)
	if err != nil {
		return diag.Errorf("error while fetching system user passwords: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("passwords", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of System User Passwords.",
		}}
	}

	getResp := resp.Data.GetValue().([]clusterConfig.SystemUserPassword)
	if err := d.Set("passwords", flattenPasswordEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	// set the resource id to random uuid
	d.SetId(utils.GenUUID())
	return nil
}

// flatten funcs
func flattenPasswordEntities(passwords []clusterConfig.SystemUserPassword) []map[string]interface{} {
	passwordList := make([]map[string]interface{}, 0)
	for _, password := range passwords {
		passwordMap := make(map[string]interface{})
		if password.ExtId != nil {
			passwordMap["ext_id"] = utils.StringValue(password.ExtId)
		}
		if password.TenantId != nil {
			passwordMap["tenant_id"] = utils.StringValue(password.TenantId)
		}
		if password.Links != nil {
			convertedLinks := flattenLinks(password.Links)
			passwordMap["links"] = convertedLinks
		}
		if password.Username != nil {
			passwordMap["username"] = utils.StringValue(password.Username)
		}
		if password.HostIp != nil {
			hostIPMap := make(map[string]interface{})
			hostIPMap["value"] = utils.StringValue(password.HostIp.Value)
			hostIPMap["prefix_length"] = utils.IntValue(password.HostIp.PrefixLength)
			passwordMap["host_ip"] = []map[string]interface{}{hostIPMap}
		}
		if password.ClusterExtId != nil {
			passwordMap["cluster_ext_id"] = utils.StringValue(password.ClusterExtId)
		}
		if password.LastUpdateTime != nil {
			passwordMap["last_update_time"] = password.LastUpdateTime.Format("2006-01-02T15:04:05Z07:00")
		}
		if password.ExpiryTime != nil {
			passwordMap["expiry_time"] = password.ExpiryTime.Format("2006-01-02T15:04:05Z07:00")
		}
		if password.Status != nil {
			passwordMap["status"] = flattenPasswordStatus(password.Status)
		}
		if password.SystemType != nil {
			passwordMap["system_type"] = flattenPasswordSystemType(password.SystemType)
		}
		if password.HasHspInUse != nil {
			passwordMap["has_hsp_in_use"] = utils.BoolValue(password.HasHspInUse)
		}

		passwordList = append(passwordList, passwordMap)
	}
	return passwordList
}

func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// flatten funcs
func flattenLinks(links []import1.ApiLink) []map[string]interface{} {
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

func flattenPasswordStatus(pr *clusterConfig.PasswordStatus) string {
	const (
		Unknown        = 0
		Redacted       = 1
		Default        = 2
		Secure         = 3
		NoPassword     = 4
		MultipleIssues = 5
	)
	if pr != nil {
		switch *pr {
		case clusterConfig.PasswordStatus(Unknown):
			return "UNKNOWN"
		case clusterConfig.PasswordStatus(Redacted):
			return "REDACTED"
		case clusterConfig.PasswordStatus(Default):
			return "DEFAULT"
		case clusterConfig.PasswordStatus(Secure):
			return "SECURE"
		case clusterConfig.PasswordStatus(NoPassword):
			return "NOPASSWD"
		case clusterConfig.PasswordStatus(MultipleIssues):
			return "MULTIPLE_ISSUES"
		default:
			return "UNKNOWN"
		}
	}
	return "UNKNOWN"
}

func flattenPasswordSystemType(pr *clusterConfig.SystemType) string {
	const (
		Unknown  = 0
		Redacted = 1
		PC       = 2
		AOS      = 3
		AHV      = 4
		IPMI     = 5
	)
	if pr != nil {
		switch *pr {
		case clusterConfig.SystemType(Unknown):
			return "UNKNOWN"
		case clusterConfig.SystemType(Redacted):
			return "REDACTED"
		case clusterConfig.SystemType(PC):
			return "PC"
		case clusterConfig.SystemType(AOS):
			return "AOS"
		case clusterConfig.SystemType(AHV):
			return "AHV"
		case clusterConfig.SystemType(IPMI):
			return "IPMI"
		default:
			return "UNKNOWN"
		}
	}
	return "UNKNOWN"
}
