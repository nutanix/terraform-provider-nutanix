package clustersv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixClusterProfileV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixClusterProfileV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allowed_overrides": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_server_ip_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      common.HashIPItem,
				Elem:     common.SchemaForIPList(false), // do not include FQDN
			},
			"ntp_server_ip_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      common.HashIPItem,
				Elem:     common.SchemaForIPList(true), // include FQDN
			},
			"smtp_server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     common.SchemaForIPList(true), // include FQDN
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"password": {
										Type:      schema.TypeString,
										Sensitive: true,
										Computed:  true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"nfs_subnet_white_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"snmp_config": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"users": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auth_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auth_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"priv_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"priv_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"transports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"traps": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     common.SchemaForIPList(false), // do not include FQDN
									},
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"should_inform": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"engine_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"receiver_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"community_string": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"rsyslog_server_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     common.SchemaForIPList(false), // do not include FQDN
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"log_severity_level": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"should_log_monitor_files": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"pulse_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"pii_scrubbing_level": {
							Type:         schema.TypeString,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(PIIScrubbingLevelStrings, false),
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixClusterProfileV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	// Fetch the Cluster Profile by UUID
	clusterProfileResp, err := conn.ClusterProfilesAPI.GetClusterProfileById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if Data is nil or empty
	if clusterProfileResp.Data == nil || clusterProfileResp.Data.GetValue() == nil {
		return diag.Errorf("ClusterProfile API returned empty data for ID %s", d.Id())
	}

	// Safe type assertion
	clusterProfile, ok := clusterProfileResp.Data.GetValue().(config.ClusterProfile)
	if !ok {
		return diag.Errorf("ClusterProfile API returned unexpected type for ID %s", d.Id())
	}

	// Set the resource data from the API response
	if err := d.Set("tenant_id", clusterProfile.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(clusterProfile.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", clusterProfile.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", clusterProfile.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", clusterProfile.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allowed_overrides", common.FlattenEnumValueList(clusterProfile.AllowedOverrides)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name_server_ip_list", flattenIPAddressList(clusterProfile.NameServerIpList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ntp_server_ip_list", flattenIPAddressOrFQDN(clusterProfile.NtpServerIpList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("smtp_server", flattenSMTPServerRef(clusterProfile.SmtpServer)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nfs_subnet_white_list", clusterProfile.NfsSubnetWhitelist); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("snmp_config", flattenSnmpConfig(clusterProfile.SnmpConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rsyslog_server_list", flattenRsyslogServerList(clusterProfile.RsyslogServerList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("pulse_status", flattenPulseStatus(clusterProfile.PulseStatus)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(clusterProfile.ExtId))
	return nil
}
