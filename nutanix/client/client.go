package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

const (
	// libraryVersion = "v3"
	defaultBaseURL = "%s://%s/"
	httpPrefix     = "http"
	httpsPrefix    = "https"
	// absolutePath   = "api/nutanix/" + libraryVersion
	// userAgent      = "nutanix/" + libraryVersion
	mediaType       = "application/json"
	formEncodedType = "application/x-www-form-urlencoded"
	octetStreamType = "application/octet-stream"
)

// Client Config Configuration of the client
type Client struct {
	Credentials *Credentials

	// HTTP client used to communicate with the Nutanix API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	Cookies []*http.Cookie

	// Optional function called after every successful request made.
	onRequestCompleted RequestCompletionCallback

	// absolutePath: for example api/nutanix/v3
	AbsolutePath string

	// error message, incase client is in error state
	ErrorMsg string
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response, interface{})

// Credentials needed username and password (or API key)
type Credentials struct {
	URL                string
	Username           string
	Password           string
	Endpoint           string
	Port               string
	Insecure           bool
	SessionAuth        bool
	ProxyURL           string
	FoundationEndpoint string              // Required field for connecting to foundation VM APIs
	FoundationPort     string              // Port for connecting to foundation VM APIs
	RequiredFields     map[string][]string // RequiredFields is client to its required fields mapping for validations and usage in every client
	NdbEndpoint        string              // Required field for connecting to Era VM APIs.
	NdbUsername        string
	NdbPassword        string
	APIKey             string            // API key for authentication (alternative to username/password)
	CustomHeaders      map[string]string // Custom headers to add to all requests (e.g., for Cloudflare Access)
}

// AdditionalFilter specification for client side filters
type AdditionalFilter struct {
	Name   string
	Values []string
}

// JoinHostPort joins a configured endpoint and port into a URL authority.
//
// net.JoinHostPort is not used directly because it brackets any host containing a
// colon, which would also bracket endpoints that were configured with a scheme
// (e.g. "https://pc.example.com") and produce an unparseable URL. Only genuine IPv6
// literals require brackets, so that case is detected explicitly; every other input
// keeps the plain "host:port" form this provider has always produced.
func JoinHostPort(host, port string) string {
	if port == "" {
		return host
	}
	if ip := net.ParseIP(host); ip != nil && ip.To4() == nil {
		return "[" + host + "]:" + port
	}
	return host + ":" + port
}

// applyAuthHeaders adds the appropriate authentication header to the request
// Priority: Cookies (session auth) > API Key > Basic Auth
func (c *Client) applyAuthHeaders(req *http.Request) {
	if c.Cookies != nil {
		// Session-based authentication using cookies
		for _, cookie := range c.Cookies {
			req.AddCookie(cookie)
		}
	} else if c.Credentials.APIKey != "" {
		// API key authentication
		req.Header.Add("X-Ntnx-Api-Key", c.Credentials.APIKey)
	} else {
		// Basic authentication (username/password)
		req.Header.Add("Authorization", "Basic "+
			base64.StdEncoding.EncodeToString([]byte(c.Credentials.Username+":"+c.Credentials.Password)))
	}
}

// applyCustomHeaders adds any custom headers from credentials to the request
func (c *Client) applyCustomHeaders(req *http.Request) {
	for key, value := range c.Credentials.CustomHeaders {
		req.Header.Add(key, value)
	}
}

// NewClient returns a wrapper around http/https (as per isHTTP flag) client with additions of proxy & session_auth if given
func NewClient(credentials *Credentials, userAgent string, absolutePath string, isHTTP bool) (*Client, error) {
	if userAgent == "" {
		return nil, fmt.Errorf("userAgent argument must be passed")
	}
	if absolutePath == "" {
		return nil, fmt.Errorf("absolutePath argument must be passed")
	}

	// create base client with basic configs
	baseClient, err := NewBaseClient(credentials, absolutePath, isHTTP)
	if err != nil {
		return nil, err
	}
	// add useragent
	baseClient.UserAgent = userAgent

	if credentials.ProxyURL != "" {
		log.Printf("[DEBUG] Using proxy: %s\n", credentials.ProxyURL)
		proxy, err := url.Parse(credentials.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing proxy url: %s", err)
		}

		transCfg := newTransport(credentials.Insecure)
		transCfg.Proxy = http.ProxyURL(proxy)
		baseClient.client.Transport = logging.NewTransport("Nutanix", transCfg)
	}

	if credentials.SessionAuth {
		log.Printf("[DEBUG] Using session_auth\n")

		ctx := context.TODO()
		req, err := baseClient.NewRequest(ctx, http.MethodGet, "/users/me", nil)
		if err != nil {
			return baseClient, err
		}

		resp, err := baseClient.client.Do(req)
		if err != nil {
			return baseClient, err
		}
		defer resp.Body.Close()

		// A failed session handshake must surface here. Continuing would hand back a
		// cookie-less client that fails later with an unrelated-looking error.
		if err := CheckResponse(resp); err != nil {
			return baseClient, fmt.Errorf("session authentication failed: %w", err)
		}

		baseClient.Cookies = resp.Cookies()
	}

	return baseClient, nil
}

