package iamv2

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixDirectoryServicesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixDirectoryServicesV2Read,
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
			"directory_services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secondary_urls": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_account": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"password": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"open_ldap_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_configuration": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user_object_class": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"user_search_base": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"username_attribute": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"user_group_configuration": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"group_object_class": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"group_search_base": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"group_member_attribute": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"group_member_attribute_value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"group_search_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"white_listed_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixDirectoryServicesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

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
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	resp, err := conn.DirectoryServiceAPIInstance.ListDirectoryServices(page, limit, filter, orderBy, selects)
	if err != nil {
		fmt.Println(err)
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching directory services: %v", errorMessage["message"])
	}

	if resp.Data == nil {
		if err := d.Set("directory_services", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of directory services.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.DirectoryService)

	if err := d.Set("directory_services", flattenDirectoryServicesEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenDirectoryServicesEntities(pr []import1.DirectoryService) []interface{} {
	if len(pr) > 0 {
		dsList := make([]interface{}, len(pr))

		for k, v := range pr {
			ds := make(map[string]interface{})

			ds["ext_id"] = v.ExtId
			ds["name"] = v.Name
			ds["domain_name"] = v.DomainName
			if v.Url != nil {
				ds["url"] = v.Url
			}
			if v.SecondaryUrls != nil {
				ds["secondary_urls"] = v.SecondaryUrls
			}
			if v.DirectoryType != nil {
				ds["directory_type"] = flattenDirectoryType(v.DirectoryType)
			}
			if v.ServiceAccount != nil {
				ds["service_account"] = flattenDsServiceAccount(v.ServiceAccount)
			}
			if v.OpenLdapConfiguration != nil {
				ds["open_ldap_configuration"] = flattenOpenLdapConfig(v.OpenLdapConfiguration)
			}
			if v.GroupSearchType != nil {
				ds["group_search_type"] = flattenGroupSearchType(v.GroupSearchType)
			}
			if v.WhiteListedGroups != nil {
				ds["white_listed_groups"] = v.WhiteListedGroups
			}
			if v.CreatedTime != nil {
				t := v.CreatedTime
				ds["created_time"] = t.String()
			}
			if v.LastUpdatedTime != nil {
				t := v.LastUpdatedTime
				ds["last_updated_time"] = t.String()
			}
			if v.CreatedBy != nil {
				ds["created_by"] = v.CreatedBy
			}

			dsList[k] = ds
		}
		return dsList
	}
	return nil
}
