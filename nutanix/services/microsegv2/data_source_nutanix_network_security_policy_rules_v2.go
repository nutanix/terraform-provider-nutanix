package microsegv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	config "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixNetworkSecurityPolicyRulesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceNutanixNetworkSecurityPolicyRulesV2Read,
		Schema: map[string]*schema.Schema{
			"policy_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ExtId of the network security policy to list rules for.",
			},
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page number for pagination.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of rules to return.",
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter expression for the list.",
			},
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Order by clause.",
			},
			"select": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Select fields to return.",
			},
			"network_security_policy_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of network security policy rules.",
				Elem: &schema.Resource{
					Schema: networkSecurityPolicyRuleSchema(),
				},
			},
		},
	}
}

func networkSecurityPolicyRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ext_id":      {Type: schema.TypeString, Computed: true},
		"description": {Type: schema.TypeString, Computed: true},
		"tenant_id":   {Type: schema.TypeString, Computed: true},
		"type":        {Type: schema.TypeString, Computed: true},
		"links":       {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{"href": {Type: schema.TypeString, Computed: true}, "rel": {Type: schema.TypeString, Computed: true}}}},
		"spec":        {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: ruleSpecSchema()}},
	}
}

func ruleSpecSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"two_env_isolation_rule_spec": {
			Type: schema.TypeList, Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"first_isolation_group":  {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
					"second_isolation_group": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
				},
			},
		},
		"application_rule_spec": {
			Type: schema.TypeList, Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"secured_group_category_references": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
					"src_allow_spec":                    {Type: schema.TypeString, Computed: true},
					"dest_allow_spec":                   {Type: schema.TypeString, Computed: true},
					"src_category_references":           {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
					"dest_category_references":          {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
					"src_subnet":                        {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{"value": {Type: schema.TypeString, Computed: true}, "prefix_length": {Type: schema.TypeInt, Computed: true}}}},
					"dest_subnet":                       {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{"value": {Type: schema.TypeString, Computed: true}, "prefix_length": {Type: schema.TypeInt, Computed: true}}}},
					"src_address_group_references":      {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
					"dest_address_group_references":     {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
					"service_group_references":          {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
					"is_all_protocol_allowed":           {Type: schema.TypeBool, Computed: true},
					"tcp_services":                      {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{"start_port": {Type: schema.TypeInt, Computed: true}, "end_port": {Type: schema.TypeInt, Computed: true}}}},
					"udp_services":                      {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{"start_port": {Type: schema.TypeInt, Computed: true}, "end_port": {Type: schema.TypeInt, Computed: true}}}},
					"icmp_services":                     {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{"is_all_allowed": {Type: schema.TypeBool, Computed: true}, "type": {Type: schema.TypeInt, Computed: true}, "code": {Type: schema.TypeInt, Computed: true}}}},
					"network_function_chain_reference":  {Type: schema.TypeString, Computed: true},
				},
			},
		},
		"intra_entity_group_rule_spec": {
			Type: schema.TypeList, Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"secured_group_action":              {Type: schema.TypeString, Computed: true},
					"secured_group_category_references": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
				},
			},
		},
		"multi_env_isolation_rule_spec": {
			Type: schema.TypeList, Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"spec": {
						Type: schema.TypeList, Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"all_to_all_isolation_group": {
									Type: schema.TypeList, Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"isolation_group": {
												Type: schema.TypeList, Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"group_category_references": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func DataSourceNutanixNetworkSecurityPolicyRulesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	policyExtID := d.Get("policy_ext_id").(string)
	var page, limit *int
	var filter, orderBy, select_ *string

	if v, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.NetworkingSecurityInstance.ListNetworkSecurityPolicyRules(
		utils.StringPtr(policyExtID), page, limit, filter, orderBy, select_)
	if err != nil {
		return diag.Errorf("error listing network security policy rules: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("network_security_policy_rules", []interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of network security policy rules.",
		}}
	}

	// Response Data is OneOfListNetworkSecurityPolicyRulesApiResponseData; oneOfType0 is []NetworkSecurityPolicyRule
	value := resp.Data.GetValue()
	rules, ok := value.([]import1.NetworkSecurityPolicyRule)
	if !ok {
		if err := d.Set("network_security_policy_rules", []interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return nil
	}

	if err := d.Set("network_security_policy_rules", flattenNetworkSecurityPolicyRules(rules)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resource.UniqueId())
	return nil
}

func flattenNetworkSecurityPolicyRules(rules []import1.NetworkSecurityPolicyRule) []interface{} {
	if len(rules) == 0 {
		return nil
	}
	out := make([]interface{}, len(rules))
	for i, r := range rules {
		m := map[string]interface{}{
			"ext_id":      utils.StringValue(r.ExtId),
			"description": utils.StringValue(r.Description),
			"tenant_id":   utils.StringValue(r.TenantId),
			"type":        flattenRuleTypeMicroseg(r.Type),
			"links":       flattenLinksMicroseg(r.Links),
			"spec":        flattenOneOfNetworkSecurityPolicyRuleSpecMicroseg(r.Spec),
		}
		out[i] = m
	}
	return out
}

func flattenLinksMicroseg(links []import2.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, len(links))
	for i, l := range links {
		out[i] = map[string]interface{}{
			"href": l.Href,
			"rel":  l.Rel,
		}
	}
	return out
}

func flattenRuleTypeMicroseg(t *import1.RuleType) string {
	if t == nil {
		return ""
	}
	return t.GetName()
}

func flattenOneOfNetworkSecurityPolicyRuleSpecMicroseg(pr *import1.OneOfNetworkSecurityPolicyRuleSpec) []map[string]interface{} {
	if pr == nil || pr.ObjectType_ == nil {
		return nil
	}
	value := pr.GetValue()
	if value == nil {
		return nil
	}
	switch *pr.ObjectType_ {
	case "microseg.v4.config.TwoEnvIsolationRuleSpec":
		v, ok := value.(import1.TwoEnvIsolationRuleSpec)
		if !ok {
			return nil
		}
		return []map[string]interface{}{{
			"two_env_isolation_rule_spec": []map[string]interface{}{{
				"first_isolation_group":  v.FirstIsolationGroup,
				"second_isolation_group": v.SecondIsolationGroup,
			}},
		}}
	case "microseg.v4.config.ApplicationRuleSpec":
		v, ok := value.(import1.ApplicationRuleSpec)
		if !ok {
			return nil
		}
		app := map[string]interface{}{
			"secured_group_category_references": v.SecuredGroupCategoryReferences,
			"src_category_references":           v.SrcCategoryReferences,
			"dest_category_references":          v.DestCategoryReferences,
			"src_address_group_references":      v.SrcAddressGroupReferences,
			"dest_address_group_references":     v.DestAddressGroupReferences,
			"service_group_references":          v.ServiceGroupReferences,
			"network_function_chain_reference":  v.NetworkFunctionChainReference,
		}
		if v.SrcAllowSpec != nil {
			app["src_allow_spec"] = (*v.SrcAllowSpec).GetName()
		}
		if v.DestAllowSpec != nil {
			app["dest_allow_spec"] = (*v.DestAllowSpec).GetName()
		}
		if v.SrcSubnet != nil {
			app["src_subnet"] = flattenIPv4Microseg(v.SrcSubnet)
		}
		if v.DestSubnet != nil {
			app["dest_subnet"] = flattenIPv4Microseg(v.DestSubnet)
		}
		if v.IsAllProtocolAllowed != nil {
			app["is_all_protocol_allowed"] = *v.IsAllProtocolAllowed
		}
		if v.TcpServices != nil {
			app["tcp_services"] = flattenTCPPortRangeMicroseg(v.TcpServices)
		}
		if v.UdpServices != nil {
			app["udp_services"] = flattenUDPPortRangeMicroseg(v.UdpServices)
		}
		if v.IcmpServices != nil {
			app["icmp_services"] = flattenIcmpTypeCodeMicroseg(v.IcmpServices)
		}
		return []map[string]interface{}{{"application_rule_spec": []map[string]interface{}{app}}}
	case "microseg.v4.config.IntraEntityGroupRuleSpec":
		v, ok := value.(import1.IntraEntityGroupRuleSpec)
		if !ok {
			return nil
		}
		intra := map[string]interface{}{
			"secured_group_category_references": v.SecuredGroupCategoryReferences,
		}
		if v.SecuredGroupAction != nil {
			intra["secured_group_action"] = v.SecuredGroupAction.GetName()
		}
		return []map[string]interface{}{{"intra_entity_group_rule_spec": []map[string]interface{}{intra}}}
	case "microseg.v4.config.MultiEnvIsolationRuleSpec":
		v, ok := value.(import1.MultiEnvIsolationRuleSpec)
		if !ok || v.Spec == nil {
			return nil
		}
		specVal := v.Spec.GetValue()
		allToAll, ok := specVal.(import1.AllToAllIsolationGroup)
		if !ok {
			return nil
		}
		isolationGroups := make([]interface{}, 0, len(allToAll.IsolationGroups))
		for _, g := range allToAll.IsolationGroups {
			isolationGroups = append(isolationGroups, map[string]interface{}{
				"group_category_references": g.GroupCategoryReferences,
			})
		}
		allToAllGroup := map[string]interface{}{"isolation_group": isolationGroups}
		specMap := map[string]interface{}{
			"all_to_all_isolation_group": []interface{}{allToAllGroup},
		}
		multiEnv := map[string]interface{}{"spec": []interface{}{specMap}}
		return []map[string]interface{}{{
			"multi_env_isolation_rule_spec": []interface{}{multiEnv},
		}}
	}
	return nil
}

func flattenIPv4Microseg(pr *config.IPv4Address) []interface{} {
	if pr == nil {
		return nil
	}
	m := map[string]interface{}{"value": utils.StringValue(pr.Value), "prefix_length": utils.IntValue(pr.PrefixLength)}
	return []interface{}{m}
}

func flattenTCPPortRangeMicroseg(pr []import1.TcpPortRangeSpec) []interface{} {
	if len(pr) == 0 {
		return nil
	}
	out := make([]interface{}, len(pr))
	for i, p := range pr {
		out[i] = map[string]interface{}{"start_port": utils.IntValue(p.StartPort), "end_port": utils.IntValue(p.EndPort)}
	}
	return out
}

func flattenUDPPortRangeMicroseg(pr []import1.UdpPortRangeSpec) []interface{} {
	if len(pr) == 0 {
		return nil
	}
	out := make([]interface{}, len(pr))
	for i, p := range pr {
		out[i] = map[string]interface{}{"start_port": utils.IntValue(p.StartPort), "end_port": utils.IntValue(p.EndPort)}
	}
	return out
}

func flattenIcmpTypeCodeMicroseg(pr []import1.IcmpTypeCodeSpec) []interface{} {
	if len(pr) == 0 {
		return nil
	}
	out := make([]interface{}, len(pr))
	for i, p := range pr {
		out[i] = map[string]interface{}{
			"is_all_allowed": utils.BoolValue(p.IsAllAllowed),
			"type":           utils.IntValue(p.Type),
			"code":           utils.IntValue(p.Code),
		}
	}
	return out
}
