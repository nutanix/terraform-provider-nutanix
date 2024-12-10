package ndb

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	Era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixEraProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraProfilesRead,
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
			"profiles": {
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
				},
			},
		},
	}
}

func dataSourceNutanixEraProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	engine := ""
	profileType := ""

	if engineType, ok := d.GetOk("engine"); ok {
		engine = engineType.(string)
	}

	if ptype, ok := d.GetOk("profile_type"); ok {
		profileType = ptype.(string)
	}

	//check for profile type and engine
	_, eror := dataSourceEraProfileEngineDiff(ctx, d, meta)
	if eror != nil {
		return diag.FromErr(eror)
	}

	resp, err := conn.Service.ListProfiles(ctx, engine, profileType)
	if err != nil {
		return diag.FromErr(err)
	}

	if e := d.Set("profiles", flattenProfilesResponse(resp)); err != nil {
		return diag.FromErr(e)
	}

	uuid, er := uuid.GenerateUUID()

	if er != nil {
		return diag.Errorf("Error generating UUID for era clusters: %+v", err)
	}
	d.SetId(uuid)
	return nil
}

func flattenVersions(erv []*Era.Versions) []map[string]interface{} {
	if len(erv) > 0 {
		res := make([]map[string]interface{}, len(erv))

		for k, v := range erv {
			ents := make(map[string]interface{})
			ents["id"] = v.ID
			ents["name"] = v.Name
			ents["description"] = v.Description

			ents["status"] = v.Status
			ents["owner"] = v.Owner
			ents["engine_type"] = v.Enginetype

			ents["type"] = v.Type
			ents["topology"] = v.Topology
			ents["db_version"] = v.Dbversion

			ents["system_profile"] = v.Systemprofile
			ents["version"] = v.Version
			ents["profile_id"] = v.Profileid

			ents["published"] = v.Published
			ents["deprecated"] = v.Deprecated

			ents["properties"] = flattenProperties(v.Properties)
			ents["properties_map"] = utils.ConvertMapString(v.Propertiesmap)
			ents["version_cluster_association"] = flattenClusterAssociation(v.VersionClusterAssociation)
			res[k] = ents
		}
		return res
	}
	return nil
}

func flattenProperties(erp []*Era.ProfileProperties) []map[string]interface{} {
	if len(erp) > 0 {
		res := make([]map[string]interface{}, len(erp))

		for k, v := range erp {
			ents := make(map[string]interface{})
			ents["name"] = v.Name
			ents["value"] = v.Value
			ents["secure"] = v.Secure
			res[k] = ents
		}
		return res
	}
	return nil
}

func flattenProfilesResponse(erp *Era.ProfileListResponse) []map[string]interface{} {
	if erp != nil {
		lst := []map[string]interface{}{}
		for _, v := range *erp {
			d := map[string]interface{}{}
			d["id"] = v.ID
			d["name"] = v.Name
			d["description"] = v.Description
			d["status"] = v.Status
			d["owner"] = v.Owner
			d["engine_type"] = v.Enginetype
			d["type"] = v.Type
			d["topology"] = v.Topology
			d["db_version"] = v.Dbversion
			d["system_profile"] = v.Systemprofile
			d["latest_version"] = v.Latestversion
			d["latest_version_id"] = v.Latestversionid
			d["versions"] = flattenVersions(v.Versions)
			if v.Clusteravailability != nil {
				d["cluster_availability"] = flattenClusterAvailability(v.Clusteravailability)
			}
			if v.Nxclusterid != nil {
				d["nx_cluster_id"] = v.Nxclusterid
			}

			lst = append(lst, d)
		}
		return lst
	}
	return nil
}

func flattenClusterAssociation(erc []*Era.VersionClusterAssociation) []map[string]interface{} {
	if len(erc) > 0 {
		res := make([]map[string]interface{}, len(erc))

		for k, v := range erc {
			ercs := map[string]interface{}{}

			ercs["nx_cluster_id"] = v.NxClusterID
			ercs["date_created"] = v.DateCreated
			ercs["date_modified"] = v.DateModified
			ercs["owner_id"] = v.OwnerID
			ercs["status"] = v.Status
			ercs["profile_version_id"] = v.ProfileVersionID
			ercs["properties"] = flattenProperties(v.Properties)
			ercs["optimized_for_provisioning"] = v.OptimizedForProvisioning

			res[k] = ercs
		}
		return res
	}
	return nil
}

func dataSourceEraProfileEngineDiff(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	// check for profile type
	if ptype, vok := d.GetOk("profile_type"); vok {
		if ptype == "Compute" {
			if _, sok := d.GetOk("engine"); sok {
				return false, fmt.Errorf("compute profile type should not be used if engine is given")
			}
		}
	}
	return true, nil
}
