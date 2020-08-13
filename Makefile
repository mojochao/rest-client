#==============================================================================
#
# Makefile for building Docker images and pushing them to AWS ECR registries.
#
#==============================================================================

# Set app identity.
APP=rest-client
VERSION := $(shell cat VERSION | tr -d '\n')

# Discover OS for default static build configuration.
OS = shell("uname")
ifeq ($(OS),Linux)
    OS=linux
else
    OS=darwin
endif

# Set GOOS for static builds using discovered OS if not provided.
GOOS ?= $(OS)

# Set Docker image build configuration.
DOCKER_FILE ?= Dockerfile

# Set Docker image identity.
IMAGE = github.com/mojochao/$(APP)
TAG ?= latest

#==============================================================================
#
# Define help targets with descriptions provided in trailing `##` comments.
#
# Note that the '## description' is used in generating documentation when 'make'
# is invoked with no arguments.
#
# See https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html for
# additional details on how this works.
#
#==============================================================================

.PHONY: help

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

#==============================================================================
#
# Define Golang targets.
#
#==============================================================================

run: ## Run the application
	@echo 'running $(APP)'
	go run main.go

prepare: ## Prepare govvv compiler wrapper to embed build identity info
	@echo 'Installing govvv'
	go get github.com/ahmetb/govvv

build: ## Build the application
	@echo 'building $(APP)'
	CGO_ENABLED=0 GOOS=$(GOOS) govvv build -a -installsuffix cgo -ldflags '-extldflags "-static"' -pkg $(IMAGE)/identity -o $(APP) .

lint: ## Lint the application
	@echo 'linting $(APP)'
	golint ./...

test: ## Run all tests
	@echo 'testing $(APP)'
	go test -v ./...

clean:  ## Clean build artifacts
	@echo 'cleaning $(APP)'
	rm -f $(APP)

#==============================================================================
#
# Define docker targets.
#
#==============================================================================

docker-build: ## Build docker image with $TAG (default: latest)
	@echo 'getting govvv'
	@echo 'building docker image $(IMAGE):$(TAG)'
	docker build -f $(DOCKER_FILE) -t $(IMAGE):$(TAG) .

docker-run: ## Run docker image with $TAG (default: latest)
	@echo 'running docker image $(IMAGE):$(TAG)'
	@echo 'site will be served at http://localhost:8080/'
	docker run -it $(IMAGE):$(TAG)
