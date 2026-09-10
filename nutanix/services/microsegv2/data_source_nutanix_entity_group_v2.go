package microsegv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/request/entitygroups"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixEntityGroupV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixEntityGroupV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allowed_config": schemaForAllowedConfig(),
			"except_config":  schemaForExceptConfig(),
			"policy_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_shared_with_all_projects": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_system_defined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixEntityGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	extID := d.Get("ext_id").(string)
	getEntityGroupByIDRequest := import3.GetEntityGroupByIdRequest{
		ExtId: utils.StringPtr(extID),
	}

	resp, err := conn.EntityGroupsAPIInstance.GetEntityGroupById(ctx, &getEntityGroupByIDRequest)
	if err != nil {
		return diag.Errorf("error while fetching Entity Group: %s", err)
	}

	if resp.Data == nil {
		return diag.Errorf("no data in GetEntityGroupById response")
	}

	getResp, ok := resp.Data.GetValue().(import2.EntityGroup)
	if !ok {
		return diag.Errorf("invalid EntityGroup type in response")
	}

	if err := d.Set("ext_id", utils.StringValue(getResp.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", utils.StringValue(getResp.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", utils.StringValue(getResp.Description)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allowed_config", flattenAllowedConfig(getResp.AllowedConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("except_config", flattenExceptConfig(getResp.ExceptConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("policy_ext_ids", getResp.PolicyExtIds); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", utils.StringValue(getResp.OwnerExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(getResp.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinksEntityGroup(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreationTime != nil {
		if err := d.Set("creation_time", getResp.CreationTime.Format("2006-01-02T15:04:05.000Z")); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdateTime != nil {
		if err := d.Set("last_update_time", getResp.LastUpdateTime.Format("2006-01-02T15:04:05.000Z")); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("project_ext_id", utils.StringValue(getResp.ProjectExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_shared_with_all_projects", utils.BoolValue(getResp.IsSharedWithAllProjects)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_system_defined", utils.BoolValue(getResp.IsSystemDefined)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
