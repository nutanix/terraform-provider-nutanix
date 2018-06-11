TEST ?= $$(go list ./... |grep -v 'vendor')
GO_FILES ?= $$(find . -name '*.go' |grep -v vendor)

default: build

build: fmtcheck
	vetcheck
	go install

test: fmtcheck
	vetcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4 -coverprofile c.out

testacc: fmtcheck
	vetcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -coverprofile c.out

vetcheck:
	go tool vet -v $(GO_FILES)

fmt:
	gofmt -s .

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@dep ensure

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc fmt fmtcheck errcheck vetcheck vendor-status test-compile
