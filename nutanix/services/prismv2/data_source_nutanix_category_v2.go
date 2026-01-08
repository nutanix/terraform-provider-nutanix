package prismv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixCategoryV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixCategoryV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"detailed_associations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixCategoryV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	extID := d.Get("ext_id")
	var expand *string
	if expandf, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandf.(string))
	} else {
		expand = nil
	}
	resp, err := conn.CategoriesAPIInstance.GetCategoryById(utils.StringPtr(extID.(string)), expand)
	if err != nil {
		return diag.Errorf("error while fetching category : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Category)

	if err := d.Set("key", getResp.Key); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("value", getResp.Value); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", flattenCategoryType(getResp.Type)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_uuid", getResp.OwnerUuid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("associations", flattenAssociationSummary(getResp.Associations)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("detailed_associations", flattenAssociationDetail(getResp.DetailedAssociations)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenCategoryType(pr *import1.CategoryType) string {
	const two, three, four = 2, 3, 4
	if pr != nil {
		if *pr == import1.CategoryType(two) {
			return "USER"
		}
		if *pr == import1.CategoryType(three) {
			return "SYSTEM"
		}
		if *pr == import1.CategoryType(four) {
			return "INTERNAL"
		}
	}
	return "UNKNOWN"
}

func flattenAssociationSummary(pr []import1.AssociationSummary) []interface{} {
	if len(pr) > 0 {
		associationList := make([]interface{}, len(pr))

		for k, v := range pr {
			assn := make(map[string]interface{})

			assn["category_id"] = v.CategoryId
			assn["count"] = v.Count
			assn["resource_group"] = flattenResourceGroup(v.ResourceGroup)
			assn["resource_type"] = flattenResourceType(v.ResourceType)

			associationList[k] = assn
		}
		return associationList
	}
	return nil
}

func flattenAssociationDetail(pr []import1.AssociationDetail) []interface{} {
	if len(pr) > 0 {
		detailList := make([]interface{}, len(pr))

		for k, v := range pr {
			detail := make(map[string]interface{})

			detail["category_id"] = v.CategoryId
			detail["resource_group"] = flattenResourceGroup(v.ResourceGroup)
			detail["resource_type"] = flattenResourceType(v.ResourceType)
			detail["resource_id"] = v.ResourceId

			detailList[k] = detail
		}
		return detailList
	}
	return nil
}

func flattenResourceGroup(pr *import1.ResourceGroup) string {
	const two, three = 2, 3
	if pr != nil {
		if *pr == import1.ResourceGroup(two) {
			return "ENTITY"
		}
		if *pr == import1.ResourceGroup(three) {
			return "POLICY"
		}
	}
	return "UNKNOWN"
}

func flattenResourceType(pr *import1.ResourceType) string {
	const (
		two, three, four, five, six, seven, eight, nine, ten, eleven, twelve, thirteen, fourteen, fifteen,
		sixteen, seventeen, eighteen, nineteen, twenty, twentyone, twentytwo, twentythree, twentyfour, twentyfive,
		twentysix = 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26
	)
	if pr != nil {
		if *pr == import1.ResourceType(two) {
			return "VM"
		}
		if *pr == import1.ResourceType(three) {
			return "MH_VM"
		}
		if *pr == import1.ResourceType(four) {
			return "IMAGE"
		}
		if *pr == import1.ResourceType(five) {
			return "SUBNET"
		}
		if *pr == import1.ResourceType(six) {
			return "CLUSTER"
		}
		if *pr == import1.ResourceType(seven) {
			return "HOST"
		}
		if *pr == import1.ResourceType(eight) {
			return "REPORT"
		}
		if *pr == import1.ResourceType(nine) {
			return "MARKETPLACE_ITEM"
		}
		if *pr == import1.ResourceType(ten) {
			return "BLUEPRINT"
		}
		if *pr == import1.ResourceType(eleven) {
			return "APP"
		}
		if *pr == import1.ResourceType(twelve) {
			return "VOLUMEGROUP"
		}
		if *pr == import1.ResourceType(thirteen) {
			return "IMAGE_PLACEMENT_POLICY"
		}
		if *pr == import1.ResourceType(fourteen) {
			return "NETWORK_SECURITY_POLICY"
		}
		if *pr == import1.ResourceType(fifteen) {
			return "NETWORK_SECURITY_RULE"
		}
		if *pr == import1.ResourceType(sixteen) {
			return "VM_HOST_AFFINITY_POLICY"
		}
		if *pr == import1.ResourceType(seventeen) {
			return "VM_VM_ANTI_AFFINITY_POLICY"
		}
		if *pr == import1.ResourceType(eighteen) {
			return "QOS_POLICY"
		}
		if *pr == import1.ResourceType(nineteen) {
			return "NGT_POLICY"
		}
		if *pr == import1.ResourceType(twenty) {
			return "PROTECTION_RULE"
		}
		if *pr == import1.ResourceType(twentyone) {
			return "ACCESS_CONTROL_POLICY"
		}
		if *pr == import1.ResourceType(twentytwo) {
			return "STORAGE_POLICY"
		}
		if *pr == import1.ResourceType(twentythree) {
			return "IMAGE_RATE_LIMIT"
		}
		if *pr == import1.ResourceType(twentyfour) {
			return "RECOVERY_PLAN"
		}
		if *pr == import1.ResourceType(twentyfive) {
			return "BUNDLE"
		}
		if *pr == import1.ResourceType(twentysix) {
			return "POLICY_SCHEMA"
		}
	}
	return "UNKNOWN"
}

func flattenLinks(pr []import2.ApiLink) []map[string]interface{} {
	if len(pr) > 0 {
		linkList := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			links := map[string]interface{}{}
			if v.Href != nil {
				links["href"] = v.Href
			}
			if v.Rel != nil {
				links["rel"] = v.Rel
			}

			linkList[k] = links
		}
		return linkList
	}
	return nil
}
