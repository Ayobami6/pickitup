build:
	@go build -o bin/pickitup cmd/main.go

test:
	@go test -v ./...

run:
	@air
