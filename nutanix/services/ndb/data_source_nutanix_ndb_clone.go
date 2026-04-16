package ndb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixNDBClone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBCloneRead,
		Schema: map[string]*schema.Schema{
			"clone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"clone_name"},
			},
			"clone_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"clone_id"},
			},
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
					},
				},
			},

			// computed

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
	}
}

func dataSourceNutanixNDBCloneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	cloneID, ok := d.GetOk("clone_id")
	cloneName, cok := d.GetOk("clone_name")

	if !ok && !cok {
		return diag.Errorf("atleast one of clone_id or clone_name is required")
	}

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
		}
	} else {
		filterParams.Detailed = "false"
		filterParams.AnyStatus = "false"
		filterParams.LoadDBServerCluster = "false"
		filterParams.TimeZone = "UTC"
	}

	resp, err := conn.Service.GetClone(ctx, cloneID.(string), cloneName.(string), filterParams)
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

	if err := d.Set("date_created", resp.Datecreated); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_modified", resp.Datemodified); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("properties", flattenDBInstanceProperties(resp.Properties)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", flattenDBTags(resp.Tags)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clone", resp.Clone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clustered", resp.Clustered); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database_name", resp.Databasename); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", resp.Type); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database_cluster_type", resp.Databaseclustertype); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", resp.Status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database_status", resp.Databasestatus); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dbserver_logical_cluster_id", resp.Dbserverlogicalclusterid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_machine_id", resp.Timemachineid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("parent_time_machine_id", resp.Parenttimemachineid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_zone", resp.Timezone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("info", flattenDBInfo(resp.Info)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("metric", resp.Metric); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("parent_database_id", resp.ParentDatabaseID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("parent_source_database_id", resp.ParentSourceDatabaseID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("lcm_config", flattenDBLcmConfig(resp.Lcmconfig)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_machine", flattenDBTimeMachine(resp.TimeMachine)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dbserver_logical_cluster", resp.Dbserverlogicalcluster); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database_nodes", flattenDBNodes(resp.Databasenodes)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("linked_databases", flattenDBLinkedDbs(resp.Linkeddatabases)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("databases", resp.Databases); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.ID)

	return nil
}
