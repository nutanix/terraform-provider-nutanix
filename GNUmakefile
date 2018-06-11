GO_FILES ?= $$(go list ./... |grep -v 'vendor')

default: build

build:
	echo "==> Starting make build"
	sanity
	echo "==> Starting go install"
	go install
	echo "==> Finished Build"

test:
	echo "==> Starting acceptance tests"
	sanity
	TF_ACC=1 go test $(GO_FILES) -v $(TESTARGS) -timeout 120m -coverprofile c.out
	echo "==> Finished acceptance tests"

fmt:
	echo "==> Starting gofmt with simplify flag"
	gofmt -s .
	echo "==> Finished gofmt with simplify flag"

sanity:
	echo "==> Running Sanity Checks"
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"
	@sh -c "'$(CURDIR)/scripts/govetcheck.sh'"
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"
	echo "==> Finishes Sanity Checks"

vendor-status:
	echo "==> Starting go dep ensure"
	@dep ensure
	@dep status
	echo "==> Finished go dep ensure"

.PHONY: build test fmt sanity vendor-status
