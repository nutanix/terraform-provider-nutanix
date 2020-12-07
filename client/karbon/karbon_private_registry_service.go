package karbon

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// PrivateRegistryOperations ...
type PrivateRegistryOperations struct {
	client *client.Client
}

// Service ...
type PrivateRegistryService interface {
	// karbon v2.1
	ListKarbonPrivateRegistries() (*KarbonPrivateRegistryListResponse, error)
	CreateKarbonPrivateRegistry(createRequest *KarbonPrivateRegistryIntentInput) (*KarbonPrivateRegistryResponse, error)
	GetKarbonPrivateRegistry(name string) (*KarbonPrivateRegistryResponse, error)
	DeleteKarbonPrivateRegistry(name string) (*KarbonPrivateRegistryOperationResponse, error)
}

func (op PrivateRegistryOperations) ListKarbonPrivateRegistries() (*KarbonPrivateRegistryListResponse, error) {
	ctx := context.TODO()
	path := "/v1-alpha.1/registries"
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonPrivateRegistryListResponse := new(KarbonPrivateRegistryListResponse)
	log.Printf("post request")
	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryListResponse, op.client.Do(ctx, req, karbonPrivateRegistryListResponse)
}

func (op PrivateRegistryOperations) CreateKarbonPrivateRegistry(createRequest *KarbonPrivateRegistryIntentInput) (*KarbonPrivateRegistryResponse, error) {
	ctx := context.TODO()
	path := "/v1-alpha.1/registries"
	req, err := op.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	karbonPrivateRegistryResponse := new(KarbonPrivateRegistryResponse)
	if err != nil {
		return nil, err
	}
	return karbonPrivateRegistryResponse, op.client.Do(ctx, req, karbonPrivateRegistryResponse)
}

func (op PrivateRegistryOperations) GetKarbonPrivateRegistry(name string) (*KarbonPrivateRegistryResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1-alpha.1/registries/%s", name)
	fmt.Printf("Path: %s", path)
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonPrivateRegistryResponse := new(KarbonPrivateRegistryResponse)

	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryResponse, op.client.Do(ctx, req, karbonPrivateRegistryResponse)
}

func (op PrivateRegistryOperations) DeleteKarbonPrivateRegistry(name string) (*KarbonPrivateRegistryOperationResponse, error) {
	ctx := context.TODO()
	log.Printf("[Debug] Deleting /v1-alpha.1/registries/%s", name)
	path := fmt.Sprintf("/v1-alpha.1/registries/%s", name)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	karbonPrivateRegistryOperationResponse := new(KarbonPrivateRegistryOperationResponse)

	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryOperationResponse, op.client.Do(ctx, req, karbonPrivateRegistryOperationResponse)
}
