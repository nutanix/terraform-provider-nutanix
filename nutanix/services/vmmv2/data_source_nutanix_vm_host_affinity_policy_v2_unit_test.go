package vmmv2

import (
	"testing"

	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/policies"
)

func TestUnitFlattenVMHostAffinityPolicyEntity(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVMHostAffinityPolicyEntity(import1.VmHostAffinityPolicy{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVMHostAffinityPolicyEntity(import1.VmHostAffinityPolicy{})
		_ = got
	})

}
