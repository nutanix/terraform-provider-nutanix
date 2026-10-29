package prismv2

import (
	"testing"

	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitExpandProviderConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandProviderConfig([]interface{}{}))
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandProviderConfig(nil))
	})

	t.Run("nil first element", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandProviderConfig([]interface{}{nil}))
	})

	t.Run("nil credentials block", func(t *testing.T) {
		t.Parallel()
		got := expandProviderConfig([]interface{}{map[string]interface{}{
			"bucket_name": "bucket-1",
			"region":      "us-west-2",
			"credentials": []interface{}{nil},
		}})
		require.NotNil(t, got)
		aws := got.GetValue().(management.AWSS3Config)
		require.Equal(t, "bucket-1", utils.StringValue(aws.BucketName))
		require.Equal(t, "us-west-2", utils.StringValue(aws.Region))
		require.Nil(t, aws.Credentials)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandProviderConfig([]interface{}{map[string]interface{}{
			"bucket_name": "bucket-1",
			"region":      "us-west-2",
			"credentials": []interface{}{map[string]interface{}{
				"access_key_id":     "AKIA",
				"secret_access_key": "secret",
			}},
		}})
		require.NotNil(t, got)
		aws := got.GetValue().(management.AWSS3Config)
		require.Equal(t, "bucket-1", utils.StringValue(aws.BucketName))
		require.Equal(t, "us-west-2", utils.StringValue(aws.Region))
		require.Equal(t, "AKIA", utils.StringValue(aws.Credentials.AccessKeyId))
		require.Equal(t, "secret", utils.StringValue(aws.Credentials.SecretAccessKey))
	})
}

func TestUnitExpandAccessKeyCredentials(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandAccessKeyCredentials([]interface{}{}))
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandAccessKeyCredentials(nil))
	})

	t.Run("nil first element", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandAccessKeyCredentials([]interface{}{nil}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandAccessKeyCredentials([]interface{}{map[string]interface{}{
			"access_key_id":     "AKIA",
			"secret_access_key": "secret",
		}})
		require.Equal(t, "AKIA", utils.StringValue(got.AccessKeyId))
		require.Equal(t, "secret", utils.StringValue(got.SecretAccessKey))
	})
}

func TestUnitExpandBackupPolicy(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandBackupPolicy([]interface{}{}))
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandBackupPolicy(nil))
	})

	t.Run("nil first element", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandBackupPolicy([]interface{}{nil}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandBackupPolicy([]interface{}{map[string]interface{}{
			"rpo_in_minutes": 45,
		}})
		require.Equal(t, 45, utils.IntValue(got.RpoInMinutes))
	})
}

func TestUnitExpandClusterSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want *config.Size
	}{
		{name: "starter", in: "STARTER", want: sizePtr(2)},
		{name: "small", in: "SMALL", want: sizePtr(3)},
		{name: "large", in: "LARGE", want: sizePtr(4)},
		{name: "extralarge", in: "EXTRALARGE", want: sizePtr(5)},
		{name: "unknown", in: "OTHER", want: nil},
		{name: "empty", in: "", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandClusterSize(tt.in))
		})
	}
}

func TestUnitExpandPCConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandPCConfig([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandPCConfig([]interface{}{map[string]interface{}{
			"should_enable_lockdown_mode": true,
			"name":                        "pc-1",
			"size":                        "SMALL",
			"build_info": []interface{}{map[string]interface{}{
				"version": "2024.1",
			}},
			"resource_config": []interface{}{map[string]interface{}{
				"num_vcpus":            4,
				"memory_size_bytes":    1024,
				"data_disk_size_bytes": 2048,
				"container_ext_ids":    []interface{}{"ctr-1"},
			}},
			"credentials": []interface{}{map[string]interface{}{
				"username": "admin",
				"password": "pass",
			}},
			"bootstrap_config": []interface{}{map[string]interface{}{
				"environment_info": []interface{}{map[string]interface{}{
					"type":              "ONPREM",
					"provider_type":     "NTNX",
					"provisioning_type": "NTNX",
				}},
			}},
		}})
		require.NotNil(t, got)
		require.True(t, utils.BoolValue(got.ShouldEnableLockdownMode))
		require.Equal(t, "pc-1", utils.StringValue(got.Name))
		require.Equal(t, config.Size(3), *got.Size)
		require.Equal(t, "2024.1", utils.StringValue(got.BuildInfo.Version))
		require.Equal(t, 4, utils.IntValue(got.ResourceConfig.NumVcpus))
		require.Equal(t, "admin", utils.StringValue(got.Credentials[0].Username))
		require.Equal(t, "ONPREM", flattenEnvironmentType(*got.BootstrapConfig.EnvironmentInfo.Type))
	})
}

func TestUnitExpandBootstrapConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandBootstrapConfig([]interface{}{}))
	})

	t.Run("populated environment", func(t *testing.T) {
		t.Parallel()
		got := expandBootstrapConfig([]interface{}{map[string]interface{}{
			"environment_info": []interface{}{map[string]interface{}{
				"type":              "NTNX_CLOUD",
				"provider_type":     "AWS",
				"provisioning_type": "NATIVE",
			}},
		}})
		require.NotNil(t, got.EnvironmentInfo)
		require.Equal(t, "NTNX_CLOUD", flattenEnvironmentType(*got.EnvironmentInfo.Type))
		require.Equal(t, "AWS", flattenEnvironmentProviderType(*got.EnvironmentInfo.ProviderType))
		require.Equal(t, "NATIVE", flattenEnvironmentProvisioningType(*got.EnvironmentInfo.ProvisioningType))
	})

	t.Run("populated cloud init userdata", func(t *testing.T) {
		t.Parallel()
		got := expandBootstrapConfig([]interface{}{map[string]interface{}{
			"cloud_init_config": []interface{}{map[string]interface{}{
				"datasource_type": "CONFIG_DRIVE_V2",
				"metadata":        "meta",
				"cloud_init_script": []interface{}{map[string]interface{}{
					"user_data": []interface{}{map[string]interface{}{
						"value": "dXNlcg==",
					}},
				}},
			}},
		}})
		require.Len(t, got.CloudInitConfig, 1)
		require.Equal(t, "meta", utils.StringValue(got.CloudInitConfig[0].Metadata))
		require.NotNil(t, got.CloudInitConfig[0].CloudInitScript)
	})
}

func TestUnitExpandCustomKeyValuesPairs(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandCustomKeyValuesPairs(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandCustomKeyValuesPairs([]interface{}{map[string]interface{}{
			"key_value_pairs": []interface{}{map[string]interface{}{
				"name": "k1",
				"value": []interface{}{map[string]interface{}{
					"string": "v1",
				}},
			}},
		}})
		require.NotNil(t, got)
		require.Len(t, got.KeyValuePairs, 1)
		require.Equal(t, "k1", utils.StringValue(got.KeyValuePairs[0].Name))
		require.Equal(t, "v1", got.KeyValuePairs[0].Value.GetValue())
	})
}

func TestUnitExpandKVPair(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandKVPair(map[string]interface{}{})
		require.Nil(t, got.Name)
		require.Nil(t, got.Value)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandKVPair(map[string]interface{}{
			"name": "k1",
			"value": []interface{}{map[string]interface{}{
				"integer": 7,
			}},
		})
		require.Equal(t, "k1", utils.StringValue(got.Name))
		require.Equal(t, 7, got.Value.GetValue())
	})
}

func TestUnitExpandValue(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		got := expandValue(nil)
		require.NotNil(t, got)
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		got := expandValue([]interface{}{map[string]interface{}{"string": "hello"}})
		require.Equal(t, "hello", got.GetValue())
	})

	t.Run("boolean", func(t *testing.T) {
		t.Parallel()
		got := expandValue([]interface{}{map[string]interface{}{"boolean": true}})
		require.Equal(t, true, got.GetValue())
	})

	t.Run("string list", func(t *testing.T) {
		t.Parallel()
		got := expandValue([]interface{}{map[string]interface{}{
			"string_list": []interface{}{"a", "b"},
		}})
		require.Equal(t, []string{"a", "b"}, got.GetValue())
	})

	t.Run("integer list", func(t *testing.T) {
		t.Parallel()
		got := expandValue([]interface{}{map[string]interface{}{
			"integer_list": []interface{}{1, 2},
		}})
		require.Equal(t, []int{1, 2}, got.GetValue())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		got := expandValue([]interface{}{map[string]interface{}{}})
		require.Nil(t, got)
	})
}

