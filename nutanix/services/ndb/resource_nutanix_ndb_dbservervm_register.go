package ndb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBRegisterDBServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBRegisterDBServerCreate,
		ReadContext:   resourceNutanixNDBRegisterDBServerRead,
		UpdateContext: resourceNutanixNDBRegisterDBServerUpdate,
		DeleteContext: resourceNutanixNDBRegisterDBServerDelete,
		Schema: map[string]*schema.Schema{
			"database_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nxcluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"working_directory": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/tmp",
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ssh_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"forced_install": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"postgres_database": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"postgres_software_home": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"credentials": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:     schema.TypeString,
							Required: true,
						},
						"label": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"update_name_description_in_cluster": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// delete values
			"delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"remove": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"soft_remove": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"delete_vgs": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"delete_vm_snapshots": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			// computed
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
			"tags": dataSourceEraDBInstanceTags(),
			"era_created": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"internal": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dbserver_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_cluster_uuid": {
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
			"mac_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"era_drive_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"era_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixNDBRegisterDBServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.DBServerRegisterInput{}

	// build request for dbServerVMs
	if err := buildRegisterDBServerVMRequest(d, req); err != nil {
		return diag.FromErr(err)
	}

	// api to register dbserver
	resp, err := conn.Service.RegisterDBServerVM(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.Entityid)

	// Get Operation ID from response of Response and poll for the operation to get completed.
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
		return diag.Errorf("error waiting for db Server VM (%s) to register: %s", resp.Entityid, errWaitTask)
	}
	log.Printf("NDB database Server VM with %s id is registered successfully", d.Id())
	return resourceNutanixNDBRegisterDBServerRead(ctx, d, meta)
}

func resourceNutanixNDBRegisterDBServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNutanixNDBServerVMRead(ctx, d, meta)
}

func resourceNutanixNDBRegisterDBServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.UpdateDBServerVMRequest{}

	// default for update request
	req.ResetName = utils.BoolPtr(false)
	req.ResetDescription = utils.BoolPtr(false)
	req.ResetCredential = utils.BoolPtr(false)
	req.ResetTags = utils.BoolPtr(true)
	req.ResetDescriptionInNxCluster = utils.BoolPtr(false)
	req.ResetNameInNxCluster = utils.BoolPtr(false)

	if d.HasChange("name") {
		req.Name = utils.StringPtr(d.Get("name").(string))
		req.ResetName = utils.BoolPtr(true)
	}

	if d.HasChange("description") {
		req.Description = utils.StringPtr(d.Get("description").(string))
		req.ResetDescription = utils.BoolPtr(true)
	}
	//nolint:staticcheck
	if _, ok := d.GetOkExists("update_name_description_in_cluster"); ok {
		req.ResetDescriptionInNxCluster = utils.BoolPtr(true)
		req.ResetNameInNxCluster = utils.BoolPtr(true)
	}

	if d.HasChange("credential") {
		req.ResetCredential = utils.BoolPtr(true)

		creds := d.Get("credentials")
		credList := creds.([]interface{})

		credArgs := []*era.VMCredentials{}

		for _, v := range credList {
			val := v.(map[string]interface{})
			cred := &era.VMCredentials{}
			if username, ok := val["username"]; ok {
				cred.Username = utils.StringPtr(username.(string))
			}

			if pass, ok := val["password"]; ok {
				cred.Password = utils.StringPtr(pass.(string))
			}

			if label, ok := val["label"]; ok {
				cred.Label = utils.StringPtr(label.(string))
			}

			credArgs = append(credArgs, cred)
		}
		req.Credentials = credArgs
	}

	resp, err := conn.Service.UpdateDBServerVM(ctx, req, d.Id())
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
	}
	log.Printf("NDB database Server VM with %s id is updated successfully", d.Id())
	return resourceNutanixNDBRegisterDBServerRead(ctx, d, meta)
}

func resourceNutanixNDBRegisterDBServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.DeleteDBServerVMRequest{}
	if delete, ok := d.GetOk("delete"); ok {
		req.Delete = delete.(bool)
	}
	if remove, ok := d.GetOk("remove"); ok {
		req.Remove = remove.(bool)
	}
	if softremove, ok := d.GetOk("soft_remove"); ok {
		req.SoftRemove = softremove.(bool)
	}
	if deleteVgs, ok := d.GetOk("delete_vgs"); ok {
		req.DeleteVgs = deleteVgs.(bool)
	}
	if deleteVMSnaps, ok := d.GetOk("delete_vm_snapshots"); ok {
		req.DeleteVMSnapshots = deleteVMSnaps.(bool)
	}

	resp, err := conn.Service.DeleteDBServerVM(ctx, req, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("Operation to delete dbserver vm with id %s has started, operation id: %s", d.Id(), resp.Operationid)
	opID := resp.Operationid
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
		return diag.Errorf("error waiting for db server VM (%s) to delete: %s", resp.Entityid, errWaitTask)
	}
	log.Printf("NDB database Server VM with %s id is deleted successfully", d.Id())
	return nil
}

func buildRegisterDBServerVMRequest(d *schema.ResourceData, req *era.DBServerRegisterInput) error {
	if dbType, ok := d.GetOk("database_type"); ok {
		req.DatabaseType = utils.StringPtr(dbType.(string))
	}

	if vmip, ok := d.GetOk("vm_ip"); ok {
		req.VMIP = utils.StringPtr(vmip.(string))
	}

	if nxcls, ok := d.GetOk("nxcluster_id"); ok {
		req.NxClusterUUID = utils.StringPtr(nxcls.(string))
	}
	if user, ok := d.GetOk("username"); ok {
		req.Username = utils.StringPtr(user.(string))
	}
	if pass, ok := d.GetOk("password"); ok {
		req.Password = utils.StringPtr(pass.(string))
	}
	if sshkey, ok := d.GetOk("ssh_key"); ok {
		req.SSHPrivateKey = utils.StringPtr(sshkey.(string))
	}
	if workd, ok := d.GetOk("working_directory"); ok {
		req.WorkingDirectory = utils.StringPtr(workd.(string))
	}
	if forcedIns, ok := d.GetOk("forced_install"); ok {
		req.ForcedInstall = utils.BoolPtr(forcedIns.(bool))
	}
	if postgresType, ok := d.GetOk("postgres_database"); ok {
		req.ActionArguments = expandPsRegisterDBServer(postgresType.([]interface{}))
	}
	return nil
}

func expandPsRegisterDBServer(ps []interface{}) []*era.Actionarguments {
	if len(ps) > 0 {
		args := make([]*era.Actionarguments, 0)

		for _, v := range ps {
			val := v.(map[string]interface{})

			if listnerPort, ok := val["listener_port"]; ok {
				args = append(args, &era.Actionarguments{
					Name:  "listener_port",
					Value: listnerPort,
				})
			}
			if psHome, ok := val["postgres_software_home"]; ok {
				args = append(args, &era.Actionarguments{
					Name:  "postgres_software_home",
					Value: psHome,
				})
			}
		}
		return args
	}
	return nil
}
