package nutanix

import (
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

	conn := meta.(*NutanixClient).API

	request := &v3.ImageIntentInput{}
	spec := &v3.Image{}
	metadata := &v3.ImageMetadata{}
	image := &v3.ImageResources{}

	n, nok := d.GetOk("name")
	desc, descok := d.GetOk("description")

	// Read Arguments and set request values
	if v, ok := d.GetOk("api_version"); ok {
		request.APIVersion = utils.String(v.(string))
	}

	if !nok {
		return fmt.Errorf("Please provide the required attribute name")
	}

	if err := getImageMetadaAttributes(d, metadata); err != nil {
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

	utils.PrintToJSON(request, "[DEBUG] Image request")

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
			"Error waiting for vm (%s) to create: %s", d.Id(), err)
	}
	return resourceNutanixImageRead(d, meta)
}

func resourceNutanixImageRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Image: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*NutanixClient).API

	// Make request to the API
	resp, err := conn.V3.GetImage(d.Id())
	if err != nil {
		return err
	}

	// set metadata values
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = resp.Metadata.LastUpdateTime.String()
	metadata["kind"] = utils.StringValue(resp.Metadata.Kind)
	metadata["uuid"] = utils.StringValue(resp.Metadata.UUID)
	metadata["creation_time"] = resp.Metadata.CreationTime.String()
	metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(resp.Metadata.SpecVersion)))
	metadata["spec_hash"] = utils.StringValue(resp.Metadata.SpecHash)
	metadata["name"] = utils.StringValue(resp.Metadata.Name)
	if err := d.Set("metadata", metadata); err != nil {
		return err
	}
	if err := d.Set("categories", resp.Metadata.Categories); err != nil {
		return err
	}

	or := make(map[string]interface{})
	or["kind"] = utils.StringValue(resp.Metadata.OwnerReference.Kind)
	or["name"] = utils.StringValue(resp.Metadata.OwnerReference.Name)
	or["uuid"] = utils.StringValue(resp.Metadata.OwnerReference.UUID)
	if err := d.Set("owner_reference", or); err != nil {
		return err
	}

	if err := d.Set("api_version", utils.StringValue(resp.APIVersion)); err != nil {
		return err
	}
	if err := d.Set("name", utils.StringValue(resp.Status.Name)); err != nil {
		return err
	}
	if err := d.Set("description", utils.StringValue(resp.Status.Description)); err != nil {
		return err
	}

	// set availability zone reference values
	availabilityZoneReference := make(map[string]interface{})
	if resp.Status.AvailabilityZoneReference != nil {
		availabilityZoneReference["kind"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Kind)
		availabilityZoneReference["name"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Name)
		availabilityZoneReference["uuid"] = utils.StringValue(resp.Status.AvailabilityZoneReference.UUID)
	}
	if err := d.Set("availability_zone_reference", availabilityZoneReference); err != nil {
		return err
	}
	// set cluster reference values
	if resp.Status.ClusterReference != nil {
		clusterReference := make(map[string]interface{})
		clusterReference["kind"] = utils.StringValue(resp.Status.ClusterReference.Kind)
		clusterReference["name"] = utils.StringValue(resp.Status.ClusterReference.Name)
		clusterReference["uuid"] = utils.StringValue(resp.Status.ClusterReference.UUID)
		if err := d.Set("cluster_reference", clusterReference); err != nil {
			return err
		}
	}

	// set message list values
	if resp.Status.MessageList != nil {
		messages := make([]map[string]interface{}, len(resp.Status.MessageList))
		for k, v := range resp.Status.MessageList {
			message := make(map[string]interface{})
			message["message"] = utils.StringValue(v.Message)
			message["reason"] = utils.StringValue(v.Reason)
			message["details"] = v.Details
			messages[k] = message
		}
		if err := d.Set("message_list", messages); err != nil {
			return err
		}
	}

	// set state value
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

	var uriList []string
	for _, uri := range resp.Status.Resources.RetrievalURIList {
		uriList = append(uriList, utils.StringValue(uri))
	}

	if err := d.Set("retrieval_uri_list", uriList); err != nil {
		return err
	}

	return nil
}

func resourceNutanixImageUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NutanixClient).API

	// get state
	request := &v3.ImageIntentInput{}
	metadata := &v3.ImageMetadata{}
	res := &v3.ImageResources{}

	if d.HasChange("metadata") ||
		d.HasChange("categories") ||
		d.HasChange("owner_reference") ||
		d.HasChange("project_reference") {
		if err := getImageMetadaAttributes(d, metadata); err != nil {
			return err
		}
		request.Metadata = metadata
	}

	if d.HasChange("name") {
		request.Spec.Name = utils.String(d.Get("name").(string))
	}
	if d.HasChange("description") {
		request.Spec.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("source_uri") || d.HasChange("checksum") {
		if err := getImageResource(d, res); err != nil {
			return err
		}
		request.Spec.Resources = res
	}
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
			"Error waiting for vm (%s) to update: %s", d.Id(), err)
	}

	return resourceNutanixImageRead(d, meta)
}

func resourceNutanixImageDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting Image: %s", d.Get("name").(string))

	conn := meta.(*NutanixClient).API
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
		"api_version": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"metadata": &schema.Schema{
			Type:     schema.TypeMap,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"creation_time": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"spec_version": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"spec_hash": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"categories": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"owner_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"project_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"state": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"availability_zone_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"message_list": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"message": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"reason": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"details": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"retrieval_uri_list": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"image_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"checksum": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"checksum_algorithm": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"checksum_value": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"source_uri": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"version": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"product_version": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"product_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"architecture": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"size_bytes": &schema.Schema{
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
	log.Printf("[DEBUG] Get Image Existance : %s", name)

	imageEntities := &v3.ImageListMetadata{}
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
