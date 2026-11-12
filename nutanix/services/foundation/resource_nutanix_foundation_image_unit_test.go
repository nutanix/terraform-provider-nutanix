package foundation

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUnitFoundationImageSchema(t *testing.T) {
	t.Parallel()

	r := ResourceNutanixFoundationImage()
	require.NoError(t, r.InternalValidate(r.Schema, true))
}

func TestUnitFoundationImageRead(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixFoundationImage().Schema, map[string]interface{}{
		"filename":       "nos.tar",
		"installer_type": "nos",
	})
	require.Nil(t, resourceNutanixFoundationImageRead(t.Context(), d, nil))
}
