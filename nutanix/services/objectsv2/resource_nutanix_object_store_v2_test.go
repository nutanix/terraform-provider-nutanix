package objectstoresv2_test

import (
	"fmt"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameObjectStore = "nutanix_object_store_v2.test"
const resourceNameObjectStore2 = "nutanix_object_store_v2.test2"
const resourceNameObjectStoreCertificate = "nutanix_object_store_certificate_v2.test"

const datasourceNameObjectStoreFetch = "data.nutanix_object_store_v2.fetch"
const datasourceNameObjectStoreList = "data.nutanix_object_stores_v2.list"
const datasourceNameObjectStoreLimit = "data.nutanix_object_stores_v2.limit"
const datasourceNameObjectStoreFilter = "data.nutanix_object_stores_v2.filter"
const datasourceNameCertificateFetch = "data.nutanix_certificate_v2.fetch"
const datasourceNameCertificatesList = "data.nutanix_certificates_v2.list"

// object store OVA resource name
const resourceNameObjectLiteStoreImage = "nutanix_images_v2.object-liteStore-img"
const resourceNameVM = "nutanix_virtual_machine_v2.vm-test"

const resourceNameVMOva = "nutanix_ova_v2.vm-ova"
const resourceNameObjectLiteSourceOva = "nutanix_ova_v2.object-liteSource-ova"

func TestAccV2NutanixObjectStoreResource_OneWorkerNode(t *testing.T) {
	r := acctest.RandIntRange(1, 99)
	objectStoreName := fmt.Sprintf("tf-test-os-%d", r)
	r2 := acctest.RandIntRange(100, 199)
	objectStoreName2 := fmt.Sprintf("tf-test-os-%d", r2)

	objectLiteSourceImgName := fmt.Sprintf("tf-object-ls-img-%d", r)
	vmName := fmt.Sprintf("tf-object-vm-%d", r)
	vmOvaName := fmt.Sprintf("tf-vm-ova-%d", r)
	objectOvaName := fmt.Sprintf("tf-object-liteStore-ova-%d", r)

	config := testAccObjectStoreWithOneWorkerNodeConfig(objectStoreName, objectStoreName2)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
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
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_dns_ip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkDNSIP[0]),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkVip[0]),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "state", "OBJECT_STORE_AVAILABLE"),
					// secand object store check
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "name", objectStoreName2),
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "description", "terraform test object store second object store"),
					resource.TestCheckResourceAttrSet(resourceNameObjectStore2, "deployment_version"),
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "domain", testVars.ObjectStore.Domain),
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "num_worker_nodes", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "total_capacity_gib", fmt.Sprintf("%d", 20*int(math.Pow(1024, 3)))),
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "public_network_ips.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "public_network_ips.0.ipv4.0.value", testVars.ObjectStore.PublicNetworkIPs[1]),
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "storage_network_dns_ip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "storage_network_dns_ip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkDNSIP[1]),
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "storage_network_vip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "storage_network_vip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkVip[1]),
					resource.TestCheckResourceAttr(resourceNameObjectStore2, "state", "OBJECT_STORE_AVAILABLE"),
				),
			},
			// list object store with filter and limit
			{
				Config: config + testAccObjectStoreDatasourceConfig(),
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
				Config: config + testAccObjectStoreCertificateResourceConfig() + testAccObjectStoreCertificateDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Fetch object store certificate check
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "ext_id", datasourceNameCertificateFetch, "object_store_ext_id"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStore, "public_network_ips.0.ipv4.0.value", datasourceNameCertificateFetch, "alternate_ips.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(resourceNameObjectStoreCertificate, "id", datasourceNameCertificateFetch, "ext_id"),

					// List object store certificate check
					resource.TestCheckResourceAttrSet(datasourceNameCertificatesList, "certificates.#"),
					resource.TestCheckResourceAttrPair(datasourceNameCertificatesList, "certificates.0.ext_id", resourceNameObjectStoreCertificate, "id"),
				),
			},
			// Lite source tests for image and ova using object store source
			{
				PreConfig: func() {
					fmt.Println("object lite source tests")
				},
				Config: config + testAccObjectStoreObjectLiteSourceConfig(objectLiteSourceImgName, vmName, vmOvaName, objectOvaName),
				Check: resource.ComposeTestCheckFunc(
					// check object store image
					resource.TestCheckResourceAttrSet(resourceNameObjectLiteStoreImage, "id"),
					resource.TestCheckResourceAttrSet(resourceNameObjectLiteStoreImage, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameObjectLiteStoreImage, "name", objectLiteSourceImgName),
					resource.TestCheckResourceAttr(resourceNameObjectLiteStoreImage, "description", "Image created from object store"),
					resource.TestCheckResourceAttr(resourceNameObjectLiteStoreImage, "type", "DISK_IMAGE"),
					resource.TestCheckResourceAttrSet(resourceNameObjectLiteStoreImage, "source.0.url_source.0.url"),
					func(s *terraform.State) error {
						attr := s.RootModule().Resources[resourceNameObjectLiteStoreImage].Primary.Attributes["source.0.url_source.0.url"]

						if !strings.Contains(attr, objectStoreName) {
							return fmt.Errorf("expected URL to contain '%s', got: %s", objectStoreName, attr)
						}

						if !strings.Contains(attr, "objects") {
							return fmt.Errorf("expected URL to contain 'objects', got: %s", attr)
						}

						return nil
					},

					// check vm
					resource.TestCheckResourceAttrSet(resourceNameVM, "id"),
					resource.TestCheckResourceAttrSet(resourceNameVM, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVM, "name", vmName),
					resource.TestCheckResourceAttr(resourceNameVM, "description", "terraform test object store vm"),
					resource.TestCheckResourceAttr(resourceNameVM, "num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceNameVM, "memory_size_bytes", fmt.Sprintf("%d", 4*1024*1024*1024)),

					// check vm ova
					resource.TestCheckResourceAttrSet(resourceNameVMOva, "id"),
					resource.TestCheckResourceAttrSet(resourceNameVMOva, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVMOva, "name", vmOvaName),

					// check object lite source ova
					resource.TestCheckResourceAttrSet(resourceNameObjectLiteSourceOva, "id"),
					resource.TestCheckResourceAttrSet(resourceNameObjectLiteSourceOva, "ext_id"),
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
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStoreUndeployedObjectStoreConfig(objectStoreName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameObjectStore, "name", objectStoreName),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "description", "terraform test object store"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "domain", testVars.ObjectStore.Domain),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "num_worker_nodes", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "total_capacity_gib", fmt.Sprintf("%d", 20*int(math.Pow(1024, 3)))),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "public_network_ips.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "public_network_ips.0.ipv4.0.value", testVars.ObjectStore.PublicNetworkIPs[0]),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_dns_ip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_dns_ip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkDNSIP[0]),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkVip[0]),
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
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixObjectStoreDestroy,
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
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_dns_ip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkDNSIP[0]),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.#", "1"),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "storage_network_vip.0.ipv4.0.value", testVars.ObjectStore.StorageNetworkVip[0]),
					resource.TestCheckResourceAttr(resourceNameObjectStore, "state", "OBJECT_STORE_AVAILABLE"),

					// delete object store bucket
					deleteObjectStoreBucket(),
				),
			},
		},
	})
}
func testAccObjectStoreWithOneWorkerNodeConfig(objectStoreName, objectStoreName2 string) string {
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
      value = local.objectStore.storage_network_dns_ip[0]
    }
  }
  storage_network_vip {
    ipv4 {
      value = local.objectStore.storage_network_vip[0]
    }
  }
}

