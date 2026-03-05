package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	vmmConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/policies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVMAntiAffinityPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVMAntiAffinityPolicyV2Create,
		ReadContext:   ResourceNutanixVMAntiAffinityPolicyV2Read,
		UpdateContext: ResourceNutanixVMAntiAffinityPolicyV2Update,
		DeleteContext: ResourceNutanixVMAntiAffinityPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"updated_by": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func ResourceNutanixVMAntiAffinityPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	body := policies.VmAntiAffinityPolicy{}

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(desc.(string))
	}
	if cats, ok := d.GetOk("categories"); ok {
		catStrings := common.ExpandListOfString(common.InterfaceToSlice(cats))
		body.Categories = expandPolicyCategoryReference(catStrings)
	}



	

	resp, err := conn.VMAntiAffinityPolicyAPIInstance.CreateVmAntiAffinityPolicy(&body)

	if err != nil {
		return diag.Errorf("error while creating Anti-affinity policy : %v", err)
	}

	TaskRef := resp.Data.GetValue().(vmmConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Policy to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM Anti-affinity policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching VM Anti-affinity policy UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)
	return ResourceNutanixVMAntiAffinityPolicyV2Read(ctx, d, meta)
}

func ResourceNutanixVMAntiAffinityPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.VMAntiAffinityPolicyAPIInstance.GetVmAntiAffinityPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching Anti-affinity policy : %v", err)
	}

	getResp := resp.Data.GetValue().(policies.VmAntiAffinityPolicy)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("create_time", utils.TimeStringValue(getResp.CreateTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("update_time", utils.TimeStringValue(getResp.UpdateTime)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedBy != nil {
		createdBy := make(map[string]string)
		if getResp.CreatedBy.ExtId != nil {
			createdBy["ext_id"] = *getResp.CreatedBy.ExtId
			if err := d.Set("created_by", createdBy); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if getResp.UpdatedBy != nil {
		updatedBy := make(map[string]string)
		if getResp.UpdatedBy.ExtId != nil {
			updatedBy["ext_id"] = *getResp.UpdatedBy.ExtId
			if err := d.Set("updated_by", updatedBy); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if err := d.Set("categories", flattenPolicyCategoryReference(getResp.Categories)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixVMAntiAffinityPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.VMAntiAffinityPolicyAPIInstance.GetVmAntiAffinityPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching Anti-affinity policy : %v", err)
	}

	respPolicy := resp.Data.GetValue().(policies.VmAntiAffinityPolicy)
	updateSpec := respPolicy

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("categories") {
		if cats, ok := d.GetOk("categories"); ok {
			catStrings := common.ExpandListOfString(common.InterfaceToSlice(cats))
			updateSpec.Categories = expandPolicyCategoryReference(catStrings)
		}
	}


	


	updateResp, err := conn.VMAntiAffinityPolicyAPIInstance.UpdateVmAntiAffinityPolicyById(utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error while updating Anti-affinity policy : %v", err)
	}
	TaskRef := updateResp.Data.GetValue().(vmmConfig.TaskReference)
	taskUUID := TaskRef.ExtId
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Policy to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM Anti-affinity policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixVMAntiAffinityPolicyV2Read(ctx, d, meta)
}

func ResourceNutanixVMAntiAffinityPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	readResp, err := conn.VMAntiAffinityPolicyAPIInstance.GetVmAntiAffinityPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while reading policy : %v", err)
	}

	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAntiAffinityPolicyAPIInstance.DeleteVmAntiAffinityPolicyById(utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while deleting Anti-affinity policy : %v", err)
	}
	TaskRef := resp.Data.GetValue().(vmmConfig.TaskReference)
	taskUUID := TaskRef.ExtId
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Policy to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Anti-affinity policy (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandPolicyCategoryReference(cats []string) []policies.CategoryReference {
	if len(cats) == 0 {
		return nil
	}
	catRefs := make([]policies.CategoryReference, len(cats))

	for i, cat := range cats {
		if cat != "" {
			catRefs[i] = policies.CategoryReference{
				ExtId: utils.StringPtr(cat),
			}
		}
	}
	return catRefs

}

func flattenPolicyCategoryReference(cats []policies.CategoryReference) []interface{} {
	if len(cats) > 0 {
		catList := make([]interface{}, len(cats))
		for k, v := range cats {
			if v.ExtId != nil {
				catList[k] = v.ExtId
			}
		}
		return catList
	}
	return nil
}
