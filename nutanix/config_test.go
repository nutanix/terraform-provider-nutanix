package nutanix

import (
	"testing"
)

func TestConfig_Client(t *testing.T) {
	type fields struct {
		Endpoint           string
		Username           string
		Password           string
		Port               string
		Insecure           bool
		FoundationPort     string
		FoundationEndpoint string
	}
	config := &Config{
		Endpoint:           "http://localhost",
		Username:           "test",
		Password:           "test",
		Port:               "8080",
		Insecure:           true,
		FoundationPort:     "8000",
		FoundationEndpoint: "0.0.0.0",
	}

	client, err := config.Client()
	if err != nil {
		t.Errorf("failed to create wanted client: %s", err)
	}

	tests := []struct {
		name    string
		fields  fields
		want    *Client
		wantErr bool
	}{
		{
			name: "new client",
			fields: fields{
				Endpoint:           "http://localhost",
				Username:           "test",
				Password:           "test",
				Port:               "8080",
				Insecure:           true,
				FoundationPort:     "8000",
				FoundationEndpoint: "0.0.0.0",
			},
			want:    client,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Endpoint:           tt.fields.Endpoint,
				Username:           tt.fields.Username,
				Password:           tt.fields.Password,
				Port:               tt.fields.Port,
				Insecure:           tt.fields.Insecure,
				FoundationEndpoint: tt.fields.FoundationEndpoint,
				FoundationPort:     tt.fields.FoundationPort,
			}
			got, err := c.Client()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Client() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.want {
				t.Errorf("Config.Client() = %v, want %v", got, tt.want)
			}
		})
	}
}
