package vmmv2

import (
	"testing"

	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// TestFlattenOneOfOvaSource_UrlSource is the regression test for a guaranteed panic.
// OneOfOvaSource only ever holds OvaUrlSource, OvaVmSource or ObjectsLiteSource, but
// the flattener asserted the similarly-named UrlSource/VmDiskSource types. Those exist
// in the same package, so the code compiled and panicked 100% of the time it ran.
func TestFlattenOneOfOvaSource_UrlSource(t *testing.T) {
	source := import1.NewOneOfOvaSource()
	urlSrc := import1.OvaUrlSource{
		ObjectType_: utils.StringPtr("vmm.v4.content.OvaUrlSource"),
		Url:         utils.StringPtr("https://example.com/image.ova"),
	}
	if err := source.SetValue(urlSrc); err != nil {
		t.Fatalf("SetValue: %v", err)
	}

	got := flattenOneOfOvaSource(source)
	if len(got) != 1 {
		t.Fatalf("flattenOneOfOvaSource returned %d entries, want 1", len(got))
	}
	inner, ok := got[0]["ova_url_source"].([]map[string]interface{})
	if !ok || len(inner) != 1 {
		t.Fatalf("ova_url_source = %#v", got[0]["ova_url_source"])
	}
	if utils.StringValue(inner[0]["url"].(*string)) != "https://example.com/image.ova" {
		t.Errorf("url = %v", inner[0]["url"])
	}
}

// TestFlattenOneOfOvaSource_VmSource also covers the field rename: the old code read
// ExtId, but OvaVmSource exposes VmExtId, and the schema key is vm_ext_id. Both
// mismatches prove the branch had never executed.
func TestFlattenOneOfOvaSource_VmSource(t *testing.T) {
	source := import1.NewOneOfOvaSource()
	diskFormat := import1.OVADISKFORMAT_QCOW2
	vmSrc := import1.OvaVmSource{
		ObjectType_:    utils.StringPtr("vmm.v4.content.OvaVmSource"),
		VmExtId:        utils.StringPtr("vm-uuid-1"),
		DiskFileFormat: &diskFormat,
	}
	if err := source.SetValue(vmSrc); err != nil {
		t.Fatalf("SetValue: %v", err)
	}

	got := flattenOneOfOvaSource(source)
	if len(got) != 1 {
		t.Fatalf("flattenOneOfOvaSource returned %d entries, want 1", len(got))
	}
	inner, ok := got[0]["ova_vm_source"].([]map[string]interface{})
	if !ok || len(inner) != 1 {
		t.Fatalf("ova_vm_source = %#v", got[0]["ova_vm_source"])
	}
	if utils.StringValue(inner[0]["vm_ext_id"].(*string)) != "vm-uuid-1" {
		t.Errorf("vm_ext_id = %v, want vm-uuid-1", inner[0]["vm_ext_id"])
	}
	if inner[0]["disk_file_format"] != "QCOW2" {
		t.Errorf("disk_file_format = %v, want QCOW2", inner[0]["disk_file_format"])
	}
}

func TestFlattenOneOfOvaSource_ObjectsLiteSource(t *testing.T) {
	source := import1.NewOneOfOvaSource()
	objSrc := import1.ObjectsLiteSource{
		ObjectType_: utils.StringPtr("vmm.v4.content.ObjectsLiteSource"),
		Key:         utils.StringPtr("bucket/key"),
	}
	if err := source.SetValue(objSrc); err != nil {
		t.Fatalf("SetValue: %v", err)
	}

	got := flattenOneOfOvaSource(source)
	if len(got) != 1 {
		t.Fatalf("flattenOneOfOvaSource returned %d entries, want 1", len(got))
	}
	inner, ok := got[0]["object_lite_source"].([]map[string]interface{})
	if !ok || len(inner) != 1 {
		t.Fatalf("object_lite_source = %#v", got[0]["object_lite_source"])
	}
	if utils.StringValue(inner[0]["key"].(*string)) != "bucket/key" {
		t.Errorf("key = %v, want bucket/key", inner[0]["key"])
	}
}

func TestFlattenOneOfOvaSource_Nil(t *testing.T) {
	if got := flattenOneOfOvaSource(nil); got != nil {
		t.Errorf("flattenOneOfOvaSource(nil) = %#v, want nil", got)
	}
}

// TestExpandOneOfOvaSource_ReturnsErrorNotExit covers the replacement of log.Fatalf,
// which called os.Exit and killed the plugin process mid-apply.
func TestExpandOneOfOvaSource_ReturnsErrorNotExit(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"ova_url_source": []interface{}{
				map[string]interface{}{
					"url":        "https://example.com/image.ova",
					"basic_auth": []interface{}{},
				},
			},
		},
	}

	got, err := expandOneOfOvaSource(input)
	if err != nil {
		t.Fatalf("expandOneOfOvaSource returned error: %v", err)
	}
	if got == nil {
		t.Fatal("expandOneOfOvaSource returned nil for a valid url source")
	}
	if utils.StringValue(got.ObjectType_) != "vmm.v4.content.OvaUrlSource" {
		t.Errorf("ObjectType_ = %v", utils.StringValue(got.ObjectType_))
	}
}

func TestExpandOneOfOvaSource_Empty(t *testing.T) {
	got, err := expandOneOfOvaSource([]interface{}{})
	if err != nil {
		t.Fatalf("expandOneOfOvaSource returned error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for empty input, got %#v", got)
	}
}

// TestExpandOneOfOvaSource_RoundTrip proves expand and flatten agree, which is what
// keeps Terraform from seeing a permanent diff on the source block.
func TestExpandOneOfOvaSource_RoundTrip(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"ova_vm_source": []interface{}{
				map[string]interface{}{
					"vm_ext_id":        "vm-uuid-1",
					"disk_file_format": "QCOW2",
				},
			},
		},
	}

	expanded, err := expandOneOfOvaSource(input)
	if err != nil {
		t.Fatalf("expandOneOfOvaSource returned error: %v", err)
	}

	flattened := flattenOneOfOvaSource(expanded)
	inner, ok := flattened[0]["ova_vm_source"].([]map[string]interface{})
	if !ok || len(inner) != 1 {
		t.Fatalf("round trip lost ova_vm_source: %#v", flattened)
	}
	if utils.StringValue(inner[0]["vm_ext_id"].(*string)) != "vm-uuid-1" {
		t.Errorf("vm_ext_id did not survive the round trip: %v", inner[0]["vm_ext_id"])
	}
	if inner[0]["disk_file_format"] != "QCOW2" {
		t.Errorf("disk_file_format did not survive the round trip: %v", inner[0]["disk_file_format"])
	}
}
