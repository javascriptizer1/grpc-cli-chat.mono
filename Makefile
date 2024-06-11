include .env

lint:
	golangci-lint run ./... --config .golangci.pipeline.yaml


generate-api:
	mkdir -p pkg/swagger
	make generate-api-access
	make generate-api-auth
	make generate-api-user

generate-api-auth:
	mkdir -p pkg/grpc/auth_v1
	protoc --proto_path api/proto/v1 --proto_path vendor.protogen \
	--go_out=pkg/grpc/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=/usr/local/bin/protoc-gen-go \
	--go-grpc_out=pkg/grpc/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=/usr/local/bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/grpc/auth_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=/usr/local/bin/protoc-gen-grpc-gateway \
	--validate_out lang=go:pkg/grpc/auth_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=/usr/local/bin/protoc-gen-validate \
	--openapiv2_out=allow_merge=true,merge_file_name=api_auth_v1:pkg/swagger \
	--plugin=protoc-gen-openapiv2=/usr/local/bin/protoc-gen-openapiv2 \
	api/proto/v1/auth.proto

generate-api-access:
	mkdir -p pkg/grpc/access_v1
	protoc --proto_path api/proto/v1 --proto_path vendor.protogen \
	--go_out=pkg/grpc/access_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=/usr/local/bin/protoc-gen-go \
	--go-grpc_out=pkg/grpc/access_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=/usr/local/bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/grpc/access_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=/usr/local/bin/protoc-gen-grpc-gateway \
	--openapiv2_out=allow_merge=true,merge_file_name=api_access_v1:pkg/swagger \
	--plugin=protoc-gen-openapiv2=/usr/local/bin/protoc-gen-openapiv2 \
	api/proto/v1/access.proto

generate-api-user:
	mkdir -p pkg/grpc/user_v1
	protoc --proto_path api/proto/v1 --proto_path vendor.protogen \
	--go_out=pkg/grpc/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=/usr/local/bin/protoc-gen-go \
	--go-grpc_out=pkg/grpc/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=/usr/local/bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/grpc/user_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=/usr/local/bin/protoc-gen-grpc-gateway \
	--validate_out lang=go:pkg/grpc/user_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=/usr/local/bin/protoc-gen-validate \
	--openapiv2_out=allow_merge=true,merge_file_name=api_user_v1:pkg/swagger \
	--plugin=protoc-gen-openapiv2=/usr/local/bin/protoc-gen-openapiv2 \
	api/proto/v1/user.proto


vendor-proto:
		@if [ ! -d vendor.protogen/buf/validate ]; then \
			mkdir -p vendor.protogen/buf/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/buf/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi

migrations-up:
	GOOSE_DBSTRING=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable goose -dir migrations/postgres postgres up
