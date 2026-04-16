package datapoliciesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/datapolicies/v4/config"
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

	// initialize query params
	var filter, orderBy, selects *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	resp, err := conn.ProtectionPolicies.ListProtectionPolicies(page, limit, filter, orderBy, selects)
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

		protectionPoliciesList = append(protectionPoliciesList, protectionPolicyMap)
	}

	return protectionPoliciesList
}
