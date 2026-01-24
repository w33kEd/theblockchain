build:
	go build -o ./bin/theblockchain

run: build
	./bin/theblockchain

test:
	go test -v ./...