package fc

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Operations ...
type Operations struct {
	client *client.Client
}

// Service ...
type Service interface {
	GetVersion()
	GetImagedCluster()
	ListImagedCluster()
	CreateCluster()
	UpdateCluster()
	DeleteCluster()
	GetImagedNode()
	ListImagedNodes()
	CreateAPIKey()
	ListAPIKeys()
	GetAPIKey()
}