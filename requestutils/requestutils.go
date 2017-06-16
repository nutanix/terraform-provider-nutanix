package requestutils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type testStruct struct {
	Test string
}

// Function checks if there is an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

var statusCodeFilter map[int]bool

func init() {
	statusMap := map[int]bool{
		200: true,
		201: true,
		202: true,
		203: true,
		204: true,
		205: true,
		206: true,
		207: true,
		208: true,
	}
	statusCodeFilter = statusMap
}

// RequestHandler  creates a connection request
func RequestHandler(url, method string, jsonStr []byte, username, password string, insecure bool) ([]byte, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	check(err)
	req.SetBasicAuth(username, password)
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	requestBody := req.Body
	requestHeader := req.Header

	tr := &http.Transport{}
	if insecure {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if !statusCodeFilter[resp.StatusCode] {
		errorstr := fmt.Sprintf("jsonStr: %v \n %v URL: %v\n request Header: %v\n request Body: %v\n response Status: %v\n response Body: %v\n", string(jsonStr), method, url, requestHeader, requestBody, resp.Status, string(body))
		errormsg := errors.New(errorstr)
		return body, errormsg
	}
	return body, nil
}
