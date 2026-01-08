package networkingv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixSubnetV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixSubnetV2Create,
		ReadContext:   ResourceNutanixSubnetV2Read,
		UpdateContext: ResourceNutanixSubnetV2Update,
		DeleteContext: ResourceNutanixSubnetV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Optional: true,
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
			"subnet_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OVERLAY", "VLAN"}, false),
			},
			"network_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dhcp_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name_servers": {
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
						"domain_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"search_domains": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tftp_server_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"boot_file_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ntp_servers": {
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
					},
				},
			},
			"ip_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_subnet": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": SchemaForValuePrefixLength(),
												"prefix_length": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"default_gateway_ip":  SchemaForValuePrefixLength(),
									"dhcp_server_address": SchemaForValuePrefixLength(),
									"pool_list": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"end_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"ipv6": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_subnet": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": SchemaForValuePrefixLength(),
												"prefix_length": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"default_gateway_ip":  SchemaForValuePrefixLength(),
									"dhcp_server_address": SchemaForValuePrefixLength(),
									"pool_list": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"end_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"cluster_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"virtual_switch_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_nat_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_external": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"reserved_ip_addresses": SchemaForValuePrefixLength(),
			"dynamic_ip_addresses": {
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
			"network_function_chain_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bridge_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_advanced_networking": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hypervisor_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"virtual_switch": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DataSourceVirtualSwitchSchemaV2(),
				},
			},
			"vpc": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DataSourceVPCSchemaV2(),
				},
			},
			"ip_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_usage": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"num_macs": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"num_free_ips": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"num_assigned_ips": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"ip_pool_usages": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_free_ips": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"num_total_ips": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"range": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"end_ip": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"migration_state": {
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
							Optional: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixSubnetV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	inputSpec := import1.Subnet{}
	if name, nok := d.GetOk("name"); nok {
		inputSpec.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		inputSpec.Description = utils.StringPtr(desc.(string))
	}
	if subType, ok := d.GetOk("subnet_type"); ok {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"OVERLAY": two,
			"VLAN":    three,
		}
		pVal := subMap[subType.(string)]

		p := import1.SubnetType(pVal.(int))
		inputSpec.SubnetType = &p
	}

	if networkID, ok := d.GetOk("network_id"); ok {
		inputSpec.NetworkId = utils.IntPtr(networkID.(int))
	}

	if dhcp, ok := d.GetOk("dhcp_options"); ok {
		inputSpec.DhcpOptions = expandDhcpOptions(dhcp.([]interface{}))
	}
	if clsRef, ok := d.GetOk("cluster_reference"); ok {
		inputSpec.ClusterReference = utils.StringPtr(clsRef.(string))
	}
	if vsRef, ok := d.GetOk("virtual_switch_reference"); ok {
		inputSpec.VirtualSwitchReference = utils.StringPtr(vsRef.(string))
	}
	if vpcRef, ok := d.GetOk("vpc_reference"); ok {
		inputSpec.VpcReference = utils.StringPtr(vpcRef.(string))
	}
	if common.IsExplicitlySet(d, "is_nat_enabled") {
		isNat := d.Get("is_nat_enabled").(bool)
		inputSpec.IsNatEnabled = utils.BoolPtr(isNat)
	}
	if common.IsExplicitlySet(d, "is_external") {
		isExt := d.Get("is_external").(bool)
		inputSpec.IsExternal = utils.BoolPtr(isExt)
	}
	if reservedIPAdd, ok := d.GetOk("reserved_ip_addresses"); ok {
		inputSpec.ReservedIpAddresses = expandIPAddress(reservedIPAdd.([]interface{}))
	}
	if dynamicIPAdd, ok := d.GetOk("dynamic_ip_addresses"); ok {
		inputSpec.DynamicIpAddresses = expandIPAddress(dynamicIPAdd.([]interface{}))
	}
	if ntwfuncRef, ok := d.GetOk("network_function_chain_reference"); ok {
		inputSpec.NetworkFunctionChainReference = utils.StringPtr(ntwfuncRef.(string))
	}
	if bridgeName, ok := d.GetOk("bridge_name"); ok {
		inputSpec.BridgeName = utils.StringPtr(bridgeName.(string))
	}
	if isAdvNet, ok := d.GetOk("is_advanced_networking"); ok {
		inputSpec.IsAdvancedNetworking = utils.BoolPtr(isAdvNet.(bool))
	}
	if clsName, ok := d.GetOk("cluster_name"); ok {
		inputSpec.ClusterName = utils.StringPtr(clsName.(string))
	}
	if hypervisorType, ok := d.GetOk("hypervisor_type"); ok {
		inputSpec.HypervisorType = utils.StringPtr(hypervisorType.(string))
	}
	if vswitch, ok := d.GetOk("virtual_switch"); ok {
		inputSpec.VirtualSwitch = expandVirtualSwitch(vswitch)
	}
	if vpc, ok := d.GetOk("vpc"); ok {
		inputSpec.Vpc = expandVpc(vpc)
	}
	if ipPrefix, ok := d.GetOk("ip_prefix"); ok {
		inputSpec.IpPrefix = utils.StringPtr(ipPrefix.(string))
	}
	if ipUsage, ok := d.GetOk("ip_usage"); ok {
		inputSpec.IpUsage = expandIPUsage(ipUsage)
	}
	if ipConfig, ok := d.GetOk("ip_config"); ok {
		inputSpec.IpConfig = expandIPConfig(ipConfig.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(inputSpec, "", " ")
	log.Printf("[DEBUG] Subnet create payload : %s", string(aJSON))

	resp, err := conn.SubnetAPIInstance.CreateSubnet(&inputSpec)
	if err != nil {
		return diag.Errorf("error while creating subnets : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the subnet to be created
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for subnet (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching subnet task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Subnet Task Details: %s", string(aJSON))

	var subnetExtID *string
	for _, entity := range taskDetails.EntitiesAffected {
		if utils.StringValue(entity.Rel) == utils.RelEntityTypeSubnet {
			subnetExtID = entity.ExtId
			break
		}
	}
	if subnetExtID != nil {
		d.SetId(utils.StringValue(subnetExtID))
	} else {
		return diag.Errorf("error while fetching subnet ExtId: subnet entity not found in EntitiesAffected")
	}
	return ResourceNutanixSubnetV2Read(ctx, d, meta)
}

func ResourceNutanixSubnetV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.SubnetAPIInstance.GetSubnetById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching subnets : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Subnet)

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
	if err := d.Set("subnet_type", flattenSubnetType(getResp.SubnetType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network_id", getResp.NetworkId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dhcp_options", flattenDhcpOptions(getResp.DhcpOptions)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ip_config", flattenIPConfig(getResp.IpConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_reference", getResp.ClusterReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("virtual_switch_reference", getResp.VirtualSwitchReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc_reference", getResp.VpcReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_nat_enabled", getResp.IsNatEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_external", getResp.IsExternal); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("reserved_ip_addresses", flattenReservedIPAddresses(getResp.ReservedIpAddresses)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dynamic_ip_addresses", flattenReservedIPAddresses(getResp.DynamicIpAddresses)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network_function_chain_reference", getResp.NetworkFunctionChainReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("bridge_name", getResp.BridgeName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_advanced_networking", getResp.IsAdvancedNetworking); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_name", getResp.ClusterName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hypervisor_type", getResp.HypervisorType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("virtual_switch", flattenVirtualSwitch(getResp.VirtualSwitch)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc", flattenVPC(getResp.Vpc)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ip_prefix", getResp.IpPrefix); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ip_usage", flattenIPUsage(getResp.IpUsage)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("migration_state", flattenMigrationState(getResp.MigrationState)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixSubnetV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI
	updateSpec := import1.Subnet{}

	readResp, err := conn.SubnetAPIInstance.GetSubnetById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching subnets : %v", err)
	}

	updateSpec = readResp.Data.GetValue().(import1.Subnet)
	// Extract E-Tag Header
	etagValue := conn.SubnetAPIInstance.ApiClient.GetEtag(readResp)

	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("subnet_type") {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"OVERLAY": two,
			"VLAN":    three,
		}
		pInt := subMap[d.Get("subnet_type").(string)]
		p := import1.SubnetType(pInt.(int))
		updateSpec.SubnetType = &p
	}
	if d.HasChange("dhcp_options") {
		updateSpec.DhcpOptions = expandDhcpOptions(d.Get("dhcp_options").([]interface{}))
	} else {
		updateSpec.DhcpOptions = nil
	}
	if d.HasChange("cluster_reference") {
		updateSpec.ClusterReference = utils.StringPtr(d.Get("cluster_reference").(string))
	}
	if d.HasChange("virtual_switch_reference") {
		updateSpec.VirtualSwitchReference = utils.StringPtr(d.Get("virtual_switch_reference").(string))
	}
	if d.HasChange("vpc_reference") {
		updateSpec.VpcReference = utils.StringPtr(d.Get("vpc_reference").(string))
	}
	if d.HasChange("is_nat_enabled") {
		updateSpec.IsNatEnabled = utils.BoolPtr(d.Get("is_nat_enabled").(bool))
	}
	if d.HasChange("is_external") {
		updateSpec.IsExternal = utils.BoolPtr(d.Get("is_external").(bool))
	}
	if d.HasChange("reserved_ip_addresses") {
		updateSpec.ReservedIpAddresses = expandIPAddress(d.Get("reserved_ip_addresses").([]interface{}))
	}
	if d.HasChange("dynamic_ip_addresses") {
		updateSpec.DynamicIpAddresses = expandIPAddress(d.Get("dynamic_ip_addresses").([]interface{}))
	}

	if d.HasChange("network_function_chain_reference") {
		updateSpec.NetworkFunctionChainReference = utils.StringPtr(d.Get("network_function_chain_reference").(string))
	}
	if d.HasChange("bridge_name") {
		updateSpec.BridgeName = utils.StringPtr(d.Get("bridge_name").(string))
	}
	if d.HasChange("is_advanced_networking") {
		updateSpec.IsAdvancedNetworking = utils.BoolPtr(d.Get("is_advanced_networking").(bool))
	}
	if d.HasChange("cluster_name") {
		updateSpec.ClusterName = utils.StringPtr(d.Get("cluster_name").(string))
	}
	if d.HasChange("hypervisor_type") {
		updateSpec.HypervisorType = utils.StringPtr(d.Get("hypervisor_type").(string))
	}
	if d.HasChange("virtual_switch") {
		updateSpec.VirtualSwitch = expandVirtualSwitch(d.Get("virtual_switch"))
	}
	if d.HasChange("vpc") {
		updateSpec.Vpc = expandVpc(d.Get("vpc"))
	}
	if d.HasChange("ip_prefix") {
		updateSpec.IpPrefix = utils.StringPtr(d.Get("ip_prefix").(string))
	}
	if d.HasChange("ip_usage") {
		updateSpec.IpUsage = expandIPUsage(d.Get("ip_usage"))
	}
	if d.HasChange("ip_config") {
		updateSpec.IpConfig = expandIPConfig(d.Get("ip_config").([]interface{}))
	} else {
		updateSpec.IpConfig = nil
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] Update Subnet Request: %s", string(aJSON))

	updateResp, err := conn.SubnetAPIInstance.UpdateSubnetById(utils.StringPtr(d.Id()), &updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating subnets : %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the subnet to be updated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for subnet (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixSubnetV2Read(ctx, d, meta)
}

func ResourceNutanixSubnetV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.SubnetAPIInstance.DeleteSubnetById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting subnet : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the subnet to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for subnet (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandDhcpOptions(pr []interface{}) *import1.DhcpOptions {
	if len(pr) > 0 {
		dhcpOps := import1.DhcpOptions{}

		val := pr[0].(map[string]interface{})

		if bootfn, ok := val["boot_file_name"]; ok && len(bootfn.(string)) > 0 {
			dhcpOps.BootFileName = utils.StringPtr(bootfn.(string))
		}
		if dns, ok := val["domain_name_servers"]; ok && len(dns.([]interface{})) > 0 {
			dhcpOps.DomainNameServers = expandIPAddress(dns.([]interface{}))
		}
		if dn, ok := val["domain_name"]; ok && len(dn.(string)) > 0 {
			dhcpOps.DomainName = utils.StringPtr(dn.(string))
		}
		if searchDomain, ok := val["search_domains"]; ok && len(searchDomain.([]interface{})) > 0 {
			dhcpOps.SearchDomains = common.ExpandListOfString(searchDomain.([]interface{}))
		}
		if tftp, ok := val["tftp_server_name"]; ok && len(tftp.(string)) > 0 {
			dhcpOps.TftpServerName = utils.StringPtr(tftp.(string))
		}
		if ntp, ok := val["ntp_servers"]; ok && len(ntp.([]interface{})) > 0 {
			dhcpOps.NtpServers = expandIPAddress(ntp.([]interface{}))
		}
		return &dhcpOps
	}
	return nil
}

func expandIPAddress(pr []interface{}) []config.IPAddress {
	if len(pr) > 0 {
		configList := make([]config.IPAddress, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			config := config.IPAddress{}

			if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
				config.Ipv4 = expandIPv4Address(ipv4)
			}
			if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
				config.Ipv6 = expandIPv6Address(ipv6)
			}

			configList[k] = config
		}
		return configList
	}
	return nil
}

func expandIPv4Address(pr interface{}) *config.IPv4Address {
	if pr == nil {
		return nil
	}

	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}

	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok || len(valMap) == 0 {
		return nil
	}

	ipv4 := &config.IPv4Address{}

	if v, ok := valMap["value"]; ok {
		if s, ok2 := v.(string); ok2 && len(s) > 0 {
			ipv4.Value = utils.StringPtr(s)
		}
	}

	if p, ok := valMap["prefix_length"]; ok {
		if n, ok2 := p.(int); ok2 {
			ipv4.PrefixLength = utils.IntPtr(n)
		}
	}

	return ipv4
}

func expandIPv6Address(pr interface{}) *config.IPv6Address {
	if pr == nil {
		return nil
	}

	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}

	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok || len(valMap) == 0 {
		return nil
	}

	ipv6 := &config.IPv6Address{}

	if v, ok := valMap["value"]; ok {
		if s, ok2 := v.(string); ok2 && len(s) > 0 {
			ipv6.Value = utils.StringPtr(s)
		}
	}

	if p, ok := valMap["prefix_length"]; ok {
		if n, ok2 := p.(int); ok2 {
			ipv6.PrefixLength = utils.IntPtr(n)
		}
	}

	return ipv6
}

func expandVirtualSwitch(pr interface{}) *import1.VirtualSwitch {
	if pr != nil {
		vSwitch := &import1.VirtualSwitch{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if extID, ok := val["ext_id"]; ok {
			vSwitch.ExtId = utils.StringPtr(extID.(string))
		}
		if name, ok := val["name"]; ok {
			vSwitch.Name = utils.StringPtr(name.(string))
		}
		if desc, ok := val["description"]; ok {
			vSwitch.Description = utils.StringPtr(desc.(string))
		}
		if isDefault, ok := val["is_default"]; ok {
			vSwitch.IsDefault = utils.BoolPtr(isDefault.(bool))
		}
		if hasDepErr, ok := val["has_deployment_error"]; ok {
			vSwitch.HasDeploymentError = utils.BoolPtr(hasDepErr.(bool))
		}
		if mtu, ok := val["mtu"]; ok {
			vSwitch.Mtu = utils.Int64Ptr(mtu.(int64))
		}
		if bondMode, ok := val["bond_mode"]; ok {
			const two, three, four, five = 2, 3, 4, 5
			bondMap := map[string]interface{}{
				"ACTIVE_BACKUP": two,
				"BALANCE_SLB":   three,
				"BALANCE_TCP":   four,
				"NONE":          five,
			}
			pInt := bondMap[bondMode.(string)]
			p := import1.BondModeType(pInt.(int))
			vSwitch.BondMode = &p
		}
		if cls, ok := val["clusters"]; ok {
			vSwitch.Clusters = expandCluster(cls.([]interface{}))
		}
		if name, ok := val["name"]; ok {
			vSwitch.Name = utils.StringPtr(name.(string))
		}

		return vSwitch
	}
	return nil
}

func expandCluster(pr []interface{}) []import1.Cluster {
	if len(pr) > 0 {
		clsList := make([]import1.Cluster, len(pr))

		for k, v := range pr {
			cls := import1.Cluster{}
			val := v.(map[string]interface{})

			if extID, ok := val["ext_id"]; ok {
				cls.ExtId = utils.StringPtr(extID.(string))
			}
			if hosts, ok := val["hosts"]; ok {
				cls.Hosts = expandHost(hosts.([]interface{}))
			}
			if gateway, ok := val["gateway_ip_address"]; ok {
				cls.GatewayIpAddress = expandIPv4Address(gateway)
			}
			clsList[k] = cls
		}
		return clsList
	}
	return nil
}

func expandHost(pr []interface{}) []import1.Host {
	if len(pr) > 0 {
		hosts := make([]import1.Host, len(pr))

		for k, v := range pr {
			host := import1.Host{}
			val := v.(map[string]interface{})

			if extID, ok := val["ext_id"]; ok {
				host.ExtId = utils.StringPtr(extID.(string))
			}
			if hostNics, ok := val["host_nics"]; ok {
				host.HostNics = utils.StringValueSlice(hostNics.([]*string))
			}
			if ipAdd, ok := val["ip_address"]; ok {
				host.IpAddress = expandIPv4Subnet(ipAdd)
			}

			hosts[k] = host
		}
		return hosts
	}
	return nil
}

func expandIPv4Subnet(pr interface{}) *import1.IPv4Subnet {
	if pr == nil {
		return nil
	}

	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}

	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok || len(valMap) == 0 {
		return nil
	}

	ipv4Subs := &import1.IPv4Subnet{}

	if ip, ok := valMap["ip"]; ok {
		ipv4Subs.Ip = expandIPv4Address(ip)
	}

	if prefix, ok := valMap["prefix_length"]; ok {
		if n, ok2 := prefix.(int); ok2 {
			ipv4Subs.PrefixLength = utils.IntPtr(n)
		}
	}

	return ipv4Subs
}

func expandIPv6Subnet(pr interface{}) *import1.IPv6Subnet {
	if pr == nil {
		return nil
	}

	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}

	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok || len(valMap) == 0 {
		return nil
	}

	ipv6Subs := &import1.IPv6Subnet{}

	if ip, ok := valMap["ip"]; ok {
		ipv6Subs.Ip = expandIPv6Address(ip)
	}

	if prefix, ok := valMap["prefix_length"]; ok {
		if n, ok2 := prefix.(int); ok2 {
			ipv6Subs.PrefixLength = utils.IntPtr(n)
		}
	}

	return ipv6Subs
}

