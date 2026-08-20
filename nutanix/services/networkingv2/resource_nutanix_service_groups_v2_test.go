package networkingv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameServiceGroup = "nutanix_service_groups_v2.test"

func TestAccV2NutanixServiceGroupResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-Service-group-%d", r)
	desc := "test Service group description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testServiceGroupV2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "tcp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "tcp_services.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "udp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "udp_services.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixServiceGroupResource_WithUpdate(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-Service-group-%d", r)
	updatedName := fmt.Sprintf("test-Service-group-%d", r+1)
	updatedDesc := "test Service group description updated"
	desc := "test Service group description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testServiceGroupV2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "tcp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "tcp_services.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "udp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "udp_services.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "ext_id"),
				),
			},
			{
				Config: testServiceGroupV2Config(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "tcp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "tcp_services.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "udp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "udp_services.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixServiceGroupResource_WithICMP(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-Service-group-%d", r)
	desc := "test Service group description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testServiceGroupV2ConfigWithICMP(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "icmp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "icmp_services.#", "1"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "icmp_services.0.type", "8"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "icmp_services.0.code", "0"),
				),
			},
		},
	})
}

// TestAccV2NutanixServiceGroupResource_WithICMPAllAllowed covers the ICMP
// wildcard case: is_all_allowed = true with no type/code. The provider must
// omit type/code from the payload, otherwise the microseg v4.2 API rejects it
// with MIC-30302 ("ICMP is set to allow all but also contains specific type/code").
func TestAccV2NutanixServiceGroupResource_WithICMPAllAllowed(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-Service-group-%d", r)
	desc := "test Service group ICMP wildcard"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testServiceGroupV2ConfigWithICMPAllAllowed(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "links.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "icmp_services.#", "1"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "icmp_services.0.is_all_allowed", "true"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "icmp_services.0.type", "0"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "icmp_services.0.code", "0"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixServiceGroupResource_WithAll(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-Service-group-%d", r)
	desc := "test Service group description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testServiceGroupV2ConfigWithAll(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "tcp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "tcp_services.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "udp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "udp_services.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "icmp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "icmp_services.#", "1"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "icmp_services.0.type", "8"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "icmp_services.0.code", "0"),
				),
			},
		},
	})
}

func TestAccV2NutanixServiceGroupResource_WithUpdateTCP(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-Service-group-%d", r)
	updatedName := fmt.Sprintf("test-Service-group-%d", r+1)
	updatedDesc := "test Service group description updated"
	desc := "test Service group description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testServiceGroupV2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "tcp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "tcp_services.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "udp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "udp_services.#", "1"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "tcp_services.0.start_port", "232"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "tcp_services.0.end_port", "232"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "ext_id"),
				),
			},
			{
				Config: testServiceGroupV2ConfigWithUpdateTCP(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "tcp_services.#"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "tcp_services.#", "1"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "tcp_services.0.start_port", "211"),
					resource.TestCheckResourceAttr(resourceNameServiceGroup, "tcp_services.0.end_port", "221"),
					resource.TestCheckResourceAttrSet(resourceNameServiceGroup, "ext_id"),
				),
			},
		},
	})
}

func testServiceGroupV2Config(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_service_groups_v2" "test" {
		name  = "%[1]s"
		description = "%[2]s"  
		tcp_services {
		  start_port = "232"
		  end_port = "232"
		}
		udp_services {
		  start_port = "232"
		  end_port = "232"
		}
	  }
`, name, desc)
}

func testServiceGroupV2ConfigWithICMP(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_service_groups_v2" "test" {
		name  = "%[1]s"
		description = "%[2]s"  
		icmp_services {
			type = 8
			code = 0
	 	}
	}
`, name, desc)
}

func testServiceGroupV2ConfigWithICMPAllAllowed(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_service_groups_v2" "test" {
		name  = "%[1]s"
		description = "%[2]s"  
		icmp_services {
			is_all_allowed = true
	 	}
	}
`, name, desc)
}

func testServiceGroupV2ConfigWithAll(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_service_groups_v2" "test" {
		name  = "%[1]s"
		description = "%[2]s"  
		tcp_services {
		  start_port = "232"
		  end_port = "232"
		}
		udp_services {
		  start_port = "232"
		  end_port = "232"
		}
		icmp_services {
			type = 8
			code = 0
	  	}
	}
`, name, desc)
}

func testServiceGroupV2ConfigWithUpdateTCP(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_service_groups_v2" "test" {
		name  = "%[1]s"
		description = "%[2]s"  
		tcp_services {
		  start_port = "211"
		  end_port = "221"
		}
	  }
`, name, desc)
}

func TestAccV2NutanixServiceGroupResource_ProjectAssociation(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-sg-projassoc-%d", r)
	desc := "service group project association test"
	projectName := fmt.Sprintf("tf-sg-pa-proj-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testServiceGroupProjectAssociationConfig(name, desc, projectName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameServiceGroup, "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_service_group_v2.test", "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttr("data.nutanix_service_groups_v2.test", "service_groups.#", "1"),
					resource.TestCheckResourceAttrPair("data.nutanix_service_groups_v2.test", "service_groups.0.ext_id", resourceNameServiceGroup, "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_service_groups_v2.test", "service_groups.0.project_ext_id", "nutanix_project_v2.test", "ext_id"),
				),
			},
			{
				Config:      testServiceGroupProjectAssociationConfig(name, desc, projectName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile("Update of project_ext_id is not supported"),
			},
		},
	})
}

func sgProjectExtIDLine(override string) string {
	if override == "" {
		return `project_ext_id = nutanix_project_v2.test.ext_id`
	}
	return fmt.Sprintf(`project_ext_id = "%s"`, override)
}

func testServiceGroupProjectAssociationConfig(name, desc, projectName, projectExtIDOverride string) string {
	return fmt.Sprintf(`
	resource "nutanix_project_v2" "test" {
		name        = "%[3]s"
		project_id  = "%[3]s"
		description = "project association test"
	}

	resource "nutanix_service_groups_v2" "test" {
		name        = "%[1]s"
		description = "%[2]s"
		tcp_services {
			start_port = "232"
			end_port   = "232"
		}
		udp_services {
			start_port = "232"
			end_port   = "232"
		}
		%[4]s
		depends_on = [nutanix_project_v2.test]
	}

	data "nutanix_service_group_v2" "test" {
		ext_id     = nutanix_service_groups_v2.test.id
		depends_on = [nutanix_service_groups_v2.test]
	}

	data "nutanix_service_groups_v2" "test" {
		filter     = "name eq '%[1]s'"
		depends_on = [nutanix_service_groups_v2.test]
	}
	`, name, desc, projectName, sgProjectExtIDLine(projectExtIDOverride))
}
