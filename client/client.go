package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	libraryVersion      = "v3"
	DefaultBaseURL      = "https://%s:%s/api/nutanix/" + libraryVersion
	UserAgent           = "nutanix/" + libraryVersion
	mediaTypeJSON       = "application/json"
	mediaTypeWSDL       = "application/wsdl+xml"
	mediaTypeURLEncoded = "application/x-www-form-urlencoded"
)

// BuildRequestHandler creates a new request and marshals the body depending on the implementation
type BuildRequestHandler func(v interface{}, method, url string) (*http.Request, io.ReadSeeker, error)

// MarshalHander marshals the incoming body to a desired format
type MarshalHander func(v interface{}, action, version string) (string, error)

// UnmarshalHandler unmarshals the body request depending on different implementations
type UnmarshalHandler func(v interface{}, req *http.Response) error

// UnmarshalErrorHandler unmarshals the errors coming from an http respose
type UnmarshalErrorHandler func(r *http.Response) error

// Client manages the communication between the Outscale API's
type Client struct {
	Config Config

	// Handlers
	MarshalHander         MarshalHander
	BuildRequestHandler   BuildRequestHandler
	UnmarshalHandler      UnmarshalHandler
	UnmarshalErrorHandler UnmarshalErrorHandler
}

// Config Configuration of the client
type Config struct {
	Target      string
	Credentials *Credentials

	// HTTP client used to communicate with the Outscale API.
	Client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Optional function called after every successful request made to the DO APIs
	onRequestCompleted RequestCompletionCallback
}

// Credentials needed access key, secret key and region
type Credentials struct {
	Endpoint string
	Username string
	Password string
	Insecure string
	Port     string
	URL      string
}

// RequestCompletionCallback defines the type of the request callback function.
type RequestCompletionCallback func(*http.Request, *http.Response)

// NewRequest creates a request and signs it
func (c *Client) NewRequest(ctx context.Context, operation, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, errp := url.Parse(urlStr)
	if errp != nil {
		return nil, errp
	}

	u := c.Config.BaseURL.ResolveReference(rel)

	req, _, err := c.BuildRequestHandler(body, method, u.String())
	if err != nil {
		return nil, err
	}

	fmt.Println(rel.Opaque)

	return req, nil
}

// SetHeaders sets the headers for the request
func (c Client) SetHeaders(req *http.Request, headers []map[string]string) {
	for _, v := range headers {
		for h := range v {
			req.Header.Add(h, v[h])
		}
	}
}

// Do sends the request to the API's
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {

	req = req.WithContext(ctx)

	resp, err := c.Config.Client.Do(req)
	requestDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\n\n[DEBUG RESP]\n")
	fmt.Println(string(requestDump))
	if err != nil {
		return err
	}

	err = c.checkResponse(resp)
	if err != nil {
		return err
	}

	return c.UnmarshalHandler(v, resp)
}

func (c Client) checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	return c.UnmarshalErrorHandler(r)
}
