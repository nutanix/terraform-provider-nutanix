package nutanix

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceNutanixImage() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceNutanixImageRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"retrieval_uri_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"size_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixImageRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Datasource Image Read: %s", d.Get("name").(string))

	client := meta.(*V3Client)
	ImageAPIInstance := ImageAPIInstance(client)
	uuid, err := client.ImageExists(d.Get("name").(string))
	if err != nil {
		return err
	}

	if uuid == "" {
		return fmt.Errorf("Image doesn't exists in given cluster.")
	}

	get_image, APIResponse, err := ImageAPIInstance.ImagesUuidGet(uuid)
	if err != nil {
		return err
	}

	err = checkAPIResponse(*APIResponse)
	if err != nil {
		return err
	}

	get_image_json, _ := json.Marshal(get_image)
	log.Printf("[DEBUG] Read Image %s", get_image_json)

	d.Set("name", get_image.Status.Name)
	d.Set("uuid", uuid)
	d.Set("image_type", get_image.Status.Resources.ImageType)
	d.Set("source_uri", get_image.Status.Resources.SourceUri)
	d.Set("size_bytes", get_image.Status.Resources.SizeBytes)
	d.Set("description", get_image.Status.Description)
	var uri_list []string
	for _, uri := range get_image.Status.Resources.RetrievalUriList {
		uri_list = append(uri_list, uri)
	}
	d.Set("retrieval_uri_list", uri_list)
	d.SetId(uuid)

	return nil
}
