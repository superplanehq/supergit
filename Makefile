.PHONY: setup setup.tools lint test check.format.go image.build

MAKEFLAGS += --no-print-directory

GO ?= go
GOTESTSUM ?= gotestsum

PKG_TEST_PACKAGES := ./...

GOTESTSUM_CMD = $(GOTESTSUM) --format short --junitfile junit-report.xml

lint:
	revive -formatter friendly -config lint.toml -exclude ./tmp/... ./...

test:
	$(GOTESTSUM_CMD) --packages="$(PKG_TEST_PACKAGES)" -- -p 1

check.format.go:
	@find . -name '*.go' -not -path './tmp/*' -print0 | xargs -0 gofmt -s -l | tee /dev/stderr | if read; then exit 1; else exit 0; fi

setup.tools:
	$(GO) install gotest.tools/gotestsum@v1.13.0
	$(GO) install github.com/mgechev/revive@v1.7.1

setup: setup.tools
	$(GO) mod download

image.build:
	bash scripts/build-image.sh "$(VERSION)" "$(ARCH)"
