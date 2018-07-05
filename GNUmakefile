TEST?=$$(go list ./... |grep -v 'vendor' |grep -v 'utils')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=nutanix
WEBSITE_REPO=github.com/hashicorp/terraform-website

default: build

build: deps
	go install

test: deps
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -coverprofile c.out
	go tool cover -html=c.out

cibuild:
	go build
#	env GOOS=darwin GOARCH=amd64 go build
#	env GOOS=windows GOARCH=amd64 go build
#	env GOOS=linux GOARCH=amd64 go build
#	env GOOS=linux GOARCH=ppc64 go build

citest:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -coverprofile c.out

fmt:
	gofmt -s .

extrasanity:
	echo "==>sanity: golangci-lint"
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $(GOPATH)/bin v1.9.1
	$(GOPATH)/bin/golangci-lint run

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep status

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.NOTPARALLEL:

.PHONY: build test cibuild citest fmt sanity deps tools website website-test