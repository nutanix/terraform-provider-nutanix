package karbon

type KarbonPrivateRegistryListResponse []KarbonPrivateRegistryResponse

type KarbonPrivateRegistryResponse struct {
	Name     *string `json:"name" mapstructure:"name, omitempty"`
	Endpoint *string `json:"endpoint" mapstructure:"endpoint, omitempty"`
	UUID     *string `json:"uuid" mapstructure:"uuid, omitempty"`
}

type KarbonPrivateRegistryIntentInput struct {
	Name *string `json:"name" mapstructure:"name, omitempty"`
	Cert *string `json:"cert" mapstructure:"cert, omitempty"`
	URL  *string `json:"url" mapstructure:"url, omitempty"`
	Port *int64  `json:"port" mapstructure:"port, omitempty"`
}

type KarbonPrivateRegistryOperationResponse struct {
	RegistryName *string `json:"registry_name" mapstructure:"registry_name, omitempty"`
}

type KarbonPrivateRegistryOperationIntentInput struct {
	RegistryName *string `json:"registry_name" mapstructure:"registry_name, omitempty"`
}
