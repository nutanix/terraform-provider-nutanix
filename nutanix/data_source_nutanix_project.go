package nutanix

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixProjectRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_hash": {
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
			"project_reference": {
				Type:     schema.TypeMap,
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
			"owner_reference": {
				Type:     schema.TypeMap,
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
			"categories": categoriesSchema(),
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixProjectRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	projectID := d.Get("project_id").(string)

	project, err := conn.V3.GetProject(projectID)
	if err != nil {
		return err
	}

	m, c := setRSEntityMetadata(project.Metadata)

	if err := d.Set("name", project.Status.Name); err != nil {
		return fmt.Errorf("error setting `name` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("description", project.Status.Descripion); err != nil {
		return fmt.Errorf("error setting `description` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("state", project.Status.State); err != nil {
		return fmt.Errorf("error setting `state` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("is_default", project.Status.Resources.IsDefault); err != nil {
		return fmt.Errorf("error setting `is_default` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("resource_domain", flattenResourceDomain(project.Spec.Resources.ResourceDomain)); err != nil {
		return fmt.Errorf("error setting `resource_domain` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("account_reference_list", flattenReferenceList(project.Spec.Resources.AccountReferenceList)); err != nil {
		return fmt.Errorf("error setting `account_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("environment_reference_list", flattenReferenceList(project.Spec.Resources.EnvironmentReferenceList)); err != nil {
		return fmt.Errorf("error setting `environment_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("default_subnet_reference", flattenReference(project.Spec.Resources.DefaultSubnetReference)); err != nil {
		return fmt.Errorf("error setting `default_subnet_reference` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("user_reference_list", flattenReferenceList(project.Spec.Resources.UserReferenceList)); err != nil {
		return fmt.Errorf("error setting `user_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("external_user_group_reference_list",
		flattenReferenceList(project.Spec.Resources.ExternalUserGroupReferenceList)); err != nil {
		return fmt.Errorf("error setting `external_user_group_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("subnet_reference_list", flattenReferenceList(project.Spec.Resources.SubnetReferenceList)); err != nil {
		return fmt.Errorf("error setting `subnet_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("external_network_list", flattenReferenceList(project.Spec.Resources.ExternalNetworkList)); err != nil {
		return fmt.Errorf("error setting `external_network_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("metadata", m); err != nil {
		return fmt.Errorf("error setting `metadata` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("project_reference", flattenReferenceValues(project.Metadata.ProjectReference)); err != nil {
		return fmt.Errorf("error setting `project_reference` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("owner_reference", flattenReferenceValues(project.Metadata.OwnerReference)); err != nil {
		return fmt.Errorf("error setting `owner_reference` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("categories", c); err != nil {
		return fmt.Errorf("error setting `categories` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("api_version", project.APIVersion); err != nil {
		return fmt.Errorf("error setting `api_version` for Project(%s): %s", d.Id(), err)
	}

	d.SetId(projectID)

	return nil
}
