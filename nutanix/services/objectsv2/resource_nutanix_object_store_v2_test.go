package objectstoresv2_test

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameObjectStore = "nutanix_object_store_v2.test"
const resourceNameObjectStoreCertificate = "nutanix_object_store_certificate_v2.test"

const datasourceNameObjectStoreFetch = "data.nutanix_object_store_v2.fetch"
const datasourceNameObjectStoreList = "data.nutanix_object_stores_v2.list"
const datasourceNameObjectStoreLimit = "data.nutanix_object_stores_v2.limit"
const datasourceNameObjectStoreFilter = "data.nutanix_object_stores_v2.filter"
const datasourceNameCertificateFetch = "data.nutanix_certificate_v2.fetch"
const datasourceNameCertificatesList = "data.nutanix_certificates_v2.list"

func TestAccV2NutanixObjectStoreResource_OneWorkerNode(t *testing.T) {
	r := acctest.RandIntRange(1, 99)
	objectStoreName := fmt.Sprintf("tf-test-os-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStoreWithOneWorkerNodeConfig(objectStoreName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameObjectStore, "name", objectStoreName),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "description", "terraform test object store"),
					resource.TestCheckResourceAttrSet(resourceNameObjectStore, "deployment_version"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "domain", testVars.ObjectStore.Domain),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "num_worker_nodes", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "total_capacity_gib", fmt.Sprintf("%d", 20*int(math.Pow(1024, 3)))),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "public_network_ips.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "public_network_ips.0.ipv4.0.value", testVars.ObjectStore.PublicNetworkIPs[0]),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_dns_ip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_dns_ip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkDNSIP),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkVip),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "state", "OBJECT_STORE_AVAILABLE"),
				),
			},
			// list object store with filter and limit
			{
				Config: testAccObjectStoreWithOneWorkerNodeConfig(objectStoreName) + testAccObjectStoreDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					// fetch object store check
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "name", datasourceNameObjectStoreFetch, "name"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "description", datasourceNameObjectStoreFetch, "description"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "deployment_version", datasourceNameObjectStoreFetch, "deployment_version"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "domain", datasourceNameObjectStoreFetch, "domain"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "num_worker_nodes", datasourceNameObjectStoreFetch, "num_worker_nodes"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "total_capacity_gib", datasourceNameObjectStoreFetch, "total_capacity_gib"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "public_network_ips.#", datasourceNameObjectStoreFetch, "public_network_ips.#"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "public_network_ips.0.ipv4.0.value", datasourceNameObjectStoreFetch, "public_network_ips.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "storage_network_dns_ip.#", datasourceNameObjectStoreFetch, "storage_network_dns_ip.#"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "storage_network_dns_ip.0.ipv4.0.value", datasourceNameObjectStoreFetch, "storage_network_dns_ip.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "storage_network_vip.#", datasourceNameObjectStoreFetch, "storage_network_vip.#"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "storage_network_vip.0.ipv4.0.value", datasourceNameObjectStoreFetch, "storage_network_vip.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "state", datasourceNameObjectStoreFetch, "state"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "cluster_ext_id", datasourceNameObjectStoreFetch, "cluster_ext_id"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "public_network_reference", datasourceNameObjectStoreFetch, "public_network_reference"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "storage_network_reference", datasourceNameObjectStoreFetch, "storage_network_reference"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "ext_id", datasourceNameObjectStoreFetch, "ext_id"),

					// list object store check
					resource.TestCheckResourceAttrSet(datasourceNameObjectStoreList, "object_stores.#"),

					// filter object store check
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "name", datasourceNameObjectStoreFilter, "object_stores.0.name"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "description", datasourceNameObjectStoreFilter, "object_stores.0.description"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "deployment_version", datasourceNameObjectStoreFilter, "object_stores.0.deployment_version"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "domain", datasourceNameObjectStoreFilter, "object_stores.0.domain"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "num_worker_nodes", datasourceNameObjectStoreFilter, "object_stores.0.num_worker_nodes"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "total_capacity_gib", datasourceNameObjectStoreFilter, "object_stores.0.total_capacity_gib"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "public_network_ips.#", datasourceNameObjectStoreFilter, "object_stores.0.public_network_ips.#"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "public_network_ips.0.ipv4.0.value", datasourceNameObjectStoreFilter, "object_stores.0.public_network_ips.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "storage_network_dns_ip.#", datasourceNameObjectStoreFilter, "object_stores.0.storage_network_dns_ip.#"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "storage_network_dns_ip.0.ipv4.0.value", datasourceNameObjectStoreFilter, "object_stores.0.storage_network_dns_ip.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "storage_network_vip.#", datasourceNameObjectStoreFilter, "object_stores.0.storage_network_vip.#"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "storage_network_vip.0.ipv4.0.value", datasourceNameObjectStoreFilter, "object_stores.0.storage_network_vip.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "state", datasourceNameObjectStoreFilter, "object_stores.0.state"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "cluster_ext_id", datasourceNameObjectStoreFilter, "object_stores.0.cluster_ext_id"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "public_network_reference", datasourceNameObjectStoreFilter, "object_stores.0.public_network_reference"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "storage_network_reference", datasourceNameObjectStoreFilter, "object_stores.0.storage_network_reference"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "ext_id", datasourceNameObjectStoreFilter, "object_stores.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameObjectStoreFilter, "object_stores.#", "1"),

					// list object store with limit check
					resource.TestCheckResourceAttrSet(datasourceNameObjectStoreList, "object_stores.#"),
					resource.TestCheckResourceAttr(datasourceNameObjectStoreLimit, "object_stores.#", "1"),
				),
			},
			// create a new certificate for object store
			{
				Config: testAccObjectStoreWithOneWorkerNodeConfig(objectStoreName) + testAccObjectStoreCertificateResourceConfig() + testAccObjectStoreCertificateDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Fetch object store certificate check
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "ext_id", datasourceNameCertificateFetch, "object_store_ext_id"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "public_network_ips.0.ipv4.0.value", datasourceNameCertificateFetch, "alternate_ips.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStoreCertificate, "id", datasourceNameCertificateFetch, "ext_id"),

					// List object store certificate check
					resource.TestCheckResourceAttrSet(datasourceNameCertificatesList, "certificates.#"),
					resource.TestCheckResourceAttrPair(datasourceNameCertificatesList, "certificates.0.ext_id", resourceNameObjectStoreCertificate, "id"),

					// delete object store bucket
					deleteObjectStoreBucket(),
				),
			},
		},
	})
}

