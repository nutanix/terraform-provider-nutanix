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
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Tag         string        `json:"tag,omitempty"`
	IsStruct    bool          `json:"is_struct,omitempty"`
	Fields      []NestedField `json:"fields,omitempty"`       // Nested fields if this is a struct type
	ImportAlias string        `json:"import_alias,omitempty"` // e.g., "import1", "import2" if type uses an import
	ImportPath  string        `json:"import_path,omitempty"`  // Full import path for the alias
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
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	File        string        `json:"file,omitempty"`
	Fields      []NestedField `json:"fields"`
	ImportAlias string        `json:"import_alias,omitempty"` // e.g., "import1", "import2" if type uses an import
	ImportPath  string        `json:"import_path,omitempty"`  // Full import path for the alias
}

// ImportMapping represents an import alias and its path
type ImportMapping struct {
	Alias string `json:"alias"` // e.g., "import1", "import2"
	Path  string `json:"path"`  // e.g., "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/common"
}

// OutputData is the unified output structure
type OutputData struct {
	Package                string                     `json:"package"`
	PackagePath            string                     `json:"package_path"`
	APIRequestResponseList []APIRequestResponseStruct `json:"api_request_response_struct"`
}

// Internal structures for extraction
type structInfo struct {
	name    string
	fields  []fieldInfo
	file    string
	imports map[string]string // import alias -> import path for this file
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

// shouldIgnoreMethod checks if a method should be ignored based on its name
func shouldIgnoreMethod(methodName string) bool {
	ignoredMethods := map[string]bool{
		"UnmarshalJSON": true,
		"MarshalJSON":   true,
		"SetValue":      true,
		"GetValue":      true,
		"SetData":       true,
		"GetData":       true,
	}
	return ignoredMethods[methodName]
}

// extractImportAliasFromType extracts import alias from a type string
// e.g., "import1.Parameter" -> "import1", "[]import2.ApiLink" -> "import2"
func extractImportAliasFromType(typeStr string) string {
	// Match pattern like "import1.Type" or "[]import2.Type" or "*import3.Type"
	re := regexp.MustCompile(`import\d+`)
	matches := re.FindString(typeStr)
	return matches
}

// collectUsedImports recursively collects all import aliases used in a NestedStruct
func collectUsedImports(struct_ *NestedStruct, usedImports map[string]bool) {
	if struct_ == nil {
		return
	}

	// Check the struct type itself
	if alias := extractImportAliasFromType(struct_.Type); alias != "" {
		usedImports[alias] = true
	}

	// Recursively check all fields
	for _, field := range struct_.Fields {
		// Check field type
		if alias := extractImportAliasFromType(field.Type); alias != "" {
			usedImports[alias] = true
		}
		// Recursively check nested fields
		if len(field.Fields) > 0 {
			nestedStruct := &NestedStruct{
				Type:   field.Type,
				Fields: field.Fields,
			}
			collectUsedImports(nestedStruct, usedImports)
		}
	}
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

	// Check if the struct type itself uses an import
	importAlias := extractImportAliasFromType(structName)
	importPath := ""
	if importAlias != "" && structInfo.imports != nil {
		if path, exists := structInfo.imports[importAlias]; exists {
			importPath = path
		}
	}

	nestedStruct := &NestedStruct{
		Name:        structName,
		Type:        structName,
		File:        structInfo.file,
		Fields:      []NestedField{},
		ImportAlias: importAlias,
		ImportPath:  importPath,
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

		// Extract import information if the field type uses an import
		// Use the imports from the file where this struct is defined
		importAlias := extractImportAliasFromType(field.typ)
		importPath := ""
		if importAlias != "" && structInfo.imports != nil {
			if path, exists := structInfo.imports[importAlias]; exists {
				importPath = path
			}
		}

		nestedField := NestedField{
			Name:        field.name,
			Type:        field.typ,
			Tag:         field.tag,
			IsStruct:    isStruct,
			Fields:      nestedFields,
			ImportAlias: importAlias,
			ImportPath:  importPath,
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

// extractImportsFromFile extracts import statements from a Go file
func extractImportsFromFile(fset *token.FileSet, file *ast.File) map[string]string {
	imports := make(map[string]string)

	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		alias := ""

		if imp.Name != nil {
			// Use explicit alias if provided
			alias = imp.Name.Name
		} else {
			// Extract alias from path (last component)
			parts := strings.Split(importPath, "/")
			if len(parts) > 0 {
				alias = parts[len(parts)-1]
			}
		}

		// Only track imports that look like import aliases (import1, import2, import3, etc.)
		// Match pattern like "import1", "import2", "import4", etc.
		if matched, _ := regexp.MatchString(`^import\d+$`, alias); matched {
			imports[alias] = importPath
		}
	}

	return imports
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

		// Extract imports from this file - store per file
		fileImports := extractImportsFromFile(fset, file)

		relPath, _ := filepath.Rel(packageDir, path)

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				if st, ok := x.Type.(*ast.StructType); ok {
					structInfo := structInfo{
						name:    x.Name.Name,
						fields:  []fieldInfo{},
						file:    relPath,
						imports: fileImports, // Store imports for this specific file
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
					// Skip ignored methods
					if shouldIgnoreMethod(x.Name.Name) {
						return true
					}

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

// filterByReceiverAndKeyword filters API methods and related structs based on receiver and/or keyword
// When receiver is provided, only includes methods where the receiver exactly matches
// When keyword is provided (and receiver is not), filters by keyword matching method name only
// When neither is provided, returns everything
func filterByReceiverAndKeyword(apiMethods []methodInfo, structsMap map[string]structInfo, receiver, keyword string) ([]methodInfo, map[string]structInfo) {
	// If neither receiver nor keyword is provided, return everything
	if receiver == "" && keyword == "" {
		return apiMethods, structsMap
	}

	var filteredMethods []methodInfo
	usedStructs := make(map[string]bool)

	if receiver != "" {
		// If receiver is provided, match exactly by receiver
		receiverLower := strings.ToLower(receiver)
		for _, method := range apiMethods {
			// Normalize receiver by removing pointer prefix and converting to lowercase
			receiverNormalized := strings.ToLower(strings.TrimPrefix(method.receiver, "*"))

			// Only include if receiver exactly matches
			if receiverNormalized == receiverLower {
				filteredMethods = append(filteredMethods, method)

				// Mark related structs as used
				requestType := findRequestType(method.params)
				if requestType != "" {
					baseRequestType := extractBaseTypeName(requestType)
					usedStructs[baseRequestType] = true
				}

				responseType := findResponseType(method.returns)
				if responseType != "" {
					baseResponseType := extractBaseTypeName(responseType)
					usedStructs[baseResponseType] = true
				}
			}
		}
	} else {
		// If only keyword is provided, match by keyword in method name only
		keywordLower := strings.ToLower(keyword)
		for _, method := range apiMethods {
			methodNameLower := strings.ToLower(method.name)

			// Only check if method name contains keyword
			if strings.Contains(methodNameLower, keywordLower) {
				filteredMethods = append(filteredMethods, method)

				// Mark related structs as used
				requestType := findRequestType(method.params)
				if requestType != "" {
					baseRequestType := extractBaseTypeName(requestType)
					usedStructs[baseRequestType] = true
				}

				responseType := findResponseType(method.returns)
				if responseType != "" {
					baseResponseType := extractBaseTypeName(responseType)
					usedStructs[baseResponseType] = true
				}
			}
		}
	}

	// Filter structs map to only include used structs and their dependencies
	filteredStructsMap := make(map[string]structInfo)
	visitedStructs := make(map[string]bool)

	var collectStructDependencies func(structName string)
	collectStructDependencies = func(structName string) {
		if visitedStructs[structName] {
			return
		}
		visitedStructs[structName] = true

		structInfo, exists := structsMap[structName]
		if !exists {
			return
		}

		filteredStructsMap[structName] = structInfo

		// Recursively collect dependencies
		for _, field := range structInfo.fields {
			baseType := extractBaseTypeName(field.typ)
			if baseType != "" && baseType != structName {
				collectStructDependencies(baseType)
			}
		}
	}

	// Collect all used structs and their dependencies
	for structName := range usedStructs {
		collectStructDependencies(structName)
	}

	return filteredMethods, filteredStructsMap
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
			// Set import info for the struct itself if it uses an import
			// Need to find which file defines this struct to get its imports
			if requestStruct != nil {
				importAlias := extractImportAliasFromType(requestType)
				if importAlias != "" {
					// Find the struct in structsMap to get its file's imports
					if structInfo, exists := structsMap[baseRequestType]; exists && structInfo.imports != nil {
						if path, exists := structInfo.imports[importAlias]; exists {
							requestStruct.ImportAlias = importAlias
							requestStruct.ImportPath = path
						}
					}
				}
			}
		}

		// Find response type
		responseType := findResponseType(method.returns)
		var responseStruct *NestedStruct
		if responseType != "" {
			baseResponseType := extractBaseTypeName(responseType)
			visited := make(map[string]bool)
			responseStruct = buildNestedStruct(baseResponseType, structsMap, visited)
			// Set import info for the struct itself if it uses an import
			// Need to find which file defines this struct to get its imports
			if responseStruct != nil {
				importAlias := extractImportAliasFromType(responseType)
				if importAlias != "" {
					// Find the struct in structsMap to get its file's imports
					if structInfo, exists := structsMap[baseResponseType]; exists && structInfo.imports != nil {
						if path, exists := structInfo.imports[importAlias]; exists {
							responseStruct.ImportAlias = importAlias
							responseStruct.ImportPath = path
						}
					}
				}
			}
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
	outputDirFlag := flag.String("output-dir", "code_gen/sdk_extract_output", "Output directory for extracted information")
	apiPackageFlag := flag.String("api-package", "", "Specific API package to extract")
	receiverFlag := flag.String("receiver", "", "Optional receiver type to filter extraction (only extract APIs/methods with matching receiver)")
	keywordFlag := flag.String("keyword", "", "Optional keyword to filter extraction (only extract APIs/methods matching this keyword, used if receiver is not provided)")
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

	// Also extract structs from models directory if it exists (to get their imports)
	modelsDir := filepath.Join(filepath.Dir(packageDir), "models")
	if fileExists(modelsDir) {
		fmt.Printf("Extracting structs from models directory...\n")
		modelStructs, _, err := extractFromPackage(modelsDir)
		if err == nil {
			// Merge model structs into structsMap (they have their own imports)
			for name, info := range modelStructs {
				structsMap[name] = info
			}
		}
	}

	// Filter API methods
	apiMethods := filterAPIMethods(methods)

	// Apply receiver/keyword filter if provided
	if *receiverFlag != "" || *keywordFlag != "" {
		if *receiverFlag != "" {
			fmt.Printf("Filtering by receiver: %s\n", *receiverFlag)
		} else {
			fmt.Printf("Filtering by keyword: %s\n", *keywordFlag)
		}
		apiMethods, structsMap = filterByReceiverAndKeyword(apiMethods, structsMap, *receiverFlag, *keywordFlag)
		fmt.Printf("  Filtered to %d API methods\n", len(apiMethods))
		fmt.Printf("  Filtered to %d structs\n", len(structsMap))
	}

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
