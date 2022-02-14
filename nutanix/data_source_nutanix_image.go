package nutanix

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixImage() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceNutanixImageRead,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNutanixDatasourceImageInstanceResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceDatasourceImageInstanceStateUpgradeV0,
				Version: 0,
			},
		},
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": categoriesSchema(),
			"owner_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

func dataSourceNutanixImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection
	conn := meta.(*Client).API

	imageID, iok := d.GetOk("image_id")
	imageName, nok := d.GetOk("image_name")

	if !iok && !nok {
		return diag.Errorf("please provide one of image_id or image_name attributes")
	}

	var reqErr error
	var resp *v3.ImageIntentResponse

	if iok {
		resp, reqErr = findImageByUUID(conn, imageID.(string))
	} else {
		resp, reqErr = findImageByName(conn, imageName.(string))
	}

	if reqErr != nil {
		return diag.FromErr(reqErr)
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("categories", c); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("availability_zone_reference", flattenReferenceValues(resp.Status.AvailabilityZoneReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := flattenClusterReference(resp.Status.ClusterReference, d); err != nil {
		return diag.FromErr(err)
	}
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
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(resp.Metadata.UUID))

	return nil
}

func findImageByUUID(conn *v3.Client, uuid string) (*v3.ImageIntentResponse, error) {
	return conn.V3.GetImage(uuid)
}

func findImageByName(conn *v3.Client, name string) (*v3.ImageIntentResponse, error) {
	filter := fmt.Sprintf("name==%s", name)
	resp, err := conn.V3.ListAllImage(filter)
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

func resourceDatasourceImageInstanceStateUpgradeV0(ctx context.Context, is map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Entering resourceDatasourceImageInstanceStateUpgradeV0")
	return resourceNutanixCategoriesMigrateState(is, meta)
}

func resourceNutanixDatasourceImageInstanceResourceV0() *schema.Resource {
	return &schema.Resource{
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
