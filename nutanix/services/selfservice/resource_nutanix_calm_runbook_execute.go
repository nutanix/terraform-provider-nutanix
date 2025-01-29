package selfservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixCalmRunbookExecute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCalmRunbookExecuteCreate,
		ReadContext:   resourceNutanixCalmRunbookExecuteRead,
		UpdateContext: resourceNutanixCalmRunbookExecuteUpdate,
		DeleteContext: resourceNutanixCalmRunbookExecuteDelete,
		Schema: map[string]*schema.Schema{
			"rb_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rb_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_variable_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"variable_list": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixCalmRunbookExecuteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Calm

	var rb_uuid string

	// Check if rb_uuid is provided
	if rbUUID, ok := d.GetOk("rb_uuid"); ok {
		rb_uuid = rbUUID.(string)
	} else {
		// Fetch rb_uuid from rb_name if rb_uuid is not provided
		rb_name := d.Get("rb_name").(string)

		rbFilter := &calm.RunbookListInput{}
		rbFilter.Filter = fmt.Sprintf("name==%s;state!=DELETED", rb_name)

		rbNameResp, err := conn.Service.ListRunbook(ctx, rbFilter)
		if err != nil {
			return diag.FromErr(err)
		}

		var RbNameStatus []interface{}
		if err := json.Unmarshal([]byte(rbNameResp.Entities), &RbNameStatus); err != nil {
			return diag.FromErr(err)
		}

		if len(RbNameStatus) == 0 {
			return diag.Errorf("No runbooks found with name %s", rb_name)
		}

		entities := RbNameStatus[0].(map[string]interface{})
		if entity, ok := entities["metadata"].(map[string]interface{}); ok {
			rb_uuid = entity["uuid"].(string)
		}
	}

	d.Set("rb_uuid", rb_uuid)

	var args []calm.RunbookArgs
	if variableList, ok := d.Get("variable_list").([]interface{}); ok {
		for _, v := range variableList {
			variable := v.(map[string]interface{})
			log.Printf("%v", variable)
			args = append(args, calm.RunbookArgs{
				Name:  variable["name"].(string),
				Value: variable["value"].(string),
			})
		}
	}

	input := &calm.RunbookProvisionInput{}
	inputSpec := &calm.RunbookProvisionSpec{}

	inputSpec.Args = args
	input.Spec = *inputSpec

	output, err := conn.Service.ExecuteRunbook(ctx, rb_uuid, input)
	if err != nil {
		return diag.Errorf("Error executing Runbook: %s", err)
	}

	runlogUUID := output.Status.RunlogUUID
	d.SetId(runlogUUID)

	// poll till action is completed
	rbStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "ABORTING"},
		Target:  []string{"SUCCESS", "FAILURE", "WARNING", "ABORTED"},
		Refresh: RbRunlogStateRefreshFunc(ctx, conn, runlogUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
		Delay:   5 * time.Second,
	}

	if _, errWaitTask := rbStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for runbook to finish execute: %s", errWaitTask)
	}

	// Fetch the final state after execution
	finalRunlog, err := conn.Service.RbRunlogs(ctx, runlogUUID)
	if err != nil {
		return diag.Errorf("Error fetching final runlog state: %s", err)
	}

	// Set the runlog state in the resource data
	if err := d.Set("state", finalRunlog.Status.State); err != nil {
		return diag.Errorf("Error setting state in resource data: %s", err)
	}

	outputVariable, _, err := RbOutputFunc(ctx, conn, runlogUUID)
	if err != nil {
		return diag.Errorf("Error fetching output variables: %s", err)
	}
	d.Set("output_variable_list", outputVariable)

	return nil

}

func resourceNutanixCalmRunbookExecuteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNutanixCalmRunbookExecuteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNutanixCalmRunbookExecuteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func RbRunlogStateRefreshFunc(ctx context.Context, client *calm.Client, runlogUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.Service.RbRunlogs(ctx, runlogUUID)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "INVALID_UUID") {
				return v, ERROR, nil
			}
			return nil, "", err
		}

		runlogstate := utils.StringValue(v.Status.State)
		return v, runlogstate, nil
	}
}

func RbOutputFunc(ctx context.Context, client *calm.Client, runlogUUID string) (interface{}, string, error) {
	v, err := client.Service.RbRunlogs(ctx, runlogUUID)
	if err != nil {
		if strings.Contains(err.Error(), "INVALID_UUID") {
			return nil, ERROR, nil
		}
		return nil, "", err
	}

	var outputVariables []map[string]interface{}
	for _, outputVar := range v.Status.OutputVariableList {
		// Append the output variable details as a map
		outputVariables = append(outputVariables, map[string]interface{}{
			"name":  outputVar.Name,
			"value": outputVar.Value,
		})
	}

	return outputVariables, "", nil
}
