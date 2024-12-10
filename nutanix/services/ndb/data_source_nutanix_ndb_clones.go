package ndb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixNDBClones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBClonesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"detailed": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"any_status": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"load_dbserver_cluster": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "UTC",
						},
						"order_by_dbserver_cluster": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"order_by_dbserver_logical_cluster": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
					},
				},
			},
			"clones": {
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
						"date_created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"properties": dataSourceEraDatabaseProperties(),
						"tags":       dataSourceEraDBInstanceTags(),
						"clustered": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"clone": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"database_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dbserver_logical_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_machine_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_time_machine_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"info": dataSourceEraDatabaseInfo(),
						"metric": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"parent_database_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_source_database_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lcm_config":   dataSourceEraLCMConfig(),
						"time_machine": dataSourceEraTimeMachine(),
						"dbserver_logical_cluster": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"database_nodes":   dataSourceEraDatabaseNodes(),
						"linked_databases": dataSourceEraLinkedDatabases(),
						"databases": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixNDBClonesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	filterParams := &era.FilterParams{}
	if filter, fok := d.GetOk("filters"); fok {
		filterList := filter.([]interface{})

		for _, v := range filterList {
			val := v.(map[string]interface{})

			if detailed, dok := val["detailed"]; dok {
				filterParams.Detailed = detailed.(string)
			}

			if anyStatus, aok := val["any_status"]; aok {
				filterParams.AnyStatus = anyStatus.(string)
			}
			if loadDB, lok := val["load_dbserver_details"]; lok {
				filterParams.LoadDBServerCluster = loadDB.(string)
			}

			if timezone, tok := val["timezone"]; tok {
				filterParams.TimeZone = timezone.(string)
			}

			if orderCls, ok := val["order_by_dbserver_cluster"]; ok {
				filterParams.OrderByDBServerCluster = orderCls.(string)
			}

			if orderLogicalCls, ok := val["order_by_dbserver_logical_cluster"]; ok {
				filterParams.OrderByDBServerLogicalCluster = orderLogicalCls.(string)
			}
		}
	} else {
		filterParams.Detailed = "false"
		filterParams.AnyStatus = "false"
		filterParams.LoadDBServerCluster = "false"
		filterParams.TimeZone = "UTC"
		filterParams.OrderByDBServerCluster = "false"
		filterParams.OrderByDBServerLogicalCluster = "false"
	}

	resp, err := conn.Service.ListClones(ctx, filterParams)
	if err != nil {
		return diag.FromErr(err)
	}

	if e := d.Set("clones", flattenDatabaseIntancesList(resp)); e != nil {
		return diag.FromErr(e)
	}

	uuid, er := uuid.GenerateUUID()
	if er != nil {
		return diag.Errorf("Error generating UUID for era clones: %+v", er)
	}
	d.SetId(uuid)
	return nil
}
