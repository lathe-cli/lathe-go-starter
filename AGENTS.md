# Agent Development Contract

This repository is CLI-first. The generated `appctl` CLI is the primary acceptance surface for application behavior.

For every API-facing change:

1. Update the Go application under `internal/app/`.
2. Update `openapi/openapi.yaml` in the same change. It is the API contract and CLI source of truth.
3. Run `make check` to regenerate, build, and test through `appctl`.
4. Commit the matching changes under `internal/generated/`, `skills/appctl/`, and `cmd/appctl/cli.yaml`.

Do not hand-edit generated files. Use direct HTTP only for transport-level checks that the generated CLI cannot express.

Before guessing a command or flag, use `appctl search "<intent>" --json`, then `appctl commands show <path...> --json`.

Keep the application on the Go standard library until a concrete requirement needs a dependency. Do not add database, authentication, deployment, or compatibility scaffolding without an explicit requirement.

## Application versioning

Use [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html) for application releases.
Apply these rules to future releases:

- Use `X.Y.Z`, with nonnegative integers and no leading zeroes (clause 2).
- Never change published contents, move published tags, or reuse published versions (clause 3). Publish changes under a new version.
- Use only `X.Y.Z-alpha.N`, `X.Y.Z-beta.N`, or `X.Y.Z-rc.N` for prereleases (clause 9). Start `N` at 1, without leading zeroes.
- Compare versions by SemVer precedence, not text order (clause 11). For the same core, `alpha < beta < rc < final`.
- Compare numeric identifiers numerically: `1.2.3-rc.2 < 1.2.3-rc.10 < 1.2.3`. Keep commit IDs and build metadata separate from release versions.

Use the Git tag as the single source of truth for the Go application version.
Use `v<semver>` for an application released as one unit, such as `v1.2.3-rc.1`.
The tag prefix is not part of the version.
Do not add a separate `VERSION` file or use the Go or Lathe dependency version as the application version.

Before the first application release, implement version reporting with standard-library support:

1. Inject the tag version into a string variable through `-ldflags -X` for release builds.
2. If the injected version is empty or `dev`, fall back to the main module version from `runtime/debug.ReadBuildInfo()`.
3. If build information is unavailable or the module version is empty or `(devel)`, report `dev`.
4. Remove the tag prefix from the reported version. Make sure that the release tag and reported version identify the same release.

If binaries or artifacts release independently, document each artifact's fixed tag prefix and independent version sequence before its first release.
For example, use `v1.2.3` for one artifact and `agent-v0.4.0` for another.
Derive each artifact's reported version from its own tag, without its fixed prefix.

