// Package clustersv2 provides resources for managing Nutanix clusters.
package clustersv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/clusters"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

//nolint:misspell // British English spelling is intentional
const (
	CANCELED = "CANCELLED"
)

// customizeClusterProfileExtIDDiff forces Terraform to detect changes when cluster_profile_ext_id
// is explicitly set to empty or removed from config. This is needed because computed+optional
// attributes may not trigger HasChange when going from a value to empty.
func customizeClusterProfileExtIDDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// Get old and new values using GetChange (old from state, new from config)
	oldValRaw, newValRaw := d.GetChange("cluster_profile_ext_id")
	oldValStr := ""
	newValStr := ""
	if oldValRaw != nil {
		oldValStr = oldValRaw.(string)
	}
	if newValRaw != nil {
		newValStr = newValRaw.(string)
	}

	// If old and new are the same, no change needed
	if oldValStr == newValStr {
		log.Printf("[DEBUG] customizeClusterProfileExtIDDiff: cluster_profile_ext_id is the same")
		return nil
	}

	// If there's no old value, nothing to disassociate
	if oldValStr == "" {
		log.Printf("[DEBUG] customizeClusterProfileExtIDDiff: cluster_profile_ext_id is empty")
		return nil
	}

	// Check if the attribute was explicitly set to empty or removed from config
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		// Config is null/unknown, can't determine - but if newValStr is empty and oldValStr is not,
		// it means we're going from a value to empty, so force the diff
		if newValStr == "" && oldValStr != "" {
			log.Printf("[DEBUG] CustomizeDiff: cluster_profile_ext_id changed to empty (config unknown), forcing diff. Old: '%s'", oldValStr)
			if err := d.SetNew("cluster_profile_ext_id", ""); err != nil {
				log.Printf("[DEBUG] customizeClusterProfileExtIDDiff: error setting new cluster_profile_ext_id: %v", err)
				return err
			}
			return nil
		}
		return nil
	}

	configMap := rawConfig.AsValueMap()
	val, exists := configMap["cluster_profile_ext_id"]

	// Case 1: Attribute was removed from config (exists in state but not in config)
	if !exists && oldValStr != "" {
		log.Printf("[DEBUG] CustomizeDiff: cluster_profile_ext_id removed from config, forcing diff. Old: '%s'", oldValStr)
		// Force a diff by setting new value to empty
		if err := d.SetNew("cluster_profile_ext_id", ""); err != nil {
			log.Printf("[DEBUG] customizeClusterProfileExtIDDiff: error setting new cluster_profile_ext_id: %v", err)
			return err
		}
		log.Printf("[DEBUG] customizeClusterProfileExtIDDiff: cluster_profile_ext_id set to empty")
		return nil
	}

	// Case 2: Attribute exists in config - check if it's explicitly set to empty
	if exists && !val.IsNull() && val.IsKnown() {
		if val.Type().Equals(cty.String) {
			strVal := val.AsString()
			if strVal == "" && oldValStr != "" {
				log.Printf("[DEBUG] CustomizeDiff: cluster_profile_ext_id explicitly set to empty, forcing diff. Old: '%s'", oldValStr)
				// Force a diff by setting new value to empty
				if err := d.SetNew("cluster_profile_ext_id", ""); err != nil {
					log.Printf("[DEBUG] customizeClusterProfileExtIDDiff: error setting new cluster_profile_ext_id: %v", err)
					return err
				}
				log.Printf("[DEBUG] customizeClusterProfileExtIDDiff: cluster_profile_ext_id set to empty")
				return nil
			}
		}
	}

	// Case 3: HasChange detected a change, but we need to ensure empty string is explicitly set
	// when going from non-empty to empty (this handles computed values that might be empty)
	if d.HasChange("cluster_profile_ext_id") && newValStr == "" && oldValStr != "" {
		log.Printf("[DEBUG] CustomizeDiff: cluster_profile_ext_id changed from non-empty to empty, forcing explicit empty string. Old: '%s'", oldValStr)
		if err := d.SetNew("cluster_profile_ext_id", ""); err != nil {
			log.Printf("[DEBUG] customizeClusterProfileExtIDDiff: error setting new cluster_profile_ext_id: %v", err)
			return err
		}
		log.Printf("[DEBUG] customizeClusterProfileExtIDDiff: cluster_profile_ext_id set to empty")
		return nil
	}

	log.Printf("[DEBUG] customizeClusterProfileExtIDDiff: cluster_profile_ext_id no change needed")
	return nil
}

func ResourceNutanixClusterV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixClusterV2Create,
		ReadContext:   ResourceNutanixClusterV2Read,
		UpdateContext: ResourceNutanixClusterV2Update,
		DeleteContext: ResourceNutanixClusterV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: clusterImportFunc,
		},
		CustomizeDiff: customizeClusterProfileExtIDDiff,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dryrun": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number_of_nodes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Set:      hashNodeItem,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"controller_vm_ip": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     common.SchemaForIPList(false),
									},
									"node_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_ip": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     common.SchemaForIPList(false),
									},

									// expand cluster with node params
									"should_skip_host_networking": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"should_skip_add_node": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"should_skip_pre_expand_checks": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"should_skip_discovery": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"should_skip_imaging": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"should_validate_rack_awareness": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"is_nos_compatible": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"is_compute_only": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"is_never_scheduleable": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"is_light_compute": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"hypervisor_hostname": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						// remove node params
						"remove_node_params": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"extra_params": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"should_skip_upgrade_check": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"skip_space_check": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"should_skip_add_check": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
											},
										},
									},

									"should_skip_remove": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"should_skip_prechecks": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
					},
				},
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_address": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     common.SchemaForIPList(false),
						},
						"external_data_services_ip": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     common.SchemaForIPList(false),
						},
						"external_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nfs_subnet_white_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name_server_ip_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     common.SchemaForIPList(true),
						},
						"ntp_server_ip_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     common.SchemaForIPList(true),
						},
						"smtp_server": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email_address": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"server": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_address": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem:     common.SchemaForIPList(true),
												},
												"port": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"username": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"password": {
													Type:      schema.TypeString,
													Sensitive: true,
													Optional:  true,
												},
											},
										},
									},
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice(SMTPTypeStrings, false),
									},
								},
							},
						},
						"masquerading_ip": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     common.SchemaForIPList(false),
						},
						"masquerading_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"management_server": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     common.SchemaForIPList(false),
									},
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice(ManagementServerTypeStrings, false),
									},
									"is_drs_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"is_registered": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"is_in_use": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"fqdn": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"key_management_server_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(KeyManagementServerTypeStrings, false),
						},
						"backplane": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_segmentation_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"vlan_tag": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"subnet":  common.SchemaForValuePrefixLengthResource(ipv4PrefixLengthDefaultValue),
									"netmask": common.SchemaForValuePrefixLengthResource(ipv4PrefixLengthDefaultValue),
								},
							},
						},
						"http_proxy_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     common.SchemaForIPList(false),
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"password": {
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"proxy_types": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice(HTTPProxyTypeStrings, false),
										},
									},
								},
							},
						},
						"http_proxy_white_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target": {
										Type:     schema.TypeString,
										Required: true,
									},
									"target_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(HTTPProxyWhiteListTargetStrings, false),
									},
								},
							},
						},
					},
				},
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"incarnation_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"build_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"build_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"full_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"commit_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"short_commit_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"hypervisor_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cluster_function": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice(ClusterFunctionStrings, false),
							},
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"authorized_public_key_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},

							DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
								log.Printf("[DEBUG] DiffSuppressFunc authorized_public_key_list k: %v, oldValue: %v, newValue: %v", k, oldValue, newValue)

								// Check if the list has changed
								if d.HasChange("config.0.authorized_public_key_list") {
									log.Printf("[DEBUG] authorized_public_key_list has changed \n")
									oldRaw, newRaw := d.GetChange("config.0.authorized_public_key_list")
									// Convert to lists of interfaces
									oldList := oldRaw.([]interface{})
									newList := newRaw.([]interface{})
									// Sort lists based on a unique field (e.g., "key") for comparison
									sort.SliceStable(oldList, func(i, j int) bool {
										return oldList[i].(map[string]interface{})["key"].(string) < oldList[j].(map[string]interface{})["key"].(string)
									})
									sort.SliceStable(newList, func(i, j int) bool {
										return newList[i].(map[string]interface{})["key"].(string) < newList[j].(map[string]interface{})["key"].(string)
									})
									aJSON, _ := json.Marshal(oldList)
									log.Printf("[DEBUG] authorized_public_key_list oldList: %v", string(aJSON))
									aJSON, _ = json.Marshal(newList)
									log.Printf("[DEBUG] authorized_public_key_list newList: %v", string(aJSON))
									// Check if lists are equal
									if reflect.DeepEqual(oldList, newList) {
										log.Printf("[DEBUG] authorized_public_key_list are  equal \n")
										return true
									}

									log.Printf("[DEBUG] authorized_public_key_list are not equal \n")
									return false
								}
								log.Printf("[DEBUG] authorized_public_key_list has not changed \n")
								return false
							},
						},
						"redundancy_factor": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"cluster_software_map": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"software_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cluster_arch": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(ClusterArchStrings, false),
						},
						"fault_tolerance_state": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"current_max_fault_tolerance": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"desired_max_fault_tolerance": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"domain_awareness_level": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice(DomainAwarenessLevelStrings, false),
									},
									"current_cluster_fault_tolerance": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"desired_cluster_fault_tolerance": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice(ClusterFaultToleranceStrings, false),
									},
									"redundancy_status": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_cassandra_preparation_done": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"is_zookeeper_preparation_done": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"is_remote_support_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"operation_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(OperationModeStrings, false),
						},
						"is_lts": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_password_remote_login_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"encryption_in_transit_status": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"encryption_option": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"encryption_scope": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"pulse_status": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"pii_scrubbing_level": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice(PIIScrubbingLevelStrings, false),
									},
								},
							},
						},
						"is_available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"upgrade_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      common.HashStringItem,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vm_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"inefficient_vm_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_profile_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress diff if both old and new are empty (treat empty string and nil as equivalent)
					oldTrimmed := strings.TrimSpace(old)
					newTrimmed := strings.TrimSpace(new)
					log.Printf("[DEBUG] cluster_profile_ext_id DiffSuppressFunc: oldTrimmed: %s, newTrimmed: %s", oldTrimmed, newTrimmed)
					if oldTrimmed == "" && newTrimmed == "" {
						return true
					}
					// Otherwise, let Terraform handle the diff normally
					return false
				},
			},
			"backup_eligibility_score": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// Computed fields
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
		},
	}
}

