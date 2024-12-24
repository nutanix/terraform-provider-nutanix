package ndb

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

var EraRefreshCloneTimeout = 15 * time.Minute

func ResourceNutanixNDBCloneRefresh() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBCloneRefreshCreate,
		ReadContext:   resourceNutanixNDBCloneRefreshRead,
		UpdateContext: resourceNutanixNDBCloneRefreshUpdate,
		DeleteContext: resourceNutanixNDBCloneRefreshDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EraRefreshCloneTimeout),
		},
		Schema: map[string]*schema.Schema{
			"clone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_pitr_timestamp"},
			},
			"user_pitr_timestamp": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"snapshot_id"},
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Asia/Calcutta",
			},
		},
	}
}

func resourceNutanixNDBCloneRefreshCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.CloneRefreshInput{}
	cloneID := ""
	if clone, ok := d.GetOk("clone_id"); ok {
		cloneID = clone.(string)
	}

	if snapshotID, ok := d.GetOk("snapshot_id"); ok {
		req.SnapshotID = utils.StringPtr(snapshotID.(string))
	}

	if userPitrTime, ok := d.GetOk("user_pitr_timestamp"); ok {
		req.UserPitrTimestamp = utils.StringPtr(userPitrTime.(string))
	}

	if timezone, ok := d.GetOk("timezone"); ok {
		req.Timezone = utils.StringPtr(timezone.(string))
	}

	resp, err := conn.Service.RefreshClone(ctx, req, cloneID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get Operation ID from response of clone refresh and poll for the operation to get completed.
	opID := resp.Operationid
	if opID == "" {
		return diag.Errorf("error: operation ID is an empty string")
	}
	opReq := era.GetOperationRequest{
		OperationID: opID,
	}

	// Poll for operation here - Operation GET Call
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FAILED"},
		Refresh: eraRefresh(ctx, conn, opReq),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   eraDelay,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for db refresh clone (%s) to create: %s", resp.Entityid, errWaitTask)
	}
	log.Printf("NDB clone Refresh with %s id is completed successfully", d.Id())
	d.SetId(resp.Operationid)
	return nil
}

func resourceNutanixNDBCloneRefreshRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBCloneRefreshUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBCloneRefreshDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
