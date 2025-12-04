package nutanix

import (
	"fmt"

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
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/networking"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/objectstores"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
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
}

// Client ...
func (c *Config) Client() (*Client, error) {
	configCreds := client.Credentials{
		URL:                fmt.Sprintf("%s:%s", c.Endpoint, c.Port),
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
	return &Client{
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
	}, nil
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
}
