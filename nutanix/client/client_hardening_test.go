package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// testCredentials builds Credentials with keyed fields so the tests keep compiling
// if the struct gains members.
func testCredentials(insecure bool) *Credentials {
	return &Credentials{
		URL:      "foo.com",
		Username: "username",
		Password: "password",
		Insecure: insecure,
	}
}

// TestNewBaseClient_DoesNotMutateDefaultClient guards the fix for clients sharing
// http.DefaultClient. Previously NewBaseClient assigned the global and overwrote its
// Transport, so the last client constructed dictated TLS and proxy behavior for every
// other client in the process.
func TestNewBaseClient_DoesNotMutateDefaultClient(t *testing.T) {
	originalTransport := http.DefaultClient.Transport

	if _, err := NewBaseClient(testCredentials(true), testAbsolutePath, false); err != nil {
		t.Fatalf("NewBaseClient returned error: %v", err)
	}

	if http.DefaultClient.Transport != originalTransport {
		t.Error("NewBaseClient mutated http.DefaultClient.Transport; each client must own its transport")
	}
}

// TestNewBaseClient_TransportsAreIndependent proves two clients built with different
// `insecure` settings do not share a transport, so one cannot weaken the other's TLS.
func TestNewBaseClient_TransportsAreIndependent(t *testing.T) {
	secure, err := NewBaseClient(testCredentials(false), testAbsolutePath, false)
	if err != nil {
		t.Fatalf("NewBaseClient(secure) returned error: %v", err)
	}
	insecure, err := NewBaseClient(testCredentials(true), testAbsolutePath, false)
	if err != nil {
		t.Fatalf("NewBaseClient(insecure) returned error: %v", err)
	}

	if secure.client == insecure.client {
		t.Fatal("two clients share the same *http.Client")
	}
	if secure.client.Transport == insecure.client.Transport {
		t.Error("two clients share the same Transport; TLS settings would leak between them")
	}
}

// TestFilter_NoEntitiesKey guards against the panic from an unchecked
// res["entities"].([]interface{}) assertion on a payload with no entities collection.
func TestFilter_NoEntitiesKey(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{"no entities key", `{"message": "something went wrong"}`},
		{"null entities", `{"entities": null}`},
		{"entities is object", `{"entities": {"unexpected": true}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			filters := []*AdditionalFilter{{Name: "name", Values: []string{"x"}}}

			out, err := filter(io.NopCloser(strings.NewReader(tc.body)), filters, nil)
			if err != nil {
				t.Fatalf("filter returned error: %v", err)
			}
			got, err := io.ReadAll(out)
			if err != nil {
				t.Fatalf("reading filtered body: %v", err)
			}
			if string(got) != tc.body {
				t.Errorf("body was modified: got %s, want %s", got, tc.body)
			}
		})
	}
}

// TestFilter_InvalidJSON verifies a malformed body surfaces an error instead of being
// silently ignored, which previously left `res` nil and panicked on the next line.
func TestFilter_InvalidJSON(t *testing.T) {
	filters := []*AdditionalFilter{{Name: "name", Values: []string{"x"}}}

	if _, err := filter(io.NopCloser(strings.NewReader(`not json`)), filters, nil); err == nil {
		t.Error("expected an error for a malformed body, got nil")
	}
}

// TestFilter_NonObjectEntity ensures a non-object entry inside entities is skipped
// rather than panicking on entity.(map[string]interface{}), and that a matching
// object alongside it still survives filtering.
//
// The search paths mirror the only production caller (ListSubnet), which always
// supplies explicit base paths.
func TestFilter_NonObjectEntity(t *testing.T) {
	filters := []*AdditionalFilter{{Name: "name", Values: []string{"keep"}}}
	body := `{"entities": ["a string", 42, {"status": {"name": "keep"}}, {"status": {"name": "drop"}}]}`

	out, err := filter(io.NopCloser(strings.NewReader(body)), filters, []string{"status"})
	if err != nil {
		t.Fatalf("filter returned error: %v", err)
	}
	got, err := io.ReadAll(out)
	if err != nil {
		t.Fatalf("reading filtered body: %v", err)
	}
	if !strings.Contains(string(got), `"keep"`) {
		t.Errorf("matching entity was dropped: %s", got)
	}
	if strings.Contains(string(got), `"drop"`) {
		t.Errorf("non-matching entity was kept: %s", got)
	}
}

// TestCheckResponse_NonMapStatus guards the unchecked status.(map[string]interface{})
// assertion, which panicked inside the error handler for scalar status values.
func TestCheckResponse_NonMapStatus(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{"numeric status", `{"status": 500}`},
		{"bool status", `{"status": true}`},
		{"array status", `{"status": ["ERROR"]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader(tc.body)),
			}
			// The assertion here is simply that this does not panic.
			_ = CheckResponse(resp)
		})
	}
}

