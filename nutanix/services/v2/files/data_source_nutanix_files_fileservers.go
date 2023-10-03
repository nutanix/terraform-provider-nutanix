package files

// import (
// 	"context"
// 	"encoding/json"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/files-go-client/v4/models/files/v4/config"
// 		conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"

// 	"github.com/terraform-providers/terraform-provider-nutanix/utils"
// )

// func DatasourceNutanixFilesServers() *schema.Resource {
// 	return &schema.Resource{
// 		ReadContext: datasourceNutanixFilesServersRead,
// 		Schema: map[string]*schema.Schema{
// 			"page": {
// 				Type:     schema.TypeInt,
// 				Optional: true,
// 			},
// 			"select": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"order_by": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"filter": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"limit": {
// 				Type:     schema.TypeInt,
// 				Optional: true,
// 			},
// 			"file_servers": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"name": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"cvm_ip_addresses": {
// 							Type:     schema.TypeList,
// 							Computed: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"prefix_length": {
// 										Type:     schema.TypeInt,
// 										Computed: true,
// 									},
// 									"value": {
// 										Type:     schema.TypeString,
// 										Computed: true,
// 									},
// 								},
// 							},
// 						},
// 						"ext_id": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"links": {
// 							Type:     schema.TypeList,
// 							Computed: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"href": {
// 										Type:     schema.TypeString,
// 										Computed: true,
// 									},
// 									"rel": {
// 										Type:     schema.TypeString,
// 										Computed: true,
// 									},
// 								},
// 							},
// 						},
// 						"memory_gib": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
// 						},
// 						"vcpus": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
// 						},
// 						"vms": {
// 							Type:     schema.TypeList,
// 							Computed: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"ext_id": {
// 										Type:     schema.TypeString,
// 										Computed: true,
// 									},
// 									"fsvm_uuid": {
// 										Type:     schema.TypeString,
// 										Computed: true,
// 									},
// 									"memory_gib": {
// 										Type:     schema.TypeInt,
// 										Computed: true,
// 									},
// 									"name": {
// 										Type:     schema.TypeString,
// 										Computed: true,
// 									},
// 									"vcpus": {
// 										Type:     schema.TypeInt,
// 										Computed: true,
// 									},
// 									"external_networks": {
// 										Type:     schema.TypeList,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"ip_addresses": {
// 													Type:     schema.TypeList,
// 													Computed: true,
// 													Elem: &schema.Resource{
// 														Schema: map[string]*schema.Schema{
// 															"ipv4": {
// 																Type:     schema.TypeList,
// 																Computed: true,
// 																Elem: &schema.Resource{
// 																	Schema: map[string]*schema.Schema{
// 																		"value": {
// 																			Type:     schema.TypeString,
// 																			Computed: true,
// 																		},
// 																	},
// 																},
// 															},
// 														},
// 													},
// 												},
// 											},
// 										},
// 									},
// 									"internal_networks": {
// 										Type:     schema.TypeList,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"ip_addresses": {
// 													Type:     schema.TypeList,
// 													Computed: true,
// 													Elem: &schema.Resource{
// 														Schema: map[string]*schema.Schema{
// 															"ipv4": {
// 																Type:     schema.TypeList,
// 																Computed: true,
// 																Elem: &schema.Resource{
// 																	Schema: map[string]*schema.Schema{
// 																		"value": {
// 																			Type:     schema.TypeString,
// 																			Computed: true,
// 																		},
// 																	},
// 																},
// 															},
// 														},
// 													},
// 												},
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 						"nvms_count": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
// 						},
// 						"size_in_gib": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
// 						},
// 						"version": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"cluster_ext_id": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"dns_domain_name": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"dns_servers": {
// 							Type:     schema.TypeList,
// 							Computed: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"value": {
// 										Type:     schema.TypeString,
// 										Computed: true,
// 									},
// 								},
// 							},
// 						},
// 						"ntp_servers": {
// 							Type:     schema.TypeList,
// 							Computed: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"fqdn": {
// 										Type:     schema.TypeList,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"value": {
// 													Type:     schema.TypeString,
// 													Computed: true,
// 												},
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 						"external_networks": {
// 							Type:     schema.TypeList,
// 							Computed: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"default_gateway": {
// 										Type:     schema.TypeList,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"ipv4": {
// 													Type:     schema.TypeList,
// 													Computed: true,
// 													Elem: &schema.Resource{
// 														Schema: map[string]*schema.Schema{
// 															"prefix_length": {
// 																Type:     schema.TypeInt,
// 																Computed: true,
// 															},
// 															"value": {
// 																Type:     schema.TypeString,
// 																Computed: true,
// 															},
// 														},
// 													},
// 												},
// 											},
// 										},
// 									},
// 									"ip_addresses": {
// 										Type:     schema.TypeList,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"ipv4": {
// 													Type:     schema.TypeList,
// 													Computed: true,
// 													Elem: &schema.Resource{
// 														Schema: map[string]*schema.Schema{
// 															"prefix_length": {
// 																Type:     schema.TypeInt,
// 																Computed: true,
// 															},
// 															"value": {
// 																Type:     schema.TypeString,
// 																Computed: true,
// 															},
// 														},
// 													},
// 												},
// 											},
// 										},
// 									},
// 									"subnet_mask": {
// 										Type:     schema.TypeList,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"ipv4": {
// 													Type:     schema.TypeList,
// 													Computed: true,
// 													Elem: &schema.Resource{
// 														Schema: map[string]*schema.Schema{
// 															"value": {
// 																Type:     schema.TypeString,
// 																Computed: true,
// 															},
// 														},
// 													},
// 												},
// 											},
// 										},
// 									},
// 									"virtual_network_name": {
// 										Type:     schema.TypeString,
// 										Computed: true,
// 									},
// 									"is_managed": {
// 										Type:     schema.TypeBool,
// 										Computed: true,
// 									},
// 									"network_ext_id": {
// 										Type:     schema.TypeString,
// 										Computed: true,
// 									},
// 								},
// 							},
// 						},
// 						"internal_networks": {
// 							Type:     schema.TypeList,
// 							Computed: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"default_gateway": {
// 										Type:     schema.TypeList,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"ipv4": {
// 													Type:     schema.TypeList,
// 													Computed: true,
// 													Elem: &schema.Resource{
// 														Schema: map[string]*schema.Schema{
// 															"prefix_length": {
// 																Type:     schema.TypeInt,
// 																Computed: true,
// 															},
// 															"value": {
// 																Type:     schema.TypeString,
// 																Computed: true,
// 															},
// 														},
// 													},
// 												},
// 											},
// 										},
// 									},
// 									"ip_addresses": {
// 										Type:     schema.TypeList,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"ipv4": {
// 													Type:     schema.TypeList,
// 													Computed: true,
// 													Elem: &schema.Resource{
// 														Schema: map[string]*schema.Schema{
// 															"prefix_length": {
// 																Type:     schema.TypeInt,
// 																Computed: true,
// 															},
// 															"value": {
// 																Type:     schema.TypeString,
// 																Computed: true,
// 															},
// 														},
// 													},
// 												},
// 											},
// 										},
// 									},
// 									"subnet_mask": {
// 										Type:     schema.TypeList,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"ipv4": {
// 													Type:     schema.TypeList,
// 													Computed: true,
// 													Elem: &schema.Resource{
// 														Schema: map[string]*schema.Schema{
// 															"value": {
// 																Type:     schema.TypeString,
// 																Computed: true,
// 															},
// 														},
// 													},
// 												},
// 											},
// 										},
// 									},
// 									"virtual_network_name": {
// 										Type:     schema.TypeString,
// 										Computed: true,
// 									},
// 									"virtual_ip_address": {
// 										Type:     schema.TypeList,
// 										Computed: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"ipv4": {
// 													Type:     schema.TypeList,
// 													Computed: true,
// 													Elem: &schema.Resource{
// 														Schema: map[string]*schema.Schema{
// 															"prefix_length": {
// 																Type:     schema.TypeInt,
// 																Computed: true,
// 															},
// 															"value": {
// 																Type:     schema.TypeString,
// 																Computed: true,
// 															},
// 														},
// 													},
// 												},
// 											},
// 										},
// 									},
// 									"is_managed": {
// 										Type:     schema.TypeBool,
// 										Computed: true,
// 									},
// 									"network_ext_id": {
// 										Type:     schema.TypeString,
// 										Computed: true,
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// func datasourceNutanixFilesServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	conn := meta.(*conns.Client).Files

// 	var pagef, limitf *int
// 	var selectf, order_byf, filterf *string

// 	if page, ok := d.GetOk("page"); ok {
// 		pagef = utils.IntPtr(page.(int))
// 	} else {
// 		pagef = nil
// 	}

// 	if limit, ok := d.GetOk("limit"); ok {
// 		limitf = utils.IntPtr(limit.(int))
// 	} else {
// 		limitf = nil
// 	}

// 	if sel, ok := d.GetOk("select"); ok {
// 		selectf = utils.StringPtr(sel.(string))
// 	} else {
// 		selectf = nil
// 	}

// 	if ord, ok := d.GetOk("order_by"); ok {
// 		order_byf = utils.StringPtr(ord.(string))
// 	} else {
// 		order_byf = nil
// 	}

// 	if fil, ok := d.GetOk("filter"); ok {
// 		filterf = utils.StringPtr(fil.(string))
// 	} else {
// 		filterf = nil
// 	}

// 	resp, err := conn.FilesServerAPI.GetFileServers(pagef, limitf, filterf, order_byf, selectf)
// 	if err != nil {
// 		var errordata map[string]interface{}
// 		e := json.Unmarshal([]byte(err.Error()), &errordata)
// 		if e != nil {
// 			return diag.FromErr(e)
// 		}
// 		data := errordata["data"].(map[string]interface{})
// 		errorList := data["error"].([]interface{})
// 		errorMessage := errorList[0].(map[string]interface{})
// 		return diag.Errorf("error while fetching fileservers: %v", errorMessage["message"])
// 	}

// 	filesResp := resp.Data.GetValue().([]import1.FileServer)

// 	if err := d.Set("file_servers", flattenFileServers(filesResp)); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	d.SetId(resource.UniqueId())
// 	return nil
// }

// func flattenFileServers(pr []import1.FileServer) []map[string]interface{} {
// 	if len(pr) > 0 {
// 		fs := make([]map[string]interface{}, len(pr))

// 		for k, v := range pr {
// 			file := make(map[string]interface{})

// 			file["name"] = v.Name
// 			file["ext_id"] = v.ExtId
// 			file["memory_gib"] = v.MemoryGib
// 			file["vcpus"] = v.Vcpus
// 			file["nvms_count"] = v.NvmsCount
// 			file["size_in_gib"] = v.SizeInGib
// 			file["version"] = v.Version
// 			file["cluster_ext_id"] = v.ClusterExtId
// 			file["dns_domain_name"] = v.DnsDomainName
// 			file["cvm_ip_addresses"] = flattenFilesIPAddress(v.CvmIpAddresses)
// 			// file["links"] = flattenLinks(v.Links)
// 			// file["vms"] = flattenVMs(v.Vms)
// 			file["dns_servers"] = flattenFilesIPAddress(v.DnsServers)
// 			file["ntp_servers"] = flattenNTPServers(v.NtpServers)
// 			file["external_networks"] = flattenExtIntNetworks(v.ExternalNetworks)
// 			file["internal_networks"] = flattenExtIntNetworks(v.InternalNetworks)

// 			fs[k] = file
// 		}
// 		return fs
// 	}
// 	return nil
// }
