package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/policies"
)

func TestUnitFlattenVMHostAffinityPolicyEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVMHostAffinityPolicyEntities(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVMHostAffinityPolicyEntities([]import1.VmHostAffinityPolicy{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVMHostAffinityPolicyEntities([]import1.VmHostAffinityPolicy{import1.VmHostAffinityPolicy{}})
		_ = got
	})

}
