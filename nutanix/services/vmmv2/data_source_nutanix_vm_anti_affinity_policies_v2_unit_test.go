package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/policies"
)

func TestUnitFlattenVMAntiAffinityPolicyEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVMAntiAffinityPolicyEntities(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVMAntiAffinityPolicyEntities([]import1.VmAntiAffinityPolicy{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVMAntiAffinityPolicyEntities([]import1.VmAntiAffinityPolicy{import1.VmAntiAffinityPolicy{}})
		_ = got
	})

}
