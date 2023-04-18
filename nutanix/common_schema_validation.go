package nutanix

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var requiredResourceFields map[string][]string = map[string][]string{
	"era_provision_database": {"databasetype", "dbparameterprofileid", "timemachineinfo", "nodes"},
}

func schemaValidation(resourceName string, d *schema.ResourceData) error {
	var diagMap []string
	if vals, ok := requiredResourceFields[resourceName]; ok {
		for _, attr := range vals {
			if _, ok := d.GetOk(attr); !ok {
				diagMap = append(diagMap, attr)
			}
		}

		if diagMap != nil {
			return fmt.Errorf("missing required fields are %s for %s", diagMap, resourceName)
		}
	}
	return nil
}
