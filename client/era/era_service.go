package era

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
	ListDatabaseTypes() (*ListDatabaseTypesResponse, error)
	ListDatabaseParams() (*ListDatabaseParamsResponse, error)
	ListDatabaseInstances() (*ListDatabaseInstancesResponse, error)
	ListDatabaseServerVMs() (*ListDatabaseServerVMResponse, error)
	GetOperation(GetOperationRequest) (*GetOperationResponse, error)
	GetDatabaseInstance(string) (*GetDatabaseResponse, error)
	UpdateDatabase(*UpdateDatabaseRequest, string) (*UpdateDatabaseResponse, error)
	DeleteDatabase(*DeleteDatabaseRequest, string) (*DeleteDatabaseResponse, error)
	ListProfiles(ctx context.Context, engine string, profileType string) (*ProfileListResponse, error)
	GetProfiles(ctx context.Context, engine string, profileType string, id string, name string) (*ListProfileResponse, error)
	GetCluster(ctx context.Context, id string, name string) (*ListClusterResponse, error)
	ListClusters(ctx context.Context) (*ClusterListResponse, error)
	GetSLA(ctx context.Context, id string, name string) (*ListSLAResponse, error)
	ListSLA(ctx context.Context) (*SLAResponse, error)
}

type ServiceClient struct {
	c *client.Client
}

func (sc ServiceClient) ListProfiles(ctx context.Context, engine string, profileType string) (*ProfileListResponse, error) {
	var httpReq *http.Request
	var err error

	path := makeListProfilePath(engine, profileType)
	httpReq, err = sc.c.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, err
	}
	res := new(ProfileListResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetProfiles(ctx context.Context, engine string, profileType string, id string, name string) (*ListProfileResponse, error) {
	var httpReq *http.Request
	var err error
	path := makePathProfiles(engine, profileType, id, name)

	httpReq, err = sc.c.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, err
	}
	res := new(ListProfileResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetCluster(ctx context.Context, id string, name string) (*ListClusterResponse, error) {
	var path string
	if id != "" {
		path = fmt.Sprintf("/clusters/%s", id)
	}
	if name != "" {
		path = fmt.Sprintf("/clusters/name/%s", name)
	}
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	res := new(ListClusterResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListClusters(ctx context.Context) (*ClusterListResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/clusters", nil)
	if err != nil {
		return nil, err
	}
	res := new(ClusterListResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetSLA(ctx context.Context, id string, name string) (*ListSLAResponse, error) {
	var path string
	if id != "" {
		path = fmt.Sprintf("/slas/%s", id)
	}
	if name != "" {
		path = fmt.Sprintf("/slas/name/%s", name)
	}
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	res := new(ListSLAResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListSLA(ctx context.Context) (*SLAResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/slas", nil)
	if err != nil {
		return nil, err
	}
	res := new(SLAResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func makeListProfilePath(engine string, profileType string) string {
	if engine != "" && profileType != "" {
		return fmt.Sprintf("/profiles?engine=%s&type=%s", engine, profileType)
	}
	if engine != "" {
		return fmt.Sprintf("/profiles?engine=%s", engine)
	} else if profileType != "" {
		return fmt.Sprintf("/profiles?type=%s", profileType)
	}
	return "/profiles"
}

func makePathProfiles(engine string, ptype string, id string, name string) string {
	if engine != "" {
		path := "/profiles?engine=" + engine
		if ptype != "" {
			path = path + "&type=" + ptype
		}
		if id != "" {
			path = path + "&id=" + id
		}
		if name != "" {
			path = path + "&name=" + name
		}
		return path
	}
	if ptype != "" {
		path := "/profiles?type=" + ptype
		if id != "" {
			path = path + "&id=" + id
		}
		if name != "" {
			path = path + "&name=" + name
		}
		return path
	}

	if id != "" {
		path := "/profiles?id=" + id
		if name != "" {
			path = path + "&name=" + name
		}
		return path
	}

	if name != "" {
		path := "/profiles?name=" + name
		return path
	}
	return ""
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
