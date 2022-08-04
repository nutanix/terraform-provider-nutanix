package Era

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

type Service interface {
	ProvisionDatabase(*ProvisionDatabaseRequest) (*ProvisionDatabaseResponse, error)
	ListProfiles() (*ListProfileResponse, error)
	ListClusters() (*ListClusterResponse, error)
	ListDatabaseTypes() (*ListDatabaseTypesResponse, error)
	ListSLA() (*ListSLAResponse, error)
	ListDatabaseParams() (*ListDatabaseParamsResponse, error)
	ListDatabaseInstances() (*ListDatabaseInstancesResponse, error)
	ListDatabaseServerVMs() (*ListDatabaseServerVMResponse, error)
	GetOperation(GetOperationRequest) (*GetOperationResponse, error)
	GetDatabaseInstance(string) (*GetDatabaseResponse, error)
	UpdateDatabase(*UpdateDatabaseRequest, string) (*UpdateDatabaseResponse, error)
	DeleteDatabase(*DeleteDatabaseRequest, string) (*DeleteDatabaseResponse, error)
}

type ServiceClient struct {
	c *client.Client
}

func (sc ServiceClient) ProvisionDatabase(req *ProvisionDatabaseRequest) (*ProvisionDatabaseResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/databases/provision", req)
	//res := new(ProvisionDatabaseResponse) // TODO: patch the response, take care of the error messages as well.
	res := new(ProvisionDatabaseResponse)

	if err != nil {
		return nil, err
	}
	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateDatabase(req *UpdateDatabaseRequest, databaseID string) (*UpdateDatabaseResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("/databases/%s", databaseID), req)
	//res := new(ProvisionDatabaseResponse) // TODO: patch the response, take care of the error messages as well.
	res := new(UpdateDatabaseResponse)

	if err != nil {
		return nil, err
	}
	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteDatabase(req *DeleteDatabaseRequest, databaseID string) (*DeleteDatabaseResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/databases/%s", databaseID), req)
	//res := new(ProvisionDatabaseResponse) // TODO: patch the response, take care of the error messages as well.
	res := new(DeleteDatabaseResponse)

	if err != nil {
		return nil, err
	}
	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListProfiles() (*ListProfileResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/profiles", nil)
	if err != nil {
		return nil, err
	}
	res := new(ListProfileResponse)

	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListClusters() (*ListClusterResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/clusters", nil)
	if err != nil {
		return nil, err
	}
	res := new(ListClusterResponse)

	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListDatabaseTypes() (*ListDatabaseTypesResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/databases/i/era-drive/tune-config", nil)
	if err != nil {
		return nil, err
	}
	res := new(ListDatabaseTypesResponse)

	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListSLA() (*ListSLAResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/slas", nil)
	if err != nil {
		return nil, err
	}
	res := new(ListSLAResponse)

	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListDatabaseParams() (*ListDatabaseParamsResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/app_types/postgres_database/provision/input-file?category=db_server;database", nil) // TODO: Check this API, is this api used to generate second page?, What is the sense of these params and do we get response of all database types ?
	if err != nil {
		return nil, err
	}
	res := new(ListDatabaseParamsResponse)

	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListDatabaseInstances() (*ListDatabaseInstancesResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/databases?detailed=true&order-by-dbserver-logical-cluster=true", nil) // TODO: Check this API, is this api used to generate second page?, What is the sense of these params and do we get response of all database types ?
	if err != nil {
		return nil, err
	}
	res := new(ListDatabaseInstancesResponse)

	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListDatabaseServerVMs() (*ListDatabaseServerVMResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/dbservers?detailed=true&load-dbserver-cluster=true", nil) // TODO: Check this API, is this api used to generate second page?, What is the sense of these params and do we get response of all database types ?
	if err != nil {
		return nil, err
	}
	res := new(ListDatabaseServerVMResponse)

	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetOperation(req GetOperationRequest) (*GetOperationResponse, error) {
	ctx := context.TODO()

	opID := req.OperationID
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/operations/%s", opID), nil) // TODO: Check this API, is this api used to generate second page?, What is the sense of these params and do we get response of all database types ?
	if err != nil {
		return nil, err
	}
	res := new(GetOperationResponse)

	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetDatabaseInstance(dbInstanceID string) (*GetDatabaseResponse, error) {
	ctx := context.TODO()

	// TODO: Use dbInstanceID in the request
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/databases/%s?detailed=true&load-dbserver-cluster=true", dbInstanceID), nil) // TODO: Check this API, is this api used to generate second page?, What is the sense of these params and do we get response of all database types ?
	if err != nil {
		return nil, err
	}
	res := new(GetDatabaseResponse)

	log.Println("Request dump in service: ")
	b, _ := httputil.DumpRequest(httpReq, true)
	log.Println(string(b))

	return res, sc.c.Do(ctx, httpReq, res)
}
