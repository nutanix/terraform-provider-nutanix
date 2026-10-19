package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitSchemas(t *testing.T) {
	t.Parallel()

	t.Run("DatasourceNutanixEffectiveImageRateLimitPoliciesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixEffectiveImageRateLimitPoliciesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixImagePlacementV4", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixImagePlacementV4()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixImagePlacementsV4", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixImagePlacementsV4()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixImageRateLimitPoliciesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixImageRateLimitPoliciesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixImageRateLimitPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixImageRateLimitPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixImageV4", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixImageV4()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixImagesV4", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixImagesV4()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixNGTConfigurationV4", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixNGTConfigurationV4()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixOvaV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixOvaV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixOvasV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixOvasV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixTemplatePlacementPoliciesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixTemplatePlacementPoliciesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixTemplatePlacementPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixTemplatePlacementPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixTemplateV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixTemplateV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixTemplatesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixTemplatesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVirtualMachineV4", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVirtualMachineV4()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVirtualMachinesV4", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVirtualMachinesV4()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVMAntiAffinityPoliciesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVMAntiAffinityPoliciesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVMAntiAffinityPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVMAntiAffinityPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmGuestCustomizationProfileV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmGuestCustomizationProfileV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmGuestCustomizationProfilesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmGuestCustomizationProfilesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVMHostAffinityPoliciesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVMHostAffinityPoliciesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVMHostAffinityPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVMHostAffinityPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmStartupPoliciesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmStartupPoliciesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmStartupPolicyDependencyConflictDependeeVmsV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmStartupPolicyDependencyConflictDependeeVmsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmStartupPolicyDependencyConflictDependentVmsV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmStartupPolicyDependencyConflictDependentVmsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmStartupPolicyStartConditionConflictDependeeVmsV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmStartupPolicyStartConditionConflictDependeeVmsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmStartupPolicyStartConditionConflictDependentVmsV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmStartupPolicyStartConditionConflictDependentVmsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmStartupPolicyDependencyConflictV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmStartupPolicyDependencyConflictV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmStartupPolicyDependencyConflictsV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmStartupPolicyDependencyConflictsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmStartupPolicyStartConditionConflictV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmStartupPolicyStartConditionConflictV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmStartupPolicyStartConditionConflictsV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmStartupPolicyStartConditionConflictsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmStartupPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmStartupPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVmStartupPolicyVmComplianceStatesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVmStartupPolicyVmComplianceStatesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("ResourceNutanixImagePlacementV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixImagePlacementV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixImageRateLimitPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixImageRateLimitPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixImageV4", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixImageV4()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixNGTInsertIsoV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixNGTInsertIsoV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixNGTInstallationV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixNGTInstallationV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixNGTUpgradeV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixNGTUpgradeV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixOvaDownloadV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixOvaDownloadV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixOvaV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixOvaV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixOvaVMDeploymentV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixOvaVMDeploymentV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixTemplateDeployV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixTemplateDeployV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixTemplateActionsV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixTemplateActionsV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixTemplatePlacementPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixTemplatePlacementPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixTemplatesV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixTemplatesV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVirtualMachineV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVirtualMachineV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVMAntiAffinityPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVMAntiAffinityPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVmsCdRomsInsertEjectV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVmsCdRomsInsertEjectV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVmGuestCustomizationProfileV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVmGuestCustomizationProfileV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVMHostAffinityPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVMHostAffinityPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixRevertVMRecoveryPointV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixRevertVMRecoveryPointV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVmStartupPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVmStartupPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVMCloneV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVMCloneV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVMGCUpdateV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVMGCUpdateV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVmsNetworkDeviceAssignIPV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVmsNetworkDeviceAssignIPV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVmsNetworkDeviceMigrateV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVmsNetworkDeviceMigrateV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVmsShutdownActionV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVmsShutdownActionV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

}
