package volumesv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	taskPoll "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	volumesPrism "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/prism/v4/config"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixVolumeAttachVMToVolumeGroupV2 Attach an AHV VM to the given Volume Group.
func ResourceNutanixVolumeAttachVMToVolumeGroupV2() *schema.Resource {
	return &schema.Resource{
		Description:   "Attaches VM to a Volume Group identified by {extId}.",
		CreateContext: ResourceNutanixVolumeAttachVMToVolumeGroupV2Create,
		ReadContext:   ResourceNutanixVolumeAttachVMToVolumeGroupV2Read,
		UpdateContext: ResourceNutanixVolumeAttachVMToVolumeGroupV2Update,
		DeleteContext: ResourceNutanixVolumeAttachVMToVolumeGroupV2Delete,

		Schema: map[string]*schema.Schema{
			"volume_group_ext_id": {
				Description: "The external identifier of the volume group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vm_ext_id": {
				Description: "A globally unique identifier of an instance that is suitable for external consumption. This Field is Required.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"index": {
				Description: "The index on the SCSI bus to attach the VM to the Volume Group. This is an optional field.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"ext_id": {
				Description: "A globally unique identifier of a task.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func ResourceNutanixVolumeAttachVMToVolumeGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id")

	body := volumesClient.VmAttachment{}

	if vmExtID, ok := d.GetOk("vm_ext_id"); ok {
		body.ExtId = utils.StringPtr(vmExtID.(string))
	}

	if index, ok := d.GetOk("index"); ok {
		body.Index = utils.IntPtr(index.(int))
	}

	resp, err := conn.VolumeAPIInstance.AttachVm(utils.StringPtr(volumeGroupExtID.(string)), &body)
	if err != nil {
		return diag.Errorf("error while Attaching Vm to Volume Group : %v", err)
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Volume group (%s) to attach to VM: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while Attaching Vm to Volume Group: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)

	uuid := taskDetails.EntitiesAffected[0].ExtId

	d.SetId(*uuid)
	d.Set("ext_id", *uuid)

	return nil
}

func ResourceNutanixVolumeAttachVMToVolumeGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVolumeAttachVMToVolumeGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVolumeAttachVMToVolumeGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id")

	body := volumesClient.VmAttachment{}

	if vmExtID, ok := d.GetOk("vm_ext_id"); ok {
		body.ExtId = utils.StringPtr(vmExtID.(string))
	}

	if index, ok := d.GetOk("index"); ok {
		body.Index = utils.IntPtr(index.(int))
	}

	resp, err := conn.VolumeAPIInstance.DetachVm(utils.StringPtr(volumeGroupExtID.(string)), &body)
	if err != nil {
		return diag.Errorf("error while Detaching Vm to Volume Group : %v", err)
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Volume group (%s) to detach from VM: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while Detaching Vm to Volume Group: %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(taskPoll.Task)

	aJSON, _ := json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Detach Vm from Volume Group Task Details: %s", string(aJSON))

	return nil
}
