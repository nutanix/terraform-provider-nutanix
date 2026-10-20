package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/request/networksecuritypolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixNetworkSecurityPoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceNutanixNetworkSecurityPoliciesV2Read,
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
			"network_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DataSourceNutanixNetworkSecurityPolicyV2(),
			},
		},
	}
}

func DataSourceNutanixNetworkSecurityPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	listNetworkSecurityPoliciesRequest := import2.ListNetworkSecurityPoliciesRequest{}

	if v, ok := d.GetOk("page"); ok {
		listNetworkSecurityPoliciesRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listNetworkSecurityPoliciesRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listNetworkSecurityPoliciesRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listNetworkSecurityPoliciesRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listNetworkSecurityPoliciesRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.NetworkingSecurityInstance.ListNetworkSecurityPolicies(ctx, &listNetworkSecurityPoliciesRequest)
	if err != nil {
		return diag.Errorf("error while fetching network security policy: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("network_policies", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of network security policies.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.NetworkSecurityPolicy)
	if err := d.Set("network_policies", flattenNetworkSecurityPolicy(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenNetworkSecurityPolicy(pr []import1.NetworkSecurityPolicy) []interface{} {
	if len(pr) > 0 {
		nets := make([]interface{}, len(pr))

		for k, v := range pr {
			net := make(map[string]interface{})

			net["ext_id"] = v.ExtId
			net["name"] = v.Name
			net["type"] = common.FlattenPtrEnum(v.Type)
			net["description"] = v.Description
			net["state"] = common.FlattenPtrEnum(v.State)
			net["rules"] = flattenNetworkSecurityPolicyRule(v.Rules)
			net["is_ipv6_traffic_allowed"] = v.IsIpv6TrafficAllowed
			net["is_hitlog_enabled"] = v.IsHitlogEnabled
			if v.Scope != nil {
				net["scope"] = common.FlattenPtrEnum(v.Scope)
			}
			if v.VpcReferences != nil {
				net["vpc_reference"] = v.VpcReferences
			}
			if v.SecuredGroups != nil {
				net["secured_groups"] = v.SecuredGroups
			}
			if v.LastUpdateTime != nil {
				t := v.LastUpdateTime
				net["last_update_time"] = t.String()
			}
			if v.CreationTime != nil {
				t := v.CreationTime
				net["creation_time"] = t.String()
			}
			net["is_system_defined"] = v.IsSystemDefined
			net["created_by"] = v.CreatedBy

			if v.TenantId != nil {
				net["tenant_id"] = v.TenantId
			}
			if v.Links != nil {
				net["links"] = flattenLinksMicroSeg(v.Links)
			}
			net["project_ext_id"] = v.ProjectExtId
			net["is_shared_with_all_projects"] = utils.BoolValue(v.IsSharedWithAllProjects)

			nets[k] = net
		}
		return nets
	}
	return nil
}
