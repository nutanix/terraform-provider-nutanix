package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

// PrintToJSON method helper to debug responses
func PrintToJSON(v interface{}, msg string) {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	log.Print("\n", msg, string(pretty))
	fmt.Print("\n", msg, string(pretty))
}
