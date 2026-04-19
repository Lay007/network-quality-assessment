.PHONY: build test lint ci php-lint shell-lint web-link-check clean

BINARY ?= build/Server_SFP_SLA

build:
	go build -o $(BINARY) ./cmd/server-sfp-sla

test:
	go test ./...

lint: shell-lint web-link-check php-lint

ci: test build lint

gofmt-check:
	@test -z "$$(gofmt -l cmd internal)"

shell-lint:
	bash -n scripts/*.sh

php-lint:
	find web/htdocs -name '*.php' -print0 | xargs -0 -n1 php -l

web-link-check:
	python scripts/check_web_links.py

clean:
	rm -rf build
