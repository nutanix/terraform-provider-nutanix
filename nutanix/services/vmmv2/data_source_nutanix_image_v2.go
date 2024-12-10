package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import5 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixImageV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixImageV4Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"checksum": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hex_digest": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url_source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"should_allow_insecure_url": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"basic_auth": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"username": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"password": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"vm_disk_source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"object_lite_source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"category_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_location_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"placement_policy_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"placement_policy_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"compliance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enforcement_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_cluster_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enforced_cluster_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"conflicting_policy_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixImageV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id")

	resp, err := conn.ImagesAPIInstance.GetImageById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching images : %v", err)
	}

	getResp := resp.Data.GetValue().(import5.Image)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", flattenImageType(getResp.Type)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("checksum", flattenOneOfImageChecksum(getResp.Checksum)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("size_bytes", getResp.SizeBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source", flattenOneOfImageSource(getResp.Source)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category_ext_ids", getResp.CategoryExtIds); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_location_ext_ids", getResp.ClusterLocationExtIds); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreateTime != nil {
		t := getResp.CreateTime
		if err := d.Set("create_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdateTime != nil {
		t := getResp.LastUpdateTime
		if err := d.Set("last_update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("placement_policy_status", flattenImagePlacementStatus(getResp.PlacementPolicyStatus)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)
	return nil
}

func flattenImageType(pr *import5.ImageType) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == import5.ImageType(two) {
			return "DISK_IMAGE"
		}
		if *pr == import5.ImageType(three) {
			return "ISO_IMAGE"
		}
	}
	return "UNKNOWN"
}

func flattenOneOfImageChecksum(pr *import5.OneOfImageChecksum) []map[string]interface{} {
	if pr != nil {
		resList := make([]map[string]interface{}, 0)

		sha := make(map[string]interface{})

		getVal := pr.ObjectType_

		if utils.StringValue(getVal) == "vmm.v4.content.ImageSha1Checksum" {
			sha1 := pr.GetValue().(import5.ImageSha1Checksum)

			sha["hex_digest"] = sha1.HexDigest
		} else {
			sha256 := pr.GetValue().(import5.ImageSha256Checksum)

			sha["hex_digest"] = sha256.HexDigest
		}
		resList = append(resList, sha)
		return resList
	}
	return nil
}

func flattenOneOfImageSource(pr *import5.OneOfImageSource) []map[string]interface{} {
	if pr != nil {
		urlSrcMap := make(map[string]interface{})
		urlSrcList := make([]map[string]interface{}, 0)

		vmDiskSrcMap := make(map[string]interface{})
		vmDiskSrcList := make([]map[string]interface{}, 0)

		objectLiteSrc := make(map[string]interface{})
		objectLiteSrcList := make([]map[string]interface{}, 0)

		if *pr.ObjectType_ == "vmm.v4.content.UrlSource" {
			urlSrc := pr.GetValue().(import5.UrlSource)

			urlSrcObj := make(map[string]interface{})
			urlSrcObjList := make([]map[string]interface{}, 0)

			if urlSrc.Url != nil {
				urlSrcObj["url"] = urlSrc.Url
			}
			if urlSrc.BasicAuth != nil {
				urlSrcObj["basic_auth"] = flattenURLBasicAuth(urlSrc.BasicAuth)
			}
			if urlSrc.ShouldAllowInsecureUrl != nil {
				urlSrcObj["should_allow_insecure_url"] = urlSrc.ShouldAllowInsecureUrl
			}

			urlSrcObjList = append(urlSrcObjList, urlSrcObj)

			urlSrcMap["url_source"] = urlSrcObjList
			urlSrcList = append(urlSrcList, urlSrcMap)

			return urlSrcList
		}

		if *pr.ObjectType_ == "vmm.v4.content.VmDiskSource" {
			vmDiskSrc := pr.GetValue().(import5.VmDiskSource)

			vmDiskObj := make(map[string]interface{})
			vmDiskObjList := make([]map[string]interface{}, 0)

			if vmDiskSrc.ExtId != nil {
				vmDiskObj["ext_id"] = vmDiskSrc.ExtId
			}

			vmDiskObjList = append(vmDiskObjList, vmDiskObj)

			vmDiskSrcMap["vm_disk_source"] = vmDiskObjList

			vmDiskSrcList = append(vmDiskSrcList, vmDiskSrcMap)

			return vmDiskSrcList
		}

		if *pr.ObjectType_ == "vmm.v4.content.ObjectsLiteSource" {
			objLiteSrc := pr.GetValue().(import5.ObjectsLiteSource)

			objLiteSrcObj := make(map[string]interface{})
			objLiteSrcObjList := make([]map[string]interface{}, 0)

			if objLiteSrc.Key != nil {
				objLiteSrcObj["key"] = objLiteSrc.Key
			}

			objLiteSrcObjList = append(objLiteSrcObjList, objLiteSrcObj)

			objectLiteSrc["object_lite_source"] = objLiteSrcObjList

			objectLiteSrcList = append(objectLiteSrcList, objectLiteSrc)

			return objectLiteSrcList
		}
	}
	return nil
}

func flattenURLBasicAuth(pr *import5.UrlBasicAuth) []map[string]interface{} {
	if pr != nil {
		auths := make([]map[string]interface{}, 0)

		auth := make(map[string]interface{})

		auth["username"] = pr.Username
		auth["password"] = pr.Password
		auths = append(auths, auth)
		return auths
	}
	return nil
}

func flattenImagePlacementStatus(pr []import5.ImagePlacementStatus) []interface{} {
	if len(pr) > 0 {
		imgList := make([]interface{}, len(pr))

		for k, v := range pr {
			img := make(map[string]interface{})

			img["placement_policy_ext_id"] = v.PlacementPolicyExtId
			img["compliance_status"] = v.ComplianceStatus
			img["enforcement_mode"] = v.EnforcementMode
			img["policy_cluster_ext_ids"] = v.PolicyClusterExtIds
			img["enforced_cluster_ext_ids"] = v.EnforcedClusterExtIds
			img["conflicting_policy_ext_ids"] = v.ConflictingPolicyExtIds

			imgList[k] = img
		}
		return imgList
	}
	return nil
}
