package nutanix

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	Era "github.com/terraform-providers/terraform-provider-nutanix/client/era"
)

func dataSourceNutanixEraProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraProfilesRead,
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"oracle_database",
					"postgres_database", "sqlserver_database", "mariadb_database",
					"mysql_database"}, false),
			},
			"profile_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"Software", "Compute",
					"Network", "Database_Parameter"}, false),
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixEraProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	engine := ""
	profile_type := ""
	if engineType, ok := d.GetOk("engine"); ok {
		engine = engineType.(string)
	}

	if ptype, ok := d.GetOk("profile_type"); ok {
		profile_type = ptype.(string)
	}

	resp, err := conn.Service.ListProfiles(ctx, engine, profile_type)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("profiles", flattenProfilesResponse(resp)); err != nil {
		return diag.FromErr(err)
	}

	log.Println("HELLLLLOOOOOO")
	aJSON, _ := json.Marshal(resp)
	fmt.Printf("JSON Print - \n%s\n", string(aJSON))

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func flattenVersions(erv []Era.Versions) []map[string]interface{} {
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
			res[k] = ents
		}
		return res
	}
	return nil
}

func flattenProperties(erp []Era.Properties) []map[string]interface{} {
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

func flattenProfilesResponse(erp *Era.ListProfileResponse) []map[string]interface{} {
	if erp != nil {
		lst := []map[string]interface{}{}
		for _, v := range *erp {
			d := map[string]interface{}{}
			if v.ID != "" {
				d["id"] = v.ID
			}
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

			lst = append(lst, d)
		}
		return lst
	}
	return nil
}
