package multidomain

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	multidomainClient "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

const ncTaskStatusPathFormat = "/api/nutanix-central/v4.0.a1/config/tasks/%s/status"

// NCTaskClient provides Nutanix Central task-status APIs that are not yet
// available through the generated multidomain SDK.
type NCTaskClient struct {
	apiClient   *multidomainClient.ApiClient
	credentials client.Credentials
	httpClient  *http.Client
}

// NCTaskStatusResponse represents the Nutanix Central task-status response.
type NCTaskStatusResponse struct {
	Status string `json:"status"`
	Data   struct {
		SubStateInfo struct {
			CurrentMilestone   string `json:"currentMilestone"`
			CurrentMilestoneID int    `json:"currentMilestoneID"`
			ErrorCode          string `json:"errorCode"`
			ErrorDetail        string `json:"errorDetail"`
			OperationType      string `json:"operationType"`
			PercentageComplete int    `json:"percentageComplete"`
			Status             string `json:"status"`
		} `json:"subStateInfo"`
	} `json:"data"`
}

// NewNCTaskClient creates a client for NC task-status polling.
func NewNCTaskClient(apiClient *multidomainClient.ApiClient, credentials client.Credentials) *NCTaskClient {
	verifySSL := false
	if apiClient != nil {
		verifySSL = apiClient.VerifySSL
	}

	return &NCTaskClient{
		apiClient:   apiClient,
		credentials: credentials,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifySSL}, //nolint:gosec // Matches the configured SDK TLS behavior.
			},
		},
	}
}

// GetStatus gets the current status of a Nutanix Central asynchronous task.
func (c *NCTaskClient) GetStatus(ctx context.Context, taskUUID string) (*NCTaskStatusResponse, error) {
	if c.apiClient == nil {
		return nil, fmt.Errorf("Nutanix Central task client is not configured")
	}
	if taskUUID == "" {
		return nil, fmt.Errorf("task UUID must not be empty")
	}

	scheme := c.apiClient.Scheme
	if scheme == "" {
		scheme = "https"
	}
	endpoint := (&url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(c.apiClient.Host, fmt.Sprintf("%d", c.apiClient.Port)),
		Path:   fmt.Sprintf(ncTaskStatusPathFormat, url.PathEscape(taskUUID)),
	}).String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating task-status request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	for key, value := range c.credentials.CustomHeaders {
		req.Header.Set(key, value)
	}
	if c.credentials.APIKey != "" {
		req.Header.Set("X-Ntnx-Api-Key", c.credentials.APIKey)
	} else {
		req.SetBasicAuth(c.credentials.Username, c.credentials.Password)
	}

	log.Printf("[DEBUG] GET Nutanix Central task status: %s", endpoint)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getting task status for %q: %w", taskUUID, err)
	}
	defer resp.Body.Close()
	log.Printf("[DEBUG] Nutanix Central task status response for %q: %s", taskUUID, resp.Status)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("getting task status for %q: unexpected HTTP status %s", taskUUID, resp.Status)
	}

	var taskStatus NCTaskStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskStatus); err != nil {
		return nil, fmt.Errorf("decoding task status for %q: %w", taskUUID, err)
	}
	if taskStatus.State() == "" {
		return nil, fmt.Errorf("task-status response for %q did not include a status", taskUUID)
	}
	log.Printf("[DEBUG] Nutanix Central task %q state: %s", taskUUID, taskStatus.State())
	return &taskStatus, nil
}

// State returns the detailed task status, falling back to the top-level status.
func (r *NCTaskStatusResponse) State() string {
	if r.Data.SubStateInfo.Status != "" {
		return r.Data.SubStateInfo.Status
	}
	return r.Status
}

// IsTerminal reports whether this task has completed successfully or unsuccessfully.
func (r *NCTaskStatusResponse) IsTerminal() bool {
	switch strings.ToUpper(r.State()) {
	case "SUCCEEDED", "COMPLETED", "FAILED", "CANCELED", "CANCELLED":
		return true
	default:
		return false
	}
}

// WaitForTerminalState polls until the task reaches a terminal state or the
// caller's context expires. A non-positive interval defaults to ten seconds.
func (c *NCTaskClient) WaitForTerminalState(ctx context.Context, taskUUID string, interval time.Duration) (*NCTaskStatusResponse, error) {
	if interval <= 0 {
		interval = 10 * time.Second
	}

	for {
		taskStatus, err := c.GetStatus(ctx, taskUUID)
		if err != nil {
			return nil, err
		}
		if taskStatus.IsTerminal() {
			return taskStatus, nil
		}

		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, fmt.Errorf("waiting for task %q: %w", taskUUID, ctx.Err())
		case <-timer.C:
		}
	}
}
