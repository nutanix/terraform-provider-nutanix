package tenancyv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Optional: true,
				Computed: true,
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

	tenantExtID := utils.StringPtr(d.Get("tenant_ext_id").(string))

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
		TenantExtId: tenantExtID,
		Body:        &body,
	}
	aJSON, _ := json.MarshalIndent(createRequest, "", "  ")
	log.Printf("[DEBUG] Tenant domain allocation create payload: %s", string(aJSON))

	resp, err := conn.TenantManagementAPIInstance.CreateTenantDomainAllocation(ctx, &createRequest)
	if err != nil {
		return diag.Errorf("error while creating tenant domain allocation: %v", err)
	}

	taskRef := resp.Data.GetValue().(import3.TaskReference)
	d.SetId(utils.StringValue(taskRef.ExtId))
	ResourceNutanixTenantDomainAllocationV2Read(ctx, d, meta)

	// Update the domain allocation with the external network config
	if v, ok := d.GetOk("external_network_config"); ok {
		if err := updateTenantDomainAllocationExternalNetworkConnection(ctx, d, meta, utils.StringValue(tenantExtID), v.([]interface{}), d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

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
	body := import1.DomainAllocation{}
	if d.HasChange("projects") {
		body.Projects = expandProjects(d.Get("projects").([]interface{}))
	} else {
		body.Projects = updateSpec.Projects
	}
	if d.HasChange("domain_ext_id") {
		body.DomainExtId = utils.StringPtr(d.Get("domain_ext_id").(string))
	} else {
		body.DomainExtId = updateSpec.DomainExtId
	}

	etagValue := conn.TenantManagementAPIInstance.ApiClient.GetEtag(resp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	updateRequest := import4.UpdateTenantDomainAllocationByIdRequest{
		TenantExtId: utils.StringPtr(tenantExtID),
		ExtId:       utils.StringPtr(d.Id()),
		Body:        &body,
	}
	aJSON, _ := json.MarshalIndent(updateRequest, "", "  ")
	log.Printf("[DEBUG] Tenant domain allocation update payload: %s", string(aJSON))

	_, err = conn.TenantManagementAPIInstance.UpdateTenantDomainAllocationById(ctx, &updateRequest, args)
	if err != nil {
		return diag.Errorf("error while updating tenant domain allocation: %v", err)
	}
	ResourceNutanixTenantDomainAllocationV2Read(ctx, d, meta)

	// Update the domain allocation with the external network config, if it has changed
	if d.HasChange("external_network_config") {
		if err := updateTenantDomainAllocationExternalNetworkConnection(ctx, d, meta, tenantExtID, d.Get("external_network_config").([]interface{}), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}
	return ResourceNutanixTenantDomainAllocationV2Read(ctx, d, meta)
}

// ResourceNutanixTenantDomainAllocationV2Delete deletes a domain allocation and polls the delete task.
func ResourceNutanixTenantDomainAllocationV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	tenantExtID := d.Get("tenant_ext_id").(string)
	getRequest := import4.GetTenantDomainAllocationByIdRequest{
		TenantExtId: utils.StringPtr(tenantExtID),
		ExtId:       utils.StringPtr(d.Id()),
	}
	readResp, err := conn.TenantManagementAPIInstance.GetTenantDomainAllocationById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant domain allocation: %v", err)
	}
	etagValue := conn.TenantManagementAPIInstance.ApiClient.GetEtag(readResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	deleteRequest := import4.DeleteTenantDomainAllocationByIdRequest{
		TenantExtId: utils.StringPtr(tenantExtID),
		ExtId:       utils.StringPtr(d.Id()),
	}
	_, err = conn.TenantManagementAPIInstance.DeleteTenantDomainAllocationById(ctx, &deleteRequest, args)
	if err != nil {
		return diag.Errorf("error while deleting tenant domain allocation: %v", err)
	}

	timeout := time.NewTimer(3 * time.Minute)
	ticker := time.NewTicker(3 * time.Second)
	defer timeout.Stop()
	defer ticker.Stop()

	for {
		getRequest := import4.GetTenantDomainAllocationByIdRequest{
			TenantExtId: utils.StringPtr(tenantExtID),
			ExtId:       utils.StringPtr(d.Id()),
		}
		readResp, err := conn.TenantManagementAPIInstance.GetTenantDomainAllocationById(ctx, &getRequest)
		if err != nil {
			return diag.Errorf("error while fetching tenant domain allocation: %v", err)
		}
		if readResp == nil || readResp.Data == nil {
			return diag.Errorf("error while fetching tenant domain allocation deletion state: API returned an empty response")
		}

		domainAllocation, ok := readResp.Data.GetValue().(import1.DomainAllocation)
		if !ok {
			return diag.Errorf("error while fetching tenant domain allocation deletion state: unexpected response type %T", readResp.Data.GetValue())
		}
		if domainAllocation.StateInfo != nil &&
			domainAllocation.StateInfo.State != nil &&
			domainAllocation.StateInfo.State.GetName() == "DELETED" {
			d.SetId("")
			return nil
		}

		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		case <-timeout.C:
			return diag.Errorf("timed out waiting up to 3 minutes for tenant domain allocation %q to reach DELETED state", d.Id())
		case <-ticker.C:
		}
	}
}

func updateTenantDomainAllocationExternalNetworkConnection(
	ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	tenantExtID string,
	externalNetworkConfig []interface{},
	timeout time.Duration,
) error {
	connection := expandExternalNetworkConnection(externalNetworkConfig)
	if connection == nil || connection.Config == nil {
		return fmt.Errorf("external_network_config.config must define a NAT or no-NAT configuration")
	}

	conn := meta.(*conns.Client).ServiceProviderManagementAPI
	getRequest := import4.GetTenantExternalNetworkConnectionRequest{
		TenantExtId:           utils.StringPtr(tenantExtID),
		DomainAllocationExtId: utils.StringPtr(d.Id()),
	}
	readResp, err := conn.TenantManagementAPIInstance.GetTenantExternalNetworkConnection(ctx, &getRequest)
	if err != nil {
		return fmt.Errorf("error while fetching tenant external network connection: %w", err)
	}

	etagValue := conn.TenantManagementAPIInstance.ApiClient.GetEtag(readResp)
	args := map[string]interface{}{
		"If-Match": utils.StringPtr(etagValue),
	}
	updateRequest := import4.UpdateTenantExternalNetworkConnectionRequest{
		TenantExtId:           utils.StringPtr(tenantExtID),
		DomainAllocationExtId: utils.StringPtr(d.Id()),
		Body:                  connection.Config,
	}
	aJSON, err := json.MarshalIndent(updateRequest, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing tenant external network connection update: %w", err)
	}
	log.Printf("[DEBUG] Tenant domain allocation external network connection update payload: %s", string(aJSON))

	updateResp, err := conn.TenantManagementAPIInstance.UpdateTenantExternalNetworkConnection(ctx, &updateRequest, args)
	if err != nil {
		return fmt.Errorf("error while updating tenant external network connection: %w", err)
	}
	if updateResp == nil || updateResp.Data == nil {
		return fmt.Errorf("error while updating tenant external network connection: API returned an empty response")
	}

	// taskRef, ok := updateResp.Data.GetValue().(import3.TaskReference)
	// if !ok {
	// 	return fmt.Errorf("error while updating tenant external network connection: unexpected response type %T", updateResp.Data.GetValue())
	// }
	// if taskRef.ExtId == nil || *taskRef.ExtId == "" {
	// 	return fmt.Errorf("error while updating tenant external network connection: task response did not include an ext_id")
	// }

	// taskconn := meta.(*conns.Client).PrismAPI
	// stateConf := &resource.StateChangeConf{
	// 	Pending: []string{"PENDING", "RUNNING", "QUEUED"},
	// 	Target:  []string{"SUCCEEDED"},
	// 	Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskRef.ExtId)),
	// 	Timeout: timeout,
	// }
	// if _, err := stateConf.WaitForStateContext(ctx); err != nil {
	// 	return fmt.Errorf("error waiting for tenant external network connection task %q: %w", utils.StringValue(taskRef.ExtId), err)
	// }
	return nil
}
