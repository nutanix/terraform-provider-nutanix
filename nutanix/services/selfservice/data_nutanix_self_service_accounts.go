package selfservice

import (
	"context"

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
				Default:  "ASC",
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
										Type:     schema.TypeString,
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
					Type:     schema.TypeString,
					Computed: true,
				},
				"sync_status": {
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
				"messages_list": buildMessageListSchema(),
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
			},
		},
	}
}

func buildAzureProviderSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
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
		acFilter.Filter = filter.(string)
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
			Summary:  "🫙 No Data found",
			Detail:   "The API returned an empty list of Accounts.",
		}}
	}

	if err := d.Set("api_version", accountResp.APIVersion); err != nil {
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
			if specVersion, ok := metadata["spec_version"].(string); ok {
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
					if typeName == "custom_provider" {
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
				if lastSyncTime, ok := resources["last_sync_time"].(string); ok {
					resouceMap["last_sync_time"] = lastSyncTime
				}

				// Handle sync_status
				if syncStatus, ok := resources["sync_status"].(string); ok {
					resouceMap["sync_status"] = syncStatus
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
				if priceItems, ok := resources["price_items"].([]interface{}); ok && len(priceItems) > 0 {
					priceItemList := make([]map[string]interface{}, 0)

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
								// Using only the first item assuming one object
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
							if messagesList, ok := itemMap["messages_list"].([]interface{}); ok && len(messagesList) > 0 {
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
								priceItem["messages_list"] = messages
							}

							priceItemList = append(priceItemList, priceItem)
						}
					}
					resouceMap["price_items"] = priceItemList
				}

				// Handle data as a map (provider-specific)
				if data, ok := resources["data"].(map[string]interface{}); ok {
					// Initialize provider_type
					dataMap := make(map[string]interface{})
					provider := resources["type"].(string)
					switch provider {
					case "aws":
						if accessKeyID, ok := data["access_key_id"].(string); ok {
							dataMap["access_key_id"] = accessKeyID
						}
						if regions, ok := data["regions"].([]interface{}); ok {
							regionNames := []string{}
							for _, region := range regions {
								if r, ok := region.(map[string]interface{}); ok {
									if name, ok := r["name"].(string); ok {
										regionNames = append(regionNames, name)
									}
								}
							}
							dataMap["regions"] = regionNames
						}

					case "nutanix_pc":
						if server, ok := data["server"].(string); ok {
							dataMap["server"] = server
						}
						if port, ok := data["port"].(int); ok {
							dataMap["port"] = port
						}
						if username, ok := data["username"].(string); ok {
							dataMap["username"] = username
						}

					case "azure":
						if tenantID, ok := data["tenant_id"].(string); ok {
							dataMap["tenant_id"] = tenantID
						}
						if clientID, ok := data["client_id"].(string); ok {
							dataMap["client_id"] = clientID
						}
						if subscriptionID, ok := data["subscription_id"].(string); ok {
							dataMap["subscription_id"] = subscriptionID
						}
						if cloudEnv, ok := data["cloud_environment"].(string); ok {
							dataMap["cloud_environment"] = cloudEnv
						}

					case "vmware":
						if server, ok := data["server"].(string); ok {
							dataMap["server"] = server
						}
						if port, ok := data["port"].(int); ok {
							dataMap["port"] = port
						}
						if datacenter, ok := data["datacenter"].(string); ok {
							dataMap["datacenter"] = datacenter
						}
						if username, ok := data["username"].(string); ok {
							dataMap["username"] = username
						}

					case "gcp":
						if accountType, ok := data["account_type"].(string); ok {
							dataMap["account_type"] = accountType
						}
						if projectID, ok := data["project_id"].(string); ok {
							dataMap["project_id"] = projectID
						}
						if privateKeyID, ok := data["private_key_id"].(string); ok {
							dataMap["private_key_id"] = privateKeyID
						}
						if clientEmail, ok := data["client_email"].(string); ok {
							dataMap["client_email"] = clientEmail
						}
						if clientID, ok := data["client_id"].(string); ok {
							dataMap["client_id"] = clientID
						}

					case "k8s":
						if k8sType, ok := data["type"].(string); ok {
							dataMap["type"] = k8sType
						}
						if server, ok := data["server"].(string); ok {
							dataMap["server"] = server
						}
						if port, ok := data["port"].(int); ok {
							dataMap["port"] = port
						}
						if auth, ok := data["authentication"].([]interface{}); ok && len(auth) > 0 {
							authMap, _ := auth[0].(map[string]interface{})
							if authType, ok := authMap["type"].(string); ok {
								dataMap["auth_type"] = authType
							}
							if username, ok := authMap["username"].(string); ok {
								dataMap["username"] = username
							}
						}
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
