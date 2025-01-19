package datapoliciesv2

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"reflect"
	"sort"
)

func expandListOfString(list []interface{}) []string {
	stringListStr := make([]string, len(list))
	for i, v := range list {
		stringListStr[i] = v.(string)
	}
	return stringListStr
}

func categoryIdsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {

	log.Printf("[DEBUG] DiffSuppressFunc for category_ids")

	if d.HasChange("category_ids") {
		oldCap, newCap := d.GetChange("category_ids")
		log.Printf("[DEBUG] oldCap : %v", oldCap)
		log.Printf("[DEBUG] newCap : %v", newCap)

		oldList := oldCap.([]interface{})
		newList := newCap.([]interface{})

		if len(oldList) != len(newList) {
			log.Printf("[DEBUG] category_ids are different")
			return false
		}

		sort.SliceStable(oldList, func(i, j int) bool {
			return oldList[i].(string) < oldList[j].(string)
		})
		sort.SliceStable(newList, func(i, j int) bool {
			return newList[i].(string) < newList[j].(string)
		})

		aJSON, _ := json.Marshal(oldList)
		log.Printf("[DEBUG] oldList : %s", aJSON)
		aJSON, _ = json.Marshal(newList)
		log.Printf("[DEBUG] newList : %s", aJSON)

		if reflect.DeepEqual(oldList, newList) {
			log.Printf("[DEBUG] category_ids are same")
			return true
		}
		log.Printf("[DEBUG] category_ids are different")
		return false
	}
	return false
}
