package clusters

import (
	"strconv"

	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	// ERROR ..
	ERROR = "ERROR"
)

func setRSEntityMetadata(v *v3.Metadata) (map[string]interface{}, []interface{}) {
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = utils.TimeValue(v.LastUpdateTime).String()
	metadata["uuid"] = utils.StringValue(v.UUID)
	metadata["creation_time"] = utils.TimeValue(v.CreationTime).String()
	metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(v.SpecVersion)))
	metadata["spec_hash"] = utils.StringValue(v.SpecHash)
	metadata["name"] = utils.StringValue(v.Name)

	return metadata, flattenCategories(v.Categories)
}

func flattenReferenceValues(r *v3.Reference) map[string]interface{} {
	reference := make(map[string]interface{})
	if r != nil {
		reference["kind"] = utils.StringValue(r.Kind)
		reference["uuid"] = utils.StringValue(r.UUID)
		if r.Name != nil {
			reference["name"] = utils.StringValue(r.Name)
		}
	}
	return reference
}
