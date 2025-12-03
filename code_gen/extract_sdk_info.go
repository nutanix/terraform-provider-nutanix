package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// NestedField represents a field with nested structure support
type NestedField struct {
	Name     string        `json:"name"`
	Type     string        `json:"type"`
	Tag      string        `json:"tag,omitempty"`
	IsStruct bool          `json:"is_struct,omitempty"`
	Fields   []NestedField `json:"fields,omitempty"` // Nested fields if this is a struct type
}

// APIRequestResponseStruct combines API method with its request and response structs
type APIRequestResponseStruct struct {
	APIMethod      APIMethodInfo `json:"api_method"`
	RequestStruct  *NestedStruct `json:"request_struct,omitempty"`
	ResponseStruct *NestedStruct `json:"response_struct,omitempty"`
}

// MethodParam represents a method parameter
type MethodParam struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// APIMethodInfo contains information about an API method
type APIMethodInfo struct {
	Name     string        `json:"name"`
	Receiver string        `json:"receiver"`
	File     string        `json:"file"`
	Params   []MethodParam `json:"params"` // Required arguments/parameters
}

// NestedStruct represents a struct with nested field structure
type NestedStruct struct {
	Name   string        `json:"name"`
	Type   string        `json:"type"`
	File   string        `json:"file,omitempty"`
	Fields []NestedField `json:"fields"`
}

// OutputData is the unified output structure
type OutputData struct {
	Package                string                     `json:"package"`
	PackagePath            string                     `json:"package_path"`
	APIRequestResponseList []APIRequestResponseStruct `json:"api_request_response_struct"`
}

// Internal structures for extraction
type structInfo struct {
	name   string
	fields []fieldInfo
	file   string
}

type fieldInfo struct {
	name string
	typ  string
	tag  string
}

type methodInfo struct {
	name     string
	receiver string
	params   []paramInfo
	returns  []string
	file     string
}

type paramInfo struct {
	name string
	typ  string
}

// formatTypeString formats a type expression using go/format
func formatTypeString(fset *token.FileSet, expr ast.Expr, src []byte) string {
	var buf strings.Builder
	err := format.Node(&buf, fset, expr)
	if err != nil {
		return normalizeTypeString(string(src[expr.Pos()-1 : expr.End()-1]))
	}
	return strings.TrimSpace(buf.String())
}

// normalizeTypeString removes newlines and normalizes whitespace in type strings
func normalizeTypeString(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.TrimSpace(s)
}

// shouldIgnoreField checks if a field should be ignored based on its name
func shouldIgnoreField(fieldName string) bool {
	if fieldName == "" {
		return false // Don't ignore embedded fields by default
	}

	// Ignore fields ending with underscore (internal/metadata fields)
	if strings.HasSuffix(fieldName, "_") {
		return true
	}

	// Ignore specific known internal fields
	ignoredFields := map[string]bool{
		"ObjectType_":    true,
		"Reserved_":      true,
		"UnknownFields_": true,
		"Discriminator":  true,
	}
	if ignoredFields[fieldName] {
		return true
	}

	return false
}

// extractBaseTypeName extracts the base type name from complex types (pointers, slices, etc.)
func extractBaseTypeName(typeStr string) string {
	// Remove pointer
	typeStr = strings.TrimPrefix(typeStr, "*")
	// Remove slice brackets
	typeStr = regexp.MustCompile(`\[\][\s]*`).ReplaceAllString(typeStr, "")
	// Remove map brackets
	typeStr = regexp.MustCompile(`map\[[^\]]+\][\s]*`).ReplaceAllString(typeStr, "")
	// Remove package prefixes (e.g., "import1.GetAlertApiResponse" -> "GetAlertApiResponse")
	if idx := strings.LastIndex(typeStr, "."); idx != -1 {
		typeStr = typeStr[idx+1:]
	}
	// Remove array brackets
	typeStr = regexp.MustCompile(`\[[0-9]+\]`).ReplaceAllString(typeStr, "")
	return strings.TrimSpace(typeStr)
}

