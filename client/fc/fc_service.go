package fc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// Operations ...
type Operations struct {
	client *client.Client
}

// Service ...
type Service interface {
	GetImagedNode(context.Context, string) (*ImagedNodeDetails, error)
	ListImagedNodes(context.Context, *ImagedNodesListInput) (*ImagedNodesListResponse, error)
	GetImagedCluster(context.Context, string) (*ImagedClusterDetails, error)
	ListImagedClusters(context.Context, *ImagedClustersListInput) (*ImagedClustersListResponse, error)
	CreateCluster(context.Context, *CreateClusterInput) (*CreateClusterResponse, error)
	UpdateCluster(context.Context, string, *UpdateClusterData) error
	DeleteCluster(context.Context, string) error
	CreateAPIKey(context.Context, *CreateAPIKeysInput) (*CreateAPIKeysResponse, error)
	GetAPIKey(context.Context, string) (*CreateAPIKeysResponse, error)
	ListAPIKeys(context.Context, *ListMetadataInput) (*ListAPIKeysResponse, error)
}

func (op Operations) GetImagedNode(ctx context.Context, nodeUUID string) (*ImagedNodeDetails, error) {
	path := fmt.Sprintf("/imaged_nodes/%s", nodeUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	imagedNodeDetails := new(ImagedNodeDetails)

	if err != nil {
		return nil, err
	}

	return imagedNodeDetails, op.client.Do(ctx, req, imagedNodeDetails)
}

func (op Operations) ListImagedNodes(ctx context.Context, input *ImagedNodesListInput) (*ImagedNodesListResponse, error) {
	path := "/imaged_nodes/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	imagedNodesListResponse := new(ImagedNodesListResponse)

	if err != nil {
		return nil, err
	}

	return imagedNodesListResponse, op.client.Do(ctx, req, imagedNodesListResponse)
}

func (op Operations) GetImagedCluster(ctx context.Context, clusterUUID string) (*ImagedClusterDetails, error) {
	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	imagedClusterDetails := new(ImagedClusterDetails)

	if err != nil {
		return nil, err
	}

	return imagedClusterDetails, op.client.Do(ctx, req, imagedClusterDetails)
}

func (op Operations) ListImagedClusters(ctx context.Context, input *ImagedClustersListInput) (*ImagedClustersListResponse, error) {
	path := "/imaged_clusters/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	imagedClustersListResponse := new(ImagedClustersListResponse)

	if err != nil {
		return nil, err
	}

	return imagedClustersListResponse, op.client.Do(ctx, req, imagedClustersListResponse)
}

func (op Operations) CreateCluster(ctx context.Context, input *CreateClusterInput) (*CreateClusterResponse, error) {
	path := "/imaged_clusters"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	createClusterResponse := new(CreateClusterResponse)

	if err != nil {
		return nil, err
	}

	return createClusterResponse, op.client.Do(ctx, req, createClusterResponse)
}

func (op Operations) UpdateCluster(ctx context.Context, clusterUUID string, updateData *UpdateClusterData) error {
	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, updateData)

	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

func (op Operations) DeleteCluster(ctx context.Context, clusterUUID string) error {
	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)

	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

//Create a new api key which will be used by remote nodes to authenticate with Foundation Central
func (op Operations) CreateAPIKey(ctx context.Context, input *CreateAPIKeysInput) (*CreateAPIKeysResponse, error) {
	path := "/api_keys"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)
	if err != nil {
		return nil, err
	}

	createAPIResponse := new(CreateAPIKeysResponse)
	return createAPIResponse, op.client.Do(ctx, req, createAPIResponse)
}

//Get an api key given its UUID.
func (op Operations) GetAPIKey(ctx context.Context, uuid string) (*CreateAPIKeysResponse, error) {
	path := fmt.Sprintf("/api_keys/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, uuid)
	if err != nil {
		return nil, err
	}

	getAPIResponse := new(CreateAPIKeysResponse)
	return getAPIResponse, op.client.Do(ctx, req, getAPIResponse)
}

//List all the api keys.
func (op Operations) ListAPIKeys(ctx context.Context, body *ListMetadataInput) (*ListAPIKeysResponse, error) {
	path := "/api_keys/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	listAPIKeysResponse := new(ListAPIKeysResponse)
	return listAPIKeysResponse, op.client.Do(ctx, req, listAPIKeysResponse)
}
