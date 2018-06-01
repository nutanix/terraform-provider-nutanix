package nutanix

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func getMetadataAttributes(d *schema.ResourceData, metadata *v3.Metadata, kind string) error {
	metadata.Kind = utils.String(kind)

	if v, ok := d.GetOk("categories"); ok {
		catl := v.([]interface{})

		if len(catl) > 0 {
			cl := make(map[string]string)
			for _, v := range catl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["value"]; kok && k.(string) != "" {
						cl[i.(string)] = k.(string)
					}
				}
			}
			metadata.Categories = cl
		} else {
			metadata.Categories = nil
		}
	}
	if p, ok := d.GetOk("project_reference"); ok {
		pr := p.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(pr["kind"].(string)),
			UUID: utils.String(pr["uuid"].(string)),
		}
		if v1, ok1 := pr["name"]; ok1 {
			r.Name = utils.String(v1.(string))
		}
		metadata.ProjectReference = r
	}
	if o, ok := d.GetOk("owner_reference"); ok {
		or := o.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(or["kind"].(string)),
			UUID: utils.String(or["uuid"].(string)),
		}
		if v1, ok1 := or["name"]; ok1 {
			r.Name = utils.String(v1.(string))
		}
		metadata.OwnerReference = r
	}

	return nil
}

func readListMetadata(d *schema.ResourceData, kind string) (*v3.DSMetadata, error) {
	metadata := &v3.DSMetadata{
		Kind: utils.String(kind),
	}

	if v, ok := d.GetOk("metadata"); ok {
		m := v.(map[string]interface{})

		if mv, mok := m["sort_attribute"]; mok {
			metadata.SortAttribute = utils.String(mv.(string))
		}
		if mv, mok := m["filter"]; mok {
			metadata.Filter = utils.String(mv.(string))
		}
		if mv, mok := m["length"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return nil, err
			}
			metadata.Length = utils.Int64(int64(i))
		}
		if mv, mok := m["sort_order"]; mok {
			metadata.SortOrder = utils.String(mv.(string))
		}
		if mv, mok := m["offset"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return nil, err
			}
			metadata.Offset = utils.Int64(int64(i))
		}
	}

	return metadata, nil
}

func setRSEntityMetadata(v *v3.Metadata) (map[string]interface{}, []map[string]interface{}) {
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = utils.TimeValue(v.LastUpdateTime).String()
	metadata["kind"] = utils.StringValue(v.Kind)
	metadata["uuid"] = utils.StringValue(v.UUID)
	metadata["creation_time"] = utils.TimeValue(v.CreationTime).String()
	metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(v.SpecVersion)))
	metadata["spec_hash"] = utils.StringValue(v.SpecHash)
	metadata["name"] = utils.StringValue(v.Name)

	c := make([]map[string]interface{}, 0)
	if v.Categories != nil {
		categories := v.Categories
		var catList []map[string]interface{}

		for name, values := range categories {
			catItem := make(map[string]interface{})
			catItem["name"] = name
			catItem["value"] = values
			catList = append(catList, catItem)
		}
		c = catList
	}

	return metadata, c
}

func getReferenceValues(r *v3.Reference) map[string]interface{} {
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
		r.Kind = utils.String(v.(string))
		hasValue = true
	}
	if v, ok := ref["uuid"]; ok {
		r.UUID = utils.String(v.(string))
		hasValue = true
	}
	// if v, ok := ref["name"]; ok {
	// 	r.Name = utils.String(v.(string))
	// 	hasValue = true
	// }

	if hasValue {
		return r
	}

	return nil
}

func validateMapStringValue(value map[string]interface{}, key string) *string {
	if v, ok := value[key]; ok && v != nil && v.(string) != "" {
		return utils.String(v.(string))
	}
	return nil
}

func validateMapIntValue(value map[string]interface{}, key string) *int64 {
	if v, ok := value[key]; ok && v != nil && v.(int) != 0 {
		return utils.Int64(int64(v.(int)))
	}
	return nil
}
