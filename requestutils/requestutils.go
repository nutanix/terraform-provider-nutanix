package requestutils

import (
	"bufio"
	"bytes"
	"crypto/tls"
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

// RequestHandler  creates a connection request
func RequestHandler(url, method string, jsonStr []byte, username, password string) []byte {
	if method == "POST" {

		file, err := os.Create("request_log")
		check(err)
		defer file.Close()

		w := bufio.NewWriter(file)
		defer w.Flush()

		fmt.Fprintf(w, "url> %v\n", url)
		req, err1 := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		check(err1)
		req.SetBasicAuth(username, password)
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")
		fmt.Fprintf(w, "-------------- POST REQUEST ---------------------\n")
		fmt.Fprintf(w, "request Header: %v\n\n", req.Header)
		fmt.Fprintf(w, "request Body: %v\n\n", req.Body)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err2 := client.Do(req)
		check(err2)
		defer resp.Body.Close()

		fmt.Fprintf(w, "-------------- POST RESPONSE ---------------------\n")
		fmt.Fprintf(w, "response Status: %v\n\n", resp.Status)
		fmt.Fprintf(w, "response Headers: %v\n\n", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "response Body: %v\n\n", string(body))
		return body

	} else if method == "DELETE" {

		file, err := os.Create("request_log")
		check(err)
		defer file.Close()

		w := bufio.NewWriter(file)
		defer w.Flush()

		fmt.Fprintf(w, "url> %v\n", url)

		req, err1 := http.NewRequest("DELETE", url, nil)
		req.SetBasicAuth(username, password)
		check(err1)
		fmt.Fprintf(w, "-------------- DELETE REQUEST ---------------------\n")
		fmt.Fprintf(w, "request Header: %v\n\n", req.Header)
		fmt.Fprintf(w, "request Body: %v\n\n", req.Body)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}
		resp, err2 := client.Do(req)
		check(err2)
		defer resp.Body.Close()

		fmt.Fprintf(w, "-------------- DELETE RESPONSE ---------------------\n")
		fmt.Fprintf(w, "response Status: %v\n\n", resp.Status)
		fmt.Fprintf(w, "response Headers: %v\n\n", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "response Body: %v\n\n", string(body))
		return body

	} else if method == "GET" {

		file, err := os.Create("request_log")
		check(err)
		defer file.Close()

		w := bufio.NewWriter(file)
		defer w.Flush()

		fmt.Fprintf(w, "url> %v\n", url)

		req, err1 := http.NewRequest("GET", url, nil)
		req.SetBasicAuth(username, password)
		check(err1)
		fmt.Fprintf(w, "-------------- GET REQUEST ---------------------\n")
		fmt.Fprintf(w, "request Header: %v\n\n", req.Header)
		fmt.Fprintf(w, "request Body: %v\n\n", req.Body)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}
		resp, err2 := client.Do(req)
		check(err2)
		defer resp.Body.Close()

		fmt.Fprintf(w, "-------------- GET RESPONSE ---------------------\n")
		fmt.Fprintf(w, "response Status: %v\n\n", resp.Status)
		fmt.Fprintf(w, "response Headers: %v\n\n", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "response Body: %v\n\n", string(body))
		return body
	}
	return []byte(`{}`)
}
