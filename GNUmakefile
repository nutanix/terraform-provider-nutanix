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
	gofmt -s .

.PHONY: sanity
sanity:
	echo "Sanity: go vet"
	go tool vet -v $(GO_FILES)
	echo "Sanity: gofmt simplify"
	gofmt -l -s $(GO_FILES)
	echo "Sanity: error check"
	go install $(TRAVIS_BUILD_DIR)/vendor/github.com/kisielk/errcheck
	errcheck -ignoretests -ignore 'github.com/hashicorp/terraform/helper/schema:Set' -ignore 'bytes:.*' -ignore 'io:Close|Write' $(GO_FILES)

.PHONY: deps
deps:
	dep ensure
	dep status