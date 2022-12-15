package nutanix

import (
	"context"
	"fmt"
	"log"
	"time"

	era "github.com/terraform-providers/terraform-provider-nutanix/client/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	eraDelay            = 1 * time.Minute
	EraProvisionTimeout = 35 * time.Minute
)

func resourceDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: createDatabaseInstance,
		ReadContext:   readDatabaseInstance,
		UpdateContext: updateDatabaseInstance,
		DeleteContext: deleteDatabaseInstance,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EraProvisionTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"database_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"databasetype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"softwareprofileid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"softwareprofileversionid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"computeprofileid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"networkprofileid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dbparameterprofileid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"newdbservertimezone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"nxclusterid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"sshpublickey": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"createdbserver": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"dbserverid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"clustered": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"autotunestagingdrive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"nodecount": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"vm_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"actionarguments": actionArgumentsSchema(),

			"timemachineinfo": timeMachineInfoSchema(),

			"nodes": nodesSchema(),

			// "tags": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"tag_id": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 			"tag_name": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 			"value": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 		},
			// 	},
			// },

			"properties": {
				Type:        schema.TypeList,
				Description: "List of all the properties",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},

						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
					},
				},
			},
			"postgresql_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"database_size": {
							Type:     schema.TypeString,
							Required: true,
						},
						"auto_tune_staging_drive": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"allocate_pg_hugepage": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"cluster_database": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"auth_method": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "md5",
						},
						"database_names": {
							Type:     schema.TypeString,
							Required: true,
						},
						"db_password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"pre_create_script": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"post_create_script": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			// Computed values

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
			"tags": dataSourceEraDBInstanceTags(),
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
			"metadata": dataSourceEraDBInstanceMetadata(),
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
		},
	}
}

func createDatabaseInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	// check for resource schema validation
	er := schemaValidation("era_provision_database", d)
	if er != nil {
		return diag.FromErr(er)
	}

	log.Println("Creating the request!!!")
	req, err := buildEraRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := conn.Service.ProvisionDatabase(ctx, req)
	if err != nil {
		return diag.Errorf("error while sending request...........:\n %s\n\n", err.Error())
	}
	d.SetId(resp.Entityid)

	// Get Operation ID from response of ProvisionDatabaseResponse and poll for the operation to get completed.
	opID := resp.Operationid
	if opID == "" {
		return diag.Errorf("error: operation ID is an empty string")
	}
	opReq := era.GetOperationRequest{
		OperationID: opID,
	}

	log.Printf("polling for operation with id: %s\n", opID)

	// Poll for operation here - Operation GET Call
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FAILED"},
		Refresh: eraRefresh(ctx, conn, opReq),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   eraDelay,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for db Instance	 (%s) to create: %s", resp.Entityid, errWaitTask)
	}

	return readDatabaseInstance(ctx, d, meta)
}

func buildEraRequest(d *schema.ResourceData) (*era.ProvisionDatabaseRequest, error) {
	return &era.ProvisionDatabaseRequest{
		Databasetype:             utils.StringPtr(d.Get("databasetype").(string)),
		Name:                     utils.StringPtr(d.Get("name").(string)),
		Databasedescription:      utils.StringPtr(d.Get("description").(string)),
		Softwareprofileid:        utils.StringPtr(d.Get("softwareprofileid").(string)),
		Softwareprofileversionid: utils.StringPtr(d.Get("softwareprofileversionid").(string)),
		Computeprofileid:         utils.StringPtr(d.Get("computeprofileid").(string)),
		Networkprofileid:         utils.StringPtr(d.Get("networkprofileid").(string)),
		Dbparameterprofileid:     utils.StringPtr(d.Get("dbparameterprofileid").(string)),
		DatabaseServerID:         utils.StringPtr(d.Get("dbserverid").(string)),
		Timemachineinfo:          buildTimeMachineFromResourceData(d.Get("timemachineinfo").(*schema.Set)),
		Actionarguments:          expandActionArguments(d),
		Createdbserver:           d.Get("createdbserver").(bool),
		Nodecount:                utils.IntPtr(d.Get("nodecount").(int)),
		Nxclusterid:              utils.StringPtr(d.Get("nxclusterid").(string)),
		Sshpublickey:             utils.StringPtr(d.Get("sshpublickey").(string)),
		Clustered:                d.Get("clustered").(bool),
		Nodes:                    buildNodesFromResourceData(d.Get("nodes").(*schema.Set)),
		Autotunestagingdrive:     d.Get("autotunestagingdrive").(bool),
		VMPassword:               utils.StringPtr(d.Get("vm_password").(string)),
		Tags:                     expandTags(d.Get("tags").([]interface{})),
	}, nil
}

