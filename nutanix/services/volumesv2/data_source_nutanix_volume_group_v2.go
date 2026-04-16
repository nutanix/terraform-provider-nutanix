package volumesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Get a Volume Group.
func DatasourceNutanixVolumeGroupV2() *schema.Resource {
	return &schema.Resource{
		Description: "Query the Volume Group identified by {extId}.",
		ReadContext: DatasourceNutanixVolumeGroupV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "The external identifier of the Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tenant_id": {
				Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"links": {
				Description: "A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Description: "The URL at which the entity described by the link can be accessed.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"rel": {
							Description: "A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of \"self\" identifies the URL for the object.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"name": {
				Description: "Volume Group name. This is an optional field.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Volume Group description. This is an optional field.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"should_load_balance_vm_attachments": {
				Description: "Indicates whether to enable Volume Group load balancing for VM attachments. This cannot be enabled if there are iSCSI client attachments already associated with the Volume Group, and vice-versa. This is an optional field.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"sharing_status": {
				Description: "Indicates whether the Volume Group can be shared across multiple iSCSI initiators. The mode cannot be changed from SHARED to NOT_SHARED on a Volume Group with multiple attachments. Similarly, a Volume Group cannot be associated with more than one attachment as long as it is in exclusive mode. This is an optional field. Possible values [SHARED, NOT_SHARED]",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"target_name": {
				Description: "Name of the external client target that will be visible and accessible to the client. This is an optional field.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enabled_authentications": {
				Description: "The authentication type enabled for the Volume Group. This is an optional field. If omitted, authentication is not configured for the Volume Group. If this is set to CHAP, the target/client secret must be provided.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"iscsi_features": {
				Description: "iSCSI specific settings for the Volume Group. This is an optional field.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled_authentications": {
							Description: "The authentication type enabled for the Volume Group. This is an optional field. If omitted, authentication is not configured for the Volume Group. If this is set to CHAP, the target/client secret must be provided.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"created_by": {
				Description: "Service/user who created this Volume Group. This is an optional field.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cluster_reference": {
				Description: "The UUID of the cluster that will host the Volume Group. This is a mandatory field for creating a Volume Group on Prism Central.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"storage_features": {
				Description: "Storage optimization features which must be enabled on the Volume Group. This is an optional field.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flash_mode": {
							Description: "Once configured, this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Description: "Indicates whether the flash mode is enabled for the Volume Group.",
										Type:        schema.TypeBool,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"usage_type": {
				Description: "Expected usage type for the Volume Group. This is an indicative hint on how the caller will consume the Volume Group. This is an optional",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_hidden": {
				Description: "Indicates whether the Volume Group is meant to be hidden or not. This is an optional field. If omitted, the VG will not be hidden.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func DatasourceNutanixVolumeGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	extID := d.Get("ext_id")

	resp, err := conn.VolumeAPIInstance.GetVolumeGroupById(utils.StringPtr(extID.(string)), nil)
	if err != nil {
		return diag.Errorf("error while fetching Volume Group : %v", err)
	}

	getResp := resp.Data.GetValue().(volumesClient.VolumeGroup)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("should_load_balance_vm_attachments", getResp.ShouldLoadBalanceVmAttachments); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("sharing_status", flattenSharingStatus(getResp.SharingStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("target_name", getResp.TargetName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled_authentications", flattenEnabledAuthentications(getResp.EnabledAuthentications)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("iscsi_features", flattenIscsiFeatures(getResp.IscsiFeatures)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_reference", getResp.ClusterReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_features", flattenStorageFeatures(getResp.StorageFeatures)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("usage_type", flattenUsageType(getResp.UsageType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_hidden", getResp.IsHidden); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
