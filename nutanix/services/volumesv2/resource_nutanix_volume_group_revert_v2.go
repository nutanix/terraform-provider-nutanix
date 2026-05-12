package volumesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	volumesPrism "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/prism/v4/config"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVolumeGroupRevertV2() *schema.Resource {
	return &schema.Resource{
		Description:   "Reverts a Volume Group identified by Volume Group external identifier. This API performs an in-place restore from a specified Volume Group recovery point.",
		CreateContext: ResourceNutanixVolumeGroupRevertV2Create,
		ReadContext:   ResourceNutanixVolumeGroupRevertV2Read,
		DeleteContext: ResourceNutanixVolumeGroupRevertV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "The external identifier of a Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"volume_group_recovery_point_ext_id": {
				Description: "The external identifier of the Volume Group recovery point. This is a mandatory field.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func ResourceNutanixVolumeGroupRevertV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	extID := d.Get("ext_id").(string)
	recoveryPointExtID := d.Get("volume_group_recovery_point_ext_id").(string)

	body := volumesClient.RevertSpec{
		VolumeGroupRecoveryPointExtId: utils.StringPtr(recoveryPointExtID),
	}

	resp, err := conn.VolumeAPIInstance.RevertVolumeGroup(utils.StringPtr(extID), &body)
	if err != nil {
		return diag.Errorf("error while reverting Volume Group : %v", err)
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for volume group revert (%s): %s", utils.StringValue(taskUUID), errWaitTask)
	}

	d.SetId(extID)
	return nil
}

func ResourceNutanixVolumeGroupRevertV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVolumeGroupRevertV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
