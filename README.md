# CLI-first Go app

A minimal Go application where the CLI is ready on day one. Change the application, update OpenAPI, and test through the generated `appctl` CLI.

Requires Go 1.25+.

```sh
make check
go run ./cmd/server
```

In another terminal:

```sh
./bin/appctl search "create task" --json
./bin/appctl commands show tasks create --json
./bin/appctl tasks create --set title="Ship from the CLI" -o json
```

`make check` regenerates the CLI from `openapi/openapi.yaml`, builds it, and runs acceptance tests through CLI commands. `make test` only builds and tests the checked-in generated CLI.

Generated output under `internal/generated/` and `skills/appctl/` is committed with API changes. The pipeline uses [Lathe](https://github.com/lathe-cli/lathe), pinned by version.

