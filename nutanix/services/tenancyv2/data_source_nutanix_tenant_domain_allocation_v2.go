package tenancyv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixTenantDomainAllocationV2 reads a single domain allocation for a tenant.
func DatasourceNutanixTenantDomainAllocationV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixTenantDomainAllocationV2Read,
		Schema: map[string]*schema.Schema{
			"tenant_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_network_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     externalNetworkConnectionDataSourceElem(),
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     projectDataSourceElem(),
			},
			"state_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     stateInfoDataSourceElem(),
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
		},
	}
}

// projectDataSourceElem returns the fully-computed datasource schema for a project block.
func projectDataSourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaDataSourceElem(),
			},
			"resource_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"placement_targets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"storage_container_uuids": {
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
				},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// externalNetworkConnectionDataSourceElem returns the fully-computed datasource schema
// for an external network connection block.
func externalNetworkConnectionDataSourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     externalNetworkConnConfigDataSourceElem(),
			},
			"state_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     stateInfoDataSourceElem(),
			},
		},
	}
}

// externalNetworkConnConfigDataSourceElem returns the fully-computed datasource schema
// for the NAT / no-NAT external network connection configuration.
func externalNetworkConnConfigDataSourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nat": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transit_vpc_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connections": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"floating_ip_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"overlay_subnet_ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"no_nat": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"overlay_subnet_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_vpc_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"externally_routable_prefixes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": ipSubnetDataSourceSchema(),
									"ipv6": ipSubnetDataSourceSchema(),
								},
							},
						},
					},
				},
			},
		},
	}
}

// ipSubnetDataSourceSchema returns the fully-computed datasource schema for an IP subnet block.
func ipSubnetDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"prefix_length": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"ip": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"value": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"prefix_length": {
								Type:     schema.TypeInt,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

// DatasourceNutanixTenantDomainAllocationV2Read fetches a domain allocation and populates state.
func DatasourceNutanixTenantDomainAllocationV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	tenantExtID := d.Get("tenant_ext_id").(string)
	extID := d.Get("ext_id").(string)
	getRequest := import4.GetTenantDomainAllocationByIdRequest{
		TenantExtId: utils.StringPtr(tenantExtID),
		ExtId:       utils.StringPtr(extID),
	}
	resp, err := conn.TenantManagementAPIInstance.GetTenantDomainAllocationById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant domain allocation: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.DomainAllocation)

	if err := d.Set("domain_ext_id", getResp.DomainExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("external_network_config", flattenExternalNetworkConnection(getResp.ExternalNetworkConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("projects", flattenProjects(getResp.Projects)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state_info", flattenDomainAllocationStateInfo(getResp.StateInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