func expandVpc(pr interface{}) *import1.Vpc {
	if pr != nil {
		vpc := &import1.Vpc{}
		prI := pr.([]interface{})

		val := prI[0].(map[string]interface{})

		if ext, ok := val["ext_id"]; ok && len(ext.(string)) > 0 {
			vpc.ExtId = utils.StringPtr(ext.(string))
		}
		if vpcType, ok := val["vpc_type"]; ok && len(vpcType.(string)) > 0 {
			const two, three = 2, 3
			vpcMap := map[string]interface{}{
				"REGULAR": two,
				"TRANSIT": three,
			}
			pInt := vpcMap[vpcType.(string)]
			p := import1.VpcType(pInt.(int))
			vpc.VpcType = &p
		}
		if desc, ok := val["description"]; ok && len(desc.(string)) > 0 {
			vpc.Description = utils.StringPtr(desc.(string))
		}
		if dhcpOps, ok := val["common_dhcp_options"]; ok && len(dhcpOps.([]interface{})) > 0 {
			vpc.CommonDhcpOptions = expandVpcDhcpOptions(dhcpOps)
		}
		if extSubs, ok := val["external_subnets"]; ok && len(extSubs.([]interface{})) > 0 {
			vpc.ExternalSubnets = expandExternalSubnet(extSubs.([]interface{}))
		}
		if extRoutingDomainRef, ok := val["external_routing_domain_reference"]; ok && len(extRoutingDomainRef.(string)) > 0 {
			vpc.ExternalRoutingDomainReference = utils.StringPtr(extRoutingDomainRef.(string))
		}
		if extRoutablePrefix, ok := val["externally_routable_prefixes"]; ok && len(extRoutablePrefix.([]interface{})) > 0 {
			vpc.ExternallyRoutablePrefixes = expandIPSubnet(extRoutablePrefix.([]interface{}))
		}
		return vpc
	}
	return nil
}

