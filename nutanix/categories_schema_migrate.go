package nutanix

import (
	"fmt"
	"log"
	"sort"
)

func resourceNutanixCategoriesMigrateState(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if len(rawState) == 0 || rawState == nil {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return rawState, nil
	}

	keys := make([]string, 0, len(rawState))
	for k := range rawState {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	log.Printf("[DEBUG] Attributes before migration: %#v", rawState)

	if l, ok := rawState["categories"]; ok {
		if asserted_l, ok := l.(map[string]interface{}); ok {
			c := make([]interface{}, 0)
			keys := make([]string, 0, len(asserted_l))
			for k := range asserted_l {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, name := range keys {
				value := asserted_l[name]
				c = append(c, map[string]interface{}{
					"name":  name,
					"value": value.(string),
				})
			}
			rawState["categories"] = c
		}
	}
	log.Printf("[DEBUG] Attributes after migration: %#v", rawState)
	return rawState, nil
}

func flattenTempCategories(categories []map[string]string, rawState map[string]interface{}) {
	for index, category := range categories {
		for key, value := range category {
			rawState[fmt.Sprintf("categories.%d.%s", index, key)] = value
		}
	}
}
