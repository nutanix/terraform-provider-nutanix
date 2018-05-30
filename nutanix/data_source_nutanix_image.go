package nutanix

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixImage() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNutanixImageRead,
		Schema: getDataSourceImageSchema(),
	}
}

func dataSourceNutanixImageRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	imageID, ok := d.GetOk("image_id")

	if !ok {
		return fmt.Errorf("please provide the required attribute vm_id")
	}

	// Make request to the API
	resp, err := conn.V3.GetImage(imageID.(string))
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

	if err := d.Set("availability_zone_reference", getReferenceValues(resp.Status.AvailabilityZoneReference)); err != nil {
		return err
	}
	if resp.Status.ClusterReference != nil {
		if err := d.Set("cluster_reference", getClusterReferenceValues(resp.Status.ClusterReference)); err != nil {
			return err
		}

		d.Set("cluster_reference_name", utils.StringValue(resp.Status.ClusterReference.Name))
	}

	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))
	d.Set("description", utils.StringValue(resp.Status.Description))
	d.Set("state", utils.StringValue(resp.Status.State))
	d.Set("image_type", resp.Status.Resources.ImageType)
	d.Set("source_uri", resp.Status.Resources.SourceURI)
	d.Set("size_bytes", utils.Int64Value(resp.Status.Resources.SizeBytes))

	var uriList []string
	for _, uri := range resp.Status.Resources.RetrievalURIList {
		uriList = append(uriList, utils.StringValue(uri))
	}

	if err := d.Set("retrieval_uri_list", uriList); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return nil
}

func getDataSourceImageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"image_id": {
			Type:     schema.TypeString,
			Required: true,
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
	}
}
