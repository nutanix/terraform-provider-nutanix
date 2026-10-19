package foundation

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFoundationNodeNetworkDetailsSchema(t *testing.T) {
	t.Parallel()

	r := DataSourceNodeNetworkDetails()
	require.NoError(t, r.InternalValidate(r.Schema, false))
}

func TestUnitExpandStringList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want []*string
	}{
		{name: "nil", in: nil, want: []*string{}},
		{name: "empty", in: []interface{}{}, want: []*string{}},
		{name: "skips empty and non-string", in: []interface{}{"", 1, nil}, want: []*string{}},
		{
			name: "populated",
			in:   []interface{}{"fe80::1", "fe80::2"},
			want: []*string{utils.StringPtr("fe80::1"), utils.StringPtr("fe80::2")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandStringList(tt.in))
		})
	}
}
