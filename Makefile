.PHONY=build

build-client:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/go-client learn/go-client.go

run-client: build-client
	@./bin/go-client

coverage:
	@go test -v -cover ./...

test:
	@go test -v ./...

