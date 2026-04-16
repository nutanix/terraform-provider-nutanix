package iam

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

var (
	subnetDelay      = 10 * time.Second
	subnetMinTimeout = 3 * time.Second
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

func getMetadataAttributesV2(d *schema.ResourceData, metadata *v3.Metadata, kind string) error {
	metadata.Kind = utils.StringPtr(kind)

	if v, ok := d.GetOk("categories"); ok {
		metadata.Categories = expandCategories(v)
	} else {
		metadata.Categories = nil
	}

	if p, ok := d.GetOk("project_reference"); ok {
		metadata.ProjectReference = validateRefList(p.([]interface{}), utils.StringPtr("project"))
	}
	if o, ok := d.GetOk("owner_reference"); ok {
		metadata.OwnerReference = validateRefList(o.([]interface{}), nil)
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

func validateArrayRef(references interface{}, kindValue *string) []*v3.Reference {
	refs := make([]*v3.Reference, 0)

	for _, s := range references.(*schema.Set).List() {
		ref := s.(map[string]interface{})
		r := v3.Reference{}

		if v, ok := ref["kind"]; ok {
			kind := v.(string)
			if kindValue != nil {
				kind = *kindValue
			}
			r.Kind = utils.StringPtr(kind)
		}

		if v, ok := ref["uuid"]; ok {
			r.UUID = utils.StringPtr(v.(string))
		}
		if v, ok := ref["name"]; ok {
			r.Name = utils.StringPtr(v.(string))
		}

		refs = append(refs, &r)
	}
	if len(refs) > 0 {
		return refs
	}

	return nil
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

func validateRefList(refs []interface{}, kindValue *string) *v3.Reference {
	r := &v3.Reference{}
	hasValue := false

	for _, v2 := range refs {
		ref := v2.(map[string]interface{})

		if v, ok := ref["kind"]; ok {
			r.Kind = utils.StringPtr(v.(string))
			hasValue = true
		}
		if kindValue != nil {
			r.Kind = kindValue
		}
		if v, ok := ref["uuid"]; ok {
			r.UUID = utils.StringPtr(v.(string))
			hasValue = true
		}
		if v, ok := ref["name"]; ok {
			r.Name = utils.StringPtr(v.(string))
			hasValue = true
		}
	}

	if hasValue {
		return r
	}

	return nil
}

func flattenReferenceValuesList(r *v3.Reference) []interface{} {
	references := make([]interface{}, 0)
	if r != nil {
		reference := make(map[string]interface{})
		reference["kind"] = utils.StringValue(r.Kind)
		reference["uuid"] = utils.StringValue(r.UUID)

		if r.Name != nil {
			reference["name"] = utils.StringValue(r.Name)
		}

		references = append(references, reference)
	}
	return references
}

func buildDataSourceListMetadata(set *schema.Set) *v3.DSMetadata {
	filters := v3.DSMetadata{}
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		if m["filter"].(string) != "" {
			filters.Filter = utils.StringPtr(m["filter"].(string))
		}
		if m["kind"].(string) != "" {
			filters.Kind = utils.StringPtr(m["kind"].(string))
		}
		if m["sort_order"].(string) != "" {
			filters.SortOrder = utils.StringPtr(m["sort_order"].(string))
		}
		if m["offset"].(int) != 0 {
			filters.Offset = utils.Int64Ptr(int64(m["offset"].(int)))
		}
		if m["length"].(int) != 0 {
			filters.Length = utils.Int64Ptr(int64(m["length"].(int)))
		}
		if m["sort_attribute"].(string) != "" {
			filters.SortAttribute = utils.StringPtr(m["sort_attribute"].(string))
		}
	}
	return &filters
}

func flattenContextList(contextList []*v3.ContextList) []interface{} {
	contexts := make([]interface{}, 0)
	for _, con := range contextList {
		if con != nil {
			scope := make(map[string]interface{})
			scope["scope_filter_expression_list"] = flattenScopeExpressionList(con.ScopeFilterExpressionList)
			scope["entity_filter_expression_list"] = flattenEntityExpressionList(con.EntityFilterExpressionList)

			contexts = append(contexts, scope)
		}
	}

	return contexts
}

func flattenScopeExpressionList(scopeList []*v3.ScopeFilterExpressionList) []interface{} {
	scopes := make([]interface{}, 0)

	for _, sco := range scopeList {
		scope := make(map[string]interface{})
		scope["left_hand_side"] = sco.LeftHandSide
		scope["operator"] = sco.Operator
		scope["right_hand_side"] = flattenRightHandSide(sco.RightHandSide)

		scopes = append(scopes, scope)
	}

	return scopes
}

func flattenEntityExpressionList(entities []v3.EntityFilterExpressionList) []interface{} {
	scopes := make([]interface{}, 0)

	for _, ent := range entities {
		scope := make(map[string]interface{})
		scope["left_hand_side_entity_type"] = utils.StringValue(ent.LeftHandSide.EntityType)
		scope["operator"] = ent.Operator
		scope["right_hand_side"] = flattenRightHandSide(ent.RightHandSide)

		scopes = append(scopes, scope)
	}

	return scopes
}

func flattenRightHandSide(right v3.RightHandSide) []interface{} {
	rightHand := make([]interface{}, 0)

	r := make(map[string]interface{})
	r["collection"] = utils.StringValue(right.Collection)
	r["uuid_list"] = right.UUIDList
	r["categories"] = flattenTightHandsideCategories(right.Categories)

	rightHand = append(rightHand, r)

	return rightHand
}

func flattenTightHandsideCategories(categories map[string][]string) []interface{} {
	c := make([]interface{}, 0)

	for name, value := range categories {
		c = append(c, map[string]interface{}{
			"name":  name,
			"value": value,
		})
	}

	return c
}
