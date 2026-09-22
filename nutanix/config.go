package nutanix

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	foundation_central "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/fc"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/foundation"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/karbon"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/selfservice"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/clusters"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/datapolicies"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/dataprotection"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/iam"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/lcm"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/microseg"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/multidomain"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/networking"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/objectstores"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/security"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/vmm"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/volumes"
)

// Version represents api version
// const Version = "3.1"

// Config ...
type Config struct {
	Endpoint           string
	Username           string
	Password           string
	Port               string
	Insecure           bool
	SessionAuth        bool
	WaitTimeout        int64
	ProxyURL           string
	FoundationEndpoint string              // Required field for connecting to foundation VM APIs
	FoundationPort     string              // Port for connecting to foundation VM APIs
	RequiredFields     map[string][]string // RequiredFields is client name to its required fields mapping for validations and usage in every client
	NdbEndpoint        string
	NdbUsername        string
	NdbPassword        string
	APIKey             string            // API key for authentication (alternative to username/password)
	CustomHeaders      map[string]string // Custom headers to add to all requests (e.g., for Cloudflare Access)
}

// clientCache memoizes fully constructed *Client instances for the lifetime of
// the process, keyed by the connection parameters in Config. Terraform core
// calls the provider's ConfigureContextFunc (which ends up here) once per graph
// walk - plan, apply, refresh and destroy, and once per test step - so without
// caching every one of those calls rebuilds all ~19 SDK sub-clients from
// scratch. Each fresh v4 *ApiClient has negotiationCompleted=false, so it
// re-runs the lazy OPTIONS ".../unversioned/info" version negotiation against
// every service endpoint again. Reusing the same *Client keeps the negotiation
// state (negotiationCompleted) on the underlying *ApiClient, so version
// negotiation happens at most once per unique config per process instead of
// once per Configure.
var (
	clientCacheMu sync.Mutex
	clientCache   = map[string]*Client{}
)

// cacheKey builds a stable fingerprint of every Config field that influences
// how the SDK clients are constructed. Two Configs with the same key are
// interchangeable, so they can safely share one *Client.
func (c *Config) cacheKey() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\x00%s\x00%s\x00%s\x00%t\x00%t\x00%d\x00%s\x00%s\x00%s\x00%s\x00%s\x00%s\x00%s\x00",
		c.Endpoint, c.Username, c.Password, c.Port, c.Insecure, c.SessionAuth,
		c.WaitTimeout, c.ProxyURL, c.FoundationEndpoint, c.FoundationPort,
		c.NdbEndpoint, c.NdbUsername, c.NdbPassword, c.APIKey)

	// Custom headers, serialized in a deterministic (sorted) order.
	headerKeys := make([]string, 0, len(c.CustomHeaders))
	for k := range c.CustomHeaders {
		headerKeys = append(headerKeys, k)
	}
	sort.Strings(headerKeys)
	for _, k := range headerKeys {
		fmt.Fprintf(&b, "%s=%s\x01", k, c.CustomHeaders[k])
	}

	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

