package prismv2

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	prismapi "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/api"
	prismclient "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/client"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// func to flatten the time to string
func flattenTime(time *time.Time) *string {
	if time == nil {
		return nil
	}
	return utils.StringPtr(time.String())
}

// schemas for links
func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// createRestoreSourceWithV42Fallback first uses the generated SDK request.
// If it fails with 404 (common in mixed-version PC/PE setups), it retries via explicit v4.2 endpoint.
func createRestoreSourceWithV42Fallback(
	ctx context.Context,
	api *prismapi.DomainManagerBackupsApi,
	body *management.RestoreSource,
) (*management.CreateRestoreSourceApiResponse, error) {
	resp, err := api.CreateRestoreSource(body)
	if err == nil || !isNotFoundError(err) {
		return resp, err
	}

	log.Printf("[WARN] CreateRestoreSource v4.3 failed with 404; retrying via v4.2 endpoint")
	uri := "/api/prism/v4.2/management/restore-sources"
	var out management.CreateRestoreSourceApiResponse

	err = callDomainManagerBackupsAPI(ctx, api.ApiClient, http.MethodPost, uri, body, url.Values{}, map[string]string{}, &out)
	if err != nil {
		log.Printf("[ERROR] CreateRestoreSource fallback v4.2 failed: %v", err)
		return nil, err
	}
	log.Printf("[DEBUG] CreateRestoreSource fallback v4.2 succeeded")
	return &out, nil
}

func getRestoreSourceByIDWithV42Fallback(
	ctx context.Context,
	api *prismapi.DomainManagerBackupsApi,
	extID *string,
	args ...map[string]interface{},
) (*management.GetRestoreSourceApiResponse, error) {
	resp, err := api.GetRestoreSourceById(extID, args...)
	if err == nil || !isNotFoundError(err) {
		return resp, err
	}

	log.Printf("[WARN] GetRestoreSourceById v4.3 failed with 404; retrying via v4.2 endpoint")
	uri := "/api/prism/v4.2/management/restore-sources/{extId}"
	uri = strings.Replace(uri, "{extId}", url.PathEscape(*extID), 1)
	var out management.GetRestoreSourceApiResponse

	err = callDomainManagerBackupsAPI(
		ctx,
		api.ApiClient,
		http.MethodGet,
		uri,
		nil,
		url.Values{},
		extractHeaderParams(args...),
		&out,
	)
	if err != nil {
		log.Printf("[ERROR] GetRestoreSourceById fallback v4.2 failed: %v", err)
		return nil, err
	}
	log.Printf("[DEBUG] GetRestoreSourceById fallback v4.2 succeeded")
	return &out, nil
}

func deleteRestoreSourceByIDWithV42Fallback(
	ctx context.Context,
	api *prismapi.DomainManagerBackupsApi,
	extID *string,
	args map[string]interface{},
) (*management.DeleteRestoreSourceApiResponse, error) {
	resp, err := api.DeleteRestoreSourceById(extID, args)
	if err == nil || !isNotFoundError(err) {
		return resp, err
	}

	log.Printf("[WARN] DeleteRestoreSourceById v4.3 failed with 404; retrying via v4.2 endpoint")
	uri := "/api/prism/v4.2/management/restore-sources/{extId}"
	uri = strings.Replace(uri, "{extId}", url.PathEscape(*extID), 1)
	var out management.DeleteRestoreSourceApiResponse

	err = callDomainManagerBackupsAPI(ctx, api.ApiClient, http.MethodDelete, uri, nil, url.Values{}, extractHeaderParams(args), &out)
	if err != nil {
		log.Printf("[ERROR] DeleteRestoreSourceById fallback v4.2 failed: %v", err)
		return nil, err
	}
	log.Printf("[DEBUG] DeleteRestoreSourceById fallback v4.2 succeeded")
	return &out, nil
}

