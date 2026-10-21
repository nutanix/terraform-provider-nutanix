package nke

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitDataSourceNutanixKarbonClusterSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, DataSourceNutanixKarbonCluster().InternalValidate(DataSourceNutanixKarbonCluster().Schema, false))
}

func TestUnitDataSourceNutanixKarbonClustersSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, DataSourceNutanixKarbonClusters().InternalValidate(DataSourceNutanixKarbonClusters().Schema, false))
}

func TestUnitDataSourceNutanixKarbonClusterKubeconfigSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, DataSourceNutanixKarbonClusterKubeconfig().InternalValidate(DataSourceNutanixKarbonClusterKubeconfig().Schema, false))
}

func TestUnitDataSourceNutanixKarbonClusterSSHSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, DataSourceNutanixKarbonClusterSSH().InternalValidate(DataSourceNutanixKarbonClusterSSH().Schema, false))
}

func TestUnitDataSourceNutanixKarbonPrivateRegistrySchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, DataSourceNutanixKarbonPrivateRegistry().InternalValidate(DataSourceNutanixKarbonPrivateRegistry().Schema, false))
}

func TestUnitDataSourceNutanixKarbonPrivateRegistriesSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, DataSourceNutanixKarbonPrivateRegistries().InternalValidate(DataSourceNutanixKarbonPrivateRegistries().Schema, false))
}
