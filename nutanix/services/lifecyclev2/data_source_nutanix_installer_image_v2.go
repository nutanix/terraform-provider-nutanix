package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	config "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixInstallerImageV2 reads the details of a single installer image
// registered with Foundation Central, identified by its external ID.
func DatasourceNutanixInstallerImageV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixInstallerImageV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A globally unique identifier of an instance that is suitable for external consumption.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the image.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the installer image. One of AOS, AHV or ESX.",
			},
			"source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source of the image. One of LOCAL or REMOTE_URL.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL from where the image can be downloaded.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the image.",
			},
			"certificate_chain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate chain for the image URL.",
			},
			"metadata_download_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL from where the image metadata can be downloaded for an AOS image.",
			},
			"checksum": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Checksum value of the image. This is applicable only for AHV images.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sha256": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "SHA-256 checksum of the image.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hex_digest": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "SHA-256 checksum value in hexadecimal format (64 characters).",
									},
								},
							},
						},
						"md5": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "MD5 checksum of the image.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hex_digest": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "MD5 checksum value in hexadecimal format (32 characters).",
									},
								},
							},
						},
					},
				},
			},
			"file_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the image file on the cluster.",
			},
			"metadata_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the image metadata file on the cluster.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the image entity was created.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
			},
			"links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A HATEOAS style link for the response.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL at which the entity described by the link can be accessed.",
						},
						"rel": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A name that identifies the relationship of the link to the object that is returned by the URL.",
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixInstallerImageV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	extID := d.Get("ext_id")

	resp, err := conn.InstallerImagesAPIInstance.GetImageById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching image : %v", err)
	}

	getResp := resp.Data.GetValue().(config.Image)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if getResp.Type != nil {
		if err := d.Set("type", getResp.Type.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.Source != nil {
		if err := d.Set("source", getResp.Source.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("url", getResp.Url); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", getResp.Version); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("certificate_chain", getResp.CertificateChain); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata_download_url", getResp.MetadataDownloadUrl); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("checksum", flattenImageChecksum(getResp.Checksum)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.FileStatus != nil {
		if err := d.Set("file_status", getResp.FileStatus.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.MetadataStatus != nil {
		if err := d.Set("metadata_status", getResp.MetadataStatus.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenImageLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
