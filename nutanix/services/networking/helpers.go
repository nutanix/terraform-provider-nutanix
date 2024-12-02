package networking

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	// ERROR ..
	ERROR              = "ERROR"
	DEFAULTWAITTIMEOUT = 60
)

func getMetadataAttributes(d *schema.ResourceData, metadata *v3.Metadata, kind string) error {
	metadata.Kind = utils.StringPtr(kind)

	if v, ok := d.GetOk("categories"); ok {
		metadata.Categories = expandCategories(v)
	} else {
		metadata.Categories = nil
	}

	if p, ok := d.GetOk("project_reference"); ok {
		pr := p.(map[string]interface{})
		r := &v3.Reference{}
		if v1, ok1 := pr["name"]; ok1 {
			r.Name = utils.StringPtr(v1.(string))
		}
		if v2, ok2 := pr["kind"]; ok2 {
			r.Kind = utils.StringPtr(v2.(string))
		}
		if v3, ok3 := pr["uuid"]; ok3 {
			r.UUID = utils.StringPtr(v3.(string))
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

func flattenClusterReference(r *v3.Reference, d *schema.ResourceData) error {
	if r != nil {
		if err := d.Set("cluster_uuid", utils.StringValue(r.UUID)); err != nil {
			return err
		}

		if err := d.Set("cluster_name", utils.StringValue(r.Name)); err != nil {
			return err
		}
	}
	return nil
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

func expandReference(ref map[string]interface{}) *v3.Reference {
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
	if v, ok := ref["name"]; ok && v.(string) != "" {
		r.Name = utils.StringPtr(v.(string))
		hasValue = true
	}

	if hasValue {
		return r
	}

	return nil
}

func buildReference(uuid, kind string) *v3.Reference {
	return &v3.Reference{
		Kind: utils.StringPtr(kind),
		UUID: utils.StringPtr(uuid),
	}
}

func validateMapStringValue(value map[string]interface{}, key string) *string {
	if v, ok := value[key]; ok && v != nil && v.(string) != "" {
		return utils.StringPtr(v.(string))
	}
	return nil
}

func taskStateRefreshFunc(client *v3.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.V3.GetTask(taskUUID)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "INVALID_UUID") {
				return v, ERROR, nil
			}
			return nil, "", err
		}

		if *v.Status == "INVALID_UUID" || *v.Status == "FAILED" {
			return v, *v.Status,
				fmt.Errorf("error_detail: %s, progress_message: %s", utils.StringValue(v.ErrorDetail), utils.StringValue(v.ProgressMessage))
		}
		return v, *v.Status, nil
	}
}

func flattenArrayReferenceValues(refs []*v3.Reference) []map[string]interface{} {
	references := make([]map[string]interface{}, 0)
	for _, r := range refs {
		reference := make(map[string]interface{})
		if r != nil {
			reference["kind"] = utils.StringValue(r.Kind)
			reference["uuid"] = utils.StringValue(r.UUID)

			if r.Name != nil {
				reference["name"] = utils.StringValue(r.Name)
			}
			references = append(references, reference)
		}
	}

	return references
}

func flattenReferenceList(references []*v3.ReferenceValues) []map[string]interface{} {
	res := make([]map[string]interface{}, len(references))
	if len(references) > 0 {
		for i, r := range references {
			res[i] = flattenReference(r)
		}
	}
	return res
}

func flattenReference(reference *v3.ReferenceValues) map[string]interface{} {
	if reference != nil {
		return map[string]interface{}{
			"kind": reference.Kind,
			"uuid": reference.UUID,
			"name": reference.Name,
		}
	}
	return map[string]interface{}{}
}