func TestAccV2NutanixObjectStoreResource_DraftObjectStore(t *testing.T) {
	r := acctest.RandIntRange(1, 99)
	objectStoreName := fmt.Sprintf("tf-test-os-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStoreUndeployedObjectStoreConfig(objectStoreName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameObjectStore, "name", objectStoreName),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "description", "terraform test object store"),
					resource.TestCheckResourceAttrSet(resourceNameObjectStore, "deployment_version"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "domain", testVars.ObjectStore.Domain),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "num_worker_nodes", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "total_capacity_gib", fmt.Sprintf("%d", 20*int(math.Pow(1024, 3)))),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "public_network_ips.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "public_network_ips.0.ipv4.0.value", testVars.ObjectStore.PublicNetworkIPs[0]),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_dns_ip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_dns_ip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkDNSIP),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkVip),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "state", "UNDEPLOYED_OBJECT_STORE"),
				),
			},
		},
	})
}

func TestAccV2NutanixObjectStoreResource_UpdateObjectStore(t *testing.T) {
	r := acctest.RandIntRange(1, 99)
	objectStoreName := fmt.Sprintf("tf-test-os-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccObjectStoreWithInvalidImageTagConfig(objectStoreName),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameObjectStore, "id"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "state", "OBJECT_STORE_DEPLOYMENT_FAILED"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Updating object store ")
				},
				Config: testAccObjectStoreWithInvalidImageTagConfig(objectStoreName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameObjectStore, "name", objectStoreName),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "description", "terraform test object store"),
					resource.TestCheckResourceAttrSet(resourceNameObjectStore, "deployment_version"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "domain", testVars.ObjectStore.Domain),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "num_worker_nodes", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "total_capacity_gib", fmt.Sprintf("%d", 20*int(math.Pow(1024, 3)))),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "public_network_ips.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "public_network_ips.0.ipv4.0.value", testVars.ObjectStore.PublicNetworkIPs[0]),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_dns_ip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_dns_ip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkDNSIP),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkVip),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "state", "OBJECT_STORE_AVAILABLE"),

					// delete object store bucket
					deleteObjectStoreBucket(),
				),
			},
		},
	})
}
func testAccObjectStoreWithOneWorkerNodeConfig(objectStoreName string) string {
	return fmt.Sprintf(`

locals {
  config = jsondecode(file("%[1]s"))
  objectStore      = local.config.object_store
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
  subnetExtId = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
}

data "nutanix_clusters_v2" "clusters" {}

data "nutanix_subnets_v2" "subnets" {
  filter = "name eq '${local.objectStore.subnet_name}'"
}

resource "nutanix_object_store_v2" "test" {
  timeouts {
    create = "120m"
  }
  name                     = "%[2]s"
  description              = "terraform test object store"
  domain                   = local.objectStore.domain
  num_worker_nodes         = 1
  cluster_ext_id           = local.clusterExtId
  total_capacity_gib       = 20 * pow(1024, 3)

  public_network_reference = local.subnetExtId
  public_network_ips {
    ipv4 {
      value = local.objectStore.public_network_ips[0]
    }
  }

  storage_network_reference = local.subnetExtId
  storage_network_dns_ip {
    ipv4 {
      value = local.objectStore.storage_network_dns_ip
    }
  }
  storage_network_vip {
    ipv4 {
      value = local.objectStore.storage_network_vip
    }
  }
}

`, filepath, objectStoreName)
}

