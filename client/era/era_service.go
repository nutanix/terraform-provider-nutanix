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
	GetProfile(ctx context.Context, filters *ProfileFilter) (*ListProfileResponse, error)
	CreateProfiles(ctx context.Context, req *ProfileRequest) (*ListProfileResponse, error)
	DeleteProfile(ctx context.Context, uuid string) (*string, error)
	GetCluster(ctx context.Context, id string, name string) (*ListClusterResponse, error)
	ListClusters(ctx context.Context) (*ClusterListResponse, error)
	GetSLA(ctx context.Context, id string, name string) (*ListSLAResponse, error)
	ListSLA(ctx context.Context) (*SLAResponse, error)
	CreateSLA(ctx context.Context, req *SLAIntentInput) (*ListSLAResponse, error)
	UpdateSLA(ctx context.Context, req *SLAIntentInput, id string) (*ListSLAResponse, error)
	DeleteSLA(ctx context.Context, uuid string) (*SLADeleteResponse, error)
	DatabaseRestore(ctx context.Context, databaseID string, req *DatabaseRestoreRequest) (*ProvisionDatabaseResponse, error)
	LogCatchUp(ctx context.Context, id string, req *LogCatchUpRequest) (*ProvisionDatabaseResponse, error)
	CreateSoftwareProfiles(ctx context.Context, req *ProfileRequest) (*SoftwareProfileResponse, error)
	UpdateProfile(ctx context.Context, req *UpdateProfileRequest, id string) (*ListProfileResponse, error)
	GetSoftwareProfileVersion(ctx context.Context, profileID string, profileVersionID string) (*Versions, error)
	CreateSoftwareProfileVersion(ctx context.Context, id string, req *ProfileRequest) (*SoftwareProfileResponse, error)
	UpdateProfileVersion(ctx context.Context, req *ProfileRequest, id string, vid string) (*ListProfileResponse, error)
	DeleteProfileVersion(ctx context.Context, profileID string, profileVersionID string) (*string, error)
	DatabaseScale(ctx context.Context, id string, req *DatabaseScale) (*ProvisionDatabaseResponse, error)
	RegisterDatabase(ctx context.Context, request *RegisterDBInputRequest) (*ProvisionDatabaseResponse, error)
	GetTimeMachine(ctx context.Context, tmsID string, tmsName string) (*TimeMachine, error)
	ListTimeMachines(ctx context.Context) (*ListTimeMachines, error)
	DatabaseSnapshot(ctx context.Context, id string, req *DatabaseSnapshotRequest) (*ProvisionDatabaseResponse, error)
	UpdateSnapshot(ctx context.Context, id string, req *UpdateSnapshotRequest) (*SnapshotResponse, error)
	GetSnapshot(ctx context.Context, id string, filter *FilterParams) (*SnapshotResponse, error)
	DeleteSnapshot(ctx context.Context, id string) (*ProvisionDatabaseResponse, error)
	ListSnapshots(ctx context.Context, tmsID string) (*ListSnapshots, error)
	CreateClone(ctx context.Context, id string, req *CloneRequest) (*ProvisionDatabaseResponse, error)
	UpdateCloneDatabase(ctx context.Context, id string, req *UpdateDatabaseRequest) (*UpdateDatabaseResponse, error)
	GetClone(ctx context.Context, id string, name string, filterParams *FilterParams) (*GetDatabaseResponse, error)
	ListClones(ctx context.Context, filter *FilterParams) (*ListDatabaseInstance, error)
	DeleteClone(ctx context.Context, id string, req *DeleteDatabaseRequest) (*ProvisionDatabaseResponse, error)
	AuthorizeDBServer(ctx context.Context, id string, req []*string) (*AuthorizeDBServerResponse, error)
	DeAuthorizeDBServer(ctx context.Context, id string, req []*string) (*AuthorizeDBServerResponse, error)
	TimeMachineCapability(ctx context.Context, tmsID string) (*TimeMachineCapability, error)
	AddRemoveDatabase(ctx context.Context, id string, req *AddRemoveDatabasesRequest)
	CreateLinkedDatabase(ctx context.Context, id string, req *CreateLinkedDatabasesRequest) (*ProvisionDatabaseResponse, error)
	DeleteLinkedDatabase(ctx context.Context, dbId string, linkeddbId string, req *DeleteLinkedDatabaseRequest) (*ProvisionDatabaseResponse, error)
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

