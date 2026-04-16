package networkingv2_test

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func checkAttributeLength(resourceName, attribute string, minLength int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		attrKey := fmt.Sprintf("%s.#", attribute)
		attr, ok := rs.Primary.Attributes[attrKey]
		if !ok {
			return fmt.Errorf("attribute %s not found", attrKey)
		}

		count, err := strconv.Atoi(attr)
		if err != nil {
			return fmt.Errorf("error converting %s to int: %s", attrKey, err)
		}

		if count < minLength {
			return fmt.Errorf("expected %s to be >= %d, got %d", attrKey, minLength, count)
		}

		return nil
	}
}

func checkAttributeLengthEqual(resourceName, attribute string, expectedLength int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		attrKey := fmt.Sprintf("%s.#", attribute)
		attr, ok := rs.Primary.Attributes[attrKey]
		if !ok {
			return fmt.Errorf("attribute %s not found", attrKey)
		}

		count, err := strconv.Atoi(attr)
		if err != nil {
			return fmt.Errorf("error converting %s to int: %s", attrKey, err)
		}

		if count != expectedLength {
			return fmt.Errorf("expected %s to be %d, got %d", attrKey, expectedLength, count)
		}

		return nil
	}
}

func waitForNetworkFunctionHealth(resourceName, attributeName, desiredValue string) resource.TestCheckFunc {
	const maxRetries = 30
	const retryInterval = 10 * time.Second

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		conn := acc.TestAccProvider.Meta().(*conns.Client).NetworkingAPI
		var lastValue string

		for i := 0; i < maxRetries; i++ {
			resp, err := conn.NetworkFunctionAPI.GetNetworkFunctionById(utils.StringPtr(rs.Primary.ID))
			if err != nil {
				return fmt.Errorf("error getting network function by id: %v", err)
			}

			raw := resp.Data.GetValue()
			var nf import1.NetworkFunction
			switch v := raw.(type) {
			case import1.NetworkFunction:
				nf = v
			case *import1.NetworkFunction:
				if v == nil {
					return fmt.Errorf("network function response was nil")
				}
				nf = *v
			default:
				return fmt.Errorf("unexpected network function response type: %T", raw)
			}

			allHealthy := len(nf.NicPairs) > 0
			for _, pair := range nf.NicPairs {
				status := common.FlattenPtrEnum(pair.DataPlaneHealthStatus)
				lastValue = status
				if status != desiredValue {
					allHealthy = false
					break
				}
			}

			if allHealthy {
				return nil
			}

			time.Sleep(retryInterval)
		}

		return fmt.Errorf(
			"network function: failed to reach desired value for attribute %q: expected %q, got %q after %d retries",
			attributeName,
			desiredValue,
			lastValue,
			maxRetries,
		)
	}
}

func testAccCheckNetworkFunctionResourcesDestroy(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	vmClient := conn.VmmAPI.VMAPIInstance
	netClient := conn.NetworkingAPI

	for _, rs := range state.RootModule().Resources {
		switch rs.Type {
		case "nutanix_virtual_machine_v2":
			readResp, err := vmClient.GetVmById(utils.StringPtr(rs.Primary.ID))
			if err != nil {
				continue
			}
			etag := vmClient.ApiClient.GetEtag(readResp)
			args := make(map[string]interface{})
			args["If-Match"] = utils.StringPtr(etag)
			if _, err = vmClient.DeleteVmById(utils.StringPtr(rs.Primary.ID), args); err != nil {
				return fmt.Errorf("error: VM still exists: %v", err)
			}
		case "nutanix_subnet_v2":
			if _, err := netClient.SubnetAPIInstance.GetSubnetById(utils.StringPtr(rs.Primary.ID)); err != nil {
				continue
			}
			if _, err := netClient.SubnetAPIInstance.DeleteSubnetById(utils.StringPtr(rs.Primary.ID)); err != nil {
				return fmt.Errorf("error: Subnet still exists: %v", err)
			}
		case "nutanix_network_function_v2":
			if _, err := netClient.NetworkFunctionAPI.GetNetworkFunctionById(utils.StringPtr(rs.Primary.ID)); err != nil {
				continue
			}
			if _, err := netClient.NetworkFunctionAPI.DeleteNetworkFunctionById(utils.StringPtr(rs.Primary.ID)); err != nil {
				return fmt.Errorf("error: Network function still exists: %v", err)
			}
		}
	}

	return nil
}

func testAccCheckNetworkFunctionPair(resourceName, vmResourceName, ingressAttr, egressAttr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %q not found in state", resourceName)
		}

		vmState, ok := s.RootModule().Resources[vmResourceName]
		if !ok {
			return fmt.Errorf("resource %q not found in state", vmResourceName)
		}

		expectedVM := vmState.Primary.Attributes["id"]
		expectedIngress := vmState.Primary.Attributes[ingressAttr]
		expectedEgress := vmState.Primary.Attributes[egressAttr]

		vmNicExtIDs := make(map[string]struct{})
		for key, value := range vmState.Primary.Attributes {
			if strings.HasPrefix(key, "nics.") && strings.HasSuffix(key, ".ext_id") && value != "" {
				vmNicExtIDs[value] = struct{}{}
			}
		}

		count, err := strconv.Atoi(resourceState.Primary.Attributes["nic_pairs.#"])
		if err != nil {
			return fmt.Errorf("invalid nic_pairs count: %w", err)
		}

		for i := 0; i < count; i++ {
			prefix := fmt.Sprintf("nic_pairs.%d.", i)
			if resourceState.Primary.Attributes[prefix+"vm_reference"] != expectedVM {
				continue
			}
			ingressValue := resourceState.Primary.Attributes[prefix+"ingress_nic_reference"]
			if ingressValue == "" {
				continue
			}
			if _, ok := vmNicExtIDs[ingressValue]; !ok {
				continue
			}
			egressValue := resourceState.Primary.Attributes[prefix+"egress_nic_reference"]
			if egressValue == "" {
				continue
			}
			if _, ok := vmNicExtIDs[egressValue]; !ok {
				continue
			}
			if resourceState.Primary.Attributes[prefix+"is_enabled"] != "true" {
				continue
			}
			return nil
		}

		return fmt.Errorf(
			"no nic_pairs entry matched vm=%q ingress=%q egress=%q from vm nics",
			expectedVM,
			expectedIngress,
			expectedEgress,
		)
	}
}