// Ensure bucket is deleted before the object store is destroyed, even if a test step fails.
// Terraform destroys resources in reverse dependency order, so this will run *before* the object store delete.
resource "terraform_data" "cleanup_bucket_test" {
  input = {
    object_store_id = nutanix_object_store_v2.test.id
    bucket_name     = local.objectStore.bucket_name
  }
  provisioner "local-exec" {
    when = destroy
    command = <<EOT
set -eu
BASE="https://$NUTANIX_ENDPOINT:$NUTANIX_PORT/oss/api/nutanix/v3/objectstore_proxy/${self.input.object_store_id}"
AUTH="$NUTANIX_USERNAME:$NUTANIX_PASSWORD"

list_buckets() {
  curl -sSk -u "$AUTH" "$BASE/buckets" || true
}

delete_bucket() {
  b="$1"
  url="$BASE/buckets/$b?force_empty_bucket=true"
  code="$(curl -sSk -u "$AUTH" -X DELETE "$url" -o /tmp/os_bucket_delete_test.out -w "%%%%{http_code}" || echo "000")"
  case "$code" in
    200|202|204|404|503) return 0 ;;
    500) # observed as "Bucket lookup failed"; treat as non-fatal so destroy can proceed
      return 0
      ;;
    *) echo "bucket delete failed (test) http_code=$code url=$url"; cat /tmp/os_bucket_delete_test.out || true; return 1 ;;
  esac
}

