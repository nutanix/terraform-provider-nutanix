package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixAddressGroupsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixAddressGroupsV2Create,
		ReadContext:   ResourceNutanixAddressGroupsV2Read,
		UpdateContext: ResourceNutanixAddressGroupsV2Update,
		DeleteContext: ResourceNutanixAddressGroupsV2Delete,
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
			},
			"ipv4_addresses": SchemaForValuePrefixLength(),
			"ip_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"end_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
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
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixAddressGroupsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	input := &import1.AddressGroup{}

	if name, ok := d.GetOk("name"); ok {
		input.Name = utils.StringPtr(name.(string))
	}

	if desc, ok := d.GetOk("description"); ok {
		input.Description = utils.StringPtr(desc.(string))
	}
	if ipv4, ok := d.GetOk("ipv4_addresses"); ok {
		input.Ipv4Addresses = expandIPv4AddressList(ipv4.([]interface{}))
	}
	if ipranges, ok := d.GetOk("ip_ranges"); ok {
		input.IpRanges = expandIPv4Range(ipranges.([]interface{}))
	}

	resp, err := conn.AddressGroupAPIInstance.CreateAddressGroup(input)
	if err != nil {
		return diag.Errorf("error while creating address groups : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Address Group to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for address groups (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching vpc UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)
	return ResourceNutanixAddressGroupsV2Read(ctx, d, meta)
}

func ResourceNutanixAddressGroupsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	resp, err := conn.AddressGroupAPIInstance.GetAddressGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching address group : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.AddressGroup)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ipv4_addresses", flattenIPv4AddressMicroSeg(getResp.Ipv4Addresses)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip_ranges", flattenIPv4Range(getResp.IpRanges)); err != nil {
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

func ResourceNutanixAddressGroupsV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	resp, err := conn.AddressGroupAPIInstance.GetAddressGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching address group : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.AddressGroup)

	updateInput := &getResp

	if d.HasChange("name") {
		updateInput.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateInput.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("ipv4_addresses") {
		updateInput.Ipv4Addresses = expandIPv4AddressList(d.Get("ipv4_addresses").([]interface{}))
	}
	if d.HasChange("ip_ranges") {
		updateInput.IpRanges = expandIPv4Range(d.Get("ip_ranges").([]interface{}))
	}

	updatedResp, err := conn.AddressGroupAPIInstance.UpdateAddressGroupById(utils.StringPtr(d.Id()), updateInput)
	if err != nil {
		return diag.Errorf("error while updating Address groups : %v", err)
	}

	TaskRef := updatedResp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Address Group to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for address groups (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixAddressGroupsV2Read(ctx, d, meta)
}

func ResourceNutanixAddressGroupsV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	resp, err := conn.AddressGroupAPIInstance.DeleteAddressGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while address service groups : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Address Group to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for address groups (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandIPv4AddressList(pr []interface{}) []config.IPv4Address {
	if len(pr) > 0 {
		ipv4s := make([]config.IPv4Address, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			ip := config.IPv4Address{}

			if v, ok := val["value"]; ok {
				if s, ok2 := v.(string); ok2 && len(s) > 0 {
					ip.Value = utils.StringPtr(s)
				}
			}

			if p, ok := val["prefix_length"]; ok {
				if n, ok2 := p.(int); ok2 {
					ip.PrefixLength = utils.IntPtr(n)
				}
			}

			ipv4s[k] = ip
		}
		return ipv4s
	}
	return nil
}

func expandIPv4Range(pr []interface{}) []import1.IPv4Range {
	if len(pr) > 0 {
		ipv4s := make([]import1.IPv4Range, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			ip := import1.IPv4Range{}

			if startPort, ok := val["start_ip"]; ok {
				ip.StartIp = utils.StringPtr(startPort.(string))
			}
			if endPort, ok := val["end_ip"]; ok {
				ip.EndIp = utils.StringPtr(endPort.(string))
			}
			ipv4s[k] = ip
		}
		return ipv4s
	}
	return nil
}
