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

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = NewClient(&Credentials{"", "username", "password", "", "", true})
	client.BaseURL, _ = url.Parse(server.URL)
}

func teardown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	u := "foo.com"
	c, err := NewClient(&Credentials{u, "username", "password", "", "", true})

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	expectedURL := fmt.Sprintf(defaultBaseURL, u)

	if c.BaseURL == nil || c.BaseURL.String() != expectedURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL, expectedURL)
	}

	if c.UserAgent != userAgent {
		t.Errorf("NewClient UserAgent = %v, expected %v", c.UserAgent, userAgent)
	}
}

func TestNewRequest(t *testing.T) {
	u := "foo.com"
	c, err := NewClient(&Credentials{u, "username", "password", "", "", true})

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	inURL, outURL := "/foo", fmt.Sprintf(defaultBaseURL+absolutePath+"/foo", u)
	inBody, outBody := map[string]interface{}{"name": "bar"}, `{"name":"bar"}`+"\n"

	req, _ := c.NewRequest(ctx, http.MethodPost, inURL, inBody)

	//test relative URL was expanded
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
		Body:       ioutil.NopCloser(strings.NewReader(`{"api_version": "3.0", "code": 400, "kind": "error", "message_list": [{"message": "This field may not be blank."}], "state": "ERROR"}`)),
	}

	err := CheckResponse(res)

	if err == nil {
		t.Fatal("Expected error response.")
	}

	if !strings.Contains(fmt.Sprint(err), "This field may not be blank.") {
		t.Errorf("Error = %#v, expected %#v", err, "This field may not be blank.")
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(`{"api_version": "3.0", "code": 400, "kind": "error", "message_list": [{"message": "This field may not be blank."}], "state": "ERROR"}`)),
	}
	err := CheckResponse(res)

	if err == nil {
		t.Fatalf("Expected error response.")
	}

	if !strings.Contains(fmt.Sprint(err), "This field may not be blank.") {
		t.Errorf("Error = %#v, expected %#v", err, "This field may not be blank.")
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

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
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

/// Test handling of an error caused by the internal http client's Do()
// function.
func TestDo_redirectLoop(t *testing.T) {
	setup()
	defer teardown()

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

// 	//var completedReq *http.Request
// 	var completedResp string

// 	client.OnRequestCompleted(func(req *http.Request, resp *http.Response, v interface{}) {
// 		//completedReq = req
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
