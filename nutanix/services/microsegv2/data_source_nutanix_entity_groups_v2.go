package microsegv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/common/v1/config"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/request/entitygroups"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixEntityGroupsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixEntityGroupsV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entity_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": schemaForLinks(),
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"allowed_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaAllowedConfig(),
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixEntityGroupsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	listEntityGroupsRequest := import2.ListEntityGroupsRequest{}

	if v, ok := d.GetOk("page"); ok {
		listEntityGroupsRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listEntityGroupsRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listEntityGroupsRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listEntityGroupsRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listEntityGroupsRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.EntityGroupsAPIInstance.ListEntityGroups(ctx, &listEntityGroupsRequest)
	if err != nil {
		return diag.Errorf("error while Listing Entity Groups: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("entity_groups", []map[string]interface{}{}); err != nil {
			return diag.Errorf("error setting Entity Groups: %s", err)
		}
		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of entity groups.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.EntityGroup)

	if err := d.Set("entity_groups", flattenEntityGroups(getResp)); err != nil {
		return diag.Errorf("error setting Entity Groups: %s", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenEntityGroups(entityGroups []import1.EntityGroup) []map[string]interface{} {
	if len(entityGroups) == 0 {
		return []map[string]interface{}{}
	}

	entityGroupsList := make([]map[string]interface{}, 0)

	for _, entityGroup := range entityGroups {
		entityGroupMap := make(map[string]interface{})

		entityGroupMap["tenant_id"] = utils.StringValue(entityGroup.TenantId)
		entityGroupMap["ext_id"] = entityGroup.ExtId
		entityGroupMap["links"] = flattenLinks(entityGroup.Links)
		entityGroupMap["name"] = utils.StringValue(entityGroup.Name)
		entityGroupMap["description"] = utils.StringValue(entityGroup.Description)
		entityGroupMap["owner_ext_id"] = utils.StringValue(entityGroup.OwnerExtId)
		entityGroupMap["policy_ext_ids"] = entityGroup.PolicyExtIds
		entityGroupMap["allowed_config"] = flattenAllowedConfig(entityGroup.AllowedConfig)

		entityGroupsList = append(entityGroupsList, entityGroupMap)
	}

	return entityGroupsList
}
