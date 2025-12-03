# SDK Information Extraction Script

This script extracts config models and API methods from Nutanix Go SDK packages.

**Note:** There are two versions available:
- **Go version** (`extract_sdk_info.go`) - Recommended, pure Go implementation
- **Python version** (`extract_sdk_info.py`) - Alternative implementation

## Requirements

### For Go Version (Recommended)
- Go 1.16+ installed and in PATH
- Internet connection (to download Go modules)

### For Python Version
- Python 3.7+
- Go 1.16+ installed and in PATH
- Internet connection (to download Go modules)

## Usage

### Building the Go Version

First, build the Go script:

```bash
cd scripts
go build -mod=mod -o extract_sdk_info extract_sdk_info.go
```

Or run it directly:

```bash
go run -mod=mod scripts/extract_sdk_info.go --package <package>
```

### Basic Usage (Go Version)

Extract information from the monitoring-go-client package:

```bash
./scripts/extract_sdk_info \
  -package github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4@v4.1.1
```

### Basic Usage (Python Version)

```bash
python3 scripts/extract_sdk_info.py \
  --package github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4@v4.1.1
```

### Extract from API Subpackage

To extract specifically from the `api` subpackage:

**Go version:**
```bash
./scripts/extract_sdk_info \
  -package github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4@v4.1.1 \
  -api-package github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/api
```

**Python version:**
```bash
python3 scripts/extract_sdk_info.py \
  --package github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4@v4.1.1 \
  --api-package github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/api
```

### Custom Output Directory

**Go version:**
```bash
go run -mod=mod extract_sdk_info.go
  package github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4@v4.1.1 -output-dir sdk_extract_output
```

## Output

The script generates two files in the output directory:

1. **sdk_info.json** - Complete structured JSON with:
   - `config_models`: Structs that appear to be configuration models
   - `api_methods`: Methods on API types
   - `all_structs`: All struct definitions found
   - `all_methods`: All method definitions found

2. **summary.txt** - Human-readable summary of the extracted information

## Example Output Structure

```json
{
  "package": "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4@v4.1.1",
  "package_path": "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4",
  "config_models": [
    {
      "name": "ConfigModel",
      "fields": [
        {
          "name": "FieldName",
          "type": "string",
          "tag": "json:\"field_name\""
        }
      ],
      "file": "models/config.go"
    }
  ],
  "api_methods": [
    {
      "type": "*MonitoringApi",
      "name": "GetMetrics",
      "receiver": "*MonitoringApi",
      "params": [
        {
          "name": "ctx",
          "type": "context.Context"
        }
      ],
      "returns": ["*Response", "error"],
      "file": "api/monitoring_api.go"
    }
  ]
}
```

## How It Works

1. The script locates the package in `$GOPATH/pkg/mod` (or downloads it if needed)
2. It uses Go's AST (Abstract Syntax Tree) parser to analyze the source code
3. It extracts:
   - **Struct definitions** (config models) using `go/format` for proper type formatting
   - **Method definitions** on API types
4. Results are saved in JSON and text formats

The Go version uses native Go tooling (`go/ast`, `go/format`, `go/parser`) for better performance and reliability.

## Troubleshooting

### Package Not Found

If you get an error about the package not being found:
- Ensure you have internet connectivity
- Check that the package path and version are correct
- Try running `go mod download` manually first

### Go Version Issues

If you encounter Go version errors:
- Ensure Go 1.16+ is installed
- Check with `go version`

### Permission Errors

If you get permission errors:
- Ensure the script is executable: `chmod +x scripts/extract_sdk_info.py`
- Check write permissions for the output directory

