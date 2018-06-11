GO_FILES ?= $$(go list ./... |grep -v 'vendor')

default: build

.PHONY: build
build: sanity
	go install

.PHONY: test
test: sanity
	TF_ACC=1 go test $(GO_FILES) -v $(TESTARGS) -timeout 120m -coverprofile c.out

.PHONY: fmt
fmt:
	@gofmt -s .

.PHONY: sanity
sanity:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"
	@sh -c "'$(CURDIR)/scripts/govetcheck.sh'"
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

.PHONY: vendor-status
vendor-status:
	@dep ensure
	@dep status
