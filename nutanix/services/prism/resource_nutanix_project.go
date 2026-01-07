package prism

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/spf13/cast"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixProjectCreate,
		ReadContext:   resourceNutanixProjectRead,
		UpdateContext: resourceNutanixProjectUpdate,
		DeleteContext: resourceNutanixProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Update: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Delete: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
		},
		CustomizeDiff: customizeDiffProjectACP,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_collab": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_domain": {
				Type:       schema.TypeList,
				Optional:   true,
				MaxItems:   1,
				Deprecated: "Deprecated since v2.4.0. Prism Central no longer supports `resource_domain` for projects; remove this block from your configuration/scripts.",
				// `resource_domain` is no longer supported by Prism Central, but we keep it in the schema to avoid
				// breaking existing customer configurations. We also suppress diffs so it does not trigger updates.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return true },
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:     schema.TypeList,
							Required: true,
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
										Required: true,
									},
									"resource_type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"account_reference_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Default:  "account",
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"environment_reference_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "environment",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"default_subnet_reference": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "subnet",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"user_reference_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "user",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"external_user_group_reference_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "user_group",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"subnet_reference_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "subnet",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return schema.HashString(m["uuid"].(string))
				},
			},
			"external_network_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"tunnel_reference_list": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"use_project_internal"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "tunnel",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"cluster_reference_list": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"use_project_internal"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "cluster",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"vpc_reference_list": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"use_project_internal"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "vpc",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"default_environment_reference": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"use_project_internal"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "environment",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
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
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": categoriesSchema(),
			"api_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"use_project_internal": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"acp": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"use_project_internal"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_reference_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
							Set: acpReferenceHash,
						},
						"user_group_reference_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
							Set: acpReferenceHash,
						},
						"role_reference": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"role"}, false),
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"context_filter_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scope_filter_expression_list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"left_hand_side": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"CATEGORY", "PROJECT"}, false),
												},
												"operator": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"IN", "IN_ALL", "NOT_IN"}, false),
												},
												"right_hand_side": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"collection": {
																Type:         schema.TypeString,
																Optional:     true,
																Computed:     true,
																ValidateFunc: validation.StringInSlice([]string{"ALL"}, false),
															},
															"categories": {
																Type:     schema.TypeList,
																MaxItems: 1,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"value": {
																			Type:     schema.TypeSet,
																			Optional: true,
																			Computed: true,
																			Elem:     &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
															"uuid_list": {
																Type:     schema.TypeSet,
																Optional: true,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
											},
										},
									},
									"entity_filter_expression_list": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"left_hand_side_entity_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													ValidateFunc: utils.StringLowerCaseValidateFunc,
												},
												"operator": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"IN", "NOT_IN"}, false),
												},
												"right_hand_side": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"collection": {
																Type:         schema.TypeString,
																Optional:     true,
																Computed:     true,
																ValidateFunc: validation.StringInSlice([]string{"ALL", "SELF_OWNED"}, false),
															},
															"categories": {
																Type:     schema.TypeList,
																MaxItems: 1,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"value": {
																			Type:     schema.TypeSet,
																			Optional: true,
																			Computed: true,
																			Elem:     &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
															"uuid_list": {
																Type:     schema.TypeSet,
																Optional: true,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
											},
										},
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
					},
				},
			},
			"user_list": {
				Type:         schema.TypeList,
				Optional:     true,
				RequiredWith: []string{"use_project_internal"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"directory_service_user": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_principal_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"directory_service_reference": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "directory_service",
												},
												"uuid": {
													Type:     schema.TypeString,
													Required: true,
												},
												"name": {
													Type:     schema.TypeString,
													Optional: true,
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
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"identity_provider_reference": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "identity_provider",
												},
												"uuid": {
													Type:     schema.TypeString,
													Required: true,
												},
												"name": {
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
						"metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"user_group_list": {
				Type:         schema.TypeList,
				Optional:     true,
				RequiredWith: []string{"use_project_internal"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"directory_service_user_group": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"distinguished_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"saml_user_group": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"idp_uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"directory_service_ou": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"distinguished_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"cluster_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNutanixProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	var uuid, taskUUID string
	// if use project internal flag is set ,  we will use projects_internal API
	//nolint:staticcheck
	if _, ok := d.GetOkExists("use_project_internal"); ok {
		req := &v3.ProjectInternalIntentInput{
			Spec:       expandProjectInternalSpec(d, meta),
			Metadata:   expandMetadata(d, "project"),
			APIVersion: d.Get("api_version").(string),
		}

		resp, err := conn.V3.CreateProjectInternal(ctx, req)
		if err != nil {
			return diag.FromErr(err)
		}

		uuid = *resp.Metadata.UUID
		taskUUID = resp.Status.ExecutionContext.TaskUUID.(string)

		// Wait for the Project to be available
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"QUEUED", "RUNNING"},
			Target:     []string{"SUCCEEDED"},
			Refresh:    taskStateRefreshFunc(conn, taskUUID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      vmDelay,
			MinTimeout: vmMinTimeout,
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for project(%s) to create: %s", uuid, errWaitTask)
		}

		d.SetId(uuid)

		// once project is created , create acp .
		// check if acp is given in resource
		if _, ok1 := d.GetOk("acp"); ok1 {
			request := &v3.ProjectInternalIntentInput{}
			spec := &v3.ProjectInternalSpec{}
			metadata := &v3.Metadata{}
			projDetails := &v3.ProjectDetails{}
			accessControlPolicy := make([]*v3.AccessControlPolicyList, 0)

			var clusterUUID string
			if clusterID, ok := d.GetOk("cluster_uuid"); ok {
				clusterUUID = clusterID.(string)
			}
			response, err := conn.V3.GetProjectInternal(ctx, d.Id())
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					d.SetId("")
					return nil
				}
				return diag.Errorf("error reading Project %s: %s", d.Id(), err)
			}
			if response.Metadata != nil {
				metadata = response.Metadata
			}

			if response.Spec != nil {
				spec = response.Spec

				if response.Spec.ProjectDetail != nil || response.Spec.AccessControlPolicyList != nil {
					projDetails = response.Spec.ProjectDetail
					accessControlPolicy = response.Spec.AccessControlPolicyList
				}

				if len(response.Spec.ProjectDetail.Resources.ResourceDomain.Resources) > 0 {
					projDetails.Resources.ResourceDomain = response.Spec.ProjectDetail.Resources.ResourceDomain
				} else {
					projDetails.Resources.ResourceDomain = nil
				}
			}

			if acp, ok := d.GetOk("acp"); ok {
				acp := acp.([]interface{})
				accessControlPolicy = expandCreateAcp(acp, d, d.Id(), clusterUUID, meta)
			}
			spec.AccessControlPolicyList = accessControlPolicy
			spec.ProjectDetail = projDetails

			request.Spec = spec
			request.Metadata = metadata
			request.APIVersion = response.APIVersion

			UpdateResp, err := conn.V3.UpdateProjectInternal(ctx, d.Id(), request)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					d.SetId("")
				}
				return diag.FromErr(err)
			}

			uuid = *UpdateResp.Metadata.UUID
			taskUUID = UpdateResp.Status.ExecutionContext.TaskUUID.(string)
			// Wait for the Project to be available
			UpstateConf := &resource.StateChangeConf{
				Pending:    []string{"QUEUED", "RUNNING"},
				Target:     []string{"SUCCEEDED"},
				Refresh:    taskStateRefreshFunc(conn, taskUUID),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      vmDelay,
				MinTimeout: vmMinTimeout,
			}

			if _, errWaitTask := UpstateConf.WaitForStateContext(ctx); errWaitTask != nil {
				return diag.Errorf("error waiting for project(%s) to update: %s", uuid, errWaitTask)
			}
		}
	} else {
		req := &v3.Project{
			Spec:       expandProjectSpec(d),
			Metadata:   expandMetadata(d, "project"),
			APIVersion: d.Get("api_version").(string),
		}

		resp, err := conn.V3.CreateProject(req)
		if err != nil {
			return diag.FromErr(err)
		}

		uuid = *resp.Metadata.UUID
		taskUUID = resp.Status.ExecutionContext.TaskUUID.(string)

		// Wait for the Project to be available
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"QUEUED", "RUNNING"},
			Target:     []string{"SUCCEEDED"},
			Refresh:    taskStateRefreshFunc(conn, taskUUID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      vmDelay,
			MinTimeout: vmMinTimeout,
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for project(%s) to create: %s", uuid, errWaitTask)
		}

		d.SetId(uuid)
	}
	return resourceNutanixProjectRead(ctx, d, meta)
}

func resourceNutanixProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API
	diags := diag.Diagnostics{}

	//nolint:staticcheck
	if _, ok := d.GetOkExists("use_project_internal"); ok {
		project, err := conn.V3.GetProjectInternal(ctx, d.Id())
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}

		m, c := setRSEntityMetadata(project.Metadata)

		if err := d.Set("name", project.Status.ProjectStatus.Name); err != nil {
			return diag.Errorf("error setting `name` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("description", project.Status.ProjectStatus.Description); err != nil {
			return diag.Errorf("error setting `description` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("state", project.Status.State); err != nil {
			return diag.Errorf("error setting `state` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("is_default", project.Status.ProjectStatus.Resources.IsDefault); err != nil {
			return diag.Errorf("error setting `is_default` for Project(%s): %s", d.Id(), err)
		}
		// Deprecated since v2.4.0. Prism Central no longer supports `resource_domain` for projects; remove this block from your configuration/scripts.
		if err := d.Set("account_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.AccountReferenceList)); err != nil {
			return diag.Errorf("error setting `account_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("environment_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.EnvironmentReferenceList)); err != nil {
			return diag.Errorf("error setting `environment_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("default_subnet_reference", []interface{}{flattenReference(project.Spec.ProjectDetail.Resources.DefaultSubnetReference)}); err != nil {
			return diag.Errorf("error setting `default_subnet_reference` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("user_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.UserReferenceList)); err != nil {
			return diag.Errorf("error setting `user_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("external_user_group_reference_list",
			flattenReferenceList(project.Spec.ProjectDetail.Resources.ExternalUserGroupReferenceList)); err != nil {
			return diag.Errorf("error setting `external_user_group_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("subnet_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.SubnetReferenceList)); err != nil {
			return diag.Errorf("error setting `subnet_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("external_network_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.ExternalNetworkList)); err != nil {
			return diag.Errorf("error setting `external_network_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("metadata", m); err != nil {
			return diag.Errorf("error setting `metadata` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("project_reference", flattenReferenceValues(project.Metadata.ProjectReference)); err != nil {
			return diag.Errorf("error setting `project_reference` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("owner_reference", flattenReferenceValues(project.Metadata.OwnerReference)); err != nil {
			return diag.Errorf("error setting `owner_reference` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("categories", c); err != nil {
			return diag.Errorf("error setting `categories` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("api_version", project.APIVersion); err != nil {
			return diag.Errorf("error setting `api_version` for Project(%s): %s", d.Id(), err)
		}

		if err := d.Set("tunnel_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.TunnelReferenceList)); err != nil {
			return diag.Errorf("error setting `tunnel_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("vpc_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.VPCReferenceList)); err != nil {
			return diag.Errorf("error setting `vpc_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("cluster_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.ClusterReferenceList)); err != nil {
			return diag.Errorf("error setting `cluster_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("default_environment_reference", flattenReferenceValuesList(project.Spec.ProjectDetail.Resources.DefaultEnvironmentReference)); err != nil {
			return diag.Errorf("error setting `default_environment_reference` for Project(%s): %s", d.Id(), err)
		}
		remoteACPFlat := flattenProjectAcp(project.Status.AccessControlPolicyListStatus)
		if err := d.Set("acp", orderACPsLikeState(d, remoteACPFlat)); err != nil {
			return diag.Errorf("error setting `acp` for Project(%s): %s", d.Id(), err)
		}

		// During plan, Terraform runs a refresh (Read). If the user explicitly set `acp`
		// in config and removed an ACP from the middle, subsequent ACPs shift indices and
		// the plan can look noisy. We detect this by comparing the ACP role order from the
		// API with the ACP role order in raw config, then emit a warning so it shows up in
		// plan output.
		if newConfigRoles, set := rawConfigACPRoleUUIDs(d); set {
			oldRoles := make([]string, 0, len(remoteACPFlat))
			for _, acp := range remoteACPFlat {
				if role := getACPRoleUUID(acp); role != "" {
					oldRoles = append(oldRoles, role)
				}
			}
			if acpRemovalFromMiddleRoles(oldRoles, newConfigRoles) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  acpIndexWarningSummary,
					Detail:   acpIndexWarningDetail,
				})
			}
		}
	} else {
		project, err := conn.V3.GetProject(d.Id())
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}

		m, c := setRSEntityMetadata(project.Metadata)

		if err := d.Set("name", project.Status.Name); err != nil {
			return diag.Errorf("error setting `name` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("description", project.Status.Descripion); err != nil {
			return diag.Errorf("error setting `description` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("state", project.Status.State); err != nil {
			return diag.Errorf("error setting `state` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("is_default", project.Status.Resources.IsDefault); err != nil {
			return diag.Errorf("error setting `is_default` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("account_reference_list", flattenReferenceList(project.Spec.Resources.AccountReferenceList)); err != nil {
			return diag.Errorf("error setting `account_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("environment_reference_list", flattenReferenceList(project.Spec.Resources.EnvironmentReferenceList)); err != nil {
			return diag.Errorf("error setting `environment_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("default_subnet_reference", []interface{}{flattenReference(project.Spec.Resources.DefaultSubnetReference)}); err != nil {
			return diag.Errorf("error setting `default_subnet_reference` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("user_reference_list", flattenReferenceList(project.Spec.Resources.UserReferenceList)); err != nil {
			return diag.Errorf("error setting `user_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("external_user_group_reference_list",
			flattenReferenceList(project.Spec.Resources.ExternalUserGroupReferenceList)); err != nil {
			return diag.Errorf("error setting `external_user_group_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("subnet_reference_list", flattenReferenceList(project.Spec.Resources.SubnetReferenceList)); err != nil {
			return diag.Errorf("error setting `subnet_reference_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("external_network_list", flattenReferenceList(project.Spec.Resources.ExternalNetworkList)); err != nil {
			return diag.Errorf("error setting `external_network_list` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("metadata", m); err != nil {
			return diag.Errorf("error setting `metadata` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("project_reference", flattenReferenceValues(project.Metadata.ProjectReference)); err != nil {
			return diag.Errorf("error setting `project_reference` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("owner_reference", flattenReferenceValues(project.Metadata.OwnerReference)); err != nil {
			return diag.Errorf("error setting `owner_reference` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("categories", c); err != nil {
			return diag.Errorf("error setting `categories` for Project(%s): %s", d.Id(), err)
		}
		if err := d.Set("api_version", project.APIVersion); err != nil {
			return diag.Errorf("error setting `api_version` for Project(%s): %s", d.Id(), err)
		}
	}
	return diags
}

func resourceNutanixProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API
	diags := diag.Diagnostics{}

	var uuid, taskUUID string
	//nolint:staticcheck
	if _, ok := d.GetOkExists("use_project_internal"); ok {
		request := &v3.ProjectInternalIntentInput{}
		spec := &v3.ProjectInternalSpec{}
		metadata := &v3.Metadata{}
		projDetails := &v3.ProjectDetails{}
		var accessControlPolicy []*v3.AccessControlPolicyList

		response, err := conn.V3.GetProjectInternal(ctx, d.Id())
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				d.SetId("")
				return nil
			}
			return diag.Errorf("error reading Project %s: %s", d.Id(), err)
		}
		if response.Metadata != nil {
			metadata = response.Metadata
		}

		if response.Spec != nil {
			spec = response.Spec

			if response.Spec.ProjectDetail != nil || response.Spec.AccessControlPolicyList != nil {
				projDetails = response.Spec.ProjectDetail
			}
		}
		var clusterUUID string
		if clusterID, ok := d.GetOk("cluster_uuid"); ok {
			clusterUUID = clusterID.(string)
		}

		if d.HasChange("name") {
			projDetails.Name = utils.StringPtr(d.Get("name").(string))
		}
		if d.HasChange("description") {
			projDetails.Description = utils.StringPtr(d.Get("description").(string))
		}
		if d.HasChange("account_reference_list") {
			projDetails.Resources.AccountReferenceList = expandReferenceList(d, "account_reference_list")
		}
		if d.HasChange("environment_reference_list") {
			projDetails.Resources.EnvironmentReferenceList = expandReferenceList(d, "environment_reference_list")
		}
		if d.HasChange("default_subnet_reference") {
			projDetails.Resources.DefaultSubnetReference = expandReferenceList(d, "default_subnet_reference")[0]
		}
		if d.HasChange("user_reference_list") {
			projDetails.Resources.UserReferenceList = expandReferenceSet(d, "user_reference_list")
		}
		if d.HasChange("external_user_group_reference_list") {
			projDetails.Resources.ExternalUserGroupReferenceList = expandReferenceSet(d, "external_user_group_reference_list")
		}
		if d.HasChange("subnet_reference_list") {
			projDetails.Resources.SubnetReferenceList = expandReferenceSet(d, "subnet_reference_list")
		}
		if d.HasChange("external_network_list") {
			projDetails.Resources.ExternalNetworkList = expandReferenceList(d, "external_network_list")
		}
		if d.HasChange("tunnel_reference_list") {
			projDetails.Resources.TunnelReferenceList = expandReferenceList(d, "tunnel_reference_list")
		}
		if d.HasChange("vpc_reference_list") {
			projDetails.Resources.VPCReferenceList = expandReferenceList(d, "vpc_reference_list")
		}
		if d.HasChange("cluster_reference_list") {
			projDetails.Resources.ClusterReferenceList = expandReferenceList(d, "cluster_reference_list")
		}
		if d.HasChange("default_environment_reference") {
			projDetails.Resources.DefaultEnvironmentReference = expandOptionalReference(d, "default_environment_reference", "environment")
		}

		if d.HasChange("metadata") || d.HasChange("project_reference") ||
			d.HasChange("owner_reference") || d.HasChange("categories") {
			if err = getMetadataAttributes(d, response.Metadata, "project"); err != nil {
				return diag.Errorf("error expanding metadata: %+v", err)
			}
		}
		if d.HasChange("api_version") {
			response.APIVersion = d.Get("api_version").(string)
		}

		if d.HasChange("acp") {
			// Emit a targeted warning when ACPs are removed from the middle of the list,
			// which causes positional index shifting and noisy diffs.
			if oldRaw, newRaw := d.GetChange("acp"); oldRaw != nil && newRaw != nil {
				if oldList, ok1 := oldRaw.([]interface{}); ok1 {
					if newList, ok2 := newRaw.([]interface{}); ok2 {
						if acpRemovalFromMiddle(oldList, newList) {
							diags = append(diags, diag.Diagnostic{
								Severity: diag.Warning,
								Summary:  acpIndexWarningSummary,
								Detail:   acpIndexWarningDetail,
							})
						}
					}
				}
			}

			acp := d.Get("acp").([]interface{})
			log.Printf("[DEBUG] acp has changed")
			aJSON, _ := json.MarshalIndent(acp, "", "  ")
			log.Printf("[DEBUG] acp: %s", string(aJSON))
			accessControlPolicy = UpdateExpandAcpRM(acp, response, d, meta, d.Id(), clusterUUID)
		} else {
			accessControlPolicy = UpdateACPNoChange(response)
		}

		spec.AccessControlPolicyList = accessControlPolicy
		spec.ProjectDetail = projDetails

		request.Spec = spec
		request.Metadata = metadata
		request.APIVersion = response.APIVersion

		resp, err := conn.V3.UpdateProjectInternal(ctx, d.Id(), request)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				d.SetId("")
			}
			return diag.FromErr(err)
		}

		uuid = *resp.Metadata.UUID
		taskUUID = resp.Status.ExecutionContext.TaskUUID.(string)
	} else {
		project, err := conn.V3.GetProject(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		project.Status = nil

		if d.HasChange("name") {
			project.Spec.Name = d.Get("name").(string)
		}
		if d.HasChange("description") {
			project.Spec.Descripion = d.Get("description").(string)
		}
		if d.HasChange("account_reference_list") {
			project.Spec.Resources.AccountReferenceList = expandReferenceList(d, "account_reference_list")
		}
		if d.HasChange("environment_reference_list") {
			project.Spec.Resources.EnvironmentReferenceList = expandReferenceList(d, "environment_reference_list")
		}
		if d.HasChange("default_subnet_reference") {
			project.Spec.Resources.DefaultSubnetReference = expandReferenceList(d, "default_subnet_reference")[0]
		}
		if d.HasChange("user_reference_list") {
			project.Spec.Resources.UserReferenceList = expandReferenceSet(d, "user_reference_list")
		}
		if d.HasChange("external_user_group_reference_list") {
			project.Spec.Resources.ExternalUserGroupReferenceList = expandReferenceSet(d, "external_user_group_reference_list")
		}
		if d.HasChange("subnet_reference_list") {
			project.Spec.Resources.SubnetReferenceList = expandReferenceSet(d, "subnet_reference_list")
		}
		if d.HasChange("external_network_list") {
			project.Spec.Resources.ExternalNetworkList = expandReferenceList(d, "external_network_list")
		}
		if d.HasChange("metadata") || d.HasChange("project_reference") ||
			d.HasChange("owner_reference") || d.HasChange("categories") {
			if err = getMetadataAttributes(d, project.Metadata, "project"); err != nil {
				return diag.Errorf("error expanding metadata: %+v", err)
			}
		}
		if d.HasChange("api_version") {
			project.APIVersion = d.Get("api_version").(string)
		}

		resp, err := conn.V3.UpdateProject(d.Id(), project)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				d.SetId("")
			}
			return diag.FromErr(err)
		}

		uuid = *resp.Metadata.UUID
		taskUUID = resp.Status.ExecutionContext.TaskUUID.(string)
	}

	// Wait for the Project to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for project(%s) to update: %s", uuid, errWaitTask)
	}

	return append(diags, resourceNutanixProjectRead(ctx, d, meta)...)
}

func resourceNutanixProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	resp, err := conn.V3.DeleteProject(d.Id())
	if err != nil {
		return diag.Errorf("error deleting project id %s): %s", d.Id(), err)
	}

	// Wait for the Project to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING", "DELETED_PENDING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, cast.ToString(resp.Status.ExecutionContext.TaskUUID)),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for project (%s) to update: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func expandProjectSpec(d *schema.ResourceData) *v3.ProjectSpec {
	return &v3.ProjectSpec{
		Name:       d.Get("name").(string),
		Descripion: d.Get("description").(string),
		Resources: &v3.ProjectResources{
			AccountReferenceList:           expandReferenceList(d, "account_reference_list"),
			EnvironmentReferenceList:       expandReferenceList(d, "environment_reference_list"),
			DefaultSubnetReference:         expandReferenceList(d, "default_subnet_reference")[0],
			UserReferenceList:              expandReferenceSet(d, "user_reference_list"),
			ExternalUserGroupReferenceList: expandReferenceSet(d, "external_user_group_reference_list"),
			SubnetReferenceList:            expandReferenceSet(d, "subnet_reference_list"),
			ExternalNetworkList:            expandReferenceList(d, "external_network_list"),
		},
	}
}

func expandResourceDomain(d *schema.ResourceData) *v3.ResourceDomain {
	resourceDomain, ok := d.GetOk("resource_domain")
	if !ok {
		return nil
	}
	resources := cast.ToStringMap(resourceDomain.([]interface{})[0])["resources"].([]interface{})

	rs := make([]*v3.Resources, len(resources))
	for i, resource := range resources {
		r := cast.ToStringMap(resource)
		rs[i] = &v3.Resources{
			Limit:        utils.Int64Ptr(cast.ToInt64(r["limit"])),
			ResourceType: cast.ToString(r["resource_type"]),
		}
	}
	return &v3.ResourceDomain{Resources: rs}
}

func flattenResourceDomain(resourceDomain *v3.ResourceDomain) (res []map[string]interface{}) {
	if resourceDomain != nil {
		if len(resourceDomain.Resources) > 0 {
			resources := make([]map[string]interface{}, len(resourceDomain.Resources))

			for i, r := range resourceDomain.Resources {
				resources[i] = map[string]interface{}{
					"units":         r.Units,
					"value":         cast.ToInt64(r.Value),
					"limit":         cast.ToInt64(r.Limit),
					"resource_type": r.ResourceType,
				}
			}
			res = append(res, map[string]interface{}{"resources": resources})
		}
	}
	return
}

func expandReferenceByMap(reference map[string]interface{}) *v3.ReferenceValues {
	return &v3.ReferenceValues{
		Kind: cast.ToString(reference["kind"]),
		Name: cast.ToString(reference["name"]),
		UUID: cast.ToString(reference["uuid"]),
	}
}

func expandReferenceList(d *schema.ResourceData, key string) []*v3.ReferenceValues {
	references := d.Get(key).([]interface{})
	list := make([]*v3.ReferenceValues, len(references))

	for i, r := range references {
		list[i] = expandReferenceByMap(cast.ToStringMap(r))
	}
	return list
}

func expandReferenceSet(d *schema.ResourceData, key string) []*v3.ReferenceValues {
	references := d.Get(key).(*schema.Set).List()
	list := make([]*v3.ReferenceValues, len(references))

	for i, r := range references {
		list[i] = expandReferenceByMap(cast.ToStringMap(r))
	}
	return list
}

func expandMetadata(d *schema.ResourceData, kind string) *v3.Metadata {
	metadata := new(v3.Metadata)

	if err := getMetadataAttributes(d, metadata, kind); err != nil {
		log.Printf("Error expanding metadata: %+v", err)
	}
	return metadata
}

func expandProjectDetails(d *schema.ResourceData) *v3.ProjectDetails {
	return &v3.ProjectDetails{
		Name:        utils.StringPtr(d.Get("name").(string)),
		Description: utils.StringPtr(d.Get("description").(string)),
		Resources: &v3.ProjectInternalResources{
			AccountReferenceList:           expandReferenceList(d, "account_reference_list"),
			EnvironmentReferenceList:       expandReferenceList(d, "environment_reference_list"),
			DefaultSubnetReference:         expandReferenceList(d, "default_subnet_reference")[0],
			UserReferenceList:              expandReferenceSet(d, "user_reference_list"),
			ExternalUserGroupReferenceList: expandReferenceSet(d, "external_user_group_reference_list"),
			SubnetReferenceList:            expandReferenceSet(d, "subnet_reference_list"),
			ExternalNetworkList:            expandReferenceList(d, "external_network_list"),
			TunnelReferenceList:            expandReferenceList(d, "tunnel_reference_list"),
			ClusterReferenceList:           expandReferenceList(d, "cluster_reference_list"),
			VPCReferenceList:               expandReferenceList(d, "vpc_reference_list"),
			DefaultEnvironmentReference:    expandOptionalReference(d, "default_environment_reference", "environment"),
		},
	}
}

func expandProjectInternalSpec(d *schema.ResourceData, meta interface{}) *v3.ProjectInternalSpec {
	proSpec := &v3.ProjectInternalSpec{}
	accessControlPolicyList := []*v3.AccessControlPolicyList{}
	userList := make([]*v3.UserList, 0)
	userGroupList := make([]*v3.UserGroupList, 0)

	projDetail := expandProjectDetails(d)

	if user, ok := d.GetOk("user_list"); ok {
		userList = expandUser(user.([]interface{}), d, projDetail.Resources.UserReferenceList, meta)
	}

	if usergroup, ok := d.GetOk("user_group_list"); ok {
		userGroupList = expandUserGroup(usergroup.([]interface{}), d, projDetail.Resources.ExternalUserGroupReferenceList, meta)
	}

	proSpec.ProjectDetail = projDetail
	proSpec.AccessControlPolicyList = accessControlPolicyList
	proSpec.UserList = userList
	proSpec.UserGroupList = userGroupList

	return proSpec
}

func expandAcp(pr []interface{}, d *schema.ResourceData) []*v3.AccessControlPolicyList {
	if len(pr) > 0 {
		acpList := make([]*v3.AccessControlPolicyList, len(pr))

		for k, val := range pr {
			acps := &v3.AccessControlPolicyList{}
			acpSpec := &v3.AccessControlPolicySpec{}
			acpRes := &v3.AccessControlPolicyResources{}

			v := val.(map[string]interface{})

			if v1, ok1 := v["operation"]; ok1 {
				acps.Operation = utils.StringPtr(v1.(string))
			}
			if v1, ok1 := v["name"]; ok1 {
				acpSpec.Name = utils.StringPtr(v1.(string))
			}
			if v1, ok1 := v["description"]; ok1 {
				acpSpec.Description = utils.StringPtr(v1.(string))
			}
			if v, ok := v["user_reference_list"]; ok {
				acpRes.UserReferenceList = validateArrayRef(v.(*schema.Set), utils.StringPtr("user"))
			}

			if v, ok := v["user_group_reference_list"]; ok {
				acpRes.UserGroupReferenceList = validateArrayRef(v.(*schema.Set), utils.StringPtr("user_group"))
			}

			if v, ok := v["role_reference"]; ok {
				acpRes.RoleReference = validateRefList(v.([]interface{}), nil)
			}

			if cfl, ok := v["context_filter_list"]; ok && len(cfl.([]interface{})) > 0 {
				filterList := &v3.FilterList{}
				filterList.ContextList = expandProjectContextFilterList(cfl)
				if filterList.ContextList != nil {
					acpRes.FilterList = filterList
				}
			}

			metadata := &v3.Metadata{}
			if err := getMetadataAttributes(d, metadata, "access_control_policy"); err != nil {
				return nil
			}

			acps.Metadata = metadata
			acpSpec.Resources = acpRes
			acps.ACP = acpSpec
			acpList[k] = acps
		}

		return acpList
	}
	return nil
}

func expandUser(pr []interface{}, d *schema.ResourceData, userListRef []*v3.ReferenceValues, meta interface{}) []*v3.UserList {
	if len(pr) > 0 {
		userList := make([]*v3.UserList, len(pr))

		for k, val := range pr {
			user := &v3.UserList{}
			userSpec := &v3.UserSpec{}
			userRes := &v3.UserResources{}
			metaData := &v3.Metadata{}

			v := val.(map[string]interface{})

			if v2, ok1 := v["directory_service_user"]; ok1 {
				userRes.DirectoryServiceUser = expandDirectoryServiceUserPI(v2.([]interface{}))
			}
			if _, ok1 := v["identity_provider_user"]; ok1 {
				userRes.IdentityProviderUser = expandIdentityProviderUser(d)
			}

			userSpec.Resources = userRes
			user.User = userSpec
			metaData.Kind = utils.StringPtr("user")
			metaData.UUID = &userListRef[k].UUID
			user.Operation = utils.StringPtr("ADD")
			user.Metadata = metaData
			userList[k] = user
		}
		return userList
	}
	return nil
}

func expandUserGroup(pr []interface{}, d *schema.ResourceData, userListRef []*v3.ReferenceValues, meta interface{}) []*v3.UserGroupList {
	if len(pr) > 0 {
		userList := make([]*v3.UserGroupList, len(pr))

		for k, val := range pr {
			user := &v3.UserGroupList{}
			userSpec := &v3.UserGroupSpec{}
			userRes := &v3.UserGroupResources{}
			metaData := &v3.Metadata{}

			v := val.(map[string]interface{})

			if ds, ok := v["directory_service_user_group"]; ok {
				userRes.DirectoryServiceUserGroup = expandDirectoryUserGroup(ds.([]interface{}))
			}

			if ds, ok := v["directory_service_ou"]; ok {
				userRes.DirectoryServiceOU = expandDirectoryUserGroup(ds.([]interface{}))
			}

			if su, ok := v["saml_user_group"]; ok {
				userRes.SamlUserGroup = expandSamlUserGroup(su.([]interface{}))
			}

			userSpec.Resources = userRes
			user.UserGroup = userSpec
			metaData.Kind = utils.StringPtr("user_group")
			metaData.UUID = &userListRef[k].UUID
			user.Operation = utils.StringPtr("ADD")
			user.Metadata = metaData
			userList[k] = user
		}
		return userList
	}
	return nil
}

func expandDirectoryServiceUserPI(pr []interface{}) *v3.DirectoryServiceUser {
	if pr != nil {
		res := &v3.DirectoryServiceUser{}
		entry := pr[0].(map[string]interface{})
		if v1, ok1 := entry["directory_service_reference"]; ok1 {
			res.DirectoryServiceReference = expandReference(v1.([]interface{})[0].(map[string]interface{}))
		}

		if v1, ok1 := entry["user_principal_name"]; ok1 {
			res.UserPrincipalName = utils.StringPtr(v1.(string))
		}
		return res
	}
	return nil
}

func flattenProjectAcp(acp []*v3.ProjectAccessControlPolicyListStatus) []map[string]interface{} {
	if len(acp) > 0 {
		extSub := make([]map[string]interface{}, len(acp))

		for k, v := range acp {
			exts := make(map[string]interface{})

			if v.ProjectAccessControlPolicyStatus != nil {
				if v.ProjectAccessControlPolicyStatus.Name != nil {
					exts["name"] = v.ProjectAccessControlPolicyStatus.Name
				}

				if v.ProjectAccessControlPolicyStatus.Description != nil {
					exts["description"] = v.ProjectAccessControlPolicyStatus.Description
				}

				if v.ProjectAccessControlPolicyStatus.Resources != nil {
					exts["user_reference_list"] = flattenArrayReferenceValues(v.ProjectAccessControlPolicyStatus.Resources.UserReferenceList)
					exts["user_group_reference_list"] = flattenArrayReferenceValues(v.ProjectAccessControlPolicyStatus.Resources.UserGroupReferenceList)
					exts["role_reference"] = flattenReferenceValuesList(v.ProjectAccessControlPolicyStatus.Resources.RoleReference)
					exts["context_filter_list"] = flattenContextList(v.ProjectAccessControlPolicyStatus.Resources.FilterList.ContextList)
				}
			}

			if v.Metadata != nil {
				m, _ := setRSEntityMetadata(v.Metadata)
				exts["metadata"] = m
			}
			extSub[k] = exts
		}
		return extSub
	}
	return nil
}

func UpdateExpandAcp(pr []interface{}, res *v3.ProjectInternalIntentResponse, kind string) []*v3.AccessControlPolicyList {
	if len(pr) > 0 {
		acpList := make([]*v3.AccessControlPolicyList, len(pr))

		for k, val := range pr {
			acps := &v3.AccessControlPolicyList{}
			acpSpec := &v3.AccessControlPolicySpec{}
			acpRes := &v3.AccessControlPolicyResources{}
			var filterList v3.FilterList

			v := val.(map[string]interface{})

			if v1, ok1 := v["operation"]; ok1 {
				if kind == "old" {
					acps.Operation = utils.StringPtr("UPDATE")
				} else {
					acps.Operation = utils.StringPtr(v1.(string))
				}
			}
			if v1, ok1 := v["name"]; ok1 {
				acpSpec.Name = utils.StringPtr(v1.(string))
			}
			if v1, ok1 := v["description"]; ok1 {
				acpSpec.Description = utils.StringPtr(v1.(string))
			}
			if v, ok := v["user_reference_list"]; ok {
				acpRes.UserReferenceList = validateArrayRef(v.(*schema.Set), utils.StringPtr("user"))
			}

			if v, ok := v["user_group_reference_list"]; ok {
				acpRes.UserGroupReferenceList = validateArrayRef(v.(*schema.Set), utils.StringPtr("user_group"))
			}

			if v, ok := v["role_reference"]; ok {
				acpRes.RoleReference = validateRefList(v.([]interface{}), nil)
			}

			if cfl, ok := v["context_filter_list"]; ok && len(cfl.([]interface{})) > 0 {
				filterList.ContextList = expandProjectContextFilterList(cfl)

				if filterList.ContextList != nil {
					acpRes.FilterList = &filterList
				}
			}

			metadata := &v3.Metadata{}
			metadata.Kind = res.Spec.AccessControlPolicyList[k].Metadata.Kind
			metadata.UUID = res.Spec.AccessControlPolicyList[k].Metadata.UUID
			metadata.Categories = nil
			metadata.ProjectReference = nil

			acps.Metadata = metadata
			acpSpec.Resources = acpRes
			acps.ACP = acpSpec
			acpList[k] = acps
		}

		return acpList
	}
	return nil
}

func expandProjectContextFilterList(pr interface{}) []*v3.ContextList {
	if pr != nil {
		contextList := make([]*v3.ContextList, 0)
		for _, a1 := range pr.([]interface{}) {
			var context v3.ContextList
			con := a1.(map[string]interface{})

			context.ScopeFilterExpressionList = expandScopeExpressionList(con)
			context.EntityFilterExpressionList = expandEntityExpressionList(con)

			contextList = append(contextList, &context)
		}
		return contextList
	}
	return nil
}

func expandProjectInternalResourceDomain(pr interface{}) *v3.ResourceDomain {
	if pr != nil {
		resources := cast.ToStringMap(pr.([]interface{})[0])["resources"].([]interface{})

		rs := make([]*v3.Resources, len(resources))
		for i, resource := range resources {
			r := cast.ToStringMap(resource)
			rs[i] = &v3.Resources{
				Limit:        utils.Int64Ptr(cast.ToInt64(r["limit"])),
				ResourceType: cast.ToString(r["resource_type"]),
			}
		}
		return &v3.ResourceDomain{Resources: rs}
	}
	return nil
}

func expandOptionalReference(d *schema.ResourceData, key string, kind string) *v3.Reference {
	if v, ok := d.GetOk(key); ok {
		val := v.([]interface{})[0].(map[string]interface{})
		return validateRef(val)
	}
	return nil
}

func UpdateACPNoChange(resp *v3.ProjectInternalIntentResponse) []*v3.AccessControlPolicyList {
	acpList := resp.Spec.AccessControlPolicyList
	acpRes := make([]*v3.AccessControlPolicyList, len(acpList))
	for k, v := range acpList {
		acps := &v3.AccessControlPolicyList{}

		acps.ACP = v.ACP
		acps.Metadata = v.Metadata
		acps.Metadata = v.Metadata
		acps.Operation = utils.StringPtr("UPDATE")

		acpRes[k] = acps
	}
	return acpRes
}

func expandCreateAcp(pr []interface{}, d *schema.ResourceData, projectUUID string, clusterUUID string, meta interface{}) []*v3.AccessControlPolicyList {
	if len(pr) > 0 {
		acpList := make([]*v3.AccessControlPolicyList, len(pr))

		for k, val := range pr {
			acps := &v3.AccessControlPolicyList{}
			acpSpec := &v3.AccessControlPolicySpec{}
			acpRes := &v3.AccessControlPolicyResources{}

			v := val.(map[string]interface{})

			if v1, ok1 := v["name"]; ok1 {
				acpSpec.Name = utils.StringPtr(v1.(string))
			}
			if v1, ok1 := v["description"]; ok1 {
				acpSpec.Description = utils.StringPtr(v1.(string))
			}
			// Handle multiple users in user_reference_list (TypeSet)
			if v1, ok := v["user_reference_list"]; ok {
				if userSet, ok := v1.(*schema.Set); ok && userSet.Len() > 0 {
					userRefs := userSet.List()
					refList := make([]*v3.Reference, 0, len(userRefs))
					for _, userRef := range userRefs {
						ref := expandReference(userRef.(map[string]interface{}))
						if ref != nil {
							ref.Kind = utils.StringPtr("user")
							refList = append(refList, ref)
						}
					}
					acpRes.UserReferenceList = refList
				}
			}

			// Handle multiple user groups in user_group_reference_list (TypeSet)
			if v1, ok := v["user_group_reference_list"]; ok {
				if groupSet, ok := v1.(*schema.Set); ok && groupSet.Len() > 0 {
					groupRefs := groupSet.List()
					refList := make([]*v3.Reference, 0, len(groupRefs))
					for _, groupRef := range groupRefs {
						ref := expandReference(groupRef.(map[string]interface{}))
						if ref != nil {
							ref.Kind = utils.StringPtr("user_group")
							refList = append(refList, ref)
						}
					}
					acpRes.UserGroupReferenceList = refList
				}
			}

			if v, ok := v["role_reference"]; ok {
				acpRes.RoleReference = validateRefList(v.([]interface{}), nil)

				// role uuid
				roleID := acpRes.RoleReference.UUID

				// check for project collaboration. default is set to true
				pcCollab := true
				//nolint:staticcheck
				if pc, ok1 := d.GetOkExists("enable_collab"); ok1 {
					pcCollab = pc.(bool)
				}
				// get permissions based on role

				conList := getRolesPermission(*roleID, meta, projectUUID, clusterUUID, pcCollab)

				// generate the filter list based on role
				filterList := &v3.FilterList{}
				filterList.ContextList = conList

				if filterList.ContextList != nil {
					acpRes.FilterList = filterList
				}
			}

			metadata := &v3.Metadata{}
			metadata.Kind = utils.StringPtr("access_control_policy")

			acps.Operation = utils.StringPtr("ADD")
			acps.Metadata = metadata
			acpSpec.Resources = acpRes
			acps.ACP = acpSpec
			acpList[k] = acps
		}

		return acpList
	}
	return nil
}

func UpdateExpandAcpRM(pr []interface{}, res *v3.ProjectInternalIntentResponse, d *schema.ResourceData, meta interface{}, projectUUID string, clusterUUID string) []*v3.AccessControlPolicyList {
	if len(pr) > 0 {
		// When "acp" is schema.TypeList, Terraform compares positionally. However, the API may
		// return ACPs in a different order, so we must update ACPs by stable identity (role UUID),
		// not by list index.
		existingByRole := make(map[string]*v3.AccessControlPolicyList)
		for _, existing := range res.Spec.AccessControlPolicyList {
			if existing == nil || existing.ACP == nil || existing.ACP.Resources == nil || existing.ACP.Resources.RoleReference == nil {
				continue
			}
			if existing.ACP.Resources.RoleReference.UUID == nil {
				continue
			}
			existingByRole[*existing.ACP.Resources.RoleReference.UUID] = existing
		}

		acpList := make([]*v3.AccessControlPolicyList, len(pr))
		for k, val := range pr {
			acps := &v3.AccessControlPolicyList{}
			acpSpec := &v3.AccessControlPolicySpec{}
			acpRes := &v3.AccessControlPolicyResources{}

			v := val.(map[string]interface{})

			if v1, ok1 := v["name"]; ok1 {
				acpSpec.Name = utils.StringPtr(v1.(string))
			}
			if v1, ok1 := v["description"]; ok1 {
				acpSpec.Description = utils.StringPtr(v1.(string))
			}
			// Handle multiple users in user_reference_list (TypeSet)
			if v, ok := v["user_reference_list"]; ok {
				if userSet, ok := v.(*schema.Set); ok && userSet.Len() > 0 {
					userRefs := userSet.List()
					refList := make([]*v3.Reference, 0, len(userRefs))
					for _, userRef := range userRefs {
						ref := expandReference(userRef.(map[string]interface{}))
						if ref != nil {
							ref.Kind = utils.StringPtr("user")
							refList = append(refList, ref)
						}
					}
					acpRes.UserReferenceList = refList
				}
			}

			// Handle multiple user groups in user_group_reference_list (TypeSet)
			if v, ok := v["user_group_reference_list"]; ok {
				if groupSet, ok := v.(*schema.Set); ok && groupSet.Len() > 0 {
					groupRefs := groupSet.List()
					refList := make([]*v3.Reference, 0, len(groupRefs))
					for _, groupRef := range groupRefs {
						ref := expandReference(groupRef.(map[string]interface{}))
						if ref != nil {
							ref.Kind = utils.StringPtr("user_group")
							refList = append(refList, ref)
						}
					}
					acpRes.UserGroupReferenceList = refList
				}
			}

			if v, ok := v["role_reference"]; ok {
				acpRes.RoleReference = validateRefList(v.([]interface{}), nil)

				// get permissions based on role
				roleID := acpRes.RoleReference.UUID

				// check for project collaboration. default is set to true
				pcCollab := true
				//nolint:staticcheck
				if pc, ok1 := d.GetOkExists("enable_collab"); ok1 {
					pcCollab = pc.(bool)
				}

				conList := getRolesPermission(*roleID, meta, projectUUID, clusterUUID, pcCollab)

				// get the filter list based on role
				filterList := &v3.FilterList{}
				filterList.ContextList = conList

				if filterList.ContextList != nil {
					acpRes.FilterList = filterList
				}
			}

			metadata := &v3.Metadata{}
			// Match existing ACP by role UUID; if found, UPDATE with correct metadata UUID.
			if acpRes.RoleReference != nil && acpRes.RoleReference.UUID != nil {
				if existing, ok := existingByRole[*acpRes.RoleReference.UUID]; ok && existing.Metadata != nil {
					acps.Operation = utils.StringPtr("UPDATE")
					metadata.Kind = existing.Metadata.Kind
					metadata.UUID = existing.Metadata.UUID
					metadata.Categories = nil
					metadata.ProjectReference = nil
				} else {
					acps.Operation = utils.StringPtr("ADD")
					metadata.Kind = utils.StringPtr("access_control_policy")
				}
			} else {
				acps.Operation = utils.StringPtr("ADD")
				metadata.Kind = utils.StringPtr("access_control_policy")
			}

			acps.Metadata = metadata
			acpSpec.Resources = acpRes
			acps.ACP = acpSpec
			acpList[k] = acps
		}

		// check for delete ACP

		ck, delacp := checkACPdelete(res, acpList)
		if ck {
			acpList = append(acpList, delacp)
		}
		return acpList
	}
	return nil
}

func checkACPdelete(resp *v3.ProjectInternalIntentResponse, acpList []*v3.AccessControlPolicyList) (bool, *v3.AccessControlPolicyList) {
	oldACP := len(resp.Status.AccessControlPolicyListStatus)
	newACP := len(acpList)

	if newACP < oldACP {
		for _, v := range resp.Spec.AccessControlPolicyList {
			oldmeta := v.Metadata.UUID
			checkk := false
			for _, val := range acpList {
				if oldmeta == val.Metadata.UUID {
					checkk = true
				}
			}
			if !checkk {
				v.Operation = utils.StringPtr("DELETE")
				return true, v
			}
		}
		return false, nil
	}
	return false, nil
}
