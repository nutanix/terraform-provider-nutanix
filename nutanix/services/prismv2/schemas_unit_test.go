package prismv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUnitPrismV2Schemas(t *testing.T) {
	t.Parallel()

	t.Run("resources", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, ResourceNutanixBackupTargetV2().InternalValidate(ResourceNutanixBackupTargetV2().Schema, true))
		require.NoError(t, ResourceNutanixRestorePcV2().InternalValidate(ResourceNutanixRestorePcV2().Schema, true))
		require.NoError(t, ResourceNutanixRestoreSourceV2().InternalValidate(ResourceNutanixRestoreSourceV2().Schema, true))
		require.NoError(t, ResourceNutanixCategoriesV2().InternalValidate(ResourceNutanixCategoriesV2().Schema, true))
		require.NoError(t, ResourceNutanixDeployPcV2().InternalValidate(ResourceNutanixDeployPcV2().Schema, true))
		require.NoError(t, ResourceNutanixUnregisterClusterV2().InternalValidate(ResourceNutanixUnregisterClusterV2().Schema, true))
		require.NoError(t, ResourceNutanixClusterPCRegistrationV2().InternalValidate(ResourceNutanixClusterPCRegistrationV2().Schema, true))
	})

	t.Run("datasources", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, DatasourceNutanixFetchPcV2().InternalValidate(DatasourceNutanixFetchPcV2().Schema, false))
		require.NoError(t, DatasourceNutanixListRestorablePcsV2().InternalValidate(DatasourceNutanixListRestorablePcsV2().Schema, false))
		require.NoError(t, DatasourceNutanixListPcsV2().InternalValidate(DatasourceNutanixListPcsV2().Schema, false))
		require.NoError(t, DatasourceNutanixRestoreSourceV2().InternalValidate(DatasourceNutanixRestoreSourceV2().Schema, false))
		require.NoError(t, DatasourceNutanixBackupTargetsV2().InternalValidate(DatasourceNutanixBackupTargetsV2().Schema, false))
		require.NoError(t, DatasourceNutanixFetchRestorePointV2().InternalValidate(DatasourceNutanixFetchRestorePointV2().Schema, false))
		require.NoError(t, DatasourceNutanixFetchRestorePointsV2().InternalValidate(DatasourceNutanixFetchRestorePointsV2().Schema, false))
		require.NoError(t, DatasourceNutanixBackupTargetV2().InternalValidate(DatasourceNutanixBackupTargetV2().Schema, false))
		require.NoError(t, DatasourceNutanixCategoriesV2().InternalValidate(DatasourceNutanixCategoriesV2().Schema, false))
		require.NoError(t, DatasourceNutanixCategoryV2().InternalValidate(DatasourceNutanixCategoryV2().Schema, false))
	})

	t.Run("schema helpers", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, schemaForLinks().Elem.(*schema.Resource).InternalValidate(schemaForLinks().Elem.(*schema.Resource).Schema, true))
		require.NoError(t, schemaForPcConfig(true).InternalValidate(schemaForPcConfig(true).Schema, true))
		require.NoError(t, schemaForPcConfig(false).InternalValidate(schemaForPcConfig(false).Schema, false))
		require.NoError(t, schemaForPcNetwork().InternalValidate(schemaForPcNetwork().Schema, true))
		require.NoError(t, schemaForIPAddress().InternalValidate(schemaForIPAddress().Schema, true))
		require.NoError(t, schemaForIPAddressOrFqdn().InternalValidate(schemaForIPAddressOrFqdn().Schema, true))
		require.NoError(t, schemaForBackupTarget().InternalValidate(schemaForBackupTarget().Schema, true))
	})
}
