package datapoliciesv2
import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/datapolicies/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixStoragePolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixStoragePolicyV2Read,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category_ext_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"compression_spec": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compression_state": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"DISABLED", "POSTPROCESS", "INLINE", "SYSTEM_DERIVED"}, false),
						},
					},
				},
			},
			"encryption_spec": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encryption_state": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"SYSTEM_DERIVED", "ENABLED"}, false),
						},
					},
				},
			},
			"qos_spec": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"throttled_iops": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"fault_tolerance_spec": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replication_factor": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Default:  "USER",
			},
		},
	}
}

func dataSourceNutanixStoragePolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	resp, err := conn.StoragePolicies.GetStoragePolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while reading Storage Policy: %v", err)
	}
	body := resp.Data.GetValue().(import1.StoragePolicy)
	return commonReadStateStoragePolicy(ctx, d, meta, body)
}