package nutanix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixUserRead,
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
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
			"categories": categoriesSchema(),
			"owner_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"directory_service_user": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_principal_name": {
							Type:     schema.TypeString,
							Computed: true,
							//ValidateFunc: validation.StringInSlice([]string{"role"}, false),
						},
						"directory_service_reference": {
							Type:     schema.TypeList,
							MaxItems: 1,
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
						"default_user_principal_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"identity_provider_user": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identity_provider_reference": {
							Type:     schema.TypeList,
							MaxItems: 1,
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
					},
				},
			},
			"user_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_reference_list": {
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
			"access_control_policy_reference_list": {
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixUserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading User: %s", d.Id())

	// Get client connection
	conn := meta.(*Client).API
	uuid := d.Get("uuid").(string)

	// Make request to the API
	resp, err := conn.V3.GetUser(uuid)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return fmt.Errorf("error reading user UUID (%s) with error %s", uuid, err)
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err = d.Set("metadata", m); err != nil {
		return fmt.Errorf("error setting metadata for image UUID(%s), %s", d.Id(), err)
	}
	if err = d.Set("categories", c); err != nil {
		return fmt.Errorf("error setting categories for image UUID(%s), %s", d.Id(), err)
	}

	if err = d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return fmt.Errorf("error setting owner_reference for image UUID(%s), %s", d.Id(), err)
	}
	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))

	if err = d.Set("state", resp.Status.State); err != nil {
		return fmt.Errorf("error setting state for image UUID(%s), %s", d.Id(), err)
	}

	if err = d.Set("directory_service_user", flattenDirectoryServiceUser(resp.Status.Resources.DirectoryServiceUser)); err != nil {
		return fmt.Errorf("error setting state for image UUID(%s), %s", d.Id(), err)
	}

	if err = d.Set("identity_provider_user", flattenIdentityProviderUser(resp.Status.Resources.IdentityProviderUser)); err != nil {
		return fmt.Errorf("error setting state for image UUID(%s), %s", d.Id(), err)
	}

	if err = d.Set("user_type", resp.Status.Resources.UserType); err != nil {
		return fmt.Errorf("error setting state for image UUID(%s), %s", d.Id(), err)
	}

	if err = d.Set("display_name", resp.Status.Resources.DisplayName); err != nil {
		return fmt.Errorf("error setting state for image UUID(%s), %s", d.Id(), err)
	}

	if err := d.Set("project_reference_list", flattenArrayReferenceValues(resp.Status.Resources.ProjectsReferenceList)); err != nil {
		return fmt.Errorf("error setting state for image UUID(%s), %s", d.Id(), err)
	}

	refe := flattenArrayReferenceValues(resp.Status.Resources.AccessControlPolicyReferenceList)
	utils.PrintToJSON(refe, "acceess")

	if err := d.Set("access_control_policy_reference_list", refe); err != nil {
		return fmt.Errorf("error setting state for image UUID(%s), %s", d.Id(), err)
	}

	d.SetId(uuid)

	return nil
}
