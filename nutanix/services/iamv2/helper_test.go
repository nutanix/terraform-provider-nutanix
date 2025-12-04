package iamv2_test

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
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

func testAccCheckNutanixUserDestroy(s *terraform.State) error {
	log.Println("Checking user destroy")
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_users_v2" || rs.Primary.Attributes["users.#"] != "" {
			continue
		}

		if _, err := conn.API.V3.GetUser(rs.Primary.ID); err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return nil
			}
			return err
		}
		_, err := conn.API.V3.DeleteUser(rs.Primary.ID)
		if err != nil {
			return err
		}
		log.Println("User Deleted")
	}
	return nil
}

func testAccCheckNutanixDirectoryServicesV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_directory_services_v2" {
			continue
		}
		fmt.Printf("Checking directory service : %s", rs.Primary.ID)
		readResp, errRead := conn.IamAPI.DirectoryServiceAPIInstance.GetDirectoryServiceById(utils.StringPtr(rs.Primary.ID))
		if errRead != nil {
			if strings.Contains(fmt.Sprint(errRead), "Directory service not found") {
				return nil
			}
			return errRead
		}
		// get etag value from read response to pass in update request If-Match header, Required for update request
		etagValue := conn.IamAPI.DirectoryServiceAPIInstance.ApiClient.GetEtag(readResp)
		headers := make(map[string]interface{})
		headers["If-Match"] = utils.StringPtr(etagValue)

		fmt.Println("Deleting directory service")

		if _, err := conn.IamAPI.DirectoryServiceAPIInstance.DeleteDirectoryServiceById(utils.StringPtr(rs.Primary.ID), headers); err != nil {
			if strings.Contains(fmt.Sprint(err), "Directory service not found") {
				return nil
			}
			return err
		}
	}
	return nil
}

func testAccCheckNutanixUserGroupsV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_user_groups_v2" {
			continue
		}

		readResp, errRead := conn.IamAPI.UserGroupsAPIInstance.GetUserGroupById(utils.StringPtr(rs.Primary.ID))
		if errRead != nil {
			if strings.Contains(fmt.Sprint(errRead), "the requested user group does not exist") {
				return nil
			}
			return errRead
		}

		// get etag value from read response to pass in update request If-Match header, Required for update request
		etagValue := conn.IamAPI.DirectoryServiceAPIInstance.ApiClient.GetEtag(readResp)
		headers := make(map[string]interface{})
		headers["If-Match"] = utils.StringPtr(etagValue)

		if _, err := conn.IamAPI.UserGroupsAPIInstance.DeleteUserGroupById(utils.StringPtr(rs.Primary.ID), headers); err != nil {
			if strings.Contains(fmt.Sprint(err), "the requested user group does not exist") {
				return nil
			}
			return err
		}
	}
	return nil
}