func listRestorableDomainManagersWithV42Fallback(
	ctx context.Context,
	api *prismapi.DomainManagerBackupsApi,
	restoreSourceExtID *string,
	page, limit *int,
	filter *string,
	args ...map[string]interface{},
) (*management.ListRestorableDomainManagersApiResponse, error) {
	resp, err := api.ListRestorableDomainManagers(restoreSourceExtID, page, limit, filter, args...)
	if err == nil || !isNotFoundError(err) {
		return resp, err
	}

	log.Printf("[WARN] ListRestorableDomainManagers v4.3 failed with 404; retrying via v4.2 endpoint")
	uri := "/api/prism/v4.2/management/restore-sources/{restoreSourceExtId}/restorable-domain-managers"
	uri = strings.Replace(uri, "{restoreSourceExtId}", url.PathEscape(*restoreSourceExtID), 1)
	queryParams := url.Values{}
	if page != nil {
		queryParams.Add("$page", prismclient.ParameterToString(*page, ""))
	}
	if limit != nil {
		queryParams.Add("$limit", prismclient.ParameterToString(*limit, ""))
	}
	if filter != nil {
		queryParams.Add("$filter", prismclient.ParameterToString(*filter, ""))
	}
	var out management.ListRestorableDomainManagersApiResponse

	err = callDomainManagerBackupsAPI(
		ctx,
		api.ApiClient,
		http.MethodGet,
		uri,
		nil,
		queryParams,
		extractHeaderParams(args...),
		&out,
	)
	if err != nil {
		log.Printf("[ERROR] ListRestorableDomainManagers fallback v4.2 failed: %v", err)
		return nil, err
	}
	log.Printf("[DEBUG] ListRestorableDomainManagers fallback v4.2 succeeded")
	return &out, nil
}

func getRestorePointByIDWithV42Fallback(
	ctx context.Context,
	api *prismapi.DomainManagerBackupsApi,
	restoreSourceExtID, restorableDomainManagerExtID, extID *string,
	args ...map[string]interface{},
) (*management.GetRestorePointApiResponse, error) {
	resp, err := api.GetRestorePointById(restoreSourceExtID, restorableDomainManagerExtID, extID, args...)
	if err == nil || !isNotFoundError(err) {
		return resp, err
	}

	log.Printf("[WARN] GetRestorePointById v4.3 failed with 404; retrying via v4.2 endpoint")
	uri := "/api/prism/v4.2/management/restore-sources/{restoreSourceExtId}/restorable-domain-managers/{restorableDomainManagerExtId}/restore-points/{extId}"
	uri = strings.Replace(uri, "{restoreSourceExtId}", url.PathEscape(*restoreSourceExtID), 1)
	uri = strings.Replace(uri, "{restorableDomainManagerExtId}", url.PathEscape(*restorableDomainManagerExtID), 1)
	uri = strings.Replace(uri, "{extId}", url.PathEscape(*extID), 1)
	var out management.GetRestorePointApiResponse

	err = callDomainManagerBackupsAPI(
		ctx,
		api.ApiClient,
		http.MethodGet,
		uri,
		nil,
		url.Values{},
		extractHeaderParams(args...),
		&out,
	)
	if err != nil {
		log.Printf("[ERROR] GetRestorePointById fallback v4.2 failed: %v", err)
		return nil, err
	}
	log.Printf("[DEBUG] GetRestorePointById fallback v4.2 succeeded")
	return &out, nil
}

func listRestorePointsWithV42Fallback(
	ctx context.Context,
	api *prismapi.DomainManagerBackupsApi,
	restoreSourceExtID, restorableDomainManagerExtID *string,
	page, limit *int,
	filter, orderBy, selects *string,
	args ...map[string]interface{},
) (*management.ListRestorePointsApiResponse, error) {
	resp, err := api.ListRestorePoints(restoreSourceExtID, restorableDomainManagerExtID, page, limit, filter, orderBy, selects, args...)
	if err == nil || !isNotFoundError(err) {
		return resp, err
	}

	log.Printf("[WARN] ListRestorePoints v4.3 failed with 404; retrying via v4.2 endpoint")
	uri := "/api/prism/v4.2/management/restore-sources/{restoreSourceExtId}/restorable-domain-managers/{restorableDomainManagerExtId}/restore-points"
	uri = strings.Replace(uri, "{restoreSourceExtId}", url.PathEscape(*restoreSourceExtID), 1)
	uri = strings.Replace(uri, "{restorableDomainManagerExtId}", url.PathEscape(*restorableDomainManagerExtID), 1)
	queryParams := url.Values{}
	if page != nil {
		queryParams.Add("$page", prismclient.ParameterToString(*page, ""))
	}
	if limit != nil {
		queryParams.Add("$limit", prismclient.ParameterToString(*limit, ""))
	}
	if filter != nil {
		queryParams.Add("$filter", prismclient.ParameterToString(*filter, ""))
	}
	if orderBy != nil {
		queryParams.Add("$orderby", prismclient.ParameterToString(*orderBy, ""))
	}
	if selects != nil {
		queryParams.Add("$select", prismclient.ParameterToString(*selects, ""))
	}
	var out management.ListRestorePointsApiResponse

	err = callDomainManagerBackupsAPI(
		ctx,
		api.ApiClient,
		http.MethodGet,
		uri,
		nil,
		queryParams,
		extractHeaderParams(args...),
		&out,
	)
	if err != nil {
		log.Printf("[ERROR] ListRestorePoints fallback v4.2 failed: %v", err)
		return nil, err
	}
	log.Printf("[DEBUG] ListRestorePoints fallback v4.2 succeeded")
	return &out, nil
}