# 1) Prefer listing actual buckets and deleting all of them
payload="$(list_buckets)"
names="$(python3 - <<'PY' 2>/dev/null || true
import json,sys
try:
  data=json.loads(sys.stdin.read() or "{}")
except Exception:
  sys.exit(0)
items = data.get("entities") or data.get("buckets") or []
out=[]
for it in items:
  if isinstance(it,str):
    out.append(it)
  elif isinstance(it,dict):
    n=it.get("name")
    if isinstance(n,str):
      out.append(n)
print("\n".join([n for n in out if n]))
PY
)"

if [ -n "$names" ]; then
  for b in $names; do
    # retry a few times in case OSS proxy is transient
    i=0
    while [ $i -lt 5 ]; do
      if delete_bucket "$b"; then break; fi
      i=$((i+1))
      sleep 5
    done
  done
else
  # 2) Fallback: try configured bucket name
  delete_bucket "${self.input.bucket_name}" || true
fi

exit 0
EOT
  }
  depends_on = [nutanix_object_store_v2.test]
}

# second object store to test
resource "nutanix_object_store_v2" "test2" {
  timeouts {
    create = "120m"
  }
  name                     = "%[3]s"
  description              = "terraform test object store second object store"
  domain                   = local.objectStore.domain
  num_worker_nodes         = 1
  cluster_ext_id           = local.clusterExtId
  total_capacity_gib       = 20 * pow(1024, 3)

  public_network_reference = local.subnetExtId
  public_network_ips {
    ipv4 {
      value = local.objectStore.public_network_ips[1]
    }
  }

  storage_network_reference = local.subnetExtId
  storage_network_dns_ip {
    ipv4 {
      value = local.objectStore.storage_network_dns_ip[1]
    }
  }
  storage_network_vip {
    ipv4 {
      value = local.objectStore.storage_network_vip[1]
    }
  }
  # wait for first object store to be created
  depends_on = [nutanix_object_store_v2.test]
}

resource "terraform_data" "cleanup_bucket_test2" {
  input = {
    object_store_id = nutanix_object_store_v2.test2.id
    bucket_name     = local.objectStore.bucket_name
  }
  provisioner "local-exec" {
    when = destroy
    command = <<EOT
set -eu
BASE="https://$NUTANIX_ENDPOINT:$NUTANIX_PORT/oss/api/nutanix/v3/objectstore_proxy/${self.input.object_store_id}"
AUTH="$NUTANIX_USERNAME:$NUTANIX_PASSWORD"

list_buckets() {
  curl -sSk -u "$AUTH" "$BASE/buckets" || true
}

delete_bucket() {
  b="$1"
  url="$BASE/buckets/$b?force_empty_bucket=true"
  code="$(curl -sSk -u "$AUTH" -X DELETE "$url" -o /tmp/os_bucket_delete_test2.out -w "%%%%{http_code}" || echo "000")"
  case "$code" in
    200|202|204|404|503) return 0 ;;
    500) return 0 ;;
    *) echo "bucket delete failed (test2) http_code=$code url=$url"; cat /tmp/os_bucket_delete_test2.out || true; return 1 ;;
  esac
}

payload="$(list_buckets)"
names="$(python3 - <<'PY' 2>/dev/null || true
import json,sys
try:
  data=json.loads(sys.stdin.read() or "{}")
except Exception:
  sys.exit(0)
items = data.get("entities") or data.get("buckets") or []
out=[]
for it in items:
  if isinstance(it,str):
    out.append(it)
  elif isinstance(it,dict):
    n=it.get("name")
    if isinstance(n,str):
      out.append(n)
print("\n".join([n for n in out if n]))
PY
)"

if [ -n "$names" ]; then
  for b in $names; do
    i=0
    while [ $i -lt 5 ]; do
      if delete_bucket "$b"; then break; fi
      i=$((i+1))
      sleep 5
    done
  done
else
  delete_bucket "${self.input.bucket_name}" || true
fi

exit 0
EOT
  }
  depends_on = [nutanix_object_store_v2.test2]
}

