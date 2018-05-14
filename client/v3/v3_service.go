package v3

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

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
	ListNetworkSecurityRule(getEntitiesRequest *ListMetadata) (*NetworkSecurityRuleListIntentResponse, error)
	GetNetworkSecurityRule(UUID string) (*NetworkSecurityRuleIntentResponse, error)
	DeleteNetworkSecurityRule(UUID string) error
	CreateNetworkSecurityRule(request *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error)
	ListCluster(getEntitiesRequest *ClusterListMetadataOutput) (*ClusterListIntentResponse, error)
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

/*GetCluster gets a CLUSTER
 * This operation gets a CLUSTER.
 *
 * @param uuid The UUID of the entity.
 * @return *ImageIntentResponse
 */
// func (op Operations) GetCluster(UUID string) (*ImageIntentResponse, error) {
// 	ctx := context.TODO()

// 	path := fmt.Sprintf("/images/%s", UUID)

// 	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
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

/*ListCluster gets a list of CLUSTERS
 * This operation gets a list of CLUSTERS, allowing for sorting and pagination. Note: Entities that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest
 * @return *ClusterListIntentResponse
 */
func (op Operations) ListCluster(getEntitiesRequest *ClusterListMetadataOutput) (*ClusterListIntentResponse, error) {
	ctx := context.TODO()
	path := "/clusters/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)

	if err != nil {
		return nil, err
	}

	clusterList := new(ClusterListIntentResponse)
	err = op.client.Do(ctx, req, clusterList)

	if err != nil {
		return nil, err
	}

	return clusterList, nil
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

//CreateOrUpdateCategoryKey ...
func (op Operations) CreateOrUpdateCategoryKey(body *CategoryKey) (*CategoryKeyStatus, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s", utils.StringValue(body.Name))

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)

	categoryKeyResponse := new(CategoryKeyStatus)

	err = op.client.Do(ctx, req, categoryKeyResponse)
	if err != nil {
		return nil, err
	}

	return categoryKeyResponse, nil
}

/*ListCategories gets a list of Categories
 * This operation gets a list of Categories, allowing for sorting and pagination. Note: Entities that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest
 * @return *ImageListIntentResponse
 */
func (op Operations) ListCategories(getEntitiesRequest *CategoryListMetadata) (*CategoryKeyListResponse, error) {
	ctx := context.TODO()
	path := "/categories/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)

	if err != nil {
		return nil, err
	}

	categoryKeyListResponse := new(CategoryKeyListResponse)
	err = op.client.Do(ctx, req, categoryKeyListResponse)

	if err != nil {
		return nil, err
	}

	return categoryKeyListResponse, nil
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
	if err != nil {
		return nil, err
	}

	categoryKeyStatusResponse := new(CategoryKeyStatus)

	err = op.client.Do(ctx, req, categoryKeyStatusResponse)

	if err != nil {
		return nil, err
	}

	return categoryKeyStatusResponse, nil
}

/*ListCategoryValues gets a list of Category values for a specific key
 * This operation gets a list of Categories, allowing for sorting and pagination. Note: Entities that have not been created successfully are not listed.
 *
 * @param name
 * @param getEntitiesRequest
 * @return *CategoryValueListResponse
 */
func (op Operations) ListCategoryValues(name string, getEntitiesRequest *CategoryListMetadata) (*CategoryValueListResponse, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/categories/%s/list", name)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	categoryValueListResponse := new(CategoryValueListResponse)
	err = op.client.Do(ctx, req, categoryValueListResponse)

	if err != nil {
		return nil, err
	}

	return categoryValueListResponse, nil
}

//CreateOrUpdateCategoryValue ...
func (op Operations) CreateOrUpdateCategoryValue(name string, body *CategoryValue) (*CategoryValueStatus, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s/%s", name, utils.StringValue(body.Value))

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	categoryValueResponse := new(CategoryValueStatus)

	err = op.client.Do(ctx, req, categoryValueResponse)

	return categoryValueResponse, nil
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
	if err != nil {
		return nil, err
	}

	categoryValueStatusResponse := new(CategoryValueStatus)

	err = op.client.Do(ctx, req, categoryValueStatusResponse)
	if err != nil {
		return nil, err
	}

	return categoryValueStatusResponse, nil
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

	err = op.client.Do(ctx, req, categoryQueryResponse)

	if err != nil {
		return nil, err
	}

	return categoryQueryResponse, nil
}

/*CreateNetworkSecurityRule Creates a Network security rule
 * This operation submits a request to create a Network security rule based on the input parameters.
 *
 * @param request
 * @return *NetworkSecurityRuleIntentResponse
 */
func (op Operations) CreateNetworkSecurityRule(request *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/network_security_rules", request)
	networkSecurityRuleIntentResponse := new(NetworkSecurityRuleIntentResponse)

	err = op.client.Do(ctx, req, networkSecurityRuleIntentResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleIntentResponse, nil
}

/*DeleteNetworkSecurityRule Deletes a Network security rule
 * This operation submits a request to delete a Network security rule.
 *
 * @param UUID The UUID of the entity.
 * @return void
 */
func (op Operations) DeleteNetworkSecurityRule(UUID string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/network_security_rules/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
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
	if err != nil {
		return nil, err
	}

	networkSecurityRuleIntentResponse := new(NetworkSecurityRuleIntentResponse)

	err = op.client.Do(ctx, req, networkSecurityRuleIntentResponse)
	if err != nil {
		return nil, err
	}

	return networkSecurityRuleIntentResponse, nil
}

/*ListNetworkSecurityRule Gets all network security rules
 * This operation gets a list of Network security rules, allowing for sorting and pagination. Note: Entities that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest
 * @return *NetworkSecurityRuleListIntentResponse
 */
func (op Operations) ListNetworkSecurityRule(getEntitiesRequest *ListMetadata) (*NetworkSecurityRuleListIntentResponse, error) {
	ctx := context.TODO()
	path := "/network_security_rules/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)

	if err != nil {
		return nil, err
	}

	networkSecurityRuleListIntentResponse := new(NetworkSecurityRuleListIntentResponse)
	err = op.client.Do(ctx, req, networkSecurityRuleListIntentResponse)
	if err != nil {
		return nil, err
	}

	return networkSecurityRuleListIntentResponse, nil
}

/*UpdateNetworkSecurityRule Updates a Network security rule
 * This operation submits a request to update a Network security rule based on the input parameters.
 *
 * @param uuid The UUID of the entity.
 * @param body
 * @return void
 */
func (op Operations) UpdateNetworkSecurityRule(UUID string, body *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/network_security_rules/%s", UUID)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	networkSecurityRuleIntentResponse := new(NetworkSecurityRuleIntentResponse)

	err = op.client.Do(ctx, req, networkSecurityRuleIntentResponse)
	if err != nil {
		return nil, err
	}

	return networkSecurityRuleIntentResponse, nil
}
