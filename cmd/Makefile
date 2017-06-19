clean:
	rm terraform.tfstate
build:
	go build -o terraform-provider-nutanix

run:
	terraform apply

all: clean build run
