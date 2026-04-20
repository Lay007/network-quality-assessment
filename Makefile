.PHONY: build test test-cover vet lint ci php-lint shell-lint web-link-check gofmt-check clean

GO ?= go
GOFMT ?= gofmt
PYTHON ?= python3
PHP ?= php
BINARY ?= build/server-sfp-sla

build:
	$(GO) build -o $(BINARY) ./cmd/server-sfp-sla

test:
	$(GO) test ./...

test-cover:
	$(GO) test -cover ./...

vet:
	$(GO) vet ./...

gofmt-check:
	@test -z "$$($(GOFMT) -l cmd internal)"

shell-lint:
	bash -n scripts/*.sh

php-lint:
	@if command -v $(PHP) >/dev/null 2>&1; then \
		find web/htdocs -name '*.php' -print0 | xargs -0 -n1 $(PHP) -l; \
	else \
		echo "skip php-lint: $(PHP) not found"; \
	fi

web-link-check:
	$(PYTHON) scripts/check_web_links.py

lint: gofmt-check vet shell-lint web-link-check php-lint

ci: test test-cover build lint

clean:
	rm -rf build
