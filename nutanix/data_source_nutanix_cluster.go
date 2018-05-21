package nutanix

import (
	"fmt"
	"strconv"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixClusterRead,

		Schema: getDataSourceClusterSchema(),
	}
}

func dataSourceNutanixClusterRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*NutanixClient).API

	c, ok := d.GetOk("cluster_id")

	if !ok {
		return fmt.Errorf("Please provide the cluster_id attribute")
	}

	// Make request to the API
	v, err := conn.V3.GetCluster(c.(string))
	if err != nil {
		return err
	}

	// set metadata values
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = utils.TimeValue(v.Metadata.LastUpdateTime).String()
	metadata["kind"] = utils.StringValue(v.Metadata.Kind)
	metadata["uuid"] = utils.StringValue(v.Metadata.UUID)
	metadata["creation_time"] = utils.TimeValue(v.Metadata.CreationTime).String()
	metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(v.Metadata.SpecVersion)))
	metadata["spec_hash"] = utils.StringValue(v.Metadata.SpecHash)
	metadata["name"] = utils.StringValue(v.Metadata.Name)
	if err := d.Set("metadata", metadata); err != nil {
		return err
	}
	if err := d.Set("categories", v.Metadata.Categories); err != nil {
		return err
	}
	if err := d.Set("api_version", utils.StringValue(v.APIVersion)); err != nil {
		return err
	}

	pr := make(map[string]interface{})
	if v.Metadata.ProjectReference != nil {
		pr["kind"] = utils.StringValue(v.Metadata.ProjectReference.Kind)
		pr["name"] = utils.StringValue(v.Metadata.ProjectReference.Name)
		pr["uuid"] = utils.StringValue(v.Metadata.ProjectReference.UUID)
	}
	if err := d.Set("project_reference", pr); err != nil {
		return err
	}

	or := make(map[string]interface{})
	if v.Metadata.OwnerReference != nil {
		or["kind"] = utils.StringValue(v.Metadata.OwnerReference.Kind)
		or["name"] = utils.StringValue(v.Metadata.OwnerReference.Name)
		or["uuid"] = utils.StringValue(v.Metadata.OwnerReference.UUID)
	}
	if err := d.Set("owner_reference", or); err != nil {
		return err
	}
	if err := d.Set("name", utils.StringValue(v.Status.Name)); err != nil {
		return err
	}
	if err := d.Set("state", utils.StringValue(v.Status.State)); err != nil {
		return err
	}

	nodes := make([]map[string]interface{}, 0)
	if v.Status.Resources.Nodes != nil {
		if v.Status.Resources.Nodes.HypervisorServerList != nil {
			nodes = make([]map[string]interface{}, len(v.Status.Resources.Nodes.HypervisorServerList))
			for k, v := range v.Status.Resources.Nodes.HypervisorServerList {
				node := make(map[string]interface{})
				node["ip"] = utils.StringValue(v.IP)
				node["version"] = utils.StringValue(v.Version)
				node["type"] = utils.StringValue(v.Type)
				nodes[k] = node
			}
		}
	}
	if err := d.Set("nodes", nodes); err != nil {
		return err
	}

	config := v.Status.Resources.Config
	if err := d.Set("gpu_driver_version", utils.StringValue(config.GpuDriverVersion)); err != nil {
		return err
	}

	clientAuth := make(map[string]interface{})
	if config.ClientAuth != nil {
		clientAuth["status"] = utils.StringValue(config.ClientAuth.Status)
		clientAuth["ca_chain"] = utils.StringValue(config.ClientAuth.CaChain)
		clientAuth["name"] = utils.StringValue(config.ClientAuth.Name)
	}
	if err := d.Set("client_auth", clientAuth); err != nil {
		return err
	}

	authPublicKey := make([]map[string]interface{}, 0)
	if config.AuthorizedPublicKeyList != nil {
		authPublicKey := make([]map[string]interface{}, len(config.AuthorizedPublicKeyList))
		for k, v := range config.AuthorizedPublicKeyList {
			auth := make(map[string]interface{})
			auth["key"] = utils.StringValue(v.Key)
			auth["name"] = utils.StringValue(v.Name)
			authPublicKey[k] = auth
		}
	}
	if err := d.Set("authorized_public_key_list", authPublicKey); err != nil {
		return err
	}

	ncc := make(map[string]interface{})
	nos := make(map[string]interface{})
	if config.SoftwareMap != nil {
		ncc["software_type"] = utils.StringValue(config.SoftwareMap.NCC.SoftwareType)
		ncc["status"] = utils.StringValue(config.SoftwareMap.NCC.Status)
		ncc["version"] = utils.StringValue(config.SoftwareMap.NCC.Version)
		nos["software_type"] = utils.StringValue(config.SoftwareMap.NOS.SoftwareType)
		nos["status"] = utils.StringValue(config.SoftwareMap.NOS.Status)
		nos["version"] = utils.StringValue(config.SoftwareMap.NOS.Version)
	}
	if err := d.Set("software_map_ncc", ncc); err != nil {
		return err
	}
	if err := d.Set("software_map_nos", nos); err != nil {
		return err
	}
	if err := d.Set("encryption_status", utils.StringValue(config.EncryptionStatus)); err != nil {
		return err
	}

	signingInfo := make(map[string]interface{})
	if config.SslKey != nil {

		if err := d.Set("ssl_key_type", utils.StringValue(config.SslKey.KeyType)); err != nil {
			return err
		}
		if err := d.Set("ssl_key_name", utils.StringValue(config.SslKey.KeyName)); err != nil {
			return err
		}

		if config.SslKey.SigningInfo != nil {
			signingInfo["city"] = utils.StringValue(config.SslKey.SigningInfo.City)
			signingInfo["common_name_suffix"] = utils.StringValue(config.SslKey.SigningInfo.CommonNameSuffix)
			signingInfo["state"] = utils.StringValue(config.SslKey.SigningInfo.State)
			signingInfo["country_code"] = utils.StringValue(config.SslKey.SigningInfo.CountryCode)
			signingInfo["common_name"] = utils.StringValue(config.SslKey.SigningInfo.CommonName)
			signingInfo["organization"] = utils.StringValue(config.SslKey.SigningInfo.Organization)
			signingInfo["email_address"] = utils.StringValue(config.SslKey.SigningInfo.EmailAddress)
		}
		if err := d.Set("ssl_key_signing_info", signingInfo); err != nil {
			return err
		}
		if err := d.Set("ssl_key_expire_datetime", utils.StringValue(config.SslKey.ExpireDatetime)); err != nil {
			return err
		}

	} else {
		if err := d.Set("ssl_key_type", ""); err != nil {
			return err
		}
		if err := d.Set("ssl_key_name", ""); err != nil {
			return err
		}
		if err := d.Set("ssl_key_signing_info", signingInfo); err != nil {
			return err
		}
		if err := d.Set("ssl_key_expire_datetime", ""); err != nil {
			return err
		}

	}

	if err := d.Set("service_list", utils.StringValueSlice(config.ServiceList)); err != nil {
		return err
	}

	if err := d.Set("supported_information_verbosity", utils.StringValue(config.SupportedInformationVerbosity)); err != nil {
		return err
	}

	certSigning := make(map[string]interface{})
	if config.CertificationSigningInfo != nil {
		certSigning["city"] = utils.StringValue(config.CertificationSigningInfo.City)
		certSigning["common_name_suffix"] = utils.StringValue(config.CertificationSigningInfo.CommonNameSuffix)
		certSigning["state"] = utils.StringValue(config.CertificationSigningInfo.State)
		certSigning["country_code"] = utils.StringValue(config.CertificationSigningInfo.CountryCode)
		certSigning["common_name"] = utils.StringValue(config.CertificationSigningInfo.CommonName)
		certSigning["organization"] = utils.StringValue(config.CertificationSigningInfo.Organization)
		certSigning["email_address"] = utils.StringValue(config.CertificationSigningInfo.EmailAddress)
	}
	if err := d.Set("certification_signing_info", certSigning); err != nil {
		return err
	}
	if err := d.Set("operation_mode", utils.StringValue(config.OperationMode)); err != nil {
		return err
	}

	caCert := make([]map[string]interface{}, 0)
	if config.CaCertificateList != nil {
		caCert = make([]map[string]interface{}, len(config.CaCertificateList))
		for k, v := range config.CaCertificateList {
			ca := make(map[string]interface{})
			ca["ca_name"] = utils.StringValue(v.CaName)
			ca["certificate"] = utils.StringValue(v.Certificate)
			caCert[k] = ca
		}
	}
	if err := d.Set("ca_certificate_list", caCert); err != nil {
		return err
	}
	if err := d.Set("enabled_feature_list", utils.StringValueSlice(config.EnabledFeatureList)); err != nil {
		return err
	}
	if err := d.Set("is_available", utils.BoolValue(config.IsAvailable)); err != nil {
		return err
	}

	build := make(map[string]interface{})
	if config.Build != nil {
		build["commit_id"] = utils.StringValue(config.Build.CommitID)
		build["full_version"] = utils.StringValue(config.Build.FullVersion)
		build["commit_date"] = utils.StringValue(config.Build.CommitDate)
		build["version"] = utils.StringValue(config.Build.Version)
		build["short_commit_id"] = utils.StringValue(config.Build.ShortCommitID)
		build["build_type"] = utils.StringValue(config.Build.BuildType)
	}
	if err := d.Set("build", build); err != nil {
		return err
	}
	if err := d.Set("timezone", utils.StringValue(config.Timezone)); err != nil {
		return err
	}
	if err := d.Set("cluster_arch", utils.StringValue(config.ClusterArch)); err != nil {
		return err
	}

	managementServer := make([]map[string]interface{}, 0)
	if config.ManagementServerList != nil {
		managementServer = make([]map[string]interface{}, len(config.ManagementServerList))
		for k, v := range config.ManagementServerList {
			manage := make(map[string]interface{})
			manage["ip"] = utils.StringValue(v.IP)
			manage["drs_enabled"] = utils.BoolValue(v.DrsEnabled)
			manage["status_list"] = utils.StringValueSlice(v.StatusList)
			manage["type"] = utils.StringValue(v.Type)
			managementServer[k] = manage
		}
	}
	if err := d.Set("management_server_list", managementServer); err != nil {
		return err
	}

	network := v.Status.Resources.Network
	if err := d.Set("masquerading_port", utils.Int64Value(network.MasqueradingPort)); err != nil {
		return err
	}
	if err := d.Set("masquerading_ip", utils.StringValue(network.MasqueradingIP)); err != nil {
		return err
	}
	if err := d.Set("external_ip", utils.StringValue(network.ExternalIP)); err != nil {
		return err
	}

	httpProxy := make([]map[string]interface{}, 0)
	if network.HTTPProxyList != nil {
		httpProxy = make([]map[string]interface{}, len(network.HTTPProxyList))
		for k, v := range network.HTTPProxyList {
			http := make(map[string]interface{})
			creds := make(map[string]interface{})
			addr := make(map[string]interface{})
			creds["username"] = utils.StringValue(v.Credentials.Username)
			creds["password"] = utils.StringValue(v.Credentials.Password)
			http["credentials"] = creds
			http["proxy_type_list"] = utils.StringValueSlice(v.ProxyTypeList)
			addr["ip"] = utils.StringValue(v.Address.IP)
			addr["fqdn"] = utils.StringValue(v.Address.FQDN)
			addr["port"] = strconv.Itoa(int(utils.Int64Value(v.Address.Port)))
			addr["ipv6"] = utils.StringValue(v.Address.IPV6)
			http["address"] = addr

			httpProxy[k] = http
		}
	}
	if err := d.Set("http_proxy_list", httpProxy); err != nil {
		return err
	}

	smtpServCreds := make(map[string]interface{})
	smtpServAddr := make(map[string]interface{})
	if network.SMTPServer != nil {
		if err := d.Set("smtp_server_type", utils.StringValue(network.SMTPServer.Type)); err != nil {
			return err
		}
		if err := d.Set("smtp_server_email_address", utils.StringValue(network.SMTPServer.EmailAddress)); err != nil {
			return err
		}

		if network.SMTPServer != nil {
			smtpServCreds["username"] = utils.StringValue(network.SMTPServer.Server.Credentials.Username)
			smtpServCreds["password"] = utils.StringValue(network.SMTPServer.Server.Credentials.Password)

			smtpServAddr["ip"] = utils.StringValue(network.SMTPServer.Server.Address.IP)
			smtpServAddr["fqdn"] = utils.StringValue(network.SMTPServer.Server.Address.FQDN)
			smtpServAddr["port"] = strconv.Itoa(int(utils.Int64Value(network.SMTPServer.Server.Address.Port)))
			smtpServAddr["ipv6"] = utils.StringValue(network.SMTPServer.Server.Address.IPV6)
		}
		if err := d.Set("smtp_server_credentials", smtpServCreds); err != nil {
			return err
		}
		if err := d.Set("smtp_server_proxy_type_list", utils.StringValueSlice(network.SMTPServer.Server.ProxyTypeList)); err != nil {
			return err
		}
		if err := d.Set("smtp_server_address", smtpServAddr); err != nil {
			return err
		}

	} else {
		if err := d.Set("smtp_server_type", ""); err != nil {
			return err
		}
		if err := d.Set("smtp_server_email_address", ""); err != nil {
			return err
		}
		if err := d.Set("smtp_server_credentials", smtpServCreds); err != nil {
			return err
		}
		if err := d.Set("smtp_server_proxy_type_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("smtp_server_address", smtpServAddr); err != nil {
			return err
		}

	}

	if err := d.Set("ntp_server_ip_list", utils.StringValueSlice(network.NameServerIPList)); err != nil {
		return err
	}
	if err := d.Set("external_subnet", utils.StringValue(network.ExternalSubnet)); err != nil {
		return err
	}
	if err := d.Set("external_data_services_ip", utils.StringValue(network.ExternalDataServicesIP)); err != nil {
		return err
	}
	if err := d.Set("internal_subnet", utils.StringValue(network.InternalSubnet)); err != nil {
		return err
	}

	domain := network.DomainServer
	domServCreds := make(map[string]interface{})
	if domain != nil {
		if err := d.Set("domain_server_nameserver", utils.StringValue(domain.Nameserver)); err != nil {
			return err
		}
		if err := d.Set("domain_server_name", utils.StringValue(domain.Name)); err != nil {
			return err
		}

		domServCreds["username"] = utils.StringValue(domain.DomainCredentials.Username)
		domServCreds["password"] = utils.StringValue(domain.DomainCredentials.Password)
		if err := d.Set("domain_server_credentials", domServCreds); err != nil {
			return err
		}
	} else {
		if err := d.Set("domain_server_nameserver", ""); err != nil {
			return err
		}
		if err := d.Set("domain_server_name", ""); err != nil {
			return err
		}
		if err := d.Set("domain_server_credentials", domServCreds); err != nil {
			return err
		}
	}

	if err := d.Set("nfs_subnet_whitelist", utils.StringValueSlice(network.NFSSubnetWhitelist)); err != nil {
		return err
	}
	if err := d.Set("name_server_ip_list", utils.StringValueSlice(network.NameServerIPList)); err != nil {
		return err
	}

	httpWhiteList := make([]map[string]interface{}, 0)
	if network.HTTPProxyWhitelist != nil {
		httpWhiteList = make([]map[string]interface{}, len(network.HTTPProxyWhitelist))
		for k, v := range network.HTTPProxyWhitelist {
			http := make(map[string]interface{})
			http["target"] = utils.StringValue(v.Target)
			http["target_type"] = utils.StringValue(v.TargetType)
			httpWhiteList[k] = http
		}
	}
	if err := d.Set("http_proxy_whitelist", httpWhiteList); err != nil {
		return err
	}

	analysis := make(map[string]interface{})
	if v.Status.Resources.Analysis != nil {
		analysis["bully_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.BullyVMNum)
		analysis["constrained_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.ConstrainedVMNum)
		analysis["dead_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.DeadVMNum)
		analysis["inefficient_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.InefficientVMNum)
		analysis["overprovisioned_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.OverprovisionedVMNum)
	}
	if err := d.Set("analysis_vm_efficiency_map", analysis); err != nil {
		return err
	}

	d.SetId(*v.Metadata.UUID)

	return nil
}

func getDataSourceClusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
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
			Type:     schema.TypeList,
			Computed: true,
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
		"software_map_ncc": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
		},
		"software_map_nos": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
		},
		"encryption_status": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"ssl_key_type": &schema.Schema{
			Type:     schema.TypeString,
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
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
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
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bully_vm_num": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"constrained_vm_num": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"dead_vm_num": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"inefficient_vm_num": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"overprovisioned_vm_num": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}
