package iamv2_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

type TestConfig struct {
	Iam struct {
		Roles struct {
			DisplayName string `json:"display_name"`
			Description string `json:"description"`
		} `json:"roles"`
		Users struct {
			Name                        string `json:"name"`
			IdpID                       string `json:"idp_id"`
			DirectoryServiceID          string `json:"directory_service_id"`
			DirectoryServiceUsername    string `json:"directory_service_username"`
			EmailID                     string `json:"email_id"`
			Locale                      string `json:"locale"`
			Region                      string `json:"region"`
			Password                    string `json:"password"`
			IsForceResetPasswordEnabled bool   `json:"is_force_reset_password"`
		} `json:"users"`
		// UserGroups config
		UserGroups struct {
			Name              string `json:"name"`
			SAMLName          string `json:"saml_name"`
			DistinguishedName string `json:"distinguished_name"`
		} `json:"user_groups"`
		AuthPolicies struct {
			DisplayName    string   `json:"display_name"`
			Description    string   `json:"description"`
			AuthPolicyType string   `json:"authorization_policy_type"`
			Identities     []string `json:"identities"`
			Entities       []string `json:"entities"`
		} `json:"auth_policies"`
		// Directory Services config
		IdentityProviders struct {
			IdpMetadataURL string `json:"idp_metadata_url"`
			IdpMetadata    struct {
				EntityID           string `json:"entity_id"`
				LoginURL           string `json:"login_url"`
				LogoutURL          string `json:"logout_url"`
				ErrorURL           string `json:"error_url"`
				Certificate        string `json:"certificate"`
				NameIDPolicyFormat string `json:"name_id_policy_format"`
			} `json:"idp_metadata"`
			IdpMetadataXML          string   `json:"idp_metadata_xml"`
			Name                    string   `json:"name"`
			UsernameAttribute       string   `json:"username_attr"`
			EmailAttribute          string   `json:"email_attr"`
			GroupsAttribute         string   `json:"groups_attr"`
			GroupsDelim             string   `json:"groups_delim"`
			EntityIssuer            string   `json:"entity_issuer"`
			CustomAttributes        []string `json:"custom_attributes"`
			IsSignedAuthnReqEnabled bool     `json:"is_signed_authn_req_enabled"`
		} `json:"identity_providers"`
		// Directory Services config
		DirectoryServices struct {
			Name           string   `json:"name"`
			URL            string   `json:"url"`
			SecondaryUrls  []string `json:"secondary_urls"`
			DomainName     string   `json:"domain_name"`
			ServiceAccount struct {
				Username string `json:"username"`
				Password string `json:"password"`
			} `json:"service_account"`
			OpenLdapConfiguration struct {
				UserConfiguration struct {
					UserObjectClass   string `json:"user_object_class"`
					UserSearchBase    string `json:"user_search_base"`
					UsernameAttribute string `json:"username_attribute"`
				} `json:"user_configuration"`
				UserGroupConfiguration struct {
					GroupObjectClass          string `json:"group_object_class"`
					GroupSearchBase           string `json:"group_search_base"`
					GroupMemberAttribute      string `json:"group_member_attribute"`
					GroupMemberAttributeValue string `json:"group_member_attribute_value"`
				} `json:"user_group_configuration"`
			} `json:"open_ldap_configuration"`

			GroupSearchType   string   `json:"group_search_type"`
			WhiteListedGroups []string `json:"white_listed_groups"`
		} `json:"directory_services"`
	} `json:"iam"`
}

var testVars TestConfig

var (
	path, _     = os.Getwd()
	filepath    = path + "/../../../test_config_v2.json"
	xmlFilePath = path + "/../../../test_idp_metadata.txt"
)

func loadVars(filepath string, varStuct interface{}) {
	// Read test_config_v2.json from home current path
	configData, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Got this error while reading config.json: %s", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(configData, varStuct)
	if err != nil {
		log.Printf("Got this error while unmarshalling config.json: %s", err.Error())
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	log.Println("Do some crazy stuff before tests!")
	loadVars("../../../test_config_v2.json", &testVars)
	downloadFile(testVars.Iam.IdentityProviders.IdpMetadataXML, xmlFilePath)
	os.Exit(m.Run())
}

// downloadFile downloads a file from a given URL and saves it to the specified path.
func downloadFile(url, destinationFilePath string) error {
	// Send HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 responses
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the destination file
	out, err := os.Create(destinationFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy data from response to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	return nil
}