// newTransport builds a dedicated HTTP transport for a single client.
//
// Each client must own its transport: sharing one (for example via http.DefaultClient)
// would let one client's TLS or proxy settings silently apply to every other client in
// the process, including ones talking to a different host.
func newTransport(insecure bool) *http.Transport {
	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		transport = &http.Transport{}
	}
	transCfg := transport.Clone()
	//nolint:gosec // InsecureSkipVerify is the documented `insecure` provider opt-in.
	transCfg.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}
	return transCfg
}

// NewBaseClient returns a basic http/https client based on isHttp flag
func NewBaseClient(credentials *Credentials, absolutePath string, isHTTP bool) (*Client, error) {
	if absolutePath == "" {
		return nil, fmt.Errorf("absolutePath argument must be passed")
	}

	httpClient := &http.Client{
		Transport: logging.NewTransport("Nutanix", newTransport(credentials.Insecure)),
	}

	protocol := httpsPrefix
	if isHTTP {
		protocol = httpPrefix
	}

	baseURL, err := url.Parse(fmt.Sprintf(defaultBaseURL, protocol, credentials.URL))
	if err != nil {
		return nil, err
	}

	c := &Client{
		Credentials:  credentials,
		client:       httpClient,
		BaseURL:      baseURL,
		AbsolutePath: absolutePath,
	}

	return c, nil
}

// NewRequest creates a request
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	// check if client exists or not
	if c.client == nil {
		return nil, fmt.Errorf("%s", c.ErrorMsg)
	}

	rel, errp := url.Parse(c.AbsolutePath + urlStr)
	if errp != nil {
		return nil, errp
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)

	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	c.applyAuthHeaders(req)
	c.applyCustomHeaders(req)
	return req, nil
}

// NewUnAuthRequest creates a request without authorisation headers
func (c *Client) NewUnAuthRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	// check if client exists or not
	if c.client == nil {
		return nil, fmt.Errorf("%s", c.ErrorMsg)
	}

	rel, err := url.Parse(c.AbsolutePath + urlStr)
	if err != nil {
		return nil, err
	}
	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		if er := json.NewEncoder(buf).Encode(body); er != nil {
			return nil, er
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

// NewUnAuthFormEncodedRequest returns content-type: application/x-www-form-urlencoded based unauth request
func (c *Client) NewUnAuthFormEncodedRequest(ctx context.Context, method, urlStr string, body map[string]string) (*http.Request, error) {
	// check if client exists or not
	if c.client == nil {
		return nil, fmt.Errorf("%s", c.ErrorMsg)
	}

	rel, err := url.Parse(c.AbsolutePath + urlStr)
	if err != nil {
		return nil, err
	}
	u := c.BaseURL.ResolveReference(rel)

	// create form data from body
	data := url.Values{}
	for k, v := range body {
		data.Set(k, v)
	}

	// create a new request based on encoded from data
	req, err := http.NewRequest(method, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", formEncodedType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

// NewUploadRequest Handles image uploads for image service
func (c *Client) NewUploadRequest(ctx context.Context, method, urlStr string, fileReader *os.File) (*http.Request, error) {
	// check if client exists or not
	if c.client == nil {
		return nil, fmt.Errorf("%s", c.ErrorMsg)
	}
	rel, errp := url.Parse(c.AbsolutePath + urlStr)
	if errp != nil {
		return nil, errp
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), fileReader)
	if err != nil {
		return nil, err
	}

	// Set req.ContentLength and req.GetBody as internally there is no implementation of such for os.File type reader
	// http.NewRequest() only sets this values for types - bytes.Buffer, bytes.Reader and strings.Reader
	// Refer https://github.com/golang/go/blob/a0f77e56b7a7ecb92dca3e2afdd56ee773c2cb07/src/net/http/request.go#L896
	fileInfo, err := fileReader.Stat()
	if err != nil {
		return nil, err
	}
	req.ContentLength = fileInfo.Size()
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(fileReader), nil
	}

	req.Header.Add("Content-Type", octetStreamType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	c.applyAuthHeaders(req)
	c.applyCustomHeaders(req)

	return req, nil
}

// NewUnAuthUploadRequest handles image uploads for image service without auth
func (c *Client) NewUnAuthUploadRequest(ctx context.Context, method, urlStr string, fileReader *os.File) (*http.Request, error) {
	// check if client exists or not
	if c.client == nil {
		return nil, fmt.Errorf("%s", c.ErrorMsg)
	}
	rel, errp := url.Parse(c.AbsolutePath + urlStr)
	if errp != nil {
		return nil, errp
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), fileReader)
	if err != nil {
		return nil, err
	}

	// Set req.ContentLength and req.GetBody as internally there is no implementation of such for os.File type reader
	// http.NewRequest() only sets this values for types - bytes.Buffer, bytes.Reader and strings.Reader
	// Refer https://github.com/golang/go/blob/a0f77e56b7a7ecb92dca3e2afdd56ee773c2cb07/src/net/http/request.go#L896
	fileInfo, err := fileReader.Stat()
	if err != nil {
		return nil, err
	}
	req.ContentLength = fileInfo.Size()
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(fileReader), nil
	}

	req.Header.Add("Content-Type", octetStreamType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// OnRequestCompleted sets the DO API request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

// Do performs request passed
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
	// check if client exists or not
	if c.client == nil {
		return fmt.Errorf("%s", c.ErrorMsg)
	}

	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := CheckResponse(resp); err != nil {
		return err
	}

	if err := decodeResponse(resp, v); err != nil {
		return err
	}

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp, v)
	}
	return nil
}

