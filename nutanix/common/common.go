package common

// ExpandListOfString to expand a list of interface{}
// which is defined in the schema
// its return a list of string
func ExpandListOfString(list []interface{}) []string {
	stringListStr := make([]string, len(list))
	for i, v := range list {
		stringListStr[i] = v.(string)
	}
	return stringListStr
}
