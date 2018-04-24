package nutanix

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"path/filepath"

// 	"github.com/hashicorp/terraform/helper/schema"
// )

// const (
// 	IMAGE_KIND = "image"
// )

// func resourceNutanixImage() *schema.Resource {

// 	return &schema.Resource{
// 		Create: resourceNutanixImageCreate,
// 		Read:   resourceNutanixImageRead,
// 		Update: resourceNutanixImageUpdate,
// 		Delete: resourceNutanixImageDelete,

// 		Schema: map[string]*schema.Schema{
// 			"name": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"uuid": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"kind": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"checksum_algorithm": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"checksum_value": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"source_uri": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"image_type": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"retrieval_uri_list": {
// 				Type:     schema.TypeList,
// 				Optional: true,
// 				Computed: true,
// 				Elem:     &schema.Schema{Type: schema.TypeString},
// 			},
// 			"size_bytes": {
// 				Type:     schema.TypeInt,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"description": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 		},
// 	}
// }

// func resourceNutanixImageCreate(d *schema.ResourceData, meta interface{}) error {

// 	log.Printf("[DEBUG] Creating Image: %s", d.Get("name").(string))

// 	client := meta.(*V3Client)
// 	ImageAPIInstance := ImageAPIInstance(client)

// 	image := nutanixV3.ImageIntentInput{
// 		ApiVersion: API_VERSION,
// 		Metadata: nutanixV3.ImageMetadata{
// 			Name: d.Get("name").(string),
// 			Kind: "image",
// 		},
// 		Spec: nutanixV3.Image{
// 			Description: d.Get("description").(string),
// 			Name:        d.Get("name").(string),
// 			Resources: nutanixV3.ImageResources{
// 				SourceUri: d.Get("source_uri").(string),
// 			},
// 		},
// 	}

// 	ext := filepath.Ext(d.Get("source_uri").(string))
// 	if ext == ".qcow2" {
// 		image.Spec.Resources.ImageType = "DISK_IMAGE"
// 		if d.Get("checksum_algorithm").(string) != "" || d.Get("checksum_value").(string) != "" {
// 			return fmt.Errorf("Checksums are not supported for images that require conversion '%s'", ext)
// 		}
// 	} else if ext == ".iso" {
// 		image.Spec.Resources.ImageType = "ISO_IMAGE"
// 	} else {
// 		// By default assuming the image to be raw disk image.
// 		image.Spec.Resources.ImageType = "DISK_IMAGE"
// 	}

// 	if d.Get("checksum_algorithm").(string) != "" || d.Get("checksum_value").(string) != "" {
// 		if d.Get("checksum_value").(string) == "" && d.Get("checksum_algorithm").(string) == "" {
// 			return fmt.Errorf("'checksum_value' or 'checksum_algorithm' is not given.")
// 		}
// 		image.Spec.Resources.Checksum.ChecksumAlgorithm = d.Get("checksum_algorithm").(string)
// 		image.Spec.Resources.Checksum.ChecksumValue = d.Get("checksum_value").(string)
// 	}

// 	image_uuid, err := client.ImageExists(d.Get("name").(string))
// 	if err != nil {
// 		return err
// 	}

// 	if image_uuid != "" {
// 		return fmt.Errorf("Image already with name %s exists in the given cluster, UUID %s", d.Get("name").(string), image_uuid)
// 	}

// 	image_json, _ := json.Marshal(image)
// 	log.Printf("[DEBUG] Image JSON :%s", image_json)

// 	ImageIntentResponse, APIResponse, err := ImageAPIInstance.ImagesPost(image)

// 	if err != nil {
// 		return err
// 	}

// 	err = checkAPIResponse(*APIResponse)
// 	if err != nil {
// 		return err
// 	}

// 	uuid := ImageIntentResponse.Metadata.Uuid
// 	status, err := client.WaitForImageProcess(uuid)
// 	for status != true {
// 		return err
// 	}
// 	d.SetId(uuid)
// 	return resourceNutanixImageRead(d, meta)
// }

// func resourceNutanixImageRead(d *schema.ResourceData, meta interface{}) error {

// 	log.Printf("[DEBUG] Reading Image: %s", d.Get("name").(string))

// 	client := meta.(*V3Client)
// 	ImageAPIInstance := ImageAPIInstance(client)
// 	uuid := d.Id()

// 	get_image, APIResponse, err := ImageAPIInstance.ImagesUuidGet(uuid)
// 	if err != nil {
// 		return err
// 	}

// 	err = checkAPIResponse(*APIResponse)
// 	if err != nil {
// 		return err
// 	}

// 	get_image_json, _ := json.Marshal(get_image)
// 	log.Printf("[DEBUG] Read Image %s", get_image_json)

// 	d.Set("name", get_image.Status.Name)
// 	d.Set("image_type", get_image.Status.Resources.ImageType)
// 	d.Set("source_uri", get_image.Status.Resources.SourceUri)
// 	d.Set("size_bytes", get_image.Status.Resources.SizeBytes)
// 	d.Set("description", get_image.Status.Description)
// 	var uri_list []string
// 	for _, uri := range get_image.Status.Resources.RetrievalUriList {
// 		uri_list = append(uri_list, uri)
// 	}
// 	d.Set("retrieval_uri_list", uri_list)

// 	return nil
// }

// func resourceNutanixImageUpdate(d *schema.ResourceData, meta interface{}) error {

// 	return nil
// }

// func resourceNutanixImageDelete(d *schema.ResourceData, meta interface{}) error {

// 	log.Printf("[DEBUG] Deleting Image: %s", d.Get("name").(string))

// 	client := meta.(*V3Client)
// 	ImageAPIInstance := ImageAPIInstance(client)
// 	uuid := d.Id()

// 	APIResponse, err := ImageAPIInstance.ImagesUuidDelete(uuid)
// 	if err != nil {
// 		return err
// 	}
// 	err = checkAPIResponse(*APIResponse)
// 	if err != nil {
// 		return err
// 	}

// 	d.SetId("")
// 	return nil
// }
