SERVICE_NAME := hosting-service
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

.PHONY: generate-gql
generate-gql:
	@echo "Generating GraphQL code..."
	@go get github.com/99designs/gqlgen@v0.17.81
	@go run github.com/99designs/gqlgen generate
	@echo "GraphQL code generated"

.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
	@go install github.com/99designs/gqlgen@latest
	@echo "Dependencies installed"