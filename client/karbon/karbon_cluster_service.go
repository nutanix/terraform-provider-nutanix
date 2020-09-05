package karbon

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// ClusterOperations ...
type ClusterOperations struct {
	client *client.Client
}

// Service ...
type ClusterService interface {
	// karbon v2.1
	ListKarbonClusters() (*KarbonClusterListIntentResponse, error)
	CreateKarbonCluster(createRequest *KarbonClusterIntentInput) (*KarbonClusterActionResponse, error)
	GetKarbonCluster(karbonClusterName string) (*KarbonClusterIntentResponse, error)
	GetKarbonClusterNodePool(karbonClusterName string, nodePoolName string) (*KarbonClusterNodePool, error)
	DeleteKarbonCluster(karbonClusterName string) (*KarbonClusterActionResponse, error)
	GetKubeConfigForKarbonCluster(karbonClusterName string) (*KarbonClusterKubeconfigResponse, error)
	GetSSHConfigForKarbonCluster(karbonClusterName string) (*KarbonClusterSSHconfig, error)
	//registries
	ListPrivateRegistries(karbonClusterName string) (*KarbonPrivateRegistryListResponse, error)
	AddPrivateRegistry(karbonClusterName string, createRequest KarbonPrivateRegistryOperationIntentInput) (*KarbonPrivateRegistryResponse, error)
	DeletePrivateRegistry(karbonClusterName string, privateRegistryName string) (*KarbonPrivateRegistryOperationResponse, error)
}

// karbon 2.1
func (op ClusterOperations) ListKarbonClusters() (*KarbonClusterListIntentResponse, error) {
	log.Printf("pre request")
	ctx := context.TODO()
	log.Printf("pre request")
	path := "/v1-beta.1/k8s/clusters"
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonClusterListIntentResponse := new(KarbonClusterListIntentResponse)
	log.Printf("post request")
	if err != nil {
		return nil, err
	}

	return karbonClusterListIntentResponse, op.client.Do(ctx, req, karbonClusterListIntentResponse)
}

func (op ClusterOperations) CreateKarbonCluster(createRequest *KarbonClusterIntentInput) (*KarbonClusterActionResponse, error) {
	ctx := context.TODO()

	path := "/v1/k8s/clusters"
	req, err := op.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	karbonClusterActionResponse := new(KarbonClusterActionResponse)

	if err != nil {
		return nil, err
	}

	return karbonClusterActionResponse, op.client.Do(ctx, req, karbonClusterActionResponse)
}

func (op ClusterOperations) GetKarbonCluster(name string) (*KarbonClusterIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1/k8s/clusters/%s", name)
	fmt.Printf("Path: %s", path)
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonClusterIntentResponse := new(KarbonClusterIntentResponse)

	if err != nil {
		return nil, err
	}

	return karbonClusterIntentResponse, op.client.Do(ctx, req, karbonClusterIntentResponse)
}

func (op ClusterOperations) GetKarbonClusterNodePool(name string, nodePoolName string) (*KarbonClusterNodePool, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1-beta.1/k8s/clusters/%s/node-pools/%s", name, nodePoolName)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonClusterNodePool := new(KarbonClusterNodePool)

	if err != nil {
		return nil, err
	}

	return karbonClusterNodePool, op.client.Do(ctx, req, karbonClusterNodePool)
}

func (op ClusterOperations) DeleteKarbonCluster(name string) (*KarbonClusterActionResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1/k8s/clusters/%s", name)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	karbonClusterActionResponse := new(KarbonClusterActionResponse)

	if err != nil {
		return nil, err
	}

	return karbonClusterActionResponse, op.client.Do(ctx, req, karbonClusterActionResponse)
}

func (op ClusterOperations) GetKubeConfigForKarbonCluster(name string) (*KarbonClusterKubeconfigResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1/k8s/clusters/%s/kubeconfig", name)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonClusterKubeconfigResponse := new(KarbonClusterKubeconfigResponse)

	if err != nil {
		return nil, err
	}

	return karbonClusterKubeconfigResponse, op.client.Do(ctx, req, karbonClusterKubeconfigResponse)
}

func (op ClusterOperations) GetSSHConfigForKarbonCluster(name string) (*KarbonClusterSSHconfig, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1/k8s/clusters/%s/ssh", name)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonClusterSSHconfig := new(KarbonClusterSSHconfig)

	if err != nil {
		return nil, err
	}

	return karbonClusterSSHconfig, op.client.Do(ctx, req, karbonClusterSSHconfig)
}

//karbon shared

func (op ClusterOperations) ListPrivateRegistries(karbonClusterName string) (*KarbonPrivateRegistryListResponse, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/v1-alpha.1/k8s/clusters/%s/registries", karbonClusterName)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonPrivateRegistryListResponse := new(KarbonPrivateRegistryListResponse)

	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryListResponse, op.client.Do(ctx, req, karbonPrivateRegistryListResponse)
}

func (op ClusterOperations) AddPrivateRegistry(karbonClusterName string, createRequest KarbonPrivateRegistryOperationIntentInput) (*KarbonPrivateRegistryResponse, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/v1-alpha.1/k8s/clusters/%s/registries", karbonClusterName)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	karbonPrivateRegistryResponse := new(KarbonPrivateRegistryResponse)

	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryResponse, op.client.Do(ctx, req, karbonPrivateRegistryResponse)
}

func (op ClusterOperations) DeletePrivateRegistry(karbonClusterName string, privateRegistryName string) (*KarbonPrivateRegistryOperationResponse, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/v1-alpha.1/k8s/clusters/%s/registries/%s", karbonClusterName, privateRegistryName)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	karbonPrivateRegistryOperationResponse := new(KarbonPrivateRegistryOperationResponse)

	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryOperationResponse, op.client.Do(ctx, req, karbonPrivateRegistryOperationResponse)
}
