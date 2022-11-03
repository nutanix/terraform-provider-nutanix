package era

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

type Service interface {
	ProvisionDatabase(ctx context.Context, req *ProvisionDatabaseRequest) (*ProvisionDatabaseResponse, error)
	ListDatabaseTypes() (*ListDatabaseTypesResponse, error)
	ListDatabaseParams() (*ListDatabaseParamsResponse, error)
	ListDatabaseServerVMs() (*ListDatabaseServerVMResponse, error)
	GetOperation(GetOperationRequest) (*GetOperationResponse, error)
	GetDatabaseInstance(ctx context.Context, uuid string) (*GetDatabaseResponse, error)
	ListDatabaseInstance(ctx context.Context) (*ListDatabaseInstance, error)
	UpdateDatabase(ctx context.Context, req *UpdateDatabaseRequest, uuid string) (*UpdateDatabaseResponse, error)
	DeleteDatabase(ctx context.Context, req *DeleteDatabaseRequest, uuid string) (*DeleteDatabaseResponse, error)
	ListProfiles(ctx context.Context, engine string, profileType string) (*ProfileListResponse, error)
	GetProfiles(ctx context.Context, engine string, profileType string, id string, name string) (*ListProfileResponse, error)
	GetCluster(ctx context.Context, id string, name string) (*ListClusterResponse, error)
	ListClusters(ctx context.Context) (*ClusterListResponse, error)
	GetSLA(ctx context.Context, id string, name string) (*ListSLAResponse, error)
	ListSLA(ctx context.Context) (*SLAResponse, error)
	CreateSLA(ctx context.Context, req *SlaIntentInput) (*ListSLAResponse, error)
	UpdateSLA(ctx context.Context, req *SlaIntentInput, id string) (*ListSLAResponse, error)
	DeleteSLA(ctx context.Context, uuid string) (*SlaDeleteResponse, error)
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

func (sc ServiceClient) ProvisionDatabase(ctx context.Context, req *ProvisionDatabaseRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/databases/provision", req)
	res := new(ProvisionDatabaseResponse)

	if err != nil {
		return nil, err
	}

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateDatabase(ctx context.Context, req *UpdateDatabaseRequest, databaseID string) (*UpdateDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("/databases/%s", databaseID), req)
	res := new(UpdateDatabaseResponse)

	if err != nil {
		return nil, err
	}

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteDatabase(ctx context.Context, req *DeleteDatabaseRequest, databaseID string) (*DeleteDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/databases/%s", databaseID), req)
	res := new(DeleteDatabaseResponse)

	if err != nil {
		return nil, err
	}

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListDatabaseTypes() (*ListDatabaseTypesResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/databases/i/era-drive/tune-config", nil)
	if err != nil {
		return nil, err
	}
	res := new(ListDatabaseTypesResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListDatabaseParams() (*ListDatabaseParamsResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/app_types/postgres_database/provision/input-file?category=db_server;database", nil)
	if err != nil {
		return nil, err
	}
	res := new(ListDatabaseParamsResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListDatabaseServerVMs() (*ListDatabaseServerVMResponse, error) {
	ctx := context.TODO()

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/dbservers?detailed=true&load-dbserver-cluster=true", nil)
	if err != nil {
		return nil, err
	}
	res := new(ListDatabaseServerVMResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetOperation(req GetOperationRequest) (*GetOperationResponse, error) {
	ctx := context.TODO()

	opID := req.OperationID
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/operations/%s", opID), nil)
	if err != nil {
		return nil, err
	}
	res := new(GetOperationResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetDatabaseInstance(ctx context.Context, dbInstanceID string) (*GetDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/databases/%s?detailed=true&load-dbserver-cluster=true", dbInstanceID), nil)
	if err != nil {
		return nil, err
	}
	res := new(GetDatabaseResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListDatabaseInstance(ctx context.Context) (*ListDatabaseInstance, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, ("/databases?detailed=false&load-dbserver-cluster=false&order-by-dbserver-cluster=false"), nil)
	if err != nil {
		return nil, err
	}
	res := new(ListDatabaseInstance)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateSLA(ctx context.Context, req *SlaIntentInput) (*ListSLAResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/slas", req)
	res := new(ListSLAResponse)

	if err != nil {
		return nil, err
	}

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteSLA(ctx context.Context, uuid string) (*SlaDeleteResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/slas/%s", uuid), nil)
	if err != nil {
		return nil, err
	}
	res := new(SlaDeleteResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateSLA(ctx context.Context, req *SlaIntentInput, id string) (*ListSLAResponse, error) {
	path := fmt.Sprintf("/slas/%s", id)
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return nil, err
	}
	res := new(ListSLAResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}
