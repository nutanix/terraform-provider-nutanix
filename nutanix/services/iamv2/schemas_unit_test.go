package iamv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUnitIamSchemas(t *testing.T) {
	t.Parallel()

	t.Run("resources", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, ResourceNutanixUserV2().InternalValidate(ResourceNutanixUserV2().Schema, true))
		require.NoError(t, ResourceNutanixUserGroupsV2().InternalValidate(ResourceNutanixUserGroupsV2().Schema, true))
		require.NoError(t, ResourceNutanixRolesV2().InternalValidate(ResourceNutanixRolesV2().Schema, true))
		require.NoError(t, ResourceNutanixRoleMembershipV2().InternalValidate(ResourceNutanixRoleMembershipV2().Schema, true))
		require.NoError(t, ResourceNutanixAuthPoliciesV2().InternalValidate(ResourceNutanixAuthPoliciesV2().Schema, true))
		require.NoError(t, ResourceNutanixDirectoryServicesV2().InternalValidate(ResourceNutanixDirectoryServicesV2().Schema, true))
		require.NoError(t, ResourceNutanixSamlIdpV2().InternalValidate(ResourceNutanixSamlIdpV2().Schema, true))
		require.NoError(t, ResourceNutanixUserKeyV2().InternalValidate(ResourceNutanixUserKeyV2().Schema, true))
		require.NoError(t, ResourceNutanixUserRevokeKeyV2().InternalValidate(ResourceNutanixUserRevokeKeyV2().Schema, true))
	})

	t.Run("datasources", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, DatasourceNutanixUserV2().InternalValidate(DatasourceNutanixUserV2().Schema, false))
		require.NoError(t, DatasourceNutanixUsersV2().InternalValidate(DatasourceNutanixUsersV2().Schema, false))
		require.NoError(t, DatasourceNutanixUserGroupV2().InternalValidate(DatasourceNutanixUserGroupV2().Schema, false))
		require.NoError(t, DatasourceNutanixUserGroupsV2().InternalValidate(DatasourceNutanixUserGroupsV2().Schema, false))
		require.NoError(t, DatasourceNutanixRoleV2().InternalValidate(DatasourceNutanixRoleV2().Schema, false))
		require.NoError(t, DatasourceNutanixRolesV2().InternalValidate(DatasourceNutanixRolesV2().Schema, false))
		require.NoError(t, DatasourceNutanixRoleMembershipV2().InternalValidate(DatasourceNutanixRoleMembershipV2().Schema, false))
		require.NoError(t, DatasourceNutanixRoleMembershipsV2().InternalValidate(DatasourceNutanixRoleMembershipsV2().Schema, false))
		require.NoError(t, DatasourceNutanixRoleMembershipSummaryV2().InternalValidate(DatasourceNutanixRoleMembershipSummaryV2().Schema, false))
		require.NoError(t, DatasourceNutanixAuthorizationPolicyV2().InternalValidate(DatasourceNutanixAuthorizationPolicyV2().Schema, false))
		require.NoError(t, DatasourceNutanixAuthorizationPoliciesV2().InternalValidate(DatasourceNutanixAuthorizationPoliciesV2().Schema, false))
		require.NoError(t, DatasourceNutanixDirectoryServiceV2().InternalValidate(DatasourceNutanixDirectoryServiceV2().Schema, false))
		require.NoError(t, DatasourceNutanixDirectoryServicesV2().InternalValidate(DatasourceNutanixDirectoryServicesV2().Schema, false))
		require.NoError(t, DatasourceNutanixDirectoryServiceUsersSearchV2().InternalValidate(DatasourceNutanixDirectoryServiceUsersSearchV2().Schema, false))
		require.NoError(t, DatasourceNutanixSamlIDPV2().InternalValidate(DatasourceNutanixSamlIDPV2().Schema, false))
		require.NoError(t, DatasourceNutanixSamlIDPsV2().InternalValidate(DatasourceNutanixSamlIDPsV2().Schema, false))
		require.NoError(t, DatasourceNutanixOperationV2().InternalValidate(DatasourceNutanixOperationV2().Schema, false))
		require.NoError(t, DatasourceNutanixOperationsV2().InternalValidate(DatasourceNutanixOperationsV2().Schema, false))
		require.NoError(t, DatasourceNutanixEntityV2().InternalValidate(DatasourceNutanixEntityV2().Schema, false))
		require.NoError(t, DatasourceNutanixEntitiesV2().InternalValidate(DatasourceNutanixEntitiesV2().Schema, false))
		require.NoError(t, DatasourceNutanixUserKeyV2().InternalValidate(DatasourceNutanixUserKeyV2().Schema, false))
		require.NoError(t, DatasourceNutanixUserKeysV2().InternalValidate(DatasourceNutanixUserKeysV2().Schema, false))
	})

	t.Run("nested schemas", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, schemaForLinks().Elem.(*schema.Resource).InternalValidate(schemaForLinks().Elem.(*schema.Resource).Schema, false))
		require.NoError(t, schemaForKeyValuePairs().Elem.(*schema.Resource).InternalValidate(schemaForKeyValuePairs().Elem.(*schema.Resource).Schema, true))
		require.NoError(t, schemaForScopeTemplateNameValues().Elem.(*schema.Resource).InternalValidate(schemaForScopeTemplateNameValues().Elem.(*schema.Resource).Schema, true))
		require.NoError(t, SchemaForLinks().Elem.(*schema.Resource).InternalValidate(SchemaForLinks().Elem.(*schema.Resource).Schema, false))
	})
}
