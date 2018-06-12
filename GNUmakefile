TEST_FILES ?= $$(go list ./... | grep -v 'vendor')
GO_FILES ?= $$(find . -name '*.go' | grep -v 'vendor')

default: build

.PHONY: build
build: sanity
	go install

.PHONY: test
test: sanity
	TF_ACC=1 go test $(TEST_FILES) -v $(TESTARGS) -timeout 120m -coverprofile c.out

.PHONY: cibuild
cibuild: 
	go install

.PHONY: citest
citest: 
	TF_ACC=1 go test $(TEST_FILES) -v $(TESTARGS) -timeout 120m -coverprofile c.out

.PHONY: fmt
fmt:
	gofmt -s .

.PHONY: sanity
sanity:
	# echo "==>Sanity: go vet"
	# go tool vet -v $(GO_FILES)
	# echo "==>Sanity: gofmt simplify"
	# gofmt -l -s $(GO_FILES)
	echo "==>Sanity: gometalinter"
	gometalinter --disable-all --enable=golint

.PHONY: deps
deps:
	dep ensure
	dep status

.PHONY: extras
extras:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/alecthomas/gometalinter
	cd $(GOPATH)
	./gometalinter --install
