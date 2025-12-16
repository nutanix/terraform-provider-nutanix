package prism

import (
	"context"
	"fmt"
	"log"
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

// customizeDiffProjectACP handles the custom diff logic for the "acp" (Access Control Policy) attribute.
// Problem: Terraform treats lists positionally, so if the API returns ACPs in a different order
// than the user specified, Terraform detects a "change" even though the content is identical.
//
// Solution: This function normalizes the ACP list by:
//  1. Matching ACPs by role_reference.uuid (the unique identifier for an ACP)
//  2. Merging computed fields (name, metadata, context_filter_list) from the OLD state
//     with user-specified fields (role_reference, user_reference_list, user_group_reference_list) from NEW config
//  3. Maintaining the old order for existing ACPs and appending new ACPs at the end
//
// Note: user_reference_list and user_group_reference_list within each ACP now use TypeSet
// with UUID-based hashing, which handles order-independent comparison automatically.
//
// This ensures Terraform only shows actual content changes, not order-based false positives.
func customizeDiffProjectACP(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	log.Printf("[DEBUG] CustomizeDiff resource_nutanix_project")

	if !diff.HasChange("acp") {
		log.Printf("[DEBUG] CustomizeDiff: no acp change detected")
		return nil
	}

	oldACP, newACP := diff.GetChange("acp")
	oldACPList, ok1 := oldACP.([]interface{})
	newACPList, ok2 := newACP.([]interface{})

	if !ok1 || !ok2 {
		log.Printf("[DEBUG] CustomizeDiff: failed to convert ACP lists (ok1=%v, ok2=%v)", ok1, ok2)
		return nil
	}

	log.Printf("[DEBUG] CustomizeDiff: oldACPList length: %d, newACPList length: %d", len(oldACPList), len(newACPList))

	if len(oldACPList) == 0 || len(newACPList) == 0 {
		log.Printf("[DEBUG] CustomizeDiff: one of the lists is empty, skipping merge")
		return nil
	}

	// Build a map from role UUID to old ACP item (for getting computed fields)
	oldRoleToItem := make(map[string]map[string]interface{})
	for i, oldItem := range oldACPList {
		oldMap, ok := oldItem.(map[string]interface{})
		if !ok {
			continue
		}
		roleUUID := getACPRoleUUID(oldItem)
		if roleUUID != "" {
			oldRoleToItem[roleUUID] = oldMap
			log.Printf("[DEBUG] CustomizeDiff: oldACP[%d] role UUID: %s", i, roleUUID)
		}
	}

	// Build a map from role UUID to new ACP item (for getting user-specified fields)
	newRoleToItem := make(map[string]map[string]interface{})
	var newRoleOrder []string
	for i, newItem := range newACPList {
		newMap, ok := newItem.(map[string]interface{})
		if !ok {
			continue
		}
		roleUUID := getACPRoleUUID(newItem)
		if roleUUID != "" {
			newRoleToItem[roleUUID] = newMap
			newRoleOrder = append(newRoleOrder, roleUUID)
			log.Printf("[DEBUG] CustomizeDiff: newACP[%d] role UUID: %s", i, roleUUID)
		}
	}

	// Create merged ACPs in old order
	// For each ACP: take computed fields from old, user-specified fields from new
	mergedACPs := make([]interface{}, 0, len(newACPList))
	usedRoles := make(map[string]bool)

	// First, process ACPs that exist in both old and new (in old order)
	for _, oldItem := range oldACPList {
		oldMap, ok := oldItem.(map[string]interface{})
		if !ok {
			continue
		}
		oldRoleUUID := getACPRoleUUID(oldItem)
		newMap, exists := newRoleToItem[oldRoleUUID]
		if !exists {
			// ACP was removed in new config - skip it
			log.Printf("[DEBUG] CustomizeDiff: ACP with role %s removed in new config", oldRoleUUID)
			continue
		}

		// Merge: start with old ACP (has computed fields), overlay new user-specified fields
		merged := make(map[string]interface{})

		// Copy all computed fields from old ACP (name, metadata, context_filter_list)
		for k, v := range oldMap {
			merged[k] = v
		}

		// Overlay user-specified fields from new ACP
		if roleRef, ok := newMap["role_reference"]; ok {
			merged["role_reference"] = roleRef
		}

		// user_reference_list and user_group_reference_list use TypeSet with hash based on UUID
		// TypeSet handles order-independent comparison automatically, so we just pass through the new values
		if newUserRefList, ok := newMap["user_reference_list"]; ok {
			merged["user_reference_list"] = newUserRefList
			log.Printf("[DEBUG] CustomizeDiff: passing through user_reference_list for role %s (TypeSet handles ordering)", oldRoleUUID)
		}

		if newUserGroupRefList, ok := newMap["user_group_reference_list"]; ok {
			merged["user_group_reference_list"] = newUserGroupRefList
			log.Printf("[DEBUG] CustomizeDiff: passing through user_group_reference_list for role %s (TypeSet handles ordering)", oldRoleUUID)
		}

		// Preserve description if specified
		if desc, ok := newMap["description"]; ok && desc != "" {
			merged["description"] = desc
		}

		mergedACPs = append(mergedACPs, merged)
		usedRoles[oldRoleUUID] = true
		log.Printf("[DEBUG] CustomizeDiff: merged ACP with role %s in position %d", oldRoleUUID, len(mergedACPs)-1)
	}

	// Add new ACPs that weren't in the old list (these are truly new)
	for _, roleUUID := range newRoleOrder {
		if !usedRoles[roleUUID] {
			mergedACPs = append(mergedACPs, newRoleToItem[roleUUID])
			log.Printf("[DEBUG] CustomizeDiff: adding new ACP with role %s at position %d", roleUUID, len(mergedACPs)-1)
		}
	}

	// Set the merged list to update the diff
	log.Printf("[DEBUG] CustomizeDiff: setting merged ACPs (count: %d)", len(mergedACPs))
	if err := diff.SetNew("acp", mergedACPs); err != nil {
		log.Printf("[DEBUG] CustomizeDiff: failed to SetNew for acp: %v", err)
		return err
	}

	return nil
}

// getACPRoleUUID extracts the role UUID from an ACP item.
// It navigates the nested structure: acp -> role_reference[0] -> uuid
// Returns empty string if the structure is invalid or UUID is not found.
func getACPRoleUUID(acpItem interface{}) string {
	acpMap, ok := acpItem.(map[string]interface{})
	if !ok {
		return ""
	}
	if roleRef, ok := acpMap["role_reference"].([]interface{}); ok && len(roleRef) > 0 {
		if roleMap, ok := roleRef[0].(map[string]interface{}); ok {
			if uuid, ok := roleMap["uuid"].(string); ok {
				return uuid
			}
		}
	}
	return ""
}

// acpReferenceHash generates a hash for user_reference_list and user_group_reference_list items
// based on UUID. This ensures order-independent comparison in TypeSet.
func acpReferenceHash(v interface{}) int {
	m, ok := v.(map[string]interface{})
	if !ok {
		return 0
	}
	if uuid, ok := m["uuid"].(string); ok {
		return schema.HashString(uuid)
	}
	return 0
}
