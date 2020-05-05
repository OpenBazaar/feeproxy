.DEFAULT_GOAL := help
.PHONY: lambda deploy_lambda docker push_docker clean_docker clean help

##
## Lambda
##
LAMBDA_FILENAME ?= update_fee_estimate.zip
LAMBDA_PATH ?= lambdas
LAMBDA_DEPLOY_BUCKET ?= deploy-bucket
LAMBDA_FUNCTION_NAME ?= update-fee-estimate

lambda: ## Build lambda package
	mkdir -p dist/lambda
	GOOS=linux GOARCH=amd64 go build -o dist/lambda/main ./lambda
	chmod +x dist/lambda/main
	cd dist/lambda && zip -r $(LAMBDA_FILENAME) main

deploy_lambda: ## Deploy built lambda artifact
	aws s3api put-object --bucket $(LAMBDA_DEPLOY_BUCKET) --key $(LAMBDA_PATH)/$(LAMBDA_FILENAME) --body dist/lambda/$(LAMBDA_FILENAME)

update_lambda: ## Update Lambda code
	aws lambda update-function-code --function-name $(LAMBDA_FUNCTION_NAME) --zip-file fileb://dist/lambda/$(LAMBDA_FILENAME)

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
