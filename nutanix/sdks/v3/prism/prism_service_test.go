package prism

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func setup() (*http.ServeMux, *client.Client, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	c, _ := client.NewClient(&client.Credentials{
		URL:      "",
		Username: "username",
		Password: "password",
		Port:     "",
		Endpoint: "0.0.0.0",
		Insecure: true,
	},
		userAgent,
		absolutePath,
		false)
	c.BaseURL, _ = url.Parse(server.URL)

	return mux, c, server
}

func TestOperations_CreateVM(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/vms", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodPost; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"kind":                   "vm",
				"should_force_translate": false,
			},
			"spec": map[string]interface{}{
				"cluster_reference": map[string]interface{}{
					"kind": "cluster",
					"uuid": "00056024-6c13-4c74-0000-00000000ecb5",
				},
				"name": "VM123.create",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "vm",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
				"should_force_translate" : false
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		createRequest *VMIntentInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *VMIntentResponse
		wantErr bool
	}{
		{
			"Test CreateVM",
			fields{
				c,
			},
			args{
				&VMIntentInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Kind:                 utils.StringPtr("vm"),
						ShouldForceTranslate: utils.BoolPtr(false),
					},
					Spec: &VM{
						ClusterReference: &Reference{
							Kind: utils.StringPtr("cluster"),
							UUID: utils.StringPtr("00056024-6c13-4c74-0000-00000000ecb5"),
						},
						Name: utils.StringPtr("VM123.create"),
					},
				},
			},
			&VMIntentResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind:                 utils.StringPtr("vm"),
					UUID:                 utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					ShouldForceTranslate: utils.BoolPtr(false),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateVM(tt.args.createRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateVM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateVM() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteVM(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/vms/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)

		fmt.Fprintf(w, `{
				"status": {
					"state": "DELETE_PENDING",
					"execution_context": {
						"task_uuid": "ff1b9547-dc9a-4ebd-a2ff-f2b718af935e"
					}
				},
				"spec": "",
				"api_version": "3.1",
				"metadata": {
					"kind": "vm",
					"categories": {
						"Project": "default"
					}
				}
			}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteVM OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			false,
		},

		{
			"Test DeleteVM Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if _, err := op.DeleteVM(tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteVM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_GetVM(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/vms/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"vm","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	vmResponse := &VMIntentResponse{}
	vmResponse.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("vm"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *VMIntentResponse
		wantErr bool
	}{
		{
			"Test GetVM OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			vmResponse,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetVM(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetVM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListVM(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/vms/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"vm","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	vmList := &VMListIntentResponse{}
	vmList.Entities = make([]*VMIntentResource, 1)
	vmList.Entities[0] = &VMIntentResource{}
	vmList.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("vm"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *VMListIntentResponse
		wantErr bool
	}{
		{
			"Test ListVM OK",
			fields{c},
			args{input},
			vmList,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListVM(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListVM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UpdateVM(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/vms/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"kind": "vm",
			},
			"spec": map[string]interface{}{
				"cluster_reference": map[string]interface{}{
					"kind": "cluster",
					"uuid": "00056024-6c13-4c74-0000-00000000ecb5",
				},
				"name": "VM123.create",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "vm",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
		body *VMIntentInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *VMIntentResponse
		wantErr bool
	}{
		{
			"Test UpdateVM",
			fields{
				c,
			},
			args{
				"cfde831a-4e87-4a75-960f-89b0148aa2cc",
				&VMIntentInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Kind: utils.StringPtr("vm"),
					},
					Spec: &VM{
						ClusterReference: &Reference{
							Kind: utils.StringPtr("cluster"),
							UUID: utils.StringPtr("00056024-6c13-4c74-0000-00000000ecb5"),
						},
						Name: utils.StringPtr("VM123.create"),
					},
				},
			},
			&VMIntentResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind: utils.StringPtr("vm"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.UpdateVM(tt.args.UUID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.UpdateVM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.UpdateVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_CreateSubnet(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/subnets", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"kind": "subnet",
			},
			"spec": map[string]interface{}{
				"cluster_reference": map[string]interface{}{
					"kind": "cluster",
					"uuid": "00056024-6c13-4c74-0000-00000000ecb5",
				},
				"name": "subnet.create",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "subnet",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		createRequest *SubnetIntentInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SubnetIntentResponse
		wantErr bool
	}{
		{
			"Test CreateSubnet",
			fields{
				c,
			},
			args{
				&SubnetIntentInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Kind: utils.StringPtr("subnet"),
					},
					Spec: &Subnet{
						ClusterReference: &Reference{
							Kind: utils.StringPtr("cluster"),
							UUID: utils.StringPtr("00056024-6c13-4c74-0000-00000000ecb5"),
						},
						Name: utils.StringPtr("subnet.create"),
					},
				},
			},
			&SubnetIntentResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind: utils.StringPtr("subnet"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateSubnet(tt.args.createRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateSubnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteSubnet(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/subnets/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)

		fmt.Fprintf(w, `{
				"status": {
					"state": "DELETE_PENDING",
					"execution_context": {
						"task_uuid": "ff1b9547-dc9a-4ebd-a2ff-f2b718af935e"
					}
				},
				"spec": "",
				"api_version": "3.1",
				"metadata": {
					"kind": "subnet",
					"categories": {
						"Project": "default"
					}
				}
			}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteSubnet OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			false,
		},

		{
			"Test DeleteSubnet Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if _, err := op.DeleteSubnet(tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteSubnet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_GetSubnet(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/subnets/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"subnet","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	subnetResponse := &SubnetIntentResponse{}
	subnetResponse.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("subnet"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SubnetIntentResponse
		wantErr bool
	}{
		{
			"Test GetSubnet OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			subnetResponse,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetSubnet(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetSubnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListSubnet(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/subnets/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"subnet","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	subnetList := &SubnetListIntentResponse{}
	subnetList.Entities = make([]*SubnetIntentResponse, 1)
	subnetList.Entities[0] = &SubnetIntentResponse{}
	subnetList.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("subnet"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SubnetListIntentResponse
		wantErr bool
	}{
		{
			"Test ListSubnet OK",
			fields{c},
			args{input},
			subnetList,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListSubnet(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListSubnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UpdateSubnet(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/subnets/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"kind": "subnet",
			},
			"spec": map[string]interface{}{
				"cluster_reference": map[string]interface{}{
					"kind": "cluster",
					"uuid": "00056024-6c13-4c74-0000-00000000ecb5",
				},
				"name": "subnet.create",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "subnet",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
		body *SubnetIntentInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SubnetIntentResponse
		wantErr bool
	}{
		{
			"Test UpdateSubnet",
			fields{
				c,
			},
			args{
				"cfde831a-4e87-4a75-960f-89b0148aa2cc",
				&SubnetIntentInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Kind: utils.StringPtr("subnet"),
					},
					Spec: &Subnet{
						ClusterReference: &Reference{
							Kind: utils.StringPtr("cluster"),
							UUID: utils.StringPtr("00056024-6c13-4c74-0000-00000000ecb5"),
						},
						Name: utils.StringPtr("subnet.create"),
					},
				},
			},
			&SubnetIntentResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind: utils.StringPtr("subnet"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.UpdateSubnet(tt.args.UUID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.UpdateSubnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.UpdateSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_CreateImage(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/images", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"kind": "image",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"image_type": "DISK_IMAGE",
					"data_source_reference": map[string]interface{}{
						"kind": "vm_disk",
						"uuid": "0005a238-f165-08ba-317e-ac1f6b6e5442",
					},
				},
				"name": "image.create",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "image",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		body *ImageIntentInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ImageIntentResponse
		wantErr bool
	}{
		{
			"Test CreateImage",
			fields{
				c,
			},
			args{
				&ImageIntentInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Kind: utils.StringPtr("image"),
					},
					Spec: &Image{
						Name: utils.StringPtr("image.create"),
						Resources: &ImageResources{
							ImageType: utils.StringPtr("DISK_IMAGE"),
							DataSourceReference: &Reference{
								Kind: utils.StringPtr("vm_disk"),
								UUID: utils.StringPtr("0005a238-f165-08ba-317e-ac1f6b6e5442"),
							},
						},
					},
				},
			},
			&ImageIntentResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind: utils.StringPtr("image"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateImage(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UploadImageError(t *testing.T) {
	_, c, server := setup()

	defer server.Close()

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID     string
		filepath string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test UploadImage ERROR (Cannot Open File)",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc", "xx"},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if err := op.UploadImage(tt.args.UUID, tt.args.filepath, nil); (err != nil) != tt.wantErr {
				t.Errorf("Operations.UploadImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_UploadImage(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/images/cfde831a-4e87-4a75-960f-89b0148aa2cc/file", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		bodyBytes, _ := ioutil.ReadAll(r.Body)
		file, _ := ioutil.ReadFile("prism.go")

		if !reflect.DeepEqual(bodyBytes, file) {
			t.Errorf("Operations.UploadImage() error: different uploaded files")
		}
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID     string
		filepath string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"TestOperations_UploadImage Upload Image",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc", "prism.go"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if err := op.UploadImage(tt.args.UUID, tt.args.filepath, nil); err != nil {
				t.Errorf("Operations.UploadImage() error = %v", err)
			}
		})
	}
}

func TestOperations_DeleteImage(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/images/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)

		fmt.Fprintf(w, `{
				"status": {
					"state": "DELETE_PENDING",
					"execution_context": {
						"task_uuid": "ff1b9547-dc9a-4ebd-a2ff-f2b718af935e"
					}
				},
				"spec": "",
				"api_version": "3.1",
				"metadata": {
					"kind": "image",
					"categories": {
						"Project": "default"
					}
				}
			}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteImage OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			false,
		},

		{
			"Test DeleteImage Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if _, err := op.DeleteImage(tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_GetImage(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/images/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"image","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	response := &ImageIntentResponse{}
	response.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("image"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ImageIntentResponse
		wantErr bool
	}{
		{
			"Test GetImage OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			response,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetImage(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListImage(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/images/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"image","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	list := &ImageListIntentResponse{}
	list.Entities = make([]*ImageIntentResponse, 1)
	list.Entities[0] = &ImageIntentResponse{}
	list.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("image"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ImageListIntentResponse
		wantErr bool
	}{
		{
			"Test ListImage OK",
			fields{c},
			args{input},
			list,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListImage(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UpdateImage(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/images/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"kind": "image",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"image_type": "DISK_IMAGE",
				},
				"name": "image.update",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "image",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
		body *ImageIntentInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ImageIntentResponse
		wantErr bool
	}{
		{
			"Test UpdateVM",
			fields{
				c,
			},
			args{
				"cfde831a-4e87-4a75-960f-89b0148aa2cc",
				&ImageIntentInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Kind: utils.StringPtr("image"),
					},
					Spec: &Image{
						Resources: &ImageResources{
							ImageType: utils.StringPtr("DISK_IMAGE"),
						},
						Name: utils.StringPtr("image.update"),
					},
				},
			},
			&ImageIntentResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind: utils.StringPtr("image"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.UpdateImage(tt.args.UUID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.UpdateImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.UpdateImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_GetCluster(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/clusters/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"cluster","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	response := &ClusterIntentResponse{}
	response.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("cluster"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ClusterIntentResponse
		wantErr bool
	}{
		{
			"Test GetCluster OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			response,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetCluster(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListCluster(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/clusters/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"cluster","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	list := &ClusterListIntentResponse{}
	list.Entities = make([]*ClusterIntentResponse, 1)
	list.Entities[0] = &ClusterIntentResponse{}
	list.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("cluster"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ClusterListIntentResponse
		wantErr bool
	}{
		{
			"Test ListCLusters OK",
			fields{c},
			args{input},
			list,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListCluster(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_CreateOrUpdateCategoryKey(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/categories/test_category_key", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"description": "Testing Keys",
			"name":        "test_category_key",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"description": "Testing Keys",
			"name": "test_category_key",
			"system_defined": false
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		body *CategoryKey
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CategoryKeyStatus
		wantErr bool
	}{
		{
			"Test CreateOrUpdateCaegoryKey OK",
			fields{c},
			args{&CategoryKey{
				Description: utils.StringPtr("Testing Keys"),
				Name:        utils.StringPtr("test_category_key"),
			}},
			&CategoryKeyStatus{
				Description:   utils.StringPtr("Testing Keys"),
				Name:          utils.StringPtr("test_category_key"),
				SystemDefined: utils.BoolPtr(false),
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateOrUpdateCategoryKey(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateOrUpdateCategoryKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateOrUpdateCategoryKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListCategories(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/categories/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{ "description": "Testing Keys", "name": "test_category_key", "system_defined": false }]}`)
	})

	list := &CategoryKeyListResponse{}
	list.Entities = make([]*CategoryKeyStatus, 1)
	list.Entities[0] = &CategoryKeyStatus{
		Description:   utils.StringPtr("Testing Keys"),
		Name:          utils.StringPtr("test_category_key"),
		SystemDefined: utils.BoolPtr(false),
	}

	input := &CategoryListMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *CategoryListMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CategoryKeyListResponse
		wantErr bool
	}{
		{
			"Test ListCategoryKey OK",
			fields{c},
			args{input},
			list,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListCategories(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListCategories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteCategoryKey(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/categories/test_category_key", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteSubnet OK",
			fields{c},
			args{"test_category_key"},
			false,
		},

		{
			"Test SubnetVM Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if err := op.DeleteCategoryKey(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteCategoryKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_GetCategoryKey(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/categories/test_category_key", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"description": "Testing Keys",
			"name": "test_category_key",
			"system_defined": false
		}`)
	})

	response := &CategoryKeyStatus{
		Description:   utils.StringPtr("Testing Keys"),
		Name:          utils.StringPtr("test_category_key"),
		SystemDefined: utils.BoolPtr(false),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CategoryKeyStatus
		wantErr bool
	}{
		{
			"Test GetCategory OK",
			fields{c},
			args{"test_category_key"},
			response,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetCategoryKey(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetCategoryKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetCategoryKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListCategoryValues(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/categories/test_category_key/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{ "description": "Testing Keys", "value": "test_category_value", "system_defined": false }]}`)
	})

	list := &CategoryValueListResponse{}
	list.Entities = make([]*CategoryValueStatus, 1)
	list.Entities[0] = &CategoryValueStatus{
		Description:   utils.StringPtr("Testing Keys"),
		Value:         utils.StringPtr("test_category_value"),
		SystemDefined: utils.BoolPtr(false),
	}

	input := &CategoryListMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		name               string
		getEntitiesRequest *CategoryListMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CategoryValueListResponse
		wantErr bool
	}{
		{
			"Test ListCategoryKey OK",
			fields{c},
			args{"test_category_key", input},
			list,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListCategoryValues(tt.args.name, tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListCategoryValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListCategoryValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_CreateOrUpdateCategoryValue(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/categories/test_category_key/test_category_value", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"description": "Testing Value",
			"value":       "test_category_value",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"description": "Testing Value",
			"name": "test_category_key",
			"value": "test_category_value",
			"system_defined": false
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		name string
		body *CategoryValue
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CategoryValueStatus
		wantErr bool
	}{
		{
			"Test CreateOrUpdateCategoryValue OK",
			fields{c},
			args{"test_category_key", &CategoryValue{
				Description: utils.StringPtr("Testing Value"),
				Value:       utils.StringPtr("test_category_value"),
			}},
			&CategoryValueStatus{
				Description:   utils.StringPtr("Testing Value"),
				Value:         utils.StringPtr("test_category_value"),
				Name:          utils.StringPtr("test_category_key"),
				SystemDefined: utils.BoolPtr(false),
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateOrUpdateCategoryValue(tt.args.name, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateOrUpdateCategoryValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateOrUpdateCategoryValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_GetCategoryValue(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/categories/test_category_key/test_category_value", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"description": "Testing Value",
			"name": "test_category_key",
			"value": "test_category_value",
			"system_defined": false
		}`)
	})

	response := &CategoryValueStatus{
		Description:   utils.StringPtr("Testing Value"),
		Name:          utils.StringPtr("test_category_key"),
		Value:         utils.StringPtr("test_category_value"),
		SystemDefined: utils.BoolPtr(false),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		name  string
		value string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CategoryValueStatus
		wantErr bool
	}{
		{
			"Test GetCategoryValue OK",
			fields{c},
			args{"test_category_key", "test_category_value"},
			response,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetCategoryValue(tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetCategoryValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetCategoryValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteCategoryValue(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/categories/test_category_key/test_category_value", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		name  string
		value string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteSubnet OK",
			fields{c},
			args{"test_category_key", "test_category_value"},
			false,
		},

		{
			"Test SubnetVM Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if err := op.DeleteCategoryValue(tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteCategoryValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_GetCategoryQuery(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/category/query", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"results":[{ "kind": "category_result" }]}`)
	})

	response := &CategoryQueryResponse{}
	response.Results = make([]*CategoryQueryResponseResults, 1)
	response.Results[0] = &CategoryQueryResponseResults{
		Kind: utils.StringPtr("category_result"),
	}

	input := &CategoryQueryInput{
		UsageType: utils.StringPtr("APPLIED_TO"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		query *CategoryQueryInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CategoryQueryResponse
		wantErr bool
	}{
		{
			"Test Category Query OK",
			fields{c},
			args{input},
			response,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetCategoryQuery(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetCategoryQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetCategoryQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_CreateNetworkSecurityRule(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/network_security_rules", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"kind": "network_security_rule",
			},
			"spec": map[string]interface{}{
				"description": "Network Create",
				"name":        "network.create",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "network_security_rule",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		request *NetworkSecurityRuleIntentInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *NetworkSecurityRuleIntentResponse
		wantErr bool
	}{
		{
			"Test CreateNetwork",
			fields{
				c,
			},
			args{
				&NetworkSecurityRuleIntentInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Kind: utils.StringPtr("network_security_rule"),
					},
					Spec: &NetworkSecurityRule{
						Name:        utils.StringPtr("network.create"),
						Description: utils.StringPtr("Network Create"),
						Resources:   nil,
					},
				},
			},
			&NetworkSecurityRuleIntentResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind: utils.StringPtr("network_security_rule"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateNetworkSecurityRule(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateNetworkSecurityRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateNetworkSecurityRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteNetworkSecurityRule(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/network_security_rules/cfde831a-4e87-4a75-960f-89b0148aa2cc",
		func(w http.ResponseWriter, r *http.Request) {
			testHTTPMethod(t, r, http.MethodDelete)

			fmt.Fprintf(w, `{
				"status": {
					"state": "DELETE_PENDING",
					"execution_context": {
						"task_uuid": "ff1b9547-dc9a-4ebd-a2ff-f2b718af935e"
					}
				},
				"spec": "",
				"api_version": "3.1",
				"metadata": {
					"kind": "network_security_rule",
					"categories": {
						"Project": "default"
					}
				}
			}`)
		})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteNetwork OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			false,
		},

		{
			"Test DeleteNetowork Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if _, err := op.DeleteNetworkSecurityRule(tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteNetworkSecurityRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_GetNetworkSecurityRule(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/network_security_rules/cfde831a-4e87-4a75-960f-89b0148aa2cc",
		func(w http.ResponseWriter, r *http.Request) {
			testHTTPMethod(t, r, http.MethodGet)
			fmt.Fprint(w, `{"metadata": {"kind":"network_security_rule","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
		})

	response := &NetworkSecurityRuleIntentResponse{}
	response.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("network_security_rule"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *NetworkSecurityRuleIntentResponse
		wantErr bool
	}{
		{
			"Test GetNetworkSecurityRule OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			response,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetNetworkSecurityRule(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetNetworkSecurityRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetNetworkSecurityRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListNetworkSecurityRule(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/network_security_rules/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"network_security_rule","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	list := &NetworkSecurityRuleListIntentResponse{}
	list.Entities = make([]*NetworkSecurityRuleIntentResource, 1)
	list.Entities[0] = &NetworkSecurityRuleIntentResource{}
	list.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("network_security_rule"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *NetworkSecurityRuleListIntentResponse
		wantErr bool
	}{
		{
			"Test ListNetwork OK",
			fields{c},
			args{input},
			list,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListNetworkSecurityRule(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListNetworkSecurityRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListNetworkSecurityRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UpdateNetworkSecurityRule(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/network_security_rules/cfde831a-4e87-4a75-960f-89b0148aa2cc",
		func(w http.ResponseWriter, r *http.Request) {
			testHTTPMethod(t, r, http.MethodPut)

			expected := map[string]interface{}{
				"api_version": "3.1",
				"metadata": map[string]interface{}{
					"kind": "network_security_rule",
				},
				"spec": map[string]interface{}{
					"description": "Network Update",
					"name":        "network.update",
				},
			}

			var v map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&v)
			if err != nil {
				t.Fatalf("decode json: %v", err)
			}

			if !reflect.DeepEqual(v, expected) {
				t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
			}

			fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "network_security_rule",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
		})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
		body *NetworkSecurityRuleIntentInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *NetworkSecurityRuleIntentResponse
		wantErr bool
	}{
		{
			"Test UpdateNetwork",
			fields{
				c,
			},
			args{
				"cfde831a-4e87-4a75-960f-89b0148aa2cc",
				&NetworkSecurityRuleIntentInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Kind: utils.StringPtr("network_security_rule"),
					},
					Spec: &NetworkSecurityRule{
						Resources:   nil,
						Description: utils.StringPtr("Network Update"),
						Name:        utils.StringPtr("network.update"),
					},
				},
			},
			&NetworkSecurityRuleIntentResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind: utils.StringPtr("network_security_rule"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.UpdateNetworkSecurityRule(tt.args.UUID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.UpdateNetworkSecurityRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.UpdateNetworkSecurityRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_CreateVolumeGroup(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/volume_groups", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"kind": "volume_group",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"flash_mode": "ON",
				},
				"name": "volume.create",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "volume_group",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		request *VolumeGroupInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *VolumeGroupResponse
		wantErr bool
	}{
		{
			"Test CreateVolumeGroup",
			fields{c},
			args{
				&VolumeGroupInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Kind: utils.StringPtr("volume_group"),
					},
					Spec: &VolumeGroup{
						Name: utils.StringPtr("volume.create"),
						Resources: &VolumeGroupResources{
							FlashMode: utils.StringPtr("ON"),
						},
					},
				},
			},
			&VolumeGroupResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind: utils.StringPtr("volume_group"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateVolumeGroup(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateVolumeGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateVolumeGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteVolumeGroup(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/volume_groups/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteVolumeGroup OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			false,
		},

		{
			"Test DeleteVolumeGroup Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if err := op.DeleteVolumeGroup(tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteVolumeGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_GetVolumeGroup(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/volume_groups/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"volume_group","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	response := &VolumeGroupResponse{}
	response.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("volume_group"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *VolumeGroupResponse
		wantErr bool
	}{
		{
			"Test GetVolumeGroup OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			response,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetVolumeGroup(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetVolumeGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetVolumeGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListVolumeGroup(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/volume_groups/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"volume_group","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	list := &VolumeGroupListResponse{}
	list.Entities = make([]*VolumeGroupResponse, 1)
	list.Entities[0] = &VolumeGroupResponse{}
	list.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("volume_group"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *VolumeGroupListResponse
		wantErr bool
	}{
		{
			"Test ListVolumeGroup OK",
			fields{c},
			args{input},
			list,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListVolumeGroup(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListVolumeGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListVolumeGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UpdateVolumeGroup(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/volume_groups/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"kind": "volume_group",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"flash_mode": "ON",
				},
				"name": "volume.update",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "volume_group",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
		body *VolumeGroupInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *VolumeGroupResponse
		wantErr bool
	}{
		{
			"Test UpdateVolumeGroup",
			fields{c},
			args{
				"cfde831a-4e87-4a75-960f-89b0148aa2cc",
				&VolumeGroupInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Kind: utils.StringPtr("volume_group"),
					},
					Spec: &VolumeGroup{
						Resources: &VolumeGroupResources{
							FlashMode: utils.StringPtr("ON"),
						},
						Name: utils.StringPtr("volume.update"),
					},
				},
			},
			&VolumeGroupResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind: utils.StringPtr("volume_group"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.UpdateVolumeGroup(tt.args.UUID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.UpdateVolumeGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.UpdateVolumeGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testHTTPMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func TestOperations_GetHost(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/hosts/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"host","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	hostResponse := &HostResponse{}
	hostResponse.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("host"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *HostResponse
		wantErr bool
	}{
		{
			"Test GetHost OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			hostResponse,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetHost(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListHost(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/hosts/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"host","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	hostList := &HostListResponse{}
	hostList.Entities = make([]*HostResponse, 1)
	hostList.Entities[0] = &HostResponse{}
	hostList.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("host"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *HostListResponse
		wantErr bool
	}{
		{
			"Test ListSubnet OK",
			fields{c},
			args{input},
			hostList,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListHost(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_CreateProject(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/projects", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"name": "project_test_name",
				"kind": "project",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"resource_domain": map[string]interface{}{
						"resources": []interface{}{
							map[string]interface{}{
								"limit":         float64(4),
								"resource_type": "resource_type_test",
							},
						},
					},
				},
				"name":        "project_name",
				"description": "description_test",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "project",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		request *Project
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Project
		wantErr bool
	}{
		{
			"Test CreateProject",
			fields{c},
			args{
				&Project{
					APIVersion: "3.1",
					Metadata: &Metadata{
						Name: utils.StringPtr("project_test_name"),
						Kind: utils.StringPtr("project"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &ProjectSpec{
						Name:       "project_name",
						Descripion: "description_test",
						Resources: &ProjectResources{
							ResourceDomain: &ResourceDomain{
								Resources: []*Resources{
									{
										Limit:        utils.Int64Ptr(4),
										ResourceType: "resource_type_test",
									},
								},
							},
						},
					},
				},
			},
			&Project{
				APIVersion: "3.1",
				Metadata: &Metadata{
					Kind: utils.StringPtr("project"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateProject(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_GetProject(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/projects/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"host","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	hostResponse := &Project{}
	hostResponse.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("host"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Project
		wantErr bool
	}{
		{
			"Test GetProject OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			hostResponse,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetProject(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListProject(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/projects/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"host","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	hostList := &ProjectListResponse{}
	hostList.Entities = make([]*Project, 1)
	hostList.Entities[0] = &Project{}
	hostList.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("host"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ProjectListResponse
		wantErr bool
	}{
		{
			"Test ListSubnet OK",
			fields{c},
			args{input},
			hostList,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListProject(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UpdateProject(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/projects/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "project_test_name",
				"kind": "project",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"resource_domain": map[string]interface{}{
						"resources": []interface{}{
							map[string]interface{}{
								"limit":         float64(4),
								"resource_type": "resource_type_test",
							},
						},
					},
				},
				"name":        "project_name",
				"description": "description_test",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "project",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
		body *Project
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Project
		wantErr bool
	}{
		{
			"Test CreateProject",
			fields{c},
			args{
				"cfde831a-4e87-4a75-960f-89b0148aa2cc",
				&Project{
					Metadata: &Metadata{
						Name: utils.StringPtr("project_test_name"),
						Kind: utils.StringPtr("project"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &ProjectSpec{
						Name:       "project_name",
						Descripion: "description_test",
						Resources: &ProjectResources{
							ResourceDomain: &ResourceDomain{
								Resources: []*Resources{
									{
										Limit:        utils.Int64Ptr(4),
										ResourceType: "resource_type_test",
									},
								},
							},
						},
					},
				},
			},
			&Project{
				APIVersion: "3.1",
				Metadata: &Metadata{
					Kind: utils.StringPtr("project"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.UpdateProject(tt.args.UUID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.UpdateProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.UpdateProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteProject(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/projects/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)

		fmt.Fprintf(w, `{
				"status": {
					"state": "DELETE_PENDING",
					"execution_context": {
						"task_uuid": "ff1b9547-dc9a-4ebd-a2ff-f2b718af935e"
					}
				},
				"spec": "",
				"api_version": "3.1",
				"metadata": {
					"kind": "projects",
					"categories": {
						"Project": "default"
					}
				}
			}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteProject OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			false,
		},

		{
			"Test DeleteProject Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if _, err := op.DeleteProject(tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_CreateAccessControlPolicy(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/access_control_policies", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"name": "access_control_policy_test_name",
				"kind": "access_control_policy",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"role_reference": map[string]interface{}{
						"name": "access_control_policy_test_name",
						"kind": "role",
						"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
					},
				},
				"name":        "access_control_policy_name",
				"description": "description_test",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "access_control_policy",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		request *AccessControlPolicy
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AccessControlPolicy
		wantErr bool
	}{
		{
			"Test CreateAccessControlPolicy",
			fields{c},
			args{
				&AccessControlPolicy{
					APIVersion: "3.1",
					Metadata: &Metadata{
						Name: utils.StringPtr("access_control_policy_test_name"),
						Kind: utils.StringPtr("access_control_policy"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &AccessControlPolicySpec{
						Name:        utils.StringPtr("access_control_policy_name"),
						Description: utils.StringPtr("description_test"),
						Resources: &AccessControlPolicyResources{
							RoleReference: &Reference{
								Kind: utils.StringPtr("role"),
								UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
								Name: utils.StringPtr("access_control_policy_test_name"),
							},
						},
					},
				},
			},
			&AccessControlPolicy{
				APIVersion: "3.1",
				Metadata: &Metadata{
					Kind: utils.StringPtr("access_control_policy"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateAccessControlPolicy(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateAccessControlPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateAccessControlPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_GetAccessControlPolicy(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/access_control_policies/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"host","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	hostResponse := &AccessControlPolicy{}
	hostResponse.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("host"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AccessControlPolicy
		wantErr bool
	}{
		{
			"Test GetAccessControlPolicy OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			hostResponse,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetAccessControlPolicy(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetAccessControlPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetAccessControlPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListAccessControlPolicy(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/access_control_policies/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"host","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	hostList := &AccessControlPolicyListResponse{}
	hostList.Entities = make([]*AccessControlPolicy, 1)
	hostList.Entities[0] = &AccessControlPolicy{}
	hostList.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("host"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AccessControlPolicyListResponse
		wantErr bool
	}{
		{
			"Test ListSubnet OK",
			fields{c},
			args{input},
			hostList,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListAccessControlPolicy(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListAccessControlPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListAccessControlPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UpdateAccessControlPolicy(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/access_control_policies/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "access_control_policy_test_name",
				"kind": "access_control_policy",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"role_reference": map[string]interface{}{
						"name": "access_control_policy_test_name_2",
						"kind": "role",
						"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
					},
				},
				"name":        "access_control_policy_name",
				"description": "description_test",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "access_control_policy",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
		body *AccessControlPolicy
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AccessControlPolicy
		wantErr bool
	}{
		{
			"Test CreateAccessControlPolicy",
			fields{c},
			args{
				"cfde831a-4e87-4a75-960f-89b0148aa2cc",
				&AccessControlPolicy{
					Metadata: &Metadata{
						Name: utils.StringPtr("access_control_policy_test_name"),
						Kind: utils.StringPtr("access_control_policy"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &AccessControlPolicySpec{
						Name:        utils.StringPtr("access_control_policy_name"),
						Description: utils.StringPtr("description_test"),
						Resources: &AccessControlPolicyResources{
							RoleReference: &Reference{
								Kind: utils.StringPtr("role"),
								UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
								Name: utils.StringPtr("access_control_policy_test_name_2"),
							},
						},
					},
				},
			},
			&AccessControlPolicy{
				APIVersion: "3.1",
				Metadata: &Metadata{
					Kind: utils.StringPtr("access_control_policy"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.UpdateAccessControlPolicy(tt.args.UUID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.UpdateAccessControlPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.UpdateAccessControlPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteAccessControlPolicy(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/access_control_policies/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)

		fmt.Fprintf(w, `{
				"status": {
					"state": "DELETE_PENDING",
					"execution_context": {
						"task_uuid": "ff1b9547-dc9a-4ebd-a2ff-f2b718af935e"
					}
				},
				"spec": "",
				"api_version": "3.1",
				"metadata": {
					"kind": "access_control_policy",
					"categories": {
						"Project": "default"
					}
				}
			}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteAccessControlPolicy OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			false,
		},

		{
			"Test DeleteAccessControlPolicy Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if _, err := op.DeleteAccessControlPolicy(tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteAccessControlPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_CreateRole(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/roles", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"name": "role_test_name",
				"kind": "role",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"permission_reference_list": []interface{}{
						map[string]interface{}{
							"name": "role_test_name",
							"kind": "role",
							"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
						},
					},
				},
				"name":        "role_name",
				"description": "description_test",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "role",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		request *Role
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Role
		wantErr bool
	}{
		{
			"Test CreateRole",
			fields{c},
			args{
				&Role{
					APIVersion: "3.1",
					Metadata: &Metadata{
						Name: utils.StringPtr("role_test_name"),
						Kind: utils.StringPtr("role"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &RoleSpec{
						Name:        utils.StringPtr("role_name"),
						Description: utils.StringPtr("description_test"),
						Resources: &RoleResources{
							PermissionReferenceList: []*Reference{
								{
									Kind: utils.StringPtr("role"),
									UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
									Name: utils.StringPtr("role_test_name"),
								},
							},
						},
					},
				},
			},
			&Role{
				APIVersion: "3.1",
				Metadata: &Metadata{
					Kind: utils.StringPtr("role"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateRole(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_GetRole(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/roles/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"host","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	hostResponse := &Role{}
	hostResponse.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("host"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Role
		wantErr bool
	}{
		{
			"Test GetRole OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			hostResponse,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetRole(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListRole(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/roles/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"host","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	hostList := &RoleListResponse{}
	hostList.Entities = make([]*Role, 1)
	hostList.Entities[0] = &Role{}
	hostList.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("host"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RoleListResponse
		wantErr bool
	}{
		{
			"Test ListSubnet OK",
			fields{c},
			args{input},
			hostList,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListRole(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UpdateRole(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/roles/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "role_test_name",
				"kind": "role",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"permission_reference_list": []interface{}{
						map[string]interface{}{
							"name": "role_test_name",
							"kind": "role",
							"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
						},
					},
				},
				"name":        "role_name",
				"description": "description_test",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "role",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
		body *Role
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Role
		wantErr bool
	}{
		{
			"Test CreateRole",
			fields{c},
			args{
				"cfde831a-4e87-4a75-960f-89b0148aa2cc",
				&Role{
					Metadata: &Metadata{
						Name: utils.StringPtr("role_test_name"),
						Kind: utils.StringPtr("role"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &RoleSpec{
						Name:        utils.StringPtr("role_name"),
						Description: utils.StringPtr("description_test"),
						Resources: &RoleResources{
							PermissionReferenceList: []*Reference{
								{
									Kind: utils.StringPtr("role"),
									UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
									Name: utils.StringPtr("role_test_name"),
								},
							},
						},
					},
				},
			},
			&Role{
				APIVersion: "3.1",
				Metadata: &Metadata{
					Kind: utils.StringPtr("role"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.UpdateRole(tt.args.UUID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.UpdateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.UpdateRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteRole(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/roles/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)

		fmt.Fprintf(w, `{
				"status": {
					"state": "DELETE_PENDING",
					"execution_context": {
						"task_uuid": "ff1b9547-dc9a-4ebd-a2ff-f2b718af935e"
					}
				},
				"spec": "",
				"api_version": "3.1",
				"metadata": {
					"kind": "role",
					"categories": {
						"Project": "default"
					}
				}
			}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteRole OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			false,
		},

		{
			"Test DeleteRole Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if _, err := op.DeleteRole(tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_CreateUser(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/users", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"name": "user_name",
				"kind": "user",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"directory_service_user": map[string]interface{}{
						"directory_service_reference": map[string]interface{}{
							"kind": "directory_service",
							"uuid": "d8b53470-c432-4556-badd-a11c937d89c9",
						},
						"user_principal_name": "user-dummy-tbd@ntnx.local",
					},
				},
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "user",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		request *UserIntentInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserIntentResponse
		wantErr bool
	}{
		{
			"Test CreateUser",
			fields{c},
			args{
				&UserIntentInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Name: utils.StringPtr("user_name"),
						Kind: utils.StringPtr("user"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &UserSpec{
						Resources: &UserResources{
							DirectoryServiceUser: &DirectoryServiceUser{
								DirectoryServiceReference: &Reference{
									Kind: utils.StringPtr("directory_service"),
									UUID: utils.StringPtr("d8b53470-c432-4556-badd-a11c937d89c9"),
								},
								UserPrincipalName: utils.StringPtr("user-dummy-tbd@ntnx.local"),
							},
						},
					},
				},
			},
			&UserIntentResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind: utils.StringPtr("user"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateUser(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_GetUser(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/users/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"user","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	hostResponse := &UserIntentResponse{}
	hostResponse.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("user"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserIntentResponse
		wantErr bool
	}{
		{
			"Test GetUser OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			hostResponse,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetUser(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListUser(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/users/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"user","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	hostList := &UserListResponse{}
	hostList.Entities = make([]*UserIntentResponse, 1)
	hostList.Entities[0] = &UserIntentResponse{}
	hostList.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("user"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserListResponse
		wantErr bool
	}{
		{
			"Test ListUser OK",
			fields{c},
			args{input},
			hostList,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListUser(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UpdateUser(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/users/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"name": "user_name",
				"kind": "user",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"directory_service_user": map[string]interface{}{
						"directory_service_reference": map[string]interface{}{
							"kind": "directory_service",
							"uuid": "d8b53470-c432-4556-badd-a11c937d89c9",
						},
						"user_principal_name": "user-dummy-tbd@ntnx.local",
					},
				},
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "user",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
		body *UserIntentInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserIntentResponse
		wantErr bool
	}{
		{
			"Test CreateUser",
			fields{c},
			args{
				"cfde831a-4e87-4a75-960f-89b0148aa2cc",
				&UserIntentInput{
					APIVersion: utils.StringPtr("3.1"),
					Metadata: &Metadata{
						Name: utils.StringPtr("user_name"),
						Kind: utils.StringPtr("user"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &UserSpec{
						Resources: &UserResources{
							DirectoryServiceUser: &DirectoryServiceUser{
								DirectoryServiceReference: &Reference{
									Kind: utils.StringPtr("directory_service"),
									UUID: utils.StringPtr("d8b53470-c432-4556-badd-a11c937d89c9"),
								},
								UserPrincipalName: utils.StringPtr("user-dummy-tbd@ntnx.local"),
							},
						},
					},
				},
			},
			&UserIntentResponse{
				APIVersion: utils.StringPtr("3.1"),
				Metadata: &Metadata{
					Kind: utils.StringPtr("user"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.UpdateUser(tt.args.UUID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.UpdateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteUser(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/users/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)

		fmt.Fprintf(w, `{
				"status": {
					"state": "DELETE_PENDING",
					"execution_context": {
						"task_uuid": "ff1b9547-dc9a-4ebd-a2ff-f2b718af935e"
					}
				},
				"spec": "",
				"api_version": "3.1",
				"metadata": {
					"kind": "user",
					"categories": {
						"Project": "default"
					}
				}
			}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteUser OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			false,
		},

		{
			"Test DeleteUser Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if _, err := op.DeleteUser(tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_CreateProtectionRule(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/protection_rules", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"name": "protection_rule_test_name",
				"kind": "protection_rule",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"start_time": "00h:00m",
					"ordered_availability_zone_list": []interface{}{
						map[string]interface{}{
							"availability_zone_url": "url test",
							"cluster_uuid":          "cfde831a-4e87-4a75-960f-89b0148aa2cc",
						},
					},
					"availability_zone_connectivity_list": []interface{}{
						map[string]interface{}{
							"destination_availability_zone_index": float64(0),
							"source_availability_zone_index":      float64(0),
							"snapshot_schedule_list": []interface{}{
								map[string]interface{}{
									"recovery_point_objective_secs": float64(0),
									"auto_suspend_timeout_secs":     float64(0),
									"snapshot_type":                 "CRASH_CONSISTENT",
									"local_snapshot_retention_policy": map[string]interface{}{
										"num_snapshots": float64(1),
										"rollup_retention_policy": map[string]interface{}{
											"snapshot_interval_type": "HOURLY",
											"multiple":               float64(1),
										},
									},
									"remote_snapshot_retention_policy": map[string]interface{}{
										"num_snapshots": float64(1),
										"rollup_retention_policy": map[string]interface{}{
											"snapshot_interval_type": "HOURLY",
											"multiple":               float64(1),
										},
									},
								},
							},
						},
					},
					"category_filter": map[string]interface{}{
						"type":      "CATEGORIES_MATCH_ALL",
						"kind_list": []interface{}{"1", "2"},
					},
				},
				"name":        "protection_rule_name",
				"description": "description_test",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		assert := assert.New(t)
		if !assert.Equal(v, expected, "The response should be the same") {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "protection_rule",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		request *ProtectionRuleInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ProtectionRuleResponse
		wantErr bool
	}{
		{
			"Test CreateProtectionRule",
			fields{c},
			args{
				&ProtectionRuleInput{
					APIVersion: "3.1",
					Metadata: &Metadata{
						Name: utils.StringPtr("protection_rule_test_name"),
						Kind: utils.StringPtr("protection_rule"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &ProtectionRuleSpec{
						Name:        "protection_rule_name",
						Description: "description_test",
						Resources: &ProtectionRuleResources{
							StartTime: "00h:00m",
							OrderedAvailabilityZoneList: []*OrderedAvailabilityZoneList{
								{
									AvailabilityZoneURL: "url test",
									ClusterUUID:         "cfde831a-4e87-4a75-960f-89b0148aa2cc",
								},
							},
							AvailabilityZoneConnectivityList: []*AvailabilityZoneConnectivityList{
								{
									DestinationAvailabilityZoneIndex: utils.Int64Ptr(0),
									SourceAvailabilityZoneIndex:      utils.Int64Ptr(0),
									SnapshotScheduleList: []*SnapshotScheduleList{
										{
											RecoveryPointObjectiveSecs: utils.Int64Ptr(0),
											AutoSuspendTimeoutSecs:     utils.Int64Ptr(0),
											SnapshotType:               "CRASH_CONSISTENT",
											LocalSnapshotRetentionPolicy: &SnapshotRetentionPolicy{
												NumSnapshots: utils.Int64Ptr(1),
												RollupRetentionPolicy: &RollupRetentionPolicy{
													SnapshotIntervalType: "HOURLY",
													Multiple:             utils.Int64Ptr(1),
												},
											},
											RemoteSnapshotRetentionPolicy: &SnapshotRetentionPolicy{
												NumSnapshots: utils.Int64Ptr(1),
												RollupRetentionPolicy: &RollupRetentionPolicy{
													SnapshotIntervalType: "HOURLY",
													Multiple:             utils.Int64Ptr(1),
												},
											},
										},
									},
								},
							},
							CategoryFilter: &CategoryFilter{
								Type:     utils.StringPtr("CATEGORIES_MATCH_ALL"),
								KindList: []*string{utils.StringPtr("1"), utils.StringPtr("2")},
							},
						},
					},
				},
			},
			&ProtectionRuleResponse{
				APIVersion: "3.1",
				Metadata: &Metadata{
					Kind: utils.StringPtr("protection_rule"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateProtectionRule(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateProtectionRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateProtectionRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_GetProtectionRule(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/protection_rules/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"protection_rule","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	response := &ProtectionRuleResponse{}
	response.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("protection_rule"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ProtectionRuleResponse
		wantErr bool
	}{
		{
			"Test GetProtectionRules OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			response,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetProtectionRule(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetProtectionRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetProtectionRules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListProtectionRules(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/protection_rules/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"protection_rule","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	responseList := &ProtectionRulesListResponse{}
	responseList.Entities = make([]*ProtectionRuleResponse, 1)
	responseList.Entities[0] = &ProtectionRuleResponse{}
	responseList.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("protection_rule"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ProtectionRulesListResponse
		wantErr bool
	}{
		{
			"Test ListProtectionRules OK",
			fields{c},
			args{input},
			responseList,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListProtectionRules(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListProtectionRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListProtectionRules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UpdateProtectionRules(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/protection_rules/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"name": "protection_rule_test_name",
				"kind": "protection_rule",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"start_time": "00h:00m",
					"ordered_availability_zone_list": []interface{}{
						map[string]interface{}{
							"availability_zone_url": "url test",
							"cluster_uuid":          "cfde831a-4e87-4a75-960f-89b0148aa2cc",
						},
					},
					"availability_zone_connectivity_list": []interface{}{
						map[string]interface{}{
							"destination_availability_zone_index": float64(0),
							"source_availability_zone_index":      float64(0),
							"snapshot_schedule_list": []interface{}{
								map[string]interface{}{
									"recovery_point_objective_secs": float64(0),
									"auto_suspend_timeout_secs":     float64(0),
									"snapshot_type":                 "CRASH_CONSISTENT",
									"local_snapshot_retention_policy": map[string]interface{}{
										"num_snapshots": float64(1),
										"rollup_retention_policy": map[string]interface{}{
											"snapshot_interval_type": "HOURLY",
											"multiple":               float64(1),
										},
									},
									"remote_snapshot_retention_policy": map[string]interface{}{
										"num_snapshots": float64(1),
										"rollup_retention_policy": map[string]interface{}{
											"snapshot_interval_type": "HOURLY",
											"multiple":               float64(1),
										},
									},
								},
							},
						},
					},
					"category_filter": map[string]interface{}{
						"type":      "CATEGORIES_MATCH_ALL",
						"kind_list": []interface{}{"1", "2"},
					},
				},
				"name":        "protection_rule_name",
				"description": "description_test",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		assert := assert.New(t)
		if !assert.Equal(v, expected, "The response should be the same") {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "protection_rule",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
		body *ProtectionRuleInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ProtectionRuleResponse
		wantErr bool
	}{
		{
			"Test CreateProtectionRule",
			fields{c},
			args{
				"cfde831a-4e87-4a75-960f-89b0148aa2cc",
				&ProtectionRuleInput{
					APIVersion: "3.1",
					Metadata: &Metadata{
						Name: utils.StringPtr("protection_rule_test_name"),
						Kind: utils.StringPtr("protection_rule"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &ProtectionRuleSpec{
						Name:        "protection_rule_name",
						Description: "description_test",
						Resources: &ProtectionRuleResources{
							StartTime: "00h:00m",
							OrderedAvailabilityZoneList: []*OrderedAvailabilityZoneList{
								{
									AvailabilityZoneURL: "url test",
									ClusterUUID:         "cfde831a-4e87-4a75-960f-89b0148aa2cc",
								},
							},
							AvailabilityZoneConnectivityList: []*AvailabilityZoneConnectivityList{
								{
									DestinationAvailabilityZoneIndex: utils.Int64Ptr(0),
									SourceAvailabilityZoneIndex:      utils.Int64Ptr(0),
									SnapshotScheduleList: []*SnapshotScheduleList{
										{
											RecoveryPointObjectiveSecs: utils.Int64Ptr(0),
											AutoSuspendTimeoutSecs:     utils.Int64Ptr(0),
											SnapshotType:               "CRASH_CONSISTENT",
											LocalSnapshotRetentionPolicy: &SnapshotRetentionPolicy{
												NumSnapshots: utils.Int64Ptr(1),
												RollupRetentionPolicy: &RollupRetentionPolicy{
													SnapshotIntervalType: "HOURLY",
													Multiple:             utils.Int64Ptr(1),
												},
											},
											RemoteSnapshotRetentionPolicy: &SnapshotRetentionPolicy{
												NumSnapshots: utils.Int64Ptr(1),
												RollupRetentionPolicy: &RollupRetentionPolicy{
													SnapshotIntervalType: "HOURLY",
													Multiple:             utils.Int64Ptr(1),
												},
											},
										},
									},
								},
							},
							CategoryFilter: &CategoryFilter{
								Type:     utils.StringPtr("CATEGORIES_MATCH_ALL"),
								KindList: []*string{utils.StringPtr("1"), utils.StringPtr("2")},
							},
						},
					},
				},
			},
			&ProtectionRuleResponse{
				APIVersion: "3.1",
				Metadata: &Metadata{
					Kind: utils.StringPtr("protection_rule"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.UpdateProtectionRule(tt.args.UUID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.UpdateProtectionRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.UpdateProtectionRules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteProtectionRules(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/protection_rules/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteProtectionRules OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			true,
		},

		{
			"Test DeleteProtectionRules Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if _, err := op.DeleteProtectionRule(tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteProtectionRules() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_CreateRecoveryPlan(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/recovery_plans", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"name": "recovery_plan_test_name",
				"kind": "recovery_plan",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"parameters": map[string]interface{}{
						"network_mapping_list": []interface{}{
							map[string]interface{}{
								"are_networks_stretched": false,
								"availability_zone_network_mapping_list": []interface{}{
									map[string]interface{}{
										"availability_zone_url": "zone url",
										"recovery_network": map[string]interface{}{
											"use_vpc_reference": true,
											"virtual_network_reference": map[string]interface{}{
												"name": "recovery_plan_test_name",
												"kind": "recovery_plan",
												"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											},
											"vpc_reference": map[string]interface{}{
												"name": "recovery_plan_test_name",
												"kind": "recovery_plan",
												"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											},
											"subnet_list": []interface{}{
												map[string]interface{}{
													"gateway_ip":                  "127.0.0.1",
													"prefix_length":               float64(16),
													"external_connectivity_state": "external",
												},
											},
										},
										"test_network": map[string]interface{}{
											"use_vpc_reference": true,
											"virtual_network_reference": map[string]interface{}{
												"name": "recovery_plan_test_name",
												"kind": "recovery_plan",
												"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											},
											"vpc_reference": map[string]interface{}{
												"name": "recovery_plan_test_name",
												"kind": "recovery_plan",
												"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											},
											"subnet_list": []interface{}{
												map[string]interface{}{
													"gateway_ip":                  "127.0.0.1",
													"prefix_length":               float64(16),
													"external_connectivity_state": "external",
												},
											},
										},
										"recovery_ip_assignment_list": []interface{}{
											map[string]interface{}{
												"vm_reference": map[string]interface{}{
													"name": "recovery_plan_test_name",
													"kind": "recovery_plan",
													"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
												},
												"ip_config_list": []interface{}{
													map[string]interface{}{
														"ip_address": "127.0.0.1",
													},
												},
											},
										},
										"test_ip_assignment_list": []interface{}{
											map[string]interface{}{
												"vm_reference": map[string]interface{}{
													"name": "recovery_plan_test_name",
													"kind": "recovery_plan",
													"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
												},
												"ip_config_list": []interface{}{
													map[string]interface{}{
														"ip_address": "127.0.0.1",
													},
												},
											},
										},
										"cluster_reference_list": []interface{}{
											map[string]interface{}{
												"name": "recovery_plan_test_name",
												"kind": "recovery_plan",
												"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											},
										},
									},
								},
							},
						},
						"floating_ip_assignment_list": []interface{}{
							map[string]interface{}{
								"availability_zone_url": "zone url",
								"vm_ip_assignment_list": []interface{}{
									map[string]interface{}{
										"vm_reference": map[string]interface{}{
											"name": "recovery_plan_test_name",
											"kind": "recovery_plan",
											"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
										},
										"vm_nic_information": map[string]interface{}{
											"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											"ip":   "127.0.0.1",
										},
										"test_floating_ip_config": map[string]interface{}{
											"should_allocate_dynamically": false,
											"ip":                          "127.0.0.1",
										},
										"recovery_floating_ip_config": map[string]interface{}{
											"should_allocate_dynamically": false,
											"ip":                          "127.0.0.1",
										},
									},
								},
							},
						},
					},
					"stage_list": []interface{}{
						map[string]interface{}{
							"stage_uuid":      "cfde831a-4e87-4a75-960f-89b0148aa2cc",
							"delay_time_secs": float64(0),
							"stage_work": map[string]interface{}{
								"recover_entities": map[string]interface{}{
									"entity_info_list": []interface{}{
										map[string]interface{}{
											"script_list": []interface{}{
												map[string]interface{}{
													"enable_script_exec": false,
													"timeout":            float64(0),
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"name":        "recovery_plan_name",
				"description": "description_test",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		assert := assert.New(t)
		if !assert.Equal(v, expected, "The response should be the same") {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "recovery_plan",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		request *RecoveryPlanInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RecoveryPlanResponse
		wantErr bool
	}{
		{
			"Test CreateRecoveryPlans",
			fields{c},
			args{
				&RecoveryPlanInput{
					APIVersion: "3.1",
					Metadata: &Metadata{
						Name: utils.StringPtr("recovery_plan_test_name"),
						Kind: utils.StringPtr("recovery_plan"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &RecoveryPlanSpec{
						Name:        "recovery_plan_name",
						Description: "description_test",
						Resources: &RecoveryPlanResources{
							StageList: []*StageList{
								{
									StageUUID:     "cfde831a-4e87-4a75-960f-89b0148aa2cc",
									DelayTimeSecs: utils.Int64Ptr(0),
									StageWork: &StageWork{
										RecoverEntities: &RecoverEntities{
											EntityInfoList: []*EntityInfoList{
												{
													ScriptList: []*ScriptList{
														{
															EnableScriptExec: utils.BoolPtr(false),
															Timeout:          utils.Int64Ptr(0),
														},
													},
												},
											},
										},
									},
								},
							},
							Parameters: &Parameters{
								FloatingIPAssignmentList: []*FloatingIPAssignmentList{
									{
										AvailabilityZoneURL: "zone url",
										VMIPAssignmentList: []*VMIPAssignmentList{
											{
												VMReference: &Reference{
													Name: utils.StringPtr("recovery_plan_test_name"),
													Kind: utils.StringPtr("recovery_plan"),
													UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
												},
												VMNICInformation: &VMNICInformation{
													UUID: "cfde831a-4e87-4a75-960f-89b0148aa2cc",
													IP:   "127.0.0.1",
												},
												TestFloatingIPConfig: &FloatingIPConfig{
													IP:                        "127.0.0.1",
													ShouldAllocateDynamically: utils.BoolPtr(false),
												},
												RecoveryFloatingIPConfig: &FloatingIPConfig{
													IP:                        "127.0.0.1",
													ShouldAllocateDynamically: utils.BoolPtr(false),
												},
											},
										},
									},
								},
								NetworkMappingList: []*NetworkMappingList{
									{
										AreNetworksStretched: utils.BoolPtr(false),
										AvailabilityZoneNetworkMappingList: []*AvailabilityZoneNetworkMappingList{
											{
												AvailabilityZoneURL: "zone url",
												RecoveryNetwork: &Network{
													UseVPCReference: utils.BoolPtr(true),
													VirtualNetworkReference: &Reference{
														Name: utils.StringPtr("recovery_plan_test_name"),
														Kind: utils.StringPtr("recovery_plan"),
														UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
													},
													VPCReference: &Reference{
														Name: utils.StringPtr("recovery_plan_test_name"),
														Kind: utils.StringPtr("recovery_plan"),
														UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
													},
													SubnetList: []*SubnetList{
														{
															GatewayIP:                 "127.0.0.1",
															PrefixLength:              utils.Int64Ptr(16),
															ExternalConnectivityState: "external",
														},
													},
												},
												TestNetwork: &Network{
													UseVPCReference: utils.BoolPtr(true),
													VirtualNetworkReference: &Reference{
														Name: utils.StringPtr("recovery_plan_test_name"),
														Kind: utils.StringPtr("recovery_plan"),
														UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
													},
													VPCReference: &Reference{
														Name: utils.StringPtr("recovery_plan_test_name"),
														Kind: utils.StringPtr("recovery_plan"),
														UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
													},
													SubnetList: []*SubnetList{
														{
															GatewayIP:                 "127.0.0.1",
															PrefixLength:              utils.Int64Ptr(16),
															ExternalConnectivityState: "external",
														},
													},
												},
												RecoveryIPAssignmentList: []*IPAssignmentList{
													{
														VMReference: &Reference{
															Name: utils.StringPtr("recovery_plan_test_name"),
															Kind: utils.StringPtr("recovery_plan"),
															UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
														},
														IPConfigList: []*IPConfigList{
															{
																IPAddress: "127.0.0.1",
															},
														},
													},
												},
												TestIPAssignmentList: []*IPAssignmentList{
													{
														VMReference: &Reference{
															Name: utils.StringPtr("recovery_plan_test_name"),
															Kind: utils.StringPtr("recovery_plan"),
															UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
														},
														IPConfigList: []*IPConfigList{
															{
																IPAddress: "127.0.0.1",
															},
														},
													},
												},
												ClusterReferenceList: []*Reference{
													{
														Name: utils.StringPtr("recovery_plan_test_name"),
														Kind: utils.StringPtr("recovery_plan"),
														UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&RecoveryPlanResponse{
				APIVersion: "3.1",
				Metadata: &Metadata{
					Kind: utils.StringPtr("recovery_plan"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateRecoveryPlan(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateRecoveryPlans() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateRecoveryPlans() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_GetRecoveryPlan(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/recovery_plans/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"recovery_plan","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	response := &RecoveryPlanResponse{}
	response.Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("recovery_plan"),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RecoveryPlanResponse
		wantErr bool
	}{
		{
			"Test GetRecoveryPlan OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			response,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetRecoveryPlan(tt.args.UUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetRecoveryPlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetRecoveryPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListRecoveryPlans(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/recovery_plans/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"recovery_plan","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	responseList := &RecoveryPlanListResponse{}
	responseList.Entities = make([]*RecoveryPlanResponse, 1)
	responseList.Entities[0] = &RecoveryPlanResponse{}
	responseList.Entities[0].Metadata = &Metadata{
		UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.StringPtr("recovery_plan"),
	}

	input := &DSMetadata{
		Length: utils.Int64Ptr(1.0),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *DSMetadata
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RecoveryPlanListResponse
		wantErr bool
	}{
		{
			"Test ListRecoveryPlans OK",
			fields{c},
			args{input},
			responseList,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListRecoveryPlans(tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListRecoveryPlans() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListRecoveryPlans() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_UpdateRecoveryPlans(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/recovery_plans/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPut)

		expected := map[string]interface{}{
			"api_version": "3.1",
			"metadata": map[string]interface{}{
				"name": "recovery_plan_test_name",
				"kind": "recovery_plan",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
			},
			"spec": map[string]interface{}{
				"resources": map[string]interface{}{
					"parameters": map[string]interface{}{
						"network_mapping_list": []interface{}{
							map[string]interface{}{
								"are_networks_stretched": false,
								"availability_zone_network_mapping_list": []interface{}{
									map[string]interface{}{
										"availability_zone_url": "zone url",
										"recovery_network": map[string]interface{}{
											"use_vpc_reference": true,
											"virtual_network_reference": map[string]interface{}{
												"name": "recovery_plan_test_name",
												"kind": "recovery_plan",
												"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											},
											"vpc_reference": map[string]interface{}{
												"name": "recovery_plan_test_name",
												"kind": "recovery_plan",
												"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											},
											"subnet_list": []interface{}{
												map[string]interface{}{
													"gateway_ip":                  "127.0.0.1",
													"prefix_length":               float64(16),
													"external_connectivity_state": "external",
												},
											},
										},
										"test_network": map[string]interface{}{
											"use_vpc_reference": true,
											"virtual_network_reference": map[string]interface{}{
												"name": "recovery_plan_test_name",
												"kind": "recovery_plan",
												"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											},
											"vpc_reference": map[string]interface{}{
												"name": "recovery_plan_test_name",
												"kind": "recovery_plan",
												"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											},
											"subnet_list": []interface{}{
												map[string]interface{}{
													"gateway_ip":                  "127.0.0.1",
													"prefix_length":               float64(16),
													"external_connectivity_state": "external",
												},
											},
										},
										"recovery_ip_assignment_list": []interface{}{
											map[string]interface{}{
												"vm_reference": map[string]interface{}{
													"name": "recovery_plan_test_name",
													"kind": "recovery_plan",
													"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
												},
												"ip_config_list": []interface{}{
													map[string]interface{}{
														"ip_address": "127.0.0.1",
													},
												},
											},
										},
										"test_ip_assignment_list": []interface{}{
											map[string]interface{}{
												"vm_reference": map[string]interface{}{
													"name": "recovery_plan_test_name",
													"kind": "recovery_plan",
													"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
												},
												"ip_config_list": []interface{}{
													map[string]interface{}{
														"ip_address": "127.0.0.1",
													},
												},
											},
										},
										"cluster_reference_list": []interface{}{
											map[string]interface{}{
												"name": "recovery_plan_test_name",
												"kind": "recovery_plan",
												"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											},
										},
									},
								},
							},
						},
						"floating_ip_assignment_list": []interface{}{
							map[string]interface{}{
								"availability_zone_url": "zone url",
								"vm_ip_assignment_list": []interface{}{
									map[string]interface{}{
										"vm_reference": map[string]interface{}{
											"name": "recovery_plan_test_name",
											"kind": "recovery_plan",
											"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
										},
										"vm_nic_information": map[string]interface{}{
											"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
											"ip":   "127.0.0.1",
										},
										"test_floating_ip_config": map[string]interface{}{
											"should_allocate_dynamically": false,
											"ip":                          "127.0.0.1",
										},
										"recovery_floating_ip_config": map[string]interface{}{
											"should_allocate_dynamically": false,
											"ip":                          "127.0.0.1",
										},
									},
								},
							},
						},
					},
					"stage_list": []interface{}{
						map[string]interface{}{
							"stage_uuid":      "cfde831a-4e87-4a75-960f-89b0148aa2cc",
							"delay_time_secs": float64(0),
							"stage_work": map[string]interface{}{
								"recover_entities": map[string]interface{}{
									"entity_info_list": []interface{}{
										map[string]interface{}{
											"script_list": []interface{}{
												map[string]interface{}{
													"enable_script_exec": false,
													"timeout":            float64(0),
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"name":        "recovery_plan_name",
				"description": "description_test",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		assert := assert.New(t)
		if !assert.Equal(v, expected, "The response should be the same") {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{
			"api_version": "3.1",
			"metadata": {
				"kind": "recovery_plan",
				"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"
			}
		}`)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
		body *RecoveryPlanInput
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RecoveryPlanResponse
		wantErr bool
	}{
		{
			"Test UpdateRecoveryPlans",
			fields{c},
			args{
				"cfde831a-4e87-4a75-960f-89b0148aa2cc",
				&RecoveryPlanInput{
					APIVersion: "3.1",
					Metadata: &Metadata{
						Name: utils.StringPtr("recovery_plan_test_name"),
						Kind: utils.StringPtr("recovery_plan"),
						UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
					},
					Spec: &RecoveryPlanSpec{
						Name:        "recovery_plan_name",
						Description: "description_test",
						Resources: &RecoveryPlanResources{
							StageList: []*StageList{
								{
									StageUUID:     "cfde831a-4e87-4a75-960f-89b0148aa2cc",
									DelayTimeSecs: utils.Int64Ptr(0),
									StageWork: &StageWork{
										RecoverEntities: &RecoverEntities{
											EntityInfoList: []*EntityInfoList{
												{
													ScriptList: []*ScriptList{
														{
															EnableScriptExec: utils.BoolPtr(false),
															Timeout:          utils.Int64Ptr(0),
														},
													},
												},
											},
										},
									},
								},
							},
							Parameters: &Parameters{
								FloatingIPAssignmentList: []*FloatingIPAssignmentList{
									{
										AvailabilityZoneURL: "zone url",
										VMIPAssignmentList: []*VMIPAssignmentList{
											{
												VMReference: &Reference{
													Name: utils.StringPtr("recovery_plan_test_name"),
													Kind: utils.StringPtr("recovery_plan"),
													UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
												},
												VMNICInformation: &VMNICInformation{
													UUID: "cfde831a-4e87-4a75-960f-89b0148aa2cc",
													IP:   "127.0.0.1",
												},
												TestFloatingIPConfig: &FloatingIPConfig{
													IP:                        "127.0.0.1",
													ShouldAllocateDynamically: utils.BoolPtr(false),
												},
												RecoveryFloatingIPConfig: &FloatingIPConfig{
													IP:                        "127.0.0.1",
													ShouldAllocateDynamically: utils.BoolPtr(false),
												},
											},
										},
									},
								},
								NetworkMappingList: []*NetworkMappingList{
									{
										AreNetworksStretched: utils.BoolPtr(false),
										AvailabilityZoneNetworkMappingList: []*AvailabilityZoneNetworkMappingList{
											{
												AvailabilityZoneURL: "zone url",
												RecoveryNetwork: &Network{
													UseVPCReference: utils.BoolPtr(true),
													VirtualNetworkReference: &Reference{
														Name: utils.StringPtr("recovery_plan_test_name"),
														Kind: utils.StringPtr("recovery_plan"),
														UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
													},
													VPCReference: &Reference{
														Name: utils.StringPtr("recovery_plan_test_name"),
														Kind: utils.StringPtr("recovery_plan"),
														UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
													},
													SubnetList: []*SubnetList{
														{
															GatewayIP:                 "127.0.0.1",
															PrefixLength:              utils.Int64Ptr(16),
															ExternalConnectivityState: "external",
														},
													},
												},
												TestNetwork: &Network{
													UseVPCReference: utils.BoolPtr(true),
													VirtualNetworkReference: &Reference{
														Name: utils.StringPtr("recovery_plan_test_name"),
														Kind: utils.StringPtr("recovery_plan"),
														UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
													},
													VPCReference: &Reference{
														Name: utils.StringPtr("recovery_plan_test_name"),
														Kind: utils.StringPtr("recovery_plan"),
														UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
													},
													SubnetList: []*SubnetList{
														{
															GatewayIP:                 "127.0.0.1",
															PrefixLength:              utils.Int64Ptr(16),
															ExternalConnectivityState: "external",
														},
													},
												},
												RecoveryIPAssignmentList: []*IPAssignmentList{
													{
														VMReference: &Reference{
															Name: utils.StringPtr("recovery_plan_test_name"),
															Kind: utils.StringPtr("recovery_plan"),
															UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
														},
														IPConfigList: []*IPConfigList{
															{
																IPAddress: "127.0.0.1",
															},
														},
													},
												},
												TestIPAssignmentList: []*IPAssignmentList{
													{
														VMReference: &Reference{
															Name: utils.StringPtr("recovery_plan_test_name"),
															Kind: utils.StringPtr("recovery_plan"),
															UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
														},
														IPConfigList: []*IPConfigList{
															{
																IPAddress: "127.0.0.1",
															},
														},
													},
												},
												ClusterReferenceList: []*Reference{
													{
														Name: utils.StringPtr("recovery_plan_test_name"),
														Kind: utils.StringPtr("recovery_plan"),
														UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&RecoveryPlanResponse{
				APIVersion: "3.1",
				Metadata: &Metadata{
					Kind: utils.StringPtr("recovery_plan"),
					UUID: utils.StringPtr("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.UpdateRecoveryPlan(tt.args.UUID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.UpdateRecoveryPlans() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.UpdateRecoveryPlans() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteRecoveryPlan(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/nutanix/v3/recovery_plans/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test DeleteRecoveryPlans OK",
			fields{c},
			args{"cfde831a-4e87-4a75-960f-89b0148aa2cc"},
			true,
		},

		{
			"Test DeleteRecoveryPlans Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if _, err := op.DeleteRecoveryPlan(tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteRecoveryPlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
