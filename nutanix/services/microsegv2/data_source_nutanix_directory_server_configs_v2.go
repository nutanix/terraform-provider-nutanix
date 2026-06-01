package microsegv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/request/directoryserverconfigs"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

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
						"directory_service_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_controllers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForIPAddressOrFQDNComputed(),
						},
						"is_default_category_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"matching_criterias": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForMatchingCriteriaComputed(),
						},
						"should_keep_default_category_on_login": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"project_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixDirectoryServerConfigsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	req := import3.ListDirectoryServerConfigsRequest{}

	if v, ok := d.GetOk("select"); ok {
		req.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.DirectoryServerConfigsAPIInstance.ListDirectoryServerConfigs(ctx, &req)
	if err != nil {
		return diag.Errorf("error while listing Directory Server Configs: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("directory_server_configs", []map[string]interface{}{}); err != nil {
			return diag.Errorf("error setting Directory Server Configs: %s", err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of directory server configs.",
		}}
	}

	getResp := resp.Data.GetValue().([]import2.DirectoryServerConfig)

	if err := d.Set("directory_server_configs", flattenDirectoryServerConfigs(getResp)); err != nil {
		return diag.Errorf("error setting Directory Server Configs: %s", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}
