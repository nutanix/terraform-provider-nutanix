package prism

import (
	"fmt"
	"reflect"
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
	ERROR = "ERROR"
)

var (
	subnetDelay      = 10 * time.Second
	subnetMinTimeout = 3 * time.Second
	vmDelay          = 3 * time.Second
	vmMinTimeout     = 3 * time.Second
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

func flattenExternalNetworkListReferenceList(references []*v3.ReferenceValues) []map[string]interface{} {
	res := make([]map[string]interface{}, len(references))
	if len(references) > 0 {
		for i, r := range references {
			res[i] = flattenExternalNetworkListReference(r)
		}
	}
	return res
}

func flattenExternalNetworkListReference(reference *v3.ReferenceValues) map[string]interface{} {
	if reference != nil {
		return map[string]interface{}{
			"uuid": reference.UUID,
			"name": reference.Name,
		}
	}
	return map[string]interface{}{}
}

func expandDirectoryUserGroup(pr []interface{}) *v3.DirectoryServiceUserGroup {
	if len(pr) > 0 {
		res := &v3.DirectoryServiceUserGroup{}
		ent := pr[0].(map[string]interface{})

		if pnum, pk := ent["distinguished_name"]; pk && len(pnum.(string)) > 0 {
			res.DistinguishedName = utils.StringPtr(pnum.(string))
		}
		return res
	}
	return nil
}

func expandSamlUserGroup(pr []interface{}) *v3.SamlUserGroup {
	if len(pr) > 0 {
		res := &v3.SamlUserGroup{}
		ent := pr[0].(map[string]interface{})

		if idp, iok := ent["idp_uuid"]; iok {
			res.IdpUUID = utils.StringPtr(idp.(string))
		}

		if name, nok := ent["name"]; nok {
			res.Name = utils.StringPtr(name.(string))
		}

		return res
	}
	return nil
}

func expandIdentityProviderUser(d *schema.ResourceData) *v3.IdentityProvider {
	identityProviderState, ok := d.GetOk("identity_provider_user")
	if !ok {
		return nil
	}

	identityProviderMap := identityProviderState.([]interface{})[0].(map[string]interface{})
	identityProvider := &v3.IdentityProvider{}

	if username, ok := identityProviderMap["username"]; ok {
		identityProvider.Username = utils.StringPtr(username.(string))
	}

	if ipr, ok := identityProviderMap["identity_provider_reference"]; ok {
		identityProvider.IdentityProviderReference = expandReference(ipr.([]interface{})[0].(map[string]interface{}))
	}

	if !reflect.DeepEqual(*identityProvider, v3.IdentityProvider{}) {
		return identityProvider
	}
	return nil
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
