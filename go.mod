module github.com/terraform-providers/terraform-provider-nutanix

require (
	github.com/PaesslerAG/jsonpath v0.1.1
	github.com/client9/misspell v0.3.4
	github.com/golangci/golangci-lint v1.25.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.10.1
	github.com/mitchellh/gox v1.0.1
	github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4 v4.0.1-beta.1
	// github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v16 v16.8.0-5295 // indirect
	//github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16 v4.0.2-beta.1
	github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16 v16.9.0-8634
	github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4 v4.0.1-beta.1
	// github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4 v4.0.3-alpha.2
	github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v16 v16.8.0-5280
	github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v16 v16.9.0-8538
	github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v16 v16.9.0-8500
	github.com/spf13/cast v1.3.1
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
)

go 1.13
