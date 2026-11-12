package networking

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceSchemas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "ResourceNutanixAddressGroup",
			fn: func() error {
				r := ResourceNutanixAddressGroup()
				return r.InternalValidate(r.Schema, true)
			},
		},
		{
			name: "ResourceNutanixFloatingIP",
			fn: func() error {
				r := ResourceNutanixFloatingIP()
				return r.InternalValidate(r.Schema, true)
			},
		},
		{
			name: "ResourceNutanixVPC",
			fn: func() error {
				r := ResourceNutanixVPC()
				return r.InternalValidate(r.Schema, true)
			},
		},
		{
			name: "ResourceNutanixStaticRoute",
			fn: func() error {
				r := ResourceNutanixStaticRoute()
				return r.InternalValidate(r.Schema, true)
			},
		},
		{
			name: "ResourceNutanixPbr",
			fn: func() error {
				r := ResourceNutanixPbr()
				return r.InternalValidate(r.Schema, true)
			},
		},
		{
			name: "ResourceNutanixSubnet",
			fn: func() error {
				r := ResourceNutanixSubnet()
				return r.InternalValidate(r.Schema, true)
			},
		},
		{
			name: "ResourceNutanixServiceGroup",
			fn: func() error {
				r := ResourceNutanixServiceGroup()
				return r.InternalValidate(r.Schema, true)
			},
		},
		{
			name: "ResourceNutanixNetworkSecurityRule",
			fn: func() error {
				r := ResourceNutanixNetworkSecurityRule()
				return r.InternalValidate(r.Schema, true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.NoError(t, tt.fn())
		})
	}
}

func TestUnitDataSourceSchemas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "DataSourceNutanixNetworkSecurityRule",
			fn: func() error {
				r := DataSourceNutanixNetworkSecurityRule()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixVPCs",
			fn: func() error {
				r := DataSourceNutanixVPCs()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixFloatingIP",
			fn: func() error {
				r := DataSourceNutanixFloatingIP()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixFloatingIPs",
			fn: func() error {
				r := DataSourceNutanixFloatingIPs()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixServiceGroups",
			fn: func() error {
				r := DataSourceNutanixServiceGroups()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixPbr",
			fn: func() error {
				r := DataSourceNutanixPbr()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixStaticRoute",
			fn: func() error {
				r := DataSourceNutanixStaticRoute()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixVPC",
			fn: func() error {
				r := DataSourceNutanixVPC()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixSubnet",
			fn: func() error {
				r := DataSourceNutanixSubnet()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixAddressGroup",
			fn: func() error {
				r := DataSourceNutanixAddressGroup()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixAddressGroups",
			fn: func() error {
				r := DataSourceNutanixAddressGroups()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixSubnets",
			fn: func() error {
				r := DataSourceNutanixSubnets()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixPbrs",
			fn: func() error {
				r := DataSourceNutanixPbrs()
				return r.InternalValidate(r.Schema, false)
			},
		},
		{
			name: "DataSourceNutanixServiceGroup",
			fn: func() error {
				r := DataSourceNutanixServiceGroup()
				return r.InternalValidate(r.Schema, false)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.NoError(t, tt.fn())
		})
	}
}

func TestUnitDataSourceFiltersSchema(t *testing.T) {
	t.Parallel()

	s := DataSourceFiltersSchema()
	require.NotNil(t, s)
	require.Equal(t, schema.TypeSet, s.Type)
}
