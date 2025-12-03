package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
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
func extractNamespace(packagePath string) string {
	re := regexp.MustCompile(`/([a-z0-9-]+)-go-client/`)
	matches := re.FindStringSubmatch(packagePath)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
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

// getTerraformType converts Go type to Terraform schema type
func getTerraformType(goType string, isStruct bool) string {
	goType = strings.TrimPrefix(goType, "*")

	if isStruct {
		return "TypeList"
	}

	if strings.HasPrefix(goType, "[]") {
		innerType := strings.TrimPrefix(goType, "[]")
		if strings.Contains(innerType, ".") || strings.Contains(innerType, "import") {
			return "TypeList" // Complex type, use List
		}
		return "TypeList" // Array of primitives
	}

	if strings.Contains(goType, "time.Time") {
		return "TypeString"
	}

	switch goType {
	case "string", "*string":
		return "TypeString"
	case "int", "*int", "int64", "*int64":
		return "TypeInt"
	case "bool", "*bool":
		return "TypeBool"
	case "float64", "*float64":
		return "TypeFloat"
	case "map[string]string", "map[string]interface{}":
		return "TypeMap"
	default:
		if strings.Contains(goType, "map") {
			return "TypeMap"
		}
		return "TypeList" // Default to List for complex types
	}
}

// generateSchemaField generates schema for a field
func generateSchemaField(field NestedField, indent string, fieldNameMap map[string]int) string {
	fieldName := toSnakeCase(field.Name)

	// Handle duplicate field names
	if count, exists := fieldNameMap[fieldName]; exists {
		fieldNameMap[fieldName] = count + 1
		fieldName = fmt.Sprintf("%s_%d", fieldName, count+1)
	} else {
		fieldNameMap[fieldName] = 1
	}

	terraformType := getTerraformType(field.Type, field.IsStruct)

	var result strings.Builder

	if field.IsStruct && len(field.Fields) > 0 {
		// Nested struct - create a Resource schema
		result.WriteString(fmt.Sprintf(`%s"%s": {
%s	Type:     schema.%s,
%s	Computed: true,
%s	Elem: &schema.Resource{
%s		Schema: map[string]*schema.Schema{
`, indent, fieldName, indent, terraformType, indent, indent, indent))

		nestedIndent := indent + "\t\t"
		for _, nestedField := range field.Fields {
			// Skip oneOf fields (they start with "oneOfType")
			if strings.HasPrefix(nestedField.Name, "oneOfType") {
				continue
			}
			result.WriteString(generateSchemaField(nestedField, nestedIndent, fieldNameMap))
		}

		result.WriteString(fmt.Sprintf(`%s		},
%s	},
%s},
`, indent, indent, indent))
	} else if strings.HasPrefix(field.Type, "[]") && field.IsStruct {
		// Array of structs
		innerType := strings.TrimPrefix(field.Type, "[]*")
		innerType = strings.TrimPrefix(innerType, "[]")
		innerType = strings.TrimPrefix(innerType, "*")

		result.WriteString(fmt.Sprintf(`%s"%s": {
%s	Type:     schema.%s,
%s	Computed: true,
%s	Elem: &schema.Resource{
%s		Schema: map[string]*schema.Schema{
`, indent, fieldName, indent, terraformType, indent, indent, indent))

		nestedIndent := indent + "\t\t"
		for _, nestedField := range field.Fields {
			if strings.HasPrefix(nestedField.Name, "oneOfType") {
				continue
			}
			result.WriteString(generateSchemaField(nestedField, nestedIndent, fieldNameMap))
		}

		result.WriteString(fmt.Sprintf(`%s		},
%s	},
%s},
`, indent, indent, indent))
	} else {
		// Simple field
		elemType := ""
		if strings.HasPrefix(field.Type, "[]") {
			innerType := strings.TrimPrefix(field.Type, "[]*")
			innerType = strings.TrimPrefix(innerType, "[]")
			if innerType == "string" {
				elemType = `\n\t\t\tType: schema.TypeString,\n`
			}
		}

		if elemType != "" {
			result.WriteString(fmt.Sprintf(`%s"%s": {
%s	Type:     schema.%s,
%s	Computed: true,
%s	Elem: &schema.Schema{
%s		},
%s},
`, indent, fieldName, indent, terraformType, indent, indent, indent, indent))
		} else {
			result.WriteString(fmt.Sprintf(`%s"%s": {
%s	Type:     schema.%s,
%s	Computed: true,
%s},
`, indent, fieldName, indent, terraformType, indent, indent))
		}
	}

	return result.String()
}

// generateFlattenFunction generates flatten function for nested structures
func generateFlattenFunction(field NestedField, structName string) string {
	funcName := fmt.Sprintf("flatten%s", toCamelCase(field.Name))

	var result strings.Builder
	result.WriteString(fmt.Sprintf(`func %s(pr *%s) []map[string]interface{} {
	if pr == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0)
	item := make(map[string]interface{})

`, funcName, field.Type))

	for _, nestedField := range field.Fields {
		if strings.HasPrefix(nestedField.Name, "oneOfType") {
			continue
		}

		fieldName := toSnakeCase(nestedField.Name)
		goFieldName := nestedField.Name

		if nestedField.IsStruct && len(nestedField.Fields) > 0 {
			// Nested struct - call flatten function
			nestedFuncName := fmt.Sprintf("flatten%s", toCamelCase(nestedField.Name))
			if strings.HasPrefix(nestedField.Type, "[]") {
				result.WriteString(fmt.Sprintf(`	if pr.%s != nil {
		item["%s"] = %s(pr.%s)
	}
`, goFieldName, fieldName, nestedFuncName, goFieldName))
			} else {
				result.WriteString(fmt.Sprintf(`	if pr.%s != nil {
		item["%s"] = %s(pr.%s)
	}
`, goFieldName, fieldName, nestedFuncName, goFieldName))
			}
		} else {
			// Simple field
			if strings.Contains(nestedField.Type, "time.Time") {
				result.WriteString(fmt.Sprintf(`	if pr.%s != nil {
		item["%s"] = pr.%s.String()
	}
`, goFieldName, fieldName, goFieldName))
			} else if strings.HasPrefix(nestedField.Type, "*") {
				result.WriteString(fmt.Sprintf(`	item["%s"] = utils.StringValue(pr.%s)
`, fieldName, goFieldName))
			} else if strings.HasPrefix(nestedField.Type, "[]") {
				result.WriteString(fmt.Sprintf(`	item["%s"] = pr.%s
`, fieldName, goFieldName))
			} else {
				result.WriteString(fmt.Sprintf(`	item["%s"] = pr.%s
`, fieldName, goFieldName))
			}
		}
	}

	result.WriteString(`	result = append(result, item)
	return result
}

`)

	return result.String()
}

func main() {
	jsonFileFlag := flag.String("json", "sdk_extract_output/sdk_info.json", "Path to sdk_info.json file")
	apiMethodFlag := flag.String("method", "GetAlertById", "API method name to generate datasource for")
	outputDirFlag := flag.String("output-dir", "../nutanix/services", "Output directory for generated files")
	namespaceFlag := flag.String("namespace", "", "Namespace for the service (e.g., monitoringv2). If empty, will be derived from package")
	flag.Parse()

	// Resolve absolute paths
	jsonFilePath, err := filepath.Abs(*jsonFileFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving JSON file path: %v\n", err)
		os.Exit(1)
	}

	outputDir, err := filepath.Abs(*outputDirFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving output directory path: %v\n", err)
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

	// Find the API method
	var targetAPI *APIRequestResponseStruct
	for i := range outputData.APIRequestResponseList {
		if outputData.APIRequestResponseList[i].APIMethod.Name == *apiMethodFlag {
			targetAPI = &outputData.APIRequestResponseList[i]
			break
		}
	}

	if targetAPI == nil {
		fmt.Fprintf(os.Stderr, "Error: API method '%s' not found in JSON\n", *apiMethodFlag)
		os.Exit(1)
	}

	if targetAPI.ResponseStruct == nil {
		fmt.Fprintf(os.Stderr, "Error: No response structure found for API method '%s'\n", *apiMethodFlag)
		os.Exit(1)
	}

	// Extract namespace
	namespace := *namespaceFlag
	if namespace == "" {
		baseNamespace := extractNamespace(outputData.PackagePath)
		namespace = baseNamespace + "v2"
	}

	// Generate datasource name
	datasourceName := fmt.Sprintf("nutanix_%s_%s", toSnakeCase(namespace), toSnakeCase(*apiMethodFlag))
	funcName := fmt.Sprintf("Datasource%s", toCamelCase(datasourceName))
	readFuncName := fmt.Sprintf("%sRead", funcName)

	// Extract receiver API name
	receiver := strings.TrimPrefix(targetAPI.APIMethod.Receiver, "*")
	apiTypeName := strings.TrimSuffix(receiver, "Api")

	// Generate schema
	fieldNameMap := make(map[string]int)
	var schemaFields strings.Builder

	// Add ext_id as required field (assuming it's a GetById method)
	schemaFields.WriteString(`	"ext_id": {
		Type:     schema.TypeString,
		Required: true,
	},
`)

	// Generate schema for response fields
	responseDataField := targetAPI.ResponseStruct.Fields[0] // Usually "Data" field
	if responseDataField.IsStruct && len(responseDataField.Fields) > 0 {
		// Find the actual data type (skip oneOfType400 error responses)
		for _, field := range responseDataField.Fields {
			if strings.HasPrefix(field.Name, "oneOfType") && !strings.Contains(field.Name, "400") {
				// This is the success response
				for _, dataField := range field.Fields {
					if !strings.HasPrefix(dataField.Name, "oneOfType") {
						schemaFields.WriteString(generateSchemaField(dataField, "\t", fieldNameMap))
					}
				}
				break
			}
		}
	}

	// Generate flatten functions for nested structures
	var flattenFuncs strings.Builder
	generateFlattenFuncsRecursive(targetAPI.ResponseStruct, &flattenFuncs)

	// Generate the Go code
	tmpl := `package {{.Namespace}}

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func {{.FuncName}}() *schema.Resource {
	return &schema.Resource{
		ReadContext: {{.ReadFuncName}},
		Schema: map[string]*schema.Schema{
{{.SchemaFields}}		},
	}
}

func {{.ReadFuncName}}(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).{{.ClientField}}

	extID := d.Get("ext_id").(string)

	resp, err := conn.{{.APIField}}.{{.MethodName}}(utils.StringPtr(extID))
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching {{.ResourceName}}: %v", errorMessage["message"])
	}

	// TODO: Extract and set response fields
	// Example: getResp := resp.Data.GetValue().(YourType)
	// if err := d.Set("field_name", getResp.FieldName); err != nil {
	//     return diag.FromErr(err)
	// }

	d.SetId(extID)
	return nil
}

{{.FlattenFuncs}}
`

	t := template.Must(template.New("datasource").Parse(tmpl))

	// Determine client field name
	clientField := toCamelCase(namespace)
	clientField = strings.TrimSuffix(clientField, "v2")
	clientField = clientField + "API"

	var buf strings.Builder
	err = t.Execute(&buf, map[string]interface{}{
		"Namespace":    namespace,
		"FuncName":     funcName,
		"ReadFuncName": readFuncName,
		"SchemaFields": schemaFields.String(),
		"ClientField":  clientField,
		"APIField":     apiTypeName + "API",
		"MethodName":   *apiMethodFlag,
		"ResourceName": datasourceName,
		"FlattenFuncs": flattenFuncs.String(),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	// Format the code
	formattedCode, err := format.Source([]byte(buf.String()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not format code: %v\n", err)
		formattedCode = []byte(buf.String())
	}

	// Create namespace directory
	namespaceDir := filepath.Join(outputDir, namespace)
	if err := os.MkdirAll(namespaceDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating namespace directory: %v\n", err)
		os.Exit(1)
	}

	// Write the generated file
	outputFile := filepath.Join(namespaceDir, fmt.Sprintf("data_source_%s.go", toSnakeCase(datasourceName)))
	if err := os.WriteFile(outputFile, formattedCode, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated datasource file: %s\n", outputFile)
	fmt.Printf("Successfully generated Terraform datasource for: %s\n", *apiMethodFlag)
}

func generateFlattenFuncsRecursive(structDef *NestedStruct, result *strings.Builder) {
	for _, field := range structDef.Fields {
		if field.IsStruct && len(field.Fields) > 0 && !strings.HasPrefix(field.Name, "oneOfType") {
			result.WriteString(generateFlattenFunction(field, field.Type))
			// Recursively generate for nested structs
			for _, nestedField := range field.Fields {
				if nestedField.IsStruct && len(nestedField.Fields) > 0 {
					nestedStruct := &NestedStruct{
						Name:   nestedField.Name,
						Type:   nestedField.Type,
						Fields: nestedField.Fields,
					}
					generateFlattenFuncsRecursive(nestedStruct, result)
				}
			}
		}
	}
}
