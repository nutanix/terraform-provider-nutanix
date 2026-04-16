package era

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
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
	CreateLinkedDatabase(ctx context.Context, id string, req *CreateLinkedDatabasesRequest) (*ProvisionDatabaseResponse, error)
	DeleteLinkedDatabase(ctx context.Context, DBID string, linkedDBID string, req *DeleteLinkedDatabaseRequest) (*ProvisionDatabaseResponse, error)
	CreateMaintenanceWindow(ctx context.Context, body *MaintenanceWindowInput) (*MaintenaceWindowResponse, error)
	ReadMaintenanceWindow(ctx context.Context, id string) (*MaintenaceWindowResponse, error)
	UpdateMaintenaceWindow(ctx context.Context, body *MaintenanceWindowInput, id string) (*MaintenaceWindowResponse, error)
	DeleteMaintenanceWindow(ctx context.Context, id string) (*AuthorizeDBServerResponse, error)
	ListMaintenanceWindow(ctx context.Context) (*ListMaintenanceWindowResponse, error)
	CreateMaintenanceTask(ctx context.Context, body *MaintenanceTasksInput) (*ListMaintenanceTasksResponse, error)
	CreateTimeMachineCluster(ctx context.Context, tmsID string, body *TmsClusterIntentInput) (*TmsClusterResponse, error)
	ReadTimeMachineCluster(ctx context.Context, tmsID string, clsID string) (*TmsClusterResponse, error)
	UpdateTimeMachineCluster(ctx context.Context, tmsID string, clsID string, body *TmsClusterIntentInput) (*TmsClusterResponse, error)
	DeleteTimeMachineCluster(ctx context.Context, tmsID string, clsID string, body *DeleteTmsClusterInput) (*ProvisionDatabaseResponse, error)
	CreateTags(ctx context.Context, body *CreateTagsInput) (*TagsIntentResponse, error)
	ReadTags(ctx context.Context, id string, name string) (*GetTagsResponse, error)
	UpdateTags(ctx context.Context, body *GetTagsResponse, id string) (*GetTagsResponse, error)
	DeleteTags(ctx context.Context, id string) (*string, error)
	ListTags(ctx context.Context, entityType string) (*ListTagsResponse, error)
	CreateNetwork(ctx context.Context, body *NetworkIntentInput) (*NetworkIntentResponse, error)
	GetNetwork(ctx context.Context, id string, name string) (*NetworkIntentResponse, error)
	UpdateNetwork(ctx context.Context, body *NetworkIntentInput, id string) (*NetworkIntentResponse, error)
	DeleteNetwork(ctx context.Context, id string) (*string, error)
	ListNetwork(ctx context.Context) (*ListNetworkResponse, error)
	CreateDBServerVM(ctx context.Context, body *DBServerInputRequest) (*ProvisionDatabaseResponse, error)
	ReadDBServerVM(ctx context.Context, id string) (*DBServerVMResponse, error)
	UpdateDBServerVM(ctx context.Context, body *UpdateDBServerVMRequest, dbserverid string) (*DBServerVMResponse, error)
	DeleteDBServerVM(ctx context.Context, req *DeleteDBServerVMRequest, dbserverid string) (*DeleteDatabaseResponse, error)
	RegisterDBServerVM(ctx context.Context, body *DBServerRegisterInput) (*ProvisionDatabaseResponse, error)
	GetDBServerVM(ctx context.Context, filter *DBServerFilterRequest) (*DBServerVMResponse, error)
	ListDBServerVM(ctx context.Context) (*ListDBServerVMResponse, error)
	CreateStretchedVlan(ctx context.Context, req *StretchedVlansInput) (*StretchedVlanResponse, error)
	GetStretchedVlan(ctx context.Context, id string) (*StretchedVlanResponse, error)
	UpdateStretchedVlan(ctx context.Context, id string, req *StretchedVlansInput) (*StretchedVlanResponse, error)
	DeleteStretchedVlan(ctx context.Context, id string) (*string, error)
	RefreshClone(ctx context.Context, body *CloneRefreshInput, id string) (*ProvisionDatabaseResponse, error)
	CreateCluster(ctx context.Context, body *ClusterIntentInput) (*ProvisionDatabaseResponse, error)
	UpdateCluster(ctx context.Context, req *ClusterUpdateInput, id string) (*ListClusterResponse, error)
	DeleteCluster(ctx context.Context, req *DeleteClusterInput, id string) (*ProvisionDatabaseResponse, error)
	GetAvailableIPs(ctx context.Context, id string) (*GetNetworkAvailableIPs, error)
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
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/databases/%s?detailed=true&load-dbserver-cluster=false", dbInstanceID), nil)
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
	res := new(ProvisionDatabaseResponse)

	if err != nil {
		return nil, err
	}
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
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/tms/%s/dbservers", tmsID), req)
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

