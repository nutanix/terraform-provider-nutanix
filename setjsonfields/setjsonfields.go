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
func SetJSONFields(d *schema.ResourceData) st.JSONstruct {
	spec := d.Get("spec").(*schema.Set).List()[0].(map[string]interface{})         // spec
	metadata := d.Get("metadata").(*schema.Set).List()[0].(map[string]interface{}) // metadata

	JSON := st.JSONstruct{
		APIVersion: convertToString(d.Get("api_version")), // api_version
		Spec:       SetSpec(spec),
		Metadata:   SetMetadata(metadata),
	}
	return JSON
}

// SetMetadata sets metadata fields in json struct
func SetMetadata(s map[string]interface{}) *(st.MetaDataStruct) {

	var categories []*string
	if s["categories"] != nil {
		for i := 0; i < len(s["categories"].([]interface{})); i++ {
			str := s["categories"].([]interface{})[i].(string)
			categories = append(categories, &str)
		}
	}

	MetadataI := st.MetaDataStruct{
		LastUpdateTime: convertToString(s["last_update_time"]),
		Kind:           convertToString(s["kind"]),
		UUID:           convertToString(s["uuid"]),
		CreationTime:   convertToString(s["creation_time"]),
		Name:           convertToString(s["name"]),
		SpecVersion:    convertToInt(s["spec_version"]),
		EntityVersion:  convertToInt(s["entity_version"]),
		OwnerReference: SetSubnetReference(s["owner_reference"].(*schema.Set).List()),
		Categories:     categories,
	}
	return &MetadataI
}

// SetSubnetReference sets owner_reference fields in json struct
func SetSubnetReference(t []interface{}) *(st.SubnetReferenceStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		SubnetReferenceI := st.SubnetReferenceStruct{
			Kind: convertToString(s["kind"]),
			UUID: convertToString(s["uuid"]),
			Name: convertToString(s["name"]),
		}
		return &SubnetReferenceI
	}
	return nil
}

// SetSpec sets spec fields in json struct
func SetSpec(s map[string]interface{}) *(st.SpecStruct) {

	resources := s["resources"].(*schema.Set).List()[0].(map[string]interface{}) // resources

	SpecI := st.SpecStruct{
		Resources: SetResources(resources),
	}
	return &SpecI
}

// SetResources sets resources fields in json struct
func SetResources(s map[string]interface{}) *(st.ResourcesStruct) {

	ResourcesI := st.ResourcesStruct{
		NumVCPUsPerSocket:     convertToInt(s["num_vcpus_per_socket"]),       // num_vcpus_per_socket
		NumSockets:            convertToInt(s["num_sockets"]),                // num_sockets
		MemorySizeMb:          convertToInt(s["memory_size_mb"]),             // memory_size_mb
		PowerState:            convertToString(s["power_state"]),             // power_state
		GuestOSID:             convertToString(s["guest_os_id"]),             // guest_os_id
		HardwareClockTimezone: convertToString(s["hardware_clock_timezone"]), // hardware_clock_timezone
	}
	return &ResourcesI
}

func main() {
}
