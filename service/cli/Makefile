lint:
	golangci-lint run ./... --config ../../.golangci.pipeline.yaml

run:
	go run ./cmd/app/main.go

build:
	go build -o bin cmd/app/main.go

build-linux:
	GOARCH=amd64 GOOS=linux go build -o bin_lnx_64 cmd/app/main.go