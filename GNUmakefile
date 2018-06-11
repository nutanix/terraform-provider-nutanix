GO_FILES ?= $$(go list ./... |grep -v 'vendor')

default: build

build: sanity
	go install

test: sanity
	go test -i $(GO_FILES) || exit 1
	echo $(GO_FILES) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4 -coverprofile c.out

testacc: sanity
	TF_ACC=1 go test $(GO_FILES) -v $(TESTARGS) -timeout 120m -coverprofile c.out

fmt:
	gofmt -s .

sanity:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"
	@sh -c "'$(CURDIR)/scripts/govetcheck.sh'"
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@dep ensure

.PHONY: build test testacc fmt sanity vendor-status
