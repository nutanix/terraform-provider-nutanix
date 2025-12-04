
# Extract the SDK information
go run -mod=mod code_gen/extract_sdk_info.go -package="github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4@v4.1.1" -keyword="EntityGroup"(Keyword is fior entity you need to autoigenerate the Terraform Code.)

== sdk_info.json file will be created with all the SDK information for the entity.

# Invoke the cursor and generate Dev code, Test Code, Documentation, Examples
cd code_gen && python3 cursor_auto_generate.py sdk_extract_output/sdk_info.json
== Generate Dev code, Examples, TestCode, Documentation