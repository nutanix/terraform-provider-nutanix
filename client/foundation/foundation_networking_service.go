package foundation

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// NetworkingService is a interface for networking apis in foundation
type NetworkingService interface {
	DiscoverNodes(context.Context) (*DiscoverNodesAPIResponse, error)
}

// NetworkingOperations implements NetworkingService interface
type NetworkingOperations struct {
	client *client.Client
}

// DiscoverNodes discovers(gets) Nutanix-imaged nodes within an IPv6 network.
func (ntw NetworkingOperations) DiscoverNodes(ctx context.Context) (*DiscoverNodesAPIResponse, error) {
	path := "/discover_nodes"
	req, err := ntw.client.NewUnAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	discoverNodesAPIResponse := new(DiscoverNodesAPIResponse)
	return discoverNodesAPIResponse, ntw.client.Do(ctx, req, discoverNodesAPIResponse)
}
