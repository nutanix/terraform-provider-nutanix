package nutanix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func categoriesMappingSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func expandCategoriesMapping(cm interface{}) map[string][]string {
	categoriesMap := cm.(*schema.Set).List()
	output := make(map[string][]string)

	for _, v := range categoriesMap {
		category := v.(map[string]interface{})
		output[category["name"].(string)] = expandMapValues((category["value"].([]interface{})))
	}
	return output
}

func flattenCategoriesMapping(categories map[string][]string) []interface{} {
	c := make([]interface{}, 0)

	for name, value := range categories {
		c = append(c, map[string]interface{}{
			"name":  name,
			"value": value,
		})
	}

	return c
}

func expandMapValues(pr []interface{}) (out []string) {
	for _, v := range pr {
		out = append(out, v.(string))
	}
	return out
}