func ResourceNutanixClusterV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	body := config.NewCluster()
	var dryRun *bool

	if dryRunVar, ok := d.GetOk("dryrun"); ok {
		dryRun = utils.BoolPtr(dryRunVar.(bool))
	} else {
		dryRun = utils.BoolPtr(false)
	}
	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if nodes, ok := d.GetOk("nodes"); ok {
		body.Nodes = expandNodeReference(nodes)
	}
	if network, ok := d.GetOk("network"); ok {
		body.Network = expandClusterNetworkReference(network)
	}

	if configVar, ok := d.GetOk("config"); ok {
		body.Config = expandClusterConfigReference(configVar, d)
	}
	if upgradeStatus, ok := d.GetOk("upgrade_status"); ok {
		body.UpgradeStatus = common.ExpandEnum(upgradeStatus, UpgradeStatusMap, "upgrade_status")
	}

	if containerName, ok := d.GetOk("container_name"); ok {
		body.ContainerName = utils.StringPtr(containerName.(string))
	}

	if categories, ok := d.GetOk("categories"); ok {
		// Convert to slice - handles both TypeList and TypeSet
		categoriesList := common.InterfaceToSlice(categories)
		categoriesListStr := common.ExpandListOfString(categoriesList)
		body.Categories = categoriesListStr
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create Cluster Request Body: %s", string(aJSON))

	resp, err := conn.ClusterEntityAPI.CreateCluster(body, dryRun)
	if err != nil {
		return diag.Errorf("error while creating clusters : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching cluster UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)
	aJSON, _ = json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Create Cluster Task Response Details: %s", string(aJSON))

	randomID := utils.GenUUID()

	d.SetId(randomID)

	return nil
}

func ResourceNutanixClusterV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	if d.Get("ext_id").(string) == "" {
		log.Printf("[DEBUG] ResourceNutanixClusterV2Read : extID is empty")
		err := getClusterExtID(d, conn)
		if err != nil {
			// Check if the error is a ClusterNotFoundError
			if _, ok := err.(*ClusterNotFoundError); ok {
				log.Printf("[DEBUG] ResourceNutanixClusterV2Read : Cluster not found, err -> %v", err)
				diags := diag.Diagnostics{
					{
						Severity: diag.Warning,
						Summary:  "Cluster not found. Please register the cluster to Prism Central if not. If deleted, then reset the state.",
						Detail:   fmt.Sprintf("Cluster %s not found: %v", d.Get("name").(string), err),
					},
				}
				return diags
			}
			log.Printf("[DEBUG] ResourceNutanixClusterV2Read : error while fetching cluster : %v", err)
			return diag.Errorf("error while fetching cluster : %v", err)
		}
	}
	// Clusters READ context
	return clusterRead(d, meta)
}

func ResourceNutanixClusterV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	var expand *string
	var nodeChanges bool

	if expandVar, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandVar.(string))
	} else {
		expand = nil
	}

	if d.Get("ext_id").(string) == "" {
		log.Printf("[DEBUG] ResourceNutanixClusterV2Update : Cluster not found, extID is empty")
		err := getClusterExtID(d, conn)

		// Check if the error is a ClusterNotFoundError
		if err != nil {
			log.Printf("[DEBUG] ResourceNutanixClusterV2Update : error while fetching cluster : %v", err)
			return diag.Errorf("error while fetching cluster : %v: Please register the cluster to Prism Central if not.", err)
		}
		log.Printf("[DEBUG] ResourceNutanixClusterV2Update : error while fetching cluster : %v", err)
	}

	log.Printf("[DEBUG] ResourceNutanixClusterV2Update : Cluster found, extID : %s", d.Id())

	// === Handle Node Add/Remove ===
	if d.HasChange("nodes") {
		if diags := handleNodeChanges(ctx, d, meta, conn, expand); diags.HasError() {
			return diags
		}
		nodeChanges = true
	}

	// === Handle Category Association Changes ===
	if d.HasChange("categories") {
		if diags := handleCategoryAssociationChanges(ctx, d, meta, conn, expand); diags.HasError() {
			return diags
		}
	}

	// === Handle Cluster Profile association update ===
	// Check if cluster_profile_ext_id changed or was explicitly set to empty
	// For computed+optional attributes, we need to check both HasChange and if it was explicitly set to empty
	hasChange := d.HasChange("cluster_profile_ext_id")
	isExplicitlySetToEmpty := isClusterProfileExtIDExplicitlySetToEmpty(d)

	if hasChange || isExplicitlySetToEmpty {
		if diags := handleClusterProfileAssociationUpdate(ctx, d, meta, conn, expand); diags.HasError() {
			return diags
		}
	}

	// === Handle other Cluster field changes ===
	updateSpec, hasClusterFieldChange := handleClusterFieldUpdate(d)
	if !hasClusterFieldChange {
		log.Printf("[DEBUG] No cluster field changes detected, skipping UpdateClusterById")
		if nodeChanges {
			// delay to allow cluster to stabilize after node changes
			log.Printf("[DEBUG] Delaying for 1 minute to allow cluster to stabilize after node changes")
			time.Sleep(1 * time.Minute)
		}
		return ResourceNutanixClusterV2Read(ctx, d, meta)
	}

	// === Apply update via UpdateClusterById ===
	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] cluster update: update payload: %s", string(aJSON))

	resp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(d.Id()), expand)
	if err != nil {
		return diag.Errorf("error fetching cluster: %v", err)
	}

	args := getEtagHeader(resp, conn)
	updateResp, err := conn.ClusterEntityAPI.UpdateClusterById(utils.StringPtr(d.Id()), &updateSpec, args)
	if err != nil {
		return diag.Errorf("error updating cluster: %v", err)
	}

	// === Wait for Task completion ===
	TaskRef := updateResp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId
	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for cluster update task (%s): %s", utils.StringValue(taskUUID), errWait)
	}

	log.Printf("[DEBUG] Cluster update completed successfully")
	time.Sleep(1 * time.Minute)
	return ResourceNutanixClusterV2Read(ctx, d, meta)
}