func expandVpcDhcpOptions(pr interface{}) *import1.VpcDhcpOptions {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		vpc := &import1.VpcDhcpOptions{}

		if dns, ok := val["domain_name_servers"]; ok {
			vpc.DomainNameServers = expandIPAddress(dns.([]interface{}))
		}
		return vpc
	}
	return nil
}

func expandExternalSubnet(pr []interface{}) []import1.ExternalSubnet {
	if len(pr) > 0 {
		extSubs := make([]import1.ExternalSubnet, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			sub := import1.ExternalSubnet{}

			if subRef, ok := val["subnet_reference"]; ok && len(subRef.(string)) > 0 {
				sub.SubnetReference = utils.StringPtr(subRef.(string))
			}
			if extips, ok := val["external_ips"]; ok && len(extips.([]interface{})) > 0 {
				sub.ExternalIps = expandIPAddress(extips.([]interface{}))
			}
			if gatewayNodes, ok := val["gateway_nodes"]; ok && len(gatewayNodes.([]interface{})) > 0 {
				sub.GatewayNodes = common.ExpandListOfString(gatewayNodes.([]interface{}))
			}
			if activeGatewayNode, ok := val["active_gateway_node"]; ok && len(activeGatewayNode.([]interface{})) > 0 {
				sub.ActiveGatewayNodes = expandGatewayNodeReference(activeGatewayNode)
			}
			if activeGatewayCount, ok := val["active_gateway_count"]; ok && activeGatewayCount.(int) > 0 {
				sub.ActiveGatewayCount = utils.IntPtr(activeGatewayCount.(int))
			}
			extSubs[k] = sub
		}
		return extSubs
	}
	return nil
}

