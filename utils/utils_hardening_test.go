package utils

import (
	"errors"
	"strings"
	"testing"
)

// TestExtractErrorFromV4APIResponse_MalformedPayloads is the regression test for the
// four unchecked type assertions in the original implementation. The function runs
// only on the failure path, so a payload that did not match the expected shape used
// to panic inside the error handler and crash the provider.
func TestExtractErrorFromV4APIResponse_MalformedPayloads(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		want string
	}{
		{
			name: "not json at all",
			err:  errors.New("connection refused"),
			want: "connection refused",
		},
		{
			name: "json but no data key",
			err:  errors.New(`{"message": "nope"}`),
			want: `{"message": "nope"}`,
		},
		{
			name: "data is not an object",
			err:  errors.New(`{"data": "oops"}`),
			want: `{"data": "oops"}`,
		},
		{
			name: "error is not an array",
			err:  errors.New(`{"data": {"error": "oops"}}`),
			want: `{"data": {"error": "oops"}}`,
		},
		{
			name: "error array is empty",
			err:  errors.New(`{"data": {"error": []}}`),
			want: `{"data": {"error": []}}`,
		},
		{
			name: "error entry is not an object",
			err:  errors.New(`{"data": {"error": ["oops"]}}`),
			want: `{"data": {"error": ["oops"]}}`,
		},
		{
			name: "message is not a string",
			err:  errors.New(`{"data": {"error": [{"message": 42}]}}`),
			want: `{"data": {"error": [{"message": 42}]}}`,
		},
		{
			name: "message key missing",
			err:  errors.New(`{"data": {"error": [{"code": 5}]}}`),
			want: `{"data": {"error": [{"code": 5}]}}`,
		},
		{
			name: "json array at top level",
			err:  errors.New(`[1,2,3]`),
			want: `[1,2,3]`,
		},
		{
			name: "nil error",
			err:  nil,
			want: "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ExtractErrorFromV4APIResponse(tc.err)
			if got != tc.want {
				t.Errorf("ExtractErrorFromV4APIResponse() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestExtractErrorFromV4APIResponse_WellFormed pins the happy path: the API message is
// still unwrapped when the payload has the documented shape.
func TestExtractErrorFromV4APIResponse_WellFormed(t *testing.T) {
	err := errors.New(`{"data": {"error": [{"message": "Entity does not exist"}]}}`)

	if got := ExtractErrorFromV4APIResponse(err); got != "Entity does not exist" {
		t.Errorf("ExtractErrorFromV4APIResponse() = %q, want %q", got, "Entity does not exist")
	}
}

// TestGenUUID checks shape and uniqueness. The point of the change under test was
// removing log.Fatal (os.Exit), which would have killed the plugin process.
func TestGenUUID(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		id := GenUUID()
		if len(id) != 36 {
			t.Fatalf("GenUUID() = %q, want a 36-character UUID", id)
		}
		if parts := strings.Split(id, "-"); len(parts) != 5 {
			t.Fatalf("GenUUID() = %q, want 5 dash-separated groups", id)
		}
		if seen[id] {
			t.Fatalf("GenUUID() returned a duplicate: %q", id)
		}
		seen[id] = true
	}
}