func ResourceNutanixClusterV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	var expand *string

	if expandVar, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandVar.(string))
	} else {
		expand = nil
	}

	if d.Get("ext_id").(string) == "" {
		err := getClusterExtID(d, conn)
		if err != nil {
			// Check if the error is a ClusterNotFoundError
			if _, ok := err.(*ClusterNotFoundError); ok {
				log.Printf("[DEBUG] ResourceNutanixClusterV2Delete : Cluster not found, err -> %v", err)
				diags := diag.Diagnostics{
					{
						Severity: diag.Warning,
						Summary:  "Cluster not found. Please register the cluster to Prism Central if not. If deleted, then reset the state.",
						Detail:   fmt.Sprintf("Cluster %s not found: %v", d.Get("name").(string), err),
					},
				}
				return diags
			}
			log.Printf("[DEBUG] ResourceNutanixClusterV2Delete : error while fetching cluster : %v", err)
			return diag.Errorf("error while fetching cluster : %v: Please register the cluster to Prism Central if not.", err)
		}
	}

	readResp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(d.Id()), expand)
	if err != nil {
		return diag.Errorf("error while reading cluster : %v", err)
	}
	// Extract E-Tag Header
	args := getEtagHeader(readResp, conn)

	var dryRun *bool

	if dryRunVar, ok := d.GetOk("dryrun"); ok {
		dryRun = utils.BoolPtr(dryRunVar.(bool))
	} else {
		dryRun = utils.BoolPtr(false)
	}

	resp, err := conn.ClusterEntityAPI.DeleteClusterById(utils.StringPtr(d.Id()), dryRun, args)
	if err != nil {
		return diag.Errorf("error while deleting cluster : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func getClusterExtID(d *schema.ResourceData, conn *clusters.Client) error {
	var filter string
	if d.HasChange("name") {
		// if name changed, get the old name to fetch the cluster, since the name will be updated after update request
		oldName, _ := d.GetChange("name")
		filter = fmt.Sprintf(`name eq '%s'`, oldName.(string))
	} else {
		filter = fmt.Sprintf(`name eq '%s'`, d.Get("name").(string))
	}

	log.Printf("[DEBUG] getClusterExtID filter : %s", filter)

	// get Cluster Ext Id
	listResp, err := conn.ClusterEntityAPI.ListClusters(nil, nil, utils.StringPtr(filter), nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("error while fetching cluster : %v", err)
	}

	if utils.IntValue(listResp.Metadata.TotalAvailableResults) == 0 {
		log.Printf("[DEBUG] getClusterExtID Cluster not found, TotalAvailableResults is 0")
		return &ClusterNotFoundError{Name: d.Get("name").(string), Err: err}
	}

	if listResp.Data == nil {
		log.Printf("[DEBUG] getClusterExtID Cluster not found, clustersResponse.Data is nil")
		return &ClusterNotFoundError{Name: d.Get("name").(string), Err: err}
	}

	cls := listResp.Data.GetValue().([]config.Cluster)

	if len(cls) == 0 {
		log.Printf("[DEBUG] getClusterExtID Cluster not found, len(clusters) is 0")
		return &ClusterNotFoundError{Name: d.Get("name").(string), Err: err}
	}
	extID := utils.StringValue(cls[0].ExtId)
	if extID == "" {
		log.Printf("[DEBUG] getClusterExtID Cluster not found, extID is empty")
		return &ClusterNotFoundError{Name: d.Get("name").(string), Err: err}
	}
	log.Printf("[DEBUG] getClusterExtID Cluster found, extId : %s", extID)
	d.SetId(extID)
	err = d.Set("ext_id", extID)
	if err != nil {
		return err
	}
	return nil
}

func expandNodeReference(pr interface{}) *config.NodeReference {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		nodeRef := config.NewNodeReference()

		if nodeList, ok := val["node_list"]; ok {
			nodeRef.NodeList = expandNodeListItemReference(common.InterfaceToSlice(nodeList))
		}

		return nodeRef
	}

	return nil
}

func expandNodeListItemReference(pr []interface{}) []config.NodeListItemReference {
	if len(pr) > 0 {
		nodeList := make([]config.NodeListItemReference, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			node := config.NewNodeListItemReference()

			if controllerVMIP, ok := val["controller_vm_ip"]; ok {
				log.Printf("[DEBUG] controller_vm_ip")
				node.ControllerVmIp = expandIPAddress(controllerVMIP)
			}
			if hostIP, ok := val["host_ip"]; ok {
				log.Printf("[DEBUG] host_ip")
				node.HostIp = expandIPAddress(hostIP)
			}

			nodeList[k] = *node
		}
		return nodeList
	}
	return nil
}

func expandClusterNetworkReference(pr interface{}) *config.ClusterNetworkReference {
	if pr != nil {
		cls := config.NewClusterNetworkReference()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if externalAddress, ok := val["external_address"]; ok {
			log.Printf("[DEBUG] external_address")
			cls.ExternalAddress = expandIPAddress(externalAddress)
		}
		if externalDataServiceIP, ok := val["external_data_services_ip"]; ok {
			log.Printf("[DEBUG] external_data_services_ip")
			cls.ExternalDataServiceIp = expandIPAddress(externalDataServiceIP)
		}
		if nfsSubnetWhite, ok := val["nfs_subnet_white_list"]; ok {
			nfsSubnetWhitelist := nfsSubnetWhite.([]interface{})
			nfsSubnetWhitelistStr := make([]string, len(nfsSubnetWhitelist))
			for i, v := range nfsSubnetWhitelist {
				nfsSubnetWhitelistStr[i] = v.(string)
			}
			log.Printf("[DEBUG] nfs_subnet_white_list: %v", nfsSubnetWhitelistStr)
			cls.NfsSubnetWhitelist = nfsSubnetWhitelistStr
		}
		if nameServerIPList, ok := val["name_server_ip_list"]; ok {
			cls.NameServerIpList = expandIPAddressOrFQDN(nameServerIPList.([]interface{}))
		}
		if ntpServerIPList, ok := val["ntp_server_ip_list"]; ok {
			log.Printf("[DEBUG] ntp_server_ip_list ")
			cls.NtpServerIpList = expandIPAddressOrFQDN(ntpServerIPList.([]interface{}))
		}
		if smtpServer, ok := val["smtp_server"]; ok {
			cls.SmtpServer = expandSMTPServerRef(smtpServer)
		}
		if masqueradingIP, ok := val["masquerading_ip"]; ok {
			log.Printf("[DEBUG] masquerading_ip ")
			cls.MasqueradingIp = expandIPAddress(masqueradingIP)
		}
		if managementServer, ok := val["management_server"]; ok {
			cls.ManagementServer = expandManagementServerRef(managementServer)
		}
		if fqdn, ok := val["fqdn"]; ok && fqdn != "" {
			log.Printf("[DEBUG] network/fqdn : %v", fqdn)
			cls.Fqdn = utils.StringPtr(fqdn.(string))
		}
		if keyManagementServerType, ok := val["key_management_server_type"]; ok {
			cls.KeyManagementServerType = common.ExpandEnum(keyManagementServerType, KeyManagementServerTypeMap, "key_management_server_type")
		}
		if backplane, ok := val["backplane"]; ok {
			cls.Backplane = expandBackplaneNetworkParams(backplane)
		}

		if httpProxyList, ok := val["http_proxy_list"]; ok {
			cls.HttpProxyList = expandHTTPProxyList(httpProxyList.([]interface{}))
		}
		if httpProxyWhiteList, ok := val["http_proxy_white_list"]; ok {
			cls.HttpProxyWhiteList = expandHTTPProxyWhiteList(httpProxyWhiteList.([]interface{}))
		}

		return cls
	}
	return nil
}

