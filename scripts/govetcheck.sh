#!/usr/bin/env bash

# Check go vet
echo "==> Checking that code complies with go vet requirements..."
govet_files="$(go vet "$(find . -name '*.go' | grep -v vendor)")"
if [[ -n ${govet_files} ]]; then
    echo "Vet found suspicious constructs. Please check the reported constructs";
    echo "and fix them if necessary"
    echo 'go vet found issues on the following files:'
    echo "${govet_files}"
    exit 1
fi

exit 0