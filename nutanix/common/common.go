// Package common provides common utilities for the Nutanix Terraform provider.
package common

import (
	"context"
	"fmt"
	"hash/crc32"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ExpandListOfString to expand a list of interface{}
// which is defined in the schema
// its return a list of string
func ExpandListOfString(list []interface{}) []string {
	stringListStr := make([]string, 0)
	for i, v := range list {
		if v == nil || v == "" {
			log.Printf("[DEBUG] Skipping nil or empty value at index %d", i)
			continue // Skip nil or empty values
		}
		stringListStr = append(stringListStr, v.(string))
	}
	return stringListStr
}

// DiffStringSets compares two string slices and returns the added and removed items.
// added contains items that are in newSet but not in oldSet.
// removed contains items that are in oldSet but not in newSet.
func DiffStringSets(oldSet, newSet []string) (added, removed []string) {
	// Create maps for easier lookup
	oldSetMap := make(map[string]bool, len(oldSet))
	for _, item := range oldSet {
		oldSetMap[item] = true
	}

	newSetMap := make(map[string]bool, len(newSet))
	for _, item := range newSet {
		newSetMap[item] = true
	}

	// Find items to add (in new but not in old)
	added = make([]string, 0)
	for _, item := range newSet {
		if !oldSetMap[item] {
			added = append(added, item)
		}
	}

	// Find items to remove (in old but not in new)
	removed = make([]string, 0)
	for _, item := range oldSet {
		if !newSetMap[item] {
			removed = append(removed, item)
		}
	}

	return added, removed
}

// IsExplicitlySet defined to determine whether a particular key (or configuration attribute) within a Terraform resource configuration has been explicitly set by the user.
// Returns a Boolean (true or false). true indicates that the key was explicitly set with a non-null value; false implies it was either not set, is unknown, or explicitly set to null.
func IsExplicitlySet(d *schema.ResourceData, key string) bool {
	rawConfig := d.GetRawConfig() // Get raw Terraform config as cty.Value
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return false // If rawConfig is null/unknown, key wasn't explicitly set
	}

	// Convert rawConfig to map and check if key exists
	configMap := rawConfig.AsValueMap()
	if val, exists := configMap[key]; exists {
		log.Printf("[DEBUG] Key: %s, Value: %s", key, val)
		return !val.IsNull() // Ensure key exists and isn't explicitly null
	}
	return false
}

func TaskStateRefreshPrismTaskGroupFunc(ctx context.Context, client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vresp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)
		if err != nil {
			return "", "", fmt.Errorf("error while polling prism task: %v", err)
		}

		// get the group results

		v := vresp.Data.GetValue().(prismConfig.Task)

		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
			return v, getTaskStatus(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStatus(v.Status), nil
	}
}

func getTaskStatus(taskStatus *prismConfig.TaskStatus) string {
	return FlattenPtrEnum(taskStatus)
}

// EnumToStrings converts any slice of enums that implement GetName() to []string
func EnumToStrings[T interface{ GetName() string }](enums []T) []string {
	names := make([]string, 0, len(enums))
	for _, e := range enums {
		names = append(names, e.GetName())
	}
	return names
}

// EnumToMap returns a map from string names to enum values for any enum type that implements GetName()
func EnumToMap[T interface{ GetName() string }](enums []T) map[string]T {
	m := make(map[string]T, len(enums))
	for _, e := range enums {
		m[e.GetName()] = e
	}
	return m
}

// ExpandEnum expands a single string value to an enum pointer
func ExpandEnum[T any](val interface{}, enumMap map[string]T, fieldName string) *T {
	if val == nil {
		return nil
	}

	if str, ok := val.(string); ok && str != "" {
		if enumVal, found := enumMap[str]; found {
			return &enumVal
		} else {
			log.Printf("[WARN] unknown %s: %s", fieldName, str)
		}
	}

	return nil
}

// ExpandEnumList expands a list of strings to enum values
func ExpandEnumList[T any](val interface{}, enumMap map[string]T, fieldName string) []T {
	if val == nil {
		return nil
	}

	list := make([]T, 0)
	for _, item := range val.([]interface{}) {
		if str, ok := item.(string); ok && str != "" {
			if enumVal, found := enumMap[str]; found {
				list = append(list, enumVal)
			} else {
				log.Printf("[WARN] unknown %s: %s", fieldName, str)
			}
		}
	}

	if len(list) == 0 {
		return nil
	}
	return list
}

// FlattenPtrEnum converts a pointer to an enum implementing GetName() to a string.
// Returns "" if the pointer is nil.
func FlattenPtrEnum[T interface{ GetName() string }](enumPtr *T) string {
	if enumPtr == nil {
		return ""
	}
	return (*enumPtr).GetName()
}

// FlattenEnumPtrList converts a slice of pointers to enums into a slice of strings.
// Nil pointers are skipped.
func FlattenEnumPtrList[T interface{ GetName() string }](enums []*T) []string {
	if len(enums) == 0 {
		return nil
	}

	names := make([]string, 0, len(enums))
	for _, e := range enums {
		if e != nil {
			names = append(names, (*e).GetName())
		}
	}
	return names
}

// FlattenEnumValueList converts a slice of enum values (not pointers) to []string
func FlattenEnumValueList[T interface{ GetName() string }](enums []T) []string {
	if len(enums) == 0 {
		return nil
	}

	names := make([]string, 0, len(enums))
	for _, e := range enums {
		names = append(names, e.GetName())
	}
	return names
}

// InterfaceToSlice converts various input types to a slice of interfaces.
func InterfaceToSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}

	switch t := v.(type) {
	case *schema.Set:
		return t.List()
	case []interface{}:
		return t
	case []map[string]interface{}:
		// unlikely, but handle it
		out := make([]interface{}, len(t))
		for i := range t {
			out[i] = t[i]
		}
		return out
	default:
		// single element provided
		return []interface{}{v}
	}
}

// HashStringItem returns a hash for a string value to ensure uniqueness in schema.TypeSet
func HashStringItem(v interface{}) int {
	if v == nil {
		return 0
	}
	str, ok := v.(string)
	if !ok {
		return 0
	}
	return int(crc32.ChecksumIEEE([]byte(fmt.Sprintf("%s-", str))))
}