func expandHTTPProxyWhiteList(proxyTypesWhiteList []interface{}) []config.HttpProxyWhiteListConfig {
	if len(proxyTypesWhiteList) > 0 {
		httpProxyWhiteList := make([]config.HttpProxyWhiteListConfig, len(proxyTypesWhiteList))

		for k, v := range proxyTypesWhiteList {
			val := v.(map[string]interface{})
			httpProxy := config.NewHttpProxyWhiteListConfig()

			if target, ok := val["target"]; ok {
				httpProxy.Target = utils.StringPtr(target.(string))
			}
			httpProxy.TargetType = common.ExpandEnum(val["target_type"], HTTPProxyWhiteListTargetMap, "target_type")

			httpProxyWhiteList[k] = *httpProxy
		}
		return httpProxyWhiteList
	}
	return nil
}

func expandHTTPProxyList(httpProxyList []interface{}) []config.HttpProxyConfig {
	if len(httpProxyList) > 0 {
		httpProxyConfig := make([]config.HttpProxyConfig, len(httpProxyList))

		for k, v := range httpProxyList {
			val := v.(map[string]interface{})
			httpProxy := config.NewHttpProxyConfig()

			if ipAddr, ok := val["ip_address"]; ok {
				httpProxy.IpAddress = expandIPAddress(ipAddr)
			}
			if port, ok := val["port"]; ok {
				httpProxy.Port = utils.IntPtr(port.(int))
			}
			if username, ok := val["username"]; ok {
				httpProxy.Username = utils.StringPtr(username.(string))
			}
			if password, ok := val["password"]; ok {
				httpProxy.Password = utils.StringPtr(password.(string))
			}
			if name, ok := val["name"]; ok {
				httpProxy.Name = utils.StringPtr(name.(string))
			}
			httpProxy.ProxyTypes = common.ExpandEnumList(val["proxy_types"], HTTPProxyTypeMap, "proxy_type")

			httpProxyConfig[k] = *httpProxy
		}
		return httpProxyConfig
	}
	return nil
}

func expandClusterConfigReference(pr interface{}, d *schema.ResourceData) *config.ClusterConfigReference {
	if pr != nil {
		clsConf := config.NewClusterConfigReference()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if buildInfo, ok := val["build_info"]; ok && d.HasChange("config.0.build_info") {
			clsConf.BuildInfo = expandBuildReference(buildInfo)
		}
		clsConf.ClusterFunction = common.ExpandEnumList(val["cluster_function"], ClusterFunctionMap, "cluster_function")

		if _, ok := val["authorized_public_key_list"]; ok && d.HasChange("config.0.authorized_public_key_list") {
			_, newObj := d.GetChange("config.0.authorized_public_key_list")
			clsConf.AuthorizedPublicKeyList = expandPublicKey(newObj.([]interface{}))
		}
		if redundancyFactor, ok := val["redundancy_factor"]; ok && d.HasChange("config.0.redundancy_factor") {
			clsConf.RedundancyFactor = utils.Int64Ptr(int64(redundancyFactor.(int)))
		}
		clsConf.ClusterArch = common.ExpandEnum(val["cluster_arch"], ClusterArchMap, "cluster_arch")

		if faultToleranceState, ok := val["fault_tolerance_state"]; ok && d.HasChange("config.0.fault_tolerance_state") {
			clsConf.FaultToleranceState = expandFaultToleranceState(faultToleranceState)
		}
		if operationMode, ok := val["operation_mode"]; ok && d.HasChange("config.0.operation_mode") {
			clsConf.OperationMode = common.ExpandEnum(operationMode, OperationModeMap, "operation_mode")
		}

		if encryptionInTransitStatus, ok := val["encryption_in_transit_status"]; ok && d.HasChange("config.0.encryption_in_transit_status") {
			clsConf.EncryptionInTransitStatus = common.ExpandEnum(encryptionInTransitStatus, EncryptionStatusMap, "encryption_in_transit_status")
		}

		if pulseStatus, ok := val["pulse_status"]; ok && d.HasChange("config.0.pulse_status") {
			clsConf.PulseStatus = expandPulseStatus(pulseStatus)
		}

		return clsConf
	}
	return nil
}

func expandPulseStatus(status interface{}) *config.PulseStatus {
	if status == nil || len(status.([]interface{})) == 0 {
		log.Printf("[DEBUG] PulseStatus is nil")
		return nil
	}

	pulse := config.NewPulseStatus()
	prI := status.([]interface{})
	val := prI[0].(map[string]interface{})

	if isEnabled, ok := val["is_enabled"]; ok {
		pulse.IsEnabled = utils.BoolPtr(isEnabled.(bool))
	}
	if piiScrubbingLevel, ok := val["pii_scrubbing_level"]; ok {
		pulse.PiiScrubbingLevel = common.ExpandEnum(piiScrubbingLevel, PIIScrubbingLevelMap, "pii_scrubbing_level")
	}

	return pulse
}

func expandIPAddress(pr interface{}) *import4.IPAddress {
	if pr != nil {
		ipAddress := import4.NewIPAddress()
		prI := pr.([]interface{})
		if len(prI) == 0 {
			return nil
		}
		val := prI[0].(map[string]interface{})

		if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
			ipAddress.Ipv4 = expandIPv4Address(ipv4)
		}
		if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
			ipAddress.Ipv6 = expandIPv6Address(ipv6)
		}
		aJSON, _ := json.Marshal(ipAddress)
		log.Printf("[DEBUG] ipAddress : %v", string(aJSON))
		return ipAddress
	}
	return nil
}

func expandIPAddressOrFQDN(pr []interface{}) []import4.IPAddressOrFQDN {
	if len(pr) > 0 {
		ips := make([]import4.IPAddressOrFQDN, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			ip := import4.NewIPAddressOrFQDN()

			if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
				ip.Ipv4 = expandIPv4Address(ipv4)
			}
			if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
				ip.Ipv6 = expandIPv6Address(ipv6)
			}
			if fqdn, ok := val["fqdn"]; ok && len(fqdn.([]interface{})) > 0 {
				ip.Fqdn = expandFQDN(fqdn.([]interface{}))
			}
			ips[k] = *ip
		}
		aJSON, _ := json.Marshal(ips)
		log.Printf("[DEBUG] ipAddressOrFQDN : %v", string(aJSON))
		return ips
	}
	return nil
}

func expandIPv4Address(pr interface{}) *import4.IPv4Address {
	// nil check
	if pr == nil {
		return nil
	}

	// safe type assert for expected slice
	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}

	// safe type assert for first element being a map
	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok || len(valMap) == 0 {
		return nil
	}
	ipv4 := import4.NewIPv4Address()

	if v, ok := valMap["value"]; ok {
		if s, ok2 := v.(string); ok2 {
			ipv4.Value = utils.StringPtr(s)
		}
	}

	if p, ok := valMap["prefix_length"]; ok {
		if n, ok2 := p.(int); ok2 {
			ipv4.PrefixLength = utils.IntPtr(n)
		}
	}

	return ipv4
}

func expandIPv6Address(pr interface{}) *import4.IPv6Address {
	// nil check
	if pr == nil {
		return nil
	}

	// safe type assert for expected slice
	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}

	// safe type assert for first element being a map
	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok || len(valMap) == 0 {
		return nil
	}

	ipv6 := import4.NewIPv6Address()

	if v, ok := valMap["value"]; ok {
		if s, ok2 := v.(string); ok2 {
			ipv6.Value = utils.StringPtr(s)
		}
	}

	if p, ok := valMap["prefix_length"]; ok {
		if n, ok2 := p.(int); ok2 {
			ipv6.PrefixLength = utils.IntPtr(n)
		}
	}

	return ipv6
}

