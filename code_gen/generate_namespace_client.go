package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

// OutputData matches the structure from extract_sdk_info.go
type OutputData struct {
	Package                string                     `json:"package"`
	PackagePath            string                     `json:"package_path"`
	APIRequestResponseList []APIRequestResponseStruct `json:"api_request_response_struct"`
}

type APIRequestResponseStruct struct {
	APIMethod      APIMethodInfo `json:"api_method"`
	RequestStruct  *NestedStruct `json:"request_struct,omitempty"`
	ResponseStruct *NestedStruct `json:"response_struct,omitempty"`
}

type APIMethodInfo struct {
	Name     string        `json:"name"`
	Receiver string        `json:"receiver"`
	File     string        `json:"file"`
	Params   []MethodParam `json:"params"`
}

type MethodParam struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type NestedStruct struct {
	Name   string        `json:"name"`
	Type   string        `json:"type"`
	File   string        `json:"file,omitempty"`
	Fields []NestedField `json:"fields"`
}

type NestedField struct {
	Name     string        `json:"name"`
	Type     string        `json:"type"`
	Tag      string        `json:"tag,omitempty"`
	IsStruct bool          `json:"is_struct,omitempty"`
	Fields   []NestedField `json:"fields,omitempty"`
}

// extractNamespace extracts the namespace from package path
// e.g., "monitoring-go-client" -> "monitoring"
func extractNamespace(packagePath string) string {
	// Extract the client name from path like:
	// github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4
	re := regexp.MustCompile(`/([a-z0-9-]+)-go-client/`)
	matches := re.FindStringSubmatch(packagePath)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractAPITypeName extracts the API type name from receiver
// e.g., "*AlertEmailConfigurationApi" -> "AlertEmailConfigurationApi"
func extractAPITypeName(receiver string) string {
	// Remove pointer prefix
	receiver = strings.TrimPrefix(receiver, "*")
	// Remove "Api" suffix if present
	if strings.HasSuffix(receiver, "Api") {
		return receiver
	}
	return receiver
}

// toCamelCase converts a string to CamelCase
func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	parts := strings.Split(s, "-")
	var result strings.Builder
	for _, part := range parts {
		if part != "" {
			result.WriteString(strings.ToUpper(part[:1]) + part[1:])
		}
	}
	return result.String()
}

// toFieldName converts API type name to a field name
// e.g., "AlertEmailConfigurationApi" -> "AlertEmailConfigurationAPI"
// or "SubnetsApi" -> "SubnetAPIInstance" (following existing patterns)
func toFieldName(apiType string) string {
	// Remove "Api" suffix
	if strings.HasSuffix(apiType, "Api") {
		baseName := apiType[:len(apiType)-3]
		// For some common patterns, use specific naming
		if strings.HasSuffix(baseName, "s") {
			// Plural to singular + APIInstance (e.g., Subnets -> SubnetAPIInstance)
			singular := baseName[:len(baseName)-1]
			return singular + "APIInstance"
		}
		return baseName + "API"
	}
	return apiType + "API"
}

