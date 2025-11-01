package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixServiceGroupsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixServiceGroupsV2Create,
		ReadContext:   ResourceNutanixServiceGroupsV2Read,
		UpdateContext: ResourceNutanixServiceGroupsV2Update,
		DeleteContext: ResourceNutanixServiceGroupsV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tcp_services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"end_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"udp_services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"end_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"icmp_services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_all_allowed": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"code": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"policy_references": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"is_system_defined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixServiceGroupsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	spec := &import1.ServiceGroup{}

	if name, ok := d.GetOk("name"); ok {
		spec.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		spec.Description = utils.StringPtr(desc.(string))
	}
	if tcp, ok := d.GetOk("tcp_services"); ok {
		spec.TcpServices = expandTCPPortRangeSpec(tcp.([]interface{}))
	}
	if udp, ok := d.GetOk("udp_services"); ok {
		spec.UdpServices = expandUDPPortRangeSpec(udp.([]interface{}))
	}
	if icmp, ok := d.GetOk("icmp_services"); ok {
		spec.IcmpServices = expandIcmpTypeCodeSpec(icmp.([]interface{}))
	}

	resp, err := conn.ServiceGroupAPIInstance.CreateServiceGroup(spec)
	if err != nil {
		return diag.Errorf("error while creating service groups : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Service Group to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for service groups (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching vpc UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)

	return ResourceNutanixServiceGroupsV2Read(ctx, d, meta)
}

func ResourceNutanixServiceGroupsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	resp, err := conn.ServiceGroupAPIInstance.GetServiceGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching service groups : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.ServiceGroup)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_system_defined", getResp.IsSystemDefined); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tcp_services", flattenTCPPortRangeSpec(getResp.TcpServices)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("udp_services", flattenUDPPortRangeSpec(getResp.UdpServices)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("icmp_services", flattenIcmpTypeCodeSpec(getResp.IcmpServices)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("policy_references", flattenListofString(getResp.PolicyReferences)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinksMicroSeg(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixServiceGroupsV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI
	updatedSpec := import1.ServiceGroup{}

	resp, err := conn.ServiceGroupAPIInstance.GetServiceGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching service groups : %v", err)
	}

	//Extract E-Tag Header
	etagValue := conn.ServiceGroupAPIInstance.ApiClient.GetEtag(resp)

	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	updatedSpec = resp.Data.GetValue().(import1.ServiceGroup)

	if d.HasChange("name") {
		updatedSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updatedSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("tcp_services") {
		updatedSpec.TcpServices = expandTCPPortRangeSpec(d.Get("tcp_services").([]interface{}))
	}
	if d.HasChange("udp_services") {
		updatedSpec.UdpServices = expandUDPPortRangeSpec(d.Get("udp_services").([]interface{}))
	}
	if d.HasChange("icmp_services") {
		updatedSpec.IcmpServices = expandIcmpTypeCodeSpec(d.Get("icmp_services").([]interface{}))
	}

	// removing read only attribute from spec
	updatedSpec.IsSystemDefined = nil

	updatedResp, err := conn.ServiceGroupAPIInstance.UpdateServiceGroupById(utils.StringPtr(d.Id()), &updatedSpec, args)
	if err != nil {
		return diag.Errorf("error while updating service groups : %v", err)
	}

	TaskRef := updatedResp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Service Group to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for service groups (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixServiceGroupsV2Read(ctx, d, meta)
}

func ResourceNutanixServiceGroupsV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	resp, err := conn.ServiceGroupAPIInstance.DeleteServiceGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting service groups : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Service Group to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for service groups (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandTCPPortRangeSpec(pr []interface{}) []import1.TcpPortRangeSpec {
	if len(pr) > 0 {
		tcps := make([]import1.TcpPortRangeSpec, len(pr))

		for k, v := range pr {
			tcp := import1.TcpPortRangeSpec{}
			val := v.(map[string]interface{})

			if start, ok := val["start_port"]; ok {
				tcp.StartPort = utils.IntPtr(start.(int))
			}
			if end, ok := val["end_port"]; ok {
				tcp.EndPort = utils.IntPtr(end.(int))
			}
			tcps[k] = tcp
		}
		return tcps
	}
	return nil
}

func expandUDPPortRangeSpec(pr []interface{}) []import1.UdpPortRangeSpec {
	if len(pr) > 0 {
		udps := make([]import1.UdpPortRangeSpec, len(pr))

		for k, v := range pr {
			udp := import1.UdpPortRangeSpec{}
			val := v.(map[string]interface{})

			if start, ok := val["start_port"]; ok {
				udp.StartPort = utils.IntPtr(start.(int))
			}
			if end, ok := val["end_port"]; ok {
				udp.EndPort = utils.IntPtr(end.(int))
			}
			udps[k] = udp
		}
		return udps
	}
	return nil
}

func expandIcmpTypeCodeSpec(pr []interface{}) []import1.IcmpTypeCodeSpec {
	if len(pr) > 0 {
		icmps := make([]import1.IcmpTypeCodeSpec, len(pr))

		for k, v := range pr {
			icmp := import1.IcmpTypeCodeSpec{}
			val := v.(map[string]interface{})

			if allAllow, ok := val["is_all_allowed"]; ok {
				icmp.IsAllAllowed = utils.BoolPtr(allAllow.(bool))
			}
			if code, ok := val["code"]; ok {
				icmp.Code = utils.IntPtr(code.(int))
			}
			if types, ok := val["type"]; ok {
				icmp.Type = utils.IntPtr(types.(int))
			}
			icmps[k] = icmp
		}
		return icmps
	}
	return nil
}
