package ndb

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

var EraDBProvisionTimeout = 30 * time.Minute

func ResourceNutanixNDBServerVM() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBServerVMCreate,
		ReadContext:   resourceNutanixNDBServerVMRead,
		UpdateContext: resourceNutanixNDBServerVMUpdate,
		DeleteContext: resourceNutanixNDBServerVMDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EraDBProvisionTimeout),
			Delete: schema.DefaultTimeout(EraDBProvisionTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"database_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"software_profile_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"time_machine_id"},
				RequiredWith:  []string{"software_profile_version_id"},
			},
			"software_profile_version_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_machine_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"software_profile_id"},
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_profile_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"compute_profile_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nx_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"latest_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"postgres_database": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_public_key": {
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

			"maintenance_tasks": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maintenance_window_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tasks": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"task_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"OS_PATCHING", "DB_PATCHING"}, false),
									},
									"pre_command": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"post_command": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			// delete arguments for database server vm
			"delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"remove": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
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

func resourceNutanixNDBServerVMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.DBServerInputRequest{}

	// build request for dbServerVMs
	if err := buildDBServerVMRequest(d, req); err != nil {
		return diag.FromErr(err)
	}

	// api to create request

	resp, err := conn.Service.CreateDBServerVM(ctx, req)
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
		return diag.Errorf("error waiting for db Server VM (%s) to create: %s", resp.Entityid, errWaitTask)
	}
	log.Printf("NDB database Server VM with %s id is created successfully", d.Id())
	return resourceNutanixNDBServerVMRead(ctx, d, meta)
}

func resourceNutanixNDBServerVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	// check if d.Id() is nil
	if d.Id() == "" {
		return diag.Errorf("id is required for read operation")
	}
	resp, err := conn.Service.ReadDBServerVM(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

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

	if err := d.Set("tags", flattenDBTags(resp.Tags)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dbserver_cluster_id", resp.DbserverClusterID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vm_cluster_name", resp.VMClusterName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vm_cluster_uuid", resp.VMClusterUUID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip_addresses", utils.StringValueSlice(resp.IPAddresses)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("fqdns", resp.Fqdns); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("mac_addresses", utils.StringValueSlice(resp.MacAddresses)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", resp.Type); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", resp.Status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("client_id", resp.ClientID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("era_drive_id", resp.EraDriveID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("era_version", resp.EraVersion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vm_timezone", resp.VMTimeZone); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNutanixNDBServerVMUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.UpdateDBServerVMRequest{}

	// setting default values
	req.ResetName = utils.BoolPtr(false)
	req.ResetDescription = utils.BoolPtr(false)
	req.ResetCredential = utils.BoolPtr(false)
	req.ResetTags = utils.BoolPtr(false)

	if d.HasChange("description") {
		req.Description = utils.StringPtr(d.Get("description").(string))
		req.ResetDescription = utils.BoolPtr(true)
	}

	if d.HasChange("postgres_database") {
		ps := d.Get("postgres_database").([]interface{})[0].(map[string]interface{})

		vmName := ps["vm_name"]
		req.Name = utils.StringPtr(vmName.(string))
		req.ResetName = utils.BoolPtr(true)
	}

	if d.HasChange("tags") {
		req.Tags = expandTags(d.Get("tags").([]interface{}))
		req.ResetTags = utils.BoolPtr(true)
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

	log.Printf("NDB database with %s id updated successfully", d.Id())
	return nil
}

func resourceNutanixNDBServerVMDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	res, err := conn.Service.DeleteDBServerVM(ctx, req, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Operation to delete dbserver vm with id %s has started, operation id: %s", d.Id(), res.Operationid)
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
		Timeout: d.Timeout(schema.TimeoutDelete),
		Delay:   eraDelay,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for db server VM (%s) to delete: %s", res.Entityid, errWaitTask)
	}
	log.Printf("NDB database Server VM with %s id is deleted successfully", d.Id())
	return nil
}

func buildDBServerVMRequest(d *schema.ResourceData, res *era.DBServerInputRequest) error {
	if dbType, ok := d.GetOk("database_type"); ok {
		res.DatabaseType = utils.StringPtr(dbType.(string))
	}

	if softwareProfile, ok := d.GetOk("software_profile_id"); ok {
		res.SoftwareProfileID = utils.StringPtr(softwareProfile.(string))
	}

	if softwareVersion, ok := d.GetOk("software_profile_version_id"); ok {
		res.SoftwareProfileVersionID = utils.StringPtr(softwareVersion.(string))
	}

	if LatestSnapshot, ok := d.GetOk("latest_snapshot"); ok {
		res.LatestSnapshot = LatestSnapshot.(bool)
	}

	if timeMachine, ok := d.GetOk("time_machine_id"); ok {
		res.TimeMachineID = utils.StringPtr(timeMachine.(string))

		// if snapshot id is provided
		if snapshotid, ok := d.GetOk("snapshot_id"); ok {
			res.SnapshotID = utils.StringPtr(snapshotid.(string))
			res.LatestSnapshot = false
		} else {
			res.LatestSnapshot = true
		}
	}

	if NetworkProfile, ok := d.GetOk("network_profile_id"); ok {
		res.NetworkProfileID = utils.StringPtr(NetworkProfile.(string))
	}

	if ComputeProfile, ok := d.GetOk("compute_profile_id"); ok {
		res.ComputeProfileID = utils.StringPtr(ComputeProfile.(string))
	}

	if ClusterID, ok := d.GetOk("nx_cluster_id"); ok {
		res.NxClusterID = utils.StringPtr(ClusterID.(string))
	}

	if VMPass, ok := d.GetOk("vm_password"); ok {
		res.VMPassword = utils.StringPtr(VMPass.(string))
	}

	if desc, ok := d.GetOk("description"); ok {
		res.Description = utils.StringPtr(desc.(string))
	}

	if postgresDatabase, ok := d.GetOk("postgres_database"); ok && len(postgresDatabase.([]interface{})) > 0 {
		res.ActionArguments = expandDBServerPostgresInput(postgresDatabase.([]interface{}))
	}

	if maintenance, ok := d.GetOk("maintenance_tasks"); ok {
		res.MaintenanceTasks = expandMaintenanceTasks(maintenance.([]interface{}))
	}
	return nil
}

func expandDBServerPostgresInput(pr []interface{}) []*era.Actionarguments {
	if len(pr) > 0 {
		args := make([]*era.Actionarguments, 0)

		for _, v := range pr {
			val := v.(map[string]interface{})

			if vmName, ok := val["vm_name"]; ok {
				args = append(args, &era.Actionarguments{
					Name:  "vm_name",
					Value: vmName,
				})
			}
			if clientKey, ok := val["client_public_key"]; ok && len(clientKey.(string)) > 0 {
				args = append(args, &era.Actionarguments{
					Name:  "client_public_key",
					Value: clientKey,
				})
			}
		}
		return args
	}
	return nil
}
