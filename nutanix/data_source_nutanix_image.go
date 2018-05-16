package nutanix

import (
	"fmt"
	"strconv"

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
	conn := meta.(*NutanixClient).API

	imageID, ok := d.GetOk("image_id")

	if !ok {
		return fmt.Errorf("please provide the required attribute vm_id")
	}

	// Make request to the API
	resp, err := conn.V3.GetImage(imageID.(string))
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
	if resp.Metadata.OwnerReference != nil {
		or["kind"] = utils.StringValue(resp.Metadata.OwnerReference.Kind)
		or["name"] = utils.StringValue(resp.Metadata.OwnerReference.Name)
		or["uuid"] = utils.StringValue(resp.Metadata.OwnerReference.UUID)

	}
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
	clusterReference := make(map[string]interface{})
	if resp.Status.ClusterReference != nil {
		clusterReference["kind"] = utils.StringValue(resp.Status.ClusterReference.Kind)
		clusterReference["name"] = utils.StringValue(resp.Status.ClusterReference.Name)
		clusterReference["uuid"] = utils.StringValue(resp.Status.ClusterReference.UUID)
	}
	if err := d.Set("cluster_reference", clusterReference); err != nil {
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
			Type:     schema.TypeMap,
			Computed: true,
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
		"message_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"message": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"reason": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"details": {
						Type:     schema.TypeMap,
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
					"name": {
						Type:     schema.TypeString,
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
