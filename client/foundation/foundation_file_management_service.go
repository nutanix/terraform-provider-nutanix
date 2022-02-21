package foundation

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

type FileManagementService interface {
	ListNOSPackages() (*ListNOSPackagesResponse, error)
	ListHypervisorISOs() (*ListHypervisorISOsResponse, error)
}

type FileManagementOperations struct {
	client *client.Client
}

//Lists the available AOS packages in Foundation
func (fileManagementOperations FileManagementOperations) ListNOSPackages() (*ListNOSPackagesResponse, error) {
	ctx := context.TODO()
	path := "/enumerate_nos_packages"
	req, err := fileManagementOperations.client.NewUnAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	listNOSPackagesResponse := new(ListNOSPackagesResponse)
	return listNOSPackagesResponse, fileManagementOperations.client.Do(ctx, req, listNOSPackagesResponse)
}

//Lists the hypervisor ISOs available in Foundation
func (fileManagementOperations FileManagementOperations) ListHypervisorISOs() (*ListHypervisorISOsResponse, error) {
	ctx := context.TODO()
	path := "/enumerate_hypervisor_isos"
	req, err := fileManagementOperations.client.NewUnAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	listHypervisorISOsResponse := new(ListHypervisorISOsResponse)
	return listHypervisorISOsResponse, fileManagementOperations.client.Do(ctx, req, listHypervisorISOsResponse)
}
