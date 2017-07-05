clean:
	rm terraform.tfstate
getSDK:
	wget http://developer.nutanix.com/documentation/v3_sdk/html/sdk/go_sdk.tgz
	mkdir $(GOPATH)/src/nutanixV3
	tar -xzf go_sdk.tgz -C $(GOPATH)/src/nutanixV3
	rm go_sdk.tgz
	mv $(GOPATH)/src/nutanixV3/go_sdk/* $(GOPATH)/src/nutanixV3
	rm -r $(GOPATH)/src/nutanixV3/go_sdk
autoGenerate:	
	go build ../jsonToStruct
	./jsonToStruct < ../jsonToStruct/exampleJSON/vmlist.json -structName=VMList
	./jsonToStruct < ../jsonToStruct/exampleJSON/vm.json
	rm jsonToStruct
build:
	go build -o terraform-provider-nutanix

run:
	terraform apply

all: getSDK build run
