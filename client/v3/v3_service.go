package v3

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Operations ...
type Operations struct {
	client *client.Client
}

// Service ...
type Service interface {
	CreateVM(createRequest *VMIntentInput) (*VMIntentResponse, error)
	DeleteVM(UUID string) (*DeleteResponse, error)
	GetVM(UUID string) (*VMIntentResponse, error)
	ListVM(getEntitiesRequest *DSMetadata) (*VMListIntentResponse, error)
	UpdateVM(UUID string, body *VMIntentInput) (*VMIntentResponse, error)
	CreateSubnet(createRequest *SubnetIntentInput) (*SubnetIntentResponse, error)
	DeleteSubnet(UUID string) (*DeleteResponse, error)
	GetSubnet(UUID string) (*SubnetIntentResponse, error)
	ListSubnet(getEntitiesRequest *DSMetadata) (*SubnetListIntentResponse, error)
	UpdateSubnet(UUID string, body *SubnetIntentInput) (*SubnetIntentResponse, error)
	CreateImage(createRequest *ImageIntentInput) (*ImageIntentResponse, error)
	DeleteImage(UUID string) (*DeleteResponse, error)
	GetImage(UUID string) (*ImageIntentResponse, error)
	ListImage(getEntitiesRequest *DSMetadata) (*ImageListIntentResponse, error)
	UpdateImage(UUID string, body *ImageIntentInput) (*ImageIntentResponse, error)
	UploadImage(UUID, filepath string) error
	CreateOrUpdateCategoryKey(body *CategoryKey) (*CategoryKeyStatus, error)
	ListCategories(getEntitiesRequest *CategoryListMetadata) (*CategoryKeyListResponse, error)
	DeleteCategoryKey(name string) error
	GetCategoryKey(name string) (*CategoryKeyStatus, error)
	ListCategoryValues(name string, getEntitiesRequest *CategoryListMetadata) (*CategoryValueListResponse, error)
	CreateOrUpdateCategoryValue(name string, body *CategoryValue) (*CategoryValueStatus, error)
	GetCategoryValue(name string, value string) (*CategoryValueStatus, error)
	DeleteCategoryValue(name string, value string) error
	GetCategoryQuery(query *CategoryQueryInput) (*CategoryQueryResponse, error)
	UpdateNetworkSecurityRule(UUID string, body *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error)
	ListNetworkSecurityRule(getEntitiesRequest *DSMetadata) (*NetworkSecurityRuleListIntentResponse, error)
	GetNetworkSecurityRule(UUID string) (*NetworkSecurityRuleIntentResponse, error)
	DeleteNetworkSecurityRule(UUID string) (*DeleteResponse, error)
	CreateNetworkSecurityRule(request *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error)
	ListCluster(getEntitiesRequest *DSMetadata) (*ClusterListIntentResponse, error)
	GetCluster(UUID string) (*ClusterIntentResponse, error)
	UpdateVolumeGroup(UUID string, body *VolumeGroupInput) (*VolumeGroupResponse, error)
	ListVolumeGroup(getEntitiesRequest *DSMetadata) (*VolumeGroupListResponse, error)
	GetVolumeGroup(UUID string) (*VolumeGroupResponse, error)
	DeleteVolumeGroup(UUID string) error
	CreateVolumeGroup(request *VolumeGroupInput) (*VolumeGroupResponse, error)
	ListAllVM() (*VMListIntentResponse, error)
	ListAllSubnet() (*SubnetListIntentResponse, error)
	ListAllNetworkSecurityRule() (*NetworkSecurityRuleListIntentResponse, error)
	ListAllImage() (*ImageListIntentResponse, error)
	ListAllCluster() (*ClusterListIntentResponse, error)
	GetTask(TaskUUID string) (*TasksResponse, error)
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
	vmIntentResponse := new(VMIntentResponse)

	if err != nil {
		return nil, err
	}

	return vmIntentResponse, op.client.Do(ctx, req, vmIntentResponse)
}

