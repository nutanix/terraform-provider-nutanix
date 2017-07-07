package main

import (
	"regexp"
	"fmt"
	"strings"
	"unicode"
	"os"
	"bufio"
	glog "log"
)	

var structGenerated = map[string]bool{}

// Field data type
type Field struct {
	name string
	gtype string
}

var structNameMap =  map[string]string {
	"cluster_reference": "ClusterReference",
	"data_source_reference":	"Reference",
	"disk_address":	"DiskAddress",
	"device_properties":	"VmDiskDeviceProperties",
	"spec":	"Vm",
	"network_function_chain_reference":	"NetworkFunctionChainReference",
	"subnet_reference":	"SubnetReference",
	"ip_endpoint_list":	"IpAddress",
	"parent_reference":	"Reference",
	"guest_customization":	"GuestCustomization",
	"cloud_init":	"GuestCustomizationCloudInit",
	"sysprep":	"GuestCustomizationSysprep",
	"owner_reference":	"UserReference",
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
func NewField(name, gtype string, bodyConfig []byte, bodyList  []byte,body ...byte) Field {
	fileConfig, err := os.OpenFile(os.ExpandEnv(configFilePath), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		glog.Fatal(err)
	}
	wConfig := bufio.NewWriter(fileConfig)
	defer fileConfig.Close()
	defer wConfig.Flush()
	if gtype == "struct" && len(body) > 0 {
		gtype = goField(name)
		if !structGenerated[gtype] {
			fmt.Fprintf(wConfig, configStruct, goFunc(name), name, goFunc(name), goStruct(name),  bodyList, goFunc(name), goStruct(name), bodyConfig, goFunc(name), goStruct(name))
			structGenerated[gtype] = true
		}	
	} else if gtype == "struct" {
		gtype = "map[string]string"
		if !structGenerated[goField(name)] {
			fmt.Fprintf(wConfig, configMap, goFunc(name), name, goFunc(name), goField(name), name, goField(name), name, goField(name), goField(name), goField(name), goField(name))
			structGenerated[goField(name)] = true
		}	
	}
	return Field{goField(name), gtype}
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
//Output: FooId
func goField(jsonfield string) string {
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

// Returns struct name for the json TypeSet
func goStruct(jsonfield string) string{
	structField := goField(jsonfield)
	structName := strings.TrimSuffix(structField, "List")
	structName = strings.TrimPrefix(structName, "Vm")
	structName = "Vm" + structName
	if structNameMap[jsonfield] != "" {
		structName = structNameMap[jsonfield]
	}
	return structName
}

// Returns name of the setconfig function for the corresponding struct
func goFunc(jsonfield string) string{
	structField := goField(jsonfield)
	return keywordsToUpper(structField, "Ip", "Uuid", "Vm", "Cpu", "Api")
}

func keywordsToUpper(src string, keywords ...string) string {
	var re = regexp.MustCompile(`(` + strings.Join(keywords, "|") + `)`)
	return re.ReplaceAllStringFunc(src, func(w string) string {
		return strings.ToUpper(w)
	})
}
						
