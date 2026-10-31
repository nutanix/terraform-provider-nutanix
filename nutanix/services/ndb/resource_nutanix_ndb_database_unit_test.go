package ndb

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

// MockEraService is a mock implementation of the era.Service interface
type MockEraService struct {
	mock.Mock
	era.Service
}

func (m *MockEraService) ReadDBServerVM(ctx context.Context, id string) (*era.DBServerVMResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*era.DBServerVMResponse), args.Error(1)
}

func (m *MockEraService) DeleteDBServerVM(ctx context.Context, req *era.DeleteDBServerVMRequest, id string) (*era.DeleteDatabaseResponse, error) {
	args := m.Called(ctx, req, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*era.DeleteDatabaseResponse), args.Error(1)
}

func (m *MockEraService) DeleteDBServerCluster(ctx context.Context, req *era.DeleteDPCRequest, id string) (*era.DeleteDatabaseResponse, error) {
	args := m.Called(ctx, req, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*era.DeleteDatabaseResponse), args.Error(1)
}

func (m *MockEraService) GetOperation(req era.GetOperationRequest) (*era.GetOperationResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*era.GetOperationResponse), args.Error(1)
}

func TestCaptureDatabaseNodeIDs(t *testing.T) {
	ctx := context.Background()
	mockService := new(MockEraService)
	conn := &era.Client{Service: mockService}

	resourceSchema := ResourceDatabaseInstance().Schema
	d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"database_nodes": []interface{}{
			map[string]interface{}{
				"dbserver_id": "vm-1",
			},
			map[string]interface{}{
				"dbserver_id": "vm-2",
			},
		},
	})

	dpcID := "dpc-123"
	mockService.On("ReadDBServerVM", ctx, "vm-1").Return(&era.DBServerVMResponse{
		DbserverClusterID: &dpcID,
	}, nil)

	vmIDs, capturedDpcID := captureDatabaseNodeIDs(ctx, d, conn)

	assert.Equal(t, []string{"vm-1", "vm-2"}, vmIDs)
	assert.Equal(t, "dpc-123", capturedDpcID)
	mockService.AssertExpectations(t)
}

func TestCleanupOrphanedNodes_DPC(t *testing.T) {
	// Set small delay for testing
	oldDelay := eraDelay
	eraDelay = 1 * time.Millisecond
	defer func() { eraDelay = oldDelay }()

	ctx := context.Background()
	mockService := new(MockEraService)
	conn := &era.Client{Service: mockService}

	dpcID := "dpc-123"
	vmIDs := []string{"vm-1", "vm-2"}
	opID := "op-456"
	statusCompleted := "5"

	// Mock DPC delete call
	mockService.On("DeleteDBServerCluster", ctx, mock.Anything, dpcID).Return(&era.DeleteDatabaseResponse{
		Operationid: opID,
	}, nil)

	// Mock polling for DPC delete operation
	mockService.On("GetOperation", era.GetOperationRequest{OperationID: opID}).Return(&era.GetOperationResponse{
		Status: &statusCompleted,
	}, nil)

	// Using a dummy ResourceData to satisfy the signature
	d := &schema.ResourceData{}

	err := cleanupOrphanedNodes(ctx, d, conn, dpcID, vmIDs)
	assert.NoError(t, err)

	// Should call DeleteDBServerCluster and NOT DeleteDBServerVM because DPC was successful
	mockService.AssertCalled(t, "DeleteDBServerCluster", ctx, mock.Anything, dpcID)
	mockService.AssertNotCalled(t, "DeleteDBServerVM", ctx, mock.Anything, mock.Anything)
}

func TestCleanupOrphanedNodes_FallbackVMs(t *testing.T) {
	// Set small delay for testing
	oldDelay := eraDelay
	eraDelay = 1 * time.Millisecond
	defer func() { eraDelay = oldDelay }()

	ctx := context.Background()
	mockService := new(MockEraService)
	conn := &era.Client{Service: mockService}

	dpcID := "" // No DPC
	vmIDs := []string{"vm-1", "vm-2"}
	opID1 := "op-1"
	opID2 := "op-2"
	statusCompleted := "5"

	mockService.On("DeleteDBServerVM", ctx, mock.Anything, "vm-1").Return(&era.DeleteDatabaseResponse{
		Operationid: opID1,
	}, nil)
	mockService.On("GetOperation", era.GetOperationRequest{OperationID: opID1}).Return(&era.GetOperationResponse{
		Status: &statusCompleted,
	}, nil)

	mockService.On("DeleteDBServerVM", ctx, mock.Anything, "vm-2").Return(&era.DeleteDatabaseResponse{
		Operationid: opID2,
	}, nil)
	mockService.On("GetOperation", era.GetOperationRequest{OperationID: opID2}).Return(&era.GetOperationResponse{
		Status: &statusCompleted,
	}, nil)

	d := &schema.ResourceData{}

	err := cleanupOrphanedNodes(ctx, d, conn, dpcID, vmIDs)
	assert.NoError(t, err)

	mockService.AssertCalled(t, "DeleteDBServerVM", ctx, mock.Anything, "vm-1")
	mockService.AssertCalled(t, "DeleteDBServerVM", ctx, mock.Anything, "vm-2")
}

func TestEraRefresh_Status17(t *testing.T) {
	ctx := context.Background()
	mockService := new(MockEraService)
	conn := &era.Client{Service: mockService}

	opID := "op-17"
	statusWarning := "17"
	msg := "Warning: Some nodes on same host"

	mockService.On("GetOperation", era.GetOperationRequest{OperationID: opID}).Return(&era.GetOperationResponse{
		Status:  &statusWarning,
		Message: &msg,
	}, nil)

	refreshFunc := eraRefresh(ctx, conn, era.GetOperationRequest{OperationID: opID})
	res, state, err := refreshFunc()

	assert.NoError(t, err)
	assert.Equal(t, "COMPLETED", state)
	assert.NotNil(t, res)
	mockService.AssertExpectations(t)
}
