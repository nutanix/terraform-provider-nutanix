package foundation

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

type NetworkingService interface {
	DiscoverNodes(context.Context) (*DiscoverNodesAPIResponse, error)
	NodeNetworkDetails(context.Context, *NodeNetworkDetailsInput) (*NodeNetworkDetailsResponse, error)
}

type NetworkingOperations struct {
	client *client.Client
}

//Discovers Nutanix-imaged nodes within an IPv6 network.
func (ntw NetworkingOperations) DiscoverNodes(ctx context.Context) (*DiscoverNodesAPIResponse, error) {
	path := "/discover_nodes"
	req, err := ntw.client.NewUnAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	discoverNodesAPIResponse := new(DiscoverNodesAPIResponse)
	return discoverNodesAPIResponse, ntw.client.Do(ctx, req, discoverNodesAPIResponse)
}

//Gets hypervisor, CVM & IPMI info of the discovered nodes using IPv6 Api
func (ntw NetworkingOperations) NodeNetworkDetails(ctx context.Context, ntwInput *NodeNetworkDetailsInput) (*NodeNetworkDetailsResponse, error) {
	path := "/node_network_details"
	req, err := ntw.client.NewUnAuthRequest(ctx, http.MethodPost, path, ntwInput)
	if err != nil {
		return nil, err
	}
	nodeNetworkDetailsResponse := new(NodeNetworkDetailsResponse)
	return nodeNetworkDetailsResponse, ntw.client.Do(ctx, req, nodeNetworkDetailsResponse)
}
