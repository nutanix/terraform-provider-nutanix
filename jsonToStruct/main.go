package main


import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	glog "log"
	"os"
	"sort"
)

// TODO: write proper Usage and README
var (
	log               = glog.New(os.Stderr, "", glog.Lshortfile)
	fstruct           = flag.String("structName", "VmIntentInput", "struct name for json object")
	debug             = false
	ErrNotValidSyntax = errors.New("Json reflection is not valid Go syntax")
	fileSchema, _ 	  = os.Create(os.ExpandEnv("$GOPATH/src/github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachineschema/virtualmachineschema.go"))
	wSchema			  = bufio.NewWriter(fileSchema)
	depth 			  = 0
)

const schemaHeader = `
package virtualmachineschema

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// VMSchema is Schema for VM
func VMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_address": &schema.Schema{
        	Type:     schema.TypeString,
        	Computed: true,
        },
        "name": &schema.Schema{
            Type:     schema.TypeString,
            Required: true,
        },
`

func main() {
	flag.Parse()
	fmt.Fprintf(wSchema, "%s\n", schemaHeader)
	depth = 2
	err := read(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(wSchema, "\t}\n}")
	wSchema.Flush()
	fileSchema.Close()
}

func tabN(n int ){
	i := 0
	for i<n {
		i++
		fmt.Fprintf(wSchema, "\t")
	}
}

func read(r io.Reader, w io.Writer) error {
	var v interface{}
	err := json.NewDecoder(r).Decode(&v)
	if err != nil {
		log.Println(err)
		return err
	}
	buf := new(bytes.Buffer)
	// Open struct
	o, bConfig, bList, err := xreflect(v)
	if err != nil {
		log.Println(err)
		return err
	}
	field := NewField(*fstruct, "struct", bConfig, bList, o...)
	fmt.Fprintf(buf, "type %s %s", field.name, field.gtype)
	if debug {
		os.Stdout.WriteString("*********DEBUG***********")
		os.Stdout.Write(buf.Bytes())
		os.Stdout.WriteString("*********DEBUG***********")
	}
	// Pass through gofmt for uniform formatting, and weak syntax check.
	bConfig, err = format.Source(buf.Bytes())
	if err != nil {
		log.Println(err)
		fmt.Println("Final Go Code")
		fmt.Println()
		os.Stderr.Write(buf.Bytes())
		fmt.Println()
		return ErrNotValidSyntax
	}
	w.Write(bConfig)
	return nil
}

func xreflect(v interface{}) ([]byte, []byte, []byte, error) {
	var (
		buf = new(bytes.Buffer)
		bufConfig = new(bytes.Buffer)
		bufList = new(bytes.Buffer)
	)
	fields := []Field{}
	switch root := v.(type) {
	case map[string]interface{}:
		for key, val := range root {
			tabN(depth)
			fmt.Fprintf(wSchema, "\"%s\": &schema.Schema{\n", key)
			tabN(depth+1)
			fmt.Fprintf(wSchema, "Optional: true,\n")
			switch j := val.(type) {
			case nil:
				// FIXME: sometimes json service will return nil even though the type is string.
				// go can not convert string to nil and vs versa. Can we assume its a string?
				continue
			case float64:
				fields = append(fields, NewField(key, "int", nil, nil))
				tabN(depth+1)
				fmt.Fprintf(wSchema, "Type: schema.TypeInt,\n")
				fmt.Fprintf(bufConfig, "\t\t\t%s:\t\tconvertToInt(s[\"%s\"]),\n", goField(key), key)
				
			case map[string]interface{}:
				// If type is map[string]interface{} then we have nested object, Recurse
				tabN(depth+1)
				rootTemp := interface{}(j)
				s := fmt.Sprintf("%v", rootTemp)
				if s == "map[]" {
					fmt.Fprintf(wSchema, "Type: schema.TypeMap,\n")
					tabN(depth+1)
					fmt.Fprintf(wSchema, "Elem:     &schema.Schema{Type: schema.TypeString},\n")
					fmt.Fprintf(bufConfig, "\t\t\t%s:\t\tSet%s(s),\n", goField(key) ,goField(key))
				} else {
					fmt.Fprintf(wSchema, "Type: schema.TypeSet,\n")
					tabN(depth+1)
					fmt.Fprintf(wSchema, "Elem: &schema.Resource{\n")
					tabN(depth+2)
					fmt.Fprintf(wSchema, "Schema: map[string]*schema.Schema{\n")
				}	
				depth = depth + 3

				o, bConfig, bList, err := xreflect(j)
				if err != nil {
					log.Println(err)
					return nil, nil, nil, err
				}
				fields = append(fields, NewField(key, "struct", bConfig, bList, o...))
				depth = depth - 3 
				if s != "map[]" {
					tabN(depth+2)
					fmt.Fprintf(wSchema, "},\n")
					tabN(depth+1)
					fmt.Fprintf(wSchema, "},\n")
					fmt.Fprintf(bufConfig, "\t\t\t%s:\t\tSet%s(s[\"%s\"].(*schema.Set).List(), 0),\n", goField(key), goField(key), key)
				}	
			case []interface{}:
				tabN(depth+1)
				fmt.Fprintf(wSchema, "Type: schema.TypeList,\n")
				tabN(depth+1)
				fmt.Fprintf(wSchema, "Elem: &schema.Resource{\n")
				tabN(depth+2)
				fmt.Fprintf(wSchema, "Schema: map[string]*schema.Schema{\n")
				depth = depth + 3
				gtype, err := sliceType(key, j)
				if err != nil {
					log.Println(err)
					return nil, nil, nil, err
				}
				fields = append(fields, NewField(key, gtype, nil, nil))
				depth = depth - 3 
				tabN(depth+2)
				fmt.Fprintf(wSchema, "},\n")
				tabN(depth+1)
				fmt.Fprintf(wSchema, "},\n")
				fmt.Fprintf(bufConfig, "\t\t\t%s:\t\t%s,\n", goField(key), goField(key))
				fmt.Fprintf(bufList, "\n\t\tvar %s []nutanixV3.%s\n", goField(key), goStruct(key))
				fmt.Fprintf(bufList, "\t\tif s[\"%s\"] != nil {\n\t\t\tfor i := 0; i< len(s[\"%s\"].([]interface{})); i++ {\n", key, key)
				fmt.Fprintf(bufList, "\t\t\t\telem := Set%s(s[\"%s\"].([]interface{}),	i)\n", goField(key), key)
				fmt.Fprintf(bufList, "\t\t\t\t%s = append(%s, elem)\n\t\t\t}\n\t\t}\n\n", goField(key), goField(key))

			default:
				fields = append(fields, NewField(key, fmt.Sprintf("%T", val), nil, nil))
				tabN(depth+1)
				fmt.Fprintf(wSchema, "Type: schema.TypeString,\n")
				if strings.HasSuffix(goField(key), "Time") {
					fmt.Fprintf(bufConfig, "\t\t\t%s:\t\t%s,\n", goField(key), goField(key))
					fmt.Fprintf(bufList, "\n\t\tvar %s time.Time\n\t\ttemp%s := convertToString(s[\"%s\"])\n", goField(key), goField(key), key)
					fmt.Fprintf(bufList, "\t\tif temp%s != \"\"{\n\t\t\t%s, _ = time.Parse(temp%s,temp%s)\n", goField(key), goField(key), goField(key), goField(key))
					fmt.Fprintf(bufList, "\t\t}\n")
				} else {	
					fmt.Fprintf(bufConfig, "\t\t\t%s:\t\tconvertToString(s[\"%s\"]),\n", goField(key), key)
				}	
			}
			tabN(depth)
			fmt.Fprintf(wSchema, "},\n")
		}
	default:
		return nil, nil, nil, fmt.Errorf("%T: unexpected type", root)
	}
	// Sort and write field buffer last to keep order and formatting.
	sort.Sort(FieldSort(fields))
	for _, f := range fields {
		fmt.Fprintf(buf, "%s %s\n", f.name, f.gtype)
	}
	return buf.Bytes(), bufConfig.Bytes(), bufList.Bytes(), nil
}

// if all entries in j are the same type, return slice of that type
func sliceType(key string, j []interface{}) (string, error) {
	dft := "[]interface{}"
	if len(j) == 0 {
		return dft, nil
	}
	var t, t2 string
	for _, v := range j {
		switch v.(type) {
		case string:
			t2 = "[]string"
		case float64:
			t2 = "[]int"
		case map[string]interface{}:
			t2 = "[]struct"
		default:
			// something else, just return default
			return dft, nil
		}
		if t != "" && t != t2 {
			return dft, nil
		}
		t = t2
	}

	if t == "[]struct" {
		o, bConfig, bList, err := xreflect(j[0])
		if err != nil {
			log.Println(err)
			return "", err
		}
		f := NewField(key, "struct", bConfig, bList, o...)
		t = "[]" + f.gtype
	}
	return t, nil
}