func expandGatewayNodeReference(pr interface{}) []import1.GatewayNodeReference {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		gatewayNodesRef := &import1.GatewayNodeReference{}

		if nodeID, ok := val["node_id"]; ok {
			gatewayNodesRef.NodeId = utils.StringPtr(nodeID.(string))
		}
		if nodeipAdd, ok := val["node_ip_address"]; ok {
			gatewayNodesRef.NodeIpAddress = expandIPAddressMap(nodeipAdd)
		}
		gatewayNodesRefList := make([]import1.GatewayNodeReference, 1)
		gatewayNodesRefList[0] = *gatewayNodesRef
		return gatewayNodesRefList
	}
	return nil
}

func expandIPAddressMap(pr interface{}) *config.IPAddress {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		ipAdd := &config.IPAddress{}

		if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
			ipAdd.Ipv4 = expandIPv4AddressMap(ipv4)
		}
		if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
			ipAdd.Ipv6 = expandIPv6AddressMap(ipv6)
		}

		return ipAdd
	}
	return nil
}

func expandIPv4AddressMap(pr interface{}) *config.IPv4Address {
	if pr == nil {
		return nil
	}

	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}

	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok || len(valMap) == 0 {
		return nil
	}

	ipv4Add := &config.IPv4Address{}

	if v, ok := valMap["value"]; ok {
		if s, ok2 := v.(string); ok2 {
			ipv4Add.Value = utils.StringPtr(s)
		}
	}

	if p, ok := valMap["prefix_length"]; ok {
		if n, ok2 := p.(int); ok2 {
			ipv4Add.PrefixLength = utils.IntPtr(n)
		}
	}

	return ipv4Add
}

