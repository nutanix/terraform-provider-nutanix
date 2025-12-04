package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

const (
	testLibraryVersion = "v3"
	testAbsolutePath   = "api/nutanix/" + testLibraryVersion
	testUserAgent      = "nutanix/" + testLibraryVersion
	fileName           = "../sdks/v3/prism/prism.go"
)

func setup() (*http.ServeMux, *Client, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client, _ := NewClient(&Credentials{"", "username", "password", "", "", true, false, "", "", "", nil, "", "", ""}, testUserAgent, testAbsolutePath, false)
	client.BaseURL, _ = url.Parse(server.URL)

	return mux, client, server
}

func TestNewClient(t *testing.T) {
	c, err := NewClient(&Credentials{"foo.com", "username", "password", "", "", true, false, "", "", "", nil, "", "", ""}, testUserAgent, testAbsolutePath, false)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	expectedURL := fmt.Sprintf(defaultBaseURL, httpsPrefix, "foo.com")

	if c.BaseURL == nil || c.BaseURL.String() != expectedURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL, expectedURL)
	}

	if c.UserAgent != testUserAgent {
		t.Errorf("NewClient UserAgent = %v, expected %v", c.UserAgent, testUserAgent)
	}
}

func TestNewBaseClient(t *testing.T) {
	c, err := NewBaseClient(&Credentials{"foo.com", "username", "password", "", "", true, false, "", "", "", nil, "", "", ""}, testAbsolutePath, true)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	expectedURL := fmt.Sprintf(defaultBaseURL, httpPrefix, "foo.com")

	if c.BaseURL == nil || c.BaseURL.String() != expectedURL {
		t.Errorf("NewBaseClient BaseURL = %v, expected %v", c.BaseURL, expectedURL)
	}

	if c.AbsolutePath != testAbsolutePath {
		t.Errorf("NewBaseClient UserAgent = %v, expected %v", c.AbsolutePath, testAbsolutePath)
	}
}

func TestNewRequest(t *testing.T) {
	c, err := NewClient(&Credentials{"foo.com", "username", "password", "", "", true, false, "", "", "", nil, "", "", ""}, testUserAgent, testAbsolutePath, false)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	inURL, outURL := "/foo", fmt.Sprintf(defaultBaseURL+testAbsolutePath+"/foo", "https", "foo.com")
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

func TestNewUploadRequest(t *testing.T) {
	c, err := NewClient(&Credentials{"foo.com", "username", "password", "", "", true, false, "", "", "", nil, "", "", ""}, testUserAgent, testAbsolutePath, true)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	inURL, outURL := "/foo", fmt.Sprintf(defaultBaseURL+testAbsolutePath+"/foo", httpPrefix, "foo.com")
	inBody, _ := os.Open(fileName)
	if err != nil {
		t.Fatalf("Error opening file %v, error : %v", fileName, err)
	}

	// expected body
	out, _ := os.Open(fileName)
	outBody, _ := ioutil.ReadAll(out)

	req, err := c.NewUploadRequest(context.TODO(), http.MethodPost, inURL, inBody)
	if err != nil {
		t.Fatalf("NewUploadRequest() errored out with error : %v", err.Error())
	}
	// test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewUploadRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	got, _ := ioutil.ReadAll(req.Body)
	if !bytes.Equal(got, outBody) {
		t.Errorf("NewUploadRequest(%v) Body = %v, expected %v", inBody, string(got), string(outBody))
	}

	// test headers.
	inHeaders := map[string]string{
		"Content-Type": octetStreamType,
		"Accept":       mediaType,
	}
	for k, v := range inHeaders {
		if v != req.Header[k][0] {
			t.Errorf("NewUploadRequest() Header value for %v = %v, expected %v", k, req.Header[k][0], v)
		}
	}
}

func TestNewUnAuthRequest(t *testing.T) {
	c, err := NewClient(&Credentials{"foo.com", "username", "password", "", "", true, false, "", "", "", nil, "", "", ""}, testUserAgent, testAbsolutePath, true)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	inURL, outURL := "/foo", fmt.Sprintf(defaultBaseURL+testAbsolutePath+"/foo", httpPrefix, "foo.com")
	inBody, outBody := map[string]interface{}{"name": "bar"}, `{"name":"bar"}`+"\n"

	req, _ := c.NewUnAuthRequest(context.TODO(), http.MethodPost, inURL, inBody)

	// test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewUnAuthRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// test body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewUnAuthRequest(%v) Body = %v, expected %v", inBody, string(body), outBody)
	}

	// test headers. Authorization header shouldn't exist
	if _, ok := req.Header["Authorization"]; ok {
		t.Errorf("Unexpected Authorization header obtained in request from NewUnAuthRequest()")
	}
	inHeaders := map[string]string{
		"Content-Type": mediaType,
		"Accept":       mediaType,
		"User-Agent":   testUserAgent,
	}
	for k, v := range req.Header {
		if v[0] != inHeaders[k] {
			t.Errorf("NewUnAuthRequest() Header value for %v = %v, expected %v", k, v[0], inHeaders[k])
		}
	}
}