func testAccObjectStoreUndeployedObjectStoreConfig(objectStoreName string) string {
	return fmt.Sprintf(`

locals {
  config = jsondecode(file("%[1]s"))
  objectStore      = local.config.object_store
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
  subnetExtId = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
}

data "nutanix_clusters_v2" "clusters" {}

data "nutanix_subnets_v2" "subnets" {
  filter = "name eq '${local.objectStore.subnet_name}'"
}

resource "nutanix_object_store_v2" "test" {
  timeouts {
    create = "120m"
  }
  name                     = "%[2]s"
  description              = "terraform test object store"
  domain                   = local.objectStore.domain
  num_worker_nodes         = 1
  cluster_ext_id           = local.clusterExtId
  total_capacity_gib       = 20 * pow(1024, 3)
  state					   = "UNDEPLOYED_OBJECT_STORE"

  public_network_reference = local.subnetExtId
  public_network_ips {
    ipv4 {
      value = local.objectStore.public_network_ips[0]
    }
  }

  storage_network_reference = local.subnetExtId
  storage_network_dns_ip {
    ipv4 {
      value = local.objectStore.storage_network_dns_ip
    }
  }
  storage_network_vip {
    ipv4 {
      value = local.objectStore.storage_network_vip
    }
  }
}

`, filepath, objectStoreName)
}

// list object store with filter and limit
func testAccObjectStoreDatasourceConfig() string {
	return `
data "nutanix_object_stores_v2" "list" {
  depends_on = [nutanix_object_store_v2.test]
}

data "nutanix_object_stores_v2" "filter" {
  filter = "name eq '${nutanix_object_store_v2.test.name}'"
  depends_on = [nutanix_object_store_v2.test]
}

data "nutanix_object_stores_v2" "limit" {
  limit      = 1
  depends_on = [nutanix_object_store_v2.test]
}

data "nutanix_object_store_v2" "fetch" {
  ext_id = nutanix_object_store_v2.test.id
}


`
}

