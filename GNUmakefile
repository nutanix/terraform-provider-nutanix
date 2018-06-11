TEST_FILES ?= $$(go list ./... | grep -v 'vendor')
GO_FILES ?= $$(find . -name '*.go' | grep -v 'vendor')

default: build

.PHONY: build
build: sanity
	go install

.PHONY: test
test: sanity
	TF_ACC=1 go test $(TEST_FILES) -v $(TESTARGS) -timeout 120m -coverprofile c.out

.PHONY: fmt
fmt:
	@gofmt -s .

.PHONY: sanity
sanity:
	go tool vet -v $(GO_FILES)
	gofmt -l -s $(GO_FILES)
	go get -u github.com/kisielk/errcheck
	errcheck -ignoretests -ignore 'github.com/hashicorp/terraform/helper/schema:Set' -ignore 'bytes:.*' -ignore 'io:Close|Write' $(GO_FILES)

.PHONY: vendor-status
vendor-status:
	@dep ensure
	@dep status
