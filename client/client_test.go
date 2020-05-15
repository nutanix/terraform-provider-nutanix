package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func setup() (*http.ServeMux, *Client, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client, _ := NewClient(&Credentials{"", "username", "password", "", "", true, false, ""})
	client.BaseURL, _ = url.Parse(server.URL)

	return mux, client, server
}

func TestNewClient(t *testing.T) {
	c, err := NewClient(&Credentials{"foo.com", "username", "password", "", "", true, false, ""})

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	expectedURL := fmt.Sprintf(defaultBaseURL, "foo.com")

	if c.BaseURL == nil || c.BaseURL.String() != expectedURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL, expectedURL)
	}

	if c.UserAgent != userAgent {
		t.Errorf("NewClient UserAgent = %v, expected %v", c.UserAgent, userAgent)
	}
}

func TestNewRequest(t *testing.T) {
	c, err := NewClient(&Credentials{"foo.com", "username", "password", "", "", true, false, ""})

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	inURL, outURL := "/foo", fmt.Sprintf(defaultBaseURL+absolutePath+"/foo", "foo.com")
	inBody, outBody := map[string]interface{}{"name": "bar"}, `{"name":"bar"}`+"\n"

	req, _ := c.NewRequest(context.TODO(), http.MethodPost, inURL, inBody)

	// test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// test body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest(%v) Body = %v, expected %v", inBody, string(body), outBody)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	messageResource := MessageResource{Message: "This field may not be blank."}
	messageList := make([]MessageResource, 1)
	messageList[0] = messageResource

	err := ErrorResponse{MessageList: messageList}

	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}
}

func TestGetResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(
			`{"api_version": "3.1", "code": 400, "kind": "error", "message_list":
				 [{"message": "This field may not be blank."}], "state": "ERROR"}`)),
	}

	err := CheckResponse(res)

	if err == nil {
		t.Fatal("Expected error response.")
	}

	if !strings.Contains(fmt.Sprint(err), "This field may not be blank.") {
		t.Errorf("error = %#v, expected %#v", err, "This field may not be blank.")
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(
			`{"api_version": "3.1", "code": 400, "kind": "error", "message_list":
				 [{"message": "This field may not be blank."}], "state": "ERROR"}`)),
	}
	err := CheckResponse(res)

	if err == nil {
		t.Fatalf("Expected error response.")
	}

	if !strings.Contains(fmt.Sprint(err), "This field may not be blank.") {
		t.Errorf("error = %#v, expected %#v", err, "This field may not be blank.")
	}
}

func TestDo(t *testing.T) {
	ctx := context.TODO()
	mux, client, server := setup()

	defer server.Close()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}

		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	body := new(foo)

	err := client.Do(context.Background(), req, body)

	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

	expected := &foo{"a"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body = %v, expected %v", body, expected)
	}
}

func TestDo_httpError(t *testing.T) {
	ctx := context.TODO()
	mux, client, server := setup()

	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

// / Test handling of an error caused by the internal http client's Do()
// function.
func TestDo_redirectLoop(t *testing.T) {
	ctx := context.TODO()
	mux, client, server := setup()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

// func TestDo_completion_callback(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	type foo struct {
// 		A string
// 	}

// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		if m := http.MethodGet; m != r.Method {
// 			t.Errorf("Request method = %v, expected %v", r.Method, m)
// 		}
// 		fmt.Fprint(w, `{"A":"a"}`)
// 	})

// 	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
// 	req = req.WithContext(ctx)
// 	body := new(foo)

// 	// var completedReq *http.Request
// 	var completedResp string

// 	client.OnRequestCompleted(func(req *http.Request, resp *http.Response, v interface{}) {
// 		// completedReq = req
// 		b, err := httputil.DumpResponse(resp, true)
// 		if err != nil {
// 			t.Errorf("Failed to dump response: %s", err)
// 		}
// 		completedResp = string(b)
// 	})
// 	err := client.Do(ctx, req, body)

// 	if err != nil {
// 		t.Fatalf("Do(): %v", err)
// 	}

// 	// if !reflect.DeepEqual(req., completedReq) {
// 	// 	t.Errorf("Completed request = %v, expected %v", completedReq, req)
// 	// }

// 	expected := `{"A":"a"}`

// 	if !strings.Contains(completedResp, expected) {
// 		t.Errorf("expected response to contain %v, Response = %v", expected, completedResp)
// 	}
// }

func TestClient_NewRequest(t *testing.T) {
	type fields struct {
		Credentials        *Credentials
		client             *http.Client
		BaseURL            *url.URL
		UserAgent          string
		onRequestCompleted RequestCompletionCallback
	}
	type args struct {
		ctx    context.Context
		method string
		urlStr string
		body   interface{}
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Credentials:        tt.fields.Credentials,
				client:             tt.fields.client,
				BaseURL:            tt.fields.BaseURL,
				UserAgent:          tt.fields.UserAgent,
				onRequestCompleted: tt.fields.onRequestCompleted,
			}
			got, err := c.NewRequest(tt.args.ctx, tt.args.method, tt.args.urlStr, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.NewRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_NewUploadRequest(t *testing.T) {
	type fields struct {
		Credentials        *Credentials
		client             *http.Client
		BaseURL            *url.URL
		UserAgent          string
		onRequestCompleted RequestCompletionCallback
	}
	type args struct {
		ctx    context.Context
		method string
		urlStr string
		body   []byte
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Credentials:        tt.fields.Credentials,
				client:             tt.fields.client,
				BaseURL:            tt.fields.BaseURL,
				UserAgent:          tt.fields.UserAgent,
				onRequestCompleted: tt.fields.onRequestCompleted,
			}
			got, err := c.NewUploadRequest(tt.args.ctx, tt.args.method, tt.args.urlStr, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.NewUploadRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.NewUploadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_OnRequestCompleted(t *testing.T) {
	type fields struct {
		Credentials        *Credentials
		client             *http.Client
		BaseURL            *url.URL
		UserAgent          string
		onRequestCompleted RequestCompletionCallback
	}
	type args struct {
		rc RequestCompletionCallback
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Credentials:        tt.fields.Credentials,
				client:             tt.fields.client,
				BaseURL:            tt.fields.BaseURL,
				UserAgent:          tt.fields.UserAgent,
				onRequestCompleted: tt.fields.onRequestCompleted,
			}
			c.OnRequestCompleted(tt.args.rc)
		})
	}
}

func TestClient_Do(t *testing.T) {
	type fields struct {
		Credentials        *Credentials
		client             *http.Client
		BaseURL            *url.URL
		UserAgent          string
		onRequestCompleted RequestCompletionCallback
	}
	type args struct {
		ctx context.Context
		req *http.Request
		v   interface{}
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Credentials:        tt.fields.Credentials,
				client:             tt.fields.client,
				BaseURL:            tt.fields.BaseURL,
				UserAgent:          tt.fields.UserAgent,
				onRequestCompleted: tt.fields.onRequestCompleted,
			}
			if err := c.Do(tt.args.ctx, tt.args.req, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Client.Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fillStruct(t *testing.T) {
	type args struct {
		data   map[string]interface{}
		result interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := fillStruct(tt.args.data, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("fillStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