func expandIPv6AddressMap(pr interface{}) *config.IPv6Address {
	if pr == nil {
		return nil
	}

	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}

	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok || len(valMap) == 0 {
		return nil
	}

	ipv6Add := &config.IPv6Address{}

	if v, ok := valMap["value"]; ok {
		if s, ok2 := v.(string); ok2 {
			ipv6Add.Value = utils.StringPtr(s)
		}
	}

	if p, ok := valMap["prefix_length"]; ok {
		if n, ok2 := p.(int); ok2 {
			ipv6Add.PrefixLength = utils.IntPtr(n)
		}
	}

	return ipv6Add
}

func expandIPSubnet(pr []interface{}) []import1.IPSubnet {
	if len(pr) > 0 {
		ips := make([]import1.IPSubnet, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			ip := import1.IPSubnet{}

			if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
				ip.Ipv4 = expandIPv4Subnet(ipv4)
			}
			if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
				ip.Ipv6 = expandIPv6Subnet(ipv6)
			}
			ips[k] = ip
		}

		return ips
	}
	return nil
}

func expandIPUsage(pr interface{}) *import1.IPUsage {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		ipUsage := &import1.IPUsage{}

		if numMacs, ok := val["num_macs"]; ok {
			ipUsage.NumMacs = utils.Int64Ptr(numMacs.(int64))
		}
		if numFreeIPS, ok := val["num_free_ips"]; ok {
			ipUsage.NumFreeIPs = utils.Int64Ptr(numFreeIPS.(int64))
		}
		if numAssgIPs, ok := val["num_assigned_ips"]; ok {
			ipUsage.NumAssignedIPs = utils.Int64Ptr(numAssgIPs.(int64))
		}
		return ipUsage
	}
	return nil
}

