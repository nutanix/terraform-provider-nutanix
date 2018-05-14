package nutanix

import (
	"strconv"

	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixClustersRead,

		Schema: getDataSourceClustersSchema(),
	}
}

func dataSourceNutanixClustersRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*NutanixClient).API

	metadata := &v3.ClusterListMetadataOutput{}

	if v, ok := d.GetOk("metadata"); ok {
		m := v.(map[string]interface{})
		metadata.Kind = utils.String("vm")
		if mv, mok := m["sort_attribute"]; mok {
			metadata.SortAttribute = utils.String(mv.(string))
		}
		if mv, mok := m["filter"]; mok {
			metadata.Filter = utils.String(mv.(string))
		}
		if mv, mok := m["length"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return err
			}
			metadata.Length = utils.Int64(int64(i))
		}
		if mv, mok := m["sort_order"]; mok {
			metadata.SortOrder = utils.String(mv.(string))
		}
		if mv, mok := m["offset"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return err
			}
			metadata.Offset = utils.Int64(int64(i))
		}
	}

	// Make request to the API
	resp, err := conn.V3.ListCluster(metadata)
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}

	entities := make([]map[string]interface{}, len(resp.Entities))
	for k, v := range resp.Entities {
		entity := make(map[string]interface{})
		// set metadata values
		metadata := make(map[string]interface{})
		metadata["last_update_time"] = utils.TimeValue(v.Metadata.LastUpdateTime).String()
		metadata["kind"] = utils.StringValue(v.Metadata.Kind)
		metadata["uuid"] = utils.StringValue(v.Metadata.UUID)
		metadata["creation_time"] = utils.TimeValue(v.Metadata.CreationTime).String()
		metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(v.Metadata.SpecVersion)))
		metadata["spec_hash"] = utils.StringValue(v.Metadata.SpecHash)
		metadata["name"] = utils.StringValue(v.Metadata.Name)
		entity["metadata"] = metadata

		entity["categories"] = v.Metadata.Categories
		entity["api_version"] = utils.StringValue(v.APIVersion)

		pr := make(map[string]interface{})
		pr["kind"] = utils.StringValue(v.Metadata.ProjectReference.Kind)
		pr["name"] = utils.StringValue(v.Metadata.ProjectReference.Name)
		pr["uuid"] = utils.StringValue(v.Metadata.ProjectReference.UUID)

		entity["project_reference"] = pr

		or := make(map[string]interface{})
		or["kind"] = utils.StringValue(v.Metadata.OwnerReference.Kind)
		or["name"] = utils.StringValue(v.Metadata.OwnerReference.Name)
		or["uuid"] = utils.StringValue(v.Metadata.OwnerReference.UUID)
		entity["owner_reference"] = or
		entity["name"] = utils.StringValue(v.Status.Name)

		// TODO: set remaining attributes

		entities[k] = entity
	}

	if err := d.Set("entities", entities); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return nil
}

func getDataSourceClustersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"sort_attribute": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"filter": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"length": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"sort_order": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"offset": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"api_version": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"entities": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"metadata": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"last_update_time": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"creation_time": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"spec_version": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"spec_hash": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"categories": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
					"project_reference": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"owner_reference": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"api_version": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"description": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"availability_zone_reference": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"cluster_reference": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},

					// COMPUTED
					"state": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"nodes": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"version": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"type": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"gpu_driver_version": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"client_auth": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"status": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"ca_chain": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"authorized_public_key_list": &schema.Schema{
						Type: schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"software_map": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
					"encryption_status": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
					"ssl_key_type": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
					"ssl_key_name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"ssl_key_signing_info": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"city": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"common_name_suffix": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"state": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"country_code": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"common_name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"organization": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"email_address": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"ssl_key_expire_datetime": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"service_list": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
					"supported_information_verbosity": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"certification_signing_info": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"city": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"common_name_suffix": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"state": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"country_code": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"common_name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"organization": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"email_address": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"operation_mode": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"ca_certificate_list": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ca_name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"certificate": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"enabled_feature_list": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"is_available": &schema.Schema{
						Type:     schema.TypeBool,
						Computed: true,
					},
					"build": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"commit_id": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"full_version": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"commit_date": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"version": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"short_commit_id": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"build_type": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"timezone": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"cluster_arch": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"management_server_list": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"drs_enabled": &schema.Schema{
									Type:     schema.TypeBool,
									Computed: true,
								},
								"status_list": &schema.Schema{
									Type:     schema.TypeList,
									Computed: true,
									Elem:     &schema.Schema{Type: schema.TypeString},
								},
								"type": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"masquerading_port": &schema.Schema{
						Type:     schema.TypeInt,
						Computed: true,
					},
					"masquerading_ip": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"external_ip": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"http_proxy_list": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"credentials": &schema.Schema{
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"username": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"password": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"proxy_type_list": &schema.Schema{
									Type:     schema.TypeList,
									Computed: true,
									Elem:     &schema.Schema{Type: schema.TypeString},
								},
								"address": &schema.Schema{
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"ip": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"fqdn": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"port": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"ipv6": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
					"smtp_server_type": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"smtp_server_email_address": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"smtp_server_credentials": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"username": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"password": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"smtp_server_proxy_type_list": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"smtp_server_address": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"fqdn": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"port": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"ipv6": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"ntp_server_ip_list": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"external_subnet": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"external_data_services_ip": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"internal_subnet": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"domain_server_nameserver": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"domain_server_name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"domain_server_credentials": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"username": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"password": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"nfs_subnet_whitelist": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"name_server_ip_list": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"http_proxy_whitelist": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"target": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"target_type": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"analysis_vm_efficiency_map": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
	}
}
