package nutanix

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	era "github.com/terraform-providers/terraform-provider-nutanix/client/era"
)

func dataSourceNutanixEraDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraDatabaseIntancesRead,
		Schema: map[string]*schema.Schema{
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
						"owner_id": {
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
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"clustered": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"clone": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"era_created": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"internal": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"placeholder": {
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
						"group_info": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"metric": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
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
						"database_group_state_info": {
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
	conn := meta.(*Client).Era

	resp, err := conn.Service.ListDatabaseInstance(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if e := d.Set("database_instances", flattenDatabaseIntancesList(resp)); err != nil {
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

			d["category"] = data.Category
			d["clone"] = data.Clone
			d["clustered"] = data.Clustered
			d["database_group_state_info"] = data.DatabaseGroupStateInfo
			d["database_cluster_type"] = data.Databaseclustertype
			d["database_name"] = data.Databasename
			d["database_nodes"] = flattenDBNodes(data.Databasenodes)
			d["databases"] = data.Databases
			d["database_status"] = data.Databasestatus
			d["date_created"] = data.Datecreated
			d["date_modified"] = data.Datemodified
			d["dbserver_logical_cluster"] = data.Dbserverlogicalcluster
			d["dbserver_logical_cluster_id"] = data.Dbserverlogicalclusterid
			d["description"] = data.Description
			d["group_info"] = data.GroupInfo
			d["id"] = data.ID
			d["info"] = flattenDBInfo(data.Info)
			d["internal"] = data.Internal
			d["lcm_config"] = flattenDBLcmConfig(data.Lcmconfig)
			d["linked_databases"] = flattenDBLinkedDbs(data.Linkeddatabases)
			d["metadata"] = flattenDBInstanceMetadata(data.Metadata)
			d["metric"] = data.Metric
			d["name"] = data.Name
			d["owner_id"] = data.Ownerid
			d["parent_database_id"] = data.ParentDatabaseID
			d["parent_source_database_id"] = data.ParentSourceDatabaseID
			d["parent_time_machine_id"] = data.Parenttimemachineid
			d["placeholder"] = data.Placeholder
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
