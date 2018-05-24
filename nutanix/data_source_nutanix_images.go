package nutanix

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixImages() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNutanixImagesRead,
		Schema: getDataSourceImagesSchema(),
	}
}

func dataSourceNutanixImagesRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	metadata := &v3.ImageListMetadata{}

	if v, ok := d.GetOk("metadata"); ok {
		m := v.(map[string]interface{})
		metadata.Kind = utils.String("image")
		if mv, mok := m["sort_attribute"]; mok {
			metadata.SortAttribute = utils.String(mv.(string))
		}
		if mv, mok := m["filter"]; mok {
			metadata.Filter = utils.String(mv.(string))
		}
		if mv, mok := m["length"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return err
			}
			metadata.Length = utils.Int64(int64(i))
		}
		if mv, mok := m["sort_order"]; mok {
			metadata.SortOrder = utils.String(mv.(string))
		}
		if mv, mok := m["offset"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return err
			}
			metadata.Offset = utils.Int64(int64(i))
		}
	}

	// Make request to the API
	resp, err := conn.V3.ListImage(metadata)
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}

	entities := make([]map[string]interface{}, len(resp.Entities))
	for k, v := range resp.Entities {
		entity := make(map[string]interface{})
		// set metadata values
		metadata := make(map[string]interface{})
		metadata["last_update_time"] = utils.TimeValue(v.Metadata.LastUpdateTime).String()
		metadata["kind"] = utils.StringValue(v.Metadata.Kind)
		metadata["uuid"] = utils.StringValue(v.Metadata.UUID)
		metadata["creation_time"] = utils.TimeValue(v.Metadata.CreationTime).String()
		metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(v.Metadata.SpecVersion)))
		metadata["spec_hash"] = utils.StringValue(v.Metadata.SpecHash)
		metadata["name"] = utils.StringValue(v.Metadata.Name)
		entity["metadata"] = metadata

		if v.Metadata.Categories != nil {
			categories := v.Metadata.Categories
			var catList []map[string]interface{}

			for name, values := range categories {
				catItem := make(map[string]interface{})
				catItem["name"] = name
				catItem["value"] = values
				catList = append(catList, catItem)
			}
			entity["categories"] = catList
		}

		entity["api_version"] = utils.StringValue(v.APIVersion)

		pr := make(map[string]interface{})
		if v.Metadata.ProjectReference != nil {
			pr["kind"] = utils.StringValue(v.Metadata.ProjectReference.Kind)
			pr["name"] = utils.StringValue(v.Metadata.ProjectReference.Name)
			pr["uuid"] = utils.StringValue(v.Metadata.ProjectReference.UUID)
		}
		entity["project_reference"] = pr

		or := make(map[string]interface{})
		if v.Metadata.OwnerReference != nil {
			or["kind"] = utils.StringValue(v.Metadata.OwnerReference.Kind)
			or["name"] = utils.StringValue(v.Metadata.OwnerReference.Name)
			or["uuid"] = utils.StringValue(v.Metadata.OwnerReference.UUID)
		}
		entity["owner_reference"] = or

		entity["name"] = utils.StringValue(v.Status.Name)
		entity["description"] = utils.StringValue(v.Status.Description)

		// set availability zone reference values
		availabilityZoneReference := make(map[string]interface{})
		if v.Status.AvailabilityZoneReference != nil {
			availabilityZoneReference["kind"] = utils.StringValue(v.Status.AvailabilityZoneReference.Kind)
			availabilityZoneReference["name"] = utils.StringValue(v.Status.AvailabilityZoneReference.Name)
			availabilityZoneReference["uuid"] = utils.StringValue(v.Status.AvailabilityZoneReference.UUID)
		}
		entity["availability_zone_reference"] = availabilityZoneReference
		// set cluster reference values
		clusterReference := make(map[string]interface{})
		if v.Status.ClusterReference != nil {
			clusterReference["kind"] = utils.StringValue(v.Status.ClusterReference.Kind)
			clusterReference["name"] = utils.StringValue(v.Status.ClusterReference.Name)
			clusterReference["uuid"] = utils.StringValue(v.Status.ClusterReference.UUID)
		}
		entity["cluster_reference"] = clusterReference
		entity["state"] = utils.StringValue(v.Status.State)

		entity["image_type"] = utils.StringValue(v.Status.Resources.ImageType)
		entity["source_uri"] = utils.StringValue(v.Status.Resources.SourceURI)
		entity["size_bytes"] = utils.Int64Value(v.Status.Resources.SizeBytes)

		var uriList []string
		for _, uri := range v.Status.Resources.RetrievalURIList {
			uriList = append(uriList, utils.StringValue(uri))
		}

		entity["retrieval_uri_list"] = uriList

		entities[k] = entity
	}

	if err := d.Set("entities", entities); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return nil
}

func getDataSourceImagesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
				},
			},
		},
	}
}