// buildNestedStruct builds a nested struct structure from a struct name
func buildNestedStruct(structName string, structsMap map[string]structInfo, visited map[string]bool) *NestedStruct {
	// Prevent infinite recursion
	if visited[structName] {
		return &NestedStruct{
			Name:   structName,
			Type:   structName,
			Fields: []NestedField{},
		}
	}
	visited[structName] = true
	defer delete(visited, structName)

	structInfo, exists := structsMap[structName]
	if !exists {
		return nil
	}

	nestedStruct := &NestedStruct{
		Name:   structName,
		Type:   structName,
		File:   structInfo.file,
		Fields: []NestedField{},
	}

	for _, field := range structInfo.fields {
		// Skip ignored fields
		if shouldIgnoreField(field.name) {
			continue
		}

		baseType := extractBaseTypeName(field.typ)
		isStruct := false
		var nestedFields []NestedField

		// Check if this field type is a struct we know about
		if nestedFieldStruct := buildNestedStruct(baseType, structsMap, visited); nestedFieldStruct != nil {
			isStruct = true
			nestedFields = nestedFieldStruct.Fields
		}

		nestedField := NestedField{
			Name:     field.name,
			Type:     field.typ,
			Tag:      field.tag,
			IsStruct: isStruct,
			Fields:   nestedFields,
		}

		nestedStruct.Fields = append(nestedStruct.Fields, nestedField)
	}

	return nestedStruct
}

// findRequestType finds the request type from method parameters
func findRequestType(params []paramInfo) string {
	// Skip context.Context and find the first meaningful parameter
	for _, param := range params {
		paramType := strings.TrimSpace(param.typ)
		// Skip context
		if strings.Contains(paramType, "context.Context") {
			continue
		}
		// Skip variadic args
		if strings.HasPrefix(paramType, "...") {
			continue
		}
		// Return the first non-context parameter
		return paramType
	}
	return ""
}

// findResponseType finds the response type from method returns
func findResponseType(returns []string) string {
	// Usually the first return is the response, second is error
	for _, ret := range returns {
		retType := strings.TrimSpace(ret)
		// Skip error
		if retType == "error" {
			continue
		}
		// Return the first non-error return
		return retType
	}
	return ""
}

func extractFromPackage(packageDir string) (map[string]structInfo, []methodInfo, error) {
	fset := token.NewFileSet()
	structsMap := make(map[string]structInfo)
	var methods []methodInfo

	err := filepath.Walk(packageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(packageDir, path)

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				if st, ok := x.Type.(*ast.StructType); ok {
					structInfo := structInfo{
						name:   x.Name.Name,
						fields: []fieldInfo{},
						file:   relPath,
					}

					if st.Fields != nil {
						for _, field := range st.Fields.List {
							fieldType := formatTypeString(fset, field.Type, src)
							for _, name := range field.Names {
								tagValue := ""
								if field.Tag != nil {
									tagValue = string(field.Tag.Value)
								}
								structInfo.fields = append(structInfo.fields, fieldInfo{
									name: name.Name,
									typ:  fieldType,
									tag:  tagValue,
								})
							}
							if len(field.Names) == 0 {
								// Embedded field
								structInfo.fields = append(structInfo.fields, fieldInfo{
									name: "",
									typ:  fieldType,
								})
							}
						}
					}

					structsMap[structInfo.name] = structInfo
				}
			case *ast.FuncDecl:
				if x.Recv != nil && len(x.Recv.List) > 0 {
					recvType := ""
					if ident, ok := x.Recv.List[0].Type.(*ast.Ident); ok {
						recvType = ident.Name
					} else if star, ok := x.Recv.List[0].Type.(*ast.StarExpr); ok {
						if ident, ok := star.X.(*ast.Ident); ok {
							recvType = "*" + ident.Name
						} else {
							recvType = formatTypeString(fset, x.Recv.List[0].Type, src)
						}
					} else {
						recvType = formatTypeString(fset, x.Recv.List[0].Type, src)
					}

					methodInfo := methodInfo{
						name:     x.Name.Name,
						receiver: recvType,
						params:   []paramInfo{},
						returns:  []string{},
						file:     relPath,
					}

					if x.Type.Params != nil {
						for _, param := range x.Type.Params.List {
							paramType := formatTypeString(fset, param.Type, src)
							for _, name := range param.Names {
								methodInfo.params = append(methodInfo.params, paramInfo{
									name: name.Name,
									typ:  paramType,
								})
							}
							if len(param.Names) == 0 {
								methodInfo.params = append(methodInfo.params, paramInfo{
									name: "",
									typ:  paramType,
								})
							}
						}
					}

					if x.Type.Results != nil {
						for _, result := range x.Type.Results.List {
							resultType := formatTypeString(fset, result.Type, src)
							methodInfo.returns = append(methodInfo.returns, resultType)
						}
					}

					methods = append(methods, methodInfo)
				}
			}
			return true
		})

		return nil
	})

	return structsMap, methods, err
}

func filterAPIMethods(methods []methodInfo) []methodInfo {
	var apiMethods []methodInfo
	for _, m := range methods {
		methodType := strings.ToLower(m.receiver)
		if strings.Contains(methodType, "api") || strings.HasSuffix(strings.ToLower(m.receiver), "api") {
			apiMethods = append(apiMethods, m)
		}
	}
	return apiMethods
}

