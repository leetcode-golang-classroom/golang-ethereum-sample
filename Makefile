.PHONY=build

build-client:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/go-client learn/go-client/go-client.go

build-wallet:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/go-wallet learn/go-wallet/go-wallet.go

build-keystore:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/go-keystore learn/go-keystore/go-keystore.go

build-transaction:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/go-transaction learn/go-transaction/go-transaction.go

run-client: build-client
	@./bin/go-client

run-wallet: build-wallet
	@./bin/go-wallet

run-keystore: build-keystore
	@./bin/go-keystore

run-transaction: build-transaction
	@./bin/go-transaction

coverage:
	@go test -v -cover ./...

test:
	@go test -v ./...

