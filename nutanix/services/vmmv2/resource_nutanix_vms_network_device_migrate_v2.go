package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVmsNetworkDeviceMigrateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVmsNetworkDeviceMigrateV2Create,
		ReadContext:   ResourceNutanixVmsNetworkDeviceMigrateV2Read,
		UpdateContext: ResourceNutanixVmsNetworkDeviceMigrateV2Update,
		DeleteContext: ResourceNutanixVmsNetworkDeviceMigrateV2Delete,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"migrate_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ASSIGN_IP", "RELEASE_IP"}, false),
			},
			"ip_address": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prefix_length": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  defaultValue,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixVmsNetworkDeviceMigrateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	extID := d.Get("ext_id")
	body := config.MigrateNicConfig{}

	if subnet, ok := d.GetOk("subnet"); ok {
		body.Subnet = expandSubnetReference(subnet)
	}
	if migrateType, ok := d.GetOk("migrate_type"); ok && len(migrateType.(string)) > 0 {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"ASSIGN_IP":  two,
			"RELEASE_IP": three,
		}
		pVal := subMap[migrateType.(string)]
		p := config.MigrateNicType(pVal.(int))
		body.MigrateType = &p
	}
	if ipAddress, ok := d.GetOk("ip_address"); ok {
		body.IpAddress = expandIPv4Address(ipAddress)
	}

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.MigrateNicById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(extID.(string)), &body, args)
	if err != nil {
		return diag.Errorf("error while migrate nic : %v", err)
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
		return diag.Errorf("error waiting for nic (%s) to migrate: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	d.SetId(*taskUUID)
	return nil
}

func ResourceNutanixVmsNetworkDeviceMigrateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVmsNetworkDeviceMigrateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	extID := d.Get("ext_id")
	body := config.MigrateNicConfig{}

	if subnet, ok := d.GetOk("subnet"); ok {
		body.Subnet = expandSubnetReference(subnet)
	}
	if migrateType, ok := d.GetOk("migrate_type"); ok && len(migrateType.(string)) > 0 {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"ASSIGN_IP":  two,
			"RELEASE_IP": three,
		}
		pVal := subMap[migrateType.(string)]
		p := config.MigrateNicType(pVal.(int))
		body.MigrateType = &p
	}

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.MigrateNicById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(extID.(string)), &body, args)
	if err != nil {
		return diag.Errorf("error while migrate nic : %v", err)
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
		return diag.Errorf("error waiting for nic (%s) to migrate: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}

func ResourceNutanixVmsNetworkDeviceMigrateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
