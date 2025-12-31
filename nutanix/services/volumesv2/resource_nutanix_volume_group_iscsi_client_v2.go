package volumesv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	taskPoll "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	config "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/config"
	volumesPrism "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/prism/v4/config"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Attach/Detach an iSCSI client to the given Volume Group.
func ResourceNutanixVolumeGroupIscsiClientV2() *schema.Resource {
	return &schema.Resource{
		Description:   "Attach iSCSI initiator to a Volume Group identified by {extId}",
		CreateContext: ResourceNutanixVolumeGroupIscsiClientV2Create,
		ReadContext:   ResourceNutanixVolumeGroupIscsiClientV2Read,
		UpdateContext: ResourceNutanixVolumeGroupIscsiClientV2Update,
		DeleteContext: ResourceNutanixVVolumeGroupIscsiClientV2Delete,
		Schema: map[string]*schema.Schema{
			"vg_ext_id": {
				Description: "The external identifier of the Volume Group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ext_id": {
				Description: "A globally unique identifier of an instance that is suitable for external consumption.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"iscsi_initiator_name": {
				Description: "iSCSI initiator name. During the attach operation, exactly one of iscsiInitiatorName and iscsiInitiatorNetworkID must be specified. This field is immutable.",
				Type:        schema.TypeString,
				Optional:    true,
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
										Description: "A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"client_secret": {
				Description: "iSCSI initiator client secret in case of CHAP authentication. This field should not be provided in case the authentication type is not set to CHAP..",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled_authentications": {
				Description:  "The authentication type enabled for the Volume Group. This is an optional field. If omitted, authentication is not configured for the Volume Group. If this is set to CHAP, the target/client secret must be provided.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CHAP", "NONE"}, false),
			},
			"num_virtual_targets": {
				Description: "Number of virtual targets generated for the iSCSI target. This field is immutable.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"attachment_site": {
				Description:  "The site where the Volume Group attach operation should be processed. This is an optional field. This field may only be set if Metro DR has been configured for this Volume Group.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"SECONDARY", "PRIMARY"}, false),
			},
		},
	}
}

func SchemaForIPV4ValuePrefixLength() *schema.Schema {
	return &schema.Schema{
		Description: "An unique address that identifies a device on the internet or a local network in IPv4 format.",
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Description: "An unique address that identifies a device on the internet or a local network in IPv4 format.",
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
				},
				"prefix_length": {
					Description: "The prefix length of the network to which this host IPv4 address belongs.",
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
				},
			},
		},
	}
}

func SchemaForIPV6ValuePrefixLength() *schema.Schema {
	return &schema.Schema{
		Description: "An unique address that identifies a device on the internet or a local network in IPv6 format.",
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Description: "An unique address that identifies a device on the internet or a local network in IPv6 format.",
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
				},
				"prefix_length": {
					Description: "The prefix length of the network to which this host IPv6 address belongs.",
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
				},
			},
		},
	}
}

