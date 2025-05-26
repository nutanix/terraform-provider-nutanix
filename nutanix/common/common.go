package common

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// defined to determine whether a particular key (or configuration attribute) within a Terraform resource configuration has been explicitly set by the user.
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

// ExtIdValidation checks if a given string is a valid UUID.
// It uses a regular expression to validate the UUID format.
var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)

func ExtIdValidation(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if !uuidRegex.MatchString(v) {
		errs = append(errs, fmt.Errorf("%q must be a valid Ext ID: got %q", key, v))
	}
	return
}
