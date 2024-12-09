package ndb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixEraDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraDatabaseIntancesRead,
		Schema: map[string]*schema.Schema{
			"database_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"oracle_database",
					"postgres_database", "sqlserver_database", "mariadb_database",
					"mysql_database", "mssql_database", "saphana_database", "mongodb_database",
				}, false),
			},
			"database_instances": {
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
						"dbserver_logical_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_machine_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"info":     dataSourceEraDatabaseInfo(),
						"metadata": dataSourceEraDBInstanceMetadata(),
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

func dataSourceNutanixEraDatabaseIntancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	var resp *era.ListDatabaseInstance
	var err error
	if dbEng, ok := d.GetOk("database_type"); ok {
		// todo : when era have query params for db egine type call , API here
		// filter the database based on db engine type provided
		respon, er := conn.Service.ListDatabaseInstance(ctx)
		if er != nil {
			return diag.FromErr(er)
		}
		resp, err = filterDatabaseBasedOnDatabaseEngine(respon, dbEng.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		resp, err = conn.Service.ListDatabaseInstance(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if e := d.Set("database_instances", flattenDatabaseIntancesList(resp)); e != nil {
		return diag.FromErr(e)
	}

	uuid, er := uuid.GenerateUUID()

	if er != nil {
		return diag.Errorf("Error generating UUID for era clusters: %+v", err)
	}
	d.SetId(uuid)
	return nil
}

func flattenDatabaseIntancesList(db *era.ListDatabaseInstance) []map[string]interface{} {
	if db != nil {
		lst := []map[string]interface{}{}
		for _, data := range *db {
			d := map[string]interface{}{}

			d["clone"] = data.Clone
			d["clustered"] = data.Clustered
			d["database_cluster_type"] = data.Databaseclustertype
			d["database_name"] = data.Databasename
			d["database_nodes"] = flattenDBNodes(data.Databasenodes)
			d["databases"] = data.Databases
			d["date_created"] = data.Datecreated
			d["date_modified"] = data.Datemodified
			d["dbserver_logical_cluster"] = data.Dbserverlogicalcluster
			d["dbserver_logical_cluster_id"] = data.Dbserverlogicalclusterid
			d["description"] = data.Description
			d["id"] = data.ID
			d["info"] = flattenDBInfo(data.Info)
			d["lcm_config"] = flattenDBLcmConfig(data.Lcmconfig)
			d["linked_databases"] = flattenDBLinkedDbs(data.Linkeddatabases)
			d["metric"] = data.Metric
			d["name"] = data.Name
			d["parent_database_id"] = data.ParentDatabaseID
			d["properties"] = flattenDBInstanceProperties(data.Properties)
			d["status"] = data.Status
			d["tags"] = flattenDBTags(data.Tags)
			d["time_machine"] = flattenDBTimeMachine(data.TimeMachine)
			d["time_machine_id"] = data.Timemachineid
			d["time_zone"] = data.Timezone
			d["type"] = data.Type

			lst = append(lst, d)
		}
		return lst
	}
	return nil
}

func filterDatabaseBasedOnDatabaseEngine(resp *era.ListDatabaseInstance, dbengine string) (*era.ListDatabaseInstance, error) {
	found := make(era.ListDatabaseInstance, 0)

	for _, v := range *resp {
		if dbengine == v.Type {
			found = append(found, v)
		}
	}
	return &found, nil
}
