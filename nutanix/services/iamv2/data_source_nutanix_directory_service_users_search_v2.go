package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	import2 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/request/directoryservices"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixDirectoryServiceUsersSearchV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixDirectoryServiceUsersSearchV2Read,
		Schema: map[string]*schema.Schema{
			"directory_service_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "External identifier of the directory service to search.",
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Query string for directory service search.",
			},
			"is_wildcard_search": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Flag indicating whether the search should be a wildcard search or not. Defaults to true.",
			},
			"searched_attributes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Attributes for search operation. By default, the search will be performed with a common name.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"returned_attributes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Attributes returned by the search operation.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain name for the directory service.",
			},
			"search_results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of search result entities.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of entity, either user or group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the entity in canonical format.",
						},
						"identity_ext_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "External identifier of the identity (if exists) in the directory service.",
						},
						"attributes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of attributes for the search entity.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the attribute.",
									},
									"values": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Values of the attribute.",
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
		},
	}
}

func datasourceNutanixDirectoryServiceUsersSearchV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	extID := d.Get("directory_service_ext_id").(string)
	query := d.Get("query").(string)

	searchQuery := import1.NewDirectoryServiceSearchQuery()
	searchQuery.Query = utils.StringPtr(query)

	if common.IsExplicitlySet(d, "is_wildcard_search") {
		searchQuery.IsWildcardSearch = utils.BoolPtr(d.Get("is_wildcard_search").(bool))
	}
	if v, ok := d.GetOk("returned_attributes"); ok {
		searchQuery.ReturnedAttributes = common.ExpandListOfString(v.([]interface{}))
	}

	req := &import2.SearchDirectoryServiceRequest{
		ExtId: utils.StringPtr(extID),
		Body:  searchQuery,
	}
	resp, err := conn.DirectoryServiceAPIInstance.SearchDirectoryService(ctx, req)
	if err != nil {
		return diag.Errorf("error searching directory service (%s): %v", extID, err)
	}
	if resp.Data != nil {
		result := resp.Data.GetValue().(import1.DirectoryServiceSearchResult)
		if err := d.Set("domain_name", result.DomainName); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("search_results", flattenDirectoryServiceSearchResults(result.SearchResults)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenDirectoryServiceSearchResults(entities []import1.DirectoryServiceSearchEntity) []map[string]interface{} {
	if len(entities) == 0 {
		return []map[string]interface{}{}
	}
	results := make([]map[string]interface{}, len(entities))
	for i, entity := range entities {
		results[i] = map[string]interface{}{
			"entity_type":     entity.EntityType,
			"name":            entity.Name,
			"identity_ext_id": entity.IdentityExtId,
			"attributes":      flattenDirectoryServiceSearchAttributes(entity.Attributes),
		}
	}
	return results
}

func flattenDirectoryServiceSearchAttributes(attrs []import1.DirectoryServiceSearchAttribute) []map[string]interface{} {
	if len(attrs) == 0 {
		return []map[string]interface{}{}
	}
	results := make([]map[string]interface{}, len(attrs))
	for i, attr := range attrs {
		results[i] = map[string]interface{}{
			"name":   attr.Name,
			"values": attr.Values,
		}
	}
	return results
}
