package v3

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

var (
	mux    *http.ServeMux
	c      *client.Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	c, _ = client.NewClient(&client.Credentials{
		URL:      "",
		Username: "username",
		Password: "password",
		Port:     "",
		Endpoint: "",
		Insecure: true})
	c.BaseURL, _ = url.Parse(server.URL)
}

func teardown() {
	server.Close()
}

func TestOperations_CreateVM(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/vms", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodPost; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}

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
					APIVersion: utils.String("3.1"),
					Metadata: &Metadata{
						Kind: utils.String("vm"),
					},
					Spec: &VM{
						ClusterReference: &Reference{
							Kind: utils.String("cluster"),
							UUID: utils.String("00056024-6c13-4c74-0000-00000000ecb5"),
						},
						Name: utils.String("VM123.create"),
					},
				},
			},
			&VMIntentResponse{
				APIVersion: utils.String("3.1"),
				Metadata: &Metadata{
					Kind: utils.String("vm"),
					UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
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
	setup()
	defer teardown()

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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/vms/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"vm","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	vmResponse := &VMIntentResponse{}
	vmResponse.Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("vm"),
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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/vms/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"vm","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	vmList := &VMListIntentResponse{}
	vmList.Entities = make([]*VMIntentResource, 1)
	vmList.Entities[0] = &VMIntentResource{}
	vmList.Entities[0].Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("vm"),
	}

	input := &DSMetadata{
		Length: utils.Int64(1.0),
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
	setup()
	defer teardown()

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
					APIVersion: utils.String("3.1"),
					Metadata: &Metadata{
						Kind: utils.String("vm"),
					},
					Spec: &VM{
						ClusterReference: &Reference{
							Kind: utils.String("cluster"),
							UUID: utils.String("00056024-6c13-4c74-0000-00000000ecb5"),
						},
						Name: utils.String("VM123.create"),
					},
				},
			},
			&VMIntentResponse{
				APIVersion: utils.String("3.1"),
				Metadata: &Metadata{
					Kind: utils.String("vm"),
					UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
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
	setup()
	defer teardown()

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
					APIVersion: utils.String("3.1"),
					Metadata: &Metadata{
						Kind: utils.String("subnet"),
					},
					Spec: &Subnet{
						ClusterReference: &Reference{
							Kind: utils.String("cluster"),
							UUID: utils.String("00056024-6c13-4c74-0000-00000000ecb5"),
						},
						Name: utils.String("subnet.create"),
					},
				},
			},
			&SubnetIntentResponse{
				APIVersion: utils.String("3.1"),
				Metadata: &Metadata{
					Kind: utils.String("subnet"),
					UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
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
	setup()
	defer teardown()

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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/subnets/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"subnet","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	subnetResponse := &SubnetIntentResponse{}
	subnetResponse.Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("subnet"),
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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/subnets/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"subnet","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	subnetList := &SubnetListIntentResponse{}
	subnetList.Entities = make([]*SubnetIntentResponse, 1)
	subnetList.Entities[0] = &SubnetIntentResponse{}
	subnetList.Entities[0].Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("subnet"),
	}

	input := &DSMetadata{
		Length: utils.Int64(1.0),
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
	setup()
	defer teardown()

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
					APIVersion: utils.String("3.1"),
					Metadata: &Metadata{
						Kind: utils.String("subnet"),
					},
					Spec: &Subnet{
						ClusterReference: &Reference{
							Kind: utils.String("cluster"),
							UUID: utils.String("00056024-6c13-4c74-0000-00000000ecb5"),
						},
						Name: utils.String("subnet.create"),
					},
				},
			},
			&SubnetIntentResponse{
				APIVersion: utils.String("3.1"),
				Metadata: &Metadata{
					Kind: utils.String("subnet"),
					UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
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
	setup()
	defer teardown()

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
					APIVersion: utils.String("3.1"),
					Metadata: &Metadata{
						Kind: utils.String("image"),
					},
					Spec: &Image{
						Name: utils.String("image.create"),
						Resources: &ImageResources{
							ImageType: utils.String("DISK_IMAGE"),
						},
					},
				},
			},
			&ImageIntentResponse{
				APIVersion: utils.String("3.1"),
				Metadata: &Metadata{
					Kind: utils.String("image"),
					UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
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

func TestOperations_UploadImage(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if err := op.UploadImage(tt.args.UUID, tt.args.filepath); (err != nil) != tt.wantErr {
				t.Errorf("Operations.UploadImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperations_DeleteImage(t *testing.T) {
	setup()
	defer teardown()

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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/images/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"image","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	response := &ImageIntentResponse{}
	response.Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("image"),
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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/images/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"image","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	list := &ImageListIntentResponse{}
	list.Entities = make([]*ImageIntentResource, 1)
	list.Entities[0] = &ImageIntentResource{}
	list.Entities[0].Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("image"),
	}

	input := &DSMetadata{
		Length: utils.Int64(1.0),
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
	setup()
	defer teardown()

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
					APIVersion: utils.String("3.1"),
					Metadata: &Metadata{
						Kind: utils.String("image"),
					},
					Spec: &Image{
						Resources: &ImageResources{
							ImageType: utils.String("DISK_IMAGE"),
						},
						Name: utils.String("image.update"),
					},
				},
			},
			&ImageIntentResponse{
				APIVersion: utils.String("3.1"),
				Metadata: &Metadata{
					Kind: utils.String("image"),
					UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/clusters/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"cluster","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	response := &ClusterIntentResponse{}
	response.Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("cluster"),
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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/clusters/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"cluster","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	list := &ClusterListIntentResponse{}
	list.Entities = make([]*ClusterIntentResource, 1)
	list.Entities[0] = &ClusterIntentResource{}
	list.Entities[0].Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("cluster"),
	}

	input := &DSMetadata{
		Length: utils.Int64(1.0),
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
	setup()
	defer teardown()

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
				Description: utils.String("Testing Keys"),
				Name:        utils.String("test_category_key")}},
			&CategoryKeyStatus{
				Description:   utils.String("Testing Keys"),
				Name:          utils.String("test_category_key"),
				SystemDefined: utils.Bool(false)},
			false,
		},
	}
	for _, tt := range tests {
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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/categories/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{ "description": "Testing Keys", "name": "test_category_key", "system_defined": false }]}`)
	})

	list := &CategoryKeyListResponse{}
	list.Entities = make([]*CategoryKeyStatus, 1)
	list.Entities[0] = &CategoryKeyStatus{
		Description:   utils.String("Testing Keys"),
		Name:          utils.String("test_category_key"),
		SystemDefined: utils.Bool(false)}

	input := &CategoryListMetadata{
		Length: utils.Int64(1.0),
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
	setup()
	defer teardown()

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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/categories/test_category_key", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"description": "Testing Keys",
			"name": "test_category_key",
			"system_defined": false
		}`)
	})

	response := &CategoryKeyStatus{
		Description:   utils.String("Testing Keys"),
		Name:          utils.String("test_category_key"),
		SystemDefined: utils.Bool(false),
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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/categories/test_category_key/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{ "description": "Testing Keys", "value": "test_category_value", "system_defined": false }]}`)
	})

	list := &CategoryValueListResponse{}
	list.Entities = make([]*CategoryValueStatus, 1)
	list.Entities[0] = &CategoryValueStatus{
		Description:   utils.String("Testing Keys"),
		Value:         utils.String("test_category_value"),
		SystemDefined: utils.Bool(false)}

	input := &CategoryListMetadata{
		Length: utils.Int64(1.0),
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
	setup()
	defer teardown()

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
				Description: utils.String("Testing Value"),
				Value:       utils.String("test_category_value")}},
			&CategoryValueStatus{
				Description:   utils.String("Testing Value"),
				Value:         utils.String("test_category_value"),
				Name:          utils.String("test_category_key"),
				SystemDefined: utils.Bool(false)},
			false,
		},
	}
	for _, tt := range tests {
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
	setup()
	defer teardown()

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
		Description:   utils.String("Testing Value"),
		Name:          utils.String("test_category_key"),
		Value:         utils.String("test_category_value"),
		SystemDefined: utils.Bool(false),
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
	setup()
	defer teardown()

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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/categories/query", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"results":[{ "kind": "category_result" }]}`)
	})

	response := &CategoryQueryResponse{}
	response.Results = make([]*CategoryQueryResponseResults, 1)
	response.Results[0] = &CategoryQueryResponseResults{
		Kind: utils.String("category_result"),
	}

	input := &CategoryQueryInput{
		UsageType: utils.String("APPLIED_TO"),
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
	setup()
	defer teardown()

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
					APIVersion: utils.String("3.1"),
					Metadata: &Metadata{
						Kind: utils.String("network_security_rule"),
					},
					Spec: &NetworkSecurityRule{
						Name:        utils.String("network.create"),
						Description: utils.String("Network Create"),
						Resources:   nil,
					},
				},
			},
			&NetworkSecurityRuleIntentResponse{
				APIVersion: utils.String("3.1"),
				Metadata: &Metadata{
					Kind: utils.String("network_security_rule"),
					UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
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
	setup()
	defer teardown()

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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/network_security_rules/cfde831a-4e87-4a75-960f-89b0148aa2cc",
		func(w http.ResponseWriter, r *http.Request) {
			testHTTPMethod(t, r, http.MethodGet)
			fmt.Fprint(w, `{"metadata": {"kind":"network_security_rule","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
		})

	response := &NetworkSecurityRuleIntentResponse{}
	response.Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("network_security_rule"),
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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/network_security_rules/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"network_security_rule","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	list := &NetworkSecurityRuleListIntentResponse{}
	list.Entities = make([]*NetworkSecurityRuleIntentResource, 1)
	list.Entities[0] = &NetworkSecurityRuleIntentResource{}
	list.Entities[0].Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("network_security_rule"),
	}

	input := &DSMetadata{
		Length: utils.Int64(1.0),
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
	setup()
	defer teardown()

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
					APIVersion: utils.String("3.1"),
					Metadata: &Metadata{
						Kind: utils.String("network_security_rule"),
					},
					Spec: &NetworkSecurityRule{
						Resources:   nil,
						Description: utils.String("Network Update"),
						Name:        utils.String("network.update"),
					},
				},
			},
			&NetworkSecurityRuleIntentResponse{
				APIVersion: utils.String("3.1"),
				Metadata: &Metadata{
					Kind: utils.String("network_security_rule"),
					UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
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
	setup()
	defer teardown()

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
					APIVersion: utils.String("3.1"),
					Metadata: &Metadata{
						Kind: utils.String("volume_group"),
					},
					Spec: &VolumeGroup{
						Name: utils.String("volume.create"),
						Resources: &VolumeGroupResources{
							FlashMode: utils.String("ON"),
						},
					},
				},
			},
			&VolumeGroupResponse{
				APIVersion: utils.String("3.1"),
				Metadata: &Metadata{
					Kind: utils.String("volume_group"),
					UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
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
	setup()
	defer teardown()

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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/volume_groups/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"metadata": {"kind":"volume_group","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}`)
	})

	response := &VolumeGroupResponse{}
	response.Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("volume_group"),
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
	setup()
	defer teardown()

	mux.HandleFunc("/api/nutanix/v3/volume_groups/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"entities":[{"metadata": {"kind":"volume_group","uuid":"cfde831a-4e87-4a75-960f-89b0148aa2cc"}}]}`)
	})

	list := &VolumeGroupListResponse{}
	list.Entities = make([]*VolumeGroupResponse, 1)
	list.Entities[0] = &VolumeGroupResponse{}
	list.Entities[0].Metadata = &Metadata{
		UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
		Kind: utils.String("volume_group"),
	}

	input := &DSMetadata{
		Length: utils.Int64(1.0),
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
	setup()
	defer teardown()

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
					APIVersion: utils.String("3.1"),
					Metadata: &Metadata{
						Kind: utils.String("volume_group"),
					},
					Spec: &VolumeGroup{
						Resources: &VolumeGroupResources{
							FlashMode: utils.String("ON"),
						},
						Name: utils.String("volume.update"),
					},
				},
			},
			&VolumeGroupResponse{
				APIVersion: utils.String("3.1"),
				Metadata: &Metadata{
					Kind: utils.String("volume_group"),
					UUID: utils.String("cfde831a-4e87-4a75-960f-89b0148aa2cc"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
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
