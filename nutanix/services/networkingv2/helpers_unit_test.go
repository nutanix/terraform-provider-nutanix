package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitIndexOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		slice  []string
		target string
		want   int
	}{
		{name: "nil", slice: nil, target: "a", want: -1},
		{name: "empty", slice: []string{}, target: "a", want: -1},
		{name: "found", slice: []string{"a", "b", "c"}, target: "b", want: 1},
		{name: "missing", slice: []string{"a", "b"}, target: "z", want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, indexOf(tt.slice, tt.target))
		})
	}
}

func TestUnitDiffProjectLists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		old  []string
		new  []string
		add  []string
		rem  []string
	}{
		{name: "nil", old: nil, new: nil, add: nil, rem: nil},
		{name: "empty", old: []string{}, new: []string{}, add: nil, rem: nil},
		{name: "populated", old: []string{"a", "b"}, new: []string{"b", "c"}, add: []string{"c"}, rem: []string{"a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			added, removed := diffProjectLists(tt.old, tt.new)
			require.Equal(t, tt.add, added)
			require.Equal(t, tt.rem, removed)
		})
	}
}
