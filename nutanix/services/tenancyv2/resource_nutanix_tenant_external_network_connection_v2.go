package tenancyv2

import (
	"context"
	"encoding/json"
	"log"

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

// ResourceNutanixTenantExternalNetworkConnectionV2 manages the external network
// connection settings of a tenant domain allocation. This entity only supports
// Read and Update operations.
func ResourceNutanixTenantExternalNetworkConnectionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixTenantExternalNetworkConnectionV2Update,
		ReadContext:   ResourceNutanixTenantExternalNetworkConnectionV2Read,
		UpdateContext: ResourceNutanixTenantExternalNetworkConnectionV2Update,
		DeleteContext: ResourceNutanixTenantExternalNetworkConnectionV2Delete,
		Schema: map[string]*schema.Schema{
			"tenant_ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_allocation_ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ext_id": {
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
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
		},
	}
}

// ResourceNutanixTenantExternalNetworkConnectionV2Read reads the external network
// connection settings for a tenant domain allocation.
func ResourceNutanixTenantExternalNetworkConnectionV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	tenantExtID := d.Get("tenant_ext_id").(string)
	domainAllocationExtID := d.Get("domain_allocation_ext_id").(string)

	getRequest := import4.GetTenantExternalNetworkConnectionRequest{
		TenantExtId:           utils.StringPtr(tenantExtID),
		DomainAllocationExtId: utils.StringPtr(domainAllocationExtID),
	}
	resp, err := conn.TenantManagementAPIInstance.GetTenantExternalNetworkConnection(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant external network connection: %v", err)
	}

	var extID, tenantID *string
	var config *import1.ExternalNetworkConnConfig
	var stateInfo *import1.ExternalNetworkConfigStateInfo
	var links interface{}

	switch data := resp.Data.GetValue().(type) {
	case import1.ExternalNetworkConnectionProjection:
		extID = data.ExtId
		tenantID = data.TenantId
		config = data.Config
		stateInfo = data.StateInfo
		links = data.Links
	case import1.ExternalNetworkConnection:
		extID = data.ExtId
		tenantID = data.TenantId
		config = data.Config
		stateInfo = data.StateInfo
		links = data.Links
	default:
		return diag.Errorf("unexpected response type while fetching tenant external network connection")
	}

	if err := d.Set("ext_id", extID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("config", flattenExternalNetworkConnConfig(config)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state_info", flattenExternalNetworkConfigStateInfo(stateInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", tenantID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(links)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// ResourceNutanixTenantExternalNetworkConnectionV2Update updates the external
// network connection settings for a tenant domain allocation.
func ResourceNutanixTenantExternalNetworkConnectionV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	tenantExtID := d.Get("tenant_ext_id").(string)
	domainAllocationExtID := d.Get("domain_allocation_ext_id").(string)

	getRequest := import4.GetTenantExternalNetworkConnectionRequest{
		TenantExtId:           utils.StringPtr(tenantExtID),
		DomainAllocationExtId: utils.StringPtr(domainAllocationExtID),
	}
	readResp, err := conn.TenantManagementAPIInstance.GetTenantExternalNetworkConnection(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant external network connection: %v", err)
	}

	body := import1.ExternalNetworkConnection{}
	switch data := readResp.Data.GetValue().(type) {
	case import1.ExternalNetworkConnectionProjection:
		body.ExtId = data.ExtId
		body.TenantId = data.TenantId
		body.Config = data.Config
	case import1.ExternalNetworkConnection:
		body.ExtId = data.ExtId
		body.TenantId = data.TenantId
		body.Config = data.Config
	}

	if v, ok := d.GetOk("config"); ok {
		body.Config = expandExternalNetworkConnConfig(v.([]interface{}))
	}

	etagValue := conn.TenantManagementAPIInstance.ApiClient.GetEtag(readResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	updateRequest := import4.UpdateTenantExternalNetworkConnectionRequest{
		TenantExtId:           utils.StringPtr(tenantExtID),
		DomainAllocationExtId: utils.StringPtr(domainAllocationExtID),
		Body:                  &body,
	}
	aJSON, _ := json.MarshalIndent(updateRequest, "", "  ")
	log.Printf("[DEBUG] Tenant external network connection update payload: %s", string(aJSON))

	updateResp, err := conn.TenantManagementAPIInstance.UpdateTenantExternalNetworkConnection(ctx, &updateRequest, args)
	if err != nil {
		return diag.Errorf("error while updating tenant external network connection: %v", err)
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
		return diag.Errorf("error waiting for tenant external network connection (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	d.SetId(tenantExtID + "/" + domainAllocationExtID)

	return ResourceNutanixTenantExternalNetworkConnectionV2Read(ctx, d, meta)
}

// ResourceNutanixTenantExternalNetworkConnectionV2Delete removes the resource from
// state. The external network connection has no delete API; it lives with the
// domain allocation, so this only clears the Terraform state.
func ResourceNutanixTenantExternalNetworkConnectionV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
