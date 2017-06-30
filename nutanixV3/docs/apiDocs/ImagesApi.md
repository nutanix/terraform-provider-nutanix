#\ImagesApi

##ImagesListPost
//  Get a list of IMAGEs
var api_instance := nutanix.ImagesApi
getEntitiesRequest := nutanix.ImageListMetadata() // ImageListMetadata | 
imageslistpost_response, api_response, err := api_instance.ImagesListPost(getEntitiesRequest)

##ImagesPost
//  Create a IMAGE
var api_instance := nutanix.ImagesApi
body := nutanix.ImageIntentInput() // ImageIntentInput | 
imagespost_response, api_response, err := api_instance.ImagesPost(body)

##ImagesUuidDelete
//  Delete a IMAGE
var api_instance := nutanix.ImagesApi
uuid := "uuid_example" // string | The UUID of the entity
imagesuuiddelete_response, api_response, err := api_instance.ImagesUuidDelete(uuid)

##ImagesUuidFileGet
//  Get Image Contents
var api_instance := nutanix.ImagesApi
uuid := "uuid_example" // string | The UUID of the entity
imagesuuidfileget_response, api_response, err := api_instance.ImagesUuidFileGet(uuid)

##ImagesUuidFilePut
//  Upload Image Contents
var api_instance := nutanix.ImagesApi
uuid := "uuid_example" // string | The UUID of the entity
image := "image_example" // string | 
imagesuuidfileput_response, api_response, err := api_instance.ImagesUuidFilePut(uuid, image)

##ImagesUuidGet
//  Get a IMAGE
var api_instance := nutanix.ImagesApi
uuid := "uuid_example" // string | The UUID of the entity
imagesuuidget_response, api_response, err := api_instance.ImagesUuidGet(uuid)

##ImagesUuidPut
//  Update a IMAGE
var api_instance := nutanix.ImagesApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.ImageIntentInput() // ImageIntentInput | 
imagesuuidput_response, api_response, err := api_instance.ImagesUuidPut(uuid, body)

