package objectstoresv2_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func testAccCheckNutanixObjectStoreDestroy(s *terraform.State) error {
	log.Println("Checking Object store destroy")
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	// loop through all the resources and delete the object store and bucket
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_object_store_v2" {
			continue
		}

		// Attempt to delete the bucket associated with the object store
		bucketResp, bucketErr := deleteBucketForObjectStore(rs.Primary.ID)
		if bucketErr != nil {
			log.Printf("[ERROR] error deleting bucket for ObjectStore %s: %v", rs.Primary.ID, bucketErr)
			// Return the first error encountered
			return bucketErr
		}
		if bucketResp != nil {
			defer bucketResp.Body.Close()
			if bucketResp.StatusCode != http.StatusOK &&
				bucketResp.StatusCode != http.StatusAccepted &&
				bucketResp.StatusCode != http.StatusNoContent &&
				bucketResp.StatusCode != http.StatusNotFound &&
				bucketResp.StatusCode != http.StatusServiceUnavailable {
				return fmt.Errorf("error deleting bucket for ObjectStore %s: %s", rs.Primary.ID, bucketResp.Status)
			}
			log.Printf("[DEBUG] Bucket for ObjectStore %s deleted (status %d)", rs.Primary.ID, bucketResp.StatusCode)
		}

		// Now check if the object store itself is deleted
		objStore, err := conn.ObjectStoreAPI.ObjectStoresAPIInstance.GetObjectstoreById(utils.StringPtr(rs.Primary.ID))
		if err == nil && objStore != nil {
			return fmt.Errorf("object store still exists: %s", rs.Primary.ID)
		}
		if err != nil && strings.Contains(err.Error(), "not found") {
			log.Printf("[DEBUG] Object store %s deleted", rs.Primary.ID)
		}
		// else: ignore other errors
	}
	return nil
}

func deleteBucketForObjectStore(objectStoreExtID string) (*http.Response, error) {
	// 1) Prepare context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	endpoint := os.Getenv("NUTANIX_ENDPOINT")
	user := os.Getenv("NUTANIX_USERNAME")
	password := os.Getenv("NUTANIX_PASSWORD")
	port := os.Getenv("NUTANIX_PORT")

	bucketName := testVars.ObjectStore.BucketName

	// 2) Prepare the URL
	url := fmt.Sprintf("https://%s:%s/oss/api/nutanix/v3/objectstore_proxy/%s/buckets/%s?force_empty_bucket=true", endpoint, port, objectStoreExtID, bucketName)

	// 3) Create the DELETE request
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// 4) Set Authentication
	req.SetBasicAuth(user, password)

	// 5) Set Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 6) custom TLS transport
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 10 * time.Second,
	}

	// 7) Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("delete failed: %w", err)
	}

	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	// re-wrap body so callers can still inspect it and close it
	resp.Body = io.NopCloser(bytes.NewReader(body))
	log.Printf("[DEBUG] delete bucket response: status=%s body=%s", resp.Status, string(body))

	return resp, nil
}

func deleteObjectStoreBucket() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "nutanix_object_store_v2" {
				continue
			}

			// get the object store ID
			objectStoreID := rs.Primary.ID

			// Delete the object store bucket
			resp, err := deleteBucketForObjectStore(objectStoreID)
			if err != nil {
				return fmt.Errorf("error deleting bucket: %s", err)
			}
			if resp.StatusCode != http.StatusAccepted {
				// allow typical success codes across different PC versions
				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
					return fmt.Errorf("error deleting bucket: %s", resp.Status)
				}
			}
			defer resp.Body.Close()
			log.Println("[DEBUG] Bucket Deleted")

			return nil
		}
		return nil
	}
}
