clean:
	if [ -a terraform.tfstate ] ;\
	then \
		rm terraform.tfstate; \
	fi;	
	if [ -a terraform.tfstate.backup ]; \
	then \
		rm terraform.tfstate.backup ; \
	fi;	
	if [ -a ../virtualmachineschema ]; \
	then \
		rm -r ../virtualmachineschema ; \
	fi;	
	if [ -a ../virtualmachineconfig ]; \
	then \
		rm -r ../virtualmachineconfig ; \
	fi;	
	if [ -a terraform-provider-nutanix ]; \
	then \
		rm terraform-provider-nutanix ; \
	fi;	
	if [ -a crash.log ]; \
	then \
		rm crash.log ; \
	fi;	
getSDK:
	wget http://developer.nutanix.com/documentation/v3_sdk/html/sdk/go_sdk.tgz
	mkdir $(GOPATH)/src/nutanixV3
	tar -xzf go_sdk.tgz -C $(GOPATH)/src/nutanixV3
	rm go_sdk.tgz
	mv $(GOPATH)/src/nutanixV3/go_sdk/* $(GOPATH)/src/nutanixV3
	rm -r $(GOPATH)/src/nutanixV3/go_sdk
autoGenerate:	
	mkdir ../virtualmachineschema
	mkdir ../virtualmachineconfig
	go build ../jsonToSchema
	./jsonToSchema
	rm jsonToSchema
build:
	go build -o terraform-provider-nutanix

run:
	terraform apply

all: getSDK build run
