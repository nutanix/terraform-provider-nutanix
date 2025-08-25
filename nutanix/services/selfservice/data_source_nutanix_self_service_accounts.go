package selfservice

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/selfservice"
)

func DatsourceNutanixSelfServiceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: datsourceNutanixSelfServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "account",
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sort_order": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ASCENDING",
			},
			"sort_attribute": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "name",
			},
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
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
									"availability_zone_reference": buildReferenceSchema("availability_zone"),
									"messages_list":               buildMessageListSchema(),
									"cluster_reference":           buildReferenceSchema("cluster"),
									"resources":                   buildResourcesSchema(),
								},
							},
						},
						"spec": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"creation_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_update_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"spec_version": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_matches": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildReferenceSchema(defaultKind string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"kind": {
					Type:     schema.TypeString,
					Default:  defaultKind,
					Optional: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"uuid": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
			},
		},
	}
}

func buildMessageListSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"message": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"reason": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
				"details": {
					Type:     schema.TypeList,
					Computed: true,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{},
					},
				},
			},
		},
	}
}

func buildResourcesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"sync_interval_secs": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"last_sync_time": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"sync_status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"sync_error": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"state": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"tunnel_reference": buildReferenceSchema("tunnel"),
				"parent_reference": buildReferenceSchema("account"),
				"price_items":      buildPriceItemsSchema(),
				"data":             buildProviderDataSchema(),
			},
		},
	}
}

func buildPriceItemsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"uuid": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"state_cost_list": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"state": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"cost_list": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"interval": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"name": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"value": {
											Type:     schema.TypeFloat,
											Computed: true,
										},
										"type": {
											Type:     schema.TypeString,
											Computed: true,
										},
									},
								},
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
				"details": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"association_type": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"occurrence": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"provider_type": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"message_list": buildMessageListSchema(),
				"description": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func buildProviderDataSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"aws":     buildAWSProviderSchema(),
				"nutanix": buildNutanixProviderSchema(),
				"azure":   buildAzureProviderSchema(),
				"vmware":  buildVMwareProviderSchema(),
				"gcp":     buildGCPProviderSchema(),
				"k8s":     buildk8sProviderSchema(),
			},
		},
	}
}

func buildAWSProviderSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access_key_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"regions": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func buildNutanixProviderSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"server": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"port": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"username": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"cluster_account_reference_list": {
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
							"message_list": buildMessageListSchema(),
							"resources": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"price_items": buildPriceItemsSchema(),
										"data": {
											Type:     schema.TypeList,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"cluster_name": {
														Type:     schema.TypeString,
														Computed: true,
													},
													"cluster_uuid": {
														Type:     schema.TypeString,
														Computed: true,
													},
													"pc_account_uuid": {
														Type:     schema.TypeString,
														Computed: true,
													},
												},
											},
										},
										"type": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"state": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"sync_interval_secs": {
											Type:     schema.TypeInt,
											Computed: true,
										},
									},
								},
							},
							"uuid": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func buildAzureProviderSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tenant_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"client_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"subscription_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"cloud_environment": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func buildVMwareProviderSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"server": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"port": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"datacenter": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"username": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func buildGCPProviderSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"account_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"project_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"private_key_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"client_email": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"client_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func buildk8sProviderSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"server": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"port": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"authentication": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"username": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func datsourceNutanixSelfServiceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(*conns.Client).CalmAPI

	acFilter := &selfservice.AccountsListInput{}
	if filter, ok := d.GetOk("filter"); ok {
		acFilter.Filter = filter.(string) + ";(type!=nutanix;type!=custom_provider)"
	} else {
		acFilter.Filter = "(type!=nutanix;type!=custom_provider)"
	}

	accountResp, err := conn.Service.ListAccounts(ctx, acFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(accountResp.Entities) == 0 {
		if err := d.Set("accounts", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resource.UniqueId())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of Accounts.",
		}}
	}

	if err := d.Set("api_version", accountResp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("total_matches", accountResp.Metadata["total_matches"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("accounts", flattenAccountEntities(accountResp.Entities)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenAccountEntities(entities []map[string]interface{}) []map[string]interface{} {
	entityList := make([]map[string]interface{}, 0)

	for _, entity := range entities {
		entityMap := make(map[string]interface{})

		// Handle metadata
		metadataMap := make(map[string]interface{})
		if metadata, ok := entity["metadata"].(map[string]interface{}); ok {
			if name, ok := metadata["name"].(string); ok {
				metadataMap["name"] = name
			}
			if uuid, ok := metadata["uuid"].(string); ok {
				metadataMap["uuid"] = uuid
			}
			if kind, ok := metadata["kind"].(string); ok {
				metadataMap["kind"] = kind
			}
			if creationTime, ok := metadata["creation_time"].(string); ok {
				metadataMap["creation_time"] = creationTime
			}
			if lastUpdateTime, ok := metadata["last_update_time"].(string); ok {
				metadataMap["last_update_time"] = lastUpdateTime
			}
			if specVersionRaw, ok := metadata["spec_version"]; ok {
				switch v := specVersionRaw.(type) {
				case float64:
					if v != 0 {
						metadataMap["spec_version"] = int(v)
					}
				case int:
					if v != 0 {
						metadataMap["spec_version"] = v
					}
				}
			}
			if specVersion, ok := metadata["spec_version"].(int); ok {
				metadataMap["spec_version"] = specVersion
			}
			entityMap["metadata"] = []interface{}{metadataMap}
		}

		// Handle api_version
		if apiVersion, ok := entity["api_version"].(string); ok {
			entityMap["api_version"] = apiVersion
		}

		// Handle status (as a list)
		statusMap := make(map[string]interface{})
		if status, ok := entity["status"].(map[string]interface{}); ok {
			// Handle description, name, state, availability_zone_reference, cluster_reference from status
			if description, ok := status["description"].(string); ok {
				statusMap["description"] = description
			}

			if name, ok := status["name"].(string); ok {
				statusMap["name"] = name
			}

			if availabilityZoneRef, ok := status["availability_zone_reference"].(map[string]interface{}); ok {
				azMap := make(map[string]interface{})
				if azKind, ok := availabilityZoneRef["kind"].(string); ok {
					azMap["kind"] = azKind
				}
				if azName, ok := availabilityZoneRef["name"].(string); ok {
					azMap["name"] = azName
				}
				if azUUID, ok := availabilityZoneRef["uuid"].(string); ok {
					azMap["uuid"] = azUUID
				}
				statusMap["availability_zone_reference"] = azMap
			}

			if clusterRef, ok := status["cluster_reference"].(map[string]interface{}); ok {
				clusterRefMap := make(map[string]interface{})
				if clusterKind, ok := clusterRef["kind"].(string); ok {
					clusterRefMap["kind"] = clusterKind
				}
				if clusterName, ok := clusterRef["name"].(string); ok {
					clusterRefMap["name"] = clusterName
				}
				if clusterUUID, ok := clusterRef["uuid"].(string); ok {
					clusterRefMap["uuid"] = clusterUUID
				}
				statusMap["cluster_reference"] = clusterRefMap
			}

			// Handle messages_list
			if messagesList, ok := status["messages_list"].([]interface{}); ok && len(messagesList) > 0 {
				messages := make([]map[string]interface{}, 0)
				for _, msg := range messagesList {
					if msgMap, ok := msg.(map[string]interface{}); ok {
						messageMap := make(map[string]interface{})
						if message, ok := msgMap["message"].(string); ok {
							messageMap["message"] = message
						}
						if reason, ok := msgMap["reason"].(string); ok {
							messageMap["reason"] = reason
						}
						if details, ok := msgMap["details"].(string); ok {
							messageMap["details"] = details
						}
						messages = append(messages, messageMap)
					}
				}
				statusMap["messages_list"] = messages
			}

			// resources (as list)
			// Skip custom_provider type
			resouceMap := make(map[string]interface{})
			if resources, ok := status["resources"].(map[string]interface{}); ok {
				if typeName, ok := resources["type"].(string); ok {
					if typeName == "custom_provider" || typeName == "nutanix" {
						continue // Skip
					}
					resouceMap["type"] = typeName
				}

				// Handle resources
				// Handle sync_interval_secs
				if syncInterval, ok := resources["sync_interval_secs"].(float64); ok {
					resouceMap["sync_interval_secs"] = int(syncInterval)
				}

				// Handle last_sync_time
				if lastSyncTime, ok := resources["last_sync_time"].(float64); ok {
					resouceMap["last_sync_time"] = lastSyncTime
				}

				// Handle sync_status
				if syncStatus, ok := resources["sync_status"].(string); ok {
					resouceMap["sync_status"] = syncStatus
				}

				// Handle sync_error
				if syncError, ok := resources["sync_error"].(string); ok {
					resouceMap["sync_error"] = syncError
				}

				// Handle state
				if state, ok := resources["state"].(string); ok {
					resouceMap["state"] = state
				}

				// Handle tunnel_reference and parent_reference
				if tunnelRef, ok := resources["tunnel_reference"].(map[string]interface{}); ok {
					tunnelRefMap := make(map[string]interface{})
					if tunnelKind, ok := tunnelRef["kind"].(string); ok {
						tunnelRefMap["kind"] = tunnelKind
					}
					if tunnelName, ok := tunnelRef["name"].(string); ok {
						tunnelRefMap["name"] = tunnelName
					}
					if tunnelUUID, ok := tunnelRef["uuid"].(string); ok {
						tunnelRefMap["uuid"] = tunnelUUID
					}
					resouceMap["tunnel_reference"] = []interface{}{tunnelRefMap}
				}

				if parentRef, ok := resources["parent_reference"].(map[string]interface{}); ok {
					parentRefMap := make(map[string]interface{})
					if parentKind, ok := parentRef["kind"].(string); ok {
						parentRefMap["kind"] = parentKind
					}
					if parentName, ok := parentRef["name"].(string); ok {
						parentRefMap["name"] = parentName
					}
					if parentUUID, ok := parentRef["uuid"].(string); ok {
						parentRefMap["uuid"] = parentUUID
					}
					resouceMap["parent_reference"] = []interface{}{parentRefMap}
				}

				// Handle price_items
				resouceMap["price_items"] = flattenPriceItems(resources["price_items"])

				// Handle data as a map (provider-specific)
				if data, ok := resources["data"].(map[string]interface{}); ok {
					// Initialize provider_type
					dataMap := make(map[string]interface{})
					provider := resources["type"].(string)
					log.Printf("[DEBUG] Akhilan Provider type: %s", provider)
					switch provider {
					case "aws":
						awsData := make(map[string]interface{})
						if accessKeyID, ok := data["access_key_id"].(string); ok {
							awsData["access_key_id"] = accessKeyID
						}
						if regions, ok := data["regions"].([]interface{}); ok {
							regionList := []map[string]interface{}{}
							for _, region := range regions {
								if r, ok := region.(map[string]interface{}); ok {
									regionMap := make(map[string]interface{})
									if name, ok := r["name"].(string); ok {
										regionMap["name"] = name
									}
									regionList = append(regionList, regionMap)
								}
							}
							awsData["regions"] = regionList
						}
						dataMap["aws"] = []interface{}{awsData}

					case "nutanix_pc":
						nutanixData := make(map[string]interface{})
						if server, ok := data["server"].(string); ok {
							nutanixData["server"] = server
						}
						if portRaw, ok := data["port"]; ok {
							switch v := portRaw.(type) {
							case float64:
								if v != 0 {
									nutanixData["port"] = int(v)
								}
							case int:
								if v != 0 {
									nutanixData["port"] = v
								}
							}
						}
						if username, ok := data["username"].(string); ok {
							nutanixData["username"] = username
						}
						nutanixData["cluster_account_reference_list"] = flattenClusterAccountReferenceList(data["cluster_account_reference_list"])
						dataMap["nutanix"] = []interface{}{nutanixData}

					case "azure":
						azureData := make(map[string]interface{})
						if tenantID, ok := data["tenant_id"].(string); ok {
							azureData["tenant_id"] = tenantID
						}
						if clientID, ok := data["client_id"].(string); ok {
							azureData["client_id"] = clientID
						}
						if subscriptionID, ok := data["subscription_id"].(string); ok {
							azureData["subscription_id"] = subscriptionID
						}
						if cloudEnv, ok := data["cloud_environment"].(string); ok {
							azureData["cloud_environment"] = cloudEnv
						}
						dataMap["azure"] = []interface{}{azureData}

					case "vmware":
						vmwData := make(map[string]interface{})
						if server, ok := data["server"].(string); ok {
							vmwData["server"] = server
						}
						if portRaw, ok := data["port"]; ok {
							switch v := portRaw.(type) {
							case float64:
								if v != 0 {
									vmwData["port"] = int(v)
								}
							case int:
								if v != 0 {
									vmwData["port"] = v
								}
							}
						}
						if datacenter, ok := data["datacenter"].(string); ok {
							vmwData["datacenter"] = datacenter
						}
						if username, ok := data["username"].(string); ok {
							vmwData["username"] = username
						}
						dataMap["vmware"] = []interface{}{vmwData}

					case "gcp":
						gcpData := make(map[string]interface{})
						if accountType, ok := data["account_type"].(string); ok {
							gcpData["account_type"] = accountType
						}
						if projectID, ok := data["project_id"].(string); ok {
							gcpData["project_id"] = projectID
						}
						if privateKeyID, ok := data["private_key_id"].(string); ok {
							gcpData["private_key_id"] = privateKeyID
						}
						if clientEmail, ok := data["client_email"].(string); ok {
							gcpData["client_email"] = clientEmail
						}
						if clientID, ok := data["client_id"].(string); ok {
							gcpData["client_id"] = clientID
						}
						dataMap["gcp"] = []interface{}{gcpData}

					case "k8s":
						k8sData := make(map[string]interface{})
						if k8sType, ok := data["type"].(string); ok {
							k8sData["type"] = k8sType
						}
						if server, ok := data["server"].(string); ok {
							k8sData["server"] = server
						}
						if portRaw, ok := data["port"]; ok {
							switch v := portRaw.(type) {
							case float64:
								if v != 0 {
									k8sData["port"] = int(v)
								}
							case int:
								if v != 0 {
									k8sData["port"] = v
								}
							}
						}

						if authMap, ok := data["authentication"].(map[string]interface{}); ok {
							authData := make(map[string]interface{})
							if authType, ok := authMap["type"].(string); ok {
								authData["type"] = authType
							}
							if username, ok := authMap["username"].(string); ok {
								authData["username"] = username
							}
							k8sData["authentication"] = []interface{}{authData}
						}
						dataMap["k8s"] = []interface{}{k8sData}
					}
					resouceMap["data"] = []interface{}{dataMap}
				}
				statusMap["resources"] = []interface{}{resouceMap}
			}
		}
		entityMap["status"] = []interface{}{statusMap}
		entityList = append(entityList, entityMap)
	}
	return entityList
}

func flattenPriceItems(priceItemsRaw interface{}) []map[string]interface{} {
	priceItemList := make([]map[string]interface{}, 0)

	priceItems, ok := priceItemsRaw.([]interface{})
	if !ok || len(priceItems) == 0 {
		return priceItemList
	}

	for _, item := range priceItems {
		if itemMap, ok := item.(map[string]interface{}); ok {
			priceItem := make(map[string]interface{})

			if uuid, ok := itemMap["uuid"].(string); ok {
				priceItem["uuid"] = uuid
			}
			if name, ok := itemMap["name"].(string); ok {
				priceItem["name"] = name
			}
			if state, ok := itemMap["state"].(string); ok {
				priceItem["state"] = state
			}
			if description, ok := itemMap["description"].(string); ok {
				priceItem["description"] = description
			}

			// Handle state_cost_list
			if stateCostList, ok := itemMap["state_cost_list"].([]interface{}); ok && len(stateCostList) > 0 {
				costList := make([]map[string]interface{}, 0)
				for _, costItem := range stateCostList {
					if costMap, ok := costItem.(map[string]interface{}); ok {
						cost := make(map[string]interface{})
						if state, ok := costMap["state"].(string); ok {
							cost["state"] = state
						}

						if costListItems, ok := costMap["cost_list"].([]interface{}); ok && len(costListItems) > 0 {
							costs := make([]map[string]interface{}, 0)
							for _, costItem := range costListItems {
								if costItemMap, ok := costItem.(map[string]interface{}); ok {
									costDetail := make(map[string]interface{})
									if interval, ok := costItemMap["interval"].(string); ok {
										costDetail["interval"] = interval
									}
									if name, ok := costItemMap["name"].(string); ok {
										costDetail["name"] = name
									}
									if value, ok := costItemMap["value"].(float64); ok {
										costDetail["value"] = value
									}
									if costType, ok := costItemMap["type"].(string); ok {
										costDetail["type"] = costType
									}
									costs = append(costs, costDetail)
								}
							}
							cost["cost_list"] = costs
						}
						costList = append(costList, cost)
					}
				}
				priceItem["state_cost_list"] = costList
			}

			// Handle details
			if details, ok := itemMap["details"].([]interface{}); ok && len(details) > 0 {
				if detailMapItem, ok := details[0].(map[string]interface{}); ok {
					detailMap := make(map[string]interface{})
					if associationType, ok := detailMapItem["association_type"].(string); ok {
						detailMap["association_type"] = associationType
					}
					if occurrence, ok := detailMapItem["occurrence"].(string); ok {
						detailMap["occurrence"] = occurrence
					}
					if providerType, ok := detailMapItem["provider_type"].(string); ok {
						detailMap["provider_type"] = providerType
					}
					priceItem["details"] = detailMap
				}
			}
			// Handle messages_list
			priceItem["message_list"] = flattenMessageList(itemMap["message_list"])
			priceItemList = append(priceItemList, priceItem)
		}
	}
	return priceItemList
}

func flattenMessageList(messageListRaw interface{}) []map[string]interface{} {
	messages := make([]map[string]interface{}, 0)

	messageList, ok := messageListRaw.([]interface{})
	if !ok || len(messageList) == 0 {
		return messages
	}

	for _, msg := range messageList {
		if msgMap, ok := msg.(map[string]interface{}); ok {
			messageMap := make(map[string]interface{})
			if message, ok := msgMap["message"].(string); ok {
				messageMap["message"] = message
			}
			if reason, ok := msgMap["reason"].(string); ok {
				messageMap["reason"] = reason
			}
			if details, ok := msgMap["details"].(string); ok {
				messageMap["details"] = details
			}
			messages = append(messages, messageMap)
		}
	}
	return messages
}

func flattenClusterAccountReferenceList(refListRaw interface{}) []map[string]interface{} {
	referenceList := make([]map[string]interface{}, 0)

	refList, ok := refListRaw.([]interface{})
	if !ok || len(refList) == 0 {
		return referenceList
	}

	for _, ref := range refList {
		if refMap, ok := ref.(map[string]interface{}); ok {
			clusterRef := make(map[string]interface{})

			if name, ok := refMap["name"].(string); ok {
				clusterRef["name"] = name
			}
			if description, ok := refMap["description"].(string); ok {
				clusterRef["description"] = description
			}
			clusterRef["message_list"] = flattenMessageList(refMap["message_list"])

			// Handle resources
			if resources, ok := refMap["resources"].(map[string]interface{}); ok {
				clusterRefResources := make(map[string]interface{})
				clusterRefResources["price_items"] = flattenPriceItems(resources["price_items"])
				// Handle data
				if clusterRefResdata, ok := resources["data"].(map[string]interface{}); ok {
					clusterRefResourceDataMap := make(map[string]interface{})
					if clusterName, ok := clusterRefResdata["cluster_name"].(string); ok {
						clusterRefResourceDataMap["cluster_name"] = clusterName
					}
					if clusterUUID, ok := clusterRefResdata["cluster_uuid"].(string); ok {
						clusterRefResourceDataMap["cluster_uuid"] = clusterUUID
					}
					if pcAccountUUID, ok := clusterRefResdata["pc_account_uuid"].(string); ok {
						clusterRefResourceDataMap["pc_account_uuid"] = pcAccountUUID
					}
					clusterRefResources["data"] = []interface{}{clusterRefResourceDataMap}
				}
				if typ, ok := resources["type"].(string); ok {
					clusterRefResources["type"] = typ
				}
				if state, ok := resources["state"].(string); ok {
					clusterRefResources["state"] = state
				}
				if syncInterval, ok := resources["sync_interval_secs"].(float64); ok {
					clusterRefResources["sync_interval_secs"] = int(syncInterval)
				}
				clusterRef["resources"] = []interface{}{clusterRefResources}
			}

			if uuid, ok := refMap["uuid"].(string); ok {
				clusterRef["uuid"] = uuid
			}

			referenceList = append(referenceList, clusterRef)
		}
	}
	return referenceList
}
