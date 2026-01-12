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
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVMHostAffinityPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVMHostAffinityPolicyV2Create,
		ReadContext:   ResourceNutanixVMHostAffinityPolicyV2Read,
		UpdateContext: ResourceNutanixVMHostAffinityPolicyV2Update,
		DeleteContext: ResourceNutanixVMHostAffinityPolicyV2Delete,
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
			"last_updated_by": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vm_categories": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_categories": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func ResourceNutanixVMHostAffinityPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	body := policies.VmHostAffinityPolicy{}

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(desc.(string))
	}
	if vmCats, ok := d.GetOk("vm_categories"); ok {
		if vmCatSet, ok := vmCats.(*schema.Set); ok {
			body.VmCategories = expandPolicyCategoryReference(vmCatSet.List())
		} else if vmCatList, ok := vmCats.([]interface{}); ok {
			body.VmCategories = expandPolicyCategoryReference(vmCatList)
		}
	}
	if hostCats, ok := d.GetOk("host_categories"); ok {
		if hostCatSet, ok := hostCats.(*schema.Set); ok {
			body.HostCategories = expandPolicyCategoryReference(hostCatSet.List())
		} else if hostCatList, ok := hostCats.([]interface{}); ok {
			body.HostCategories = expandPolicyCategoryReference(hostCatList)
		}
	}

	resp, err := conn.VMHostAffinityPolicyAPIInstance.CreateVmHostAffinityPolicy(&body)

	if err != nil {
		return diag.Errorf("error while creating VM-Host Affinity policy : %v", err)
	}

	TaskRef := resp.Data.GetValue().(vmmConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Policy to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM-Host Affinity policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching VM-Host Affinity policy UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)
	return ResourceNutanixVMHostAffinityPolicyV2Read(ctx, d, meta)
}

func ResourceNutanixVMHostAffinityPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.VMHostAffinityPolicyAPIInstance.GetVmHostAffinityPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching VM-Host Affinity policy : %v", err)
	}

	getResp := resp.Data.GetValue().(policies.VmHostAffinityPolicy)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreateTime != nil {
		t := getResp.CreateTime
		if err := d.Set("create_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.UpdateTime != nil {
		t := getResp.UpdateTime
		if err := d.Set("update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
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
	if getResp.LastUpdatedBy != nil {
		updatedBy := make(map[string]string)
		if getResp.LastUpdatedBy.ExtId != nil {
			updatedBy["ext_id"] = *getResp.LastUpdatedBy.ExtId
			if err := d.Set("last_updated_by", updatedBy); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if err := d.Set("vm_categories", flattenPolicyCategoryReference(getResp.VmCategories)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("host_categories", flattenPolicyCategoryReference(getResp.HostCategories)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixVMHostAffinityPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.VMHostAffinityPolicyAPIInstance.GetVmHostAffinityPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching VM-Host Affinity policy : %v", err)
	}

	respPolicy := resp.Data.GetValue().(policies.VmHostAffinityPolicy)
	updateSpec := respPolicy

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("vm_categories") {
		if vmCats := d.Get("vm_categories"); vmCats != nil {
			if vmCatSet, ok := vmCats.(*schema.Set); ok {
				updateSpec.VmCategories = expandPolicyCategoryReference(vmCatSet.List())
			} else if vmCatList, ok := vmCats.([]interface{}); ok {
				updateSpec.VmCategories = expandPolicyCategoryReference(vmCatList)
			}
		}
	}
	if d.HasChange("host_categories") {
		if hostCats := d.Get("host_categories"); hostCats != nil {
			if hostCatSet, ok := hostCats.(*schema.Set); ok {
				updateSpec.HostCategories = expandPolicyCategoryReference(hostCatSet.List())
			} else if hostCatList, ok := hostCats.([]interface{}); ok {
				updateSpec.HostCategories = expandPolicyCategoryReference(hostCatList)
			}
		}
	}

	updateResp, err := conn.VMHostAffinityPolicyAPIInstance.UpdateVmHostAffinityPolicyById(utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error while updating VM-Host Affinity policy : %v", err)
	}
	TaskRef := updateResp.Data.GetValue().(vmmConfig.TaskReference)
	taskUUID := TaskRef.ExtId
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Policy to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM-Host Affinity policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixVMHostAffinityPolicyV2Read(ctx, d, meta)
}

func ResourceNutanixVMHostAffinityPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	readResp, err := conn.VMHostAffinityPolicyAPIInstance.GetVmHostAffinityPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while reading policy : %v", err)
	}

	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMHostAffinityPolicyAPIInstance.DeleteVmHostAffinityPolicyById(utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while deleting VM-Host Affinity policy : %v", err)
	}
	TaskRef := resp.Data.GetValue().(vmmConfig.TaskReference)
	taskUUID := TaskRef.ExtId
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Policy to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM-Host Affinity policy (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}
