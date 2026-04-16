package ndb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	Era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixEraClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraClustersRead,
		Schema: map[string]*schema.Schema{
			"clusters": {
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
						"unique_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fqdns": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nx_cluster_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_type": {
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
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ref_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
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
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"reference_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage_threshold_percentage": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"memory_threshold_percentage": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
								},
							},
						},
						"management_server_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_counts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_servers": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"engine_counts": engineCountSchema(),
								},
							},
						},
						"healthy": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixEraClustersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	resp, err := conn.Service.ListClusters(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if e := d.Set("clusters", flattenClustersResponse(resp)); err != nil {
		return diag.FromErr(e)
	}

	uuid, er := uuid.GenerateUUID()

	if er != nil {
		return diag.Errorf("Error generating UUID for era clusters: %+v", err)
	}
	d.SetId(uuid)
	return nil
}

func flattenClustersResponse(crsp *Era.ClusterListResponse) []map[string]interface{} {
	if crsp != nil {
		lst := []map[string]interface{}{}
		for _, v := range *crsp {
			d := map[string]interface{}{}
			d["id"] = v.ID
			d["name"] = v.Name
			d["unique_name"] = v.Uniquename
			d["ip_addresses"] = utils.StringValueSlice(v.Ipaddresses)
			d["fqdns"] = v.Fqdns
			d["nx_cluster_uuid"] = v.Nxclusteruuid
			d["description"] = v.Description
			d["cloud_type"] = v.Cloudtype
			d["date_created"] = v.Datecreated
			d["date_modified"] = v.Datemodified
			d["owner_id"] = v.Ownerid
			d["status"] = v.Status
			d["version"] = v.Version
			d["hypervisor_type"] = v.Hypervisortype
			d["hypervisor_version"] = v.Hypervisorversion
			d["properties"] = flattenClusterProperties(v.Properties)
			d["reference_count"] = v.Referencecount
			d["username"] = v.Username
			d["password"] = v.Password
			d["cloud_info"] = v.Cloudinfo
			d["resource_config"] = flattenResourceConfig(v.Resourceconfig)
			d["management_server_info"] = v.Managementserverinfo
			d["entity_counts"] = flattenEntityCounts(v.EntityCounts)
			d["healthy"] = v.Healthy
			lst = append(lst, d)
		}
		return lst
	}
	return nil
}

func flattenClusterProperties(erp []*Era.Properties) []map[string]interface{} {
	if len(erp) > 0 {
		res := make([]map[string]interface{}, len(erp))

		for k, v := range erp {
			ents := make(map[string]interface{})
			ents["name"] = v.Name
			ents["value"] = v.Value
			ents["secure"] = v.Secure
			ents["ref_id"] = v.RefID
			ents["description"] = v.Description
			res[k] = ents
		}
		return res
	}
	return nil
}

func flattenResourceConfig(rcfg *Era.Resourceconfig) []map[string]interface{} {
	specList := make([]map[string]interface{}, 0)

	if rcfg != nil {
		specs := make(map[string]interface{})

		specs["memory_threshold_percentage"] = rcfg.Memorythresholdpercentage
		specs["storage_threshold_percentage"] = rcfg.Storagethresholdpercentage

		specList = append(specList, specs)
		return specList
	}
	return nil
}
