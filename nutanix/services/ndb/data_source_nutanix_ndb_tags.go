package ndb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixNDBTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBTagsRead,
		Schema: map[string]*schema.Schema{
			"entity_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DATABASE", "TIME_MACHINE",
					"CLONE", "DATABASE_SERVER",
				}, false),
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"entity_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"date_created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixNDBTagsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	entityType := ""
	if entity, eok := d.GetOk("entity_type"); eok {
		entityType = entity.(string)
	}

	resp, err := conn.Service.ListTags(ctx, entityType)
	if err != nil {
		return diag.FromErr(err)
	}

	if e := d.Set("tags", flattenTagsList(resp)); e != nil {
		return diag.FromErr(e)
	}

	uuid, er := uuid.GenerateUUID()
	if er != nil {
		return diag.Errorf("Error generating UUID for era tags: %+v", er)
	}
	d.SetId(uuid)
	return nil
}

func flattenTagsList(pr *era.ListTagsResponse) []interface{} {
	if pr != nil {
		tagsList := make([]interface{}, 0)

		for _, v := range *pr {
			tag := map[string]interface{}{}

			tag["id"] = v.ID
			tag["name"] = v.Name
			tag["description"] = v.Description
			tag["required"] = v.Required
			tag["entity_type"] = v.EntityType
			tag["status"] = v.Status
			tag["owner"] = v.Owner
			tag["values"] = v.Values
			tag["date_created"] = v.DateCreated
			tag["date_modified"] = v.DateModified

			tagsList = append(tagsList, tag)
		}
		return tagsList
	}
	return nil
}
