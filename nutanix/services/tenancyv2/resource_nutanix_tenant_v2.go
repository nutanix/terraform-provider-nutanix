package tenancyv2

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/prism/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixTenantV2 manages the lifecycle of a tenant in the tenancy namespace.
func ResourceNutanixTenantV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixTenantV2Create,
		ReadContext:   ResourceNutanixTenantV2Read,
		DeleteContext: ResourceNutanixTenantV2Delete,
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
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"added_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"admin_info": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"email": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							// WriteOnly cannot combine with ForceNew; password is create-time only.
							WriteOnly: true,
						},
					},
				},
			},
			"aggregated_quota": {
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
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
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

// ResourceNutanixTenantV2Create creates a new tenant and polls the create task.
func ResourceNutanixTenantV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	body := import1.Tenant{}
	if v, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("admin_info"); ok {
		body.AdminInfo = expandLocalUserInfo(v.([]interface{}))
	}

	createTenantRequest := import4.CreateTenantRequest{
		Body: &body,
	}
	aJSON, _ := json.MarshalIndent(createTenantRequest, "", "  ")
	log.Printf("[DEBUG] Tenant create payload: %s", string(aJSON))

	resp, err := conn.TenantManagementAPIInstance.CreateTenant(ctx, &createTenantRequest)
	if err != nil {
		return diag.Errorf("error while creating tenant: %v", err)
	}

	taskRef := resp.Data.GetValue().(import3.TaskReference)
	d.SetId(utils.StringValue(taskRef.ExtId))

	return ResourceNutanixTenantV2Read(ctx, d, meta)
}

// ResourceNutanixTenantV2Read reads a tenant by its external identifier.
func ResourceNutanixTenantV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	getTenantRequest := import4.GetTenantByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.TenantManagementAPIInstance.GetTenantById(ctx, &getTenantRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Tenant)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("url", getResp.Url); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("added_by", getResp.AddedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domains", getResp.Domains); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("aggregated_quota", flattenQuota(getResp.AggregatedQuota)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("state_info", flattenTenantStateInfo(getResp.StateInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	// admin_info: the API never returns the password, so preserve the configured
	// value while updating the other fields from the response.
	adminInfo := flattenLocalUserInfo(getResp.AdminInfo)
	if adminInfo != nil {
		if v, ok := d.GetOk("admin_info"); ok {
			configured := v.([]interface{})
			if len(configured) > 0 && configured[0] != nil {
				if pwd, ok := configured[0].(map[string]interface{})["password"]; ok {
					adminInfo[0]["password"] = pwd
				}
			}
		}
	}
	if err := d.Set("admin_info", adminInfo); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// ResourceNutanixTenantV2Delete deletes a tenant and polls the delete task.
func ResourceNutanixTenantV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	deleteTenantRequest := import4.DeleteTenantByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}

	_, err := conn.TenantManagementAPIInstance.DeleteTenantById(ctx, &deleteTenantRequest)
	if err != nil {
		return diag.Errorf("error while deleting tenant: %v", err)
	}

	timeout := time.NewTimer(3 * time.Minute)
	ticker := time.NewTicker(3 * time.Second)
	defer timeout.Stop()
	defer ticker.Stop()

	for {
		getTenantRequest := import4.GetTenantByIdRequest{
			ExtId: utils.StringPtr(d.Id()),
		}
		readResp, err := conn.TenantManagementAPIInstance.GetTenantById(ctx, &getTenantRequest)
		if err != nil {
			return diag.Errorf("error while fetching tenant deletion state: %v", err)
		}
		if readResp == nil || readResp.Data == nil {
			return diag.Errorf("error while fetching tenant deletion state: API returned an empty response")
		}

		tenant, ok := readResp.Data.GetValue().(import1.Tenant)
		if !ok {
			return diag.Errorf("error while fetching tenant deletion state: unexpected response type %T", readResp.Data.GetValue())
		}
		if tenant.StateInfo.State.GetName() == "DELETED" {
			d.SetId("")
			return nil
		}

		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		case <-timeout.C:
			return diag.Errorf("timed out waiting up to 3 minutes for tenant %q to reach DELETED state", d.Id())
		case <-ticker.C:
		}
	}
}
