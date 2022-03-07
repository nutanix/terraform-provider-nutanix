package v3

import (
	"reflect"
	"testing"

	"github.com/nutanix/terraform-provider-nutanix/client"
)

func TestNewV3Client(t *testing.T) {
	cred := client.Credentials{URL: "foo.com", Username: "username", Password: "password", Port: "", Endpoint: "", Insecure: true}
	c, _ := NewV3Client(cred)

	cred2 := client.Credentials{URL: "^^^", Username: "username", Password: "password", Port: "", Endpoint: "", Insecure: true}
	c2, _ := NewV3Client(cred2)

	type args struct {
		credentials client.Credentials
	}

	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			"test one",
			args{cred},
			c,
			false,
		},
		{
			"test one",
			args{cred2},
			c2,
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
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