func TestUnitExpandCredentials(t *testing.T) {
	t.Parallel()

	t.Run("nil empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandCredentials(nil))
		require.Nil(t, expandCredentials([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandCredentials([]interface{}{map[string]interface{}{
			"username": "admin",
			"password": "secret",
		}})
		require.Len(t, got, 1)
		require.Equal(t, "admin", utils.StringValue(got[0].Username))
		require.Equal(t, "secret", utils.StringValue(got[0].Password))
	})
}

func TestUnitExpandResourceConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandResourceConfig([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandResourceConfig([]interface{}{map[string]interface{}{
			"num_vcpus":            8,
			"memory_size_bytes":    4096,
			"data_disk_size_bytes": 8192,
			"container_ext_ids":    []interface{}{"c1", "c2"},
		}})
		require.Equal(t, 8, utils.IntValue(got.NumVcpus))
		require.Equal(t, int64(4096), utils.Int64Value(got.MemorySizeBytes))
		require.Equal(t, int64(8192), utils.Int64Value(got.DataDiskSizeBytes))
		require.Equal(t, []string{"c1", "c2"}, got.ContainerExtIds)
	})
}

func TestUnitExpandIPv4Address(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv4Address(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Address([]interface{}{map[string]interface{}{
			"value":         "10.0.0.1",
			"prefix_length": 24,
		}})
		require.Equal(t, "10.0.0.1", utils.StringValue(got.Value))
		require.Equal(t, 24, utils.IntValue(got.PrefixLength))
	})
}

func TestUnitExpandIPv6Address(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv6Address(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv6Address([]interface{}{map[string]interface{}{
			"value":         "2001:db8::1",
			"prefix_length": 128,
		}})
		require.Equal(t, "2001:db8::1", utils.StringValue(got.Value))
		require.Equal(t, 128, utils.IntValue(got.PrefixLength))
	})
}

func TestUnitExpandFQDN(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandFQDN(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandFQDN([]interface{}{map[string]interface{}{"value": "pc.example.test"}})
		require.Equal(t, "pc.example.test", utils.StringValue(got.Value))
	})
}

func TestUnitExpandIPAddress(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPAddress(map[string]interface{}{})
		require.NotNil(t, got)
		require.Nil(t, got.Ipv4)
		require.Nil(t, got.Ipv6)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPAddress(map[string]interface{}{
			"ipv4": []interface{}{map[string]interface{}{
				"value":         "10.0.0.5",
				"prefix_length": 32,
			}},
		})
		require.Equal(t, "10.0.0.5", utils.StringValue(got.Ipv4.Value))
	})
}

func TestUnitExpandIPAddressOrFqdn(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPAddressOrFqdn(map[string]interface{}{})
		require.Nil(t, got.Ipv4)
		require.Nil(t, got.Fqdn)
	})

	t.Run("populated fqdn", func(t *testing.T) {
		t.Parallel()
		got := expandIPAddressOrFqdn(map[string]interface{}{
			"fqdn": []interface{}{map[string]interface{}{"value": "dns.example.test"}},
		})
		require.Equal(t, "dns.example.test", utils.StringValue(got.Fqdn.Value))
	})
}

func TestUnitExpandIPRange(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPRange(map[string]interface{}{})
		require.Nil(t, got.Begin)
		require.Nil(t, got.End)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPRange(map[string]interface{}{
			"begin": []interface{}{map[string]interface{}{
				"ipv4": []interface{}{map[string]interface{}{
					"value":         "10.0.0.10",
					"prefix_length": 32,
				}},
			}},
			"end": []interface{}{map[string]interface{}{
				"ipv4": []interface{}{map[string]interface{}{
					"value":         "10.0.0.20",
					"prefix_length": 32,
				}},
			}},
		})
		require.Equal(t, "10.0.0.10", utils.StringValue(got.Begin.Ipv4.Value))
		require.Equal(t, "10.0.0.20", utils.StringValue(got.End.Ipv4.Value))
	})
}

func TestUnitExpandExternalNetwork(t *testing.T) {
	t.Parallel()

	got := expandExternalNetwork(map[string]interface{}{
		"network_ext_id": "net-1",
		"default_gateway": []interface{}{map[string]interface{}{
			"ipv4": []interface{}{map[string]interface{}{
				"value":         "10.0.0.1",
				"prefix_length": 32,
			}},
		}},
		"subnet_mask": []interface{}{map[string]interface{}{
			"ipv4": []interface{}{map[string]interface{}{
				"value":         "255.255.255.0",
				"prefix_length": 32,
			}},
		}},
		"ip_ranges": []interface{}{},
	})
	require.Equal(t, "net-1", utils.StringValue(got.NetworkExtId))
	require.Equal(t, "10.0.0.1", utils.StringValue(got.DefaultGateway.Ipv4.Value))
}

