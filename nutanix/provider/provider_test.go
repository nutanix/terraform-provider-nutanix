package provider

import (
	"net/http"
	"os"
	"testing"
)

func TestParseHeadersFromEnv(t *testing.T) {
	// Save original environment and restore after test
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			for i := 0; i < len(env); i++ {
				if env[i] == '=' {
					os.Setenv(env[:i], env[i+1:])
					break
				}
			}
		}
	}()

	tests := []struct {
		name           string
		envVars        map[string]string
		expectedHeader string
		expectedValue  string
	}{
		{
			name: "cloudflare client id",
			envVars: map[string]string{
				"NUTANIX_HEADER_CF_ACCESS_CLIENT_ID": "test-client-id",
			},
			expectedHeader: "Cf-Access-Client-Id",
			expectedValue:  "test-client-id",
		},
		{
			name: "cloudflare client secret",
			envVars: map[string]string{
				"NUTANIX_HEADER_CF_ACCESS_CLIENT_SECRET": "test-secret",
			},
			expectedHeader: "Cf-Access-Client-Secret",
			expectedValue:  "test-secret",
		},
		{
			name: "custom x header",
			envVars: map[string]string{
				"NUTANIX_HEADER_X_CUSTOM_HEADER": "custom-value",
			},
			expectedHeader: "X-Custom-Header",
			expectedValue:  "custom-value",
		},
		{
			name: "authorization header",
			envVars: map[string]string{
				"NUTANIX_HEADER_AUTHORIZATION": "Bearer token123",
			},
			expectedHeader: "Authorization",
			expectedValue:  "Bearer token123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear relevant env vars
			for key := range tt.envVars {
				os.Unsetenv(key)
			}

			// Set test env vars
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Parse headers using the same logic as providerConfigure
			customHeaders := parseHeadersFromEnv()

			// Check the result
			if val, ok := customHeaders[tt.expectedHeader]; !ok {
				t.Errorf("Expected header %q not found in parsed headers", tt.expectedHeader)
			} else if val != tt.expectedValue {
				t.Errorf("Header %q = %q, expected %q", tt.expectedHeader, val, tt.expectedValue)
			}

			// Cleanup
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

func TestParseHeadersFromEnv_MultipleHeaders(t *testing.T) {
	// Save original environment and restore after test
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			for i := 0; i < len(env); i++ {
				if env[i] == '=' {
					os.Setenv(env[:i], env[i+1:])
					break
				}
			}
		}
	}()

	// Set multiple header env vars
	os.Setenv("NUTANIX_HEADER_CF_ACCESS_CLIENT_ID", "client-id")
	os.Setenv("NUTANIX_HEADER_CF_ACCESS_CLIENT_SECRET", "client-secret")
	os.Setenv("NUTANIX_HEADER_X_CUSTOM", "custom-value")

	customHeaders := parseHeadersFromEnv()

	expected := map[string]string{
		"Cf-Access-Client-Id":     "client-id",
		"Cf-Access-Client-Secret": "client-secret",
		"X-Custom":                "custom-value",
	}

	for header, expectedValue := range expected {
		if val, ok := customHeaders[header]; !ok {
			t.Errorf("Expected header %q not found", header)
		} else if val != expectedValue {
			t.Errorf("Header %q = %q, expected %q", header, val, expectedValue)
		}
	}

	// Cleanup
	os.Unsetenv("NUTANIX_HEADER_CF_ACCESS_CLIENT_ID")
	os.Unsetenv("NUTANIX_HEADER_CF_ACCESS_CLIENT_SECRET")
	os.Unsetenv("NUTANIX_HEADER_X_CUSTOM")
}

func TestParseHeadersFromEnv_IgnoresNonHeaderEnvVars(t *testing.T) {
	// Save original environment and restore after test
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			for i := 0; i < len(env); i++ {
				if env[i] == '=' {
					os.Setenv(env[:i], env[i+1:])
					break
				}
			}
		}
	}()

	// Clear environment so ambient NUTANIX_HEADER_* vars don't leak in.
	os.Clearenv()

	// Set various env vars, only one should be picked up
	os.Setenv("NUTANIX_HEADER_X_VALID", "valid")
	os.Setenv("NUTANIX_USERNAME", "should-ignore")
	os.Setenv("NUTANIX_PASSWORD", "should-ignore")
	os.Setenv("OTHER_HEADER_X_TEST", "should-ignore")

	customHeaders := parseHeadersFromEnv()

	if len(customHeaders) != 1 {
		t.Errorf("Expected 1 header, got %d: %v", len(customHeaders), customHeaders)
	}

	if val, ok := customHeaders["X-Valid"]; !ok || val != "valid" {
		t.Errorf("Expected X-Valid header with value 'valid', got %v", customHeaders)
	}

	// Cleanup
	os.Unsetenv("NUTANIX_HEADER_X_VALID")
	os.Unsetenv("NUTANIX_USERNAME")
	os.Unsetenv("NUTANIX_PASSWORD")
	os.Unsetenv("OTHER_HEADER_X_TEST")
}

// parseHeadersFromEnv extracts the header parsing logic for testing
// This mirrors the logic in providerConfigure
func parseHeadersFromEnv() map[string]string {
	customHeaders := make(map[string]string)
	const headerPrefix = "NUTANIX_HEADER_"

	for _, env := range os.Environ() {
		if len(env) > len(headerPrefix) && env[:len(headerPrefix)] == headerPrefix {
			for i := 0; i < len(env); i++ {
				if env[i] == '=' {
					envName := env[:i]
					envValue := env[i+1:]
					// Strip prefix, replace underscores with dashes, and title-case
					headerName := envName[len(headerPrefix):]
					// Replace underscores with dashes
					result := make([]byte, len(headerName))
					for j := 0; j < len(headerName); j++ {
						if headerName[j] == '_' {
							result[j] = '-'
						} else {
							result[j] = headerName[j]
						}
					}
					headerName = string(result)
					headerName = http.CanonicalHeaderKey(headerName)
					customHeaders[headerName] = envValue
					break
				}
			}
		}
	}
	return customHeaders
}