// Client ...
func (c *Config) Client() (*Client, error) {
	key := c.cacheKey()

	// Hold the lock across the whole construction. Construction itself is cheap
	// (no network I/O - v4 version negotiation is lazy and only fires on the
	// first real API call), and serializing it means concurrent Configure calls
	// build the client once and everyone else reuses it.
	clientCacheMu.Lock()
	defer clientCacheMu.Unlock()

	if cached, ok := clientCache[key]; ok {
		log.Printf("[DEBUG] reusing cached Nutanix API client; skipping SDK client construction and version negotiation")
		return cached, nil
	}

	configCreds := client.Credentials{
		URL:                client.JoinHostPort(c.Endpoint, c.Port),
		Endpoint:           c.Endpoint,
		Username:           c.Username,
		Password:           c.Password,
		Port:               c.Port,
		Insecure:           c.Insecure,
		SessionAuth:        c.SessionAuth,
		ProxyURL:           c.ProxyURL,
		FoundationEndpoint: c.FoundationEndpoint,
		FoundationPort:     c.FoundationPort,
		NdbEndpoint:        c.NdbEndpoint,
		NdbUsername:        c.NdbUsername,
		NdbPassword:        c.NdbPassword,
		RequiredFields:     c.RequiredFields,
		APIKey:             c.APIKey,
		CustomHeaders:      c.CustomHeaders,
	}

	v3Client, err := v3.NewV3Client(configCreds)
	if err != nil {
		return nil, err
	}
	karbonClient, err := karbon.NewKarbonAPIClient(configCreds)
	if err != nil {
		return nil, err
	}
	foundationClient, err := foundation.NewFoundationAPIClient(configCreds)
	if err != nil {
		return nil, err
	}
	fcClient, err := foundation_central.NewFoundationCentralClient(configCreds)
	if err != nil {
		return nil, err
	}
	eraClient, err := era.NewEraClient(configCreds)
	if err != nil {
		return nil, err
	}
	iamClient, err := iam.NewIamClient(configCreds)
	if err != nil {
		return nil, err
	}
	networkingClient, err := networking.NewNetworkingClient(configCreds)
	if err != nil {
		return nil, err
	}
	prismClient, err := prism.NewPrismClient(configCreds)
	if err != nil {
		return nil, err
	}
	microsegClient, err := microseg.NewMicrosegClient(configCreds)
	if err != nil {
		return nil, err
	}
	volumeClient, err := volumes.NewVolumeClient(configCreds)
	if err != nil {
		return nil, err
	}
	clustersClient, err := clusters.NewClustersClient(configCreds)
	if err != nil {
		return nil, err
	}
	dataprotectionClient, err := dataprotection.NewDataProtectionClient(configCreds)
	if err != nil {
		return nil, err
	}
	vmmClient, err := vmm.NewVmmClient(configCreds)
	if err != nil {
		return nil, err
	}
	dataPoliciesClient, err := datapolicies.NewDataPoliciesClient(configCreds)
	if err != nil {
		return nil, err
	}
	LcmClient, err := lcm.NewLcmClient(configCreds)
	if err != nil {
		return nil, err
	}
	calmClient, err := selfservice.NewCalmClient(configCreds)
	if err != nil {
		return nil, err
	}
	ObjectStoreClient, err := objectstores.NewObjectStoresClient(configCreds)
	if err != nil {
		return nil, err
	}
	SecurityClient, err := security.NewSecurityClient(configCreds)
	if err != nil {
		return nil, err
	}
	MultidomainClient, err := multidomain.NewMultidomainClient(configCreds)
	if err != nil {
		return nil, err
	}

	builtClient := &Client{
		WaitTimeout:         c.WaitTimeout,
		API:                 v3Client,
		KarbonAPI:           karbonClient,
		FoundationClientAPI: foundationClient,
		FoundationCentral:   fcClient,
		Era:                 eraClient,
		NetworkingAPI:       networkingClient,
		PrismAPI:            prismClient,
		MicroSegAPI:         microsegClient,
		IamAPI:              iamClient,
		ClusterAPI:          clustersClient,
		VolumeAPI:           volumeClient,
		DataProtectionAPI:   dataprotectionClient,
		VmmAPI:              vmmClient,
		DataPoliciesAPI:     dataPoliciesClient,
		LcmAPI:              LcmClient,
		CalmAPI:             calmClient,
		ObjectStoreAPI:      ObjectStoreClient,
		SecurityAPI:         SecurityClient,
		MultidomainAPI:      MultidomainClient,
	}

	clientCache[key] = builtClient
	return builtClient, nil
}

// Client represents the nutanix API client
type Client struct {
	API                 *v3.Client
	KarbonAPI           *karbon.Client
	FoundationClientAPI *foundation.Client
	WaitTimeout         int64
	FoundationCentral   *foundation_central.Client
	Era                 *era.Client
	NetworkingAPI       *networking.Client
	PrismAPI            *prism.Client
	MicroSegAPI         *microseg.Client
	IamAPI              *iam.Client
	ClusterAPI          *clusters.Client
	VolumeAPI           *volumes.Client
	DataProtectionAPI   *dataprotection.Client
	VmmAPI              *vmm.Client
	DataPoliciesAPI     *datapolicies.Client
	LcmAPI              *lcm.Client
	CalmAPI             *selfservice.Client
	ObjectStoreAPI      *objectstores.Client
	SecurityAPI         *security.Client
	MultidomainAPI      *multidomain.Client
}

// GetWaitTimeout returns the provider-level wait_timeout in minutes.
// A return value of 0 means the operator did not configure one and
// callers should fall back to per-resource defaults.
func (c *Client) GetWaitTimeout() int64 {
	if c == nil {
		return 0
	}
	return c.WaitTimeout
}
