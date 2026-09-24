package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/request/nicprofiles"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixNicProfileV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceNutanixNicProfileV2Read,
		Schema:      nicProfileSchema(true),
	}
}

func DataSourceNutanixNicProfileV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	getRequest := import2.GetNicProfileByIdRequest{
		ExtId: utils.StringPtr(d.Get("ext_id").(string)),
	}
	resp, err := conn.NicProfilesAPI.GetNicProfileById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching NIC profile: %v", err)
	}

	raw := resp.Data.GetValue()
	var getResp import1.NicProfile
	switch v := raw.(type) {
	case import1.NicProfile:
		getResp = v
	case *import1.NicProfile:
		if v == nil {
			return diag.Errorf("NIC profile response was nil")
		}
		getResp = *v
	default:
		return diag.Errorf("unexpected NIC profile response type: %T", raw)
	}

	if err := d.Set("name", utils.StringValue(getResp.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", utils.StringValue(getResp.Description)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("capability_config", flattenNicProfileCapabilityConfig(getResp.CapabilityConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_nic_references", flattenNicProfileHostNicReferences(getResp.HostNicReferences)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nic_family", utils.StringValue(getResp.NicFamily)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("operating_mode", common.FlattenPtrEnum(getResp.OperatingMode)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_type", common.FlattenPtrEnum(getResp.OwnerType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_ext_id", utils.StringValue(getResp.ProjectExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(getResp.TenantId)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
