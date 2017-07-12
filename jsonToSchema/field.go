package main

import (
	"regexp"
	"fmt"
	"github.com/fatih/camelcase"
	"strings"
	"unicode"
	"os"
	"bufio"
	glog "log"
)	

var structGenerated = map[string]bool{}
var structNameMap =  map[string]string {
	"VmIntentInput": "VmIntentInput",
}

func init() {
		fileConfig, err := os.Create(os.ExpandEnv(configFilePath))
		if err != nil {
			glog.Fatal(err)
		}
		wConfig := bufio.NewWriter(fileConfig)
		defer fileConfig.Close()
		defer wConfig.Flush()
		fmt.Fprintf(wConfig, "%s\n", configHeader)
}

// NewField simplifies Field construction
func NewField(name, gtype string, bodyConfig []byte, bodyList  []byte){
	fileConfig, err := os.OpenFile(os.ExpandEnv(configFilePath), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		glog.Fatal(err)
	}
	wConfig := bufio.NewWriter(fileConfig)
	defer fileConfig.Close()
	defer wConfig.Flush()
	if gtype == "struct" {
		gtype = structNameMap[name]
		if !structGenerated[name] {
			fmt.Fprintf(wConfig, configStruct, goFunc(name), fromCamelcase(name), goFunc(name), structNameMap[name],  bodyList, goFunc(name), structNameMap[name], bodyConfig, goFunc(name), structNameMap[name])
			structGenerated[name] = true
		}	
	} else if gtype == "map[string]string" {
		if !structGenerated[name] {
			fmt.Fprintf(wConfig, configMap, goFunc(name), fromCamelcase(name), goFunc(name), name, fromCamelcase(name), name, fromCamelcase(name), name, name, name, name)
			structGenerated[name] = true
		}	
	}
}

// Returns lower_case json fields to camel case fields
// Example :
//		toCamelcase("foo_id")
//Output: FooId
func toCamelcase(jsonfield string) string {
	mkUpper := true
	structField := ""
	for _, c := range jsonfield {
		if mkUpper {
			c = unicode.ToUpper(c)
			mkUpper = false
		}
		if c == '_' {
			mkUpper = true
			continue
		}
		if c == '-' {
			mkUpper = true
			continue
		}
		structField += string(c)
	}
	return fmt.Sprintf("%s", structField)
}

//converts camelcase to delimiter-separeted words
func fromCamelcase(s string) string{
	split := camelcase.Split(s)
	name := ""
	for i := range split {
		name = name + strings.ToLower(split[i]) + "_"
	}
	name = strings.TrimSuffix(name, "_")
	return name
}

// Returns name of the setconfig function for the corresponding struct
func goFunc(jsonfield string) string{
	structField := toCamelcase(jsonfield)
	return keywordsToUpper(structField, "Ip", "Uuid", "Vm", "Cpu", "Api")
}

func keywordsToUpper(src string, keywords ...string) string {
	var re = regexp.MustCompile(`(` + strings.Join(keywords, "|") + `)`)
	return re.ReplaceAllStringFunc(src, func(w string) string {
		return strings.ToUpper(w)
	})
}
