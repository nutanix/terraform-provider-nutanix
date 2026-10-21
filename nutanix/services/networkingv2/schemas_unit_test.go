package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitSchemas(t *testing.T) {
	t.Parallel()

	t.Run("DatasourceNutanixAddressGroupV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixAddressGroupV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixAddressGroupsV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixAddressGroupsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixFloatingIPV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixFloatingIPV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixFloatingIPsV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixFloatingIPsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DataSourceNutanixNetworkFunctionV2", func(t *testing.T) {
		t.Parallel()
		r := DataSourceNutanixNetworkFunctionV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DataSourceNutanixNetworkFunctionsV2", func(t *testing.T) {
		t.Parallel()
		r := DataSourceNutanixNetworkFunctionsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DataSourceNutanixNetworkSecurityPoliciesV2", func(t *testing.T) {
		t.Parallel()
		r := DataSourceNutanixNetworkSecurityPoliciesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DataSourceNutanixNetworkSecurityPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := DataSourceNutanixNetworkSecurityPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DataSourceNutanixNicProfileV2", func(t *testing.T) {
		t.Parallel()
		r := DataSourceNutanixNicProfileV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DataSourceNutanixNicProfilesV2", func(t *testing.T) {
		t.Parallel()
		r := DataSourceNutanixNicProfilesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixNodeSchedulableStatusesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixNodeSchedulableStatusesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixPbrV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixPbrV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixPbrsV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixPbrsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixRouteTableV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixRouteTableV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixRouteTablesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixRouteTablesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixRouteV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixRouteV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixRoutesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixRoutesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixServiceGroupV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixServiceGroupV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixServiceGroupsV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixServiceGroupsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DataSourceNutanixSubnetV2", func(t *testing.T) {
		t.Parallel()
		r := DataSourceNutanixSubnetV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DataSourceNutanixSubnetsV2", func(t *testing.T) {
		t.Parallel()
		r := DataSourceNutanixSubnetsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVirtualSwitchV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVirtualSwitchV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVirtualSwitchesV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVirtualSwitchesV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DataSourceNutanixVPCv2", func(t *testing.T) {
		t.Parallel()
		r := DataSourceNutanixVPCv2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DatasourceNutanixVpcVirtualSwitchMappingsV2", func(t *testing.T) {
		t.Parallel()
		r := DatasourceNutanixVpcVirtualSwitchMappingsV2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("DataSourceNutanixVPCsv2", func(t *testing.T) {
		t.Parallel()
		r := DataSourceNutanixVPCsv2()
		require.NoError(t, r.InternalValidate(r.Schema, false))
	})

	t.Run("ResourceNutanixAddressGroupsV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixAddressGroupsV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixFloatingIPv2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixFloatingIPv2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixNetworkFunctionV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixNetworkFunctionV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixNetworkSecurityPolicyV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixNetworkSecurityPolicyV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixNicProfileV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixNicProfileV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixPbrsV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixPbrsV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixRoutesV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixRoutesV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixServiceGroupsV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixServiceGroupsV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixSubnetV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixSubnetV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVirtualSwitchV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVirtualSwitchV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVpcVirtualSwitchMappingV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVpcVirtualSwitchMappingV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

	t.Run("ResourceNutanixVPCsV2", func(t *testing.T) {
		t.Parallel()
		r := ResourceNutanixVPCsV2()
		require.NoError(t, r.InternalValidate(r.Schema, true))
	})

}