func expandSMTPServerRef(pr interface{}) *config.SmtpServerRef {
	if len(pr.([]interface{})) == 0 {
		return nil
	}

	if pr != nil {
		smtp := config.NewSmtpServerRef()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if emailAddress, ok := val["email_address"]; ok {
			smtp.EmailAddress = utils.StringPtr(emailAddress.(string))
		}
		if server, ok := val["server"]; ok {
			smtp.Server = expandSMTPNetwork(server.([]interface{}))
		}
		if smtpType, ok := val["type"]; ok {
			smtp.Type = common.ExpandEnum(smtpType, SMTPTypeMap, "smtp_type")
		}

		return smtp
	}
	return nil
}

func expandBackplaneNetworkParams(pr interface{}) *config.BackplaneNetworkParams {
	if len(pr.([]interface{})) == 0 {
		return nil
	}

	if pr != nil && len(pr.([]interface{})) > 0 {
		backPlane := config.NewBackplaneNetworkParams()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if isSegmentationEnabled, ok := val["is_segmentation_enabled"]; ok {
			backPlane.IsSegmentationEnabled = utils.BoolPtr(isSegmentationEnabled.(bool))
		}
		if subnet, ok := val["subnet"]; ok {
			backPlane.Subnet = expandIPv4Address(subnet)
		}
		if netmask, ok := val["netmask"]; ok {
			backPlane.Netmask = expandIPv4Address(netmask)
		}
		if vlanTag, ok := val["vlan_tag"]; ok {
			backPlane.VlanTag = utils.Int64Ptr(int64(vlanTag.(int)))
		}

		return backPlane
	}
	return nil
}

func expandManagementServerRef(pr interface{}) *config.ManagementServerRef {
	if pr != nil && len(pr.([]interface{})) > 0 {
		mgm := config.NewManagementServerRef()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if ip, ok := val["ip"]; ok {
			log.Printf("[DEBUG] management server ip")
			mgm.Ip = expandIPAddress(ip.([]interface{}))
		}
		if mgmType, ok := val["type"]; ok {
			mgm.Type = common.ExpandEnum(mgmType, ManagementServerTypeMap, "management_server_type")
		}
		if drsEnabled, ok := val["is_drs_enabled"]; ok {
			mgm.IsDrsEnabled = utils.BoolPtr(drsEnabled.(bool))
		}
		if isRegistered, ok := val["is_registered"]; ok {
			mgm.IsRegistered = utils.BoolPtr(isRegistered.(bool))
		}
		if inUse, ok := val["is_in_use"]; ok {
			mgm.IsInUse = utils.BoolPtr(inUse.(bool))
		}

		return mgm
	}
	return nil
}

func expandSMTPNetwork(pr []interface{}) *config.SmtpNetwork {
	if len(pr) > 0 {
		smtp := config.NewSmtpNetwork()
		val := pr[0].(map[string]interface{})

		if ipAddress, ok := val["ip_address"]; ok {
			smtp.IpAddress = &(expandIPAddressOrFQDN(ipAddress.([]interface{})))[0]
		}
		if port, ok := val["port"]; ok {
			smtp.Port = utils.IntPtr(port.(int))
		}
		if username, ok := val["username"]; ok {
			smtp.Username = utils.StringPtr(username.(string))
		}
		if password, ok := val["password"]; ok {
			smtp.Password = utils.StringPtr(password.(string))
		}

		return smtp
	}
	return nil
}

func expandFQDN(pr interface{}) *import4.FQDN {
	// nil check
	if pr == nil {
		return nil
	}

	// safe type assert for expected slice
	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}

	// safe type assert for first element being a map
	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok || len(valMap) == 0 {
		return nil
	}

	fqdn := import4.NewFQDN()

	if v, ok := valMap["value"]; ok {
		if s, ok2 := v.(string); ok2 && s != "" {
			fqdn.Value = utils.StringPtr(s)
		}
	}

	return fqdn
}

func expandBuildReference(buildInfo interface{}) *config.BuildReference {
	if buildInfo == nil || len(buildInfo.([]interface{})) == 0 {
		log.Printf("[DEBUG] buildInfo is nil")
		return nil
	}

	buildReference := config.NewBuildReference()
	buildInfoI := buildInfo.([]interface{})
	buildInfoVal := buildInfoI[0].(map[string]interface{})

	if buildType, ok := buildInfoVal["build_type"]; ok {
		buildReference.BuildType = utils.StringPtr(buildType.(string))
	}
	if version, ok := buildInfoVal["version"]; ok {
		buildReference.Version = utils.StringPtr(version.(string))
	}
	if fullVersion, ok := buildInfoVal["full_version"]; ok {
		buildReference.FullVersion = utils.StringPtr(fullVersion.(string))
	}
	if commitID, ok := buildInfoVal["commit_id"]; ok {
		buildReference.CommitId = utils.StringPtr(commitID.(string))
	}
	if shortCommitID, ok := buildInfoVal["short_commit_id"]; ok {
		buildReference.ShortCommitId = utils.StringPtr(shortCommitID.(string))
	}

	return buildReference
}

func expandPublicKey(pr []interface{}) []config.PublicKey {
	if len(pr) > 0 {
		pubKey := make([]config.PublicKey, len(pr))
		aJSON, _ := json.Marshal(pr)
		log.Printf("[DEBUG] PublicKey : %v", string(aJSON))

		for k, v := range pr {
			val := v.(map[string]interface{})
			pub := config.NewPublicKey()

			if key, ok := val["key"]; ok {
				pub.Key = utils.StringPtr(key.(string))
			}
			if name, ok := val["name"]; ok {
				pub.Name = utils.StringPtr(name.(string))
			}
			pubKey[k] = *pub
		}
		return pubKey
	}
	return nil
}

func expandFaultToleranceState(pr interface{}) *config.FaultToleranceState {
	if pr != nil && len(pr.([]interface{})) > 0 {
		fts := config.NewFaultToleranceState()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if domainAwarenessLevel, ok := val["domain_awareness_level"]; ok {
			fts.DomainAwarenessLevel = common.ExpandEnum(domainAwarenessLevel, DomainAwarenessLevelMap, "domain_awareness_level")
		}

		if currentClusterFaultTolerance, ok := val["current_cluster_fault_tolerance"]; ok {
			fts.CurrentClusterFaultTolerance = common.ExpandEnum(
				currentClusterFaultTolerance,
				ClusterFaultToleranceMap,
				"current_cluster_fault_tolerance",
			)
		}

		if desiredClusterFaultTolerance, ok := val["desired_cluster_fault_tolerance"]; ok {
			fts.DesiredClusterFaultTolerance = common.ExpandEnum(
				desiredClusterFaultTolerance,
				ClusterFaultToleranceMap,
				"desired_cluster_fault_tolerance",
			)
		}

		return fts
	}
	return nil
}

func handleNodeChanges(ctx context.Context, d *schema.ResourceData, meta interface{}, conn *clusters.Client, expand *string) diag.Diagnostics {
	log.Printf("[DEBUG] Handling node changes for cluster: %s", d.Id())

	resp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(d.Id()), expand)
	if err != nil {
		return diag.Errorf("error fetching cluster for node diff: %v", err)
	}

	existingNodes := resp.Data.GetValue().(config.Cluster).Nodes.NodeList
	rawNodes := expandNodeReference(d.Get("nodes")).NodeList
	added, removed, changed := DiffNodes(d, existingNodes, rawNodes)

	// === Add Nodes ===
	for _, nodeWithFlags := range added {
		diags, unconfiguredNodeDetails := discoverUnconfiguredNode(ctx, d, meta, *conn, nodeWithFlags.Node)
		if diags.HasError() {
			return diags
		}
		diags, networkDetails := fetchNetworkDetailsForNodes(ctx, d, meta, *conn, *unconfiguredNodeDetails)
		if diags.HasError() {
			return diags
		}

		if diags := expandClusterWithNewNode(ctx, d, meta, *conn, *unconfiguredNodeDetails, *networkDetails,
			nodeWithFlags.Flags); diags.HasError() {
			return diags
		}
	}

	// === Remove Nodes ===
	for _, nodeToRemove := range removed {
		if diags := removeNodeFromCluster(ctx, d, meta, *conn, nodeToRemove); diags.HasError() {
			return diags
		}
	}

	// === Log Changed Nodes (no direct API call, just informational) ===
	if len(changed) > 0 {
		b, _ := json.MarshalIndent(changed, "", "  ")
		log.Printf("[DEBUG] Nodes changed (informational only): %s", string(b))
	}

	return nil
}

