package microsegv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/request/directoryserverconfigs"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixDirectoryServerConfigsV2 lists the Directory Server
// configurations.
func DatasourceNutanixDirectoryServerConfigsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixDirectoryServerConfigsV2Read,
		Schema: map[string]*schema.Schema{
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"directory_server_configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: datasourceDirectoryServerConfigSchema(),
				},
			},
		},
	}
}

func DatasourceNutanixDirectoryServerConfigsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	listRequest := import3.ListDirectoryServerConfigsRequest{}
	if v, ok := d.GetOk("select"); ok {
		listRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.DirectoryServerConfigsAPIInstance.ListDirectoryServerConfigs(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while listing Directory Server Configs: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("directory_server_configs", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of Directory Server Configs.",
		}}
	}

	listVal, ok := resp.Data.GetValue().([]import2.DirectoryServerConfig)
	if !ok {
		if err := d.Set("directory_server_configs", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of Directory Server Configs.",
		}}
	}

	if err := d.Set("directory_server_configs", flattenDirectoryServerConfigs(listVal)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenDirectoryServerConfigs(configs []import2.DirectoryServerConfig) []map[string]interface{} {
	if len(configs) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(configs))
	for _, c := range configs {
		m := map[string]interface{}{
			"ext_id":                                utils.StringValue(c.ExtId),
			"directory_service_reference":           utils.StringValue(c.DirectoryServiceReference),
			"domain_controllers":                    flattenIPAddressOrFQDNList(c.DomainControllers),
			"matching_criterias":                    flattenMatchingCriterias(c.MatchingCriterias),
			"is_default_category_enabled":           utils.BoolValue(c.IsDefaultCategoryEnabled),
			"should_keep_default_category_on_login": utils.BoolValue(c.ShouldKeepDefaultCategoryOnLogin),
			"project_ext_id":                        utils.StringValue(c.ProjectExtId),
			"tenant_id":                             utils.StringValue(c.TenantId),
			"links":                                 common.FlattenLinks(c.Links),
		}
		result = append(result, m)
	}
	return result
}
