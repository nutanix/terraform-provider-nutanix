package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/policies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVMHostAffinityPoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVMHostAffinityPoliciesV2Read,
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
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixVMHostAffinityPolicyV2(),
			},
		},
	}
}

func DatasourceNutanixVMHostAffinityPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	var filter, orderBy *string
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

	resp, err := conn.VMHostAffinityPolicyAPIInstance.ListVmHostAffinityPolicies(page, limit, filter, orderBy)

	if err != nil {
		return diag.Errorf("error while fetching policies : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("policies", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of policies.",
		}}
	}

	getResp := resp.Data.GetValue().([]policies.VmHostAffinityPolicy)

	if err := d.Set("policies", flattenVMHostAffinityPolicyEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenVMHostAffinityPolicyEntities(pr []policies.VmHostAffinityPolicy) []interface{} {
	if len(pr) == 0 {
		return nil
	}
	policies := make([]interface{}, len(pr))

	for k, v := range pr {
		policies[k] = flattenVMHostAffinityPolicyEntity(v)
	}
	return policies
}