/*DeleteVM Deletes a VM
 * This operation submits a request to delete a op.
 *
 * @param UUID The UUID of the entity.
 * @return error
 */
func (op Operations) DeleteVM(UUID string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/vms/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)
	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
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
	vmIntentResponse := new(VMIntentResponse)

	if err != nil {
		return nil, err
	}

	return vmIntentResponse, op.client.Do(ctx, req, vmIntentResponse)
}

/*ListVM Get a list of VMs This operation gets a list of VMs, allowing for sorting and pagination. Note: Entities that have not been created
 * successfully are not listed.
 *
 * @param getEntitiesRequest @return *VmListIntentResponse
 */
func (op Operations) ListVM(getEntitiesRequest *DSMetadata) (*VMListIntentResponse, error) {
	ctx := context.TODO()
	path := "/vms/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	vmListIntentResponse := new(VMListIntentResponse)

	if err != nil {
		return nil, err
	}

	return vmListIntentResponse, op.client.Do(ctx, req, vmListIntentResponse)
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
	vmIntentResponse := new(VMIntentResponse)

	if err != nil {
		return nil, err
	}

	return vmIntentResponse, op.client.Do(ctx, req, vmIntentResponse)
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
	subnetIntentResponse := new(SubnetIntentResponse)

	if err != nil {
		return nil, err
	}

	return subnetIntentResponse, op.client.Do(ctx, req, subnetIntentResponse)
}

/*DeleteSubnet Deletes a subnet
 * This operation submits a request to delete a subnet.
 *
 * @param uuid The UUID of the entity.
 * @return error if exist error
 */
func (op Operations) DeleteSubnet(UUID string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/subnets/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)
	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
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
	subnetIntentResponse := new(SubnetIntentResponse)

	if err != nil {
		return nil, err
	}

	// Recheck subnet already exist error
	// if *subnetIntentResponse.Status.State == "ERROR" {
	// 	pretty, _ := json.MarshalIndent(subnetIntentResponse.Status.MessageList, "", "  ")
	// 	return nil, fmt.Errorf("error: %s", string(pretty))
	// }

	return subnetIntentResponse, op.client.Do(ctx, req, subnetIntentResponse)
}

/*ListSubnet Gets a list of subnets This operation gets a list of subnets, allowing for sorting and pagination. Note: Entities that have not
 * been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *SubnetListIntentResponse
 */
func (op Operations) ListSubnet(getEntitiesRequest *DSMetadata) (*SubnetListIntentResponse, error) {
	ctx := context.TODO()
	path := "/subnets/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	subnetListIntentResponse := new(SubnetListIntentResponse)

	if err != nil {
		return nil, err
	}

	return subnetListIntentResponse, op.client.Do(ctx, req, subnetListIntentResponse)
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
	subnetIntentResponse := new(SubnetIntentResponse)

	if err != nil {
		return nil, err
	}

	return subnetIntentResponse, op.client.Do(ctx, req, subnetIntentResponse)
}

/*CreateImage Creates a IMAGE Images are raw ISO, QCOW2, or VMDK files that are uploaded by a user can be attached to a op. An ISO image is
 * attached as a virtual CD-ROM drive, and QCOW2 and VMDK files are attached as SCSI disks. An image has to be explicitly added to the
 * self-service catalog before users can create VMs from it.
 *
 * @param body @return *ImageIntentResponse
 */
func (op Operations) CreateImage(body *ImageIntentInput) (*ImageIntentResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/images", body)
	imageIntentResponse := new(ImageIntentResponse)

	if err != nil {
		return nil, err
	}

	return imageIntentResponse, op.client.Do(ctx, req, imageIntentResponse)
}

/*UploadImage Uplloads a Image Binary file Images are raw ISO, QCOW2, or VMDK files that are uploaded by a user can be attached to a op. An
 * ISO image is attached as a virtual CD-ROM drive, and QCOW2 and VMDK files are attached as SCSI disks. An image has to be explicitly added
 * to the self-service catalog before users can create VMs from it.
 *
 * @param UUID @param filepath
 */
