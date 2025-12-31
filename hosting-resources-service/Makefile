SERVICE_NAME := resources-service
CONTRACT_PATH := ../hosting-contracts
OPENAPI_SPEC := $(CONTRACT_PATH)/$(SERVICE_NAME)/openapi/openapi.yaml
RESOURCES_PROTO_DIR := $(CONTRACT_PATH)/resources-service/grpc
SERVER_GRPC_GEN_DIR := cmd/server/grpc/handlers/poolgrp/gen

.PHONY: generate
generate: generate-api generate-grpc-client
	@echo "All code generated"

.PHONY: generate-api
generate-api:
	@echo "Generating OpenAPI code..."
	@oapi-codegen -config oapi-codegen.yaml $(OPENAPI_SPEC)
	@echo "OpenAPI code generated"

.PHONY: generate-grpc-server
generate-grpc-server:
	@echo "Generating gRPC server code..."
	@protoc \
		--proto_path=$(RESOURCES_PROTO_DIR) \
		--go_out=$(SERVER_GRPC_GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(SERVER_GRPC_GEN_DIR) --go-grpc_opt=paths=source_relative \
		resources.proto
	@echo "gRPC server code generated"


.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"