func TestNewUnAuthFormEncodedRequest(t *testing.T) {
	c, err := NewClient(&Credentials{"foo.com", "username", "password", "", "", true, false, "", "", "", nil, "", "", ""}, testUserAgent, testAbsolutePath, true)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	inURL, outURL := "/foo", fmt.Sprintf(defaultBaseURL+testAbsolutePath+"/foo", httpPrefix, "foo.com")
	inBody := map[string]string{"name": "bar", "fullname": "foobar"}
	outBody := map[string][]string{"name": {"bar"}, "fullname": {"foobar"}}

	req, _ := c.NewUnAuthFormEncodedRequest(context.TODO(), http.MethodPost, inURL, inBody)

	// test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewUnAuthFormEncodedRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// test body
	// Parse the body form data to a map structure which can be accessed by req.PostForm
	req.ParseForm()

	// check form encoded key-values as compared to input values
	if !reflect.DeepEqual(outBody, (map[string][]string)(req.PostForm)) {
		t.Errorf("NewUnAuthFormEncodedRequest(%v) Form encoded k-v, got = %v, expected %v", inBody, req.PostForm, outBody)
	}

	// test headers. Authorization header shouldn't exist
	if _, ok := req.Header["Authorization"]; ok {
		t.Errorf("Unexpected Authorization header obtained in request from NewUnAuthFormEncodedRequest()")
	}
	inHeaders := map[string]string{
		"Content-Type": formEncodedType,
		"Accept":       mediaType,
		"User-Agent":   testUserAgent,
	}
	for k, v := range req.Header {
		if v[0] != inHeaders[k] {
			t.Errorf("NewUnAuthFormEncodedRequest() Header value for %v = %v, expected %v", k, v[0], inHeaders[k])
		}
	}
}

