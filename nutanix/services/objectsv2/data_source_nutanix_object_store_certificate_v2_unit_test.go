package objectstoresv2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitObjectStoreCertificateDataSourceSchema(t *testing.T) {
	t.Parallel()

	r := DatasourceNutanixObjectStoreCertificateV2()
	require.NoError(t, r.InternalValidate(r.Schema, false))
}
