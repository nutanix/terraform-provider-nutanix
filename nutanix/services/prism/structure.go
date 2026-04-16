package prism

import (
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	// CDROM ...
	CDROM = "CDROM"
)

func expandStringList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, utils.StringPtr(v.(string)))
		}
	}
	return vs
}