func (op Operations) UploadImage(UUID, filepath string) error {
	ctx := context.Background()

	path := fmt.Sprintf("/images/%s/file", UUID)

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error: cannot open file: %s", err)
	}
	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error: Cannot read file %s", err)
	}

	req, err := op.client.NewUploadRequest(ctx, http.MethodPut, path, fileContents)
	if err != nil {
		return fmt.Errorf("error: Creating request %s", err)
	}
	err = op.client.Do(ctx, req, nil)

	return err

}

/*DeleteImage deletes a IMAGE
 * This operation submits a request to delete a IMAGE.
 *
 * @param uuid The UUID of the entity.
 * @return error if error exists
 */
func (op Operations) DeleteImage(UUID string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/images/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)
	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
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
	imageIntentResponse := new(ImageIntentResponse)

	if err != nil {
		return nil, err
	}
	return imageIntentResponse, op.client.Do(ctx, req, imageIntentResponse)
}

/*ListImage gets a list of IMAGEs This operation gets a list of IMAGEs, allowing for sorting and pagination. Note: Entities that have not
 * been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *ImageListIntentResponse
 */
func (op Operations) ListImage(getEntitiesRequest *DSMetadata) (*ImageListIntentResponse, error) {
	ctx := context.TODO()
	path := "/images/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	imageListIntentResponse := new(ImageListIntentResponse)

	if err != nil {
		return nil, err
	}

	return imageListIntentResponse, op.client.Do(ctx, req, imageListIntentResponse)
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
	imageIntentResponse := new(ImageIntentResponse)

	if err != nil {
		return nil, err
	}

	return imageIntentResponse, op.client.Do(ctx, req, imageIntentResponse)
}

/*GetCluster gets a CLUSTER
 * This operation gets a CLUSTER.
 *
 * @param uuid The UUID of the entity.
 * @return *ImageIntentResponse
 */
func (op Operations) GetCluster(UUID string) (*ClusterIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/clusters/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	clusterIntentResponse := new(ClusterIntentResponse)

	if err != nil {
		return nil, err
	}

	return clusterIntentResponse, op.client.Do(ctx, req, clusterIntentResponse)
}

/*ListCluster gets a list of CLUSTERS This operation gets a list of CLUSTERS, allowing for sorting and pagination. Note: Entities that have
 * not been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *ClusterListIntentResponse
 */
func (op Operations) ListCluster(getEntitiesRequest *DSMetadata) (*ClusterListIntentResponse, error) {
	ctx := context.TODO()
	path := "/clusters/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	clusterList := new(ClusterListIntentResponse)

	if err != nil {
		return nil, err
	}

	return clusterList, op.client.Do(ctx, req, clusterList)
}

/*UpdateImage updates a CLUSTER
 * This operation submits a request to update a CLUSTER based on the input parameters.
 *
 * @param uuid The UUID of the entity.
 * @param body
 * @return *ImageIntentResponse
 */
// func (op Operations) UpdateImage(UUID string, body *ImageIntentInput) (*ImageIntentResponse, error) {
// 	ctx := context.TODO()

// 	path := fmt.Sprintf("/images/%s", UUID)

// 	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	imageIntentResponse := new(ImageIntentResponse)

// 	err = op.client.Do(ctx, req, imageIntentResponse)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return imageIntentResponse, nil
// }

// CreateOrUpdateCategoryKey ...
func (op Operations) CreateOrUpdateCategoryKey(body *CategoryKey) (*CategoryKeyStatus, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s", utils.StringValue(body.Name))
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	categoryKeyResponse := new(CategoryKeyStatus)

	if err != nil {
		return nil, err
	}

	return categoryKeyResponse, op.client.Do(ctx, req, categoryKeyResponse)
}

