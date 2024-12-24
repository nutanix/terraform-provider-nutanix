package ndb

import (
	"context"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBMaintenanceTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBMaintenanceTaskCreate,
		ReadContext:   resourceNutanixNDBMaintenanceTaskRead,
		UpdateContext: resourceNutanixNDBMaintenanceTaskUpdate,
		DeleteContext: resourceNutanixNDBMaintenanceTaskDelete,
		Schema: map[string]*schema.Schema{
			"dbserver_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dbserver_cluster": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"maintenance_window_id": {
				Type:     schema.TypeString,
				Required: true,
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
			//computed
			"entity_task_association": EntityTaskAssocSchema(),
		},
	}
}

func resourceNutanixNDBMaintenanceTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.MaintenanceTasksInput{}

	entities := &era.MaintenanceEntities{}
	if dbserver, ok := d.GetOk("dbserver_id"); ok {
		st := dbserver.([]interface{})
		sublist := make([]*string, len(st))

		for a := range st {
			sublist[a] = utils.StringPtr(st[a].(string))
		}
		entities.EraDBServer = sublist
	}
	if dbserverCls, ok := d.GetOk("dbserver_cluster"); ok {
		st := dbserverCls.([]interface{})
		sublist := make([]*string, len(st))

		for a := range st {
			sublist[a] = utils.StringPtr(st[a].(string))
		}
		entities.EraDBServerCluster = sublist
	}

	req.Entities = entities

	if windowID, ok := d.GetOk("maintenance_window_id"); ok {
		req.MaintenanceWindowID = utils.StringPtr(windowID.(string))
	}

	taskList := make([]*era.Tasks, 0)
	if task, ok := d.GetOk("tasks"); ok {
		tasks := task.([]interface{})

		for _, v := range tasks {
			out := &era.Tasks{}
			value := v.(map[string]interface{})

			if taskType, ok := value["task_type"]; ok {
				out.TaskType = utils.StringPtr(taskType.(string))
			}

			payload := &era.Payload{}
			prepostCommand := &era.PrePostCommand{}
			if preCommand, ok := value["pre_command"]; ok {
				prepostCommand.PreCommand = utils.StringPtr(preCommand.(string))
			}
			if postCommand, ok := value["post_command"]; ok {
				prepostCommand.PostCommand = utils.StringPtr(postCommand.(string))
			}

			payload.PrePostCommand = prepostCommand
			out.Payload = payload

			taskList = append(taskList, out)
		}
	}
	req.Tasks = taskList

	_, err := conn.Service.CreateMaintenanceTask(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, er := uuid.GenerateUUID()

	if er != nil {
		return diag.Errorf("error generating UUID for ndb maintenance tasks: %+v", err)
	}
	d.SetId(uuid)
	log.Printf("NDB maintenance task with %s id is performed", d.Id())
	return resourceNutanixNDBMaintenanceTaskRead(ctx, d, meta)
}

func resourceNutanixNDBMaintenanceTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	maintenanceID := d.Get("maintenance_window_id")

	// check if maintenance id is nil
	if maintenanceID == "" {
		return diag.Errorf("id is required for read operation")
	}

	resp, err := conn.Service.ReadMaintenanceWindow(ctx, maintenanceID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("entity_task_association", flattenEntityTaskAssoc(resp.EntityTaskAssoc))

	return nil
}

func resourceNutanixNDBMaintenanceTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNutanixNDBMaintenanceTaskCreate(ctx, d, meta)
}

func resourceNutanixNDBMaintenanceTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
