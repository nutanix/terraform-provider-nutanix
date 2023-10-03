package files

// import (
// 	"context"
// 	"encoding/json"
// 	"time"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/files-go-client/v4/api"
// 	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/files-go-client/v4/client"
// 	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/files-go-client/v4/models/common/v1/config"
// 	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/files-go-client/v4/models/common/v1/response"
// 	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/files-go-client/v4/models/files/v4/config"
// 	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/files-go-client/v4/models/prism/v4/config"
// 		conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"

// 	"github.com/terraform-providers/terraform-provider-nutanix/utils"
// )

// var (
// 	FilesTimeout      = 30 * time.Minute
// 	ApiClientInstance *client.ApiClient
// 	GetFilesAPI       *api.FileServerApi
// )

// func ResourceNutanixFilesServer() *schema.Resource {
// 	return &schema.Resource{
// 		CreateContext: resourceNutanixFilesServerCreate,
// 		ReadContext:   resourceNutanixFilesServerRead,
// 		UpdateContext: resourceNutanixFilesServerUpdate,
// 		DeleteContext: resourceNutanixFilesServerDelete,
// 		Timeouts: &schema.ResourceTimeout{
// 			Create: schema.DefaultTimeout(FilesTimeout),
// 			Delete: schema.DefaultTimeout(FilesTimeout),
// 		},
// 		Schema: map[string]*schema.Schema{
// 			"name": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"version": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"memory_gib": {
// 				Type:     schema.TypeInt,
// 				Required: true,
// 			},
// 			"vcpus": {
// 				Type:     schema.TypeInt,
// 				Required: true,
// 			},
// 			"size_in_gib": {
// 				Type:     schema.TypeInt,
// 				Required: true,
// 			},
// 			"nvms_count": {
// 				Type:     schema.TypeInt,
// 				Required: true,
// 			},
// 			"cvm_ip_address": {
// 				Type:     schema.TypeList,
// 				Optional: true,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"value": {
// 							Type:     schema.TypeString,
// 							Optional: true,
// 						},
// 					},
// 				},
// 			},
// 			"dns_domain_name": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"dns_servers": {
// 				Type:     schema.TypeList,
// 				Optional: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"value": {
// 							Type:     schema.TypeString,
// 							Optional: true,
// 						},
// 					},
// 				},
// 			},
// 			"ntp_servers": {
// 				Type:     schema.TypeList,
// 				Optional: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"value": {
// 							Type:     schema.TypeString,
// 							Optional: true,
// 						},
// 					},
// 				},
// 			},
// 			"cluster_ext_id": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"external_networks": {
// 				Type:     schema.TypeList,
// 				Optional: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"is_managed": {
// 							Type:     schema.TypeBool,
// 							Optional: true,
// 						},
// 						"network_ext_id": {
// 							Type:     schema.TypeString,
// 							Optional: true,
// 						},
// 						"subnet_mask": {
// 							Type:     schema.TypeList,
// 							Optional: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"ipv4": {
// 										Type:     schema.TypeList,
// 										Optional: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"value": {
// 													Type:     schema.TypeString,
// 													Optional: true,
// 												},
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 						"default_gateway": {
// 							Type:     schema.TypeList,
// 							Optional: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"ipv4": {
// 										Type:     schema.TypeList,
// 										Optional: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"value": {
// 													Type:     schema.TypeString,
// 													Optional: true,
// 												},
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 						"ip_addresses": {
// 							Type:     schema.TypeList,
// 							Optional: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"ipv4": {
// 										Type:     schema.TypeList,
// 										Optional: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"value": {
// 													Type:     schema.TypeString,
// 													Optional: true,
// 												},
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			"internal_networks": {
// 				Type:     schema.TypeList,
// 				Optional: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"is_managed": {
// 							Type:     schema.TypeBool,
// 							Optional: true,
// 						},
// 						"network_ext_id": {
// 							Type:     schema.TypeString,
// 							Optional: true,
// 						},
// 						"subnet_mask": {
// 							Type:     schema.TypeList,
// 							Optional: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"ipv4": {
// 										Type:     schema.TypeList,
// 										Optional: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"value": {
// 													Type:     schema.TypeString,
// 													Optional: true,
// 												},
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 						"default_gateway": {
// 							Type:     schema.TypeList,
// 							Optional: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"ipv4": {
// 										Type:     schema.TypeList,
// 										Optional: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"value": {
// 													Type:     schema.TypeString,
// 													Optional: true,
// 												},
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 						"ip_addresses": {
// 							Type:     schema.TypeList,
// 							Optional: true,
// 							Elem: &schema.Resource{
// 								Schema: map[string]*schema.Schema{
// 									"ipv4": {
// 										Type:     schema.TypeList,
// 										Optional: true,
// 										Elem: &schema.Resource{
// 											Schema: map[string]*schema.Schema{
// 												"value": {
// 													Type:     schema.TypeString,
// 													Optional: true,
// 												},
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			"file_blocking_extensions": {
// 				Type:     schema.TypeList,
// 				Optional: true,
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 			},
// 			"delete_fs": {
// 				Type:     schema.TypeList,
// 				Optional: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"force_delete": {
// 							Type:     schema.TypeBool,
// 							Optional: true,
// 						},
// 						"delete_pd_snapshots_schedules": {
// 							Type:     schema.TypeBool,
// 							Optional: true,
// 						},
// 						"delete_container": {
// 							Type:     schema.TypeBool,
// 							Optional: true,
// 						},
// 					},
// 				},
// 			},
// 			"links": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"href": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"rel": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 					},
// 				},
// 			},
// 			"vms": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"ext_id": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"fsvm_uuid": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"memory_gib": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
// 						},
// 						"name": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"vcpus": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
// 						},
// 					},
// 				},
// 			},
// 			"ext_id": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 		},
// 	}
// }

// func resourceNutanixFilesServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	conn := meta.(*conns.Client).Files

// 	filesSpec := &import1.FileServer{}
// 	fileName := ""
// 	clusterExtID := ""

// 	if name, nameok := d.GetOk("name"); nameok {
// 		filesSpec.Name = utils.StringPtr(name.(string))
// 		fileName = name.(string)
// 	}
// 	if version, vok := d.GetOk("version"); vok {
// 		filesSpec.Version = utils.StringPtr(version.(string))
// 	}
// 	if mem, mok := d.GetOk("memory_gib"); mok {
// 		filesSpec.MemoryGib = utils.Int64Ptr(int64(mem.(int)))
// 	}
// 	if vcpu, vcpuOk := d.GetOk("vcpus"); vcpuOk {
// 		filesSpec.Vcpus = utils.Int64Ptr(int64(vcpu.(int)))
// 	}
// 	if size, sok := d.GetOk("size_in_gib"); sok {
// 		filesSpec.SizeInGib = utils.Float64Ptr(float64(size.(int)))
// 	}
// 	if nvm, nok := d.GetOk("nvms_count"); nok {
// 		filesSpec.NvmsCount = utils.IntPtr(nvm.(int))
// 	}
// 	if dnsDomain, dok := d.GetOk("dns_domain_name"); dok {
// 		filesSpec.DnsDomainName = utils.StringPtr(dnsDomain.(string))
// 	}
// 	if clsID, clsok := d.GetOk("cluster_id"); clsok {
// 		filesSpec.ClusterExtId = utils.StringPtr(clsID.(string))
// 		clusterExtID = clsID.(string)
// 	}

// 	if cvm, ok := d.GetOk("cvm_ip_address"); ok {
// 		cvmlist := cvm.([]interface{})
// 		cvmsAdd := make([]config.IPv4Address, len(cvmlist))

// 		for k, v := range cvmlist {
// 			val := v.(map[string]interface{})

// 			if value, vok := val["value"]; vok {
// 				cvmsAdd[k].Value = utils.StringPtr(value.(string))
// 			}
// 		}
// 		filesSpec.CvmIpAddresses = cvmsAdd
// 	}

// 	if ntp, ok := d.GetOk("ntp_servers"); ok {
// 		ntplist := ntp.([]interface{})
// 		ntps := make([]config.IPAddressOrFQDN, len(ntplist))

// 		for k, v := range ntplist {
// 			fqdn := config.FQDN{}

// 			val := v.(map[string]interface{})

// 			if value, vok := val["value"]; vok {
// 				fqdn.Value = utils.StringPtr(value.(string))
// 			}

// 			ntps[k].Fqdn = &fqdn
// 		}
// 		filesSpec.NtpServers = ntps
// 	}

// 	if dns, ok := d.GetOk("dns_servers"); ok {
// 		dnslist := dns.([]interface{})
// 		dnsSpec := make([]config.IPv4Address, len(dnslist))

// 		for k, v := range dnslist {
// 			val := v.(map[string]interface{})

// 			if value, vok := val["value"]; vok {
// 				dnsSpec[k].Value = utils.StringPtr(value.(string))
// 			}
// 		}
// 		filesSpec.DnsServers = dnsSpec
// 	}

// 	if extNet, ok := d.GetOk("external_networks"); ok {
// 		filesSpec.ExternalNetworks = expandFileNetwork(extNet.([]interface{}))
// 	}

// 	if extNet, ok := d.GetOk("internal_networks"); ok {
// 		filesSpec.InternalNetworks = expandFileNetwork(extNet.([]interface{}))
// 	}

// 	if extID, ok := d.GetOk("ext_id"); ok {
// 		filesSpec.ExtId = utils.StringPtr(extID.(string))
// 	}

// 	resp, err := conn.FilesServerAPI.CreateFileServer(filesSpec)
// 	// resp, err := FilesAPI.CreateFileServer(filesSpec)

// 	if err != nil {
// 		var errordata map[string]interface{}
// 		e := json.Unmarshal([]byte(err.Error()), &errordata)
// 		if e != nil {
// 			return diag.FromErr(e)
// 		}
// 		data := errordata["data"].(map[string]interface{})
// 		errorList := data["error"].([]interface{})
// 		errorMessage := errorList[0].(map[string]interface{})
// 		return diag.Errorf("error while creating fileserver: %v", errorMessage["message"])
// 	}
// 	taskRef := resp.Data.GetValue().(import4.TaskReference)
// 	taskUUID := taskRef.ExtId
// 	// taskUUID := "a4399883-58ef-4fb6-53be-12df6d5d1b0a"

// 	// Wait for the data protection policy to be available
// 	stateConf := &resource.StateChangeConf{
// 		Pending:    []string{"RUNNING", "QUEUED"},
// 		Target:     []string{"SUCCEEDED"},
// 		Refresh:    taskStateRefreshPrismTaskGroupFunc(ctx, utils.StringValue(taskUUID), meta),
// 		Timeout:    d.Timeout(schema.TimeoutCreate),
// 		Delay:      5 * time.Second,
// 		MinTimeout: 10 * time.Second,
// 	}

// 	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
// 		return diag.Errorf("error waiting for file server (%s) to create: %s", taskUUID, errWaitTask)
// 	}

// 	fileServerID, er := getNewlyFileServerID(ctx, meta, fileName, clusterExtID)
// 	if er != nil {
// 		return er
// 	}
// 	d.SetId(fileServerID)
// 	return resourceNutanixFilesServerRead(ctx, d, meta)
// }

// func resourceNutanixFilesServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	conn := meta.(*conns.Client).Files

// 	ApiClientInstance := client.NewApiClient()
// 	ApiClientInstance.Host = conn.FilesServerAPI.ApiClient.Host
// 	ApiClientInstance.Port = conn.FilesServerAPI.ApiClient.Port
// 	ApiClientInstance.Username = conn.FilesServerAPI.ApiClient.Username
// 	ApiClientInstance.Password = conn.FilesServerAPI.ApiClient.Password
// 	ApiClientInstance.VerifySSL = false

// 	GetFilesAPI := api.NewFileServerApi(ApiClientInstance)
// 	if d.Id() == "" {
// 		return diag.Errorf("file server id cannot be empty")
// 	}

// 	resp, err := GetFilesAPI.GetFileServerByExtId(utils.StringPtr(d.Id()))
// 	if err != nil {
// 		var errordata map[string]interface{}
// 		e := json.Unmarshal([]byte(err.Error()), &errordata)
// 		if e != nil {
// 			return diag.FromErr(e)
// 		}
// 		data := errordata["data"].(map[string]interface{})
// 		errorList := data["error"].([]interface{})
// 		errorMessage := errorList[0].(map[string]interface{})
// 		return diag.Errorf("error while fetching fileserver: %v", errorMessage["message"])
// 	}

// 	filesResp := resp.Data.GetValue().(import1.FileServer)

// 	if err := d.Set("name", filesResp.Name); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("ext_id", filesResp.ExtId); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("memory_gib", filesResp.MemoryGib); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	if err := d.Set("vcpus", filesResp.Vcpus); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("nvms_count", filesResp.NvmsCount); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("size_in_gib", filesResp.SizeInGib); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	if err := d.Set("version", filesResp.Version); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("cluster_ext_id", filesResp.ClusterExtId); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("dns_domain_name", filesResp.DnsDomainName); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	if err := d.Set("cvm_ip_address", flattenFilesIPAddress(filesResp.CvmIpAddresses)); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("links", flattenLinks(filesResp.Links)); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("vms", flattenVMs(filesResp.Vms)); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	if err := d.Set("dns_servers", flattenFilesIPAddress(filesResp.DnsServers)); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	if err := d.Set("ntp_servers", flattenNTPServers(filesResp.NtpServers)); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("external_networks", flattenExtIntNetworks(filesResp.ExternalNetworks)); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("internal_networks", flattenExtIntNetworks(filesResp.InternalNetworks)); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	return nil
// }

// func resourceNutanixFilesServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	conn := meta.(*conns.Client).Files
// 	updatefilesSpec := &import1.FileServer{}

// 	getResp, er := conn.FilesServerAPI.GetFileServerByExtId(utils.StringPtr(d.Id()))
// 	if er != nil {
// 		return diag.FromErr(er)
// 	}

// 	getFileServerResp := getResp.Data.GetValue().(import1.FileServer)

// 	updatefilesSpec = &getFileServerResp

// 	if d.HasChange("name") {
// 		updatefilesSpec.Name = utils.StringPtr(d.Get("name").(string))
// 	}
// 	if d.HasChange("dns_domain_name") {
// 		updatefilesSpec.DnsDomainName = utils.StringPtr(d.Get("dns_domain_name").(string))
// 	}

// 	if d.HasChange("memory_gib") {
// 		updatefilesSpec.MemoryGib = utils.Int64Ptr(int64(d.Get("memory_gib").(int)))
// 	}

// 	if d.HasChange("vcpus") {
// 		updatefilesSpec.Vcpus = utils.Int64Ptr(int64(d.Get("vcpus").(int)))
// 	}

// 	if d.HasChange("size_in_gib") {
// 		updatefilesSpec.SizeInGib = utils.Float64Ptr(float64(d.Get("size_in_gib").(int)))
// 	}

// 	if d.HasChange("file_blocking_extensions") {
// 		blockList := make([]string, 0)
// 		if fileBlock, ok := d.GetOk("file_blocking_extensions"); ok && len(fileBlock.([]interface{})) > 0 {
// 			for _, v := range fileBlock.([]interface{}) {
// 				blockList = append(blockList, v.(string))
// 			}
// 			updatefilesSpec.FileBlockingExtensions = blockList
// 		} else {
// 			updatefilesSpec.FileBlockingExtensions = []string{""}
// 		}
// 	}

// 	if d.HasChange("ntp_servers") {
// 		ntplist := d.Get("ntp_servers").([]interface{})
// 		ntps := make([]config.IPAddressOrFQDN, len(ntplist))

// 		for k, v := range ntplist {
// 			fqdn := config.FQDN{}

// 			val := v.(map[string]interface{})

// 			if value, vok := val["value"]; vok {
// 				fqdn.Value = utils.StringPtr(value.(string))
// 			}

// 			ntps[k].Fqdn = &fqdn
// 		}
// 		updatefilesSpec.NtpServers = ntps
// 	}

// 	if d.HasChange("dns_servers") {
// 		dnslist := d.Get("dns_servers").([]interface{})
// 		dnsSpec := make([]config.IPv4Address, len(dnslist))

// 		for k, v := range dnslist {
// 			val := v.(map[string]interface{})

// 			if value, vok := val["value"]; vok {
// 				dnsSpec[k].Value = utils.StringPtr(value.(string))
// 			}
// 		}
// 		updatefilesSpec.DnsServers = dnsSpec
// 	}

// 	// Extract E-Tag Header
// 	etagValue := ApiClientInstance.GetEtag(getResp)

// 	args := make(map[string]interface{})
// 	args["If-Match"] = etagValue

// 	resp, err := conn.FilesServerAPI.UpdateFileServer(utils.StringPtr(d.Id()), updatefilesSpec, args)
// 	if err != nil {
// 		var errordata map[string]interface{}
// 		e := json.Unmarshal([]byte(err.Error()), &errordata)
// 		if e != nil {
// 			return diag.FromErr(e)
// 		}
// 		data := errordata["data"].(map[string]interface{})
// 		errorList := data["error"].([]interface{})
// 		errorMessage := errorList[0].(map[string]interface{})
// 		return diag.Errorf("error while updating fileserver: %v", errorMessage["message"])
// 	}

// 	taskRef := resp.Data.GetValue().(import4.TaskReference)
// 	taskUUID := taskRef.ExtId

// 	// Wait for the data protection policy to be available
// 	stateConf := &resource.StateChangeConf{
// 		Pending:    []string{"RUNNING", "QUEUED"},
// 		Target:     []string{"SUCCEEDED"},
// 		Refresh:    taskStateRefreshPrismTaskGroupFunc(ctx, utils.StringValue(taskUUID), meta),
// 		Timeout:    d.Timeout(schema.TimeoutCreate),
// 		Delay:      5 * time.Second,
// 		MinTimeout: 10 * time.Second,
// 	}

// 	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
// 		return diag.Errorf("error waiting for file server (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
// 	}

// 	// update the networks of file server.

// 	if d.HasChange("external_networks") || d.HasChange("internal_networks") {
// 		// get the current file server response
// 		getResponse, er := conn.FilesServerAPI.GetFileServerByExtId(utils.StringPtr(d.Id()))
// 		if er != nil {
// 			return diag.FromErr(er)
// 		}

// 		getFileServerResponse := getResponse.Data.GetValue().(import1.FileServer)

// 		updatefilesSpec = &getFileServerResponse

// 		if d.HasChange("external_networks") {
// 			updatefilesSpec.ExternalNetworks = expandFileNetwork(d.Get("external_networks").([]interface{}))
// 		}
// 		if d.HasChange("internal_networks") {
// 			updatefilesSpec.InternalNetworks = expandFileNetwork(d.Get("internal_networks").([]interface{}))
// 		}

// 		// Extract E-Tag Header
// 		etagValueResp := ApiClientInstance.GetEtag(getResponse)

// 		argsVal := make(map[string]interface{})
// 		argsVal["If-Match"] = etagValueResp

// 		updateNetresp, err := conn.FilesServerAPI.UpdateFileServerNetworkConfig(utils.StringPtr(d.Id()), updatefilesSpec, nil, argsVal)
// 		if err != nil {
// 			var errordata map[string]interface{}
// 			e := json.Unmarshal([]byte(err.Error()), &errordata)
// 			if e != nil {
// 				return diag.FromErr(e)
// 			}
// 			data := errordata["data"].(map[string]interface{})
// 			errorList := data["error"].([]interface{})
// 			errorMessage := errorList[0].(map[string]interface{})
// 			return diag.Errorf("error while updating fileserver network config: %v", errorMessage["message"])
// 		}

// 		uptaskRef := updateNetresp.Data.GetValue().(import4.TaskReference)
// 		uptaskUUID := uptaskRef.ExtId

// 		// Wait for the data protection policy to be available
// 		networkstateConf := &resource.StateChangeConf{
// 			Pending:    []string{"RUNNING", "QUEUED"},
// 			Target:     []string{"SUCCEEDED"},
// 			Refresh:    taskStateRefreshPrismTaskGroupFunc(ctx, utils.StringValue(taskUUID), meta),
// 			Timeout:    d.Timeout(schema.TimeoutCreate),
// 			Delay:      5 * time.Second,
// 			MinTimeout: 10 * time.Second,
// 		}

// 		if _, errWaitTask := networkstateConf.WaitForStateContext(ctx); errWaitTask != nil {
// 			return diag.Errorf("error waiting for file server (%s) to update network configuration: %s", utils.StringValue(uptaskUUID), errWaitTask)
// 		}
// 	}

// 	return resourceNutanixFilesServerRead(ctx, d, meta)
// }

// func resourceNutanixFilesServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	conn := meta.(*conns.Client).Files

// 	// forcedelete := true
// 	deletePDSnaps := true
// 	deleteContainer := true

// 	if deletefs, dok := d.GetOk("delete_fs"); dok && len(deletefs.([]interface{})) > 0 {
// 		ds := deletefs.([]interface{})

// 		for _, v := range ds {
// 			val := v.(map[string]interface{})

// 			// if forceds, ok := val["force_delete"]; ok {
// 			// 	forcedelete = forceds.(bool)
// 			// }

// 			if pdsnaps, ok := val["delete_pd_snapshots_schedules"]; ok {
// 				deletePDSnaps = pdsnaps.(bool)
// 			}

// 			if deletecont, ok := val["delete_container"]; ok {
// 				deleteContainer = deletecont.(bool)
// 			}
// 		}
// 	}

// 	getResponse, err := conn.FilesServerAPI.GetFileServerByExtId(utils.StringPtr(d.Id()))
// 	if err != nil {
// 		var errordata map[string]interface{}
// 		e := json.Unmarshal([]byte(err.Error()), &errordata)
// 		if e != nil {
// 			return diag.FromErr(e)
// 		}
// 		data := errordata["data"].(map[string]interface{})
// 		errorList := data["error"].([]interface{})
// 		errorMessage := errorList[0].(map[string]interface{})
// 		return diag.Errorf("error while deleting fileserver: %v", errorMessage["message"])
// 	}
// 	// Extract E-Tag Header
// 	etagValue := ApiClientInstance.GetEtag(getResponse)

// 	args := make(map[string]interface{})
// 	args["If-Match"] = etagValue

// 	resp, err := conn.FilesServerAPI.DeleteFileServer(utils.StringPtr(d.Id()), nil, utils.BoolPtr(deletePDSnaps), utils.BoolPtr(deleteContainer), args)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	taskRef := resp.Data.GetValue().(import4.TaskReference)
// 	taskUUID := taskRef.ExtId

// 	// Wait for the data protection policy to be available
// 	stateConf := &resource.StateChangeConf{
// 		Pending:    []string{"RUNNING", "QUEUED"},
// 		Target:     []string{"SUCCEEDED"},
// 		Refresh:    taskStateRefreshPrismTaskGroupFunc(ctx, utils.StringValue(taskUUID), meta),
// 		Timeout:    d.Timeout(schema.TimeoutCreate),
// 		Delay:      5 * time.Second,
// 		MinTimeout: 10 * time.Second,
// 	}

// 	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
// 		return diag.Errorf("error waiting for file server (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
// 	}
// 	return nil
// }

// func expandIPv4Value(pr []interface{}) *config.IPAddress {
// 	if len(pr) > 0 {
// 		ipSpec := config.IPAddress{}
// 		for _, v := range pr {
// 			val := v.(map[string]interface{})
// 			if ipv, ok := val["ipv4"]; ok {
// 				ipSpec.Ipv4 = expandvalue(ipv.([]interface{}))
// 			}
// 		}
// 		return &ipSpec
// 	}
// 	return nil
// }

// func expandvalue(pr []interface{}) *config.IPv4Address {
// 	if len(pr) > 0 {
// 		ipv4Add := config.IPv4Address{}
// 		for _, v := range pr {
// 			val := v.(map[string]interface{})
// 			if ipv, ok := val["value"]; ok {
// 				ipv4Add.Value = utils.StringPtr(ipv.(string))
// 			}
// 		}
// 		return &ipv4Add
// 	}
// 	return nil
// }

// func expandFileNetwork(pr []interface{}) []import1.Network {
// 	if len(pr) > 0 {
// 		extNw := make([]import1.Network, len(pr))
// 		for k, v := range pr {
// 			val := v.(map[string]interface{})
// 			net := import1.Network{}

// 			if manage, ok := val["is_managed"]; ok {
// 				net.IsManaged = utils.BoolPtr(manage.(bool))
// 			}
// 			if extid, eok := val["network_ext_id"]; eok {
// 				net.NetworkExtId = utils.StringPtr(extid.(string))
// 			}
// 			if sub, ok := val["subnet_mask"]; ok {
// 				net.SubnetMask = expandIPv4Value(sub.([]interface{}))
// 			}
// 			if sub, ok := val["default_gateway"]; ok {
// 				net.DefaultGateway = expandIPv4Value(sub.([]interface{}))
// 			}
// 			if sub, ok := val["ip_addresses"]; ok {
// 				sublist := sub.([]interface{})
// 				ipAdd := make([]config.IPAddress, len(sublist))

// 				for k, v := range sublist {
// 					val := v.(map[string]interface{})

// 					if value, vok := val["ipv4"]; vok {
// 						ipAdd[k].Ipv4 = expandvalue(value.([]interface{}))
// 					}
// 				}
// 				net.IpAddresses = ipAdd
// 			}
// 			extNw[k] = net
// 		}
// 		return extNw
// 	}
// 	return nil
// }

// func flattenLinks(pr []import2.ApiLink) []map[string]interface{} {
// 	if len(pr) > 0 {
// 		linkList := make([]map[string]interface{}, len(pr))

// 		for k, v := range pr {
// 			links := map[string]interface{}{}
// 			if v.Href != nil {
// 				links["href"] = v.Href
// 			}
// 			if v.Rel != nil {
// 				links["rel"] = v.Rel
// 			}

// 			linkList[k] = links
// 		}
// 		return linkList
// 	}
// 	return nil
// }

// func flattenVMs(pr []import1.VM) []map[string]interface{} {
// 	if len(pr) > 0 {
// 		vmsList := make([]map[string]interface{}, len(pr))

// 		for k, v := range pr {
// 			vms := map[string]interface{}{}

// 			if v.ExtId != nil {
// 				vms["ext_id"] = v.ExtId
// 			}
// 			if v.FsvmUuid != nil {
// 				vms["fsvm_uuid"] = v.FsvmUuid
// 			}
// 			if v.MemoryGib != nil {
// 				vms["memory_gib"] = v.MemoryGib
// 			}
// 			if v.Vcpus != nil {
// 				vms["vcpus"] = v.Vcpus
// 			}
// 			if v.Name != nil {
// 				vms["name"] = v.Name
// 			}

// 			vmsList[k] = vms
// 		}
// 		return vmsList
// 	}
// 	return nil
// }

// func getNewlyFileServerID(ctx context.Context, meta interface{}, fileName, clusterExtID string) (string, diag.Diagnostics) {
// 	conn := meta.(*conns.Client).Files

// 	ApiClientInstance = client.NewApiClient()
// 	ApiClientInstance.Host = conn.FilesServerAPI.ApiClient.Host
// 	ApiClientInstance.Port = conn.FilesServerAPI.ApiClient.Port
// 	ApiClientInstance.Username = conn.FilesServerAPI.ApiClient.Username
// 	ApiClientInstance.Password = conn.FilesServerAPI.ApiClient.Password
// 	ApiClientInstance.VerifySSL = false

// 	GetFilesAPI := api.NewFileServerApi(ApiClientInstance)

// 	filesResp, er := GetFilesAPI.GetFileServers(nil, nil, nil, nil, nil)
// 	if er != nil {
// 		return "", diag.FromErr(er)
// 	}
// 	fileServerID := ""
// 	listofFiles := filesResp.Data.GetValue().([]import1.FileServer)

// 	for _, v := range listofFiles {
// 		if fileName == *v.Name && clusterExtID == *v.ClusterExtId {
// 			fileServerID = *v.ExtId
// 			break
// 		}
// 	}
// 	return fileServerID, nil
// }