func testAccCheckNetworkFunctionDataSourcePair(resourceName, vmResourceName, ingressAttr, egressAttr string) resource.TestCheckFunc {
	return testAccCheckNetworkFunctionDataSourcePairWithPrefix(resourceName, "", vmResourceName, ingressAttr, egressAttr)
}

func testAccCheckNetworkFunctionDataSourcePairWithPrefix(resourceName, prefix, vmResourceName, ingressAttr, egressAttr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %q not found in state", resourceName)
		}

		vmState, ok := s.RootModule().Resources[vmResourceName]
		if !ok {
			return fmt.Errorf("resource %q not found in state", vmResourceName)
		}

		expectedVM := vmState.Primary.Attributes["id"]
		expectedIngress := vmState.Primary.Attributes[ingressAttr]
		expectedEgress := vmState.Primary.Attributes[egressAttr]

		vmNicExtIDs := make(map[string]struct{})
		for key, value := range vmState.Primary.Attributes {
			if strings.HasPrefix(key, "nics.") && strings.HasSuffix(key, ".ext_id") && value != "" {
				vmNicExtIDs[value] = struct{}{}
			}
		}

		count, err := strconv.Atoi(resourceState.Primary.Attributes[prefix+"nic_pairs.#"])
		if err != nil {
			return fmt.Errorf("invalid nic_pairs count: %w", err)
		}

		for i := 0; i < count; i++ {
			pairPrefix := fmt.Sprintf("%snic_pairs.%d.", prefix, i)
			if resourceState.Primary.Attributes[pairPrefix+"vm_reference"] != expectedVM {
				continue
			}
			ingressValue := resourceState.Primary.Attributes[pairPrefix+"ingress_nic_reference"]
			if ingressValue == "" {
				continue
			}
			if _, ok := vmNicExtIDs[ingressValue]; !ok {
				continue
			}
			egressValue := resourceState.Primary.Attributes[pairPrefix+"egress_nic_reference"]
			if egressValue == "" {
				continue
			}
			if _, ok := vmNicExtIDs[egressValue]; !ok {
				continue
			}
			if resourceState.Primary.Attributes[pairPrefix+"is_enabled"] != "true" {
				continue
			}
			return nil
		}

		return fmt.Errorf(
			"no data source nic_pairs entry matched vm=%q ingress=%q egress=%q from vm nics",
			expectedVM,
			expectedIngress,
			expectedEgress,
		)
	}
}

func testAccCheckNetworkFunctionNICPairHAStateCounts(resourceName string, expected map[string]int) resource.TestCheckFunc {
	return testAccCheckNetworkFunctionNICPairHAStateCountsWithPrefix(resourceName, "", expected)
}

func testAccCheckNetworkFunctionNICPairHAStateCountsWithPrefix(resourceName, prefix string, expected map[string]int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %q not found in state", resourceName)
		}

		count, err := strconv.Atoi(resourceState.Primary.Attributes[prefix+"nic_pairs.#"])
		if err != nil {
			return fmt.Errorf("invalid nic_pairs count: %w", err)
		}

		actual := make(map[string]int)
		for i := 0; i < count; i++ {
			state := resourceState.Primary.Attributes[fmt.Sprintf("%snic_pairs.%d.high_availability_state", prefix, i)]
			if state == "" {
				continue
			}
			actual[state]++
		}

		totalExpected := 0
		totalActual := 0
		for _, value := range actual {
			totalActual += value
		}
		for state, expectedCount := range expected {
			totalExpected += expectedCount
			if actual[state] != expectedCount {
				return fmt.Errorf(
					"expected %d nic_pairs with high_availability_state=%q, got %d",
					expectedCount,
					state,
					actual[state],
				)
			}
		}

		if totalActual != totalExpected {
			return fmt.Errorf("unexpected high_availability_state values: %v", actual)
		}

		return nil
	}
}

// testAccCheckNutanixNSPDataSourceRulesContainExpectedContent verifies the data source
// rules contain our expected rule types and, when present, the new computed attributes.
// Rule order and count are not guaranteed (API may add default rules).
func testAccCheckNutanixNSPDataSourceRulesContainExpectedContent(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("data source not found: %s", name)
		}
		nStr, ok := rs.Primary.Attributes["rules.#"]
		if !ok || nStr == "0" {
			return fmt.Errorf("rules.# missing or zero")
		}
		n, _ := strconv.Atoi(nStr)
		var hasAppRule, hasIntraRule bool
		for i := 0; i < n; i++ {
			prefix := "rules." + strconv.Itoa(i) + "."
			if rs.Primary.Attributes[prefix+"type"] == "APPLICATION" {
				hasAppRule = true
			}
			if rs.Primary.Attributes[prefix+"type"] == "INTRA_GROUP" {
				hasIntraRule = true
			}
		}
		if !hasAppRule {
			return fmt.Errorf("expected at least one APPLICATION rule in data source rules")
		}
		if !hasIntraRule {
			return fmt.Errorf("expected at least one INTRA_GROUP rule in data source rules")
		}
		return nil
	}
}
