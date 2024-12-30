package datapoliciesv2

func expandListOfString(list []interface{}) []string {
	stringListStr := make([]string, len(list))
	for i, v := range list {
		stringListStr[i] = v.(string)
	}
	return stringListStr
}
