package selfservice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitDataSourceNutanixSnapshotPolicySchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, DataSourceNutanixSnapshotPolicy().InternalValidate(DataSourceNutanixSnapshotPolicy().Schema, false))
}