func restoreDomainManagerWithV42Fallback(
	ctx context.Context,
	api *prismapi.DomainManagerBackupsApi,
	restoreSourceExtID, restorableDomainManagerExtID, restorePointExtID *string,
	body *management.RestoreSpec,
	args ...map[string]interface{},
) (*management.RestoreApiResponse, error) {
	resp, err := api.Restore(restoreSourceExtID, restorableDomainManagerExtID, restorePointExtID, body, args...)
	if err == nil || !isNotFoundError(err) {
		return resp, err
	}

	log.Printf("[WARN] Restore v4.3 failed with 404; retrying via v4.2 endpoint")
	uri := "/api/prism/v4.2/management/restore-sources/{restoreSourceExtId}/restorable-domain-managers/{restorableDomainManagerExtId}/restore-points/{extId}/$actions/restore"
	uri = strings.Replace(uri, "{restoreSourceExtId}", url.PathEscape(*restoreSourceExtID), 1)
	uri = strings.Replace(uri, "{restorableDomainManagerExtId}", url.PathEscape(*restorableDomainManagerExtID), 1)
	uri = strings.Replace(uri, "{extId}", url.PathEscape(*restorePointExtID), 1)
	var out management.RestoreApiResponse

	err = callDomainManagerBackupsAPI(
		ctx,
		api.ApiClient,
		http.MethodPost,
		uri,
		body,
		url.Values{},
		extractHeaderParams(args...),
		&out,
	)
	if err != nil {
		log.Printf("[ERROR] Restore fallback v4.2 failed: %v", err)
		return nil, err
	}
	log.Printf("[DEBUG] Restore fallback v4.2 succeeded")
	return &out, nil
}

func callDomainManagerBackupsAPI(
	ctx context.Context,
	apiClient *prismclient.ApiClient,
	method, uri string,
	body interface{},
	queryParams url.Values,
	headerParams map[string]string,
	out interface{},
) error {
	contentTypes := []string{}
	if method == http.MethodPost || method == http.MethodPut {
		contentTypes = []string{"application/json"}
	}

	formParams := url.Values{}
	accepts := []string{"application/json"}
	authNames := []string{"apiKeyAuthScheme", "basicAuthScheme"}

	apiClientResponse, err := apiClient.CallApiWithContext(
		ctx,
		&uri,
		method,
		body,
		queryParams,
		headerParams,
		formParams,
		accepts,
		contentTypes,
		authNames,
	)
	if err != nil || apiClientResponse == nil {
		return err
	}

	return json.Unmarshal(apiClientResponse.([]byte), out)
}

func extractHeaderParams(args ...map[string]interface{}) map[string]string {
	headers := map[string]string{}
	if len(args) == 0 || args[0] == nil {
		return headers
	}

	for key, val := range args[0] {
		if headerVal, ok := val.(*string); ok && headerVal != nil {
			headers[key] = *headerVal
		}
	}

	return headers
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	openAPIErr, ok := err.(prismclient.GenericOpenAPIError)
	if !ok {
		return false
	}

	return strings.Contains(strings.ToUpper(openAPIErr.Status), "404")
}
