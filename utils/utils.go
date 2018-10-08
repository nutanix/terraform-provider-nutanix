package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// PrintToJSON method helper to debug responses
func PrintToJSON(v interface{}, msg string) {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	log.Print("\n", msg, string(pretty))
	fmt.Print("\n", msg, string(pretty))
}

// DebugRequest ...
func DebugRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string("####################"))
	fmt.Println(string("###### REQUEST #######"))
	fmt.Println(string(requestDump))
}

// DebugResponse ...
func DebugResponse(req *http.Response) {
	requestDump, err := httputil.DumpResponse(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string("####################"))
	fmt.Println(string("###### RESPONSE #######"))
	fmt.Println(string(requestDump))
}
