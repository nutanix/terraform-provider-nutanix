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
	CreateVM(createRequest *VMIntentInput) (*VMIntentResponse, error)
	DeleteVM(UUID string) error
	GetVM(UUID string) (*VMIntentResponse, error)
	ListVM(getEntitiesRequest *VMListMetadata) (*VMListIntentResponse, error)
	UpdateVM(UUID string, body *VMIntentInput) (*VMIntentResponse, error)
	CreateSubnet(createRequest *SubnetIntentInput) (*SubnetIntentResponse, error)
	DeleteSubnet(UUID string) error
	GetSubnet(UUID string) (*SubnetIntentResponse, error)
	ListSubnet(getEntitiesRequest *SubnetListMetadata) (*SubnetListIntentResponse, error)
	UpdateSubnet(UUID string, body *SubnetIntentInput) (*SubnetIntentResponse, error)
	CreateImage(createRequest *ImageIntentInput) (*ImageIntentResponse, error)
	DeleteImage(UUID string) error
	GetImage(UUID string) (*ImageIntentResponse, error)
	ListImage(getEntitiesRequest *ImageListMetadata) (*ImageListIntentResponse, error)
	UpdateImage(UUID string, body *ImageIntentInput) (*ImageIntentResponse, error)
}

/*CreateVM Creates a VM
 * This operation submits a request to create a VM based on the input parameters.
 *
 * @param body
 * @return *VMIntentResponse
 */
func (op Operations) CreateVM(createRequest *VMIntentInput) (*VMIntentResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/vms", createRequest)
	if err != nil {
		return nil, err
	}

	vmIntentResponse := new(VMIntentResponse)

	err = op.client.Do(ctx, req, vmIntentResponse)

	if err != nil {
		return nil, err
	}

	return vmIntentResponse, nil
}

/*DeleteVM Deletes a VM
 * This operation submits a request to delete a op.
 *
 * @param UUID The UUID of the entity.
 * @return error
 */
func (op Operations) DeleteVM(UUID string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/vms/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*GetVM Gets a VM
 * This operation gets a op.
 *
 * @param UUID The UUID of the entity.
 * @return *VMIntentResponse
 */
func (op Operations) GetVM(UUID string) (*VMIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/vms/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	vmIntentResponse := new(VMIntentResponse)

	err = op.client.Do(ctx, req, vmIntentResponse)
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
func (op Operations) ListVM(getEntitiesRequest *VMListMetadata) (*VMListIntentResponse, error) {
	ctx := context.TODO()
	path := "/vms/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}
	vmListIntentResponse := new(VMListIntentResponse)
	err = op.client.Do(ctx, req, vmListIntentResponse)
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
func (op Operations) UpdateVM(UUID string, body *VMIntentInput) (*VMIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/vms/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	vmIntentResponse := new(VMIntentResponse)

	err = op.client.Do(ctx, req, vmIntentResponse)
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
func (op Operations) CreateSubnet(createRequest *SubnetIntentInput) (*SubnetIntentResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/subnets", createRequest)
	if err != nil {
		return nil, err
	}

	subnetIntentResponse := new(SubnetIntentResponse)

	err = op.client.Do(ctx, req, subnetIntentResponse)

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
func (op Operations) DeleteSubnet(UUID string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/subnets/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*GetSubnet Gets a subnet entity
 * This operation gets a subnet.
 *
 * @param uuid The UUID of the entity.
 * @return *SubnetIntentResponse
 */
func (op Operations) GetSubnet(UUID string) (*SubnetIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/subnets/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	subnetIntentResponse := new(SubnetIntentResponse)

	err = op.client.Do(ctx, req, subnetIntentResponse)
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
func (op Operations) ListSubnet(getEntitiesRequest *SubnetListMetadata) (*SubnetListIntentResponse, error) {
	ctx := context.TODO()
	path := "/subnets/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)

	if err != nil {
		return nil, err
	}

	subnetListIntentResponse := new(SubnetListIntentResponse)
	err = op.client.Do(ctx, req, subnetListIntentResponse)

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
func (op Operations) UpdateSubnet(UUID string, body *SubnetIntentInput) (*SubnetIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/subnets/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	subnetIntentResponse := new(SubnetIntentResponse)

	err = op.client.Do(ctx, req, subnetIntentResponse)
	if err != nil {
		return nil, err
	}

	return subnetIntentResponse, nil
}

/*CreateImage Creates a IMAGE
 * Images are raw ISO, QCOW2, or VMDK files that are uploaded by a user can be attached to a op. An ISO image is attached as a virtual CD-ROM drive, and QCOW2 and VMDK files are attached as SCSI disks. An image has to be explicitly added to the self-service catalog before users can create VMs from it.
 *
 * @param body
 * @return *ImageIntentResponse
 */
func (op Operations) CreateImage(body *ImageIntentInput) (*ImageIntentResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/images", body)
	if err != nil {
		return nil, err
	}

	imageIntentResponse := new(ImageIntentResponse)

	err = op.client.Do(ctx, req, imageIntentResponse)

	if err != nil {
		return nil, err
	}

	return imageIntentResponse, nil
}

/*DeleteImage deletes a IMAGE
 * This operation submits a request to delete a IMAGE.
 *
 * @param uuid The UUID of the entity.
 * @return error if error exists
 */
func (op Operations) DeleteImage(UUID string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/images/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*GetImage gets a IMAGE
 * This operation gets a IMAGE.
 *
 * @param uuid The UUID of the entity.
 * @return *ImageIntentResponse
 */
func (op Operations) GetImage(UUID string) (*ImageIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/images/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	imageIntentResponse := new(ImageIntentResponse)

	err = op.client.Do(ctx, req, imageIntentResponse)
	if err != nil {
		return nil, err
	}

	return imageIntentResponse, nil
}

/*ListImage gets a list of IMAGEs
 * This operation gets a list of IMAGEs, allowing for sorting and pagination. Note: Entities that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest
 * @return *ImageListIntentResponse
 */
func (op Operations) ListImage(getEntitiesRequest *ImageListMetadata) (*ImageListIntentResponse, error) {
	ctx := context.TODO()
	path := "/images/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)

	if err != nil {
		return nil, err
	}

	imageListIntentResponse := new(ImageListIntentResponse)
	err = op.client.Do(ctx, req, imageListIntentResponse)

	if err != nil {
		return nil, err
	}

	return imageListIntentResponse, nil
}

/*UpdateImage updates a IMAGE
 * This operation submits a request to update a IMAGE based on the input parameters.
 *
 * @param uuid The UUID of the entity.
 * @param body
 * @return *ImageIntentResponse
 */
func (op Operations) UpdateImage(UUID string, body *ImageIntentInput) (*ImageIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/images/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	imageIntentResponse := new(ImageIntentResponse)

	err = op.client.Do(ctx, req, imageIntentResponse)
	if err != nil {
		return nil, err
	}

	return imageIntentResponse, nil
}

//TODO: Ask for images file put & get requests.
