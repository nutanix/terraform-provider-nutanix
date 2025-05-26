package dataprotectionv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/config"
	dataprtotectionPrismConfig "github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRecoveryPointReplicateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRecoveryPointReplicateV2Create,
		ReadContext:   ResourceNutanixRecoveryPointReplicateV2Read,
		UpdateContext: ResourceNutanixRecoveryPointReplicateV2Update,
		DeleteContext: ResourceNutanixRecoveryPointReplicateV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pc_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replicated_rp_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// ResourceNutanixRecoveryPointReplicateV2Create to Replicate Recovery Points
func ResourceNutanixRecoveryPointReplicateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ResourceNutanixRecoveryPointReplicateV2Create \n")

	conn := meta.(*conns.Client).DataProtectionAPI

	body := config.RecoveryPointReplicationSpec{}
	rpExtID := d.Get("ext_id").(string)

	if pcExtID, ok := d.GetOk("pc_ext_id"); ok {
		body.PcExtId = utils.StringPtr(pcExtID.(string))
	}
	if clusterExtID, ok := d.GetOk("cluster_ext_id"); ok {
		body.ClusterExtId = utils.StringPtr(clusterExtID.(string))
	}

	resp, err := conn.RecoveryPoint.ReplicateRecoveryPoint(utils.StringPtr(rpExtID), &body)
	if err != nil {
		return diag.Errorf("error while replicating recovery point: %v", err)
	}

	TaskRef := resp.Data.GetValue().(dataprtotectionPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for recovery point: (%s) to replicate: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching recovery point UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	aJSON, _ := json.Marshal(rUUID)
	log.Printf("[DEBUG] Replicate Recovery Point Task Details: %v", string(aJSON))

	// set the UUID of the replicated recovery point
	uuid := rUUID.CompletionDetails[0].Value
	d.SetId(uuid.GetValue().(string))
	d.Set("replicated_rp_ext_id", uuid.GetValue().(string))

	return ResourceNutanixRecoveryPointReplicateV2Read(ctx, d, meta)
}

func ResourceNutanixRecoveryPointReplicateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixRecoveryPointReplicateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixRecoveryPointReplicateV2Read(ctx, d, meta)
}

func ResourceNutanixRecoveryPointReplicateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
