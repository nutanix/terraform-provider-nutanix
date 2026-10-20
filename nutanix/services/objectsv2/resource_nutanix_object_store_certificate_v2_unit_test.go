package objectstoresv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	objectsCommon "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/common/v1/config"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitObjectStoreCertificateSchema(t *testing.T) {
	t.Parallel()

	r := ResourceNutanixObjectStoreCertificateV2()
	require.NoError(t, r.InternalValidate(r.Schema, true))
}

func TestUnitObjectStoreCertificateDelete(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixObjectStoreCertificateV2().Schema, map[string]interface{}{
		"object_store_ext_id": "os-1",
	})
	require.Nil(t, ResourceNutanixObjectStoreCertificateV2Delete(t.Context(), d, nil))
}

func TestUnitFlattenFQDNs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []objectsCommon.FQDN
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []objectsCommon.FQDN{}, want: nil},
		{
			name: "populated",
			in: []objectsCommon.FQDN{{
				Value: utils.StringPtr("objects.example.test"),
			}},
			want: []map[string]interface{}{{
				"value": "objects.example.test",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFQDNs(tt.in))
		})
	}
}
