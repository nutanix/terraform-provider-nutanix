package foundation

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DEFAULTWAITTIMEOUT = 30
)

func resourceFoundationImageNodes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFoundationImageNodesCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"xs_master_label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipmi_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cvm_gateway": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hyperv_external_vnic": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"xen_config_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ucsm_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ucsm_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hypervisor_iso": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hyperv": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filename": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"checksum": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"kvm": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filename": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"checksum": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"xen": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filename": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"checksum": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"esx": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filename": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"checksum": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"unc_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hypervisor_netmask": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fc_settings": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fc_metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fc_ip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"api_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"foundation_central": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"xs_master_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"svm_rescue_args": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cvm_netmask": {
				Type:     schema.TypeString,
				Required: true,
			},
			"xs_master_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_ns": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"backplane_subnet": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cluster_init_successful": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"backplane_netmask": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"redundancy_factor": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"backplane_vlan": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cluster_external_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cvm_ntp_servers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"single_node_cluster": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"cluster_members": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cvm_dns_servers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cluster_init_now": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"hypervisor_ntp_servers": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"hyperv_external_vswitch": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hypervisor_nameserver": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hyperv_sku": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"eos_metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"account_name": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"email": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tests": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"run_syscheck": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"run_ncc": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"blocks": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nodes": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv6_address": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"node_position": {
										Type:     schema.TypeString,
										Required: true,
									},
									"image_delay": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"ucsm_params": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"hypervisor_hostname": {
										Type:     schema.TypeString,
										Required: true,
									},
									"cvm_gb_ram": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"device_hint": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"bond_mode": {
										Type:     schema.TypeString,
										Required: true,
									},
									"rdma_passthrough": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"cluster_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ucsm_node_serial": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"hypervisor_ip": {
										Type:     schema.TypeString,
										Required: true,
									},
									"node_serial": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipmi_configure_now": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"image_successful": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"ipv6_interface": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cvm_num_vcpus": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ipmi_mac": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"rdma_mac_addr": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"bond_uplinks": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"current_network_interface": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"hypervisor": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"vswitches": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"bond_lacp_rate": {
										Type:     schema.TypeString,
										Required: true,
									},
									"image_now": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"ucsm_managed_mode": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipmi_ip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"current_cvm_vlan_tag": {
										Type:     schema.TypeString,
										Required: true,
									},
									"cvm_ip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"exlude_boot_serial": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"mitigate_low_boot_space": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"block_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"hyperv_product_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"unc_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"install_script": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipmi_user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hypervisor_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"unc_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"xs_master_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"skip_hypervisor": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"hypervisor_gateway": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nos_package": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"ucsm_user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"session_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFoundationImageNodesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}
