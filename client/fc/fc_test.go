package foundationcentral

import (
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

func TestNewFoundationCentralClient(t *testing.T) {
	cred := client.Credentials{URL: "foo.com", Username: "username", Password: "password", Port: "", Endpoint: "0.0.0.0", Insecure: true}
	c, _ := NewFoundationCentralClient(cred)

	cred2 := client.Credentials{URL: "^^^", Username: "username", Password: "password", Port: "", Endpoint: "0.0.0.0", Insecure: true}
	c2, _ := NewFoundationCentralClient(cred2)

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
			got, err := NewFoundationCentralClient(tt.args.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFoundationCentralClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFoundationCentralClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
