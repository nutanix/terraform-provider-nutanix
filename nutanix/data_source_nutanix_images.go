package nutanix

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixImagesRead,
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sort_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"filter": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"length": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sort_order": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"offset": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"categories": categoriesSchema(),
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
				},
			},
		},
	}
}

func dataSourceNutanixImagesRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	resp, err := conn.V3.ListAllImage()
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}

	entities := make([]map[string]interface{}, len(resp.Entities))
	for k, v := range resp.Entities {
		entity := make(map[string]interface{})
		m, c := setRSEntityMetadata(v.Metadata)

		entity["metadata"] = m
		entity["project_reference"] = getReferenceValues(v.Metadata.ProjectReference)
		entity["owner_reference"] = getReferenceValues(v.Metadata.OwnerReference)
		entity["categories"] = c
		entity["name"] = utils.StringValue(v.Status.Name)
		entity["description"] = utils.StringValue(v.Status.Description)
		entity["availability_zone_reference"] = getReferenceValues(v.Status.AvailabilityZoneReference)
		if v.Status.ClusterReference != nil {
			entity["cluster_reference"] = getClusterReferenceValues(v.Status.ClusterReference)
			entity["cluster_reference_name"] = utils.StringValue(v.Status.ClusterReference.Name)
		}
		entity["state"] = utils.StringValue(v.Status.State)
		entity["image_type"] = utils.StringValue(v.Status.Resources.ImageType)
		entity["source_uri"] = utils.StringValue(v.Status.Resources.SourceURI)
		entity["size_bytes"] = utils.Int64Value(v.Status.Resources.SizeBytes)
		entity["retrieval_uri_list"] = utils.StringValueSlice(v.Status.Resources.RetrievalURIList)

		entities[k] = entity
	}

	if err := d.Set("entities", entities); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return nil
}
