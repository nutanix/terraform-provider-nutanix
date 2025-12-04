package prism

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var requiredResourceFields map[string]map[string][]string = map[string]map[string][]string{
	"ndb_provision_database": {
		"createdbserver": {
			"databasetype", "softwareprofileid", "softwareprofileversionid", "computeprofileid",
			"networkprofileid", "dbparameterprofileid", "nxclusterid", "sshpublickey", "timemachineinfo", "nodes",
		},
		"registerdbserver": {"databasetype", "dbparameterprofileid", "timemachineinfo", "nodes"},
	},
}

func SchemaValidation(resourceName string, d *schema.ResourceData) error {
	var diagMap []string
	if vals, ok := requiredResourceFields[resourceName]; ok {
		//nolint:staticcheck
		if dbVal, ok := d.GetOkExists("createdbserver"); ok {
			if dbVal.(bool) {
				createVals := vals["createdbserver"]
				for _, attr := range createVals {
					if _, ok := d.GetOk(attr); !ok {
						diagMap = append(diagMap, attr)
					}
				}
			} else {
				registerVals := vals["registerdbserver"]
				for _, attr := range registerVals {
					if _, ok := d.GetOk(attr); !ok {
						diagMap = append(diagMap, attr)
					}
				}
			}
		}
		if diagMap != nil {
			return fmt.Errorf("missing required fields are %s for %s", diagMap, resourceName)
		}
	}
	return nil
}
