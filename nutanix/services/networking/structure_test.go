package networking

import (
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// based on AWS Provider we should change for us
func testConf() map[string]string {
	return map[string]string{
		"listener.#":                   "1",
		"listener.0.lb_port":           "80",
		"listener.0.lb_protocol":       "http",
		"listener.0.instance_port":     "8000",
		"listener.0.instance_protocol": "http",
		"availability_zones.#":         "2",
		"availability_zones.0":         "us-east-1a",
		"availability_zones.1":         "us-east-1b",
		"ingress.#":                    "1",
		"ingress.0.protocol":           "icmp",
		"ingress.0.from_port":          "1",
		"ingress.0.to_port":            "-1",
		"ingress.0.cidr_blocks.#":      "1",
		"ingress.0.cidr_blocks.0":      "0.0.0.0/0",
		"ingress.0.security_groups.#":  "2",
		"ingress.0.security_groups.0":  "sg-11111",
		"ingress.0.security_groups.1":  "foo/sg-22222",
	}
}

func TestExpandStringList(t *testing.T) {
	expanded := utils.Expand(testConf(), "availability_zones").([]interface{})
	stringList := expandStringList(expanded)
	expected := []*string{
		utils.StringPtr("us-east-1a"),
		utils.StringPtr("us-east-1b"),
	}

	if !reflect.DeepEqual(stringList, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			stringList,
			expected)
	}
}

func TestExpandStringListEmptyItems(t *testing.T) {
	initialList := []string{"foo", "bar", "", "baz"}
	l := make([]interface{}, len(initialList))
	for i, v := range initialList {
		l[i] = v
	}
	stringList := expandStringList(l)
	expected := []*string{
		utils.StringPtr("foo"),
		utils.StringPtr("bar"),
		utils.StringPtr("baz"),
	}

	if !reflect.DeepEqual(stringList, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			stringList,
			expected)
	}
}