func buildAPIRequestResponseList(apiMethods []methodInfo, structsMap map[string]structInfo) []APIRequestResponseStruct {
	var result []APIRequestResponseStruct

	for _, method := range apiMethods {
		// Convert method parameters to MethodParam format
		params := make([]MethodParam, 0, len(method.params))
		for _, param := range method.params {
			params = append(params, MethodParam{
				Name: param.name,
				Type: param.typ,
			})
		}

		apiMethodInfo := APIMethodInfo{
			Name:     method.name,
			Receiver: method.receiver,
			File:     method.file,
			Params:   params,
		}

		// Find request type
		requestType := findRequestType(method.params)
		var requestStruct *NestedStruct
		if requestType != "" {
			baseRequestType := extractBaseTypeName(requestType)
			visited := make(map[string]bool)
			requestStruct = buildNestedStruct(baseRequestType, structsMap, visited)
		}

		// Find response type
		responseType := findResponseType(method.returns)
		var responseStruct *NestedStruct
		if responseType != "" {
			baseResponseType := extractBaseTypeName(responseType)
			visited := make(map[string]bool)
			responseStruct = buildNestedStruct(baseResponseType, structsMap, visited)
		}

		result = append(result, APIRequestResponseStruct{
			APIMethod:      apiMethodInfo,
			RequestStruct:  requestStruct,
			ResponseStruct: responseStruct,
		})
	}

	return result
}

func main() {
	packageFlag := flag.String("package", "", "Go package path with optional version (e.g., github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4@v4.1.1)")
	outputDirFlag := flag.String("output-dir", "./sdk_extract_output", "Output directory for extracted information")
	apiPackageFlag := flag.String("api-package", "", "Specific API package to extract")
	flag.Parse()

	if *packageFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --package is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Parse package path
	packagePath := *packageFlag
	var version string
	if idx := strings.LastIndex(packagePath, "@"); idx != -1 {
		version = packagePath[idx+1:]
		packagePath = packagePath[:idx]
	}

	// Create output directory
	outputDir := *outputDirFlag
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Determine package directory
	fmt.Printf("Locating package %s...\n", packagePath)
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home, _ := os.UserHomeDir()
		gopath = filepath.Join(home, "go")
	}
	modCache := filepath.Join(gopath, "pkg", "mod")

	packageDir := ""
	if version != "" {
		packageDir = filepath.Join(modCache, strings.ReplaceAll(packagePath, "/", string(filepath.Separator))+"@"+version)
	} else {
		basePath := filepath.Join(modCache, strings.ReplaceAll(packagePath, "/", string(filepath.Separator)))
		if entries, err := os.ReadDir(basePath); err == nil {
			var versions []string
			for _, entry := range entries {
				if entry.IsDir() && strings.HasPrefix(entry.Name(), "v") {
					versions = append(versions, entry.Name())
				}
			}
			if len(versions) > 0 {
				packageDir = filepath.Join(basePath, versions[0])
			}
		}
	}

	if packageDir == "" || !fileExists(packageDir) {
		fmt.Fprintf(os.Stderr, "Error: Could not locate package %s\n", packagePath)
		fmt.Fprintf(os.Stderr, "Please ensure the package is downloaded: go mod download %s\n", *packageFlag)
		os.Exit(1)
	}

	fmt.Printf("Found package at: %s\n", packageDir)

	// If api-package is specified, look for api subdirectory
	if *apiPackageFlag != "" {
		apiDir := filepath.Join(packageDir, "api")
		if fileExists(apiDir) {
			packageDir = apiDir
			packagePath = *apiPackageFlag
		}
	}

	// Extract structs and methods
	fmt.Printf("Extracting structs and methods from %s...\n", packageDir)
	structsMap, methods, err := extractFromPackage(packageDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting from package: %v\n", err)
		os.Exit(1)
	}

	// Filter API methods
	apiMethods := filterAPIMethods(methods)

	// Build unified structure
	fmt.Printf("Building API request/response structures...\n")
	apiRequestResponseList := buildAPIRequestResponseList(apiMethods, structsMap)

	// Create output data
	outputData := OutputData{
		Package:                *packageFlag,
		PackagePath:            packagePath,
		APIRequestResponseList: apiRequestResponseList,
	}

	// Write JSON output
	jsonPath := filepath.Join(outputDir, "sdk_info.json")
	jsonData, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSON file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nExtraction complete!\n")
	fmt.Printf("  API methods processed: %d\n", len(apiMethods))
	fmt.Printf("  Structs found: %d\n", len(structsMap))
	fmt.Printf("  API request/response structures: %d\n", len(apiRequestResponseList))
	fmt.Printf("\nResults saved to: %s\n", jsonPath)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
