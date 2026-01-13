package networkingv2

import (
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
)

var (
	networkFunctionFailureHandlingMap = common.EnumToMap([]import1.FailureHandling{
		import1.FAILUREHANDLING_NO_ACTION,
		import1.FAILUREHANDLING_FAIL_CLOSE,
		import1.FAILUREHANDLING_FAIL_OPEN,
	})

	networkFunctionHighAvailabilityModeMap = common.EnumToMap([]import1.HighAvailabilityMode{
		import1.HIGHAVAILABILITYMODE_ACTIVE_PASSIVE,
	})

	networkFunctionTrafficForwardingModeMap = common.EnumToMap([]import1.TrafficForwardingMode{
		import1.TRAFFICFORWARDINGMODE_INLINE,
		import1.TRAFFICFORWARDINGMODE_VTAP,
	})
)

func flattenDataPlaneHealthCheckConfig(cfg *import1.DataPlaneHealthCheckConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}
	m := make(map[string]interface{})
	m["failure_threshold"] = cfg.FailureThreshold
	m["interval_secs"] = cfg.IntervalSecs
	m["success_threshold"] = cfg.SuccessThreshold
	m["timeout_secs"] = cfg.TimeoutSecs
	return []map[string]interface{}{m}
}

func expandDataPlaneHealthCheckConfig(val interface{}) *import1.DataPlaneHealthCheckConfig {
	if val == nil {
		return nil
	}
	l := val.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	cfg := import1.DataPlaneHealthCheckConfig{}

	if v, ok := m["failure_threshold"]; ok {
		if i, ok := v.(int); ok {
			cfg.FailureThreshold = &i
		}
	}
	if v, ok := m["interval_secs"]; ok {
		if i, ok := v.(int); ok {
			cfg.IntervalSecs = &i
		}
	}
	if v, ok := m["success_threshold"]; ok {
		if i, ok := v.(int); ok {
			cfg.SuccessThreshold = &i
		}
	}
	if v, ok := m["timeout_secs"]; ok {
		if i, ok := v.(int); ok {
			cfg.TimeoutSecs = &i
		}
	}

	return &cfg
}

func flattenNicPairs(pairs []import1.NicPair) []interface{} {
	if len(pairs) == 0 {
		return nil
	}
	out := make([]interface{}, len(pairs))
	for i, p := range pairs {
		m := make(map[string]interface{})
		m["ingress_nic_reference"] = p.IngressNicReference
		m["egress_nic_reference"] = p.EgressNicReference
		m["is_enabled"] = p.IsEnabled
		m["vm_reference"] = p.VmReference
		m["data_plane_health_status"] = common.FlattenPtrEnum(p.DataPlaneHealthStatus)
		m["high_availability_state"] = common.FlattenPtrEnum(p.HighAvailabilityState)
		out[i] = m
	}
	return out
}

func expandNicPairs(val interface{}) []import1.NicPair {
	if val == nil {
		return nil
	}
	l := val.([]interface{})
	if len(l) == 0 {
		return nil
	}

	out := make([]import1.NicPair, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		m := raw.(map[string]interface{})
		p := import1.NicPair{}

		if v, ok := m["ingress_nic_reference"]; ok {
			if s, ok := v.(string); ok && s != "" {
				p.IngressNicReference = &s
			}
		}
		if v, ok := m["egress_nic_reference"]; ok {
			if s, ok := v.(string); ok && s != "" {
				p.EgressNicReference = &s
			}
		}
		if v, ok := m["vm_reference"]; ok {
			if s, ok := v.(string); ok && s != "" {
				p.VmReference = &s
			}
		}
		if v, ok := m["is_enabled"]; ok {
			if b, ok := v.(bool); ok {
				p.IsEnabled = &b
			}
		}

		out = append(out, p)
	}
	return out
}
