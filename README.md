# CLI-first Go app

A minimal Go application where the CLI is ready on day one. Change the application, update OpenAPI, and test through the generated `appctl` CLI.

Requires Go 1.26+.

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

## Application versioning

Application releases follow [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html).
These rules apply to future releases:

- The core version is `X.Y.Z`, with nonnegative integers and no leading zeroes (clause 2).
- Published versions are immutable (clause 3). Changes require a new version. Published tags must not move or be reused.
- Prereleases use only `X.Y.Z-alpha.N`, `X.Y.Z-beta.N`, or `X.Y.Z-rc.N` (clause 9). `N` starts at 1 and has no leading zeroes.
- Versions follow SemVer precedence, not text order (clause 11). For the same core, `alpha < beta < rc < final`.
- Numeric identifiers compare numerically: `1.2.3-rc.2 < 1.2.3-rc.10 < 1.2.3`. Commit IDs and build metadata stay separate from release versions.

The Git tag is the single source of truth for the Go application version.
An application released as one unit uses `v<semver>`, such as `v1.2.3-rc.1`.
The `v` prefix belongs to the tag, not the version.
A separate `VERSION` file is unnecessary. The Go and Lathe dependency versions do not define the application version.

This starter does not include release automation or application version reporting.
Before the first application release, implement version reporting with a string variable set through `-ldflags -X`.
If that value is empty or `dev`, use the main module version from `runtime/debug.ReadBuildInfo()` as the fallback.
If build information is unavailable or the module version is empty or `(devel)`, report `dev`.
The release tag and reported binary version must identify the same release after removal of the tag prefix.

If binaries or artifacts release independently, each artifact needs a documented fixed tag prefix and an independent version sequence.
For example, one artifact can use `v1.2.3`, while another uses `agent-v0.4.0`.
Each artifact's reported version comes from its own tag, without its fixed prefix.

