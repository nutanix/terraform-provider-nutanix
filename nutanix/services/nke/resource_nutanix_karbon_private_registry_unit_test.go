package nke

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitResourceNutanixKarbonPrivateRegistrySchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, ResourceNutanixKarbonPrivateRegistry().InternalValidate(ResourceNutanixKarbonPrivateRegistry().Schema, true))
}