func TestNewUnAuthUploadRequest(t *testing.T) {
	c, err := NewClient(&Credentials{"foo.com", "username", "password", "", "", true, false, "", "", "", nil, "", "", ""}, testUserAgent, testAbsolutePath, true)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	inURL, outURL := "/foo", fmt.Sprintf(defaultBaseURL+testAbsolutePath+"/foo", httpPrefix, "foo.com")
	inBody, _ := os.Open(fileName)
	if err != nil {
		t.Fatalf("Error opening fiele %v, error : %v", fileName, err)
	}

	// expected body
	out, _ := os.Open(fileName)
	outBody, _ := ioutil.ReadAll(out)

	req, err := c.NewUnAuthUploadRequest(context.TODO(), http.MethodPost, inURL, inBody)
	if err != nil {
		t.Fatalf("NewUnAuthUploadRequest() errored out with error : %v", err.Error())
	}
	// test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewUnAuthUploadRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	got, _ := ioutil.ReadAll(req.Body)
	if !bytes.Equal(got, outBody) {
		t.Errorf("NewUnAuthUploadRequest(%v) Body = %v, expected %v", inBody, string(got), string(outBody))
	}

	// test headers. Authorization header shouldn't exist
	if _, ok := req.Header["Authorization"]; ok {
		t.Errorf("Unexpected Authorization header obtained in request from NewUnAuthUploadRequest()")
	}
	inHeaders := map[string]string{
		"Content-Type": octetStreamType,
		"Accept":       mediaType,
	}
	for k, v := range inHeaders {
		if v != req.Header[k][0] {
			t.Errorf("NewUploadRequest() Header value for %v = %v, expected %v", k, req.Header[k][0], v)
		}
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
				 [{"message": "bad Request"}], "state": "ERROR"}`)),
	}

	err := CheckResponse(res)

	if err == nil {
		t.Fatal("Expected error response.")
	}

	if !strings.Contains(fmt.Sprint(err), "bad Request") {
		t.Errorf("error = %#v, expected %#v", err, "bad Request")
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(
			`{"api_version": "3.1", "code": 400, "kind": "error", "message_list":
				 [{"message": "bad Request"}], "state": "ERROR"}`)),
	}
	err := CheckResponse(res)

	if err == nil {
		t.Fatalf("Expected error response.")
	}

	if !strings.Contains(fmt.Sprint(err), "bad Request") {
		t.Errorf("error = %#v, expected %#v", err, "bad Request")
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

// *********** Filters tests ***********

func getEntity(name string, vlanID string, uuid string) string {
	return fmt.Sprintf(`{"spec":{"cluster_reference":{"uuid":"%s"},"name":"%s","resources":{"vlan_id":%s}}}`, uuid, name, vlanID)
}

func removeWhiteSpace(input string) string {
	whitespacePattern := regexp.MustCompile(`\s+`)
	return whitespacePattern.ReplaceAllString(input, "")
}

func getFilter(name string, values []string) []*AdditionalFilter {
	return []*AdditionalFilter{
		{
			Name:   name,
			Values: values,
		},
	}
}

func runTest(filters []*AdditionalFilter, inputString string, expected string) bool {
	input := io.NopCloser(strings.NewReader(inputString))
	fmt.Println(expected)
	baseSearchPaths := []string{"spec", "spec.resources"}
	filteredBody, err := filter(input, filters, baseSearchPaths)
	if err != nil {
		panic(err)
	}
	actualBytes, _ := io.ReadAll(filteredBody)
	actual := string(actualBytes)
	fmt.Println(actual)
	return actual == expected
}

func TestDoWithFilters_filter(t *testing.T) {
	entity1 := getEntity("subnet-01", "111", "012345-111")
	entity2 := getEntity("subnet-01", "112", "012345-112")
	entity3 := getEntity("subnet-02", "112", "012345-111")
	input := fmt.Sprintf(`{"entities":[%s,%s,%s]}`, entity1, entity2, entity3)

	filtersList := [][]*AdditionalFilter{
		getFilter("name", []string{"subnet-01", "subnet-03"}),
		getFilter("vlan_id", []string{"111", "subnet-03"}),
		getFilter("cluster_reference.uuid", []string{"111", "012345-112"}),
	}
	expectedList := []string{
		removeWhiteSpace(fmt.Sprintf(`{"entities":[%s,%s]}`, entity1, entity2)),
		removeWhiteSpace(fmt.Sprintf(`{"entities":[%s]}`, entity1)),
		removeWhiteSpace(fmt.Sprintf(`{"entities":[%s]}`, entity2)),
	}

	for i := 0; i < len(filtersList); i++ {
		if ok := runTest(filtersList[i], input, expectedList[i]); !ok {
			t.Fatal()
		}
	}
}

// *************************************

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
