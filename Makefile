.PHONY: cli-sync cli-build test check

cli-sync:
	cp cli.yaml cmd/appctl/cli.yaml
	go tool lathe bootstrap
	go mod tidy

cli-build:
	go build -o bin/appctl ./cmd/appctl

test: cli-build
	go test ./...

check: cli-sync test
	go vet ./...

