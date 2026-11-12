package utils

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSanitizeJSONForLogging_RedactsSensitiveFields(t *testing.T) {
	input := []byte(`{
		"adminInfo": {
			"username": "tenant-admin",
			"password": "Nutanix.123",
			"email": "a@example.com"
		},
		"api_key": "raw-api-key",
		"accessKey": "ak",
		"secret_key": "sk",
		"private_key": "pk",
		"client_secret": "cs",
		"refresh_token": "rt",
		"authorization": "Bearer xyz",
		"ssl_key_passwd": "sslpass",
		"auth": "LOCAL",
		"key": "category-key",
		"name": "example-tenant",
		"items": [
			{"passwd": "p1", "label": "ok"},
			{"nested": {"user_password": "p2"}}
		]
	}`)

	out := SanitizeJSONForLogging(input)
	if string(out) == redactedPlaceholder {
		t.Fatalf("expected sanitized JSON, got fail-closed placeholder")
	}

	var v map[string]interface{}
	if err := json.Unmarshal(out, &v); err != nil {
		t.Fatalf("sanitized output is not JSON: %v\n%s", err, out)
	}

	admin := v["adminInfo"].(map[string]interface{})
	if admin["password"] != redactedPlaceholder {
		t.Errorf("password not redacted: %v", admin["password"])
	}
	if admin["username"] != "tenant-admin" {
		t.Errorf("username should be preserved, got %v", admin["username"])
	}
	if admin["email"] != "a@example.com" {
		t.Errorf("email should be preserved, got %v", admin["email"])
	}

	for _, k := range []string{
		"api_key", "accessKey", "secret_key", "private_key",
		"client_secret", "refresh_token", "authorization", "ssl_key_passwd",
	} {
		if v[k] != redactedPlaceholder {
			t.Errorf("%s not redacted: %v", k, v[k])
		}
	}

	// Non-sensitive exact keys must stay.
	if v["auth"] != "LOCAL" {
		t.Errorf("auth enum should be preserved, got %v", v["auth"])
	}
	if v["key"] != "category-key" {
		t.Errorf("key should be preserved, got %v", v["key"])
	}
	if v["name"] != "example-tenant" {
		t.Errorf("name should be preserved, got %v", v["name"])
	}

	items := v["items"].([]interface{})
	item0 := items[0].(map[string]interface{})
	if item0["passwd"] != redactedPlaceholder {
		t.Errorf("items[0].passwd not redacted: %v", item0["passwd"])
	}
	if item0["label"] != "ok" {
		t.Errorf("items[0].label should be preserved, got %v", item0["label"])
	}
	nested := items[1].(map[string]interface{})["nested"].(map[string]interface{})
	if nested["user_password"] != redactedPlaceholder {
		t.Errorf("nested user_password not redacted: %v", nested["user_password"])
	}

	// Must not contain cleartext secrets in the serialized form.
	s := string(out)
	for _, secret := range []string{
		"Nutanix.123", "raw-api-key", "sslpass", "Bearer xyz", "p1", "p2",
	} {
		if strings.Contains(s, secret) {
			t.Errorf("sanitized output still contains secret %q", secret)
		}
	}
}

func TestSanitizeJSONForLogging_FailClosedOnInvalidJSON(t *testing.T) {
	out := SanitizeJSONForLogging([]byte(`{"password": "leak",`))
	if string(out) != redactedPlaceholder {
		t.Fatalf("expected fail-closed placeholder, got %q", out)
	}
}

func TestSanitizeForLogging(t *testing.T) {
	payload := map[string]interface{}{
		"username": "admin",
		"password": "Nutanix.123",
	}
	out := SanitizeForLogging(payload)
	if !strings.Contains(out, `"username": "admin"`) {
		t.Fatalf("expected username preserved, got %s", out)
	}
	if !strings.Contains(out, `"password": "*******"`) {
		t.Fatalf("expected password redacted, got %s", out)
	}
	if strings.Contains(out, "Nutanix.123") {
		t.Fatalf("cleartext password leaked: %s", out)
	}
}

func TestIsSensitiveKey(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"password", true},
		{"Password", true},
		{"admin_password", true},
		{"adminPassword", true},
		{"passwordHash", true},
		{"api_key", true},
		{"api-key", true},
		{"apiKey", true},
		{"client_secret", true},
		{"refresh_token", true},
		{"authorization", true},
		// Newly added sensitive keys
		{"privkey", true},
		{"privKey", true},
		{"private_key", true}, // already covered via privatekey substring
		{"certificate", true},
		{"Certificate", true},
		{"ssl_certificate", true},
		{"tls_certificate", true},
		{"passphrase", true},
		{"Passphrase", true},
		{"pass_phrase", true},
		{"passPhrase", true},
		// Existing negative tests (should still work)
		{"auth", false},
		{"key", false},
		{"username", false},
		{"name", false},
		{"email", false},
	}
	for _, tc := range cases {
		if got := isSensitiveKey(tc.key); got != tc.want {
			t.Errorf("isSensitiveKey(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}
