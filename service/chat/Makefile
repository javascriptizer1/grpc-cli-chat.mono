include .env

PROTOC := protoc
PROTOC_GEN_GO := ../../bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := ../../bin/protoc-gen-go-grpc
PROTOC_GEN_GRPC_GATEWAY := ../../bin/protoc-gen-grpc-gateway
PROTOC_GEN_VALIDATE := ../../bin/protoc-gen-validate
PROTOC_GEN_OPENAPIV2 := ../../bin/protoc-gen-openapiv2

PROTO_SRC_DIR := api/proto/v1
PROTO_VENDOR_DIR := ../../vendor.protogen
SWAGGER_DST_DIR := ../../pkg/swagger

lint:
	golangci-lint run ./... --config ../../.golangci.pipeline.yaml

generate-api:
	mkdir -p $(SWAGGER_DST_DIR)
	make generate-api-chat

generate-api-chat:
	mkdir -p ../../pkg/grpc/chat_v1
	$(PROTOC) --proto_path=$(PROTO_SRC_DIR) --proto_path=$(PROTO_VENDOR_DIR) \
		--go_out=../../pkg/grpc/chat_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--go-grpc_out=../../pkg/grpc/chat_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) \
		--grpc-gateway_out=../../pkg/grpc/chat_v1 --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(PROTOC_GEN_GRPC_GATEWAY) \
		--validate_out lang=go:../../pkg/grpc/chat_v1 --validate_opt=paths=source_relative \
		--plugin=protoc-gen-validate=$(PROTOC_GEN_VALIDATE) \
		$(PROTO_SRC_DIR)/chat.proto

migrations-up:
	GOOSE_DBSTRING=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable goose -dir migrations/postgres postgres up
