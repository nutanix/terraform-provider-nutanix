package utils

// IsStringSetAndNotEmpty checks if a value is a non-empty string.
// Returns the string value and true if the value is a string and not empty,
// otherwise returns empty string and false.
//
// This is useful for conditionally setting optional string fields that should
// only be included when they have a non-empty value.
// Usage: if value, ok := utils.IsStringSetAndNotEmpty(d.Get("key")); ok { ... }
func IsStringSetAndNotEmpty(val interface{}) (string, bool) {
	if val == nil {
		return "", false
	}
	strVal, ok := val.(string)
	if !ok {
		return "", false
	}
	if strVal == "" {
		return "", false
	}
	return strVal, true
}
