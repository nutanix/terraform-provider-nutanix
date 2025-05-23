package prism

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixProjectsRead,
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"resource_domain": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resources": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"units": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"limit": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"resource_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"account_reference_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"environment_reference_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"default_subnet_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"user_reference_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"external_user_group_reference_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"subnet_reference_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"external_network_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"project_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"owner_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"categories": categoriesSchema(),
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixProjectsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection
	conn := meta.(*conns.Client).API

	req := &v3.DSMetadata{}

	metadata, filtersOk := d.GetOk("metadata")
	if filtersOk {
		req = buildDataSourceListMetadata(metadata.(*schema.Set))
	}

	resp, err := conn.V3.ListAllProject(utils.StringValue(req.Filter))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entities", flattenProjectEntities(resp.Entities)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenProjectEntities(projects []*v3.Project) []map[string]interface{} {
	entities := make([]map[string]interface{}, len(projects))

	for i, project := range projects {
		metadata, categories := setRSEntityMetadata(project.Metadata)

		entities[i] = map[string]interface{}{
			"name":                               project.Status.Name,
			"description":                        project.Status.Descripion,
			"state":                              project.Status.State,
			"is_default":                         project.Status.Resources.IsDefault,
			"resource_domain":                    flattenResourceDomain(project.Status.Resources.ResourceDomain),
			"account_reference_list":             flattenReferenceList(project.Spec.Resources.AccountReferenceList),
			"environment_reference_list":         flattenReferenceList(project.Spec.Resources.EnvironmentReferenceList),
			"default_subnet_reference":           flattenReference(project.Spec.Resources.DefaultSubnetReference),
			"user_reference_list":                flattenReferenceList(project.Spec.Resources.UserReferenceList),
			"external_user_group_reference_list": flattenReferenceList(project.Spec.Resources.ExternalUserGroupReferenceList),
			"subnet_reference_list":              flattenReferenceList(project.Spec.Resources.SubnetReferenceList),
			"external_network_list":              flattenExternalNetworkListReferenceList(project.Spec.Resources.ExternalNetworkList),
			"metadata":                           metadata,
			"categories":                         categories,
			"project_reference":                  flattenReferenceValues(project.Metadata.ProjectReference),
			"owner_reference":                    flattenReferenceValues(project.Metadata.OwnerReference),
			"api_version":                        project.APIVersion,
		}
	}
	return entities
}
