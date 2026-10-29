package foundation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitFoundationNOSPackagesSchema(t *testing.T) {
	t.Parallel()

	r := DataSourceFoundationNOSPackages()
	require.NoError(t, r.InternalValidate(r.Schema, false))
}
