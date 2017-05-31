package setjsonfields

import (
	//	"bufio"
	//	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	st "github.com/ideadevice/terraform-ahv-provider-plugin/jsonstruct"
	//	"os"
)

func convertToInt(a interface{}) int {
	if a != nil {
		return a.(int)
	}
	return 0
}

func convertToString(a interface{}) string {
	if a != nil {
		return a.(string)
	}
	return ""
}

// SetJSONFields function sets fields in struct from ResourceData
func SetJSONFields(JSON *(st.JSONstruct), d *schema.ResourceData) {

	JSON.APIVersion = convertToString(d.Get("api_version"))                // api_version
	spec := d.Get("spec").(*schema.Set).List()[0].(map[string]interface{}) // spec
	SetSpec(JSON, spec)
	metadata := d.Get("metadata").(*schema.Set).List()[0].(map[string]interface{}) // metadata
	SetMetadata(JSON, metadata)
}

// SetMetadata sets metadata fields in json struct
func SetMetadata(JSON *(st.JSONstruct), s map[string]interface{}) {
	JSON.Metadata.LastUpdateTime = convertToString(s["last_update_time"])
	JSON.Metadata.Kind = convertToString(s["kind"])
	JSON.Metadata.UUID = convertToString(s["uuid"])
	JSON.Metadata.CreationTime = convertToString(s["creation_time"])
	JSON.Metadata.Name = convertToString(s["name"])
	JSON.Metadata.SpecVersion = convertToInt(s["spec_version"])
	JSON.Metadata.EntityVersion = convertToInt(s["entity_version"])

	if s["owner_reference"] != nil {
		SetSubnetReference(JSON.Metadata.OwnerReference, s["owner_reference"].(*schema.Set).List()[0].(map[string]interface{}))
	}

	if s["categories"] != nil {
		for i := 0; i < len(s["categories"].([]interface{})); i++ {
			str := s["categories"].([]interface{})[i].(string)
			JSON.Metadata.Categories = append(JSON.Metadata.Categories, &str)
		}
	}

}

// SetSubnetReference sets owner_reference fields in json struct
func SetSubnetReference(a *(st.SubnetReferenceStruct), s map[string]interface{}) {
	a.Kind = convertToString(s["kind"])
	a.UUID = convertToString(s["uuid"])
	a.Name = convertToString(s["name"])
}

// SetSpec sets spec fields in json struct
func SetSpec(JSON *(st.JSONstruct), s map[string]interface{}) {

	resources := s["resources"].(*schema.Set).List()[0].(map[string]interface{}) // resources
	SetResources(JSON, resources)

}

// SetResources sets resources fields in json struct
func SetResources(JSON *(st.JSONstruct), s map[string]interface{}) {

	JSON.Spec.Resources.NumVCPUsPerSocket = convertToInt(s["num_vcpus_per_socket"])           // num_vcpus_per_socket
	JSON.Spec.Resources.NumSockets = convertToInt(s["num_sockets"])                           // num_sockets
	JSON.Spec.Resources.MemorySizeMb = convertToInt(s["memory_size_mb"])                      // memory_size_mb
	JSON.Spec.Resources.PowerState = convertToString(s["power_state"])                        // power_state
	JSON.Spec.Resources.GuestOSID = convertToString(s["guest_os_id"])                         // guest_os_id
	JSON.Spec.Resources.HardwareClockTimezone = convertToString(s["hardware_clock_timezone"]) // hardware_clock_timezone

}

func main() {
}
