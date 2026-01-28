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
# Output streams to terminal in real time and to ACC_TEST_LOG (default: test_output.log); summary appended at end.
# Usage:
#   make acc-test TestAccV2NutanixOvaVmDeployResource_DeployVMFromOva   # single test (vmmv2)
#   make acc-test PKG=vmmv2                                             # all tests in vmmv2
#   make acc-test v4                                                   # all TestAccV2Nutanix* tests
#   make acc-test v3                                                   # all TestAccNutanix* tests
ACC_ARG := $(firstword $(filter-out acc-test,$(MAKECMDGOALS)))
ACC_TEST_LOG ?= test_output.log
acc-test:
	@bash -c '\
		arg="$(ACC_ARG)"; logfile="$(ACC_TEST_LOG)"; \
		: > "$$logfile"; \
		echo "==> Loading .env and running acceptance tests (output to $$logfile only; summary at end)..." >> "$$logfile"; \
		[ -f .env ] && set -a && . ./.env && set +a; \
		export TF_ACC=1 GOTRACEBACK=all; \
		run_pattern="."; pkg=""; \
		case "$$arg" in \
			v4) run_pattern="TestAccV2Nutanix*" ;; \
			v3) run_pattern="TestAccNutanix*" ;; \
			foundation) run_pattern="TestAccFoundation*" ;; \
			foundation_central) run_pattern="TestAccFC*" ;; \
			karbon) run_pattern="TestAccKarbon*" ;; \
			lcm) run_pattern="TestAccV2NutanixLcm*" ;; \
			era) run_pattern="TestAccEra*" ;; \
			PKG=*) pkg="$${arg#PKG=}"; run_pattern="." ;; \
			"") run_pattern="." ;; \
			*) run_pattern="$$arg"; [ -d nutanix/services/vmmv2 ] && pkg="vmmv2" ;; \
		esac; \
		if [ -n "$$pkg" ]; then \
			(cd nutanix/services/$$pkg && go test . -v -run="$$run_pattern" -timeout 500m -count=1 2>&1) | while IFS= read -r line; do echo "$$line" >> "$$logfile"; done; \
		else \
			go test ./... -v -run="$$run_pattern" -timeout 500m -count=1 2>&1 | while IFS= read -r line; do echo "$$line" >> "$$logfile"; done; \
		fi; \
		if [ -f "$$logfile" ] && grep -qE "^--- (PASS|FAIL|SKIP):" "$$logfile" 2>/dev/null; then \
			"$(CURDIR)/scripts/acc-test-summary.sh" "$$logfile"; \
		fi'
	@echo "==> Log file: $(ACC_TEST_LOG)"
# Dummy target so "make acc-test v4" etc. do not fail (argument is consumed by ACC_ARG)
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
