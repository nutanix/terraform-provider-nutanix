package nutanix

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixImageRead,
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"image_name"},
			},
			"image_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"image_id"},
			},
			"api_version": {
				Type:     schema.TypeString,
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
						"kind": {
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
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
			"project_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
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
			"cluster_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cluster_reference_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retrieval_uri_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"image_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"checksum": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"checksum_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"checksum_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"source_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixImageRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	imageID, iok := d.GetOk("image_id")
	imageName, nok := d.GetOk("image_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of image_id or image_name attributes")
	}

	var reqErr error
	var resp *v3.ImageIntentResponse

	if iok {
		resp, reqErr = findImageByUUID(conn, imageID.(string))
	} else {
		resp, reqErr = findImageByName(conn, imageName.(string))
	}

	if reqErr != nil {
		return reqErr
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

	if err := d.Set("availability_zone_reference", getReferenceValues(resp.Status.AvailabilityZoneReference)); err != nil {
		return err
	}
	cluster := make(map[string]interface{})
	cl := ""
	if resp.Status.ClusterReference != nil {
		cluster = getClusterReferenceValues(resp.Status.ClusterReference)
		cl = utils.StringValue(resp.Status.ClusterReference.Name)
	}
	if err := d.Set("cluster_reference", cluster); err != nil {
		return err
	}
	d.Set("cluster_reference_name", cl)
	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))
	d.Set("description", utils.StringValue(resp.Status.Description))
	d.Set("state", utils.StringValue(resp.Status.State))
	d.Set("image_type", resp.Status.Resources.ImageType)
	d.Set("source_uri", resp.Spec.Resources.SourceURI)
	d.Set("size_bytes", utils.Int64Value(resp.Status.Resources.SizeBytes))

	uriList := make([]string, 0, len(resp.Status.Resources.RetrievalURIList))
	for _, uri := range resp.Status.Resources.RetrievalURIList {
		uriList = append(uriList, utils.StringValue(uri))
	}

	if err := d.Set("retrieval_uri_list", uriList); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return nil
}

func findImageByUUID(conn *v3.Client, uuid string) (*v3.ImageIntentResponse, error) {
	return conn.V3.GetImage(uuid)
}

func findImageByName(conn *v3.Client, name string) (*v3.ImageIntentResponse, error) {
	resp, err := conn.V3.ListAllImage()
	if err != nil {
		return nil, err
	}

	entities := resp.Entities

	found := make([]*v3.ImageIntentResponse, 0)
	for _, v := range entities {
		if *v.Spec.Name == name {
			found = append(found, v)
		}
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("your query returned more than one result. Please use image_id argument instead")
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("image with the given name, not found")
	}

	return findImageByUUID(conn, *found[0].Metadata.UUID)

}
