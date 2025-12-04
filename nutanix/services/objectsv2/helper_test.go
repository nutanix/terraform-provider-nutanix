package objectstoresv2_test

import (
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

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_object_store_v2" {
			continue
		}

		bucketResp, bucketErr := deleteBucketForObjectStore(rs.Primary.ID)
		if bucketErr != nil {
			log.Printf("[ERROR] error deleting bucket: %v", bucketErr)
			return bucketErr
		}
		if bucketResp.StatusCode != http.StatusOK && bucketResp.StatusCode != 503 {
			return fmt.Errorf("error deleting bucket: %s", bucketResp.Status)
		}

		defer bucketResp.Body.Close()
		log.Println("[DEBUG] Bucket Deleted")

		// Check if the object store is deleted
		_, err := conn.ObjectStoreAPI.ObjectStoresAPIInstance.GetObjectstoreById(utils.StringPtr(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("object store still exists: %s", rs.Primary.ID)
		}
		if strings.Contains(err.Error(), "not found") {
			log.Println("[DEBUG] Object store deleted")
		}
		return nil
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
		log.Fatalf("creating request: %v", err)
		return nil, err
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
		log.Fatalf("delete failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %s\nBody: %s\n", resp.Status, body)

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
				return fmt.Errorf("error deleting bucket: %s", resp.Status)
			}
			defer resp.Body.Close()
			log.Println("[DEBUG] Bucket Deleted")

			return nil
		}
		return nil
	}
}
