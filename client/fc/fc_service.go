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
	// GetImagedCluster()
	// ListImagedCluster()
	// CreateCluster()
	// UpdateCluster()
	// DeleteCluster()
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
