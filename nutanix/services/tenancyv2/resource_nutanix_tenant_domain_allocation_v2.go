package tenancyv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	prismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/tasks"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/prism/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixTenantDomainAllocationV2 manages the lifecycle of a tenant domain allocation.
func ResourceNutanixTenantDomainAllocationV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixTenantDomainAllocationV2Create,
		ReadContext:   ResourceNutanixTenantDomainAllocationV2Read,
		UpdateContext: ResourceNutanixTenantDomainAllocationV2Update,
		DeleteContext: ResourceNutanixTenantDomainAllocationV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_network_config": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem:     externalNetworkConnectionResourceElem(),
			},
			"projects": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"memory": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"storage": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"vcpu": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"resource_group": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"placement_targets": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster_uuid": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"storage_container_uuids": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"state_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
		},
	}
}

// externalNetworkConnectionResourceElem returns the resource-side schema for an
// external network connection block (input semantics on the config sub-block).
func externalNetworkConnectionResourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     externalNetworkConnConfigResourceElem(),
			},
			"state_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// externalNetworkConnConfigResourceElem returns the resource-side schema for the
// NAT / no-NAT external network connection configuration.
func externalNetworkConnConfigResourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nat": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transit_vpc_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"floating_ip_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"overlay_subnet_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"no_nat": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"overlay_subnet_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"transit_vpc_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"externally_routable_prefixes": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

// ResourceNutanixTenantDomainAllocationV2Create creates a new domain allocation for a tenant.
func ResourceNutanixTenantDomainAllocationV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	tenantExtID := d.Get("tenant_ext_id").(string)

	body := import1.DomainAllocation{}
	if v, ok := d.GetOk("domain_ext_id"); ok {
		body.DomainExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("external_network_config"); ok {
		body.ExternalNetworkConfig = expandExternalNetworkConnection(v.([]interface{}))
	}
	if v, ok := d.GetOk("projects"); ok {
		body.Projects = expandProjects(v.([]interface{}))
	}

	createRequest := import4.CreateTenantDomainAllocationRequest{
		TenantExtId: utils.StringPtr(tenantExtID),
		Body:        &body,
	}
	aJSON, _ := json.MarshalIndent(createRequest, "", "  ")
	log.Printf("[DEBUG] Tenant domain allocation create payload: %s", string(aJSON))

	resp, err := conn.TenantManagementAPIInstance.CreateTenantDomainAllocation(ctx, &createRequest)
	if err != nil {
		return diag.Errorf("error while creating tenant domain allocation: %v", err)
	}

	taskRef := resp.Data.GetValue().(import3.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for tenant domain allocation (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	getTaskByIDRequest := import5.GetTaskByIdRequest{
		ExtId: utils.StringPtr(*taskUUID),
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIDRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant domain allocation create task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeTenantDomainAllocation, "TenantDomainAllocation")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixTenantDomainAllocationV2Read(ctx, d, meta)
}

// ResourceNutanixTenantDomainAllocationV2Read reads a domain allocation by its external identifier.
func ResourceNutanixTenantDomainAllocationV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	tenantExtID := d.Get("tenant_ext_id").(string)
	getRequest := import4.GetTenantDomainAllocationByIdRequest{
		TenantExtId: utils.StringPtr(tenantExtID),
		ExtId:       utils.StringPtr(d.Id()),
	}
	resp, err := conn.TenantManagementAPIInstance.GetTenantDomainAllocationById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant domain allocation: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.DomainAllocation)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_ext_id", getResp.DomainExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("external_network_config", flattenExternalNetworkConnection(getResp.ExternalNetworkConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("projects", flattenProjects(getResp.Projects)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state_info", flattenDomainAllocationStateInfo(getResp.StateInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// ResourceNutanixTenantDomainAllocationV2Update updates a domain allocation using a read-modify-write flow.
func ResourceNutanixTenantDomainAllocationV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	tenantExtID := d.Get("tenant_ext_id").(string)
	getRequest := import4.GetTenantDomainAllocationByIdRequest{
		TenantExtId: utils.StringPtr(tenantExtID),
		ExtId:       utils.StringPtr(d.Id()),
	}
	resp, err := conn.TenantManagementAPIInstance.GetTenantDomainAllocationById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant domain allocation: %v", err)
	}

	updateSpec := resp.Data.GetValue().(import1.DomainAllocation)

	if d.HasChange("external_network_config") {
		updateSpec.ExternalNetworkConfig = expandExternalNetworkConnection(d.Get("external_network_config").([]interface{}))
	}
	if d.HasChange("projects") {
		updateSpec.Projects = expandProjects(d.Get("projects").([]interface{}))
	}
	if d.HasChange("domain_ext_id") {
		updateSpec.DomainExtId = utils.StringPtr(d.Get("domain_ext_id").(string))
	}

	etagValue := conn.TenantManagementAPIInstance.ApiClient.GetEtag(resp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	updateRequest := import4.UpdateTenantDomainAllocationByIdRequest{
		TenantExtId: utils.StringPtr(tenantExtID),
		ExtId:       utils.StringPtr(d.Id()),
		Body:        &updateSpec,
	}
	aJSON, _ := json.MarshalIndent(updateRequest, "", "  ")
	log.Printf("[DEBUG] Tenant domain allocation update payload: %s", string(aJSON))

	updateResp, err := conn.TenantManagementAPIInstance.UpdateTenantDomainAllocationById(ctx, &updateRequest, args)
	if err != nil {
		return diag.Errorf("error while updating tenant domain allocation: %v", err)
	}

	taskRef := updateResp.Data.GetValue().(import3.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for tenant domain allocation (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return ResourceNutanixTenantDomainAllocationV2Read(ctx, d, meta)
}

// ResourceNutanixTenantDomainAllocationV2Delete deletes a domain allocation and polls the delete task.
func ResourceNutanixTenantDomainAllocationV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	tenantExtID := d.Get("tenant_ext_id").(string)
	deleteRequest := import4.DeleteTenantDomainAllocationByIdRequest{
		TenantExtId: utils.StringPtr(tenantExtID),
		ExtId:       utils.StringPtr(d.Id()),
	}
	resp, err := conn.TenantManagementAPIInstance.DeleteTenantDomainAllocationById(ctx, &deleteRequest)
	if err != nil {
		return diag.Errorf("error while deleting tenant domain allocation: %v", err)
	}

	taskRef := resp.Data.GetValue().(import3.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for tenant domain allocation (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}