// decodeResponse writes the response body into v, which may be an io.Writer or any
// json-decodable value. A nil v discards the body.
func decodeResponse(resp *http.Response, v interface{}) error {
	if v == nil {
		return nil
	}
	if w, ok := v.(io.Writer); ok {
		if _, err := io.Copy(w, resp.Body); err != nil {
			return fmt.Errorf("error copying response body: %w", err)
		}
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("error unmarshalling json: %s", err)
	}
	return nil
}

func searchSlice(slice []string, key string) bool {
	for _, v := range slice {
		if v == key {
			return true
		}
	}
	return false
}

// DoWithFilters performs request passed and filters entities in json response
func (c *Client) DoWithFilters(ctx context.Context, req *http.Request, v interface{}, filters []*AdditionalFilter, baseSearchPaths []string) error {
	// check if client exists or not
	if c.client == nil {
		return fmt.Errorf("%s", c.ErrorMsg)
	}
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err = CheckResponse(resp); err != nil {
		return err
	}

	resp.Body, err = filter(resp.Body, filters, baseSearchPaths)
	if err != nil {
		return err
	}

	if err := decodeResponse(resp, v); err != nil {
		return err
	}

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp, v)
	}

	return nil
}

func filter(body io.ReadCloser, filters []*AdditionalFilter, baseSearchPaths []string) (io.ReadCloser, error) {
	if filters == nil {
		return body, nil
	}
	if len(baseSearchPaths) == 0 {
		baseSearchPaths = []string{"$."}
	}

	var res map[string]interface{}
	b, err := io.ReadAll(body)
	if err != nil {
		return body, err
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return body, fmt.Errorf("error unmarshalling response for filtering: %w", err)
	}

	// Nothing to filter if the payload has no entities collection (error pages, empty
	// results). Returning the body unchanged lets the caller decode and report properly.
	entities, ok := res["entities"].([]interface{})
	if !ok {
		return io.NopCloser(bytes.NewReader(b)), nil
	}

	// Full search paths
	searchPaths := map[string][]string{}
	filterMap := map[string]*AdditionalFilter{}
	for _, filter := range filters {
		filterMap[filter.Name] = filter
		//Build search paths by appending target search paths to base paths
		filterSearchPaths := []string{}
		for _, baseSearchPath := range baseSearchPaths {
			searchPath := fmt.Sprintf("%s.%s", baseSearchPath, filter.Name)
			filterSearchPaths = append(filterSearchPaths, searchPath)
		}
		searchPaths[filter.Name] = filterSearchPaths
	}

	// Entities that pass filters
	var filteredEntities []interface{}

	for _, entity := range entities {
		searchTarget, ok := entity.(map[string]interface{})
		if !ok {
			continue
		}
		filtersPassed := 0
	filter_loop:
		for filter, filterSearchPaths := range searchPaths {
			for _, searchPath := range filterSearchPaths {
				val, err := jsonpath.Get(searchPath, searchTarget)
				if err != nil {
					continue
				}
				// Stringify leaf value since we support only string values in filter
				value := fmt.Sprint(val)
				if searchSlice(filterMap[filter].Values, value) {
					filtersPassed++
					continue filter_loop
				}
			}
		}

		// Value must pass all filters since we perform logical AND b/w filters
		if filtersPassed == len(filters) {
			filteredEntities = append(filteredEntities, entity)
		}
	}
	res["entities"] = filteredEntities

	// Convert filtered result back to io.ReadCloser
	filteredBody, jsonErr := json.Marshal(res)
	if jsonErr != nil {
		return body, jsonErr
	}
	return io.NopCloser(bytes.NewReader(filteredBody)), nil
}

