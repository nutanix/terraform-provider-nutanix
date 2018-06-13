package v3

import (
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

func TestNewV3Client(t *testing.T) {
	type args struct {
		credentials client.Credentials
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewV3Client(tt.args.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewV3Client() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewV3Client() = %v, want %v", got, tt.want)
			}
		})
	}
}
