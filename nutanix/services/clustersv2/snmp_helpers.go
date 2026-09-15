package clustersv2

import (
	"context"
	"fmt"

	cmgmtConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	cmgmtRequest "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/request/clusters"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/clusters"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// fetchSnmpConfig retrieves the SNMP configuration for a given cluster.
// CreateSnmp* APIs return only a TaskReference, so we resolve newly created
// child entities (users/traps) by querying the parent SNMP config and matching
// on naturally unique attributes.
func fetchSnmpConfig(ctx context.Context, conn *clusters.Client, clusterExtID string) (*cmgmtConfig.SnmpConfig, error) {
	req := cmgmtRequest.GetSnmpConfigByClusterIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
	}
	resp, err := conn.ClusterEntityAPI.GetSnmpConfigByClusterId(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error while fetching SNMP config for cluster (%s): %w", clusterExtID, err)
	}
	cfg, ok := resp.Data.GetValue().(cmgmtConfig.SnmpConfig)
	if !ok {
		return nil, fmt.Errorf("unexpected response data type when fetching SNMP config")
	}
	return &cfg, nil
}

// lookupSnmpUserExtIDByUsername resolves the ext_id of a newly created SNMP user
// by matching against the unique username on the cluster.
func lookupSnmpUserExtIDByUsername(ctx context.Context, conn *clusters.Client, clusterExtID, username string) (string, error) {
	cfg, err := fetchSnmpConfig(ctx, conn, clusterExtID)
	if err != nil {
		return "", err
	}
	for _, u := range cfg.Users {
		if utils.StringValue(u.Username) == username {
			return utils.StringValue(u.ExtId), nil
		}
	}
	return "", fmt.Errorf("SNMP user with username %q was not found on cluster %s", username, clusterExtID)
}

// snmpTrapMatch holds the disambiguating attributes used to resolve a
// freshly-created SNMP trap to its server-assigned ext_id. The upstream
// CreateSnmpTrap API only returns a TaskReference, so after the task
// succeeds we re-read the SNMP config and pick the entry matching these
// fields. (Cluster, version, address.ipv4|ipv6 value, port and protocol
// together are unique in practice.)
type snmpTrapMatch struct {
	Version  string
	IPv4     string
	IPv6     string
	Port     int
	Protocol string
}

// lookupSnmpTrapExtIDByAttrs resolves the ext_id of a newly created SNMP
// trap by matching against version + address + port + protocol on the
// cluster's SNMP config. This is needed because CreateSnmpTrap returns
// only a TaskReference.
func lookupSnmpTrapExtIDByAttrs(ctx context.Context, conn *clusters.Client, clusterExtID string, m snmpTrapMatch) (string, error) {
	cfg, err := fetchSnmpConfig(ctx, conn, clusterExtID)
	if err != nil {
		return "", err
	}
	for _, t := range cfg.Traps {
		if t.Version == nil || t.Version.GetName() != m.Version {
			continue
		}
		if m.Port != 0 && utils.IntValue(t.Port) != m.Port {
			continue
		}
		if m.Protocol != "" {
			if t.Protocol == nil || t.Protocol.GetName() != m.Protocol {
				continue
			}
		}
		if m.IPv4 != "" {
			if t.Address == nil || t.Address.Ipv4 == nil ||
				utils.StringValue(t.Address.Ipv4.Value) != m.IPv4 {
				continue
			}
		}
		if m.IPv6 != "" {
			if t.Address == nil || t.Address.Ipv6 == nil ||
				utils.StringValue(t.Address.Ipv6.Value) != m.IPv6 {
				continue
			}
		}
		return utils.StringValue(t.ExtId), nil
	}
	return "", fmt.Errorf("SNMP trap matching version=%q ipv4=%q ipv6=%q port=%d protocol=%q was not found on cluster %s",
		m.Version, m.IPv4, m.IPv6, m.Port, m.Protocol, clusterExtID)
}
