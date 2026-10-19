package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/policies"
)

func TestUnitFlattenVmReferences(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVmReferences(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVmReferences([]import1.VmReference{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVmReferences([]import1.VmReference{import1.VmReference{}})
		_ = got
	})

}
