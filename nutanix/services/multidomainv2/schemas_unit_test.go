package multidomainv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUnitMultidomainSchemas(t *testing.T) {
	t.Parallel()

	t.Run("resources", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, ResourceNutanixProjectV2().InternalValidate(ResourceNutanixProjectV2().Schema, true))
		require.NoError(t, ResourceNutanixResourceGroupV2().InternalValidate(ResourceNutanixResourceGroupV2().Schema, true))
		require.NoError(t, ResourceNutanixLocalDomainV2().InternalValidate(ResourceNutanixLocalDomainV2().Schema, true))
		require.NoError(t, ResourceNutanixRegisteredDomainV2().InternalValidate(ResourceNutanixRegisteredDomainV2().Schema, true))
		require.NoError(t, ResourceNutanixRegisteredDomainRefreshApiCredentialsV2().InternalValidate(ResourceNutanixRegisteredDomainRefreshApiCredentialsV2().Schema, true))
	})

	t.Run("datasources", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, DatasourceNutanixProjectV2().InternalValidate(DatasourceNutanixProjectV2().Schema, false))
		require.NoError(t, DatasourceNutanixProjectsV2().InternalValidate(DatasourceNutanixProjectsV2().Schema, false))
		require.NoError(t, DatasourceNutanixResourceGroupV2().InternalValidate(DatasourceNutanixResourceGroupV2().Schema, false))
		require.NoError(t, DatasourceNutanixResourceGroupsV2().InternalValidate(DatasourceNutanixResourceGroupsV2().Schema, false))
		require.NoError(t, DatasourceNutanixLocalDomainV2().InternalValidate(DatasourceNutanixLocalDomainV2().Schema, false))
		require.NoError(t, DatasourceNutanixRegisteredDomainV2().InternalValidate(DatasourceNutanixRegisteredDomainV2().Schema, false))
		require.NoError(t, DatasourceNutanixRegisteredDomainsV2().InternalValidate(DatasourceNutanixRegisteredDomainsV2().Schema, false))
		require.NoError(t, DatasourceNutanixLocationsV2().InternalValidate(DatasourceNutanixLocationsV2().Schema, false))
	})

	t.Run("nested schemas", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, schemaRegisteredDomainPlatformData().InternalValidate(schemaRegisteredDomainPlatformData().Schema, true))
		require.NoError(t, schemaRegisteredDomainRegistrationConfig().InternalValidate(schemaRegisteredDomainRegistrationConfig().Schema, true))
		require.NoError(t, schemaRegisteredDomainTunnel().InternalValidate(schemaRegisteredDomainTunnel().Schema, true))
		require.NoError(t, schemaDatasourceRegisteredDomainPlatformData().InternalValidate(schemaDatasourceRegisteredDomainPlatformData().Schema, false))
		require.NoError(t, schemaDatasourceRegisteredDomainRegistrationConfig().InternalValidate(schemaDatasourceRegisteredDomainRegistrationConfig().Schema, false))
		require.NoError(t, schemaDatasourceRegisteredDomainTunnel().InternalValidate(schemaDatasourceRegisteredDomainTunnel().Schema, false))
		require.NoError(t, schemaResourceGroupPlacementTargets().InternalValidate(schemaResourceGroupPlacementTargets().Schema, true))
		require.NoError(t, schemaDatasourceResourceGroupPlacementTargets().InternalValidate(schemaDatasourceResourceGroupPlacementTargets().Schema, false))
		require.NoError(t, schemaLocalDomainConnectionStatus().InternalValidate(schemaLocalDomainConnectionStatus().Schema, false))
		require.NoError(t, schemaLocalDomainRegistrationStatus().InternalValidate(schemaLocalDomainRegistrationStatus().Schema, false))
		require.NoError(t, schemaLocalDomainLocation().InternalValidate(schemaLocalDomainLocation().Schema, false))
		require.NoError(t, schemaForLinks().Elem.(*schema.Resource).InternalValidate(schemaForLinks().Elem.(*schema.Resource).Schema, false))
		require.NoError(t, schemaForGeographicCoordinates().Elem.(*schema.Resource).InternalValidate(schemaForGeographicCoordinates().Elem.(*schema.Resource).Schema, false))
	})
}
