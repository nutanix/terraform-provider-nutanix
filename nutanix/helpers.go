package nutanix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func getMetadataAttributes(d *schema.ResourceData, metadata *v3.Metadata, kind string) error {
	metadata.Kind = utils.StringPtr(kind)

	if v, ok := d.GetOk("categories"); ok {
		catl := v.(map[string]interface{})
		metadata.Categories = expandCategories(catl)
	} else {
		metadata.Categories = nil
	}

	if p, ok := d.GetOk("project_reference"); ok {
		pr := p.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.StringPtr(pr["kind"].(string)),
			UUID: utils.StringPtr(pr["uuid"].(string)),
		}
		if v1, ok1 := pr["name"]; ok1 {
			r.Name = utils.StringPtr(v1.(string))
		}
		metadata.ProjectReference = r
	}
	if o, ok := d.GetOk("owner_reference"); ok {
		or := o.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.StringPtr(or["kind"].(string)),
			UUID: utils.StringPtr(or["uuid"].(string)),
		}
		if v1, ok1 := or["name"]; ok1 {
			r.Name = utils.StringPtr(v1.(string))
		}
		metadata.OwnerReference = r
	}

	return nil
}

func setRSEntityMetadata(v *v3.Metadata) (map[string]interface{}, map[string]interface{}) {
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = utils.TimeValue(v.LastUpdateTime).String()
	metadata["kind"] = utils.StringValue(v.Kind)
	metadata["uuid"] = utils.StringValue(v.UUID)
	metadata["creation_time"] = utils.TimeValue(v.CreationTime).String()
	metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(v.SpecVersion)))
	metadata["spec_hash"] = utils.StringValue(v.SpecHash)
	metadata["name"] = utils.StringValue(v.Name)

	c := make(map[string]interface{}, len(v.Categories))
	if v.Categories != nil {
		for name, values := range v.Categories {
			c[name] = values
		}
	}

	return metadata, c
}

func flattenReferenceValues(r *v3.Reference) map[string]interface{} {
	reference := make(map[string]interface{})
	if r != nil {
		reference["kind"] = utils.StringValue(r.Kind)
		reference["name"] = utils.StringValue(r.Name)
		reference["uuid"] = utils.StringValue(r.UUID)
	}

	return reference
}

func getClusterReferenceValues(r *v3.Reference) map[string]interface{} {
	reference := make(map[string]interface{})
	if r != nil {
		reference["kind"] = utils.StringValue(r.Kind)
		reference["uuid"] = utils.StringValue(r.UUID)
	}

	return reference
}

func validateRef(ref map[string]interface{}) *v3.Reference {
	r := &v3.Reference{}
	hasValue := false
	if v, ok := ref["kind"]; ok {
		r.Kind = utils.StringPtr(v.(string))
		hasValue = true
	}
	if v, ok := ref["uuid"]; ok {
		r.UUID = utils.StringPtr(v.(string))
		hasValue = true
	}
	if v, ok := ref["name"]; ok {
		r.Name = utils.StringPtr(v.(string))
		hasValue = true
	}

	if hasValue {
		return r
	}

	return nil
}

func validateShortRef(ref map[string]interface{}) *v3.Reference {
	r := &v3.Reference{}
	hasValue := false
	if v, ok := ref["kind"]; ok {
		r.Kind = utils.StringPtr(v.(string))
		hasValue = true
	}
	if v, ok := ref["uuid"]; ok {
		r.UUID = utils.StringPtr(v.(string))
		hasValue = true
	}

	if hasValue {
		return r
	}

	return nil
}

func validateMapStringValue(value map[string]interface{}, key string) *string {
	if v, ok := value[key]; ok && v != nil && v.(string) != "" {
		return utils.StringPtr(v.(string))
	}
	return nil
}

func validateMapIntValue(value map[string]interface{}, key string) *int64 {
	if v, ok := value[key]; ok && v != nil && v.(int) != 0 {
		return utils.Int64Ptr(int64(v.(int)))
	}
	return nil
}

func taskStateRefreshFunc(client *v3.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.V3.GetTask(taskUUID)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "INVALID_UUID") {
				return v, "ERROR", nil
			}
			return nil, "", err
		}

		if *v.Status == "INVALID_UUID" || *v.Status == "FAILED" {
			utils.PrintToJSON(v, "TASKS Validation")
			return v, *v.Status, fmt.Errorf(utils.StringValue(v.ErrorDetail))
		}
		return v, *v.Status, nil
	}
}
