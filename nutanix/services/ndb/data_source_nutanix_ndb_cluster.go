package ndb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixEraCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraClusterRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"cluster_name"},
			},
			"cluster_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"cluster_id"},
			},
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
	}
}

func dataSourceNutanixEraClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	clusterID, iok := d.GetOk("cluster_id")
	clusterName, nok := d.GetOk("cluster_name")

	if !iok && !nok {
		return diag.Errorf("please provide one of cluster_id or cluster_name attributes")
	}

	resp, err := conn.Service.GetCluster(ctx, clusterID.(string), clusterName.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", resp.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("unique_name", resp.Uniquename); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip_addresses", resp.Ipaddresses); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("fqdns", resp.Fqdns); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nx_cluster_uuid", resp.Nxclusteruuid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_type", resp.Cloudtype); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_created", resp.Datecreated); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_modified", resp.Datemodified); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owner_id", resp.Ownerid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status", resp.Status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("version", resp.Version); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("hypervisor_type", resp.Hypervisortype); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("hypervisor_version", resp.Hypervisorversion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("properties", flattenClusterProperties(resp.Properties)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("reference_count", resp.Referencecount); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("username", resp.Username); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("password", resp.Password); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_info", resp.Cloudinfo); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("resource_config", flattenResourceConfig(resp.Resourceconfig)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("management_server_info", resp.Managementserverinfo); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("entity_counts", flattenEntityCounts(resp.EntityCounts)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("healthy", resp.Healthy); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.ID)
	return nil
}

func flattenEntityCounts(pr *era.EntityCounts) []interface{} {
	if pr != nil {
		res := make([]interface{}, 0)

		entity := map[string]interface{}{}

		entity["db_servers"] = pr.DBServers
		entity["engine_counts"] = flattenEngineCounts(pr.EngineCounts)

		res = append(res, entity)
		return res
	}
	return nil
}

func flattenEngineCounts(pr *era.EngineCounts) []interface{} {
	if pr != nil {
		engineCounts := make([]interface{}, 0)
		engine := map[string]interface{}{}

		engine["mariadb_database"] = flattenProfileTmsCount(pr.MariadbDatabase)
		engine["mongodb_database"] = flattenProfileTmsCount(pr.MongodbDatabase)
		engine["mysql_database"] = flattenProfileTmsCount(pr.MySQLDatabase)
		engine["oracle_database"] = flattenProfileTmsCount(pr.OracleDatabase)
		engine["postgres_database"] = flattenProfileTmsCount(pr.PostgresDatabase)
		engine["saphana_database"] = flattenProfileTmsCount(pr.SaphanaDatabase)
		engine["sqlserver_database"] = flattenProfileTmsCount(pr.SqlserverDatabase)

		engineCounts = append(engineCounts, engine)
		return engineCounts
	}
	return nil
}

func flattenProfileTmsCount(pr *era.ProfileTimeMachinesCount) []interface{} {
	if pr != nil {
		engineCounts := make([]interface{}, 0)
		count := map[string]interface{}{}

		count["profiles"] = flattenProfilesCount(pr.Profiles)
		count["time_machines"] = pr.TimeMachines
		engineCounts = append(engineCounts, count)
		return engineCounts
	}
	return nil
}

func flattenProfilesCount(pr *era.ProfilesEntity) []interface{} {
	if pr != nil {
		profileCounts := make([]interface{}, 0)
		count := map[string]interface{}{}

		count["compute"] = pr.Compute
		count["database_parameter"] = pr.DatabaseParameter
		count["software"] = pr.Software
		count["network"] = pr.Network
		count["storage"] = pr.Storage
		count["windows_domain"] = pr.WindowsDomain

		profileCounts = append(profileCounts, count)
		return profileCounts
	}
	return nil
}

func engineCountSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"oracle_database":    profileTimeMachineCountSchema(),
				"postgres_database":  profileTimeMachineCountSchema(),
				"mongodb_database":   profileTimeMachineCountSchema(),
				"sqlserver_database": profileTimeMachineCountSchema(),
				"saphana_database":   profileTimeMachineCountSchema(),
				"mariadb_database":   profileTimeMachineCountSchema(),
				"mysql_database":     profileTimeMachineCountSchema(),
			},
		},
	}
}

func profileTimeMachineCountSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"profiles": profilesCountSchema(),
				"time_machines": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func profilesCountSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"windows_domain": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"software": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"compute": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"network": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"storage": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"database_parameter": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}
