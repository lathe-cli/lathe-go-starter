.PHONY: cli-sync cli-build test check

cli-sync:
	cp cli.yaml cmd/appctl/cli.yaml
	go run github.com/lathe-cli/lathe/cmd/lathe@v0.5.3-0.20260816043253-e33580932a42 bootstrap
	go mod tidy

cli-build:
	go build -o bin/appctl ./cmd/appctl

test: cli-build
	go test ./...

check: cli-sync test
	go vet ./...

