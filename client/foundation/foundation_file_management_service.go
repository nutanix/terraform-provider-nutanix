package foundation

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// Interface for file management apis of foundation
type FileManagementService interface {
	ListNOSPackages(context.Context) (*ListNOSPackagesResponse, error)
	ListHypervisorISOs(context.Context) (*ListHypervisorISOsResponse, error)
}

// FileManagementOperations implements FileManagementService interface
type FileManagementOperations struct {
	client *client.Client
}

//ListNOSPackages lists the available AOS packages in Foundation
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
