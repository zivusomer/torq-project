#!/usr/bin/env bash
set -euo pipefail

apply_standard_fixes() {
	echo "Applying go fmt and go mod tidy..."
	go fmt ./...
	go mod tidy
}

run_checks() {
	./scripts/lint-custom.sh || return 1
	go vet ./... || return 1
	go test ./... || return 1
	return 0
}

apply_standard_fixes

set +e
run_checks
status=$?
set -e

if [ "$status" -eq 0 ]; then
	exit 0
fi

echo "Lint failed; applying automatic fixes and retrying once..."
go run ./scripts/fixfuncorder
apply_standard_fixes

set +e
run_checks
status=$?
set -e
exit "$status"
