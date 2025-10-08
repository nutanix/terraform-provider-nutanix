package clustersv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var allowedOverridesEnum = []string{"NFS_SUBNET_WHITELIST_CONFIG", "NTP_SERVER_CONFIG", "SNMP_SERVER_CONFIG", "SMTP_SERVER_CONFIG", "PULSE_CONFIG", "NAME_SERVER_CONFIG", "RSYSLOG_SERVER_CONFIG"}

func ResourceNutanixClusterProfilesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixClusterProfilesV2Create,
		ReadContext:   ResourceNutanixClusterProfilesV2Read,
		UpdateContext: ResourceNutanixClusterProfilesV2Update,
		DeleteContext: ResourceNutanixClusterProfilesV2Delete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"allowed_overrides": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(allowedOverridesEnum, false),
				},
			},
			"name_server_ip_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     SchemaForIPList(false),
			},
			"ntpServerIpList": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     SchemaForIPList(false),
			},
			// Computed fields
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixClusterProfilesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixClusterProfilesV2Read(ctx, d, meta)
}

func ResourceNutanixClusterProfilesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixClusterProfilesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixClusterProfilesV2Read(ctx, d, meta)
}

func ResourceNutanixClusterProfilesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
