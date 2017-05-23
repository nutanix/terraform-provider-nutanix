package requestutils

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Function checks if there is an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// RequestHandler  creates a connection request
func RequestHandler(url, method string, jsonStr []byte) {
	if method == "POST" {

		file, err := os.Create("request_log")
		check(err)
		defer file.Close()

		w := bufio.NewWriter(file)
		defer w.Flush()

		fmt.Fprintf(w, "url> %v\n", url)

		req, err1 := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		check(err1)
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err2 := client.Do(req)
		check(err2)
		defer resp.Body.Close()

		fmt.Fprintf(w, "-------------- POST RESPONSE ---------------------\n")
		fmt.Fprintf(w, "response Status: %v\n", resp.Status)
		fmt.Fprintf(w, "response Headers: %v\n", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "response Body: %v\n", string(body))

	} else if method == "DELETE" {

		file, err := os.Create("request_log")
		check(err)
		defer file.Close()

		w := bufio.NewWriter(file)
		defer w.Flush()

		fmt.Fprintf(w, "url> %v\n", url)

		req, err1 := http.NewRequest("DELETE", url, nil)
		check(err1)

		client := &http.Client{}
		resp, err2 := client.Do(req)
		check(err2)
		defer resp.Body.Close()

		fmt.Fprintf(w, "-------------- DELETE RESPONSE ---------------------\n")
		fmt.Fprintf(w, "response Status: %v\n", resp.Status)
		fmt.Fprintf(w, "response Headers: %v\n", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "response Body: %v\n", string(body))

	} else if method == "GET" {

		file, err := os.Create("request_log")
		check(err)
		defer file.Close()

		w := bufio.NewWriter(file)
		defer w.Flush()

		fmt.Fprintf(w, "url> %v\n", url)

		req, err1 := http.NewRequest("GET", url, nil)
		check(err1)

		client := &http.Client{}
		resp, err2 := client.Do(req)
		check(err2)
		defer resp.Body.Close()

		fmt.Fprintf(w, "-------------- GET RESPONSE ---------------------\n")
		fmt.Fprintf(w, "response Status: %v\n", resp.Status)
		fmt.Fprintf(w, "response Headers: %v\n", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "response Body: %v\n", string(body))
	}

}
