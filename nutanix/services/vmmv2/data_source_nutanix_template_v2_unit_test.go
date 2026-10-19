package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/common/v1/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/iam/v4/authn"
	import6 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/config"
	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/content"
)

func TestUnitFlattenGuestUpdateStatus(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenGuestUpdateStatus(&import5.GuestUpdateStatus{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenGuestUpdateStatus(&import5.GuestUpdateStatus{})
		_ = got
	})

}

func TestUnitFlattenTemplateUser(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenTemplateUser(&import5.TemplateUser{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenTemplateUser(&import5.TemplateUser{})
		_ = got
	})

}

func TestUnitFlattenUserType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenUserType(func() *authn.UserType { var v authn.UserType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenUserType(func() *authn.UserType { v := authn.UserType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenCustomKVPair(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenCustomKVPair(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCustomKVPair([]config.KVPair{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCustomKVPair([]config.KVPair{config.KVPair{}})
		_ = got
	})

}

func TestUnitFlattenKVValue(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenKVValue([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenKVValue([]interface{}{map[string]interface{}{"ext_id": "ext-1"}})
		_ = got
	})

}

func TestUnitFlattenUserStatusType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenUserStatusType(func() *authn.UserStatusType { var v authn.UserStatusType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenUserStatusType(func() *authn.UserStatusType { v := authn.UserStatusType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenBucketsAccessKey(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenBucketsAccessKey(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenBucketsAccessKey([]authn.BucketsAccessKey{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenBucketsAccessKey([]authn.BucketsAccessKey{authn.BucketsAccessKey{}})
		_ = got
	})

}

func TestUnitFlattenTemplateVersionSpec(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenTemplateVersionSpec(&import5.TemplateVersionSpec{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenTemplateVersionSpec(&import5.TemplateVersionSpec{})
		_ = got
	})

}

func TestUnitFlattenTemplateVersionSource(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenTemplateVersionSource(nil))
	})
}

func TestUnitFlattenTemplateVMRefOverrideVMConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenTemplateVMRefOverrideVMConfig(&import5.VmConfigOverride{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenTemplateVMRefOverrideVMConfig(&import5.VmConfigOverride{})
		_ = got
	})

}

func TestUnitExpandVmGcProfileConfigOverride(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandVmGcProfileConfigOverride(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVmGcProfileConfigOverride([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandConfigOverrideSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandConfigOverrideSpec(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandConfigOverrideSpec([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandSysprepParamsOverrideSpec(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandSysprepParamsOverrideSpec(map[string]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandSysprepParamsOverrideSpec(map[string]interface{}{"name": "UNKNOWN", "value": "UNKNOWN", "prefix_length": 24})
		_ = got
	})

}

func TestUnitFlattenVmGcProfileConfigOverride(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVmGcProfileConfigOverride(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVmGcProfileConfigOverride(&import6.VmGcProfileConfig{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVmGcProfileConfigOverride(&import6.VmGcProfileConfig{})
		_ = got
	})

}

func TestUnitFlattenConfigOverrideSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenConfigOverrideSpec(nil))
	})

}

func TestUnitFlattenSysprepConfigOverrideSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenSysprepConfigOverrideSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSysprepConfigOverrideSpec(&import6.VmGcProfileSysprepConfigOverrideSpec{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSysprepConfigOverrideSpec(&import6.VmGcProfileSysprepConfigOverrideSpec{})
		_ = got
	})

}

func TestUnitFlattenSysprepParamsOverrideSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenSysprepParamsOverrideSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSysprepParamsOverrideSpec(&import6.VmGcProfileSysprepParamsOverrideSpec{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSysprepParamsOverrideSpec(&import6.VmGcProfileSysprepParamsOverrideSpec{})
		_ = got
	})

}

func TestUnitFlattenGeneralSettingsOverrideSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenGeneralSettingsOverrideSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenGeneralSettingsOverrideSpec(&import6.VmGcProfileGeneralSettingsOverrideSpec{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenGeneralSettingsOverrideSpec(&import6.VmGcProfileGeneralSettingsOverrideSpec{})
		_ = got
	})

}

func TestUnitFlattenLocaleSettingsOverrideSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenLocaleSettingsOverrideSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenLocaleSettingsOverrideSpec(&import6.VmGcProfileLocaleSettingsOverrideSpec{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenLocaleSettingsOverrideSpec(&import6.VmGcProfileLocaleSettingsOverrideSpec{})
		_ = got
	})

}

func TestUnitFlattenNetworkSettingsOverrideSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenNetworkSettingsOverrideSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkSettingsOverrideSpec(&import6.VmGcProfileNetworkSettingsOverrideSpec{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkSettingsOverrideSpec(&import6.VmGcProfileNetworkSettingsOverrideSpec{})
		_ = got
	})

}

func TestUnitExpandVmGcProfileReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandVmGcProfileReference(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVmGcProfileReference([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenVmGcProfileReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVmGcProfileReference(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVmGcProfileReference(&import6.VmGcProfileReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVmGcProfileReference(&import6.VmGcProfileReference{})
		_ = got
	})

}

func TestUnitFlattenVM(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVM(&import6.Vm{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVM(&import6.Vm{})
		_ = got
	})

}
