package nutanix

// import (
// 	"fmt"
// 	"log"
// )

// // SubnetAPIInstance sets the nutanixV3.VmApi from the V3Client
// func SubnetAPIInstance(c *V3Client) *nutanixV3.SubnetApi {
// 	APIInstance := nutanixV3.NewSubnetApi()
// 	APIInstance.Configuration.Username = c.Username
// 	APIInstance.Configuration.Password = c.Password
// 	APIInstance.Configuration.BasePath = c.URL
// 	APIInstance.Configuration.APIClient.Insecure = c.Insecure
// 	return APIInstance
// }

// func (c *V3Client) SubnetExists(name string) (string, error) {

// 	log.Printf("[DEBUG] Get Subnet Existance : %s", name)

// 	APIInstance := SubnetAPIInstance(c)
// 	subnet_entities := nutanixV3.SubnetListMetadata{}
// 	var subnet_uuid string

// 	subnet_list, APIResponse, err := APIInstance.SubnetsListPost(subnet_entities)
// 	if err != nil {
// 		return "", err
// 	}

// 	err = checkAPIResponse(*APIResponse)
// 	if err != nil {
// 		return "", err
// 	}

// 	for _, subnet := range subnet_list.Entities {
// 		if subnet.Status.Name == name {
// 			subnet_uuid = subnet.Metadata.Uuid
// 		}
// 	}
// 	return subnet_uuid, nil
// }

// func (c *V3Client) WaitForSubnetProcess(uuid string) (bool, error) {
// 	APIInstance := SubnetAPIInstance(c)
// 	for {
// 		SubnetIntentResponse, APIresponse, err := APIInstance.SubnetsUuidGet(uuid)
// 		if err != nil {
// 			return false, err
// 		}
// 		err = checkAPIResponse(*APIresponse)
// 		if err != nil {
// 			return false, err
// 		}
// 		if SubnetIntentResponse.Status.State == "COMPLETE" {
// 			return true, nil
// 		} else if SubnetIntentResponse.Status.State == "ERROR" {
// 			return false, fmt.Errorf("%s", SubnetIntentResponse.Status.MessageList[0].Message)
// 		}
// 	}
// 	return false, nil
// }
