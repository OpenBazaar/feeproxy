.DEFAULT_GOAL := help

##
## Docker
##

DOCKER_PROFILE ?= openbazaar
DOCKER_IMAGE_NAME ?= $(DOCKER_PROFILE)/feeproxy

docker: ## Create Docker image
	docker build -t $(DOCKER_IMAGE_NAME) .

push_docker: ## Push Docker image to registry
	docker push $(DOCKER_IMAGE_NAME)

##
## Cleanup
##
clean_docker: ## Remove Docker image
	docker rmi -f $(DOCKER_IMAGE_NAME); true

clean:  ## Clean Docker resources
	clean_docker

##
## General
##
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'