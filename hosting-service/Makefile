SERVICE_NAME := hosting-service
CONTRACT_PATH := ../hosting-contracts
OPENAPI_SPEC := $(CONTRACT_PATH)/$(SERVICE_NAME)/openapi/openapi.yaml
RESOURCES_PROTO_DIR := $(CONTRACT_PATH)/resources-service/grpc
SERVER_GRPC_GEN_DIR := internal/server/stores/servergrpc/gen

.PHONY: generate
generate: generate-api generate-gql generate-grpc-client
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

.PHONY: generate-grpc-client
generate-grpc-client:
	@echo "Generating gRPC client code..."
	@protoc \
		--proto_path=$(RESOURCES_PROTO_DIR) \
		--go_out=$(SERVER_GRPC_GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(SERVER_GRPC_GEN_DIR) --go-grpc_opt=paths=source_relative \
		resources.proto
	@echo "gRPC client code generated"

.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
	@go install github.com/99designs/gqlgen@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Dependencies installed"