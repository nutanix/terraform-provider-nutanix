package main

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
    spec := d.Get("spec").(*schema.Set).List()           // spec
    metadata := d.Get("metadata").(*schema.Set).List()                     // metadata
    machine := nutanixV3.VmIntentInput{
        ApiVersion: "3.0", // api_version
        Spec:       SetSpec(spec, 0),   //Spec
        Metadata:   SetMetadata(metadata, 0),   //Metadata
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

const configFilePath = "$GOPATH/src/github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachineconfig/virtualmachineconfig.go"

const configStruct = `

// Set%s sets %s fields in json struct
func Set%s (t []interface{}, i int) nutanixV3.%s {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		%s 

		%s := nutanixV3.%s{
			%s
		}
		return %s
	}
	return nutanixV3.%s{}
}
`
const configMap = `

// Set%s sets %s fields in json struct
func Set%s(s map[string]interface{}) map[string]string {
	var %sI map[string]interface{}
	if s["%s"] != nil {
		%sI = s["%s"].(map[string]interface{})
	}
	%s := make(map[string]string)
	for key, value := range %sI {
		switch value := value.(type){
		case string:
			%s[key]	= value
		}
	}
	return %s
}
`

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
const configList = `

		var  %s []nutanixV3.%s
		if s["%s"] != nil {
			for i := 0; i< len(s["%s"].([]interface{})); i++ {
				elem := Set%s(s["%s"].([]interface{}), i)
				%s = append(%s, elem)
			}
		}

`
const configTime = `

		var %s time.Time
		temp%s := convertToString(s["%s"])
		if temp%s != ""{
			%s, _ = time.Parse(temp%s, temp%s)
		}

`

const schemaFilePath = "$GOPATH/src/github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachineschema/virtualmachineschema.go"
