package v3

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

//Operations ...
type Operations struct {
	client *client.Client
}

// Service ...
type Service interface {
	CreateVM(createRequest VMIntentInput) (*VMIntentResponse, error)
	DeleteVM(UUID string) error
	GetVM(UUID string) (*VMIntentResponse, error)
	ListVM(getEntitiesRequest VMListMetadata) (*VMListIntentResponse, error)
	UpdateVM(UUID string, body VMIntentInput) (*VMIntentResponse, error)
	CreateSubnet(createRequest SubnetIntentInput) (*SubnetIntentResponse, error)
	DeleteSubnet(UUID string) error
	GetSubnet(UUID string) (*SubnetIntentResponse, error)
	ListSubnet(getEntitiesRequest SubnetListMetadata) (*SubnetListIntentResponse, error)
	UpdateSubnet(UUID string, body SubnetIntentInput) (*SubnetIntentResponse, error)
}

/*CreateVM Creates a VM
 * This operation submits a request to create a VM based on the input parameters.
 *
 * @param body
 * @return *VMIntentResponse
 */
func (vm Operations) CreateVM(createRequest VMIntentInput) (*VMIntentResponse, error) {
	ctx := context.TODO()

	req, err := vm.client.NewRequest(ctx, http.MethodPost, "/vms", createRequest)
	if err != nil {
		return nil, err
	}

	vmIntentResponse := new(VMIntentResponse)

	err = vm.client.Do(ctx, req, vmIntentResponse)

	if err != nil {
		return nil, err
	}

	return vmIntentResponse, nil
}

/*DeleteVM Deletes a VM
 * This operation submits a request to delete a VM.
 *
 * @param UUID The UUID of the entity.
 * @return error
 */
func (vm Operations) DeleteVM(UUID string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/vms/%s", UUID)

	req, err := vm.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return vm.client.Do(ctx, req, nil)
}

/*GetVM Gets a VM
 * This operation gets a VM.
 *
 * @param UUID The UUID of the entity.
 * @return *VMIntentResponse
 */
func (vm Operations) GetVM(UUID string) (*VMIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/vms/%s", UUID)

	req, err := vm.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	vmIntentResponse := new(VMIntentResponse)

	err = vm.client.Do(ctx, req, vmIntentResponse)
	if err != nil {
		return nil, err
	}

	return vmIntentResponse, nil
}

/*ListVM Get a list of VMs
 * This operation gets a list of VMs, allowing for sorting and pagination. Note: Entities that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest
 * @return *VmListIntentResponse
 */
func (vm Operations) ListVM(getEntitiesRequest VMListMetadata) (*VMListIntentResponse, error) {
	ctx := context.TODO()
	path := "/vms/list"

	req, err := vm.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}
	vmListIntentResponse := new(VMListIntentResponse)
	err = vm.client.Do(ctx, req, vmListIntentResponse)
	if err != nil {
		return nil, err
	}

	return vmListIntentResponse, nil
}

/*UpdateVM Updates a VM
 * This operation submits a request to update a VM based on the input parameters.
 *
 * @param uuid The UUID of the entity.
 * @param body
 * @return *VMIntentResponse
 */
func (vm Operations) UpdateVM(UUID string, body VMIntentInput) (*VMIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/vms/%s", UUID)

	req, err := vm.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	vmIntentResponse := new(VMIntentResponse)

	err = vm.client.Do(ctx, req, vmIntentResponse)
	if err != nil {
		return nil, err
	}

	return vmIntentResponse, nil
}

/*CreateSubnet Creates a subnet
 * This operation submits a request to create a subnet based on the input parameters. A subnet is a block of IP addresses.
 *
 * @param body
 * @return *SubnetIntentResponse
 */
func (vm Operations) CreateSubnet(createRequest SubnetIntentInput) (*SubnetIntentResponse, error) {
	ctx := context.TODO()

	req, err := vm.client.NewRequest(ctx, http.MethodPost, "/subnets", createRequest)
	if err != nil {
		return nil, err
	}

	subnetIntentResponse := new(SubnetIntentResponse)

	err = vm.client.Do(ctx, req, subnetIntentResponse)

	if err != nil {
		return nil, err
	}

	return subnetIntentResponse, nil
}

/*DeleteSubnet Deletes a subnet
 * This operation submits a request to delete a subnet.
 *
 * @param uuid The UUID of the entity.
 * @return error if exist error
 */
func (vm Operations) DeleteSubnet(UUID string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/subnets/%s", UUID)

	req, err := vm.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return vm.client.Do(ctx, req, nil)
}

/*GetSubnet Gets a subnet entity
 * This operation gets a subnet.
 *
 * @param uuid The UUID of the entity.
 * @return *SubnetIntentResponse
 */
func (vm Operations) GetSubnet(UUID string) (*SubnetIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/subnets/%s", UUID)

	req, err := vm.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	subnetIntentResponse := new(SubnetIntentResponse)

	err = vm.client.Do(ctx, req, subnetIntentResponse)
	if err != nil {
		return nil, err
	}

	return subnetIntentResponse, nil
}

/*ListSubnet Gets a list of subnets
 * This operation gets a list of subnets, allowing for sorting and pagination. Note: Entities that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest
 * @return *SubnetListIntentResponse
 */
func (vm Operations) ListSubnet(getEntitiesRequest SubnetListMetadata) (*SubnetListIntentResponse, error) {
	ctx := context.TODO()
	path := "/subnets/list"

	req, err := vm.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)

	if err != nil {
		return nil, err
	}

	subnetListIntentResponse := new(SubnetListIntentResponse)
	err = vm.client.Do(ctx, req, subnetListIntentResponse)

	if err != nil {
		return nil, err
	}

	return subnetListIntentResponse, nil
}

/*UpdateSubnet Updates a subnet
 * This operation submits a request to update a subnet based on the input parameters.
 *
 * @param uuid The UUID of the entity.
 * @param body
 * @return *SubnetIntentResponse
 */
func (vm Operations) UpdateSubnet(UUID string, body SubnetIntentInput) (*SubnetIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/subnets/%s", UUID)

	req, err := vm.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	subnetIntentResponse := new(SubnetIntentResponse)

	err = vm.client.Do(ctx, req, subnetIntentResponse)
	if err != nil {
		return nil, err
	}

	return subnetIntentResponse, nil
}
