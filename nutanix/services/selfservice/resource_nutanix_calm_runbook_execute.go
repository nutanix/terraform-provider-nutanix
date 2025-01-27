package selfservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"strings"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
// 	ERROR              = "ERROR"
// 	DEFAULTWAITTIMEOUT = 60
)

func ResourceNutanixCalmRunbookExecute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCalmRunbookExecuteCreate,
		ReadContext:   resourceNutanixCalmRunbookExecuteRead,
		UpdateContext: resourceNutanixCalmRunbookExecuteUpdate,
		DeleteContext: resourceNutanixCalmRunbookExecuteDelete,
		Schema: map[string]*schema.Schema{
			"rb_name": {
				Type:          schema.TypeString,
				Optional:      true,
			},
			"rb_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
			},
			"spec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeString,
				Computed: true,
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
			fmt.Println("Error unmarshalling RBName:", err)
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

    // execute runbook using the uuid
    input := &calm.RunbookProvisionInput{}

	output, err := conn.Service.ExecuteRunbook(ctx, rb_uuid, input)
	if err != nil {
		return diag.Errorf("Error executing Runbook: %s", err)
	}

	runlogUUID := output.Status.RunlogUUID
	fmt.Println("Response:", runlogUUID)
    d.SetId(runlogUUID)

    // poll till action is completed
	rbStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING"},
		Target:  []string{"SUCCESS","FAILURE","WARNING"},
		Refresh: RbRunlogStateRefreshFunc(ctx, conn, runlogUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
		Delay:   5 * time.Second,
	}

	if _, errWaitTask := rbStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for runbook to finish execute(%s): %s", errWaitTask)
	}

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
		fmt.Println("V State: ", v.Status.State)
		fmt.Println("V: ", *v)

		runlogstate := utils.StringValue(v.Status.State)

		fmt.Printf("Runlog State: %s\n", runlogstate)
		log.Printf("Runlog State: %s\n", runlogstate)

		return v, runlogstate, nil
	}
}

