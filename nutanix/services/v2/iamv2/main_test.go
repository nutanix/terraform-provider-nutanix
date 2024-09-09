package iamv2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	Iam struct {
		Roles struct {
			Limit       int    `json:"limit"`
			DisplayName string `json:"display_name"`
			Description string `json:"description"`
		} `json:"roles"`
		Users struct {
			Limit                       int    `json:"limit"`
			Username                    string `json:"username"`
			IdpId                       string `json:"idp_id"`
			DirectoryServiceId          string `json:"directory_service_id"`
			DirectoryServiceUsername    string `json:"directory_service_username"`
			DisplayName                 string `json:"display_name"`
			FirstName                   string `json:"first_name"`
			MiddleInitial               string `json:"middle_initial"`
			LastName                    string `json:"last_name"`
			EmailId                     string `json:"email_id"`
			Locale                      string `json:"locale"`
			Region                      string `json:"region"`
			Password                    string `json:"password"`
			IsForceResetPasswordEnabled bool   `json:"is_force_reset_password"`
		} `json:"users"`
		// UserGroups config
		UserGroups struct {
			Limit              int    `json:"limit"`
			IdpId              string `json:"idp_id"`
			DirectoryServiceId string `json:"directory_service_id"`
			Name               string `json:"name"`
			SAMLName           string `json:"saml_name"`
			DistinguishedName  string `json:"distinguished_name"`
		} `json:"user_groups"`
		AuthPolicies struct {
			Limit          int      `json:"limit"`
			DisplayName    string   `json:"display_name"`
			Description    string   `json:"description"`
			AuthPolicyType string   `json:"authorization_policy_type"`
			Identities     []string `json:"identities"`
			Entities       []string `json:"entities"`
		} `json:"auth_policies"`
		// Directory Services config
		IdentityProviders struct {
			Limit          int    `json:"limit"`
			IdpMetadataUrl string `json:"idp_metadata_url"`
			IdpMetadata    struct {
				EntityId           string `json:"entity_id"`
				LoginUrl           string `json:"login_url"`
				LogoutUrl          string `json:"logout_url"`
				ErrorUrl           string `json:"error_url"`
				Certificate        string `json:"certificate"`
				NameIdPolicyFormat string `json:"name_id_policy_format"`
			} `json:"idp_metadata"`
			IdpMetadataXml          string   `json:"idp_metadata_xml"`
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
			Limit          int      `json:"limit"`
			Name           string   `json:"name"`
			Url            string   `json:"url"`
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
	loadVars("../../../../test_config_v2.json", &testVars)
	os.Exit(m.Run())
}