func readDatabaseInstance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client).Era
	if c == nil {
		return diag.Errorf("era is nil")
	}

	databaseInstanceID := d.Id()

	resp, err := c.Service.GetDatabaseInstance(ctx, databaseInstanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp != nil {
		if err = d.Set("description", resp.Description); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set("name", resp.Name); err != nil {
			return diag.FromErr(err)
		}

		props := []interface{}{}
		for _, prop := range resp.Properties {
			props = append(props, map[string]interface{}{
				"name":  prop.Name,
				"value": prop.Value,
			})
		}
		if err := d.Set("properties", props); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("date_created", resp.Datecreated); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("date_modified", resp.Datemodified); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("tags", flattenDBTags(resp.Tags)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("clone", resp.Clone); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("internal", resp.Internal); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("placeholder", resp.Placeholder); err != nil {
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

		if err := d.Set("group_info", resp.GroupInfo); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("metadata", flattenDBInstanceMetadata(resp.Metadata)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("metric", resp.Metric); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("category", resp.Category); err != nil {
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
	}
	return nil
}

func updateDatabaseInstance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client).Era
	if c == nil {
		return diag.Errorf("era is nil")
	}

	dbID := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	tags := make([]*era.Tags, 0)
	if d.HasChange("tags") {
		tags = expandTags(d.Get("tags").([]interface{}))
	}

	updateReq := era.UpdateDatabaseRequest{
		Name:             name,
		Description:      description,
		Tags:             tags,
		Resetname:        true,
		Resetdescription: true,
		Resettags:        true,
	}

	res, err := c.Service.UpdateDatabase(ctx, &updateReq, dbID)
	if err != nil {
		return diag.FromErr(err)
	}

	if res != nil {
		if err = d.Set("description", res.Description); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set("name", res.Name); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func deleteDatabaseInstance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(*Client).Era
	if conn == nil {
		return diag.Errorf("era is nil")
	}

	dbID := d.Id()

	req := era.DeleteDatabaseRequest{
		Delete:               true,
		Remove:               false,
		Softremove:           false,
		Forced:               false,
		Deletetimemachine:    true,
		Deletelogicalcluster: true,
	}
	res, err := conn.Service.DeleteDatabase(ctx, &req, dbID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Operation to delete instance with id %s has started, operation id: %s", dbID, res.Operationid)
	opID := res.Operationid
	if opID == "" {
		return diag.Errorf("error: operation ID is an empty string")
	}
	opReq := era.GetOperationRequest{
		OperationID: opID,
	}

	log.Printf("polling for operation with id: %s\n", opID)

	// Poll for operation here - Cluster GET Call
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FAILED"},
		Refresh: eraRefresh(ctx, conn, opReq),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   eraDelay,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for db Instance (%s) to delete: %s", res.Entityid, errWaitTask)
	}
	return nil
}

func expandActionArguments(d *schema.ResourceData) []*era.Actionarguments {
	args := []*era.Actionarguments{}
	if post, ok := d.GetOk("postgresql_info"); ok {
		brr := post.([]interface{})

		for _, arg := range brr {
			val := arg.(map[string]interface{})
			var values interface{}
			if plist, pok := val["listener_port"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "listener_port",
					Value: values,
				})
			}
			if plist, pok := val["database_size"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "database_size",
					Value: values,
				})
			}
			if plist, pok := val["db_password"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "db_password",
					Value: values,
				})
			}
			if plist, pok := val["database_names"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "database_names",
					Value: values,
				})
			}
			if plist, pok := val["auto_tune_staging_drive"]; pok && plist.(bool) {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "auto_tune_staging_drive",
					Value: values,
				})
			}
			if plist, pok := val["allocate_pg_hugepage"]; pok {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "allocate_pg_hugepage",
					Value: values,
				})
			}
			if plist, pok := val["auth_method"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "auth_method",
					Value: values,
				})
			}
			if plist, clok := val["cluster_database"]; clok {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "cluster_database",
					Value: values,
				})
			}
			if plist, clok := val["pre_create_script"]; clok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "pre_create_script",
					Value: values,
				})
			}
			if plist, clok := val["post_create_script"]; clok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "post_create_script",
					Value: values,
				})
			}
		}
	}
	resp := buildActionArgumentsFromResourceData(d.Get("actionarguments").(*schema.Set), args)

	return resp
}

func eraRefresh(ctx context.Context, conn *era.Client, opID era.GetOperationRequest) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opRes, err := conn.Service.GetOperation(opID)
		if err != nil {
			return nil, "FAILED", err
		}
		if *opRes.Status == "5" || *opRes.Status == "4" {
			if *opRes.Status == "5" {
				return opRes, "COMPLETED", nil
			}
			return opRes, "FAILED",
				fmt.Errorf("error_detail: %s, percentage_complete: %s", utils.StringValue(opRes.Message), utils.StringValue(opRes.Percentagecomplete))
		}
		return opRes, "PENDING", nil
	}
}

func expandTags(pr []interface{}) []*era.Tags {
	if len(pr) > 0 {
		tags := make([]*era.Tags, 0)

		for _, v := range pr {
			tag := &era.Tags{}
			val := v.(map[string]interface{})

			if tagName, ok := val["tag_name"]; ok {
				tag.TagName = tagName.(string)
			}

			if tagID, ok := val["tag_id"]; ok {
				tag.TagID = tagID.(string)
			}

			if tagVal, ok := val["value"]; ok {
				tag.Value = tagVal.(string)
			}
			tags = append(tags, tag)
		}
		return tags
	}
	return nil
}
