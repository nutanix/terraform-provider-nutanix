// Package securityv2 provides resources for managing security-related configurations in Nutanix.
package securityv2

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commonCfg "github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/security/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixKeyManagementServerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixKeyManagementServerV2Read,
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
			"access_information": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"azure_key_vault": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tenant_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"truncated_client_secret": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"credential_expiry_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"kmip_key_vault": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_pem": {
										Type:      schema.TypeString,
										Computed:  true,
										Sensitive: true,
									},
									"private_key": {
										Type:      schema.TypeString,
										Computed:  true,
										Sensitive: true,
									},
									"ca_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ca_pem": {
										Type:      schema.TypeString,
										Computed:  true,
										Sensitive: true,
									},
									"endpoint_url": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_address": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     common.SchemaForIPList(true),
												},
												"port": {
													Type:     schema.TypeInt,
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
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixKeyManagementServerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.KeyManagementServersAPIInstance.GetKeyManagementServerById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching key management server : %v", err)
	}

	getRespValue, ok := resp.Data.GetValue().(config.KeyManagementServer)
	if !ok {
		return diag.Errorf("error: unexpected response type from get API, expected KeyManagementServer")
	}
	getResp := getRespValue

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	accessInfo, flattenErr := flattenAccessInformation(getResp.GetAccessInformation())
	if flattenErr != nil {
		return diag.FromErr(flattenErr)
	}
	if err := d.Set("access_information", accessInfo); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creation_timestamp", utils.TimeStringValue(getResp.CreationTimestamp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenLinks(links []response.ApiLink) []interface{} {
	if len(links) > 0 {
		flattenedLinks := make([]interface{}, len(links))

		for k, v := range links {
			link := make(map[string]interface{})

			if v.Href != nil {
				link["href"] = v.Href
			}
			if v.Rel != nil {
				link["rel"] = v.Rel
			}
			flattenedLinks[k] = link
		}
		return flattenedLinks
	}
	return nil
}

func flattenAccessInformation(accessInfo interface{}) ([]map[string]interface{}, error) {
	if accessInfo == nil {
		log.Printf("[DEBUG] flattenAccessInformation: accessInfo is nil")
		return nil, fmt.Errorf("access information is nil")
	}

	log.Printf("[DEBUG] flattenAccessInformation: input type=%T", accessInfo)

	switch v := accessInfo.(type) {
	case *config.OneOfKeyManagementServerAccessInformation:
		log.Printf("[DEBUG] flattenAccessInformation: handling OneOfKeyManagementServerAccessInformation")
		return flattenAccessInformation(v.GetValue())
	case *config.AzureAccessInformation:
		if v == nil {
			log.Printf("[DEBUG] flattenAccessInformation: *AzureAccessInformation is nil")
			return nil, fmt.Errorf("access information is nil")
		}
		log.Printf("[DEBUG] flattenAccessInformation: handling *AzureAccessInformation")
		return flattenAccessInformation(*v)
	case config.AzureAccessInformation:
		// NOTE: Do not log sensitive fields (e.g. client_secret).
		log.Printf("[DEBUG] flattenAccessInformation: handling AzureAccessInformation")
		azure := map[string]interface{}{
			"endpoint_url":            utils.StringValue(v.EndpointUrl),
			"key_id":                  utils.StringValue(v.KeyId),
			"tenant_id":               utils.StringValue(v.TenantId),
			"client_id":               utils.StringValue(v.ClientId),
			"truncated_client_secret": utils.StringValue(v.TruncatedClientSecret),
			"credential_expiry_date":  utils.TimeValue(v.CredentialExpiryDate).Format("2006-01-02"),
		}
		return []map[string]interface{}{{
			"azure_key_vault": []map[string]interface{}{azure},
		}}, nil
	case *config.KmipAccessInformation:
		if v == nil {
			log.Printf("[DEBUG] flattenAccessInformation: *KmipAccessInformation is nil")
			return nil, fmt.Errorf("access information is nil")
		}
		log.Printf("[DEBUG] flattenAccessInformation: handling *KmipAccessInformation (endpoints=%d)", len(v.Endpoints))
		return flattenAccessInformation(*v)
	case config.KmipAccessInformation:
		// NOTE: Do not log sensitive fields (e.g. cert_pem/private_key/ca_pem).
		log.Printf("[DEBUG] flattenAccessInformation: handling KmipAccessInformation (endpoints=%d)", len(v.Endpoints))
		kmip := map[string]interface{}{
			"ca_name":     utils.StringValue(v.CaName),
			"ca_pem":      utils.StringValue(v.CaPem),
			"cert_pem":    utils.StringValue(v.CertPem),
			"private_key": utils.StringValue(v.PrivateKey),
			"endpoint_url": func() []interface{} {
				result := make([]interface{}, 0, len(v.Endpoints))
				for _, e := range v.Endpoints {
					endpoint := map[string]interface{}{
						"ip_address": flattenIPAddressOrFQDN(e.IpAddress),
						"port":       utils.IntValue(e.Port),
					}
					result = append(result, endpoint)
				}
				return result
			}(),
		}
		return []map[string]interface{}{{
			"kmip_key_vault": []map[string]interface{}{kmip},
		}}, nil
	default:
		log.Printf("[DEBUG] flattenAccessInformation: unsupported type=%T", accessInfo)
		return nil, fmt.Errorf("unsupported access information type %T", accessInfo)
	}
}

func flattenIPAddressOrFQDN(addr *commonCfg.IPAddressOrFQDN) []map[string]interface{} {
	if addr == nil {
		return nil
	}
	return []map[string]interface{}{{
		"ipv4": flattenIPv4Address(addr.Ipv4),
		"ipv6": flattenIPv6Address(addr.Ipv6),
		"fqdn": flattenFQDN(addr.Fqdn),
	}}
}

func flattenIPv4Address(addr *commonCfg.IPv4Address) []map[string]interface{} {
	if addr == nil {
		return nil
	}
	return []map[string]interface{}{{
		"value":         utils.StringValue(addr.Value),
		"prefix_length": utils.IntValue(addr.PrefixLength),
	}}
}

func flattenIPv6Address(addr *commonCfg.IPv6Address) []map[string]interface{} {
	if addr == nil {
		return nil
	}
	return []map[string]interface{}{{
		"value":         utils.StringValue(addr.Value),
		"prefix_length": utils.IntValue(addr.PrefixLength),
	}}
}

func flattenFQDN(f *commonCfg.FQDN) []map[string]interface{} {
	if f == nil {
		return nil
	}
	return []map[string]interface{}{{
		"value": utils.StringValue(f.Value),
	}}
}
