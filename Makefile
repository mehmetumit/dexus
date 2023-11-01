SCRIPTS=./scripts
CMDS="help build exec run live-reload docker-build docker-run"
IMG="dexus"
TAG="latest"

all:help
help:
	@echo 'Commands:'
	@echo 'make $(CMDS)'

.PHONY: build
build:
	@$(SCRIPTS)/build.sh

.PHONY: exec
exec: build
	@./build/dexus

.PHONY: run
run:
	@$(SCRIPTS)/run.sh

.PHONY: test
test:
	@$(SCRIPTS)/test.sh
test-coverage-html: test
	@$(SCRIPTS)/coverage_html.sh

.PHONY: live-reload
live-reload:
	@find . -type f -name '*.go' | entr -r $(SCRIPTS)/run.sh

.PHONY: docker-build
docker-build:
	@docker build -t $(IMG):$(TAG) . --build-arg="VERSION=$$(git tag -l | tail -n 1)" --build-arg="COMMIT=$$(git rev-parse --short HEAD)"

.PHONY: docker-run
docker-run:
	@docker run --rm $(IMG):$(TAG)