// ErrNotFound reports that the API answered 404. Callers use errors.Is to tell
// "this resource is gone" apart from "this request failed", which decides whether
// Terraform drops the resource from state or surfaces an error to the operator.
var ErrNotFound = errors.New("resource not found")

// CheckResponse checks errors if exist errors in request
func CheckResponse(r *http.Response) error {
	c := r.StatusCode

	if c >= 200 && c <= 299 {
		return nil
	}

	// Nutanix returns non-json response with code 401 when
	// invalid credentials are used
	if c == http.StatusUnauthorized {
		return fmt.Errorf("invalid auth Credentials")
	}

	if c == http.StatusBadRequest {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("bad Request: failed to read body: %w", err)
		}
		return fmt.Errorf("bad Request: %s", string(bodyBytes))
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	rdr2 := io.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr2

	if c == http.StatusNotFound {
		if len(buf) == 0 {
			return ErrNotFound
		}
		return fmt.Errorf("%w: %s", ErrNotFound, string(buf))
	}

	// if has entities -> return nil
	// if has message_list -> check_error["state"]
	// if has status -> check_error["status.state"]
	if len(buf) == 0 {
		// A non-2xx with no body still failed. Returning nil here used to let callers
		// carry on with a zero-value response and nil-dereference it.
		return fmt.Errorf("request failed with status %d", c)
	}

	var res map[string]interface{}
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return fmt.Errorf("unmarshalling error response %s for response body %s", err, string(buf))
	}

	errRes := &ErrorResponse{}
	if status, ok := res["status"]; ok {
		// A string `status` is a successful non-error payload for some endpoints.
		if _, sok := status.(string); sok {
			return nil
		}

		// Any other scalar carries no structured error detail. Skip the struct fill
		// rather than asserting blindly, which used to panic on numeric/bool/array
		// values inside the error handler itself.
		if statusMap, mok := status.(map[string]interface{}); mok {
			err = fillStruct(statusMap, errRes)
		}
	} else if _, ok := res["state"]; ok {
		err = fillStruct(res, errRes)
	} else if _, ok := res["entities"]; ok {
		return nil
	}

	if err != nil {
		return err
	}

	// karbon error check
	if messageInfo, ok := res["message_info"]; ok {
		return fmt.Errorf("error: %s", messageInfo)
	}

	// This check is also used for some foundation api errors
	if message, ok := res["message"]; ok {
		log.Print(message)
		return fmt.Errorf("error: %s", message)
	}
	if errRes.State != "ERROR" {
		return nil
	}

	pretty, _ := json.MarshalIndent(errRes, "", "  ")
	return fmt.Errorf("error: %s", string(pretty))
}

// ErrorResponse ...
type ErrorResponse struct {
	APIVersion  string            `json:"api_version,omitempty"`
	Code        int64             `json:"code,omitempty"`
	Kind        string            `json:"kind,omitempty"`
	MessageList []MessageResource `json:"message_list"`
	State       string            `json:"state"`
}

// MessageResource ...
type MessageResource struct {
	// Custom key-value details relevant to the status.
	Details map[string]interface{} `json:"details,omitempty"`

	// If state is ERROR, a message describing the error.
	Message string `json:"message"`

	// If state is ERROR, a machine-readable snake-cased *string.
	Reason string `json:"reason"`
}

func (r *ErrorResponse) Error() string {
	err := ""
	for key, value := range r.MessageList {
		err = fmt.Sprintf("%d: {message:%s, reason:%s }", key, value.Message, value.Reason)
	}

	return err
}

func fillStruct(data map[string]interface{}, result interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(j, result)
}