// isClusterProfileExtIDExplicitlySetToEmpty checks if cluster_profile_ext_id was explicitly set to empty string in the config
// or removed from the config entirely. This is needed because for computed+optional attributes, Terraform may not detect
// a change from a value to empty or when the attribute is removed.
func isClusterProfileExtIDExplicitlySetToEmpty(d *schema.ResourceData) bool {
	// Check if there's a current value in state (meaning we might be changing from non-empty to empty/removed)
	currentVal := d.Get("cluster_profile_ext_id").(string)
	if currentVal == "" {
		// No current value, so nothing to disassociate
		return false
	}

	// Get the raw config to check if the attribute exists and what its value is
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		// Config is null/unknown, can't determine if attribute was removed
		return false
	}

	configMap := rawConfig.AsValueMap()
	val, exists := configMap["cluster_profile_ext_id"]

	if !exists {
		// Attribute not in config at all (removed from config)
		log.Printf("[DEBUG] cluster_profile_ext_id removed from config, current value in state: '%s'", currentVal)
		return true
	}

	// Attribute exists in config - check if it's explicitly set to empty string
	if !val.IsNull() && val.IsKnown() {
		// Check if it's a string type and get its value
		if val.Type().Equals(cty.String) {
			strVal := val.AsString()
			if strVal == "" {
				// Explicitly set to empty string
				log.Printf("[DEBUG] cluster_profile_ext_id explicitly set to empty in config, current value in state: '%s'", currentVal)
				return true
			}
		}
	}

	return false
}

func handleClusterProfileAssociationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}, conn *clusters.Client, expand *string) diag.Diagnostics {
	log.Printf("[DEBUG] Handling cluster profile association update for cluster: %s", d.Id())

	// Get old and new cluster profile ext IDs
	oldProfileExtIDRaw, newProfileExtIDRaw := d.GetChange("cluster_profile_ext_id")
	oldProfileExtID := ""
	newProfileExtID := ""
	if oldProfileExtIDRaw != nil {
		oldProfileExtID = oldProfileExtIDRaw.(string)
	}
	if newProfileExtIDRaw != nil {
		newProfileExtID = newProfileExtIDRaw.(string)
	}

	// Check if the attribute was explicitly set to empty or removed from config
	// This is needed because GetChange might return the computed value instead of empty
	// when the attribute is removed or set to empty in the config
	if isClusterProfileExtIDExplicitlySetToEmpty(d) {
		log.Printf("[DEBUG] cluster_profile_ext_id was explicitly set to empty or removed, treating new value as empty string")
		newProfileExtID = ""
	}

	log.Printf("[DEBUG] Cluster profile association change - old: '%s', new: '%s'", oldProfileExtID, newProfileExtID)

	clusterUUID := d.Id()
	taskconn := meta.(*conns.Client).PrismAPI

	// Build cluster reference for this cluster
	clusterRef := config.ClusterReference{
		Uuid: utils.StringPtr(clusterUUID),
	}
	clustersRef := []config.ClusterReference{clusterRef}

	ClusterReferenceListSpec := &config.ClusterReferenceListSpec{
		Clusters: clustersRef,
	}

	// Disassociate from old profile if it was set and is different from new
	if oldProfileExtID != "" && oldProfileExtID != newProfileExtID {
		log.Printf("[DEBUG] Disassociating cluster %s from cluster profile %s", clusterUUID, oldProfileExtID)
		disassociateResp, disassociateErr := conn.ClusterProfilesAPI.DisassociateClusterFromClusterProfile(utils.StringPtr(oldProfileExtID), ClusterReferenceListSpec)
		if disassociateErr != nil {
			return diag.Errorf("error disassociating cluster from cluster profile: %v", disassociateErr)
		}

		TaskRef := disassociateResp.Data.GetValue().(import1.TaskReference)
		taskUUID := TaskRef.ExtId

		// Wait for the disassociation task to complete
		stateConf := &resource.StateChangeConf{
			Pending: []string{"QUEUED", "RUNNING", "PENDING"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutUpdate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for cluster profile (%s) to disassociate: %s", utils.StringValue(taskUUID), errWaitTask)
		}

		log.Printf("[DEBUG] Cluster profile disassociation task %s completed", utils.StringValue(taskUUID))
		// Get Task Details
		taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return diag.Errorf("error while fetching cluster profile disassociation task UUID : %v", err)
		}
		taskDetails := taskResp.Data.GetValue().(import2.Task)
		aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
		log.Printf("[DEBUG] Disassociate Cluster Profile Task Details: %s", string(aJSON))
	}

	// Associate with new profile if it's set and different from old
	if newProfileExtID != "" && newProfileExtID != oldProfileExtID {
		log.Printf("[DEBUG] Associating cluster %s with cluster profile %s", clusterUUID, newProfileExtID)
		associateResp, associateErr := conn.ClusterProfilesAPI.ApplyClusterProfile(utils.StringPtr(newProfileExtID), ClusterReferenceListSpec, utils.BoolPtr(false))
		if associateErr != nil {
			return diag.Errorf("error associating cluster to cluster profile: %v", associateErr)
		}

		TaskRef := associateResp.Data.GetValue().(import1.TaskReference)
		taskUUID := TaskRef.ExtId

		// Wait for the association task to complete
		stateConf := &resource.StateChangeConf{
			Pending: []string{"QUEUED", "RUNNING", "PENDING"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutUpdate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for cluster profile (%s) to associate: %s", utils.StringValue(taskUUID), errWaitTask)
		}
		log.Printf("[DEBUG] Cluster profile association task %s completed", utils.StringValue(taskUUID))

		// Get Task Details
		taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return diag.Errorf("error while fetching cluster profile association task UUID : %v", err)
		}
		taskDetails := taskResp.Data.GetValue().(import2.Task)
		aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
		log.Printf("[DEBUG] Associate Cluster Profile Task Details: %s", string(aJSON))
	}

	return nil
}

