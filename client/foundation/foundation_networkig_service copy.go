package foundation

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

type NetworkingService interface {
	DiscoverNodes(context.Context) (*DiscoverNodesAPIResponse, error)
}

type NetworkingOperations struct {
	client *client.Client
}

//Discovers Nutanix-imaged nodes within an IPv6 network.
func (ntw NetworkingOperations) DiscoverNodes(ctx context.Context) (*DiscoverNodesAPIResponse, error) {
	path := "/dicover_nodes"
	req, err := ntw.client.NewUnAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	discoverNodesAPIResponse := new(DiscoverNodesAPIResponse)
	return discoverNodesAPIResponse, ntw.client.Do(ctx, req, discoverNodesAPIResponse)
}
