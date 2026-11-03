package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// PrintToJSON method helper to debug responses
func PrintToJSON(v interface{}, msg string) {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	log.Print("\n", msg, string(pretty))
	fmt.Print("\n", msg, string(pretty))
}

func ToJSONString(v interface{}) string {
	pretty, _ := json.MarshalIndent(v, "", "  ")

	return string(pretty)
}

// DebugRequest ...
func DebugRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Printf("[WARN] Error getting request's dump: %s\n", err)
	}

	log.Printf("[DEBUG] %s\n", string(requestDump))
}

// DebugResponse ...
func DebugResponse(res *http.Response) {
	requestDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Printf("[WARN] Error getting response's dump: %s\n", err)
	}

	log.Printf("[DEBUG] %s\n", string(requestDump))
}

func ConvertMapString(o map[string]interface{}) map[string]string {
	converted := make(map[string]string)
	for k, v := range o {
		converted[k] = fmt.Sprintf("%s", v.(string))
	}

	return converted
}

func StringLowerCaseValidateFunc(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if !(strings.ToLower(v) == v) {
		errs = append(errs, fmt.Errorf("%q must be in lowercase, got: %s", key, v))
	}
	return
}

func GenUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Never call os.Exit (log.Fatal) from a Terraform plugin: it kills the RPC
		// server, so Terraform reports an opaque "plugin crashed" with no diagnostic
		// and no state write. A panic at least surfaces through the plugin framework.
		panic(fmt.Sprintf("nutanix: cannot read random bytes for UUID: %v", err))
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func HashcodeString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// Strings hashes a list of strings to a unique hashcode.
func HashcodeStrings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", HashcodeString(buf.String()))
}

// ExtractErrorFromV4APIResponse pulls the human-readable message out of a v4 SDK error,
// whose Error() string is normally a JSON document shaped like
// {"data": {"error": [{"message": "..."}]}}.
//
// Every step is checked: this runs only on the failure path, where the payload is least
// likely to match expectations. Falling back to the raw error text is always better than
// panicking inside an error handler and taking the whole provider down.
func ExtractErrorFromV4APIResponse(err error) string {
	if err == nil {
		return ""
	}
	raw := err.Error()

	var errordata map[string]interface{}
	if e := json.Unmarshal([]byte(raw), &errordata); e != nil {
		return raw
	}

	data, ok := errordata["data"].(map[string]interface{})
	if !ok {
		return raw
	}
	errorList, ok := data["error"].([]interface{})
	if !ok || len(errorList) == 0 {
		return raw
	}
	first, ok := errorList[0].(map[string]interface{})
	if !ok {
		return raw
	}
	message, ok := first["message"].(string)
	if !ok {
		return raw
	}

	return message
}

// waitTimeoutGetter lets utils consume *conns.Client without importing the
// nutanix package (which would create an import cycle). Implementers must
// return the value in MINUTES; 0 (or negative) means "unset".
type waitTimeoutGetter interface {
	GetWaitTimeout() int64
}

// ResolveWaitTimeout returns the effective wait duration for a long-running
// operation. If meta exposes GetWaitTimeout() and the configured value is
// positive, it is interpreted as minutes and returned. Otherwise the caller's
// defaultTimeout is returned unchanged (typically d.Timeout(schema.TimeoutXxx)).
//
// Unit contract: provider wait_timeout is expressed in MINUTES.
func ResolveWaitTimeout(meta interface{}, defaultTimeout time.Duration) time.Duration {
	if c, ok := meta.(waitTimeoutGetter); ok {
		if wt := c.GetWaitTimeout(); wt > 0 {
			return time.Duration(wt) * time.Minute
		}
	}
	return defaultTimeout
}
