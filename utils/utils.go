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
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
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

// Extract error from v4 API response
func ExtractErrorFromV4APIResponse(err error) string {
	var errordata map[string]interface{}
	e := json.Unmarshal([]byte(err.Error()), &errordata)
	if e != nil {
		return e.Error()
	}
	data := errordata["data"].(map[string]interface{})
	errorList := data["error"].([]interface{})
	errorMessage := errorList[0].(map[string]interface{})["message"]
	return errorMessage.(string)
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