func handleClusterFieldUpdate(d *schema.ResourceData) (config.Cluster, bool) {
	var updateSpec config.Cluster
	var hasChanges bool

	if d.HasChange("name") {
		hasChanges = true
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("network") {
		hasChanges = true
		updateSpec.Network = expandClusterNetworkReference(d.Get("network"))
	}
	if d.HasChange("config") {
		hasChanges = true
		updateSpec.Config = expandClusterConfigReference(d.Get("config"), d)
	}
	if d.HasChange("upgrade_status") {
		hasChanges = true
		updateSpec.UpgradeStatus = common.ExpandEnum(d.Get("upgrade_status"), UpgradeStatusMap, "upgrade_status")
	}
	if d.HasChange("container_name") {
		hasChanges = true
		updateSpec.ContainerName = utils.StringPtr(d.Get("container_name").(string))
	}
	// Note: categories are handled separately via handleCategoryAssociationChanges
	// Categories should use Associate/Disassociate APIs, not UpdateClusterById

	log.Printf("[DEBUG] handleClusterFieldUpdate: hasChanges=%v", hasChanges)
	return updateSpec, hasChanges
}

// handleCategoryAssociationChanges handles category association and disassociation changes
// Categories should be managed via Associate/Disassociate APIs, not through UpdateClusterById
func handleCategoryAssociationChanges(ctx context.Context, d *schema.ResourceData, meta interface{}, conn *clusters.Client, expand *string) diag.Diagnostics {
	// Get cluster ext_id - use d.Id() if available, otherwise try to get from state
	clusterExtID := d.Id()
	if clusterExtID == "" {
		// Try to get from ext_id field
		if extID, ok := d.GetOk("ext_id"); ok {
			clusterExtID = extID.(string)
		}
		if clusterExtID == "" {
			return diag.Errorf("cluster ext_id is required for category operations")
		}
	}

	log.Printf("[DEBUG] Handling category association changes for cluster: %s", clusterExtID)

	// Get old and new category values
	oldCategoriesRaw, newCategoriesRaw := d.GetChange("categories")

	// Use shared function to handle category updates
	return UpdateClusterCategories(ctx, d, meta, clusterExtID, oldCategoriesRaw, newCategoriesRaw)
}

func removeNodeFromCluster(ctx context.Context, d *schema.ResourceData, meta interface{},
	conn clusters.Client, nodeToRemove NodeWithFlags) diag.Diagnostics {
	body := &config.NodeRemovalParams{}

	nodeUUIDList := make([]string, 0)

	// set node UUID
	nodeUUIDList = append(nodeUUIDList, utils.StringValue(nodeToRemove.Node.NodeUuid))

	if len(nodeUUIDList) > 0 {
		body.NodeUuids = nodeUUIDList
	} else {
		return diag.Errorf("error while removing node : Node UUID is required for remove node")
	}

	body.ShouldSkipRemove = utils.BoolPtr(nodeToRemove.Flags.ShouldSkipRemove)
	body.ShouldSkipPrechecks = utils.BoolPtr(nodeToRemove.Flags.ShouldSkipPrechecks)
	body.ExtraParams = &config.NodeRemovalExtraParam{
		ShouldSkipUpgradeCheck: utils.BoolPtr(nodeToRemove.Flags.ShouldSkipUpgradeCheck),
		ShouldSkipSpaceCheck:   utils.BoolPtr(nodeToRemove.Flags.SkipSpaceCheck),
		ShouldSkipAddCheck:     utils.BoolPtr(nodeToRemove.Flags.ShouldSkipAddCheck),
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] cluster update: remove node request body: %s", string(aJSON))
	resp, err := conn.ClusterEntityAPI.RemoveNode(utils.StringPtr(d.Id()), body)
	if err != nil {
		return diag.Errorf("error while Removing node : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the node to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		resourceUUID, _ := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		rUUID := resourceUUID.Data.GetValue().(import2.Task)
		aJSON, _ := json.MarshalIndent(rUUID, "", "  ")
		log.Printf("Error Remove Node Task Details : %s", string(aJSON))
		return diag.Errorf("error waiting for  node (%s) to Remove: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching  node UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	bJSON, _ := json.MarshalIndent(rUUID, "", "  ")
	log.Printf("cluster update: remove node task details : %s", string(bJSON))
	return nil
}

func expandClusterWithNewNode(ctx context.Context, d *schema.ResourceData, meta interface{}, conn clusters.Client,
	unconfigureNodeDetails config.UnconfigureNodeDetails,
	nodeNetworkingDetails config.NodeNetworkingDetails,
	flags NodeFlags) diag.Diagnostics {
	unConfNode := unconfigureNodeDetails.NodeList[0]
	nodeNetInfo := nodeNetworkingDetails

	upLinks := config.NewUplinks()
	activeUplinkField := config.NewUplinksField()
	standbyUplinkField := config.NewUplinksField()
	activeUplinkField.Name = nodeNetInfo.Uplinks[0].UplinkList[0].Name
	activeUplinkField.Mac = nodeNetInfo.Uplinks[0].UplinkList[0].Mac
	activeUplinkField.Value = nodeNetInfo.Uplinks[0].UplinkList[0].Name
	standbyUplinkField.Name = nodeNetInfo.Uplinks[0].UplinkList[1].Name
	standbyUplinkField.Mac = nodeNetInfo.Uplinks[0].UplinkList[1].Mac
	standbyUplinkField.Value = nodeNetInfo.Uplinks[0].UplinkList[1].Name

	upLinks.Active = append(upLinks.Active, *activeUplinkField)
	upLinks.Standby = append(upLinks.Standby, *standbyUplinkField)

	networks := make([]config.UplinkNetworkItem, 0)
	uplinksNetworkItem := config.NewUplinkNetworkItem()
	uplinksNetworkItem.Name = nodeNetInfo.NetworkInfo.Hci[0].Name
	uplinksNetworkItem.Networks = nodeNetInfo.NetworkInfo.Hci[0].Networks
	uplinksNetworkItem.Uplinks = upLinks

	networks = append(networks, *uplinksNetworkItem)

	nodeItem := config.NewNodeItem()
	nodeItem.NodeUuid = unConfNode.NodeUuid
	nodeItem.NodePosition = unConfNode.NodePosition
	nodeItem.Model = unConfNode.RackableUnitModel
	nodeItem.BlockId = unConfNode.RackableUnitSerial
	nodeItem.HypervisorType = unConfNode.HypervisorType
	nodeItem.HypervisorVersion = unConfNode.HypervisorVersion
	nodeItem.NosVersion = unConfNode.NosVersion
	nodeItem.CurrentNetworkInterface = nodeNetInfo.Uplinks[0].UplinkList[0].Name
	nodeItem.HypervisorIp = unConfNode.HypervisorIp
	nodeItem.CvmIp = unConfNode.CvmIp
	nodeItem.IpmiIp = unConfNode.IpmiIp
	nodeItem.IsRoboMixedHypervisor = unConfNode.Attributes.IsRoboMixedHypervisor
	nodeItem.Networks = networks
	nodeItem.IsLightCompute = utils.BoolPtr(flags.IsLightCompute)
	if flags.HypervisorHostname != "" {
		nodeItem.HypervisorHostname = utils.StringPtr(flags.HypervisorHostname)
	}

	nodeList := make([]config.NodeItem, 0)
	nodeList = append(nodeList, *nodeItem)

	nodeParam := config.NewNodeParam()
	if flags.ShouldSkipHostNetworking != nil {
		nodeParam.ShouldSkipHostNetworking = flags.ShouldSkipHostNetworking
	}
	nodeParam.NodeList = nodeList

	hypervisorIsoMap := config.NewHypervisorIsoMap()
	hypervisorIsoMap.Type = unConfNode.HypervisorType

	nodeParam.HypervisorIsos = append(nodeParam.HypervisorIsos, *hypervisorIsoMap)

	configParams := config.NewConfigParams()
	configParams.TargetHypervisor = utils.StringPtr(unConfNode.HypervisorType.GetName())
	configParams.ShouldSkipDiscovery = utils.BoolPtr(flags.ShouldSkipDiscovery)
	configParams.ShouldSkipImaging = utils.BoolPtr(flags.ShouldSkipImaging)
	configParams.ShouldValidateRackAwareness = utils.BoolPtr(flags.ShouldValidateRackAwareness)
	configParams.IsNosCompatible = utils.BoolPtr(flags.IsNosCompatible)
	configParams.IsComputeOnly = utils.BoolPtr(flags.IsComputeOnly)
	configParams.IsNeverScheduleable = utils.BoolPtr(flags.IsNeverScheduleable)

	body := config.NewExpandClusterParams()
	body.ShouldSkipAddNode = utils.BoolPtr(flags.ShouldSkipAddNode)
	body.ShouldSkipPreExpandChecks = utils.BoolPtr(flags.ShouldSkipPreExpandChecks)
	body.NodeParams = nodeParam
	body.ConfigParams = configParams

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Add Node Request Body: %s", string(aJSON))

	resp, err := conn.ClusterEntityAPI.ExpandCluster(utils.StringPtr(d.Id()), body)
	if err != nil {
		return diag.Errorf("error while adding node : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the  node to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for  node (%s) to add: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching  node UUID : %v", err)
	}

	aJSON, _ = json.Marshal(resourceUUID)
	log.Printf("[DEBUG] Add Node Response: %s", string(aJSON))
	return nil
}

func fetchNetworkDetailsForNodes(ctx context.Context, d *schema.ResourceData, meta interface{},
	conn clusters.Client, node config.UnconfigureNodeDetails) (diag.Diagnostics, *config.NodeNetworkingDetails) {
	readResp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(d.Id()), nil)
	if err != nil {
		return diag.Errorf("error while reading cluster : %v", err), nil
	}
	// Extract E-Tag Header
	args := getEtagHeader(readResp, &conn)

	unconfiguredNodeDetail := node.NodeList[0]

	nodeListNetworkingDetails := make([]config.NodeListNetworkingDetails, 0)
	nodeListItem := config.NodeListNetworkingDetails{
		CurrentNetworkInterface: unconfiguredNodeDetail.CurrentNetworkInterface,
		HypervisorType:          unconfiguredNodeDetail.HypervisorType,
		HypervisorVersion:       unconfiguredNodeDetail.HypervisorVersion,
		IpmiIp:                  unconfiguredNodeDetail.IpmiIp,
		NodePosition:            unconfiguredNodeDetail.NodePosition,
		NodeUuid:                unconfiguredNodeDetail.NodeUuid,
		NosVersion:              unconfiguredNodeDetail.NosVersion,
		CvmIp:                   unconfiguredNodeDetail.CvmIp,
		HypervisorIp:            unconfiguredNodeDetail.HypervisorIp,
	}

	nodeListNetworkingDetails = append(nodeListNetworkingDetails, nodeListItem)

	nodeNetworkDetailsBody := config.NodeDetails{
		NodeList:    nodeListNetworkingDetails,
		RequestType: utils.StringPtr("expand_cluster"),
	}

	aJSON, _ := json.MarshalIndent(nodeNetworkDetailsBody, "", " ")
	log.Printf("[DEBUG] Fetch Network Info for Node to be added body : %s", string(aJSON))

	networkDetailsResp, err := conn.ClusterEntityAPI.FetchNodeNetworkingDetails(utils.StringPtr(d.Id()), &nodeNetworkDetailsBody, args)
	if err != nil {
		return diag.Errorf("error while Fetching Node Networking Details : %v", err), nil
	}

	TaskRef := networkDetailsResp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the  node to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for  node (%s) to add: %s", utils.StringValue(taskUUID), errWaitTask), nil
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching task : %v", err), nil
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	bJSON, _ := json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Fetch Network Info Task Details: %s", string(bJSON))

	uuid := strings.Split(utils.StringValue(rUUID.ExtId), "=:")[1]

	const networkingDetails = 3
	taskResponseType := config.TaskResponseType(networkingDetails)
	networkDetailsTaskResp, taskErr := conn.ClusterEntityAPI.FetchTaskResponse(utils.StringPtr(uuid), &taskResponseType)
	if taskErr != nil {
		return diag.Errorf("error while fetching Task Response for Unconfigured Nodes : %v", taskErr), nil
	}

	taskResp := networkDetailsTaskResp.Data.GetValue().(config.TaskResponse)

	if *taskResp.TaskResponseType != config.TaskResponseType(networkingDetails) {
		return diag.Errorf("error while fetching Task Response for Network Detail Nodes : %v", "task response type mismatch"), nil
	}

	nodeNetworkDetails := taskResp.Response.GetValue().(config.NodeNetworkingDetails)

	aJSON, _ = json.MarshalIndent(networkDetailsTaskResp, "", " ")
	log.Printf("[DEBUG] fetching Network Details for Node to be added task details: %s", string(aJSON))
	return nil, &nodeNetworkDetails
}

func discoverUnconfiguredNode(ctx context.Context, d *schema.ResourceData, meta interface{},
	conn clusters.Client, node config.NodeListItemReference) (diag.Diagnostics, *config.UnconfigureNodeDetails) {
	ipType := getIPType(node.ControllerVmIp)

	var addressType config.AddressType
	switch ipType {
	case "IPV4":
		addressType = config.ADDRESSTYPE_IPV4
	case "IPV6":
		addressType = config.ADDRESSTYPE_IPV6
	}

	unconfiguredNodeBody := &config.NodeDiscoveryParams{
		AddressType:  &addressType,
		IpFilterList: []import4.IPAddress{*node.ControllerVmIp},
	}

	aJSON, _ := json.MarshalIndent(unconfiguredNodeBody, "", " ")
	log.Printf("[DEBUG] Discover Unconfigured Nodes body : %s", string(aJSON))

	discoverUnconfiguredNodesResp, err := conn.ClusterEntityAPI.DiscoverUnconfiguredNodes(utils.StringPtr(d.Id()), unconfiguredNodeBody)
	if err != nil {
		return diag.Errorf("error while Discover Unconfigured Nodes : %v", err), nil
	}

	TaskRef := discoverUnconfiguredNodesResp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Nodes Trap to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Unconfigured Nodes (%s) to fetch: %s", utils.StringValue(taskUUID), errWaitTask), nil
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching Unconfigured Nodes UUID : %v", err), nil
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	jsonBody, _ := json.MarshalIndent(resourceUUID, "", "  ")
	log.Printf("[DEBUG] fetching Unconfigured Nodes resourceUUID : %s", string(jsonBody))

	uuid := strings.Split(utils.StringValue(rUUID.ExtId), "=:")[1]

	const unconfiguredNodes = 2
	taskResponseType := config.TaskResponseType(unconfiguredNodes)
	unconfiguredNodesResp, taskErr := conn.ClusterEntityAPI.FetchTaskResponse(utils.StringPtr(uuid), &taskResponseType)
	if taskErr != nil {
		return diag.Errorf("error while fetching Task Response for Unconfigured Nodes : %v", taskErr), nil
	}

	taskResp := unconfiguredNodesResp.Data.GetValue().(config.TaskResponse)

	if *taskResp.TaskResponseType != config.TaskResponseType(unconfiguredNodes) {
		return diag.Errorf("error while fetching Task Response for Unconfigured Nodes : %v", "task response type mismatch"), nil
	}

	unconfiguredNodeDetails := taskResp.Response.GetValue().(config.UnconfigureNodeDetails)

	aJSON, _ = json.MarshalIndent(unconfiguredNodeDetails, "", " ")
	log.Printf("[DEBUG] cluster expand: unconfigured node details: %s", string(aJSON))

	return nil, &unconfiguredNodeDetails
}

func clusterImportFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// calling read logic here
	diags := clusterRead(d, meta)
	if diags.HasError() {
		// convert diagnostics to error
		return nil, fmt.Errorf("failed to import cluster: %v", diags)
	}
	return []*schema.ResourceData{d}, nil
}

func clusterRead(d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	var expand *string

	if expandVar, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandVar.(string))
	} else {
		expand = nil
	}
	resp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(d.Id()), expand)
	if err != nil {
		return diag.Errorf("error while fetching cluster : %v: Please register the cluster to Prism Central if not.", err)
	}

	getResp := resp.Data.GetValue().(config.Cluster)
	aJSON, _ := json.MarshalIndent(getResp, "", "  ")
	log.Printf("[DEBUG] Read Cluster Response Details: %s", string(aJSON))

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nodes", flattenNodeReference(getResp.Nodes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network", flattenClusterNetworkReference(getResp.Network)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("config", flattenClusterConfigReference(getResp.Config)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("upgrade_status", flattenUpgradeStatus(getResp.UpgradeStatus)); err != nil {
		return diag.FromErr(err)
	}
	// Convert categories from API response to Terraform schema format (schema.Set)
	// The API returns []string, which needs to be converted to schema.Set for TypeSet schema
	categoriesList := make([]interface{}, 0)
	if len(getResp.Categories) > 0 {
		for _, category := range getResp.Categories {
			if category != "" {
				categoriesList = append(categoriesList, category)
			}
		}
	}
	// Create a schema.Set from the list using the same hash function as the schema
	categoriesSet := schema.NewSet(common.HashStringItem, categoriesList)
	if err := d.Set("categories", categoriesSet); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_count", getResp.VmCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("inefficient_vm_count", getResp.InefficientVmCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("container_name", getResp.ContainerName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_profile_ext_id", getResp.ClusterProfileExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("backup_eligibility_score", getResp.BackupEligibilityScore); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
