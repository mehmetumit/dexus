SCRIPTS=./scripts
CMDS="help build run live-reload docker-build docker-run"
IMG="dexus"
TAG="latest"

all:help
help:
	@echo 'Commands:'
	@echo 'make $(CMDS)'
.PHONY: build
build:
	@$(SCRIPTS)/build.sh
.PHONY: run
run:
	@go run cmd/main.go
.PHONY: test
test:
	@$(SCRIPTS)/test.sh
test-coverage-html: test
	@$(SCRIPTS)/coverage_html.sh
.PHONY: live-reload
live-reload:
	@find . -type f -name '*.go' | entr -r go run cmd/main.go
.PHONY: docker-build
docker-build:
	@docker build -t $(IMG):$(TAG) .
.PHONY: docker-run
docker-run:
	@docker run --rm $(IMG):$(TAG)
