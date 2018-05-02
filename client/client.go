package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	libraryVersion = "v3"
	defaultBaseURL = "https://%s/"
	absolutePath   = "api/nutanix/" + libraryVersion
	userAgent      = "nutanix/" + libraryVersion
	mediaType      = "application/json"
)

//Client Config Configuration of the client
type Client struct {
	Credentials *Credentials

	// HTTP client used to communicate with the Nutanix API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string
}

// Credentials needed username and password
type Credentials struct {
	URL      string
	Username string
	Password string
	Endpoint string
	Port     string
	Insecure bool
}

// NewClient returns a new Nutanix API client.
func NewClient(credentials *Credentials) (*Client, error) {

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}

	httpClient := http.DefaultClient

	httpClient.Transport = transCfg

	baseURL, err := url.Parse(fmt.Sprintf(defaultBaseURL, credentials.URL))

	if err != nil {
		return nil, err
	}

	c := &Client{credentials, httpClient, baseURL, userAgent}

	return c, nil
}

// NewRequest creates a request
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, errp := url.Parse(absolutePath + urlStr)
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
	req.Header.Add("Authorization", "Basic "+
		base64.StdEncoding.EncodeToString([]byte(c.Credentials.Username+":"+c.Credentials.Password)))

	log.Printf("REQUEST BODY \n")
	utils.PrintToJSON(req, "REQUEST BODY")

	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("################")
	fmt.Println("REQUEST")
	fmt.Println(string(requestDump))

	return req, nil
}

//Do performs request passed
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println("################")
	fmt.Println("RESPONSE")

	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(responseDump))

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	err = CheckResponse(resp)

	if err != nil {
		return err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return err
			}

			log.Printf("RESPONSE BODY \n")
			utils.PrintToJSON(v, "RESPONSE BODY")
		}
	}

	return err
}

//CheckResponse checks errors if exist errors in request
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	data, err := ioutil.ReadAll(r.Body)
	res := &ErrorResponse{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return err
	}

	pretty, _ := json.MarshalIndent(res, "", "  ")
	return fmt.Errorf("Error: %s", string(pretty))
}

// ErrorResponse ...
type ErrorResponse struct {
	ApiVersion  string            `json:"api_version"`
	Code        int64             `json:"code"`
	Kind        string            `json:"kind"`
	MessageList []MessageResource `json:"message_list"`
	State       string            `json:"state"`
}

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
		err = fmt.Sprintf("%s: %s", key, value)
	}
	return err
}
