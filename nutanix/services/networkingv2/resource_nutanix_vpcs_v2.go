package networkingv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVPCsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVPCsV2Create,
		ReadContext:   ResourceNutanixVPCsV2Read,
		UpdateContext: ResourceNutanixVPCsV2Update,
		DeleteContext: ResourceNutanixVPCsV2Delete,
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
				Computed: true,
				Optional: true,
			},
			"vpc_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"REGULAR", "TRANSIT"}, false),
			},
			"common_dhcp_options": {
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
					},
				},
			},
			"external_subnets": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_reference": {
							Type:     schema.TypeString,
							Required: true,
						},
						"external_ips": {
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
						"gateway_nodes": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"active_gateway_node": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"node_ip_address": {
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
						"active_gateway_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"external_routing_domain_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"externally_routable_prefixes": {
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
									"ip": SchemaForValuePrefixLength(),
									"prefix_length": {
										Type:     schema.TypeInt,
										Required: true,
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
									"ip": SchemaForValuePrefixLength(),
									"prefix_length": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
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
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DatasourceMetadataSchemaV2(),
				},
			},
			"snat_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(),
						"ipv6": SchemaForValuePrefixLength(),
					},
				},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixVPCsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	inputSpec := import1.Vpc{}

	if metadata, ok := d.GetOk("metadata"); ok {
		inputSpec.Metadata = expandMetadata(metadata.([]interface{}))
	}
	if name, ok := d.GetOk("name"); ok {
		inputSpec.Name = utils.StringPtr(name.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		inputSpec.Description = utils.StringPtr(description.(string))
	}
	if vpcType, ok := d.GetOk("vpc_type"); ok {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"REGULAR": two,
			"TRANSIT": three,
		}
		pVal := subMap[vpcType.(string)]

		p := import1.VpcType(pVal.(int))
		inputSpec.VpcType = &p
	}

	if dhcp, ok := d.GetOk("common_dhcp_options"); ok {
		inputSpec.CommonDhcpOptions = expandVpcDhcpOptions(dhcp.([]interface{}))
	}

	if externalSubnets, ok := d.GetOk("external_subnets"); ok {
		inputSpec.ExternalSubnets = expandExternalSubnet(externalSubnets.([]interface{}))
	}

	if externalRoutingDomainReference, ok := d.GetOk("external_routing_domain_reference"); ok {
		inputSpec.ExternalRoutingDomainReference = utils.StringPtr(externalRoutingDomainReference.(string))
	}

	if externallyRoutablePrefixes, ok := d.GetOk("externally_routable_prefixes"); ok {
		inputSpec.ExternallyRoutablePrefixes = expandIPSubnet(externallyRoutablePrefixes.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(inputSpec, "", " ")
	log.Printf("[DEBUG] VPC create payload : %s", string(aJSON))

	resp, err := conn.VpcAPIInstance.CreateVpc(&inputSpec)
	if err != nil {
		return diag.Errorf("error while creating floating IPs : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VPC to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vpc (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching vpc UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)
	return ResourceNutanixVPCsV2Read(ctx, d, meta)
}

func ResourceNutanixVPCsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.VpcAPIInstance.GetVpcById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching vpc : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Vpc)

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
	if err := d.Set("vpc_type", getResp.VpcType.GetName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("common_dhcp_options", flattenCommonDhcpOptions(getResp.CommonDhcpOptions)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("snat_ips", flattenNtpServer(getResp.SnatIps)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("external_subnets", flattenExternalSubnets(getResp.ExternalSubnets)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("external_routing_domain_reference", getResp.ExternalRoutingDomainReference); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("externally_routable_prefixes", flattenExternallyRoutablePrefixes(getResp.ExternallyRoutablePrefixes)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixVPCsV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.VpcAPIInstance.GetVpcById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching vpcs : %v", err)
	}

	respVpc := resp.Data.GetValue().(import1.Vpc)

	updateSpec := respVpc

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("vpc_type") {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"REGULAR": two,
			"TRANSIT": three,
		}
		pVal := subMap[d.Get("vpc_type").(string)]
		if pVal == nil {
			updateSpec.VpcType = nil
		} else {
			p := import1.VpcType(pVal.(int))
			updateSpec.VpcType = &p
		}
	}
	if d.HasChange("common_dhcp_options") {
		updateSpec.CommonDhcpOptions = expandVpcDhcpOptions(d.Get("common_dhcp_options").([]interface{}))
	}
	if d.HasChange("external_subnets") {
		updateSpec.ExternalSubnets = expandExternalSubnet(d.Get("external_subnets").([]interface{}))
	}
	if d.HasChange("external_routing_domain_reference") {
		updateSpec.ExternalRoutingDomainReference = utils.StringPtr(d.Get("external_routing_domain_reference").(string))
	}
	if d.HasChange("externally_routable_prefixes") {
		updateSpec.ExternallyRoutablePrefixes = expandIPSubnet(d.Get("externally_routable_prefixes").([]interface{}))
	}

	updateResp, err := conn.VpcAPIInstance.UpdateVpcById(utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error while updating vpcs : %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(import4.TaskReference)
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
		return diag.Errorf("error waiting for vpc (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixVPCsV2Read(ctx, d, meta)
}

func ResourceNutanixVPCsV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	resp, err := conn.VpcAPIInstance.DeleteVpcById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting vpc : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Subnet to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vpc (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}
