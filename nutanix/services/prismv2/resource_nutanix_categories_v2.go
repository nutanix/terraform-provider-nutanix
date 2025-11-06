package prismv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixCategoriesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixCategoriesV2Create,
		ReadContext:   ResourceNutanixCategoriesV2Read,
		UpdateContext: ResourceNutanixCategoriesV2Update,
		DeleteContext: ResourceNutanixCategoriesV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"USER", "INTERNAL", "SYSTEM"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"owner_uuid": {
				Type:     schema.TypeString,
				Optional: true,
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
							Type:     schema.TypeString,
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
		},
	}
}

func ResourceNutanixCategoriesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	input := &import1.Category{}

	if key, ok := d.GetOk("key"); ok {
		input.Key = utils.StringPtr(key.(string))
	}
	if val, ok := d.GetOk("value"); ok {
		input.Value = utils.StringPtr(val.(string))
	}
	if types, ok := d.GetOk("type"); ok {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"USER":     two,
			"SYSTEM":   three,
			"INTERNAL": four,
		}

		pInt := subMap[types.(string)]
		p := import1.CategoryType(pInt.(int))

		input.Type = &p
	}
	if desc, ok := d.GetOk("description"); ok {
		input.Description = utils.StringPtr(desc.(string))
	}
	if ownerUUID, ok := d.GetOk("owner_uuid"); ok {
		input.OwnerUuid = utils.StringPtr(ownerUUID.(string))
	}

	resp, err := conn.CategoriesAPIInstance.CreateCategory(input)
	if err != nil {
		return diag.Errorf("error while creating category: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Category)

	d.SetId(*getResp.ExtId)
	return ResourceNutanixCategoriesV2Read(ctx, d, meta)
}

func ResourceNutanixCategoriesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	resp, err := conn.CategoriesAPIInstance.GetCategoryById(utils.StringPtr(d.Id()), nil)
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
	return nil
}

func ResourceNutanixCategoriesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI
	updatedInput := import1.Category{}
	resp, err := conn.CategoriesAPIInstance.GetCategoryById(utils.StringPtr(d.Id()), nil)
	if err != nil {
		return diag.Errorf("error while fetching categories : %v", err)
	}

	updatedInput = resp.Data.GetValue().(import1.Category)

	if d.HasChange("value") {
		updatedInput.Value = utils.StringPtr(d.Get("value").(string))
	}
	if d.HasChange("description") {
		updatedInput.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("type") {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"USER":     two,
			"SYSTEM":   three,
			"INTERNAL": four,
		}

		pInt := subMap[d.Get("type").(string)]
		p := import1.CategoryType(pInt.(int))
		updatedInput.Type = &p
	}
	if d.HasChange("owner_uuid") {
		updatedInput.OwnerUuid = utils.StringPtr(d.Get("owner_uuid").(string))
	}

	_, er := conn.CategoriesAPIInstance.UpdateCategoryById(utils.StringPtr(d.Id()), &updatedInput)
	if er != nil {
		return diag.Errorf("error while updating categories : %v", err)
	}
	log.Println("[DEBUG] Category updated successfully")
	return ResourceNutanixCategoriesV2Read(ctx, d, meta)
}

func ResourceNutanixCategoriesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	resp, err := conn.CategoriesAPIInstance.DeleteCategoryById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting category : %v", err)
	}

	if resp == nil {
		log.Println("[DEBUG] Category deleted successfully.")
	}

	return nil
}
