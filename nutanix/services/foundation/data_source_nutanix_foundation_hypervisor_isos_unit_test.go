package foundation

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/foundation"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFoundationHypervisorIsosSchema(t *testing.T) {
	t.Parallel()

	r := DataSourceFoundationHypervisorIsos()
	require.NoError(t, r.InternalValidate(r.Schema, false))
	require.NoError(t, hypervisorSchema().Elem.(*schema.Resource).InternalValidate(
		hypervisorSchema().Elem.(*schema.Resource).Schema, true))
}

func TestUnitFlattenHypervisorIsos(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*foundation.HypervisorISOReference
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*foundation.HypervisorISOReference{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []*foundation.HypervisorISOReference{{
				Filename:  "kvm.iso",
				Supported: utils.BoolPtr(true),
			}},
			want: []map[string]interface{}{{
				"filename":  "kvm.iso",
				"supported": utils.BoolPtr(true),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenHypervisorIsos(tt.in))
		})
	}
}
