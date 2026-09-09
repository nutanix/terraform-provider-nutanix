package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	prismConfigLc "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixRefreshNodeV2 refreshes a node's information by fetching the latest
// data from the hardware provider. This is a standalone mutating action resource.
func ResourceNutanixRefreshNodeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRefreshNodeV2Create,
		ReadContext:   ResourceNutanixRefreshNodeV2Read,
		DeleteContext: ResourceNutanixRefreshNodeV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceNutanixRefreshNodeV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.NodesAPIInstance.RefreshNode(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while refreshing node : %v", err)
	}

	taskRef := resp.Data.GetValue().(prismConfigLc.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for node (%s) refresh to complete: %s", extID, errWaitTask)
	}

	d.SetId(extID)
	return ResourceNutanixRefreshNodeV2Read(ctx, d, meta)
}

// ResourceNutanixRefreshNodeV2Read is a no-op for this action resource; refresh has no persistent state.
func ResourceNutanixRefreshNodeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// ResourceNutanixRefreshNodeV2Delete clears the resource from state; refresh cannot be undone.
func ResourceNutanixRefreshNodeV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
