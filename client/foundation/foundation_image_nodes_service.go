package foundation

import (
	"context"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

type ImageNodesService interface {
	ImageNodes(*ImageNodesInput) (*ImageNodesAPIResponse, error)
	NodeImagingProgress(string) (*ProgressResponse, error)
}

type ImageNodesOperations struct {
	client *client.Client
}

func (imageNodesOperations ImageNodesOperations) ImageNodes(imageNodeRequest *ImageNodesInput) (*ImageNodesAPIResponse, error) {
	ctx := context.TODO()
	path := "image_nodes"
	req, err := imageNodesOperations.client.NewUnAuthRequest(ctx, http.MethodPost, path, imageNodeRequest)
	if err != nil {
		return nil, err
	}

	//Check for invalid requests
	imageNodesAPIResponse := new(ImageNodesAPIResponse)
	errd := imageNodesOperations.client.Do(ctx, req, imageNodesAPIResponse)
	return imageNodesAPIResponse, errd
}

//Gets progress of imaging session.
func (imageNodesOperations ImageNodesOperations) NodeImagingProgress(session_id string) (*ProgressResponse, error) {
	ctx := context.TODO()
	path := "/progress?session_id=" + session_id
	req, err := imageNodesOperations.client.NewUnAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	progressResponse := new(ProgressResponse)
	return progressResponse, imageNodesOperations.client.Do(ctx, req, progressResponse)
}