// generateClientCode generates the Go code for the namespace client
func generateClientCode(namespace string, apiTypes []string, packagePath string) (string, error) {
	// Extract client import alias (first part of namespace)
	namespaceParts := strings.Split(namespace, "-")
	clientAlias := namespaceParts[0]

	// Build import path
	apiImportPath := packagePath + "/api"
	clientImportPath := packagePath + "/client"

	// Sort API types for consistent output
	sort.Strings(apiTypes)

	// Generate field names and initialization
	type fieldInfo struct {
		FieldName string
		APIType   string
		NewFunc   string
	}

	fields := make([]fieldInfo, 0, len(apiTypes))
	seen := make(map[string]bool)

	for _, apiType := range apiTypes {
		// Skip ApiClient as it's handled separately
		if apiType == "ApiClient" {
			continue
		}

		fieldName := toFieldName(apiType)
		// Ensure unique field names
		if seen[fieldName] {
			// Add a suffix if duplicate
			counter := 1
			originalName := fieldName
			for seen[fieldName] {
				fieldName = fmt.Sprintf("%s%d", originalName, counter)
				counter++
			}
		}
		seen[fieldName] = true

		newFunc := fmt.Sprintf("New%s", apiType)
		fields = append(fields, fieldInfo{
			FieldName: fieldName,
			APIType:   apiType,
			NewFunc:   newFunc,
		})
	}

	// Template for the generated code
	tmpl := `package {{.Namespace}}

import (
	"{{.APIImportPath}}"
	{{.ClientAlias}} "{{.ClientImportPath}}"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
{{range .Fields}}	{{.FieldName}} *api.{{.APIType}}
{{end}}	APIClientInstance *{{.ClientAlias}}.ApiClient
}

func New{{.ClientName}}Client(credentials client.Credentials) (*Client, error) {
	var baseClient *{{.ClientAlias}}.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := {{.ClientAlias}}.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = 9440
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
{{range .Fields}}		{{.FieldName}}: api.{{.NewFunc}}(baseClient),
{{end}}		APIClientInstance: {{.ClientAlias}}.NewApiClient(),
	}

	return f, nil
}
`

	t := template.Must(template.New("client").Parse(tmpl))

	var buf strings.Builder
	err := t.Execute(&buf, map[string]interface{}{
		"Namespace":        namespace,
		"APIImportPath":    apiImportPath,
		"ClientImportPath": clientImportPath,
		"ClientAlias":      clientAlias,
		"ClientName":       toCamelCase(namespace),
		"Fields":           fields,
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func main() {
	jsonFileFlag := flag.String("json", "sdk_extract_output/sdk_info.json", "Path to sdk_info.json file")
	outputDirFlag := flag.String("output-dir", "../nutanix/sdks/v4", "Output directory for generated namespace files")
	flag.Parse()

	// Resolve absolute paths to ensure we use the correct locations regardless of where script is run from
	outputDir, err := filepath.Abs(*outputDirFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving output directory path: %v\n", err)
		os.Exit(1)
	}

	// Resolve JSON file path to absolute path
	jsonFilePath, err := filepath.Abs(*jsonFileFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving JSON file path: %v\n", err)
		os.Exit(1)
	}

	// Read JSON file
	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading JSON file: %v\n", err)
		os.Exit(1)
	}

	// Parse JSON
	var outputData OutputData
	if err := json.Unmarshal(jsonData, &outputData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Extract namespace
	namespace := extractNamespace(outputData.PackagePath)
	if namespace == "" {
		fmt.Fprintf(os.Stderr, "Error: Could not extract namespace from package path: %s\n", outputData.PackagePath)
		os.Exit(1)
	}

	fmt.Printf("Extracted namespace: %s\n", namespace)

	// Extract unique API types from receivers
	apiTypesMap := make(map[string]bool)
	for _, apiStruct := range outputData.APIRequestResponseList {
		receiver := apiStruct.APIMethod.Receiver
		apiType := extractAPITypeName(receiver)
		// Only include types that end with "Api" (actual API clients)
		if strings.HasSuffix(apiType, "Api") || apiType == "ApiClient" {
			apiTypesMap[apiType] = true
		}
	}

	apiTypes := make([]string, 0, len(apiTypesMap))
	for apiType := range apiTypesMap {
		apiTypes = append(apiTypes, apiType)
	}

	fmt.Printf("Found %d unique API types\n", len(apiTypes))

	// Generate the client code
	code, err := generateClientCode(namespace, apiTypes, outputData.PackagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	// Create namespace directory (use resolved absolute path)
	namespaceDir := filepath.Join(outputDir, namespace)
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Printf("Namespace directory: %s\n", namespaceDir)
	if err := os.MkdirAll(namespaceDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating namespace directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Printf("Namespace directory: %s\n", namespaceDir)

	// Format the code using gofmt
	formattedCode, err := format.Source([]byte(code))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not format code with gofmt: %v\n", err)
		fmt.Fprintf(os.Stderr, "Writing unformatted code...\n")
		formattedCode = []byte(code)
	}

	// Write the generated file
	outputFile := filepath.Join(namespaceDir, namespace+".go")
	if err := os.WriteFile(outputFile, formattedCode, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	// Try to run gofmt on the file to ensure it's properly formatted
	if err := exec.Command("gofmt", "-w", outputFile).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not run gofmt on output file: %v\n", err)
	}

	fmt.Printf("Generated client file: %s\n", outputFile)
	fmt.Printf("Successfully generated namespace client for: %s\n", namespace)
}
