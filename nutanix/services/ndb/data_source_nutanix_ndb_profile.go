package ndb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	Era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixEraProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraProfileRead,
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"oracle_database",
					"postgres_database", "sqlserver_database", "mariadb_database",
					"mysql_database",
				}, false),
			},
			"profile_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Software", "Compute",
					"Network", "Database_Parameter",
				}, false),
			},
			"profile_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"profile_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topology": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_profile": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"assoc_db_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"assoc_databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"latest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topology": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_profile": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"profile_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"published": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deprecated": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"secure": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"properties_map": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"version_cluster_association": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nx_cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"date_created": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"date_modified": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"owner_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"profile_version_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"properties": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"secure": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"optimized_for_provisioning": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"cluster_availability": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nx_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date_created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"profile_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"nx_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixEraProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	engine := ""
	profileType := ""
	pID := ""
	pName := ""
	profileFilters := &Era.ProfileFilter{}

	if engineType, ok := d.GetOk("engine"); ok {
		engine = engineType.(string)
		profileFilters.Engine = engine
	}

	if ptype, ok := d.GetOk("profile_type"); ok {
		profileType = ptype.(string)
		profileFilters.ProfileType = profileType
	}

	profileID, pIDOk := d.GetOk("profile_id")

	profileName, pNameOk := d.GetOk("profile_name")

	if !pIDOk && !pNameOk {
		return diag.Errorf("please provide one of profile_id or profile_name attributes")
	}
	if pIDOk {
		pID = profileID.(string)
		profileFilters.ProfileID = pID
	}
	if pNameOk {
		pName = profileName.(string)
		profileFilters.ProfileName = pName
	}

	resp, err := conn.Service.GetProfile(ctx, profileFilters)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", resp.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status", resp.Status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owner", resp.Owner); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("engine_type", resp.Enginetype); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", resp.Type); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("topology", resp.Topology); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("db_version", resp.Dbversion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("system_profile", resp.Systemprofile); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("assoc_db_servers", resp.Assocdbservers); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("assoc_databases", resp.Assocdatabases); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("latest_version", resp.Latestversion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("latest_version_id", resp.Latestversionid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("versions", flattenVersions(resp.Versions)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cluster_availability", flattenClusterAvailability(resp.Clusteravailability)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nx_cluster_id", resp.Nxclusterid); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(resp.ID))
	return nil
}

func flattenClusterAvailability(erc []*Era.Clusteravailability) []map[string]interface{} {
	if len(erc) > 0 {
		res := make([]map[string]interface{}, len(erc))

		for k, v := range erc {
			clsAv := map[string]interface{}{}

			clsAv["nx_cluster_id"] = v.Nxclusterid
			clsAv["date_created"] = v.Datecreated
			clsAv["date_modified"] = v.Datemodified
			clsAv["owner_id"] = v.Ownerid
			clsAv["profile_id"] = v.Profileid
			clsAv["status"] = v.Status

			res[k] = clsAv
		}
		return res
	}
	return nil
}
