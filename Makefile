  
BIN_DIR         ?= ./bin
PKG_NAME        ?= nozerodays
LDFLAGS         ?= "-X"
VERSION         ?=


default: build

.PHONY: build
build:
	@echo "---> Building"
	go build -o $(BIN_DIR)/$(PKG_NAME) ./cmd

.PHONY: clean
clean:
	@echo "---> Cleaning"
	rm -rf $(BIN_DIR) ./vendor

.PHONY: install
install:
	@echo "---> Installing dependencies"
	dep ensure -vendor-only

.PHONY: lint
lint:
	@echo "---> Linting"
	$(BIN_DIR)/golangci-lint run

.PHONY: setup
setup:
	@echo "--> Installing development tools"
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.16.0
	go get -u $(GOTOOLS)

.PHONY: run
run:
	@echo "---> Running"
	docker-compose -f ./deployments/docker-compose.yml up --build