// TestCheckResponse_StringStatusStillNil pins the pre-existing behavior that a string
// `status` is treated as a non-error payload.
func TestCheckResponse_StringStatusStillNil(t *testing.T) {
	resp := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(`{"status": "COMPLETE"}`)),
	}
	if err := CheckResponse(resp); err != nil {
		t.Errorf("expected nil for a string status, got %v", err)
	}
}

// TestCheckResponse_EmptyBodyErrorStatus covers the case where a non-2xx response
// carried no body. It used to return nil, so callers proceeded with a zero-value
// response struct and dereferenced it.
func TestCheckResponse_EmptyBodyErrorStatus(t *testing.T) {
	resp := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("expected an error for a bodiless 500, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should mention the status code, got %v", err)
	}
}

// TestCheckResponse_NotFoundIsTyped verifies 404 is distinguishable via errors.Is so
// resources can tell "deleted" apart from "request failed".
func TestCheckResponse_NotFoundIsTyped(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{"empty body", ""},
		{"with body", `{"message": "not found"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader(tc.body)),
			}
			err := CheckResponse(resp)
			if err == nil {
				t.Fatal("expected an error for 404, got nil")
			}
			if !errors.Is(err, ErrNotFound) {
				t.Errorf("404 error should match ErrNotFound, got %v", err)
			}
		})
	}
}

// TestCheckResponse_404PreservesEntityNotFoundText protects a compatibility contract
// that is easy to break by accident.
//
// Around thirty v3 resources detect deletion with
// strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND"). If a not-found error stops
// matching that string, those Read functions fall through to a branch that calls
// Delete on the resource and then hard-errors. Adding ErrNotFound must therefore keep
// the raw API body — which carries the reason code — inside the error text.
func TestCheckResponse_404PreservesEntityNotFoundText(t *testing.T) {
	body := `{"api_version":"3.1","code":404,"kind":"subnet",` +
		`"message_list":[{"message":"The specified subnet was not found.","reason":"ENTITY_NOT_FOUND"}],` +
		`"state":"ERROR"}`

	resp := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("expected an error for 404, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("new mechanism: errors.Is(err, ErrNotFound) = false, err = %v", err)
	}
	if !strings.Contains(err.Error(), "ENTITY_NOT_FOUND") {
		t.Errorf("existing v3 convention broken: %q no longer contains ENTITY_NOT_FOUND", err)
	}
}

// TestCheckResponse_ServerErrorIsNotNotFound makes sure a 500 is never mistaken for a
// deletion, which would silently drop resources from Terraform state.
func TestCheckResponse_ServerErrorIsNotNotFound(t *testing.T) {
	resp := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(`{"message": "boom"}`)),
	}
	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("expected an error for 500, got nil")
	}
	if errors.Is(err, ErrNotFound) {
		t.Error("a 500 must not be reported as ErrNotFound")
	}
}

// TestNewClient_SessionAuthFailureSurfaces covers the discarded CheckResponse result
// during the session-auth handshake: bad credentials produced a cookie-less client and
// a nil error.
func TestNewClient_SessionAuthFailureSurfaces(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/nutanix/v3/users/me", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parsing test server URL: %v", err)
	}

	creds := &Credentials{
		URL:         serverURL.Host,
		Username:    "username",
		Password:    "wrong",
		Insecure:    true,
		SessionAuth: true,
	}

	c, err := NewClient(creds, testUserAgent, testAbsolutePath, true)
	if err == nil {
		t.Fatal("expected an error when session authentication fails, got nil")
	}
	if c != nil && c.Cookies != nil {
		t.Error("cookies must not be set when the session handshake failed")
	}
}

// TestNewUnAuthRequest_EncodeErrorSurfaces covers the wrong-variable check that made
// body encoding failures silently produce an empty request body.
func TestNewUnAuthRequest_EncodeErrorSurfaces(t *testing.T) {
	_, c, server := setup()
	defer server.Close()

	// A channel cannot be marshaled to JSON.
	_, err := c.NewUnAuthRequest(context.Background(), http.MethodPost, "/x", make(chan int))
	if err == nil {
		t.Error("expected an error when the body cannot be encoded, got nil")
	}
}

func TestJoinHostPort(t *testing.T) {
	for _, tc := range []struct {
		name string
		host string
		port string
		want string
	}{
		{"ipv4", "10.0.0.1", "9440", "10.0.0.1:9440"},
		{"hostname", "pc.example.com", "9440", "pc.example.com:9440"},
		{"ipv6 gets brackets", "2001:db8::1", "9440", "[2001:db8::1]:9440"},
		{"ipv6 loopback", "::1", "9440", "[::1]:9440"},
		{"scheme prefixed stays untouched", "http://localhost", "8080", "http://localhost:8080"},
		{"empty port", "10.0.0.1", "", "10.0.0.1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := JoinHostPort(tc.host, tc.port); got != tc.want {
				t.Errorf("JoinHostPort(%q, %q) = %q, want %q", tc.host, tc.port, got, tc.want)
			}
		})
	}
}