func expandIPConfig(pr []interface{}) []import1.IPConfig {
	if len(pr) > 0 {
		ipConfigs := make([]import1.IPConfig, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			ipConfig := import1.IPConfig{}

			if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
				ipConfig.Ipv4 = expandIPv4Config(ipv4)
			}
			if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
				ipConfig.Ipv6 = expandIPv6Config(ipv6)
			}
			ipConfigs[k] = ipConfig
		}
		return ipConfigs
	}
	return nil
}

func expandIPv4Config(pr interface{}) *import1.IPv4Config {
	if pr != nil {
		ipv4Config := &import1.IPv4Config{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if ipSub, ok := val["ip_subnet"]; ok && len(ipSub.([]interface{})) > 0 {
			ipv4Config.IpSubnet = expandIPv4Subnet(ipSub)
		}
		if defaultGateway, ok := val["default_gateway_ip"]; ok && len(defaultGateway.([]interface{})) > 0 {
			ipv4Config.DefaultGatewayIp = expandIPv4Address(defaultGateway)
		}
		if dhcpServer, ok := val["dhcp_server_address"]; ok && len(dhcpServer.([]interface{})) > 0 {
			ipv4Config.DhcpServerAddress = expandIPv4Address(dhcpServer)
		}
		if pool, ok := val["pool_list"]; ok && len(pool.([]interface{})) > 0 {
			ipv4Config.PoolList = expandIPv4Pool(pool.([]interface{}))
		}
		return ipv4Config
	}
	return nil
}

func expandIPv6Config(pr interface{}) *import1.IPv6Config {
	if pr != nil {
		ipv4Config := &import1.IPv6Config{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if ipSub, ok := val["ip_subnet"]; ok && len(ipSub.([]interface{})) > 0 {
			ipv4Config.IpSubnet = expandIPv6Subnet(ipSub)
		}
		if defaultGateway, ok := val["default_gateway_ip"]; ok && len(defaultGateway.([]interface{})) > 0 {
			ipv4Config.DefaultGatewayIp = expandIPv6Address(defaultGateway)
		}
		if dhcpServer, ok := val["dhcp_server_address"]; ok && len(dhcpServer.([]interface{})) > 0 {
			ipv4Config.DhcpServerAddress = expandIPv6Address(dhcpServer)
		}
		if pool, ok := val["pool_list"]; ok && len(pool.([]interface{})) > 0 {
			ipv4Config.PoolList = expandIPv6Pool(pool.([]interface{}))
		}
		return ipv4Config
	}
	return nil
}

func expandIPv4Pool(pr []interface{}) []import1.IPv4Pool {
	if len(pr) > 0 {
		pools := make([]import1.IPv4Pool, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			pool := import1.IPv4Pool{}

			if startIP, ok := val["start_ip"]; ok && len(startIP.([]interface{})) > 0 {
				pool.StartIp = expandIPv4Address(startIP)
			}
			if endIP, ok := val["end_ip"]; ok && len(endIP.([]interface{})) > 0 {
				pool.EndIp = expandIPv4Address(endIP)
			}
			pools[k] = pool
		}
		return pools
	}
	return nil
}

func expandIPv6Pool(pr []interface{}) []import1.IPv6Pool {
	if len(pr) > 0 {
		pools := make([]import1.IPv6Pool, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			pool := import1.IPv6Pool{}

			if startIP, ok := val["start_ip"]; ok && len(startIP.([]interface{})) > 0 {
				pool.StartIp = expandIPv6Address(startIP)
			}
			if endIP, ok := val["end_ip"]; ok && len(endIP.([]interface{})) > 0 {
				pool.EndIp = expandIPv6Address(endIP)
			}
			pools[k] = pool
		}
		return pools
	}
	return nil
}
