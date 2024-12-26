package networkingv2_test

import (
	"fmt"
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
