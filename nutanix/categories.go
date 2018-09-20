package nutanix

import "github.com/hashicorp/terraform/helper/schema"

func categoriesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
	}
}

func expandCategories(categories map[string]interface{}) map[string]string {
	output := make(map[string]string, len(categories))

	for i, v := range categories {
		output[i] = v.(string)
	}

	return output
}
