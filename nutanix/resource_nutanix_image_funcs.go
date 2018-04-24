package nutanix

// import (
// 	"fmt"
// 	"log"
// )

// // ImageAPIInstance sets the nutanixV3.VmApi from the V3Client
// func ImageAPIInstance(c *V3Client) *nutanixV3.ImagesApi {
// 	APIInstance := nutanixV3.NewImagesApi()
// 	APIInstance.Configuration.Username = c.Username
// 	APIInstance.Configuration.Password = c.Password
// 	APIInstance.Configuration.BasePath = c.URL
// 	APIInstance.Configuration.APIClient.Insecure = c.Insecure
// 	return APIInstance
// }

// func (c *V3Client) ImageExists(name string) (string, error) {

// 	log.Printf("[DEBUG] Get Image Existance : %s", name)

// 	ImageAPIInstance := ImageAPIInstance(c)
// 	image_entities := nutanixV3.ImageListMetadata{}
// 	var image_uuid string

// 	image_list, APIResponse, err := ImageAPIInstance.ImagesListPost(image_entities)
// 	if err != nil {
// 		return "", err
// 	}

// 	err = checkAPIResponse(*APIResponse)
// 	if err != nil {
// 		return "", err
// 	}

// 	for _, img := range image_list.Entities {
// 		if img.Status.Name == name {
// 			image_uuid = img.Metadata.Uuid
// 		}
// 	}
// 	return image_uuid, nil
// }

// func (c *V3Client) WaitForImageProcess(uuid string) (bool, error) {
// 	APIInstance := ImageAPIInstance(c)
// 	for {
// 		ImageIntentResponse, APIresponse, err := APIInstance.ImagesUuidGet(uuid)
// 		if err != nil {
// 			return false, err
// 		}
// 		err = checkAPIResponse(*APIresponse)
// 		if err != nil {
// 			return false, err
// 		}
// 		if ImageIntentResponse.Status.State == "COMPLETE" {
// 			return true, nil
// 		} else if ImageIntentResponse.Status.State == "ERROR" {
// 			return false, fmt.Errorf("%s", ImageIntentResponse.Status.MessageList[0].Message)
// 		}
// 	}
// 	return false, nil
// }
