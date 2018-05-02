package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

// PrintToJSON method helper to debug responses
func PrintToJSON(v interface{}, msg string) {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	fmt.Print("\n\n[DEBUG] ", msg, string(pretty))
	log.Print("\n", msg, string(pretty))
}
