package lifecyclev2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixPatchedImageV2 reads details of a single patched image by its
// external ID. It exposes the claim_token_ext_id cross-resource reference.
func DatasourceNutanixPatchedImageV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixPatchedImageV2Read,
		Schema:      patchedImageDatasourceSchema(),
	}
}

func patchedImageDatasourceSchema() map[string]*schema.Schema {
	valuePrefix := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value":         {Type: schema.TypeString, Computed: true},
			"prefix_length": {Type: schema.TypeInt, Computed: true},
		},
	}
	ipAddress := &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": {Type: schema.TypeList, Computed: true, Elem: valuePrefix},
				"ipv6": {Type: schema.TypeList, Computed: true, Elem: valuePrefix},
			},
		},
	}
	ipOrFqdn := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ipv4": {Type: schema.TypeList, Computed: true, Elem: valuePrefix},
			"ipv6": {Type: schema.TypeList, Computed: true, Elem: valuePrefix},
			"fqdn": {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{"value": {Type: schema.TypeString, Computed: true}},
			}},
		},
	}
	return map[string]*schema.Schema{
		"ext_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "A globally unique identifier of an instance that is suitable for external consumption.",
		},
		"claim_token_ext_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ExtId of the claim token used in the patched image",
		},
		"name":      {Type: schema.TypeString, Computed: true},
		"host_type": {Type: schema.TypeString, Computed: true},
		"version":   {Type: schema.TypeString, Computed: true},
		"image_details": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"local_host_image_ext_id": {Type: schema.TypeString, Computed: true},
				},
			},
		},
		"node_configurations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"node_ext_id": {Type: schema.TypeString, Computed: true},
					"host_configuration": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"hostname":    {Type: schema.TypeString, Computed: true},
								"nameservers": {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{"ipv4": {Type: schema.TypeList, Computed: true, Elem: valuePrefix}, "ipv6": {Type: schema.TypeList, Computed: true, Elem: valuePrefix}}}},
								"ntpservers":  {Type: schema.TypeList, Computed: true, Elem: ipOrFqdn},
								"network_details": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"management": {
												Type:     schema.TypeList,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"ip":        ipAddress,
														"gateway":   ipAddress,
														"vlan_id":   {Type: schema.TypeInt, Computed: true},
														"mtu_bytes": {Type: schema.TypeInt, Computed: true},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"patched_iso_url":             {Type: schema.TypeString, Computed: true},
		"patched_iso_sha256_checksum": {Type: schema.TypeString, Computed: true},
		"created_time":                {Type: schema.TypeString, Computed: true},
		"owner_ext_id":                {Type: schema.TypeString, Computed: true},
		"tenant_id":                   {Type: schema.TypeString, Computed: true},
		"links":                       datasourceLifecycleLinksSchema(),
	}
}

func DatasourceNutanixPatchedImageV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.PatchedImagesAPIInstance.GetPatchedImageById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching patched image : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.PatchedImage)

	if err := d.Set("claim_token_ext_id", getResp.ClaimTokenExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if getResp.HostType != nil {
		if err := d.Set("host_type", getResp.HostType.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("version", getResp.Version); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("image_details", flattenImageDetails(getResp.ImageDetails)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_configurations", flattenNodeConfigurations(getResp.NodeConfigurations)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("patched_iso_url", getResp.PatchedIsoUrl); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("patched_iso_sha256_checksum", getResp.PatchedIsoSha256Checksum); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.Format(time.RFC3339)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLifecycleLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
