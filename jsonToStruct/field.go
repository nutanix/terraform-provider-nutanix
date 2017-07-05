package main

import (
	"strconv"
	"strings"
	"fmt"
	"unicode"
	"os"
	"bufio"
	glog "log"
)	

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"GPU":	 true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"NTP":   true,
	"DB":    true,
}

var intToWordMap = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

var structGenerated = map[string]bool{}

// Field data type
type Field struct {
	name string
	gtype string
	tag string
}

func init() {

		fileConfig, err := os.Create(os.ExpandEnv("$GOPATH/src/github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachineconfig/virtualmachineconfig.go"))
		if err != nil {
			glog.Fatal(err)
		}
		wConfig := bufio.NewWriter(fileConfig)
		defer fileConfig.Close()
		defer wConfig.Flush()
		fmt.Fprintf(wConfig, "package virtualmachineconfig\n\nimport (\n\t\"github.com/hashicorp/terraform/helper/schema\"\n")
		fmt.Fprintf(wConfig, "\tvm \"github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachine\"\n)\n\n")
		fmt.Fprintf(wConfig, "func convertToBool(a interface{}) bool {\n\tif a != nil {\n\t\treturn a.(bool)\n\t}\n\treturn false\n}\n")
		fmt.Fprintf(wConfig, "func convertToInt(a interface{}) int {\n\tif a != nil {\n\t\treturn a.(int)\n\t}\n\treturn 0\n}\n")
		fmt.Fprintf(wConfig, "func convertToString(a interface{}) string {\n\tif a != nil {\n\t\treturn a.(string)\n\t}\n\treturn \"\"\n}\n")
}

// NewField simplifies Field construction
func NewField(name, gtype string, bodyConfig []byte, bodyList  []byte,body ...byte) Field {
	fileConfig, err := os.OpenFile(os.ExpandEnv("$GOPATH/src/github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachineconfig/virtualmachineconfig.go"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		glog.Fatal(err)
	}
	wConfig := bufio.NewWriter(fileConfig)
	defer fileConfig.Close()
	defer wConfig.Flush()
	if gtype == "struct" && len(body) > 0 {
		gtype = goField(name)
		if !structGenerated[gtype] {
			file, err := os.Create(os.ExpandEnv("$GOPATH/src/github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachine/") +  gtype + ".go")
			if err != nil {
				glog.Fatal(err)
			}
			w := bufio.NewWriter(file)
			defer file.Close()
			defer w.Flush()
			fmt.Fprintf(w, "package virtualmachine\n\n// %s struct\ntype %s struct {\n\n%s\n}", gtype, gtype, body)
			fmt.Fprintf(wConfig, "\n\n// Set%s sets %s fields in  json struct\n", gtype, name)
			fmt.Fprintf(wConfig, "func Set%s(t []interface{}, i int) vm.%s {\n\tif len(t) > 0 {\n", gtype, gtype)
			fmt.Fprintf(wConfig, "\t\ts := t[i].(map[string]interface{})\n%s\n\t\t%s := vm.%s{\n", bodyList, gtype, gtype)
			fmt.Fprintf(wConfig, "%s\t\t}\n\t\treturn %s\n\t}\n\treturn vm.%s{}\n}",bodyConfig, gtype, gtype)
			structGenerated[gtype] = true
		}	
	} else if gtype == "struct" {
		gtype = "map[string]string"
		if !structGenerated[goField(name)] {
			fmt.Fprintf(wConfig, "\n\n// Set%s sets %s fields in  json struct\n", goField(name), name)
			fmt.Fprintf(wConfig, "func Set%s(s map[string]interface{}) map[string]string {\n\tvar %sI map[string]interface{}\n\tif s[\"%s\"] != nil{\n", goField(name), goField(name), name)
			fmt.Fprintf(wConfig, "\t\t%sI = s[\"%s\"].(map[string]interface{})\n\t}\n\t%s := make(map[string]string)\n", goField(name), name, goField(name))
			fmt.Fprintf(wConfig, "\tfor key, value := range %sI {\n\t\t switch value := value.(type) {\n\t\tcase string:\n\t\t\t%s[key] = value\n\t\t}\n\t}\n", goField(name), goField(name))
			fmt.Fprintf(wConfig, "\treturn %s\n}\n", goField(name))
			structGenerated[goField(name)] = true
		}	
	}
	return Field{goField(name), gtype, goTag(name)}
}

// FieldSort Provides Sorter interface so we can keep field order
type FieldSort []Field

func (s FieldSort) Len() int { return len(s) }

func (s FieldSort) Swap(i, j int) { s[i], s[j] = s[j], s[i]}

func (s FieldSort) Less(i, j int) bool {
	return s[i].name < s[j].name
}


// Returns lower_case json fields to camel case fields
// Example :
//		goField("foo_id")
//Output: FooID
func goField(jsonfield string) string {
	runes := []rune(jsonfield)
	for len(runes) > 0 && !unicode.IsLetter(runes[0]) && !unicode.IsDigit(runes[0]) {
		runes = runes[1:]
	}
	if len(runes) == 0 {
		return "_"
	}
	
	s := stringifyFirstChar(string(runes))
	name := fieldName(s)
	runes = []rune(name)
	for i,c := range runes {
		ok := unicode.IsLetter(c) || unicode.IsDigit(c)
		if i == 0 {
			ok = unicode.IsLetter(c)
		}
		if !ok {
			runes[i] = '_'
		}
	}
	s = string(runes)
	s = strings.Trim(s, "_")
	if len(s) == 0 {
		return "_"
	}
	return s
}

func fieldName(name string) string {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}
	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		runes := []rune(name)
		if u := strings.ToUpper(name); commonInitialisms[u]{
			copy(runes[0:], []rune(u))
		} else {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return string(runes)
	}
	allUpperWithUnderscore := true
	for _, r := range name {
		if !unicode.IsUpper(r) && r != '_' {
			allUpperWithUnderscore = false
			break
		}
	}
	if allUpperWithUnderscore {
		name = strings.ToLower(name)
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false 

		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if teh underscore is betwee two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if  unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			eow = true
		}
		i++
		if !eow{
			continue
		}

		// [w, i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			copy(runes[w:], []rune(u))
		} else if strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}

// convert first character ints to strings
func stringifyFirstChar(str string) string {
	first := str[:1]

	i, err := strconv.ParseInt(first, 10, 8)

	if err != nil {
		return str
	}

	return intToWordMap[i] + "_" + str[1:]
}

// Returns the json tag from a json field.
func goTag(jsonfield string) string {
	return fmt.Sprintf("`json:\"%s,omitempty\"bson:\"%s,omitempty\"`", jsonfield, jsonfield)
}
