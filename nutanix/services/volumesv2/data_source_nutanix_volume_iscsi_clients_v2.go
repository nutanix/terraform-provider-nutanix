package volumesv2

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/config"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// List all the iSCSI clients.
func DatasourceNutanixVolumeIscsiClientsV2() *schema.Resource {
	return &schema.Resource{
		Description: "Fetches the list of iSCSI clients.",
		ReadContext: DatasourceNutanixVolumeIscsiClientsV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Description: "A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"limit": {
				Description: "A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"filter": {
				Description: "A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions. For example, filter '$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields: clusterReference, extId",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"orderby": {
				Description: "A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields: clusterReference, extId",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"expand": {
				Description: "A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved. Each expanded item is evaluated relative to the entity containing the property being expanded. Other query options can be applied to an expanded property by appending a semicolon-separated list of query options, enclosed in parentheses, to the property name. Permissible system query options are $filter, $select and $orderby. The following expansion keys are supported. The expand can be applied to the following fields: cluster",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"select": {
				Description: "A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. If a $select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. The select can be applied to the following fields: clusterReference, extId",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"iscsi_clients": {
				Description: "List of iSCSI clients.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"ext_id": {
							Description: "A globally unique identifier of an instance that is suitable for external consumption.",
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
									"ipv4": SchemaForIpV4ValuePrefixLength(),
									"ipv6": SchemaForIpV6ValuePrefixLength(),
									"fqdn": {
										Description: "A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.",
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Description: "The fully qualified domain name.",
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
						"enabled_authentications": {
							Description: "The authentication type enabled for the Volume Group. This is an optional field. If omitted, authentication is not configured for the Volume Group. If this is set to CHAP, the target/client secret must be provided.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"attached_targets": {
							Description: "associated with each iSCSI target corresponding to the iSCSI client)",
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
				},
			},
		},
	}
}

func DatasourceNutanixVolumeIscsiClientsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	var filter, orderBy, expand, selects *string
	var page, limit *int

	// initialize the query parameters
	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if expandf, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandf.(string))
	} else {
		expand = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	// get the volume group iscsi clients
	resp, err := conn.IscsiClientAPIInstance.ListIscsiClients(page, limit, filter, orderBy, expand, selects)

	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching Iscsi Clients : %v", errorMessage["message"])
	}

	// // Check if resp is nil before accessing its data
	// if resp != nil {
	// 	diskResp := resp.Data

	// 	// extract the volume groups data from the response
	// 	if diskResp != nil {

	// 		// set the volume groups iscsi clients  data in the terraform resource
	// 		if err := d.Set("iscsi_clients", flattenIscsiClientsEntities(diskResp.GetValue().([]volumesClient.IscsiClient))); err != nil {
	// 			return diag.FromErr(err)
	// 		}
	// 	}
	// }

	iscsiClientsResp := resp.Data

	// extract the volume groups data from the response
	if iscsiClientsResp != nil {

		// set the volume groups iscsi clients  data in the terraform resource
		if err := d.Set("iscsi_clients", flattenIscsiClientsEntities(iscsiClientsResp.GetValue().([]volumesClient.IscsiClient))); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(resource.UniqueId())
	return nil

}

func flattenIscsiClientsEntities(pr []volumesClient.IscsiClient) []interface{} {
	if len(pr) > 0 {
		iscsi_clients := make([]interface{}, len(pr))

		for k, v := range pr {
			iscsi_client := make(map[string]interface{})

			if v.TenantId != nil {
				iscsi_client["tenant_id"] = v.TenantId
			}
			if v.ExtId != nil {
				iscsi_client["ext_id"] = v.ExtId
			}
			if v.Links != nil {
				iscsi_client["links"] = flattenLinks(v.Links)
			}
			if v.IscsiInitiatorName != nil {
				iscsi_client["iscsi_initiator_name"] = v.IscsiInitiatorName
			}
			if v.IscsiInitiatorNetworkId != nil {
				iscsi_client["iscsi_initiator_network_id"] = flattenIscsiInitiatorNetworkId(v.IscsiInitiatorNetworkId)
			}
			if v.EnabledAuthentications != nil {
				iscsi_client["enabled_authentications"] = flattenEnabledAuthentications(v.EnabledAuthentications)
			}
			if v.AttachedTargets != nil {
				iscsi_client["attached_targets"] = flattenAttachedTargets(v.AttachedTargets)
			}
			if v.AttachmentSite != nil {
				iscsi_client["attachment_site"] = flattenAttachmentSite(v.AttachmentSite)
			}
			if v.ClusterReference != nil {
				iscsi_client["cluster_reference"] = v.ClusterReference
			}
			// Attribute not present in the response of GA SDK
			// if v.TargetParams != nil {
			// 	iscsi_client["attached_targets"] = flattenAttachedTargets(v.TargetParams)
			// }

			iscsi_clients[k] = iscsi_client

		}
		return iscsi_clients
	}
	return nil
}

func flattenAttachmentSite(iscsiClientAttachmentSite *volumesClient.VolumeGroupAttachmentSite) string {
	const two, three = 2, 3
	if iscsiClientAttachmentSite != nil {
		if *iscsiClientAttachmentSite == volumesClient.VolumeGroupAttachmentSite(two) {
			return "PRIMARY"
		}
		if *iscsiClientAttachmentSite == volumesClient.VolumeGroupAttachmentSite(two) {
			return "SECONDARY"
		}
	}
	return "UNKNOWN"
}

func flattenAttachedTargets(targetParam []volumesClient.TargetParam) []interface{} {
	if len(targetParam) > 0 {
		targetParamList := make([]interface{}, len(targetParam))
		for k, v := range targetParam {
			target := make(map[string]interface{})

			if v.NumVirtualTargets != nil {
				target["num_virtual_targets"] = v.NumVirtualTargets
			}
			if v.IscsiTargetName != nil {
				target["iscsi_target_name"] = v.IscsiTargetName
			}
			targetParamList[k] = target
		}
		return targetParamList
	}
	return nil
}

func flattenIscsiInitiatorNetworkId(iPAddressOrFQDN *config.IPAddressOrFQDN) []interface{} {
	if iPAddressOrFQDN != nil {
		ipAddressOrFQDN := make(map[string]interface{})
		if iPAddressOrFQDN.Ipv4 != nil {
			ipAddressOrFQDN["ipv4"] = flattenIp4Address(iPAddressOrFQDN.Ipv4)
		}
		if iPAddressOrFQDN.Ipv6 != nil {
			ipAddressOrFQDN["ipv6"] = flattenIp6Address(iPAddressOrFQDN.Ipv6)
		}
		if iPAddressOrFQDN.Fqdn != nil {
			ipAddressOrFQDN["fqdn"] = flattenFQDN(iPAddressOrFQDN.Fqdn)
		}
		return []interface{}{ipAddressOrFQDN}
	}
	return nil
}

func flattenIp6Address(iPv6Address *config.IPv6Address) []interface{} {
	if iPv6Address != nil {
		ipv6 := make([]interface{}, 0)

		ip := make(map[string]interface{})

		ip["value"] = iPv6Address.Value
		ip["prefix_length"] = iPv6Address.PrefixLength

		ipv6 = append(ipv6, ip)

		return ipv6
	}
	return nil
}

func flattenIp4Address(iPv4Address *config.IPv4Address) []interface{} {
	if iPv4Address != nil {
		ipv4 := make([]interface{}, 0)

		ip := make(map[string]interface{})

		ip["value"] = iPv4Address.Value
		ip["prefix_length"] = iPv4Address.PrefixLength

		ipv4 = append(ipv4, ip)

		return ipv4
	}
	return nil
}

func flattenFQDN(fQDN *config.FQDN) []interface{} {
	if fQDN != nil {
		fqdn := make([]interface{}, 0)

		ip := make(map[string]interface{})

		ip["value"] = fQDN.Value
		fqdn = append(fqdn, ip)
		return fqdn
	}
	return nil
}

// func flattenValuePrefixLength(iPv4Address *config.IPv4Address) {
// 	panic("unimplemented")
// }
