package networking

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixFloatingIPv4() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixFloatingIPv4Create,
		ReadContext:   ResourceNutanixFloatingIPv4Read,
		UpdateContext: ResourceNutanixFloatingIPv4Update,
		DeleteContext: ResourceNutanixFloatingIPv4Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"association": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm_nic_association": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_nic_reference": {
										Type:     schema.TypeString,
										Required: true,
									},
									"vpc_reference": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"private_ip_association": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_reference": {
										Type:     schema.TypeString,
										Required: true,
									},
									"private_ip": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": SchemaForValuePrefixLength(),
												"ipv6": SchemaForValuePrefixLength(),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"floating_ip": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(),
						"ipv6": SchemaForValuePrefixLength(),
					},
				},
			},
			"external_subnet_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_subnet": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     DataSourceNutanixSubnetv4(),
			},
			"vpc_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vm_nic_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"load_balancer_session_reference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DataSourceVPCSchemaV4(),
				},
			},
			"vm_nic": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"floating_ip_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"association_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DatasourceMetadataSchemaV4(),
				},
			},
		},
	}
}

func ResourceNutanixFloatingIPv4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	inputSpec := import1.FloatingIp{}
	fipName := ""

	if name, ok := d.GetOk("name"); ok {
		inputSpec.Name = utils.StringPtr(name.(string))
		fipName = name.(string)
	}

	if desc, ok := d.GetOk("description"); ok {
		inputSpec.Description = utils.StringPtr(desc.(string))
	}
	if association, ok := d.GetOk("association"); ok {
		inputSpec.Association = expandOneOfFloatingIPAssociation(association)
	}
	if fip, ok := d.GetOk("floating_ip"); ok {
		inputSpec.FloatingIp = expandFloatingIPAddress(fip)
	}
	if extSubRef, ok := d.GetOk("external_subnet_reference"); ok {
		inputSpec.ExternalSubnetReference = utils.StringPtr(extSubRef.(string))
	}
	if extSub, ok := d.GetOk("external_subnet"); ok {
		inputSpec.ExternalSubnet = expandSubnet(extSub)
	}
	if vpcRef, ok := d.GetOk("vpc_reference"); ok {
		inputSpec.VpcReference = utils.StringPtr(vpcRef.(string))
	}
	if vmNICRef, ok := d.GetOk("vm_nic_reference"); ok {
		inputSpec.VmNicReference = utils.StringPtr(vmNICRef.(string))
	}
	if vpc, ok := d.GetOk("vpc"); ok {
		inputSpec.Vpc = expandVpc(vpc)
	}
	if vmNic, ok := d.GetOk("vm_nic"); ok {
		inputSpec.VmNic = expandVMNic(vmNic)
	}

	resp, err := conn.FloatingIPAPIInstance.CreateFloatingIp(&inputSpec)
	if err != nil {
		return diag.Errorf("error while creating floating IPs : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the FileServer to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for floating IP (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	filter := fmt.Sprintf("name eq  '%s'", fipName)
	readResp, err := conn.FloatingIPAPIInstance.ListFloatingIps(nil, nil, &filter, nil, nil, nil)
	if err != nil {
		return diag.Errorf("error while fetching fips : %v", err)
	}

	getAllFipResp := readResp.Data.GetValue().([]import1.FloatingIp)

	d.SetId(*getAllFipResp[0].ExtId)
	return ResourceNutanixFloatingIPv4Read(ctx, d, meta)
}

func ResourceNutanixFloatingIPv4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.FloatingIPAPIInstance.GetFloatingIpById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching floating ips : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.FloatingIp)
	fmt.Println(getResp)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("association", flattenAssociation(getResp.Association)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("floating_ip", flattenFloatingIP(getResp.FloatingIp)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("external_subnet_reference", getResp.ExternalSubnetReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("external_subnet", flattenExternalSubnet(getResp.ExternalSubnet)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("private_ip", getResp.PrivateIp); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("floating_ip_value", getResp.FloatingIpValue); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("association_status", getResp.AssociationStatus); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vpc_reference", getResp.VpcReference); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vm_nic_reference", getResp.VmNicReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc", flattenVpc(getResp.Vpc)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_nic", flattenVMNic(getResp.VmNic)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixFloatingIPv4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.FloatingIPAPIInstance.GetFloatingIpById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching floating ips : %v", err)
	}

	respFloatingIP := resp.Data.GetValue().(import1.FloatingIp)

	updateSpec := respFloatingIP

	// Extract E-Tag Header
	// etagValue := ApiClientInstance.GetEtag(getResp)

	// args := make(map[string]interface{})
	// args["If-Match"] = etagValue

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("association") {
		updateSpec.Association = expandOneOfFloatingIPAssociation(d.Get("association"))
	}
	if d.HasChange("floating_ip") {
		updateSpec.FloatingIp = expandFloatingIPAddress(d.Get("floating_ip"))
	}
	if d.HasChange("external_subnet_reference") {
		updateSpec.ExternalSubnetReference = utils.StringPtr(d.Get("external_subnet_reference").(string))
	}
	if d.HasChange("external_subnet") {
		updateSpec.ExternalSubnet = expandSubnet(d.Get("external_subnet"))
	}
	if d.HasChange("vpc_reference") {
		updateSpec.VpcReference = utils.StringPtr(d.Get("vpc_reference").(string))
	}
	if d.HasChange("vm_nic_reference") {
		updateSpec.VmNicReference = utils.StringPtr(d.Get("vm_nic_reference").(string))
	}
	if d.HasChange("vpc") {
		updateSpec.Vpc = expandVpc(d.Get("vpc"))
	}
	if d.HasChange("vm_nic") {
		updateSpec.VmNic = expandVMNic(d.Get("vm_nic"))
	}

	getResp, err := conn.FloatingIPAPIInstance.UpdateFloatingIpById(utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.FromErr(err)
	}
	TaskRef := getResp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the FileServer to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for floating IP (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func ResourceNutanixFloatingIPv4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.FloatingIPAPIInstance.DeleteFloatingIpById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting floating ip : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the FileServer to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for floating IP (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandFloatingIPAddress(pr interface{}) *import1.FloatingIPAddress {
	if pr != nil {
		fip := &import1.FloatingIPAddress{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if ipv4, ok := val["ipv4"]; ok {
			fip.Ipv4 = expandFloatingIPv4Address(ipv4)
		}
		if ipv6, ok := val["ipv6"]; ok {
			fip.Ipv6 = expandFloatingIPv6Address(ipv6)
		}

		return fip
	}
	return nil
}

func expandFloatingIPv4Address(pr interface{}) *import1.FloatingIPv4Address {
	if pr != nil {
		ipv4 := &import1.FloatingIPv4Address{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv4.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv4.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv4
	}
	return nil
}

func expandFloatingIPv6Address(pr interface{}) *import1.FloatingIPv6Address {
	if pr != nil {
		ipv6 := &import1.FloatingIPv6Address{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv6.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv6.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv6
	}
	return nil
}

func expandSubnet(pr interface{}) *import1.Subnet {
	if pr != nil {
		sub := &import1.Subnet{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if name, ok := val["name"]; ok {
			sub.Name = utils.StringPtr(name.(string))
		}
		if desc, ok := val["description"]; ok {
			sub.Description = utils.StringPtr(desc.(string))
		}
		if subType, ok := val["subnet_type"]; ok {
			subMap := map[string]interface{}{
				"OVERLAY": "2",
				"VLAN":    "3",
			}
			pInt := subMap[subType.(string)]
			p := import1.SubnetType(pInt.(int))
			sub.SubnetType = &p
		}
		if dhcp, ok := val["dhcp_options"]; ok {
			sub.DhcpOptions = expandDhcpOptions(dhcp.([]interface{}))
		}
		if clsRef, ok := val["cluster_reference"]; ok {
			sub.ClusterReference = utils.StringPtr(clsRef.(string))
		}
		if vsRef, ok := val["virtual_switch_reference"]; ok {
			sub.VirtualSwitchReference = utils.StringPtr(vsRef.(string))
		}
		if vpcRef, ok := val["vpc_reference"]; ok {
			sub.VirtualSwitchReference = utils.StringPtr(vpcRef.(string))
		}
		if isNat, ok := val["is_nat_enabled"]; ok {
			sub.IsNatEnabled = utils.BoolPtr(isNat.(bool))
		}
		if isExt, ok := val["is_external"]; ok {
			sub.IsExternal = utils.BoolPtr(isExt.(bool))
		}
		if reservedIPAdd, ok := val["reserved_ip_addresses"]; ok {
			sub.ReservedIpAddresses = expandIPAddress(reservedIPAdd.([]interface{}))
		}
		if dynamicIPAdd, ok := val["dynamic_ip_addresses"]; ok {
			sub.DynamicIpAddresses = expandIPAddress(dynamicIPAdd.([]interface{}))
		}
		if ntwfuncRef, ok := val["network_function_chain_reference"]; ok {
			sub.NetworkFunctionChainReference = utils.StringPtr(ntwfuncRef.(string))
		}
		if bridgeName, ok := val["bridge_name"]; ok {
			sub.BridgeName = utils.StringPtr(bridgeName.(string))
		}
		if isAdvNet, ok := val["is_advanced_networking"]; ok {
			sub.IsAdvancedNetworking = utils.BoolPtr(isAdvNet.(bool))
		}
		if clsName, ok := val["cluster_name"]; ok {
			sub.ClusterName = utils.StringPtr(clsName.(string))
		}
		if hypervisorType, ok := val["hypervisor_type"]; ok {
			sub.HypervisorType = utils.StringPtr(hypervisorType.(string))
		}
		if vswitch, ok := val["virtual_switch"]; ok {
			sub.VirtualSwitch = expandVirtualSwitch(vswitch)
		}
		if vpc, ok := val["vpc"]; ok {
			sub.Vpc = expandVpc(vpc)
		}
		if ipPrefix, ok := val["ip_prefix"]; ok {
			sub.IpPrefix = utils.StringPtr(ipPrefix.(string))
		}
		if ipUsage, ok := val["ip_usage"]; ok {
			sub.IpUsage = exapndIPUsage(ipUsage)
		}
		return sub
	}
	return nil
}

func expandOneOfFloatingIPAssociation(pr interface{}) *import1.OneOfFloatingIpAssociation {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		fip := &import1.OneOfFloatingIpAssociation{}

		if vmNic, ok := val["vm_nic_association"]; ok && len(vmNic.([]interface{})) > 0 {
			nic := import1.NewVmNicAssociation()
			prI := vmNic.([]interface{})
			val := prI[0].(map[string]interface{})

			if vmNicRef, ok := val["vm_nic_reference"]; ok {
				nic.VmNicReference = utils.StringPtr(vmNicRef.(string))
			}
			fip.SetValue(*nic)
		}

		if privateIP, ok := val["private_ip_association"]; ok && len(privateIP.([]interface{})) > 0 {
			pip := import1.NewPrivateIpAssociation()
			prI := privateIP.([]interface{})
			val := prI[0].(map[string]interface{})

			if vpcRef, ok := val["vpc_reference"]; ok && len(vpcRef.(string)) > 0 {
				pip.VpcReference = utils.StringPtr(vpcRef.(string))
			}
			if pIP, ok := val["private_ip"]; ok && len(pIP.([]interface{})) > 0 {
				pip.PrivateIp = expandIPAddressMap(pIP)
			}
			fip.SetValue(*pip)
		}
		return fip
	}
	return nil
}

func expandVMNic(pr interface{}) *import1.VmNic {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		nics := &import1.VmNic{}

		if privateIP, ok := val["private_ip"]; ok {
			nics.PrivateIp = utils.StringPtr(privateIP.(string))
		}
		return nics
	}
	return nil
}

func taskStateRefreshPrismTaskGroupFunc(ctx context.Context, client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vresp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID))

		if err != nil {
			return "", "", (fmt.Errorf("error while polling prism task: %v", err))
		}

		// get the group results

		v := vresp.Data.GetValue().(import2.Task)

		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
			return v, getTaskStatus(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStatus(v.Status), nil
	}
}

func getTaskStatus(pr *import2.TaskStatus) string {
	if pr != nil {
		const two, three, four, five, six, seven = 2, 3, 4, 5, 6, 7
		if *pr == import2.TaskStatus(two) {
			return "QUEUED"
		}
		if *pr == import2.TaskStatus(three) {
			return "RUNNING"
		}
		if *pr == import2.TaskStatus(four) {
			return "CANCELING"
		}
		if *pr == import2.TaskStatus(five) {
			return "SUCCEEDED"
		}
		if *pr == import2.TaskStatus(six) {
			return "FAILED"
		}
		if *pr == import2.TaskStatus(seven) {
			return "CANCELED"
		}
	}
	return "UNKNOWN"
}
