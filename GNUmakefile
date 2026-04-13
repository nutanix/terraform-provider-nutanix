TEST?=$$(go list ./... |grep -v 'vendor' |grep -v 'utils')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=nutanix
WEBSITE_REPO=github.com/hashicorp/terraform-website

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test --tags=unit $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	@echo "==> Running testcases..."
	@echo "TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 500m -coverprofile c.out -covermode=count"
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 500m -coverprofile c.out -covermode=count

# Acceptance tests with .env loaded (same as /ok-to-test). Loads .env from repo root before running.
# Output to ACC_TEST_LOG only; summary appended at end. Matches workflow logic from acceptance-test.yml.
# Usage:
#   make acc-test networkingv2                                         # all tests in package networkingv2 (auto-detected)
#   make acc-test networkingv2 TestAccV2NutanixSubnetResource_Basic    # single test in specific package
#   make acc-test p=networkingv2                                       # all tests in package (explicit)
#   make acc-test p=networkingv2 TestAccV2NutanixSubnetResource_Basic  # single test in specific package (explicit)
#   make acc-test v4                                                   # all TestAccV2Nutanix* tests
#   make acc-test v3                                                   # all TestAccNutanix* tests
#   make acc-test TestAccV2NutanixSubnetResource_Basic                 # single test (searches all packages)
#   make acc-test TestAccV2NutanixSubnet TestAccV2NutanixVpc           # multiple tests (combined with |)
# Note: Don't use -p flag (it's a Make built-in). Use p= or just the package name directly.
ACC_TEST_LOG ?= test_output.log
acc-test:
	@bash -c '\
		logfile="$(ACC_TEST_LOG)"; \
		: > "$$logfile"; \
		echo "==> Loading .env and running acceptance tests (output to $$logfile only; summary at end)..." >> "$$logfile"; \
		[ -f .env ] && set -a && . ./.env && set +a; \
		export TF_ACC=1 GOTRACEBACK=all; \
		args="$(filter-out acc-test,$(MAKECMDGOALS))"; \
		package_path="./..."; \
		run_flag=""; \
		if [ -n "$(p)" ]; then \
			package_path="./nutanix/services/$(p)"; \
			echo "📦 Running tests only in package: $$package_path" >> "$$logfile"; \
		fi; \
		for arg in $$args; do \
			if [ -d "nutanix/services/$$arg" ]; then \
				package_path="./nutanix/services/$$arg"; \
				echo "📦 Running tests only in package: $$package_path (auto-detected)" >> "$$logfile"; \
				continue; \
			fi; \
			case "$$arg" in \
				foundation) pattern="TestAccFoundation*" ;; \
				foundation_central) pattern="TestAccFC*" ;; \
				karbon) pattern="TestAccKarbon*" ;; \
				v3) pattern="TestAccNutanix*" ;; \
				v4) pattern="TestAccV2Nutanix*" ;; \
				lcm) pattern="TestAccV2NutanixLcm*" ;; \
				era) pattern="TestAccEra*" ;; \
				*) pattern="$$arg" ;; \
			esac; \
			if [ -n "$$run_flag" ]; then \
				run_flag="$$run_flag|$$pattern"; \
			else \
				run_flag="$$pattern"; \
			fi; \
		done; \
		if [ "$$package_path" != "./..." ] && [ -z "$$run_flag" ]; then \
			test_args="-run=."; \
		elif [ -n "$$run_flag" ]; then \
			test_args="-run=$$run_flag"; \
		else \
			test_args="-run=."; \
		fi; \
		echo "==> TESTARGS = $$test_args" >> "$$logfile"; \
		echo "==> Package path = $$package_path" >> "$$logfile"; \
		go test "$$package_path" -v $$test_args -timeout 500m -count=1 2>&1 | while IFS= read -r line; do echo "$$line" >> "$$logfile"; done; \
		if [ -f "$$logfile" ] && grep -qE "^--- (PASS|FAIL|SKIP):" "$$logfile" 2>/dev/null; then \
			"$(CURDIR)/scripts/acc-test-summary.sh" "$$logfile"; \
		fi'
	@echo "==> Log file: $(ACC_TEST_LOG)"
# Dummy target so "make acc-test v4" etc. do not fail (arguments are consumed)
%: acc-test
	@true

fmt:
	@echo "==> Fixing source code with gofmt..."
	goimports -w ./$(PKG_NAME)
	goimports -w ./client
	goimports -w ./utils


fmtcheck:
	@echo "Running fmtcheck"
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"
	@echo "fmtcheck done"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

lint: fmtcheck
	@echo "==> Checking source code against linters..."
	@GOGC=30 golangci-lint cache clean
	@GOGC=30 golangci-lint run --timeout=30m

tools:
	@echo "make: Installing tools..."
# 	GO111MODULE=on go install github.com/YakDriver/tfproviderdocs
	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell@latest
	GO111MODULE=on go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	GO111MODULE=on go install github.com/hashicorp/copywrite@latest
	GO111MODULE=on go install github.com/hashicorp/go-changelog/cmd/changelog-build@latest
	GO111MODULE=on go install github.com/katbyte/terrafmt@latest
	GO111MODULE=on go install github.com/pavius/impi/cmd/impi@latest
	GO111MODULE=on go install github.com/rhysd/actionlint/cmd/actionlint@latest
	GO111MODULE=on go install github.com/terraform-linters/tflint@latest
	GO111MODULE=on go install golang.org/x/tools/cmd/stringer@latest
	GO111MODULE=on go install mvdan.cc/gofumpt@latest

# 	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
# 	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
# 	GO111MODULE=on go install github.com/mitchellh/gox

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

# test:
# 	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -coverprofile c.out
# 	go tool cover -html=c.out

cibuild: tools
	rm -rf pkg/
	gox -output "pkg/{{.OS}}_{{.Arch}}/terraform-provider-nutanix"


citest:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -coverprofile c.out


website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.NOTPARALLEL:

.PHONY: default build test testacc acc-test fmt fmtcheck errcheck lint tools vet test-compile cibuild citest website website-lint website-test
