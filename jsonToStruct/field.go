package main

import (
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

const configHeader = `
package virtualmachineconfig

import (
	"github.com/hashicorp/terraform/helper/schema"
	nutanixV3 "nutanixV3"
	"time"
	"strings"
)

func convertToBool(a interface{}) bool {
	if a != nil {
		return a.(bool)
	}
	return false
}

func convertToInt(a interface{}) int64 {
	if a != nil {
		i :=  a.(int)
		return int64(i)
	}
	return 0
}

func convertToString(a interface{}) string {
	if a != nil {
		return a.(string)
	}
	return ""
}

// SetMachineConfig function sets fields in struct from ResourceData
func SetMachineConfig(d *schema.ResourceData) nutanixV3.VmIntentInput {
	spec := d.Get("spec").(*schema.Set).List()			 // spec
	metadata := d.Get("metadata").(*schema.Set).List()                     // metadata
	machine := nutanixV3.VmIntentInput{
		ApiVersion: "3.0", // api_version
		Spec:       SetSpec(spec, 0), 	//Spec
		Metadata:   SetMetadata(metadata, 0), 	//Metadata
	}
	if strings.ToUpper(machine.Spec.Resources.PowerState) == "ON" {
		machine.Spec.Resources.PowerState = "POWERED_ON"
	} else {
		machine.Spec.Resources.PowerState = "POWERED_OFF"
	}
	machine.Metadata.Kind = "vm"
	machine.Spec.Name = d.Get("name").(string)
	machine.Metadata.Name = d.Get("name").(string)
	return machine
}
`

func init() {

		fileConfig, err := os.Create(os.ExpandEnv("$GOPATH/src/github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachineconfig/virtualmachineconfig.go"))
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
			fmt.Fprintf(wConfig, "\n\n// Set%s sets %s fields in  json struct\n", gtype, name)
			fmt.Fprintf(wConfig, "func Set%s(t []interface{}, i int) nutanixV3.%s {\n\tif len(t) > 0 {\n", gtype, goStruct(name))
			fmt.Fprintf(wConfig, "\t\ts := t[i].(map[string]interface{})\n%s\n\t\t%s := nutanixV3.%s{\n", bodyList, gtype, goStruct(name))
			fmt.Fprintf(wConfig, "%s\t\t}\n\t\treturn %s\n\t}\n\treturn nutanixV3.%s{}\n}",bodyConfig, gtype, goStruct(name))
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
