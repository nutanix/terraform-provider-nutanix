package handler

import (
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	input := "Action=DescribeInstances&InstanceId.1=i-76536489&Version=2017-12-15"
	inputURL := "http://localhost"

	// Test Post
	_, _, err := BuildURLEncodedRequest(input, http.MethodPost, inputURL)
	if err != nil {
		t.Fatalf("Got error(%s)", err)
	}

	// Test Get
	req, _, err := BuildURLEncodedRequest(input, http.MethodGet, inputURL)
	if err != nil {
		t.Fatalf("Got error(%s)", err)
	}

	if req.URL.RawQuery != input {
		t.Fatalf("req.URL.RawQuery(%s) Got(%s)", req.URL.RawQuery, input)
	}

	// Test Unsupported
	_, _, err = BuildURLEncodedRequest(input, http.MethodDelete, inputURL)
	if err == nil {
		t.Fatalf("Got error(%s)", err)
	}

}
