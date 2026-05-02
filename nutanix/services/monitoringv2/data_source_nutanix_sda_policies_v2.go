package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixSdaPoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixSdaPoliciesV2Read,
		Schema: map[string]*schema.Schema{
			"sda_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of System-Defined Alert Policies.",
				Elem:        DatasourceNutanixSdaPolicyV2(),
			},
		},
	}
}

func DatasourceNutanixSdaPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	resp, err := conn.SystemDefinedPoliciesAPI.ListSdaPolicies(nil, nil, nil, nil, nil)
	if err != nil {
		return diag.Errorf("error while listing System-Defined Alert Policies: %v", err)
	}

	if resp.Data == nil {
		if setErr := d.Set("sda_policies", []map[string]interface{}{}); setErr != nil {
			return diag.FromErr(setErr)
		}
		d.SetId(utils.GenUUID())
		return nil
	}

	listResp, ok := resp.Data.GetValue().([]serviceability.SystemDefinedPolicy)
	if !ok {
		return diag.Errorf("error: unexpected response type from ListSdaPolicies, expected []SystemDefinedPolicy")
	}

	if len(listResp) == 0 {
		if setErr := d.Set("sda_policies", []map[string]interface{}{}); setErr != nil {
			return diag.FromErr(setErr)
		}
		d.SetId(utils.GenUUID())
		return nil
	}

	policies := make([]map[string]interface{}, len(listResp))
	for i, p := range listResp {
		policy := flattenSdaPolicyToMap(&p)
		policies[i] = policy
	}

	if err := d.Set("sda_policies", policies); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenSdaPolicyToMap(p *serviceability.SystemDefinedPolicy) map[string]interface{} {
	m := map[string]interface{}{}
	if p.ExtId != nil {
		m["ext_id"] = utils.StringValue(p.ExtId)
	}
	if p.Name != nil {
		m["name"] = utils.StringValue(p.Name)
	}
	if p.Description != nil {
		m["description"] = utils.StringValue(p.Description)
	}
	if p.Title != nil {
		m["title"] = utils.StringValue(p.Title)
	}
	if p.PolicyId != nil {
		m["policy_id"] = utils.StringValue(p.PolicyId)
	}
	if p.Publisher != nil {
		m["publisher"] = utils.StringValue(p.Publisher)
	}
	if p.TenantId != nil {
		m["tenant_id"] = utils.StringValue(p.TenantId)
	}
	if p.EntityType != nil {
		m["entity_type"] = p.EntityType.GetName()
	}
	if p.Scope != nil {
		m["scope"] = p.Scope.GetName()
	}
	if p.SubType != nil {
		m["sub_type"] = p.SubType.GetName()
	}
	if p.Type != nil {
		m["sda_type"] = p.Type.GetName()
	}

	classifications := make([]string, 0)
	for _, c := range p.Classifications {
		classifications = append(classifications, c)
	}
	m["classifications"] = classifications

	kbArticles := make([]string, 0)
	for _, k := range p.KbArticles {
		kbArticles = append(kbArticles, k)
	}
	m["kb_articles"] = kbArticles

	impactTypes := make([]string, 0)
	for _, it := range p.ImpactTypes {
		impactTypes = append(impactTypes, it.GetName())
	}
	m["impact_types"] = impactTypes

	targetClusters := make([]string, 0)
	for _, tc := range p.TargetClusters {
		targetClusters = append(targetClusters, tc.GetName())
	}
	m["target_clusters"] = targetClusters

	m["links"] = flattenLinks(p.Links)
	m["cluster_configs"] = flattenClusterConfigs(p.ClusterConfigs)

	return m
}
