package nutanix

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func Test_getMetadataAttributes(t *testing.T) {
	type args struct {
		d        *schema.ResourceData
		metadata *v3.Metadata
		kind     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getMetadataAttributes(tt.args.d, tt.args.metadata, tt.args.kind); (err != nil) != tt.wantErr {
				t.Errorf("getMetadataAttributes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_readListMetadata(t *testing.T) {
	type args struct {
		d    *schema.ResourceData
		kind string
	}
	tests := []struct {
		name    string
		args    args
		want    *v3.DSMetadata
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readListMetadata(tt.args.d, tt.args.kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("readListMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readListMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setRSEntityMetadata(t *testing.T) {
	type args struct {
		v *v3.Metadata
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]interface{}
		want1 []map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := setRSEntityMetadata(tt.args.v)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setRSEntityMetadata() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("setRSEntityMetadata() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getReferenceValues(t *testing.T) {
	type args struct {
		r *v3.Reference
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getReferenceValues(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getReferenceValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClusterReferenceValues(t *testing.T) {
	type args struct {
		r *v3.Reference
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getClusterReferenceValues(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getClusterReferenceValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateRef(t *testing.T) {
	type args struct {
		ref map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want *v3.Reference
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateRef(tt.args.ref); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateMapStringValue(t *testing.T) {
	type args struct {
		value map[string]interface{}
		key   string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateMapStringValue(tt.args.value, tt.args.key); got != tt.want {
				t.Errorf("validateMapStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateMapIntValue(t *testing.T) {
	type args struct {
		value map[string]interface{}
		key   string
	}
	tests := []struct {
		name string
		args args
		want *int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateMapIntValue(tt.args.value, tt.args.key); got != tt.want {
				t.Errorf("validateMapIntValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
