package lifecyclev2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
)

// -----------------------------------------------------------------------------
// Shared DATASOURCE (fully Computed) schema helpers for hardware-provider entities.
// These are used ONLY by datasources (singular + plural). Resources declare their
// own input-oriented schema independently (see HARD RULE 8).
// -----------------------------------------------------------------------------

// datasourceValuePrefixLengthSchema returns a Computed value/prefix_length block
// used for ipv4/ipv6 addresses in datasource schemas.
func datasourceValuePrefixLengthSchema() *schema.Schema {
	return &schema.Schema{
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
	}
}

// datasourceIPAddressSchema returns a Computed ipv4/ipv6 block.
func datasourceIPAddressSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": datasourceValuePrefixLengthSchema(),
				"ipv6": datasourceValuePrefixLengthSchema(),
			},
		},
	}
}

// datasourceConnectionAccessDetailsSchema returns the Computed access_details block.
func datasourceConnectionAccessDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auth": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"api_key_auth": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"api_key_id": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"api_key_secret": {
											Type:      schema.TypeString,
											Computed:  true,
											Sensitive: true,
										},
									},
								},
							},
							"basic_auth": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"username": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"password": {
											Type:      schema.TypeString,
											Computed:  true,
											Sensitive: true,
										},
									},
								},
							},
						},
					},
				},
				"endpoint": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"url_endpoint": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"url": {
											Type:     schema.TypeString,
											Computed: true,
										},
									},
								},
							},
							"ip_range_endpoint": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"ip_ranges": {
											Type:     schema.TypeList,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"start_ip": datasourceIPAddressSchema(),
													"end_ip":   datasourceIPAddressSchema(),
												},
											},
										},
									},
								},
							},
							"ip_address_endpoint": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"ip_addresses": {
											Type:     schema.TypeList,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"ipv4": datasourceValuePrefixLengthSchema(),
													"ipv6": datasourceValuePrefixLengthSchema(),
													"fqdn": {
														Type:     schema.TypeList,
														Computed: true,
														Elem: &schema.Resource{
															Schema: map[string]*schema.Schema{
																"value": {
																	Type:     schema.TypeString,
																	Computed: true,
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
						},
					},
				},
			},
		},
	}
}

// datasourcePoolGroupsSchema returns the Computed groups block shared by pools.
func datasourcePoolGroupsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// datasourcePoolDisplayNameMappingSchema returns the Computed display_name_mapping block for pools.
func datasourcePoolDisplayNameMappingSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"group": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// datasourceCPUInfoSchema returns the Computed cpu_info block for discovered nodes.
func datasourceCPUInfoSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"capacity_ghz": {
					Type:     schema.TypeFloat,
					Computed: true,
				},
				"logical_core_count": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"manufacturer": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"model": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"socket_count": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

// datasourceNodeIdentifiersSchema returns the Computed identifiers block for discovered nodes.
func datasourceNodeIdentifiersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// datasourceNodeProviderDataSchema returns the Computed provider_data block for discovered nodes.
func datasourceNodeProviderDataSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name_mapping": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"domain": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"groups": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"mode": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"tags": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"domain": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"groups": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
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
				"is_available": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"is_configured": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"mode": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"tags": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"value": {
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

// -----------------------------------------------------------------------------
// RESOURCE (input-oriented) schema helpers for the Connection resource. These carry
// Optional/Required semantics and MUST NOT be shared with the datasource schema.
// -----------------------------------------------------------------------------

// resourceValuePrefixLengthSchema returns an Optional value/prefix_length block
// used for ipv4/ipv6 addresses in resource schemas.
func resourceValuePrefixLengthSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
				"prefix_length": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

// resourceIPAddressSchema returns an Optional ipv4/ipv6 block for resource schemas.
func resourceIPAddressSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": resourceValuePrefixLengthSchema(),
				"ipv6": resourceValuePrefixLengthSchema(),
			},
		},
	}
}

