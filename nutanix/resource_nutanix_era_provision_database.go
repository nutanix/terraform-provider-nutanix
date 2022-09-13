package nutanix

import (
	"context"
	"log"
	"time"

	era "github.com/terraform-providers/terraform-provider-nutanix/client/era"

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
		CreateContext: createDatabaseInstance, // TODO: Use CreateContext etc functions
		ReadContext:   readDatabaseInstance,
		UpdateContext: updateDatabaseInstance,
		DeleteContext: deleteDatabaseInstance,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EraProvisionTimeout),
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
				Optional: true,
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"createdbserver": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"dbserverid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"clustered": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"autotunestagingdrive": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"nodecount": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"actionarguments": actionArgumentsSchema(),

			"timemachineinfo": timeMachineInfoSchema(),

			"nodes": nodesSchema(),

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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"database_size": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_tune_staging_drive": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"allocate_pg_hugepage": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cluster_database": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auth_method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"database_names": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"db_password": {
							Type:     schema.TypeString,
							Optional: true,
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
		},
	}
}

func createDatabaseInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era
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
		Databasetype:             d.Get("databasetype").(string),
		Name:                     d.Get("name").(string),
		Databasedescription:      d.Get("description").(string),
		Softwareprofileid:        d.Get("softwareprofileid").(string),
		Softwareprofileversionid: d.Get("softwareprofileversionid").(string),
		Computeprofileid:         d.Get("computeprofileid").(string),
		Networkprofileid:         d.Get("networkprofileid").(string),
		Dbparameterprofileid:     d.Get("dbparameterprofileid").(string),
		Newdbservertimezone:      d.Get("newdbservertimezone").(string),
		DatabaseServerID:         d.Get("dbserverid").(string),
		Timemachineinfo:          *buildTimeMachineFromResourceData(d.Get("timemachineinfo").(*schema.Set)),
		Actionarguments:          expandActionArguments(d),
		Createdbserver:           d.Get("createdbserver").(bool),
		Nodecount:                d.Get("nodecount").(int),
		Nxclusterid:              d.Get("nxclusterid").(string),
		Sshpublickey:             d.Get("sshpublickey").(string),
		Clustered:                d.Get("clustered").(bool),
		Nodes:                    buildNodesFromResourceData(d.Get("nodes").(*schema.Set)),
		Autotunestagingdrive:     d.Get("autotunestagingdrive").(bool),
	}, nil
}

func readDatabaseInstance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*Client).Era
	if c == nil {
		return diag.Errorf("era is nil")
	}

	databaseInstanceID := d.Id()

	res, err := c.Service.GetDatabaseInstance(ctx, databaseInstanceID)
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

		props := []interface{}{}
		for _, prop := range res.Properties {
			props = append(props, map[string]interface{}{
				"name":  prop.Name,
				"value": prop.Value,
			})
		}
		if err := d.Set("properties", props); err != nil {
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

	updateReq := era.UpdateDatabaseRequest{
		Name:             name,
		Description:      description,
		Tags:             []interface{}{},
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

func expandActionArguments(d *schema.ResourceData) []era.Actionarguments {
	args := []era.Actionarguments{}
	if post, ok := d.GetOk("postgresql_info"); ok {
		brr := post.([]interface{})

		for _, arg := range brr {
			val := arg.(map[string]interface{})
			var values interface{}
			if plist, pok := val["listener_port"]; pok {
				values = plist
				b, ok := tryToConvertBool(plist)
				if ok {
					values = b
				}

				args = append(args, era.Actionarguments{
					Name:  "listener_port",
					Value: values,
				})
			}
			if plist, pok := val["database_size"]; pok {
				values = plist
				b, ok := tryToConvertBool(plist)
				if ok {
					values = b
				}

				args = append(args, era.Actionarguments{
					Name:  "database_size",
					Value: values,
				})
			}
			if plist, pok := val["db_password"]; pok {
				values = plist
				b, ok := tryToConvertBool(plist)
				if ok {
					values = b
				}

				args = append(args, era.Actionarguments{
					Name:  "db_password",
					Value: values,
				})
			}
			if plist, pok := val["database_names"]; pok {
				values = plist
				b, ok := tryToConvertBool(plist)
				if ok {
					values = b
				}

				args = append(args, era.Actionarguments{
					Name:  "database_names",
					Value: values,
				})
			}
			if plist, pok := val["auto_tune_staging_drive"]; pok {
				values = plist
				b, ok := tryToConvertBool(plist)
				if ok {
					values = b
				}

				args = append(args, era.Actionarguments{
					Name:  "auto_tune_staging_drive",
					Value: values,
				})
			}
			if plist, pok := val["allocate_pg_hugepage"]; pok {
				values = plist
				b, ok := tryToConvertBool(plist)
				if ok {
					values = b
				}

				args = append(args, era.Actionarguments{
					Name:  "allocate_pg_hugepage",
					Value: values,
				})
			}
			if plist, pok := val["auth_method"]; pok && len(plist.(string)) > 0 {
				values = plist
				b, ok := tryToConvertBool(plist)
				if ok {
					values = b
				}

				args = append(args, era.Actionarguments{
					Name:  "auth_method",
					Value: values,
				})
			}
			if plist, clok := val["cluster_database"]; clok && len(plist.(string)) > 0 {
				values = plist
				b, ok := tryToConvertBool(plist)
				if ok {
					values = b
				}

				args = append(args, era.Actionarguments{
					Name:  "cluster_database",
					Value: values,
				})
			}
		}
	}
	resp := buildActionArgumentsFromResourceData(d.Get("actionarguments").(*schema.Set), args)

	return resp
}

func eraRefresh(ctx context.Context, conn *era.Client, opId era.GetOperationRequest) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opRes, err := conn.Service.GetOperation(opId)
		if err != nil {
			return nil, "FAILED", err
		}
		if opRes.Status == "5" {
			return opRes, "COMPLETED", nil
		}
		return opRes, "PENDING", nil
	}
}
