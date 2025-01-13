package dataprotectionv2_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"time"
)

func waitForVmToBeProtected(resourceName, attributeName, desiredValue string, maxRetries int, retryInterval time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var lastValue string
		for i := 0; i < maxRetries; i++ {
			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("resource not found: %s", resourceName)
			}

			lastValue = rs.Primary.Attributes[attributeName]
			if lastValue == desiredValue {
				return nil // Desired value reached
			}

			// Wait before retrying
			time.Sleep(retryInterval)
		}

		return fmt.Errorf("failed to reach desired value for attribute %q: expected %q, got %q after %d retries", attributeName, desiredValue, lastValue, maxRetries)
	}
}
