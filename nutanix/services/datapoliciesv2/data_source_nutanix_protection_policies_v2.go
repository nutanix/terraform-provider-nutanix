package datapoliciesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/datapolicies-go-client/v17/models/datapolicies/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/datapolicies-go-client/v17/models/datapolicies/v4/request/protectionpolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixProtectionPoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixProtectionPoliciesV2Read,
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
			"protection_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixProtectionPolicyV2(),
			},
		},
	}
}

func DatasourceNutanixProtectionPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	listProtectionPoliciesRequest := import1.ListProtectionPoliciesRequest{}

	if v, ok := d.GetOk("page"); ok {
		listProtectionPoliciesRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listProtectionPoliciesRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listProtectionPoliciesRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listProtectionPoliciesRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listProtectionPoliciesRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.ProtectionPolicies.ListProtectionPolicies(ctx, &listProtectionPoliciesRequest)
	if err != nil {
		return diag.Errorf("error while Listing Protection Policies: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("protection_policies", []map[string]interface{}{}); err != nil {
			return diag.Errorf("error setting Protection Policies: %s", err)
		}
		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of protection policies.",
		}}
	}

	getResp := resp.Data.GetValue().([]config.ProtectionPolicy)

	if err := d.Set("protection_policies", flattenProtectionPolicies(getResp)); err != nil {
		return diag.Errorf("error setting Protection Policies: %s", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenProtectionPolicies(protectionPolicies []config.ProtectionPolicy) []map[string]interface{} {
	if len(protectionPolicies) == 0 {
		return []map[string]interface{}{}
	}

	protectionPoliciesList := make([]map[string]interface{}, 0)

	for _, protectionPolicy := range protectionPolicies {
		protectionPolicyMap := make(map[string]interface{})

		protectionPolicyMap["tenant_id"] = utils.StringValue(protectionPolicy.TenantId)
		protectionPolicyMap["ext_id"] = protectionPolicy.ExtId
		protectionPolicyMap["links"] = flattenLinks(protectionPolicy.Links)
		protectionPolicyMap["name"] = utils.StringValue(protectionPolicy.Name)
		protectionPolicyMap["description"] = utils.StringValue(protectionPolicy.Description)
		protectionPolicyMap["replication_locations"] = flattenReplicationLocations(protectionPolicy.ReplicationLocations)
		protectionPolicyMap["replication_configurations"] = flattenReplicationConfigurations(protectionPolicy.ReplicationConfigurations)
		protectionPolicyMap["category_ids"] = protectionPolicy.CategoryIds
		protectionPolicyMap["is_approval_policy_needed"] = utils.BoolValue(protectionPolicy.IsApprovalPolicyNeeded)
		protectionPolicyMap["owner_ext_id"] = utils.StringValue(protectionPolicy.OwnerExtId)
		protectionPolicyMap["project_ext_id"] = protectionPolicy.ProjectExtId

		protectionPoliciesList = append(protectionPoliciesList, protectionPolicyMap)
	}

	return protectionPoliciesList
}
