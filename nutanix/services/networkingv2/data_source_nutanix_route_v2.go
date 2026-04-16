package networkingv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixRouteV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRouteV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"route_table_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DatasourceMetadataSchemaV4(),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValueRequiredPrefixLengthRequired(),
						"ipv6": SchemaForValueRequiredPrefixLengthRequired(),
					},
				},
			},
			"next_hop": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"next_hop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_ip_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValueRequiredPrefixLength(),
									"ipv6": SchemaForValueRequiredPrefixLength(),
								},
							},
						},
						"next_hop_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"route_table_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_routing_domain_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixRouteV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DatasourceNutanixRouteV2Read")

	conn := meta.(*conns.Client).NetworkingAPI

	routeTableExtID := d.Get("route_table_ext_id").(string)
	extID := d.Get("ext_id").(string)

	resp, err := conn.Routes.GetRouteForRouteTableById(&extID, &routeTableExtID)
	if err != nil {
		return diag.Errorf("error while fetching route : %v", err)
	}

	getResp := resp.Data.GetValue().(config.Route)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("destination", flattenDestination(getResp.Destination)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("next_hop", flattenNextHop(getResp.Nexthop)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("route_table_reference", getResp.RouteTableReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc_reference", getResp.VpcReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("external_routing_domain_reference", getResp.ExternalRoutingDomainReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("route_type", flattenRouteType(getResp.RouteType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_active", getResp.IsActive); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("priority", getResp.Priority); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))

	return nil
}

func DatasourceMetadataSchemaV4() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"owner_reference_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"owner_user_name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"project_reference_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"project_name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"category_ids": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func flattenDestination(destination *config.IPSubnet) interface{} {
	if destination == nil {
		return nil
	}
	destinationMap := make(map[string]interface{})
	if destination.Ipv4 != nil {
		destinationMap["ipv4"] = flattenIPv4Subnet(destination.Ipv4)
	}
	if destination.Ipv6 != nil {
		destinationMap["ipv6"] = flattenIPv6Subnet(destination.Ipv6)
	}
	return []interface{}{destinationMap}
}

func flattenNextHop(nextHops *config.Nexthop) interface{} {
	if nextHops != nil {
		nextHop := make(map[string]interface{})
		aJSON, _ := json.Marshal(nextHops)
		log.Printf("[DEBUG] NextHops: %s", string(aJSON))
		if nextHops.NexthopType != nil {
			nextHop["next_hop_type"] = flattenNextHopType(nextHops.NexthopType)
		}
		if nextHops.NexthopReference != nil {
			nextHop["next_hop_reference"] = nextHops.NexthopReference
		}
		if nextHops.NexthopIpAddress != nil {
			nextHop["next_hop_ip_address"] = flattenIPAddress(nextHops.NexthopIpAddress)
		}
		if nextHops.NexthopName != nil {
			nextHop["next_hop_name"] = nextHops.NexthopName
		}
		return []interface{}{nextHop}
	}
	return nil
}

func flattenNextHopType(nextHopType *config.NexthopType) string {
	if nextHopType != nil {
		const two, three, four, five, six = 2, 3, 4, 5, 6
		if *nextHopType == config.NexthopType(two) {
			return "IP_ADDRESS"
		}
		if *nextHopType == config.NexthopType(three) {
			return "DIRECT_CONNECT_VIF"
		}
		if *nextHopType == config.NexthopType(four) {
			return "LOCAL_SUBNET"
		}
		if *nextHopType == config.NexthopType(five) {
			return "EXTERNAL_SUBNET"
		}
		if *nextHopType == config.NexthopType(six) {
			return "VPN_CONNECTION"
		}
	}
	return "UNKNOWN"
}

func flattenRouteType(routeType *config.RouteType) interface{} {
	if routeType != nil {
		const two, three, four = 2, 3, 4
		if *routeType == config.RouteType(two) {
			return "DYNAMIC"
		}
		if *routeType == config.RouteType(three) {
			return "LOCAL"
		}
		if *routeType == config.RouteType(four) {
			return "STATIC"
		}
	}
	return "UNKNOWN"
}
