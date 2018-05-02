package nutanix

import (
	"fmt"
	"strconv"

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

func getDataSourceImageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"image_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
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
			ForceNew: true,
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
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
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
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
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