/*ListCategories gets a list of Categories This operation gets a list of Categories, allowing for sorting and pagination. Note: Entities
 * that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *ImageListIntentResponse
 */
func (op Operations) ListCategories(getEntitiesRequest *CategoryListMetadata) (*CategoryKeyListResponse, error) {
	ctx := context.TODO()
	path := "/categories/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	categoryKeyListResponse := new(CategoryKeyListResponse)

	if err != nil {
		return nil, err
	}

	return categoryKeyListResponse, op.client.Do(ctx, req, categoryKeyListResponse)
}

/*DeleteCategoryKey Deletes a Category
 * This operation submits a request to delete a op.
 *
 * @param name The name of the entity.
 * @return error
 */
func (op Operations) DeleteCategoryKey(name string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s", name)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*GetCategoryKey gets a Category
 * This operation gets a Category.
 *
 * @param name The name of the entity.
 * @return *CategoryKeyStatus
 */
func (op Operations) GetCategoryKey(name string) (*CategoryKeyStatus, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s", name)
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	categoryKeyStatusResponse := new(CategoryKeyStatus)

	if err != nil {
		return nil, err
	}

	return categoryKeyStatusResponse, op.client.Do(ctx, req, categoryKeyStatusResponse)
}

/*ListCategoryValues gets a list of Category values for a specific key This operation gets a list of Categories, allowing for sorting and
 * pagination. Note: Entities that have not been created successfully are not listed.
 *
 * @param name @param getEntitiesRequest @return *CategoryValueListResponse
 */
func (op Operations) ListCategoryValues(name string, getEntitiesRequest *CategoryListMetadata) (*CategoryValueListResponse, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/categories/%s/list", name)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	categoryValueListResponse := new(CategoryValueListResponse)

	if err != nil {
		return nil, err
	}

	return categoryValueListResponse, op.client.Do(ctx, req, categoryValueListResponse)
}

// CreateOrUpdateCategoryValue ...
func (op Operations) CreateOrUpdateCategoryValue(name string, body *CategoryValue) (*CategoryValueStatus, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s/%s", name, utils.StringValue(body.Value))
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	categoryValueResponse := new(CategoryValueStatus)

	if err != nil {
		return nil, err
	}

	return categoryValueResponse, op.client.Do(ctx, req, categoryValueResponse)
}

/*GetCategoryValue gets a Category Value
 * This operation gets a Category Value.
 *
 * @param name The name of the entity.
 * @params value the value of entity that belongs to category key
 * @return *CategoryValueStatus
 */
func (op Operations) GetCategoryValue(name string, value string) (*CategoryValueStatus, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s/%s", name, value)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	categoryValueStatusResponse := new(CategoryValueStatus)

	if err != nil {
		return nil, err
	}

	return categoryValueStatusResponse, op.client.Do(ctx, req, categoryValueStatusResponse)
}

/*DeleteCategoryValue Deletes a Category Value
 * This operation submits a request to delete a op.
 *
 * @param name The name of the entity.
 * @params value the value of entity that belongs to category key
 * @return error
 */
func (op Operations) DeleteCategoryValue(name string, value string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s/%s", name, value)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*GetCategoryQuery gets list of entities attached to categories or policies in which categories are used as defined by the filter criteria.
 *
 * @param query Categories query input object.
 * @return *CategoryQueryResponse
 */
func (op Operations) GetCategoryQuery(query *CategoryQueryInput) (*CategoryQueryResponse, error) {
	ctx := context.TODO()

	path := "/categories/query"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, query)
	categoryQueryResponse := new(CategoryQueryResponse)

	if err != nil {
		return nil, err
	}

	return categoryQueryResponse, op.client.Do(ctx, req, categoryQueryResponse)
}

/*CreateNetworkSecurityRule Creates a Network security rule
 * This operation submits a request to create a Network security rule based on the input parameters.
 *
 * @param request
 * @return *NetworkSecurityRuleIntentResponse
 */
func (op Operations) CreateNetworkSecurityRule(request *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error) {
	ctx := context.TODO()

	networkSecurityRuleIntentResponse := new(NetworkSecurityRuleIntentResponse)
	req, err := op.client.NewRequest(ctx, http.MethodPost, "/network_security_rules", request)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleIntentResponse, op.client.Do(ctx, req, networkSecurityRuleIntentResponse)
}

/*DeleteNetworkSecurityRule Deletes a Network security rule
 * This operation submits a request to delete a Network security rule.
 *
 * @param UUID The UUID of the entity.
 * @return void
 */
func (op Operations) DeleteNetworkSecurityRule(UUID string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/network_security_rules/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)
	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

/*GetNetworkSecurityRule Gets a Network security rule
 * This operation gets a Network security rule.
 *
 * @param UUID The UUID of the entity.
 * @return *NetworkSecurityRuleIntentResponse
 */
func (op Operations) GetNetworkSecurityRule(UUID string) (*NetworkSecurityRuleIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/network_security_rules/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	networkSecurityRuleIntentResponse := new(NetworkSecurityRuleIntentResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleIntentResponse, op.client.Do(ctx, req, networkSecurityRuleIntentResponse)
}

/*ListNetworkSecurityRule Gets all network security rules This operation gets a list of Network security rules, allowing for sorting and
 * pagination. Note: Entities that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *NetworkSecurityRuleListIntentResponse
 */
func (op Operations) ListNetworkSecurityRule(getEntitiesRequest *DSMetadata) (*NetworkSecurityRuleListIntentResponse, error) {
	ctx := context.TODO()
	path := "/network_security_rules/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	networkSecurityRuleListIntentResponse := new(NetworkSecurityRuleListIntentResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleListIntentResponse, op.client.Do(ctx, req, networkSecurityRuleListIntentResponse)
}

/*UpdateNetworkSecurityRule Updates a Network security rule
 * This operation submits a request to update a Network security rule based on the input parameters.
 *
 * @param uuid The UUID of the entity.
 * @param body
 * @return void
 */
func (op Operations) UpdateNetworkSecurityRule(
	UUID string,
	body *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/network_security_rules/%s", UUID)
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	networkSecurityRuleIntentResponse := new(NetworkSecurityRuleIntentResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleIntentResponse, op.client.Do(ctx, req, networkSecurityRuleIntentResponse)
}

/*CreateVolumeGroup Creates a Volume group
 * This operation submits a request to create a Volume group based on the input parameters.
 *
 * @param request
 * @return *VolumeGroupResponse
 */
func (op Operations) CreateVolumeGroup(request *VolumeGroupInput) (*VolumeGroupResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/volume_groups", request)
	networkSecurityRuleResponse := new(VolumeGroupResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleResponse, op.client.Do(ctx, req, networkSecurityRuleResponse)
}

/*DeleteVolumeGroup Deletes a Volume group
 * This operation submits a request to delete a Volume group.
 *
 * @param UUID The UUID of the entity.
 * @return void
 */
func (op Operations) DeleteVolumeGroup(UUID string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/volume_groups/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*GetVolumeGroup Gets a Volume group
 * This operation gets a Volume group.
 *
 * @param UUID The UUID of the entity.
 * @return *VolumeGroupResponse
 */
func (op Operations) GetVolumeGroup(UUID string) (*VolumeGroupResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/volume_groups/%s", UUID)
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	networkSecurityRuleResponse := new(VolumeGroupResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleResponse, op.client.Do(ctx, req, networkSecurityRuleResponse)
}

/*ListVolumeGroup Gets all network security rules This operation gets a list of Volume groups, allowing for sorting and pagination. Note:
 * Entities that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *VolumeGroupListResponse
 */
func (op Operations) ListVolumeGroup(getEntitiesRequest *DSMetadata) (*VolumeGroupListResponse, error) {
	ctx := context.TODO()
	path := "/volume_groups/list"
	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	networkSecurityRuleListResponse := new(VolumeGroupListResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleListResponse, op.client.Do(ctx, req, networkSecurityRuleListResponse)
}

/*UpdateVolumeGroup Updates a Volume group
 * This operation submits a request to update a Volume group based on the input parameters.
 *
 * @param uuid The UUID of the entity.
 * @param body
 * @return void
 */
func (op Operations) UpdateVolumeGroup(UUID string, body *VolumeGroupInput) (*VolumeGroupResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/volume_groups/%s", UUID)
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	networkSecurityRuleResponse := new(VolumeGroupResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleResponse, op.client.Do(ctx, req, networkSecurityRuleResponse)
}

const itemsPerPage = int64(100)

func hasNext(ri *int64) bool {
	*ri = *ri - itemsPerPage
	return *ri >= (0 - itemsPerPage)
}

// ListAllVM ...
func (op Operations) ListAllVM() (*VMListIntentResponse, error) {
	entities := make([]*VMIntentResource, 0)

	resp, err := op.ListVM(&DSMetadata{
		Kind:   utils.String("vm"),
		Length: utils.Int64(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListVM(&DSMetadata{
				Kind:   utils.String("vm"),
				Length: utils.Int64(itemsPerPage),
				Offset: utils.Int64(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
		}
		resp.Entities = entities
	}

	return resp, nil
}

// ListAllSubnet ...
func (op Operations) ListAllSubnet() (*SubnetListIntentResponse, error) {
	entities := make([]*SubnetIntentResponse, 0)

	resp, err := op.ListSubnet(&DSMetadata{
		Kind:   utils.String("subnet"),
		Length: utils.Int64(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListSubnet(&DSMetadata{
				Kind:   utils.String("subnet"),
				Length: utils.Int64(itemsPerPage),
				Offset: utils.Int64(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}
		resp.Entities = entities
	}

	return resp, nil
}

// ListAllNetworkSecurityRule ...
func (op Operations) ListAllNetworkSecurityRule() (*NetworkSecurityRuleListIntentResponse, error) {
	entities := make([]*NetworkSecurityRuleIntentResource, 0)

	resp, err := op.ListNetworkSecurityRule(&DSMetadata{
		Kind:   utils.String("network_security_rule"),
		Length: utils.Int64(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListNetworkSecurityRule(&DSMetadata{
				Kind:   utils.String("network_security_rule"),
				Length: utils.Int64(itemsPerPage),
				Offset: utils.Int64(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}
		resp.Entities = entities
	}

	return resp, nil
}

// ListAllImage ...
func (op Operations) ListAllImage() (*ImageListIntentResponse, error) {
	entities := make([]*ImageIntentResponse, 0)

	resp, err := op.ListImage(&DSMetadata{
		Kind:   utils.String("image"),
		Length: utils.Int64(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListImage(&DSMetadata{
				Kind:   utils.String("image"),
				Length: utils.Int64(itemsPerPage),
				Offset: utils.Int64(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}
		resp.Entities = entities
	}

	return resp, nil
}

// ListAllCluster ...
func (op Operations) ListAllCluster() (*ClusterListIntentResponse, error) {
	entities := make([]*ClusterIntentResource, 0)

	resp, err := op.ListCluster(&DSMetadata{
		Kind:   utils.String("cluster"),
		Length: utils.Int64(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListCluster(&DSMetadata{
				Kind:   utils.String("cluster"),
				Length: utils.Int64(itemsPerPage),
				Offset: utils.Int64(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}
		resp.Entities = entities
	}

	return resp, nil
}

//GetTask ...
func (op Operations) GetTask(TaskUUID string) (*TasksResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/tasks/%s", TaskUUID)
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	tasksTesponse := new(TasksResponse)

	if err != nil {
		return nil, err
	}

	return tasksTesponse, op.client.Do(ctx, req, tasksTesponse)
}
