package datapoliciesv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Required: true,
			},
			"links": schemaForLinks(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category_ext_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"compression_spec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compression_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"encryption_spec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encryption_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"qos_spec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"throttled_iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"fault_tolerance_spec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replication_factor": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixStoragePolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI
	resp, err := conn.StoragePolicies.GetStoragePolicyById(utils.StringPtr(d.Get("ext_id").(string)))
	if err != nil {
		return diag.Errorf("error while reading Storage Policy: %v", err)
	}
	if resp == nil || resp.Data == nil {
		return diag.Errorf("No Storage Policy found with the given ext_id: %v", d.Get("ext_id").(string))
	}
	body := resp.Data.GetValue().(import1.StoragePolicy)
	metadata := resp.Metadata
	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Get Storage Policy Response: %s", string(aJSON))
	d.SetId(*body.ExtId)
	return commonReadStateStoragePolicy(d, body, metadata)
}