`, filepath, objectStoreName, objectStoreName2)
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
      value = local.objectStore.storage_network_dns_ip[0]
    }
  }
  storage_network_vip {
    ipv4 {
      value = local.objectStore.storage_network_vip[0]
    }
  }
}

// Ensure bucket is deleted before the object store is destroyed, even if a test step fails.
resource "terraform_data" "cleanup_bucket_test" {
  input = {
    object_store_id = nutanix_object_store_v2.test.id
    bucket_name     = local.objectStore.bucket_name
  }
  provisioner "local-exec" {
    when = destroy
    command = <<EOT
set -eu
BASE="https://$NUTANIX_ENDPOINT:$NUTANIX_PORT/oss/api/nutanix/v3/objectstore_proxy/${self.input.object_store_id}"
AUTH="$NUTANIX_USERNAME:$NUTANIX_PASSWORD"

list_buckets() {
  curl -sSk -u "$AUTH" "$BASE/buckets" || true
}

delete_bucket() {
  b="$1"
  url="$BASE/buckets/$b?force_empty_bucket=true"
  code="$(curl -sSk -u "$AUTH" -X DELETE "$url" -o /tmp/os_bucket_delete_test.out -w "%%%%{http_code}" || echo "000")"
  case "$code" in
    200|202|204|404|503) return 0 ;;
    500) return 0 ;;
    *) echo "bucket delete failed (test) http_code=$code url=$url"; cat /tmp/os_bucket_delete_test.out || true; return 1 ;;
  esac
}

payload="$(list_buckets)"
names="$(python3 - <<'PY' 2>/dev/null || true
import json,sys
try:
  data=json.loads(sys.stdin.read() or "{}")
except Exception:
  sys.exit(0)
items = data.get("entities") or data.get("buckets") or []
out=[]
for it in items:
  if isinstance(it,str):
    out.append(it)
  elif isinstance(it,dict):
    n=it.get("name")
    if isinstance(n,str):
      out.append(n)
print("\n".join([n for n in out if n]))
PY
)"

if [ -n "$names" ]; then
  for b in $names; do
    i=0
    while [ $i -lt 5 ]; do
      if delete_bucket "$b"; then break; fi
      i=$((i+1))
      sleep 5
    done
  done
else
  delete_bucket "${self.input.bucket_name}" || true
fi

exit 0
EOT
  }
  depends_on = [nutanix_object_store_v2.test]
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

  pcSSHPassword    = local.config.ssh_pc_password
  username         = local.config.ssh_pc_username
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
      value = local.objectStore.storage_network_dns_ip[0]
    }
  }
  storage_network_vip {
    ipv4 {
      value = local.objectStore.storage_network_vip[0]
    }
  }
  depends_on = [terraform_data.pre_update_hook]
}

// Ensure bucket is deleted before the object store is destroyed, even if a test step fails.
resource "terraform_data" "cleanup_bucket_test" {
  input = {
    object_store_id = nutanix_object_store_v2.test.id
    bucket_name     = local.objectStore.bucket_name
  }
  provisioner "local-exec" {
    when = destroy
    command = <<EOT
set -eu
URL="https://$NUTANIX_ENDPOINT:$NUTANIX_PORT/oss/api/nutanix/v3/objectstore_proxy/${self.input.object_store_id}/buckets/${self.input.bucket_name}?force_empty_bucket=true"
CODE="$(curl -sSk -u "$NUTANIX_USERNAME:$NUTANIX_PASSWORD" -X DELETE "$URL" -o /tmp/os_bucket_delete_test.out -w "%%%%{http_code}" || echo "000")"
if [ "$CODE" != "200" ] && [ "$CODE" != "202" ] && [ "$CODE" != "204" ] && [ "$CODE" != "404" ] && [ "$CODE" != "503" ]; then
  echo "bucket delete failed (test) http_code=$CODE url=$URL"
  cat /tmp/os_bucket_delete_test.out || true
  exit 1
fi
exit 0
EOT
  }
  depends_on = [nutanix_object_store_v2.test]
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

func testAccObjectStoreObjectLiteSourceConfig(objectLiteSourceImgName, vmName, vmOvaName, objectOvaName string) string {
	nutanixUsername := os.Getenv("NUTANIX_USERNAME")
	nutanixPassword := os.Getenv("NUTANIX_PASSWORD")
	nutanixEndpoint := os.Getenv("NUTANIX_ENDPOINT")
	nutanixPort := os.Getenv("NUTANIX_PORT")
	return fmt.Sprintf(`