func TestUnitExpandInternalNetwork(t *testing.T) {
	t.Parallel()

	got := expandInternalNetwork(map[string]interface{}{
		"default_gateway": []interface{}{map[string]interface{}{
			"ipv4": []interface{}{map[string]interface{}{
				"value":         "10.0.0.1",
				"prefix_length": 32,
			}},
		}},
		"subnet_mask": []interface{}{map[string]interface{}{
			"ipv4": []interface{}{map[string]interface{}{
				"value":         "255.255.255.0",
				"prefix_length": 32,
			}},
		}},
		"ip_ranges": []interface{}{},
	})
	require.Equal(t, "10.0.0.1", utils.StringValue(got.DefaultGateway.Ipv4.Value))
	require.Equal(t, "255.255.255.0", utils.StringValue(got.SubnetMask.Ipv4.Value))
}

func TestUnitExpandPCNetwork(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandPCNetwork([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandPCNetwork([]interface{}{map[string]interface{}{
			"external_address": []interface{}{map[string]interface{}{
				"ipv4": []interface{}{map[string]interface{}{
					"value":         "10.0.0.50",
					"prefix_length": 32,
				}},
			}},
			"name_servers": []interface{}{map[string]interface{}{
				"ipv4": []interface{}{map[string]interface{}{
					"value":         "8.8.8.8",
					"prefix_length": 32,
				}},
			}},
			"ntp_servers": []interface{}{map[string]interface{}{
				"fqdn": []interface{}{map[string]interface{}{"value": "ntp.example.test"}},
			}},
			"external_networks": []interface{}{map[string]interface{}{
				"network_ext_id": "net-1",
				"default_gateway": []interface{}{map[string]interface{}{
					"ipv4": []interface{}{map[string]interface{}{
						"value":         "10.0.0.1",
						"prefix_length": 32,
					}},
				}},
				"subnet_mask": []interface{}{map[string]interface{}{
					"ipv4": []interface{}{map[string]interface{}{
						"value":         "255.255.255.0",
						"prefix_length": 32,
					}},
				}},
				"ip_ranges": []interface{}{},
			}},
		}})
		require.NotNil(t, got)
		require.Equal(t, "10.0.0.50", utils.StringValue(got.ExternalAddress.Ipv4.Value))
		require.Len(t, got.NameServers, 1)
		require.Len(t, got.NtpServers, 1)
		require.Len(t, got.ExternalNetworks, 1)
	})
}

func TestUnitExpandDomainManagerCloudType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want *management.DomainManagerCloudType
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: "", want: nil},
		{name: "unknown", in: "OTHER", want: nil},
		{name: "onprem", in: "ONPREM_CLOUD", want: cloudTypePtr(2)},
		{name: "hosted", in: "NUTANIX_HOSTED_CLOUD", want: cloudTypePtr(3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandDomainManagerCloudType(tt.in))
		})
	}
}

func TestUnitExpandAuthentication(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandAuthentication(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandAuthentication([]interface{}{map[string]interface{}{
			"username": "admin",
			"password": "secret",
		}})
		require.Equal(t, "admin", utils.StringValue(got.Username))
		require.Equal(t, "secret", utils.StringValue(got.Password))
	})
}

func TestUnitExpandRegisterCredentials(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandRegisterCredentials(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandRegisterCredentials([]interface{}{map[string]interface{}{
			"authentication": []interface{}{map[string]interface{}{
				"username": "admin",
				"password": "secret",
			}},
		}})
		require.NotNil(t, got.Authentication)
		require.Equal(t, "admin", utils.StringValue(got.Authentication.Username))
	})
}

func TestUnitExpandDomainManagerRemoteCluster(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandDomainManagerRemoteCluster(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandDomainManagerRemoteCluster(map[string]interface{}{})
		require.NotNil(t, got)
		require.Nil(t, got.Address)
		require.Nil(t, got.Credentials)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandDomainManagerRemoteCluster(map[string]interface{}{
			"address": []interface{}{map[string]interface{}{
				"ipv4": []interface{}{map[string]interface{}{
					"value":         "10.0.0.9",
					"prefix_length": 32,
				}},
			}},
			"credentials": []interface{}{map[string]interface{}{
				"authentication": []interface{}{map[string]interface{}{
					"username": "admin",
					"password": "secret",
				}},
			}},
		})
		require.Equal(t, "10.0.0.9", utils.StringValue(got.Address.Ipv4.Value))
		require.Equal(t, "admin", utils.StringValue(got.Credentials.Authentication.Username))
	})
}

func sizePtr(v int) *config.Size {
	s := config.Size(v)
	return &s
}

func cloudTypePtr(v int) *management.DomainManagerCloudType {
	c := management.DomainManagerCloudType(v)
	return &c
}
