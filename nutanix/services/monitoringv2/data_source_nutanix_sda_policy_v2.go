package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixSdaPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixSdaPolicyV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique ID of the System-Defined Alert Policy.",
			},
			"classifications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Various categories into which this alert type can be classified.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"cluster_configs": schemaForClusterConfig(),
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "System-defined alert policy description.",
			},
			"entity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"impact_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Impact types to which this rule applies.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"kb_articles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of knowledge base article links.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"links": linksSchema(),
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the System-Defined Alert Policy.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID associated with the policy.",
			},
			"publisher": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Publisher of the policy.",
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Indicates the cluster type against which this policy can be executed.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Title of a System-Defined Alert Policy.",
			},
			"sda_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixSdaPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.SystemDefinedPoliciesAPI.GetSdaPolicyById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching System-Defined Alert Policy: %v", err)
	}

	getResp, ok := resp.Data.GetValue().(serviceability.SystemDefinedPolicy)
	if !ok {
		return diag.Errorf("error: unexpected response type from GetSdaPolicyById, expected SystemDefinedPolicy")
	}

	if err := flattenSdaPolicy(d, &getResp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenSdaPolicy(d *schema.ResourceData, p *serviceability.SystemDefinedPolicy) error {
	if err := d.Set("ext_id", p.ExtId); err != nil {
		return err
	}
	if err := d.Set("name", p.Name); err != nil {
		return err
	}
	if err := d.Set("description", p.Description); err != nil {
		return err
	}
	if err := d.Set("title", p.Title); err != nil {
		return err
	}
	if err := d.Set("policy_id", p.PolicyId); err != nil {
		return err
	}
	if err := d.Set("publisher", p.Publisher); err != nil {
		return err
	}
	if err := d.Set("tenant_id", p.TenantId); err != nil {
		return err
	}
	if p.EntityType != nil {
		if err := d.Set("entity_type", p.EntityType.GetName()); err != nil {
			return err
		}
	}
	if p.Scope != nil {
		if err := d.Set("scope", p.Scope.GetName()); err != nil {
			return err
		}
	}
	if p.SubType != nil {
		if err := d.Set("sub_type", p.SubType.GetName()); err != nil {
			return err
		}
	}
	if p.Type != nil {
		if err := d.Set("sda_type", p.Type.GetName()); err != nil {
			return err
		}
	}

	classifications := make([]string, 0)
	for _, c := range p.Classifications {
		classifications = append(classifications, c)
	}
	if err := d.Set("classifications", classifications); err != nil {
		return err
	}

	kbArticles := make([]string, 0)
	for _, k := range p.KbArticles {
		kbArticles = append(kbArticles, k)
	}
	if err := d.Set("kb_articles", kbArticles); err != nil {
		return err
	}

	impactTypes := make([]string, 0)
	for _, it := range p.ImpactTypes {
		impactTypes = append(impactTypes, it.GetName())
	}
	if err := d.Set("impact_types", impactTypes); err != nil {
		return err
	}

	targetClusters := make([]string, 0)
	for _, tc := range p.TargetClusters {
		targetClusters = append(targetClusters, tc.GetName())
	}
	if err := d.Set("target_clusters", targetClusters); err != nil {
		return err
	}

	if err := d.Set("links", flattenLinks(p.Links)); err != nil {
		return err
	}
	if err := d.Set("cluster_configs", flattenClusterConfigs(p.ClusterConfigs)); err != nil {
		return err
	}
	return nil
}