locals {
			        # nutanix_username:nutanix_password
  aws_access_key  = base64encode("%[1]s:%[2]s")
  aws_secret_key  = base64encode("%[1]s:%[2]s")
  objects_API     = "https://%[3]s:%[4]s/api/prism/v4.0/objects/"
  disk_image_dest = "${path.module}/CentOS-7-cloudinit-os-img.qcow2"
}


# Configure AWS CLI access key
resource "terraform_data" "config_aws_access_key" {
  provisioner "local-exec" {
    when       = create
    command    = "aws configure set aws_access_key_id ${local.aws_access_key}"
    on_failure = fail
  }
}

# Configure AWS CLI secret key
resource "terraform_data" "config_aws_secret_key" {
  provisioner "local-exec" {
    when       = create
    command    = "aws configure set aws_secret_access_key ${local.aws_secret_key}"
    on_failure = fail
  }
  depends_on = [terraform_data.config_aws_access_key]
}

# Configure AWS Endpoint URL
resource "terraform_data" "config_aws_endpoint_url" {
  provisioner "local-exec" {
    when       = create
    command    = "aws configure set endpoint_url ${local.objects_API}"
    on_failure = fail
  }
  depends_on = [terraform_data.config_aws_secret_key]
}

# Download disk image from remote URL
resource "terraform_data" "download_disk_image" {
  provisioner "local-exec" {
    when       = create
    command    = "curl -o ${local.disk_image_dest} ${local.objectStore.img_url}"
    on_failure = fail
  }
  depends_on = [terraform_data.config_aws_endpoint_url]
}

# Upload image to object store bucket using AWS CLI
resource "terraform_data" "upload_image_to_object_store" {
  provisioner "local-exec" {
    when       = create
    command    = "aws s3api put-object --bucket vmm-images --body ${local.disk_image_dest} --key ${nutanix_object_store_v2.test.name} --no-verify-ssl"
    on_failure = fail
  }
  depends_on = [terraform_data.download_disk_image]
}

# Create image using object lite source
resource "nutanix_images_v2" "object-liteStore-img" {
  name        = "%[5]s"
  description = "Image created from object store"
  type        = "DISK_IMAGE"
  source {
    object_lite_source {
      key = nutanix_object_store_v2.test.name
    }
  }
  lifecycle {
    ignore_changes = [
      source
    ]
  }
  depends_on = [terraform_data.upload_image_to_object_store]
}


# Create VM with some specific requirements
resource "nutanix_virtual_machine_v2" "vm-test" {
  name              = "%[6]s"
  description       = "terraform test object store vm"
  num_sockets       = 2
  memory_size_bytes = 4 * 1024 * 1024 * 1024
  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
}

# Create Ova from the VM
resource "nutanix_ova_v2" "vm-ova" {
  name = "%[7]s"
  source {
    ova_vm_source {
      vm_ext_id        = nutanix_virtual_machine_v2.vm-test.id
      disk_file_format = "QCOW2"
    }
  }
}

# Download Ova
resource "nutanix_ova_download_v2" "test" {
  ova_ext_id = nutanix_ova_v2.vm-ova.id
}


# Upload Ova to object store using AWS CLI
resource "terraform_data" "upload_ova_to_object_store" {
  provisioner "local-exec" {
    when       = create
    command    = "aws s3api put-object --bucket vmm-ovas --body ${nutanix_ova_download_v2.test.ova_file_path} --key ${nutanix_object_store_v2.test.name} --no-verify-ssl"
    on_failure = fail
  }
}

# Sleep 1 min before uploading ova
resource "terraform_data" "delay" {
  provisioner "local-exec" {
    when       = create
    command    = "sleep 60"
    on_failure = fail
  }
  depends_on = [terraform_data.upload_ova_to_object_store]
}

# Create ova using object store source
resource "nutanix_ova_v2" "object-liteSource-ova" {
  name = "tf-object-ova"
  source {
    object_lite_source {
      key = nutanix_object_store_v2.test.name
    }
  }
  cluster_location_ext_ids = [local.clusterExtId]
  depends_on               = [terraform_data.delay]
}

# Download Ova from object store
resource "nutanix_ova_download_v2" "test" {
  ova_ext_id = nutanix_ova_v2.object-liteSource-ova.id
}

`, nutanixUsername, nutanixPassword, nutanixEndpoint, nutanixPort, objectLiteSourceImgName, vmName, vmOvaName, objectOvaName)
}
