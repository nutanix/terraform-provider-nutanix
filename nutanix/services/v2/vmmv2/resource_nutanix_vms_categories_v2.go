package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/prism/v4/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVmsCategoriesV4() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVmsCategoriesV4Create,
		ReadContext:   ResourceNutanixVmsCategoriesV4Read,
		UpdateContext: ResourceNutanixVmsCategoriesV4Update,
		DeleteContext: ResourceNutanixVmsCategoriesV4Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"categories": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixVmsCategoriesV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("ext_id")
	body := config.AssociateVmCategoriesParams{}

	if category, ok := d.GetOk("categories"); ok {
		body.Categories = expandCategoryReference(category.([]interface{}))
	}

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.AssociateCategories(utils.StringPtr(vmExtID.(string)), &body, args)
	if err != nil {
		return diag.Errorf("error while associating categories : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for categories (%s) to attach: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	d.SetId(resource.UniqueId())
	return nil
}

func ResourceNutanixVmsCategoriesV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVmsCategoriesV4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("ext_id")

	if d.HasChange("categories") {
		old, new := d.GetChange("categories")

		oldCats := expandCategoryReference(old.([]interface{}))
		newCats := expandCategoryReference(new.([]interface{}))

		// check if categories already exists in new
		deleteCats := config.DisassociateVmCategoriesParams{}
		addCats := config.AssociateVmCategoriesParams{}

		oldcatMap := make(map[string]interface{})
		newcatMap := make(map[string]interface{})

		for _, ov := range oldCats {
			oldcatMap[*ov.ExtId] = ov
		}

		for _, nv := range newCats {
			newcatMap[*nv.ExtId] = nv
		}

		for _, ov := range oldCats {
			if _, exists := newcatMap[*ov.ExtId]; !exists {
				deleteCats.Categories = append(deleteCats.Categories, ov)
			}
		}

		for _, nv := range newCats {
			if _, exists := oldcatMap[*nv.ExtId]; !exists {
				addCats.Categories = append(addCats.Categories, nv)
			}
		}

		if len(deleteCats.Categories) > 0 {
			readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
			if err != nil {
				return diag.Errorf("error while reading vm : %v", err)
			}
			// Extract E-Tag Header
			args := make(map[string]interface{})
			args["If-Match"] = getEtagHeader(readResp, conn)

			resp, err := conn.VMAPIInstance.DisassociateCategories(utils.StringPtr(vmExtID.(string)), &deleteCats, args)
			if err != nil {
				return diag.Errorf("error while diassociate categories : %v", err)
			}

			TaskRef := resp.Data.GetValue().(import1.TaskReference)
			taskUUID := TaskRef.ExtId

			taskconn := meta.(*conns.Client).PrismAPI
			// Wait for the VM to be available
			stateConf := &resource.StateChangeConf{
				Pending: []string{"QUEUED", "RUNNING"},
				Target:  []string{"SUCCEEDED"},
				Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
				Timeout: d.Timeout(schema.TimeoutCreate),
			}

			if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
				return diag.Errorf("error waiting for categories (%s) to diassociate: %s", utils.StringValue(taskUUID), errWaitTask)
			}
		}
		if len(addCats.Categories) > 0 {
			readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
			if err != nil {
				return diag.Errorf("error while reading vm : %v", err)
			}
			// Extract E-Tag Header
			args := make(map[string]interface{})
			args["If-Match"] = getEtagHeader(readResp, conn)

			resp, err := conn.VMAPIInstance.AssociateCategories(utils.StringPtr(vmExtID.(string)), &addCats, args)
			if err != nil {
				return diag.Errorf("error while associating categories : %v", err)
			}

			TaskRef := resp.Data.GetValue().(import1.TaskReference)
			taskUUID := TaskRef.ExtId

			taskconn := meta.(*conns.Client).PrismAPI
			// Wait for the VM to be available
			stateConf := &resource.StateChangeConf{
				Pending: []string{"QUEUED", "RUNNING"},
				Target:  []string{"SUCCEEDED"},
				Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
				Timeout: d.Timeout(schema.TimeoutCreate),
			}

			if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
				return diag.Errorf("error waiting for categories (%s) to attach: %s", utils.StringValue(taskUUID), errWaitTask)
			}
		}
	}
	return nil
}

func ResourceNutanixVmsCategoriesV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("ext_id")
	body := config.DisassociateVmCategoriesParams{}

	if category, ok := d.GetOk("categories"); ok {
		body.Categories = expandCategoryReference(category.([]interface{}))
	}

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.DisassociateCategories(utils.StringPtr(vmExtID.(string)), &body, args)
	if err != nil {
		return diag.Errorf("error while diassociate categories : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for categories (%s) to diassociate: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}
