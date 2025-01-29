package selfservice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func ResourceNutanixCalmRunbook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCalmRunbookCreate,
		ReadContext:   resourceNutanixCalmRunbookRead,
		UpdateContext: resourceNutanixCalmRunbookUpdate,
		DeleteContext: resourceNutanixCalmRunbookDelete,
		Schema: map[string]*schema.Schema{
			"runbook_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"runbook_description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_endpoint_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"task_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"task_script_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"task_script": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixCalmRunbookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Calm
	rbName := d.Get("runbook_name").(string)
	rbDesc := d.Get("runbook_description").(string)
	projectUUID := d.Get("project_uuid").(string)
	endpointName := d.Get("default_endpoint_name").(string)
	taskList := d.Get("task_list").([]interface{})

	runbookInput := &calm.RunbookImportInput{}

	runbookSpec := createRunbookSpec(rbName, rbDesc, endpointName, taskList)
	runbookMetadata := createRunbookMetadata(rbName, projectUUID)

	runbookInput.Spec = runbookSpec
	runbookInput.Metadata = runbookMetadata
	runbookInput.APIVersion = "3.0"

	createResp, err := conn.Service.RunbookImport(ctx, runbookInput)
	if err != nil {
		return diag.FromErr(err)
	}

	rbState := createResp.Status["state"].(string)
	rbUUID := createResp.Metadata["uuid"].(string)

	if err := d.Set("state", rbState); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rbUUID)

	return nil
}

func resourceNutanixCalmRunbookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixCalmRunbookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixCalmRunbookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func createRunbookSpec(rbName string, rbDesc string, endpointName string, taskList []interface{}) calm.RunbookSpec {
	rbSpec := calm.RunbookSpec{}
	rbResources := calm.RunbookResources{}
	runbookDef := calm.RunbookDefinition{}

	rbTasks := []calm.TaskDef{}
	rbTasks = append(rbTasks, createDagTask(rbName, taskList))
	for _, task := range taskList {
		taskMap := task.(map[string]interface{})
		rbTasks = append(rbTasks, createExecTask(taskMap))
	}
	runbookDef.Name = rbName + "_runbook"
	runbookDef.MainTaskLocalRef = calm.RefObject{
		Kind: "app_task",
		Name: rbTasks[0].Name,
	}
	runbookDef.TaskDefList = rbTasks
	runbookDef.Description = ""
	runbookDef.VariableList = []interface{}{}
	runbookDef.OutputVariableList = []interface{}{}

	rbResources.Runbook = runbookDef
	rbResources.DefaultTargetRef = calm.RefObject{
		Kind: "app_endpoint",
		Name: endpointName,
	}
	rbResources.ClientAttrs = map[string]interface{}{}
	rbResources.CredentialDefList = []interface{}{}
	rbResources.EndpointDefList = []interface{}{}

	rbSpec.Name = rbName
	rbSpec.Description = rbDesc
	rbSpec.Resources = rbResources

	return rbSpec
}

func createRunbookMetadata(rbName string, projectUUID string) map[string]interface{} {
	rbMetadata := map[string]interface{}{}
	projectRef := map[string]interface{}{}

	rbMetadata["spec_version"] = 1
	rbMetadata["kind"] = "runbook"
	rbMetadata["name"] = rbName

	projectRef["kind"] = "project"
	projectRef["uuid"] = projectUUID
	rbMetadata["project_reference"] = projectRef

	return rbMetadata
}

func createDagTask(rbName string, taskList []interface{}) calm.TaskDef {
	dagTask := &calm.TaskDef{}
	dagTask.Name = rbName + "_dag"
	dagTask.Type = "DAG"

	var edges []map[string]interface{}
	var childTaskRefList []calm.RefObject
	numEdges := len(taskList)
	for ind, task := range taskList {
		taskMap := task.(map[string]interface{})
		if ind < numEdges-1 {
			fromTaskRef := &calm.RefObject{
				Kind: "app_task",
				Name: taskMap["task_name"].(string),
			}
			toTaskRef := calm.RefObject{
				Kind: "app_task",
				Name: taskList[ind+1].(map[string]interface{})["task_name"].(string),
			}
			edge := map[string]interface{}{
				"from_task_reference": fromTaskRef,
				"to_task_reference":   toTaskRef,
			}
			edges = append(edges, edge)
		}
		currTaskRef := calm.RefObject{
			Kind: "app_task",
			Name: taskMap["task_name"].(string),
		}
		childTaskRefList = append(childTaskRefList, currTaskRef)
	}
	dagTask.Attrs = map[string]interface{}{
		"edges": edges,
	}
	dagTask.ChildTaskRefList = childTaskRefList
	dagTask.VariableList = []interface{}{}
	dagTask.StatusMapList = []interface{}{}
	dagTask.Retries = ""
	dagTask.Timeout = ""

	return *dagTask
}

func createExecTask(task map[string]interface{}) calm.TaskDef {
	execTask := &calm.TaskDef{}

	execTask.Type = "EXEC"
	execTask.Name = task["task_name"].(string)
	execTask.Attrs = map[string]interface{}{
		"script_type": "static_py3",
		"script":      task["task_script"].(string),
	}
	execTask.ChildTaskRefList = []calm.RefObject{}
	execTask.VariableList = []interface{}{}
	execTask.StatusMapList = []interface{}{}
	execTask.Retries = ""
	execTask.Timeout = ""

	return *execTask
}