func (sc ServiceClient) CreateLinkedDatabase(ctx context.Context, id string, req *CreateLinkedDatabasesRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/databases/%s/linked-databases", id), req)
	if err != nil {
		return nil, err
	}

	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteLinkedDatabase(ctx context.Context, id string, linkDBID string, req *DeleteLinkedDatabaseRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/databases/%s/linked-databases/%s", id, linkDBID), req)
	if err != nil {
		return nil, err
	}
	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateMaintenanceWindow(ctx context.Context, body *MaintenanceWindowInput) (*MaintenaceWindowResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/maintenance", body)
	if err != nil {
		return nil, err
	}
	res := new(MaintenaceWindowResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ReadMaintenanceWindow(ctx context.Context, id string) (*MaintenaceWindowResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/maintenance/%s?load-task-associations=true", id), nil)
	if err != nil {
		return nil, err
	}
	res := new(MaintenaceWindowResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateMaintenaceWindow(ctx context.Context, body *MaintenanceWindowInput, id string) (*MaintenaceWindowResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("/maintenance/%s", id), body)
	if err != nil {
		return nil, err
	}
	res := new(MaintenaceWindowResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteMaintenanceWindow(ctx context.Context, id string) (*AuthorizeDBServerResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/maintenance/%s", id), nil)
	if err != nil {
		return nil, err
	}
	res := new(AuthorizeDBServerResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListMaintenanceWindow(ctx context.Context) (*ListMaintenanceWindowResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/maintenance?load-task-associations=true", nil)
	if err != nil {
		return nil, err
	}
	res := new(ListMaintenanceWindowResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateMaintenanceTask(ctx context.Context, req *MaintenanceTasksInput) (*ListMaintenanceTasksResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/maintenance/tasks", req)
	if err != nil {
		return nil, err
	}

	res := new(ListMaintenanceTasksResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateTags(ctx context.Context, body *CreateTagsInput) (*TagsIntentResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/tags", body)
	if err != nil {
		return nil, err
	}
	res := new(TagsIntentResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateTimeMachineCluster(ctx context.Context, tmsID string, body *TmsClusterIntentInput) (*TmsClusterResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/tms/%s/clusters", tmsID), body)
	if err != nil {
		return nil, err
	}
	res := new(TmsClusterResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ReadTags(ctx context.Context, id, name string) (*GetTagsResponse, error) {
	path := "/tags"
	if id != "" {
		path += fmt.Sprintf("?id=%s", id)
	}
	if name != "" {
		path += fmt.Sprintf("?name=%s", name)
	}

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	res := new(GetTagsResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ReadTimeMachineCluster(ctx context.Context, tmsID string, clsID string) (*TmsClusterResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/tms/%s/clusters/%s", tmsID, clsID), nil)
	if err != nil {
		return nil, err
	}

	res := new(TmsClusterResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateTimeMachineCluster(ctx context.Context, tmsID string, clsID string, body *TmsClusterIntentInput) (*TmsClusterResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("/tms/%s/clusters/%s", tmsID, clsID), body)
	if err != nil {
		return nil, err
	}

	res := new(TmsClusterResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) RefreshClone(ctx context.Context, body *CloneRefreshInput, id string) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/clones/%s/refresh", id), body)
	if err != nil {
		return nil, err
	}

	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteTimeMachineCluster(ctx context.Context, tmsID string, clsID string, body *DeleteTmsClusterInput) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/tms/%s/clusters/%s", tmsID, clsID), body)
	if err != nil {
		return nil, err
	}
	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateCluster(ctx context.Context, body *ClusterIntentInput) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/clusters", body)
	if err != nil {
		return nil, err
	}
	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateDBServerVM(ctx context.Context, body *DBServerInputRequest) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/dbservers/provision", body)
	if err != nil {
		return nil, err
	}
	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteCluster(ctx context.Context, body *DeleteClusterInput, id string) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/clusters/%s", id), body)
	if err != nil {
		return nil, err
	}
	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateTags(ctx context.Context, body *GetTagsResponse, id string) (*GetTagsResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPut, fmt.Sprintf("/tags/%s", id), body)
	if err != nil {
		return nil, err
	}
	res := new(GetTagsResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteTags(ctx context.Context, id string) (*string, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/tags/%s", id), nil)
	if err != nil {
		return nil, err
	}
	res := new(string)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateNetwork(ctx context.Context, body *NetworkIntentInput) (*NetworkIntentResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/resources/networks", body)
	if err != nil {
		return nil, err
	}

	res := new(NetworkIntentResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetNetwork(ctx context.Context, id, name string) (*NetworkIntentResponse, error) {
	path := "/resources/networks?detailed=true&"
	if name != "" {
		path += fmt.Sprintf("name=%s", name)
	} else {
		path += fmt.Sprintf("id=%s", id)
	}

	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	res := new(NetworkIntentResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateNetwork(ctx context.Context, body *NetworkIntentInput, id string) (*NetworkIntentResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPut, fmt.Sprintf("/resources/networks/%s", id), body)
	if err != nil {
		return nil, err
	}
	res := new(NetworkIntentResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteNetwork(ctx context.Context, id string) (*string, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/resources/networks/%s", id), nil)
	if err != nil {
		return nil, err
	}
	res := new(string)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListTags(ctx context.Context, entityType string) (*ListTagsResponse, error) {
	path := "/tags"
	if entityType != "" {
		path += fmt.Sprintf("?entityType=%s", entityType)
	}
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	res := new(ListTagsResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListNetwork(ctx context.Context) (*ListNetworkResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/resources/networks?detailed=true", nil)
	if err != nil {
		return nil, err
	}

	res := new(ListNetworkResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ReadDBServerVM(ctx context.Context, id string) (*DBServerVMResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/dbservers/%s", id), nil)
	if err != nil {
		return nil, err
	}
	res := new(DBServerVMResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateDBServerVM(ctx context.Context, body *UpdateDBServerVMRequest, dbServerID string) (*DBServerVMResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("/dbservers/%s", dbServerID), body)
	if err != nil {
		return nil, err
	}

	res := new(DBServerVMResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteDBServerVM(ctx context.Context, req *DeleteDBServerVMRequest, dbServerVMID string) (*DeleteDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/dbservers/%s", dbServerVMID), req)
	if err != nil {
		return nil, err
	}

	res := new(DeleteDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) RegisterDBServerVM(ctx context.Context, body *DBServerRegisterInput) (*ProvisionDatabaseResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/dbservers/register", body)
	if err != nil {
		return nil, err
	}
	res := new(ProvisionDatabaseResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetDBServerVM(ctx context.Context, filter *DBServerFilterRequest) (*DBServerVMResponse, error) {
	var httpReq *http.Request
	var err error

	path := makeDBServerPath(filter.DBServerClusterID, filter.ID, filter.IP, filter.Name, filter.NxClusterID, filter.VMClusterID, filter.VMClusterName)

	httpReq, err = sc.c.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, err
	}
	res := new(DBServerVMResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func makeDBServerPath(dbserverClsid, id, ip, name, nxClsid, vmclsid, vmclsName *string) string {
	path := "/dbservers"

	if dbserverClsid != nil {
		path += fmt.Sprintf("/%s?value-type=dbserver-cluster-id", *dbserverClsid)
	}
	if id != nil {
		path += fmt.Sprintf("/%s?value-type=id", *id)
	}
	if ip != nil {
		path += fmt.Sprintf("/%s?value-type=ip", *ip)
	}
	if name != nil {
		path += fmt.Sprintf("/%s?value-type=name", *name)
	}
	if nxClsid != nil {
		path += fmt.Sprintf("/%s?value-type=nx-cluster-id", *nxClsid)
	}
	if vmclsid != nil {
		path += fmt.Sprintf("/%s?value-type=vm-cluster-id", *vmclsid)
	}
	if vmclsName != nil {
		path += fmt.Sprintf("/%s?value-type=vm-cluster-name", *vmclsName)
	}

	path += "&load-dbserver-cluster=false&load-databases=false&load-clones=false&load-metrics=false&detailed=false&curator=false&time-zone=UTC"

	return path
}

func (sc ServiceClient) ListDBServerVM(ctx context.Context) (*ListDBServerVMResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/dbservers?load-dbserver-cluster=false&load-databases=false&load-clones=false&detailed=false&load-metrics=false&time-zone=UTC", nil)
	if err != nil {
		return nil, err
	}
	res := new(ListDBServerVMResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) CreateStretchedVlan(ctx context.Context, req *StretchedVlansInput) (*StretchedVlanResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPost, "/resources/networks/stretched-vlan", req)
	if err != nil {
		return nil, err
	}
	res := new(StretchedVlanResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetStretchedVlan(ctx context.Context, id string) (*StretchedVlanResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/resources/networks/stretched-vlan/%s", id), nil)
	if err != nil {
		return nil, err
	}
	res := new(StretchedVlanResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateStretchedVlan(ctx context.Context, id string, req *StretchedVlansInput) (*StretchedVlanResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPut, fmt.Sprintf("/resources/networks/stretched-vlan/%s", id), req)
	if err != nil {
		return nil, err
	}
	res := new(StretchedVlanResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) DeleteStretchedVlan(ctx context.Context, id string) (*string, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/resources/networks/stretched-vlan/%s", id), nil)
	if err != nil {
		return nil, err
	}
	res := new(string)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) UpdateCluster(ctx context.Context, req *ClusterUpdateInput, id string) (*ListClusterResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("/clusters/%s", id), req)
	if err != nil {
		return nil, err
	}
	res := new(ListClusterResponse)
	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetAvailableIPs(ctx context.Context, id string) (*GetNetworkAvailableIPs, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/profiles/%s/get-available-ips", id), nil)
	if err != nil {
		return nil, err
	}
	res := new(GetNetworkAvailableIPs)
	return res, sc.c.Do(ctx, httpReq, res)
}
