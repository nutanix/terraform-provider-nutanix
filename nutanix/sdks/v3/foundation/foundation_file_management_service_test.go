package foundation

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestFMOperations_ListNOSPackages(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()
	mux.HandleFunc("/foundation/enumerate_nos_packages", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)

		// mock response
		fmt.Fprintf(w, `[
			"package1",
			"package2"
		]`)
	})
	ctx := context.TODO()

	out := &ListNOSPackagesResponse{
		"package1",
		"package2",
	}

	op := FileManagementOperations{
		client: c,
	}

	// checks
	got, err := op.ListNOSPackages(ctx)
	if err != nil {
		t.Fatalf("FileManagementOperations.ListNOSPackages() error = %v", err)
	}
	if !reflect.DeepEqual(got, out) {
		t.Errorf("FileManagementOperations.ListNOSPackages() got = %#v, want = %#v", got, out)
	}
}

func TestFMOperations_ListHypervisorISOs(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()
	mux.HandleFunc("/foundation/enumerate_hypervisor_isos", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)

		// mock response
		fmt.Fprintf(w, `{
			"hyperv": [{
					"filename": "hyperv1.iso",
					"supported": true
				},
				{
					"filename": "hyperv2.iso",
					"supported": false
				}
			],
			"kvm": [{
					"filename": "kvm1.iso",
					"supported": true
				},
				{
					"filename": "kvm2.iso",
					"supported": false
				}
			]
		}`)
	})
	ctx := context.TODO()

	out := &ListHypervisorISOsResponse{
		Hyperv: []*HypervisorISOReference{
			{
				Supported: utils.BoolPtr(true),
				Filename:  "hyperv1.iso",
			},
			{
				Supported: utils.BoolPtr(false),
				Filename:  "hyperv2.iso",
			},
		},
		Kvm: []*HypervisorISOReference{
			{
				Supported: utils.BoolPtr(true),
				Filename:  "kvm1.iso",
			},
			{
				Supported: utils.BoolPtr(false),
				Filename:  "kvm2.iso",
			},
		},
	}

	op := FileManagementOperations{
		client: c,
	}

	// checks
	got, err := op.ListHypervisorISOs(ctx)
	if err != nil {
		t.Fatalf("FileManagementOperations.ListHypervisorISOs() error = %v", err)
	}
	if !reflect.DeepEqual(got, out) {
		t.Errorf("FileManagementOperations.ListHypervisorISOs() got = %#v, want = %#v", got, out)
	}
}

func TestFMOperations_UploadImage(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()
	installerType := "kvm"
	filename := "test_ahv.iso"
	source := "foundation_api.go"
	mux.HandleFunc("/foundation/upload", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expectedURL := fmt.Sprintf("/foundation/upload?installer_type=%v&filename=%v", installerType, filename)
		if expectedURL != r.URL.String() {
			t.Errorf("FileManagementOperations.UploadImage() expected URL %v, got %v", expectedURL, r.URL.String())
		}

		body, _ := ioutil.ReadAll(r.Body)
		file, _ := ioutil.ReadFile(source)

		if !reflect.DeepEqual(body, file) {
			t.Errorf("FileManagementOperations.UploadImage() error: different uploaded files")
		}

		// mock response
		fmt.Fprintf(w, `{
			"md5sum": "1234QAA",
			"name": "/home/foundation/kvm/%v",
			"in_whitelist": false
		  }`, filename)
	})
	ctx := context.TODO()

	out := &UploadImageResponse{
		Md5Sum:      "1234QAA",
		Name:        "/home/foundation/kvm/" + filename,
		InWhitelist: false,
	}

	op := FileManagementOperations{
		client: c,
	}

	// checks
	got, err := op.UploadImage(ctx, installerType, filename, source)
	if err != nil {
		t.Fatalf("FileManagementOperations.UploadImage() error = %v", err)
	}
	if !reflect.DeepEqual(got, out) {
		t.Errorf("FileManagementOperations.UploadImage() got = %#v, want = %#v", got, out)
	}
}

func TestFMOperations_DeleteImage(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()
	installerType := "kvm"
	filename := "test_ahv.iso"
	mux.HandleFunc("/foundation/delete/", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("FileManagementOperations.DeleteImage() error reading request body = %v", err)
		}

		// check form encoded body
		expected := fmt.Sprintf("filename=%v&installer_type=%v", filename, installerType)
		if string(body) != expected {
			t.Errorf("FileManagementOperations.DeleteImage() request body expected = %v, got = %v", expected, string(body))
		}
	})
	ctx := context.TODO()

	op := FileManagementOperations{
		client: c,
	}

	// checks
	err := op.DeleteImage(ctx, installerType, filename)
	if err != nil {
		t.Fatalf("FileManagementOperations.DeleteImage() error = %v", err)
	}
}
