package dataprotectionv2_test

import (
	"fmt"
	"os"
	"strconv"
)

// remotePCProviderConfig returns a `provider "nutanix-2"` block aimed at a
// remote Prism Central by IP/FQDN (availability_zone.remote_pc_ip).
//
// Promote/restore must run on the recovery PC. On NC acceptance runs the
// process env still has NUTANIX_CENTRAL_ENDPOINT / NUTANIX_DOMAIN_NAME, and
// EnvDefaultFunc would pull those into nutanix-2 unless cleared — NC mode
// would then rewrite the host to <domain>.domains.<central> (the *source*
// PC) and the promote call would fail (DP-10000). Classic PC runs leave
// those env vars empty, so clearing them is a no-op there.
//
// Remote AZ is always a direct PC (port 9440 by default). When not in NC
// mode, NUTANIX_PORT is used so classic dual-PC setups keep working.
func remotePCProviderConfig(remotePcIP string) string {
	username := testVars.PCUsername
	password := testVars.PCPassword
	if username == "" {
		username = os.Getenv("NUTANIX_USERNAME")
		password = os.Getenv("NUTANIX_PASSWORD")
	}

	port := 9440
	if os.Getenv("NUTANIX_CENTRAL_ENDPOINT") == "" {
		if p, err := strconv.Atoi(os.Getenv("NUTANIX_PORT")); err == nil && p > 0 {
			port = p
		}
	}

	insecure, _ := strconv.ParseBool(os.Getenv("NUTANIX_INSECURE"))

	return fmt.Sprintf(`
provider "nutanix-2" {
  username                 = "%[1]s"
  password                 = "%[2]s"
  endpoint                 = "%[3]s"
  insecure                 = %[4]t
  port                     = %[5]d
  nutanix_central_endpoint = ""
  domain_name              = ""
  tenant_name              = ""
  use_nc_iam               = false
}
`, username, password, remotePcIP, insecure, port)
}
