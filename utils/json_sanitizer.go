package utils

import (
	"encoding/json"
	"strings"
)

const redactedPlaceholder = "*******"

// SanitizeForLogging marshals v to indented JSON, redacts sensitive fields, and
// returns a string suitable for log.Printf. On marshal failure it fail-closes.
//
//	log.Printf("[DEBUG] create payload: %s", utils.SanitizeForLogging(req))
func SanitizeForLogging(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return redactedPlaceholder
	}
	return string(SanitizeJSONForLogging(data))
}

// SanitizeJSONForLogging redacts sensitive fields in JSON data for safe logging.
// On parse or marshal failure it fail-closes and returns a redacted placeholder
// so secrets in malformed payloads are not logged.
func SanitizeJSONForLogging(data []byte) []byte {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return []byte(redactedPlaceholder)
	}
	sanitized := sanitizeValue(v)
	out, err := json.MarshalIndent(sanitized, "", "  ")
	if err != nil {
		return []byte(redactedPlaceholder)
	}
	return out
}

func sanitizeValue(v interface{}) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		for k, val := range v {
			if isSensitiveKey(k) {
				v[k] = redactedPlaceholder
			} else {
				v[k] = sanitizeValue(val)
			}
		}
	case []interface{}:
		for i, val := range v {
			v[i] = sanitizeValue(val)
		}
	default:
		// Do nothing for other types (string, number, bool, nil)
	}
	return v
}

// normalizeKey lowercases the key and strips '_' / '-' so api_key, api-key,
// and apiKey compare as the same token shape where applicable.
func normalizeKey(key string) string {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ReplaceAll(key, "-", "")
	return key
}

func isSensitiveKey(key string) bool {
	n := normalizeKey(key)

	// High-confidence substrings (covers passwordHash, adminPassword, api_key, …).
	for _, s := range []string{
		"password",
		"passwd",
		"privatekey",
		"privkey",
		"apikey",
		"accesskey",
		"secretkey",
		"sslkeypasswd",
		"credential",
		"certificate",
		"passphrase",
	} {
		if strings.Contains(n, s) {
			return true
		}
	}

	// Broader tokens: exact or suffix only (avoid matching bare "auth"/"key",
	// and reduce false positives like "secretary").
	for _, s := range []string{
		"secret",
		"token",
		"authorization",
		"bearer",
	} {
		if n == s || strings.HasSuffix(n, s) {
			return true
		}
	}

	return false
}
