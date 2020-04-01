package nutanix

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

func resourceNutanixCategoriesMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found Nutanix State v0; migrating to v1")
		return migrateNutanixCategoriesV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateNutanixCategoriesV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	keys := make([]string, 0, len(is.Attributes))
	for k := range is.Attributes {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	if l, ok := is.Attributes["categories.%"]; ok && l != "" {
		var tempCat []map[string]string

		for _, k := range keys {
			v := is.Attributes[k]

			if k == "categories.%" {
				is.Attributes["categories.#"] = v
				delete(is.Attributes, "categories.%")
				continue
			}

			if strings.HasPrefix(k, "categories.") {
				path := strings.Split(k, ".")
				if len(path) != 2 {
					return is, fmt.Errorf("found unexpected categories field: %#v", k)
				}

				if path[1] == "#" {
					continue
				}

				log.Printf("[DEBUG] key=%s", k)

				tempCat = append(tempCat, map[string]string{
					"name":  path[1],
					"value": v,
				})

				delete(is.Attributes, k)
			}
		}
		flattenTempCategories(tempCat, is)
	}
	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}

func flattenTempCategories(categories []map[string]string, is *terraform.InstanceState) {
	for index, category := range categories {
		for key, value := range category {
			is.Attributes[fmt.Sprintf("categories.%d.%s", index, key)] = value
		}
	}
}