// Attach an iSCSI client to the given Volume Group.
func ResourceNutanixVolumeGroupIscsiClientV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("vg_ext_id")

	body := volumesClient.IscsiClient{}

	if iscsiInitiatorName, ok := d.GetOk("iscsi_initiator_name"); ok {
		body.IscsiInitiatorName = utils.StringPtr(iscsiInitiatorName.(string))
	}
	if iscsiInitiatorNetworkID, ok := d.GetOk("iscsi_initiator_network_id"); ok {
		body.IscsiInitiatorNetworkId = expandiscsiInitiatorNetworkID(iscsiInitiatorNetworkID.([]interface{}))
	}
	if clientSecret, ok := d.GetOk("client_secret"); ok {
		body.ClientSecret = utils.StringPtr(clientSecret.(string))
	}
	if enabledAuthentications, ok := d.GetOk("enabled_authentications"); ok {
		const two, three = 2, 3
		enabledAuthenticationsMap := map[string]interface{}{
			"CHAP": two,
			"NONE": three,
		}
		pInt := enabledAuthenticationsMap[enabledAuthentications.(string)]
		p := volumesClient.AuthenticationType(pInt.(int))
		body.EnabledAuthentications = &p
	}
	if numVirtualTargets, ok := d.GetOk("num_virtual_targets"); ok {
		body.NumVirtualTargets = utils.IntPtr(numVirtualTargets.(int))
	}
	if attachmentSite, ok := d.GetOk("attachment_site"); ok {
		const two, three = 2, 3
		attachmentSiteMap := map[string]interface{}{
			"SECONDARY": two,
			"PRIMARY":   three,
		}
		pInt := attachmentSiteMap[attachmentSite.(string)]
		p := volumesClient.VolumeGroupAttachmentSite(pInt.(int))
		body.AttachmentSite = &p
	}

	resp, err := conn.VolumeAPIInstance.AttachIscsiClient(utils.StringPtr(volumeGroupExtID.(string)), &body)
	if err != nil {
		return diag.Errorf("error while Attaching Iscsi Client to Volume Group: %v", err)
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the iSCSI client to be attached to the Volume Group
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for iSCSI client (%s) to attach to Volume Group: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while Attaching Iscsi Client to Volume Group: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)

	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Attach Iscsi Client to Volume Group Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeIscsiClient, "iSCSI client")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(uuid))
	return nil
}

func ResourceNutanixVolumeGroupIscsiClientV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVolumeGroupIscsiClientV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// Detach an iSCSi client from the given Volume Group.
func ResourceNutanixVVolumeGroupIscsiClientV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("vg_ext_id")

	body := volumesClient.IscsiClientAttachment{}

	if extID, ok := d.GetOk("ext_id"); ok {
		body.ExtId = utils.StringPtr(extID.(string))
	}

	resp, err := conn.VolumeAPIInstance.DetachIscsiClient(utils.StringPtr(volumeGroupExtID.(string)), &body)
	if err != nil {
		return diag.Errorf("error while Detaching Iscsi Client to Volume Group: %v", err)
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the iSCSI client to be detached from the Volume Group
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for iSCSI client (%s) to detach from Volume Group: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while Detaching Iscsi Client to Volume Group: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)

	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Detach Iscsi Client from Volume Group Task Details: %s", string(aJSON))

	_, err = common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeIscsiClient, "iSCSI client")
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandiscsiInitiatorNetworkID(ipAddressOrFQDN interface{}) *config.IPAddressOrFQDN {
	if ipAddressOrFQDN != nil {
		fip := &config.IPAddressOrFQDN{}
		prI := ipAddressOrFQDN.([]interface{})
		val := prI[0].(map[string]interface{})

		if ipv4, ok := val["ipv4"]; ok {
			fip.Ipv4 = expandFloatingIPV4Address(ipv4)
		}
		if ipv6, ok := val["ipv6"]; ok {
			fip.Ipv6 = expandFloatingIPV6Address(ipv6)
		}
		if fqdn, ok := val["fqdn"]; ok {
			fip.Fqdn = expandFQDN(fqdn)
		}

		return fip
	}
	return nil
}

func expandFloatingIPV4Address(ipv4I interface{}) *config.IPv4Address {
	if ipv4I != nil {
		ipv4 := &config.IPv4Address{}
		prI := ipv4I.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv4.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv4.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv4
	}
	return nil
}

func expandFloatingIPV6Address(ipv6I interface{}) *config.IPv6Address {
	if ipv6I != nil {
		ipv6 := &config.IPv6Address{}
		prI := ipv6I.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv6.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv6.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv6
	}
	return nil
}

func expandFQDN(fqdnObj interface{}) *config.FQDN {
	if fqdnObj != nil {
		fqdn := &config.FQDN{}
		prI := fqdnObj.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			fqdn.Value = utils.StringPtr(value.(string))
		}
		return fqdn
	}
	return nil
}
