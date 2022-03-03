package foundation

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// Interface for file management apis of foundation
type FileManagementService interface {
	ListNOSPackages(context.Context) (*ListNOSPackagesResponse, error)
	ListHypervisorISOs(context.Context) (*ListHypervisorISOsResponse, error)
	UploadHypervisor(context.Context, *UploadHypervisorInput) (*UploadHypervisorResponse, error)
	DeleteHypervisorAOS(ctx context.Context, id string) error
}

// FileManagementOperations implements FileManagementService interface
type FileManagementOperations struct {
	client *client.Client
}

//ListNOSPackages lists the available AOS packages file names in Foundation
func (fmo FileManagementOperations) ListNOSPackages(ctx context.Context) (*ListNOSPackagesResponse, error) {
	path := "/enumerate_nos_packages"
	req, err := fmo.client.NewUnAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	listNOSPackagesResponse := new(ListNOSPackagesResponse)
	return listNOSPackagesResponse, fmo.client.Do(ctx, req, listNOSPackagesResponse)
}

//ListHypervisorISOs lists the hypervisor ISOs available in Foundation
func (fmo FileManagementOperations) ListHypervisorISOs(ctx context.Context) (*ListHypervisorISOsResponse, error) {
	path := "/enumerate_hypervisor_isos"
	req, err := fmo.client.NewUnAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	listHypervisorISOsResponse := new(ListHypervisorISOsResponse)
	return listHypervisorISOsResponse, fmo.client.Do(ctx, req, listHypervisorISOsResponse)
}

//Uploads hypervisor or AOS image to foundation.

func (fmo FileManagementOperations) UploadHypervisor(ctx context.Context, uploadAOSInput *UploadHypervisorInput) (*UploadHypervisorResponse, error) {
	path := fmt.Sprintf("/upload?filename=%s&installer_type=%s", uploadAOSInput.Filename, uploadAOSInput.Installer_type)
	req, err := fmo.client.NewUnAuthRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}
	uploadHypervisorResponse := new(UploadHypervisorResponse)
	return uploadHypervisorResponse, fmo.client.Do(ctx, req, uploadHypervisorResponse)
}

//Deletes hypervisor or AOS images uploaded to foundation.
func (fmo FileManagementOperations) DeleteHypervisorAOS(ctx context.Context, id string) error {
	path := "/delete"

	req, err := fmo.client.NewUnAuthRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return fmo.client.Do(ctx, req, nil)
}
