package networking

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitExpandCategories(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		set := schema.NewSet(categoriesSchema().Set, []interface{}{})
		require.Equal(t, map[string]string{}, expandCategories(set))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		set := schema.NewSet(categoriesSchema().Set, []interface{}{
			map[string]interface{}{"name": "Environment", "value": "Prod"},
		})
		require.Equal(t, map[string]string{"Environment": "Prod"}, expandCategories(set))
	})
}

func TestUnitFlattenCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[string]string
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{name: "empty", in: map[string]string{}, want: []interface{}{}},
		{
			name: "populated",
			in:   map[string]string{"Environment": "Prod"},
			want: []interface{}{map[string]interface{}{"name": "Environment", "value": "Prod"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenCategories(tt.in))
		})
	}
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
		{
			name: "populated",
			in:   []interface{}{"a", "", "b", 1},
			want: []*string{utils.StringPtr("a"), utils.StringPtr("b")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandStringList(tt.in))
		})
	}
}

func TestUnitBuildFiltersDataSource(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		set := schema.NewSet(schema.HashResource(DataSourceFiltersSchema().Elem.(*schema.Resource)), []interface{}{})
		require.Equal(t, []*client.AdditionalFilter{}, BuildFiltersDataSource(set))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		set := schema.NewSet(schema.HashResource(DataSourceFiltersSchema().Elem.(*schema.Resource)), []interface{}{
			map[string]interface{}{
				"name":   "name",
				"values": []interface{}{"vpc-1", "vpc-2"},
			},
		})
		got := BuildFiltersDataSource(set)
		require.Len(t, got, 1)
		require.Equal(t, "name", got[0].Name)
		require.Equal(t, []string{"vpc-1", "vpc-2"}, got[0].Values)
	})
}

func TestUnitReplaceFilterPrefixes(t *testing.T) {
	t.Parallel()

	filters := []*client.AdditionalFilter{
		{Name: "source.path", Values: []string{"v1"}},
		{Name: "keep.me", Values: []string{"v2"}},
	}

	t.Run("nil mappings", func(t *testing.T) {
		t.Parallel()
		in := []*client.AdditionalFilter{{Name: "source.path", Values: []string{"v1"}}}
		got := ReplaceFilterPrefixes(in, nil)
		require.Equal(t, "source.path", got[0].Name)
	})

	t.Run("empty mappings", func(t *testing.T) {
		t.Parallel()
		in := []*client.AdditionalFilter{{Name: "source.path", Values: []string{"v1"}}}
		got := ReplaceFilterPrefixes(in, map[string]string{})
		require.Equal(t, "source.path", got[0].Name)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		in := []*client.AdditionalFilter{
			{Name: filters[0].Name, Values: filters[0].Values},
			{Name: filters[1].Name, Values: filters[1].Values},
		}
		got := ReplaceFilterPrefixes(in, map[string]string{"source": "target"})
		require.Equal(t, "target.path", got[0].Name)
		require.Equal(t, "keep.me", got[1].Name)
	})
}

func TestUnitFlattenReferenceValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.Reference
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{
			name: "empty",
			in:   &v3.Reference{},
			want: map[string]interface{}{"kind": "", "uuid": ""},
		},
		{
			name: "populated",
			in: &v3.Reference{
				Kind: utils.StringPtr("subnet"),
				UUID: utils.StringPtr("uuid-1"),
				Name: utils.StringPtr("sub-1"),
			},
			want: map[string]interface{}{"kind": "subnet", "uuid": "uuid-1", "name": "sub-1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenReferenceValues(tt.in))
		})
	}
}

func TestUnitExpandReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[string]interface{}
		want *v3.Reference
	}{
		{name: "empty", in: map[string]interface{}{}, want: nil},
		{
			name: "populated",
			in: map[string]interface{}{
				"kind": "vpc",
				"uuid": "vpc-uuid",
				"name": "my-vpc",
			},
			want: &v3.Reference{
				Kind: utils.StringPtr("vpc"),
				UUID: utils.StringPtr("vpc-uuid"),
				Name: utils.StringPtr("my-vpc"),
			},
		},
		{
			name: "empty name omitted",
			in: map[string]interface{}{
				"kind": "vpc",
				"uuid": "vpc-uuid",
				"name": "",
			},
			want: &v3.Reference{
				Kind: utils.StringPtr("vpc"),
				UUID: utils.StringPtr("vpc-uuid"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandReference(tt.in))
		})
	}
}

func TestUnitFlattenArrayReferenceValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.Reference
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*v3.Reference{}, want: []map[string]interface{}{}},
		{
			name: "skips nil entries",
			in:   []*v3.Reference{nil},
			want: []map[string]interface{}{},
		},
		{
			name: "populated",
			in: []*v3.Reference{{
				Kind: utils.StringPtr("address_group"),
				UUID: utils.StringPtr("ag-1"),
				Name: utils.StringPtr("ag"),
			}},
			want: []map[string]interface{}{{
				"kind": "address_group",
				"uuid": "ag-1",
				"name": "ag",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenArrayReferenceValues(tt.in))
		})
	}
}

func TestUnitFlattenReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.ReferenceValues
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{name: "empty", in: &v3.ReferenceValues{}, want: map[string]interface{}{"kind": "", "uuid": "", "name": ""}},
		{
			name: "populated",
			in:   &v3.ReferenceValues{Kind: "policy", UUID: "p-1", Name: "pol"},
			want: map[string]interface{}{"kind": "policy", "uuid": "p-1", "name": "pol"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenReference(tt.in))
		})
	}
}

func TestUnitFlattenReferenceList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.ReferenceValues
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*v3.ReferenceValues{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in:   []*v3.ReferenceValues{{Kind: "policy", UUID: "p-1", Name: "pol"}},
			want: []map[string]interface{}{{"kind": "policy", "uuid": "p-1", "name": "pol"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenReferenceList(tt.in))
		})
	}
}

func TestUnitFlattenClusterReference(t *testing.T) {
	t.Parallel()

	schemaMap := map[string]*schema.Schema{
		"cluster_uuid": {Type: schema.TypeString, Optional: true},
		"cluster_name": {Type: schema.TypeString, Optional: true},
	}

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{})
		require.NoError(t, flattenClusterReference(nil, d))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{})
		require.NoError(t, flattenClusterReference(&v3.Reference{
			UUID: utils.StringPtr("c-uuid"),
			Name: utils.StringPtr("c-name"),
		}, d))
		require.Equal(t, "c-uuid", d.Get("cluster_uuid"))
		require.Equal(t, "c-name", d.Get("cluster_name"))
	})
}
