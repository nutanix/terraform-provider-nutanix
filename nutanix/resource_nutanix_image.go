package nutanix

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//ImageKind Represents kind of resource
	ImageKind = "image"
)

func resourceNutanixImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixImageCreate,
		Read:   resourceNutanixImageRead,
		Update: resourceNutanixImageUpdate,
		Delete: resourceNutanixImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getImageSchema(),
	}
}

func resourceNutanixImageCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Image: %s", d.Get("name").(string))

	conn := meta.(*Client).API

	request := &v3.ImageIntentInput{}
	spec := &v3.Image{}
	metadata := &v3.Metadata{}
	image := &v3.ImageResources{}

	n, nok := d.GetOk("name")
	desc, descok := d.GetOk("description")

	_, iok := d.GetOk("source_uri")
	_, pok := d.GetOk("source_path")

	// if both path and uri are provided, return an error
	if iok && pok {
		return errors.New("Both source_uri and source_path provided")
	}

	// Read Arguments and set request values
	if v, ok := d.GetOk("api_version"); ok {
		request.APIVersion = utils.String(v.(string))
	}

	if !nok {
		return fmt.Errorf("Please provide the required attribute name")
	}

	if err := getMetadataAttributes(d, metadata, "image"); err != nil {
		return err
	}

	if descok {
		spec.Description = utils.String(desc.(string))
	}

	if err := getImageResource(d, image); err != nil {
		return err
	}

	spec.Name = utils.String(n.(string))
	spec.Resources = image

	request.Metadata = metadata
	request.Spec = spec

	imageUUID, err := resourceNutanixImageExists(conn, n.(string))

	if err != nil {
		return err
	}

	if imageUUID != nil {
		return fmt.Errorf("Image already with name %s exists in the given cluster, UUID %s", d.Get("name").(string), *imageUUID)
	}

	//Make request to the API
	resp, err := conn.V3.CreateImage(request)
	if err != nil {
		return err
	}

	UUID := *resp.Metadata.UUID
	//set terraform state
	d.SetId(UUID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"COMPLETE"},
		Refresh:    imageStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for image (%s) to create: %s", d.Id(), err)
	}

	// if we need to upload an image, we do it now
	if pok {
		path := d.Get("source_path")

		err = conn.V3.UploadImage(UUID, path.(string))
		if err != nil {

			resourceNutanixImageDelete(d, meta)

			return fmt.Errorf("Failed uploading image: %s", err)
		}
	}

	return resourceNutanixImageRead(d, meta)
}

func resourceNutanixImageRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Image: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, err := conn.V3.GetImage(d.Id())
	if err != nil {
		return err
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}

	if err := d.Set("owner_reference", getReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))
	d.Set("description", utils.StringValue(resp.Status.Description))

	if err := d.Set("availability_zone_reference", getReferenceValues(resp.Status.AvailabilityZoneReference)); err != nil {
		return err
	}
	if err := d.Set("cluster_reference", getClusterReferenceValues(resp.Status.ClusterReference)); err != nil {
		return err
	}

	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}

	if err := d.Set("image_type", resp.Status.Resources.ImageType); err != nil {
		return err
	}

	if err := d.Set("source_uri", resp.Status.Resources.SourceURI); err != nil {
		return err
	}

	if err := d.Set("size_bytes", resp.Status.Resources.SizeBytes); err != nil {
		return err
	}

	checksum := make(map[string]string)
	if resp.Status.Resources.Checksum != nil {
		checksum["checksum_algorithm"] = utils.StringValue(resp.Status.Resources.Checksum.ChecksumAlgorithm)
		checksum["checksum_value"] = utils.StringValue(resp.Status.Resources.Checksum.ChecksumValue)
	}

	if err := d.Set("checksum", checksum); err != nil {
		return err
	}

	version := make(map[string]string)
	if resp.Status.Resources.Version != nil {
		version["product_version"] = utils.StringValue(resp.Status.Resources.Version.ProductVersion)
		version["product_name"] = utils.StringValue(resp.Status.Resources.Version.ProductName)
	}

	if err := d.Set("version", version); err != nil {
		return err
	}

	var uriList []string
	for _, uri := range resp.Status.Resources.RetrievalURIList {
		uriList = append(uriList, utils.StringValue(uri))
	}

	return d.Set("retrieval_uri_list", uriList)
}

func resourceNutanixImageUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	// get state
	request := &v3.ImageIntentInput{}
	metadata := &v3.Metadata{}
	spec := &v3.Image{}
	res := &v3.ImageResources{}

	response, err := conn.V3.GetImage(d.Id())

	if err != nil {
		return err
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if response.Spec != nil {
		spec = response.Spec

		if response.Spec.Resources != nil {
			res = response.Spec.Resources
		}
	}

	if d.HasChange("categories") {
		catl := d.Get("categories").([]interface{})

		if len(catl) > 0 {
			cl := make(map[string]string)
			for _, v := range catl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["value"]; kok && k.(string) != "" {
						cl[i.(string)] = k.(string)
					}
				}
			}
			metadata.Categories = cl
		} else {
			metadata.Categories = nil
		}
	}

	if d.HasChange("owner_reference") {
		or := d.Get("owner_reference").(map[string]interface{})
		metadata.OwnerReference = validateRef(or)
	}

	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").(map[string]interface{})
		metadata.ProjectReference = validateRef(pr)
	}

	if d.HasChange("name") {
		spec.Name = utils.String(d.Get("name").(string))
	}
	if d.HasChange("description") {
		spec.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("source_uri") || d.HasChange("checksum") {
		if err := getImageResource(d, res); err != nil {
			return err
		}
		spec.Resources = res
	}

	request.Metadata = metadata
	request.Spec = spec

	_, errUpdate := conn.V3.UpdateImage(d.Id(), request)
	if errUpdate != nil {
		return errUpdate
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"COMPLETE"},
		Refresh:    imageStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for image (%s) to update: %s", d.Id(), err)
	}

	return resourceNutanixImageRead(d, meta)
}

func resourceNutanixImageDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting Image: %s", d.Get("name").(string))

	conn := meta.(*Client).API
	UUID := d.Id()

	if err := conn.V3.DeleteImage(UUID); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "DELETE_IN_PROGRESS", "COMPLETE"},
		Target:     []string{"DELETED"},
		Refresh:    imageStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for image (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getImageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"metadata": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"creation_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_hash": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"categories": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"owner_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"project_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"availability_zone_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"retrieval_uri_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"image_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"checksum": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"checksum_algorithm": {
						Type:     schema.TypeString,
						Required: true,
					},
					"checksum_value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"source_uri": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"source_path": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"product_version": {
						Type:     schema.TypeString,
						Required: true,
					},
					"product_name": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"architecture": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"size_bytes": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func getImageMetadaAttributes(d *schema.ResourceData, metadata *v3.ImageMetadata) error {
	m, mok := d.GetOk("metadata")
	metad := m.(map[string]interface{})

	if !mok {
		return fmt.Errorf("please provide metadata required attributes")
	}

	metadata.Kind = utils.String(metad["kind"].(string))

	if v, ok := metad["uuid"]; ok && v != "" {
		metadata.UUID = utils.String(v.(string))
	}
	if v, ok := metad["spec_version"]; ok && v != 0 {
		i, err := strconv.Atoi(v.(string))
		if err != nil {
			return err
		}
		metadata.SpecVersion = utils.Int64(int64(i))
	}
	if v, ok := metad["spec_hash"]; ok && v != "" {
		metadata.SpecHash = utils.String(v.(string))
	}
	if v, ok := metad["name"]; ok {
		metadata.Name = utils.String(v.(string))
	}
	if v, ok := d.GetOk("categories"); ok {
		p := v.([]interface{})
		if len(p) > 0 {
			c := p[0].(map[string]interface{})
			labels := map[string]string{}

			for k, v := range c {
				labels[k] = v.(string)
			}
			metadata.Categories = labels
		}
	}
	if p, ok := d.GetOk("project_reference"); ok {
		pr := p.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(pr["kind"].(string)),
			UUID: utils.String(pr["uuid"].(string)),
		}
		if v1, ok1 := pr["name"]; ok1 {
			r.Name = utils.String(v1.(string))
		}
		metadata.ProjectReference = r
	}
	if o, ok := metad["owner_reference"]; ok {
		or := o.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(or["kind"].(string)),
			UUID: utils.String(or["uuid"].(string)),
		}
		if v1, ok1 := or["name"]; ok1 {
			r.Name = utils.String(v1.(string))
		}
		metadata.OwnerReference = r
	}

	return nil
}

func getImageResource(d *schema.ResourceData, image *v3.ImageResources) error {
	cs, csok := d.GetOk("checksum")
	checks := &v3.Checksum{}

	if su, suok := d.GetOk("source_uri"); suok {
		ext := filepath.Ext(su.(string))
		if ext == ".qcow2" {
			image.ImageType = utils.String("DISK_IMAGE")
		} else if ext == ".iso" {
			image.ImageType = utils.String("ISO_IMAGE")
		} else {
			// By default assuming the image to be raw disk image.
			image.ImageType = utils.String("DISK_IMAGE")
		}
		// set source uri
		image.SourceURI = utils.String(su.(string))
	}

	if csok {
		checksum := cs.(map[string]interface{})
		ca, caok := checksum["checksum_algorithm"]
		cv, cvok := checksum["checksum_value"]

		if caok {
			if ca.(string) == "" {
				return fmt.Errorf("'checksum_algorithm' is not given")
			}
			checks.ChecksumAlgorithm = utils.String(ca.(string))
		}
		if cvok {
			if cv.(string) == "" {
				return fmt.Errorf("'checksum_value' is not given")
			}
			checks.ChecksumValue = utils.String(cv.(string))
		}
		image.Checksum = checks
	}

	return nil
}

func resourceNutanixImageExists(conn *v3.Client, name string) (*string, error) {
	log.Printf("[DEBUG] Get Image Existence : %s", name)

	imageEntities := &v3.DSMetadata{}
	var imageUUID *string

	imageList, err := conn.V3.ListImage(imageEntities)

	if err != nil {
		return nil, err
	}

	for _, image := range imageList.Entities {
		if image.Status.Name == utils.String(name) {
			imageUUID = image.Metadata.UUID
		}
	}
	return imageUUID, nil
}

func imageStateRefreshFunc(client *v3.Client, uuid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.V3.GetImage(uuid)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return v, "DELETED", nil
			}
			log.Printf("ERROR %s", err)
			return nil, "", err
		}

		return v, *v.Status.State, nil
	}
}
