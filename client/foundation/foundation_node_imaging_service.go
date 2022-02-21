package foundation

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

type NodeImagingService interface {
	ImageNodes(*ImageNodesInput) (*ImageNodesAPIResponse, error)
	ImageNodesProgress(string) (*ImageNodesProgressResponse, error)
}

type NodeImagingOperations struct {
	client *client.Client
}

func (nodeImagingOperations NodeImagingOperations) ImageNodes(imageNodeInput *ImageNodesInput) (*ImageNodesAPIResponse, error) {
	ctx := context.TODO()
	path := "/image_nodes"
	req, err := nodeImagingOperations.client.NewUnAuthRequest(ctx, http.MethodPost, path, imageNodeInput)
	if err != nil {
		return nil, err
	}

	imageNodesAPIResponse := new(ImageNodesAPIResponse)
	return imageNodesAPIResponse, nodeImagingOperations.client.Do(ctx, req, imageNodesAPIResponse)
}

//Gets progress of imaging session.
func (nodeImagingOperations NodeImagingOperations) ImageNodesProgress(session_id string) (*ImageNodesProgressResponse, error) {
	ctx := context.TODO()
	path := "/progress?session_id=" + session_id
	req, err := nodeImagingOperations.client.NewUnAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	imageNodesProgressResponse := new(ImageNodesProgressResponse)
	return imageNodesProgressResponse, nodeImagingOperations.client.Do(ctx, req, imageNodesProgressResponse)
}
