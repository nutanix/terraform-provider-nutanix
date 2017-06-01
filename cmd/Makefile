clean:
	rm terraform.tfstate
	rm request_log
build:
	go build -o terraform-provider-example

run:
	terraform apply

all: clean build run
