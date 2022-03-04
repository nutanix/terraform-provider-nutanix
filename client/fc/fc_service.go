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
	GetImagedNode(uuid string) (*ImagedNodeDetails, error)
	ListImagedNodes(req *ImagedNodesListInput) (*ImagedNodesListResponse, error)
	GetImagedCluster(uuid string) (*ImagedClusterDetails, error)
	ListImagedClusters(input *ImagedClustersListInput) (*ImagedClustersListResponse, error)
	CreateCluster(input *CreateClusterInput) (*CreateClusterResponse, error)
	UpdateCluster(clusterUUID string, updateData *UpdateClusterData) error
	DeleteCluster(clusterUUID string) error
	// CreateAPIKey()
	// GetAPIKey()
	// ListAPIKeys()
}

func (op Operations) GetImagedNode(nodeUUID string) (*ImagedNodeDetails, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/imaged_nodes/%s", nodeUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	imagedNodeDetails := new(ImagedNodeDetails)

	if err != nil {
		return nil, err
	}

	return imagedNodeDetails, op.client.Do(ctx, req, imagedNodeDetails)
}

func (op Operations) ListImagedNodes(input *ImagedNodesListInput) (*ImagedNodesListResponse, error) {
	ctx := context.TODO()
	path := "/imaged_nodes/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	imagedNodesListResponse := new(ImagedNodesListResponse)

	if err != nil {
		return nil, err
	}

	return imagedNodesListResponse, op.client.Do(ctx, req, imagedNodesListResponse)
}

func (op Operations) GetImagedCluster(clusterUUID string) (*ImagedClusterDetails, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	imagedClusterDetails := new(ImagedClusterDetails)

	if err != nil {
		return nil, err
	}

	return imagedClusterDetails, op.client.Do(ctx, req, imagedClusterDetails)
}

func (op Operations) ListImagedClusters(input *ImagedClustersListInput) (*ImagedClustersListResponse, error) {
	ctx := context.TODO()
	path := "/imaged_clusters/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	imagedClustersListResponse := new(ImagedClustersListResponse)

	if err != nil {
		return nil, err
	}

	return imagedClustersListResponse, op.client.Do(ctx, req, imagedClustersListResponse)
}

func (op Operations) CreateCluster(input *CreateClusterInput) (*CreateClusterResponse, error) {
	ctx := context.TODO()
	path := "/imaged_clusters"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	createClusterResponse := new(CreateClusterResponse)

	if err != nil {
		return nil, err
	}

	return createClusterResponse, op.client.Do(ctx, req, createClusterResponse)
}

func (op Operations) UpdateCluster(clusterUUID string, updateData *UpdateClusterData) error {
	ctx := context.TODO()
	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, updateData)

	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

func (op Operations) DeleteCluster(clusterUUID string) error {
	ctx := context.TODO()
	path := fmt.Sprintf("/imaged_clusters/%s", clusterUUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)

	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}