func testAccObjectStoreCertificateResourceConfig() string {
	return fmt.Sprintf(`
resource "nutanix_object_store_certificate_v2" "test" {
  object_store_ext_id = nutanix_object_store_v2.test.id
  path                = "%s"
}

data "nutanix_certificate_v2" "fetch" {
  object_store_ext_id = nutanix_object_store_v2.test.id
  ext_id              = nutanix_object_store_certificate_v2.test.id
  depends_on = [nutanix_object_store_certificate_v2.test]
}
`, certificateJSONFile)
}

func testAccObjectStoreCertificateDatasourceConfig() string {
	return `
data "nutanix_certificates_v2" "list" {
  object_store_ext_id = nutanix_object_store_v2.test.id
  depends_on = [nutanix_object_store_certificate_v2.test]
}
`
}

func testAccObjectStoreWithInvalidImageTagConfig(objectStoreName string) string {
	endpoint := os.Getenv("NUTANIX_ENDPOINT")

	return fmt.Sprintf(`

locals {
  config = jsondecode(file("%[1]s"))
  objectStore      = local.config.object_store
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
  subnetExtId = data.nutanix_subnets_v2.subnets.subnets[0].ext_id

  pcSSHPassword    = local.objectStore.ssh_pc_password
  username         = local.objectStore.ssh_pc_username
  ip               = "%[2]s"

}

locals {
  pre_update_hook_command = <<EOT
  sshpass -p '${local.pcSSHPassword}' ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${local.username}@${local.ip} \
  "docker exec aoss_service_manager sh -c 'cd /home/nutanix/config/poseidon_master && \
  cp buckets_tools_template.yaml buckets_tools_template_backup.yml && \
  sed -i -E \"s|(image: .+/[^:]+:)[^ ]+|\\\\1invalid-version|\" buckets_tools_template.yaml'"
  EOT
}

locals {
  restore_command = <<EOT
  sshpass -p '${local.pcSSHPassword}' ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${local.username}@${local.ip} \
  "docker exec aoss_service_manager sh -c 'cd /home/nutanix/config/poseidon_master && mv buckets_tools_template_backup.yml buckets_tools_template.yaml'"
  EOT
}

data "nutanix_clusters_v2" "clusters" {}

data "nutanix_subnets_v2" "subnets" {
  filter = "name eq '${local.objectStore.subnet_name}'"
}

# this resource to change image tag on pc to incorrect one
# to make sure object store deployment fails
resource "terraform_data" "pre_update_hook" {
  provisioner "local-exec" {
    when       = create
    command    = local.pre_update_hook_command
    on_failure = continue
  }
}


resource "nutanix_object_store_v2" "test" {
  timeouts {
    create = "120m"
    update = "120m"
  }
  name                     = "%[3]s"
  description              = "terraform test object store"
  domain                   = local.objectStore.domain
  num_worker_nodes         = 1
  cluster_ext_id           = local.clusterExtId
  total_capacity_gib       = 20 * pow(1024, 3)

  public_network_reference = local.subnetExtId
  public_network_ips {
    ipv4 {
      value = local.objectStore.public_network_ips[0]
    }
  }

  storage_network_reference = local.subnetExtId
  storage_network_dns_ip {
    ipv4 {
      value = local.objectStore.storage_network_dns_ip
    }
  }
  storage_network_vip {
    ipv4 {
      value = local.objectStore.storage_network_vip
    }
  }
  depends_on = [terraform_data.pre_update_hook]
}

# this resource to change image tag on pc to correct one
# to make sure object store deployment succeeds
resource "terraform_data" "post_update_hook" {
  provisioner "local-exec" {
    command    = local.restore_command
    on_failure = continue
	when	   = create
  }
  depends_on = [nutanix_object_store_v2.test]
}

`, filepath, endpoint, objectStoreName)
}
