package volumesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Fetch an iSCSI client details.
func DatasourceNutanixVolumeIscsiClientV2() *schema.Resource {
	return &schema.Resource{
		Description: "Fetches the iSCSI client details identified by {extId}.",
		ReadContext: DatasourceNutanixVolumeIscsiClientV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "The external identifier of the iSCSI client.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tenant_id": {
				Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server)",
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
			"iscsi_initiator_name": {
				Description: "iSCSI initiator name. During the attach operation, exactly one of iscsiInitiatorName and iscsiInitiatorNetworkId must be specified. This field is immutable.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"iscsi_initiator_network_id": {
				Description: "An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForIPV4ValuePrefixLength(),
						"ipv6": SchemaForIPV6ValuePrefixLength(),
						"fqdn": {
							Description: "A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.",
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"enabled_authentications": {
				Description: "The authentication type enabled for the Volume Group. This is an optional field. If omitted, authentication is not configured for the Volume Group. If this is set to CHAP, the target/client secret must be provided.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"attached_targets": {
				Description: "with each iSCSI target corresponding to the iSCSI client)",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"num_virtual_targets": {
							Description: "Number of virtual targets generated for the iSCSI target. This field is immutable.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"iscsi_target_name": {
							Description: "Name of the iSCSI target that the iSCSI client is connected to. This is a read-only field.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"cluster_reference": {
				Description: "The UUID of the cluster that will host the iSCSI client. This field is read-only.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"attachment_site": {
				Description: "The site where the Volume Group attach operation should be processed. This is an optional field. This field may only be set if Metro DR has been configured for this Volume Group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func DatasourceNutanixVolumeIscsiClientV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	extID := d.Get("ext_id")

	// get the volume group iscsi clients
	resp, err := conn.IscsiClientAPIInstance.GetIscsiClientById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching Iscsi Client : %v", err)
	}

	getResp := resp.Data.GetValue().(volumesClient.IscsiClient)

	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("iscsi_initiator_name", getResp.IscsiInitiatorName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("iscsi_initiator_network_id", flattenIscsiInitiatorNetworkID(getResp.IscsiInitiatorNetworkId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled_authentications", getResp.EnabledAuthentications); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("attached_targets", flattenAttachedTargets(getResp.AttachedTargets)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_reference", getResp.ClusterReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("attachment_site", getResp.AttachmentSite); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
