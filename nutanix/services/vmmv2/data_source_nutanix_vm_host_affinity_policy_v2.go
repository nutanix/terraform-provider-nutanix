package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/policies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVMHostAffinityPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVMHostAffinityPolicyV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_updated_by": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_categories": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vm_categories": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func DatasourceNutanixVMHostAffinityPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id")

	resp, err := conn.VMHostAffinityPolicyAPIInstance.GetVmHostAffinityPolicyById(utils.StringPtr(extID.(string)))

	if err != nil {
		return diag.Errorf("error while fetching policy : %v", err)
	}

	getResp := resp.Data.GetValue().(policies.VmHostAffinityPolicy)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreateTime != nil {
		t := getResp.CreateTime
		if err := d.Set("create_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.UpdateTime != nil {
		t := getResp.UpdateTime
		if err := d.Set("update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.CreatedBy != nil {
		createdBy := make(map[string]string)
		if getResp.CreatedBy.ExtId != nil {
			createdBy["ext_id"] = *getResp.CreatedBy.ExtId
			if err := d.Set("created_by", createdBy); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if getResp.LastUpdatedBy != nil {
		updatedBy := make(map[string]string)
		if getResp.LastUpdatedBy.ExtId != nil {
			updatedBy["ext_id"] = *getResp.LastUpdatedBy.ExtId
			if err := d.Set("last_updated_by", updatedBy); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if err := d.Set("host_categories", flattenPolicyCategoryReference(getResp.HostCategories)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_categories", flattenPolicyCategoryReference(getResp.VmCategories)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)

	return nil
}