// resourceConnectionAccessDetailsSchema returns the Optional access_details input block
// for the Connection resource.
func resourceConnectionAccessDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auth": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"api_key_auth": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"api_key_id": {
											Type:     schema.TypeString,
											Required: true,
										},
										"api_key_secret": {
											Type:      schema.TypeString,
											Optional:  true,
											Sensitive: true,
										},
									},
								},
							},
							"basic_auth": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"username": {
											Type:     schema.TypeString,
											Required: true,
										},
										"password": {
											Type:      schema.TypeString,
											Optional:  true,
											Sensitive: true,
										},
									},
								},
							},
						},
					},
				},
				"endpoint": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"url_endpoint": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"url": {
											Type:     schema.TypeString,
											Required: true,
										},
									},
								},
							},
							"ip_range_endpoint": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"ip_ranges": {
											Type:     schema.TypeList,
											Optional: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"start_ip": resourceIPAddressSchema(),
													"end_ip":   resourceIPAddressSchema(),
												},
											},
										},
									},
								},
							},
							"ip_address_endpoint": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"ip_addresses": {
											Type:     schema.TypeList,
											Optional: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"ipv4": resourceValuePrefixLengthSchema(),
													"ipv6": resourceValuePrefixLengthSchema(),
													"fqdn": {
														Type:     schema.TypeList,
														Optional: true,
														MaxItems: 1,
														Elem: &schema.Resource{
															Schema: map[string]*schema.Schema{
																"value": {
																	Type:     schema.TypeString,
																	Required: true,
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
						},
					},
				},
			},
		},
	}
}

// -----------------------------------------------------------------------------
// Top-level Computed schema maps for each datasource entity.
// -----------------------------------------------------------------------------

// connectionDatasourceSchema returns the fully-Computed schema for a Connection object.
func connectionDatasourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": common.LinksSchema(),
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"region": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deployment_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"access_details": datasourceConnectionAccessDetailsSchema(),
	}
}

// hardwareProviderDatasourceSchema returns the fully-Computed schema for a HardwareProvider.
func hardwareProviderDatasourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": common.LinksSchema(),
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vendor": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"available_connection_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"auth_types": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

// poolDatasourceSchema returns the fully-Computed schema shared by IP / MAC / server-identity
// pools. The countAttr differs per pool type (e.g. ipv4/ipv6 available counts vs. a single
// available count), so those are appended by the individual datasource files.
func poolDatasourceBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": common.LinksSchema(),
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"reference_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_updated_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"groups":               datasourcePoolGroupsSchema(),
		"display_name_mapping": datasourcePoolDisplayNameMappingSchema(),
	}
}

// ipPoolDatasourceSchema returns the fully-Computed schema for an IpPool.
func ipPoolDatasourceSchema() map[string]*schema.Schema {
	s := poolDatasourceBaseSchema()
	s["ipv4_available_count"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	s["ipv6_available_count"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	return s
}

// macPoolDatasourceSchema returns the fully-Computed schema for a MacPool.
func macPoolDatasourceSchema() map[string]*schema.Schema {
	s := poolDatasourceBaseSchema()
	s["available_count"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	return s
}

// serverIdentityPoolDatasourceSchema returns the fully-Computed schema for a ServerIdentityPool.
func serverIdentityPoolDatasourceSchema() map[string]*schema.Schema {
	s := poolDatasourceBaseSchema()
	s["available_count"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	return s
}

// discoveredNodeDatasourceSchema returns the fully-Computed schema for a DiscoveredNode.
func discoveredNodeDatasourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links":               common.LinksSchema(),
		"bmc_ip":              datasourceIPAddressSchema(),
		"cpu_info":            datasourceCPUInfoSchema(),
		"identifiers":         datasourceNodeIdentifiersSchema(),
		"provider_data":       datasourceNodeProviderDataSchema(),
		"last_refreshed_time": {Type: schema.TypeString, Computed: true},
		"managed_by": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"manufacturer": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"memory_gb": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"model": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