func (sc ServiceClient) GetProfile(ctx context.Context, filter *ProfileFilter) (*ListProfileResponse, error) {
	var httpReq *http.Request
	var err error

	path := makePathProfiles(filter.Engine, filter.ProfileType, filter.ProfileID, filter.ProfileName)

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
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/databases/%s?detailed=false&load-dbserver-cluster=false", dbInstanceID), nil)
	if err != nil {
		return nil, err
	}
	res := new(GetDatabaseResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListDatabaseInstance(ctx context.Context) (*ListDatabaseInstance, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, ("/databases?detailed=true&load-dbserver-cluster=true&order-by-dbserver-cluster=false"), nil)
	if err != nil {
		return nil, err
	}
	res := new(ListDatabaseInstance)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateSLA(ctx context.Context, req *SLAIntentInput) (*ListSLAResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/slas", req)
	res := new(ListSLAResponse)
	if err != nil {
		return nil, err
	}

	return res, sc.c.Do(ctx, httpReq, res)
}
func (sc ServiceClient) CreateProfiles(ctx context.Context, req *ProfileRequest) (*ListProfileResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/profiles", req)
	res := new(ListProfileResponse)

	if err != nil {
		return nil, err
	}

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) RegisterDatabase(ctx context.Context, req *RegisterDBInputRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/databases/register", req)
	res := new(ProvisionDatabaseResponse)

	if err != nil {
		return nil, err
	}

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteSLA(ctx context.Context, uuid string) (*SLADeleteResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/slas/%s", uuid), nil)
	if err != nil {
		return nil, err
	}
	res := new(SLADeleteResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}
func (sc ServiceClient) DeleteProfile(ctx context.Context, uuid string) (*string, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/profiles/%s", uuid), nil)
	if err != nil {
		return nil, err
	}
	res := new(string)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateSLA(ctx context.Context, req *SLAIntentInput, id string) (*ListSLAResponse, error) {
	path := fmt.Sprintf("/slas/%s", id)
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return nil, err
	}
	res := new(ListSLAResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}
func (sc ServiceClient) UpdateProfile(ctx context.Context, req *UpdateProfileRequest, id string) (*ListProfileResponse, error) {
	path := fmt.Sprintf("/profiles/%s", id)
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return nil, err
	}
	res := new(ListProfileResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DatabaseRestore(ctx context.Context, databaseID string, req *DatabaseRestoreRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/databases/%s/restore", databaseID), req)
	if err != nil {
		return nil, err
	}

	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}
func (sc ServiceClient) CreateLinkedDatabase(ctx context.Context, id string, req *CreateLinkedDatabasesRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/databases/%s/linked-databases", id), req)
	if err != nil {
		return nil, err
	}

	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}
func (sc ServiceClient) DatabaseSnapshot(ctx context.Context, id string, req *DatabaseSnapshotRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/tms/%s/snapshots", id), req)
	if err != nil {
		return nil, err
	}

	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) LogCatchUp(ctx context.Context, tmsID string, req *LogCatchUpRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/tms/%s/log-catchups", tmsID), req)
	if err != nil {
		return nil, err
	}

	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteLinkedDatabase(ctx context.Context, id string, linkDbId string, req *DeleteLinkedDatabaseRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/databases/%s/linked-databases/%s", id, linkDbId), req)

	if err != nil {
		return nil, err
	}
	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DatabaseScale(ctx context.Context, databaseID string, req *DatabaseScale) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/databases/%s/update/extend-storage", databaseID), req)
	res := new(ProvisionDatabaseResponse)

	if err != nil {
		return nil, err
	}
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateProfileVersion(ctx context.Context, req *ProfileRequest, id string, vid string) (*ListProfileResponse, error) {
	path := fmt.Sprintf("/profiles/%s/versions/%s", id, vid)
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPut, path, req)
	if err != nil {
		return nil, err
	}
	res := new(ListProfileResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateSoftwareProfiles(ctx context.Context, req *ProfileRequest) (*SoftwareProfileResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/profiles", req)
	res := new(SoftwareProfileResponse)

	if err != nil {
		return nil, err
	}

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetSoftwareProfileVersion(ctx context.Context, profileID string, profileVersionID string) (*Versions, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/profiles/%s/versions/%s", profileID, profileVersionID), nil)
	res := new(Versions)

	if err != nil {
		return nil, err
	}

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateSoftwareProfileVersion(ctx context.Context, id string, req *ProfileRequest) (*SoftwareProfileResponse, error) {
	path := fmt.Sprintf("/profiles/%s/versions", id)
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, err
	}
	res := new(SoftwareProfileResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteProfileVersion(ctx context.Context, profileID string, profileVersionID string) (*string, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/profiles/%s/versions/%s", profileID, profileVersionID), nil)
	if err != nil {
		return nil, err
	}
	res := new(string)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateSnapshot(ctx context.Context, snapshotID string, req *UpdateSnapshotRequest) (*SnapshotResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("/snapshots/i/%s", snapshotID), req)
	if err != nil {
		return nil, err
	}

	res := new(SnapshotResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteSnapshot(ctx context.Context, snapshotID string) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/snapshots/%s", snapshotID), nil)
	if err != nil {
		return nil, err
	}

	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetSnapshot(ctx context.Context, snapshotID string, filter *FilterParams) (*SnapshotResponse, error) {
	path := fmt.Sprintf("/snapshots/%s", snapshotID)
	if filter != nil {
		path = path + "?load-replicated-child-snapshots=" + filter.LoadReplicatedChildSnapshots + "&time-zone=" + filter.TimeZone
	}
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	res := new(SnapshotResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListSnapshots(ctx context.Context, tmsID string) (*ListSnapshots, error) {
	path := ("/snapshots?all=false&time-zone=UTC")
	if tmsID != "" {
		path = path + "&value-type=time-machine&value=" + tmsID
	}
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	res := new(ListSnapshots)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetTimeMachine(ctx context.Context, tmsID string, tmsName string) (*TimeMachine, error) {
	path := ""

	if len(tmsName) > 0 {
		path = fmt.Sprintf("/tms/%s?value-type=name&detailed=false&load-database=false&load-clones=false&time-zone=UTC", tmsName)
	} else {
		path = fmt.Sprintf("/tms/%s?value-type=id&detailed=false&load-database=false&load-clones=false&time-zone=UTC", tmsID)
	}
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	res := new(TimeMachine)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListTimeMachines(ctx context.Context) (*ListTimeMachines, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/tms", nil)
	if err != nil {
		return nil, err
	}

	res := new(ListTimeMachines)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateClone(ctx context.Context, id string, req *CloneRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/tms/%s/clones", id), req)
	if err != nil {
		return nil, err
	}

	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetClone(ctx context.Context, id string, name string, filter *FilterParams) (*GetDatabaseResponse, error) {
	path := ""

	if name != "" {
		path = fmt.Sprintf("/clones/%s?value-type=name&detailed=%s&any-status=%s&load-dbserver-cluster=%s&time-zone=%s", name, filter.Detailed, filter.AnyStatus, filter.LoadDBServerCluster, filter.TimeZone)
	} else {
		path = fmt.Sprintf("/clones/%s?value-type=id&detailed=%s&any-status=%s&load-dbserver-cluster=%s&time-zone=%s", id, filter.Detailed, filter.AnyStatus, filter.LoadDBServerCluster, filter.TimeZone)
	}
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	res := new(GetDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListClones(ctx context.Context, filter *FilterParams) (*ListDatabaseInstance, error) {
	path := fmt.Sprintf("/clones?detailed=%s&any-status=%s&load-dbserver-cluster=%s&order-by-dbserver-cluster=%s&order-by-dbserver-logical-cluster=%s&time-zone=%s",
		filter.Detailed, filter.AnyStatus, filter.LoadDBServerCluster, filter.OrderByDBServerCluster, filter.OrderByDBServerLogicalCluster, filter.TimeZone)
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	res := new(ListDatabaseInstance)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateCloneDatabase(ctx context.Context, id string, req *UpdateDatabaseRequest) (*UpdateDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("/clones/%s", id), req)
	res := new(UpdateDatabaseResponse)

	if err != nil {
		return nil, err
	}

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteClone(ctx context.Context, cloneID string, req *DeleteDatabaseRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/clones/%s", cloneID), req)
	if err != nil {
		return nil, err
	}

	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) AuthorizeDBServer(ctx context.Context, tmsID string, req []*string) (*AuthorizeDBServerResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("/tms/%s/dbservers", tmsID), req)
	if err != nil {
		return nil, err
	}

	res := new(AuthorizeDBServerResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeAuthorizeDBServer(ctx context.Context, tmsID string, req []*string) (*AuthorizeDBServerResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/tms/%s/dbservers", tmsID), req)
	if err != nil {
		return nil, err
	}

	res := new(AuthorizeDBServerResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) TimeMachineCapability(ctx context.Context, tmsID string) (*TimeMachineCapability, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/tms/%s/capability?time-zone=UTC&type=detailed&load-db-logs=true&load-snapshots=true", tmsID), "")
	if err != nil {
		return nil, err
	}

	res := new(TimeMachineCapability)
	return res, sc.c.Do(ctx, httpReq, res)
}
