package objectstoresv2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitObjectStoreDataSourceSchema(t *testing.T) {
	t.Parallel()

	r := DatasourceNutanixObjectStoreV2()
	require.NoError(t, r.InternalValidate(r.Schema, false))
}
