package volumesv2

import (
	"testing"

	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
)

// TestExpandIscsiFeatures_OmittedOptionalFields is the regression test for a panic on
// a perfectly valid configuration.
//
// Elements of a nested TypeList always carry every schema key, so an omitted optional
// attribute arrives as "" rather than being absent. The old code did
// enabledAuthenticationsMap[""] (nil) followed by nil.(int), which panicked during
// plan/apply for:
//
//	iscsi_features { target_secret = "..." }
func TestExpandIscsiFeatures_OmittedOptionalFields(t *testing.T) {
	for _, tc := range []struct {
		name       string
		input      map[string]interface{}
		wantSecret bool
		wantAuth   *volumesClient.AuthenticationType
	}{
		{
			name:       "enabled_authentications omitted",
			input:      map[string]interface{}{"target_secret": "s3cret", "enabled_authentications": ""},
			wantSecret: true,
		},
		{
			name:       "both omitted",
			input:      map[string]interface{}{"target_secret": "", "enabled_authentications": ""},
			wantSecret: false,
		},
		{
			name:       "unknown enum value",
			input:      map[string]interface{}{"target_secret": "", "enabled_authentications": "NOT_A_REAL_VALUE"},
			wantSecret: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := expandIscsiFeatures([]interface{}{tc.input})
			if got == nil {
				t.Fatal("expandIscsiFeatures returned nil")
			}
			if tc.wantSecret && got.TargetSecret == nil {
				t.Error("target_secret should have been set")
			}
			if !tc.wantSecret && got.TargetSecret != nil {
				t.Errorf("target_secret should be unset, got %q", *got.TargetSecret)
			}
			if got.EnabledAuthentications != nil {
				t.Errorf("enabled_authentications should be unset, got %v", *got.EnabledAuthentications)
			}
		})
	}
}

// TestExpandIscsiFeatures_ValidEnum pins the happy path so the guard did not break
// legitimate CHAP configuration.
func TestExpandIscsiFeatures_ValidEnum(t *testing.T) {
	for _, name := range []string{"CHAP", "NONE"} {
		t.Run(name, func(t *testing.T) {
			got := expandIscsiFeatures([]interface{}{
				map[string]interface{}{"target_secret": "s3cret", "enabled_authentications": name},
			})
			if got == nil {
				t.Fatal("expandIscsiFeatures returned nil")
			}
			if got.EnabledAuthentications == nil {
				t.Fatalf("enabled_authentications should be set for %q", name)
			}
			if got.EnabledAuthentications.GetName() != name {
				t.Errorf("enabled_authentications = %q, want %q", got.EnabledAuthentications.GetName(), name)
			}
			if got.TargetSecret == nil || *got.TargetSecret != "s3cret" {
				t.Error("target_secret was not preserved")
			}
		})
	}
}

func TestExpandIscsiFeatures_EmptyList(t *testing.T) {
	if got := expandIscsiFeatures([]interface{}{}); got != nil {
		t.Errorf("expandIscsiFeatures([]) = %#v, want nil", got)
	}
}
