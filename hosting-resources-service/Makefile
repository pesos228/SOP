SERVICE_NAME := resources-service
CONTRACT_PATH := ../hosting-contracts
OPENAPI_SPEC := $(CONTRACT_PATH)/$(SERVICE_NAME)/openapi/openapi.yaml

.PHONY: generate
generate: generate-api generate-gql
	@echo "All code generated"

.PHONY: generate-api
generate-api:
	@echo "Generating OpenAPI code..."
	@oapi-codegen -config oapi-codegen.yaml $(OPENAPI_SPEC)
	@echo "OpenAPI code generated"

.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"