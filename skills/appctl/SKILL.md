---
name: appctl
description: >
  Use when operating the appctl generated CLI. Discover commands, inspect parameters,
  check auth state, and execute API operations safely.
---

# appctl CLI

Use this skill when a user asks you to operate `appctl`, inspect its API commands, or find the right generated command for an API task.

## Workflow

1. Search for candidates with `appctl search "<intent>" --json`; use `--limit` when needed. Search is only candidate discovery.
2. Inspect the exact command with `appctl commands show <path...> --json` before executing an unfamiliar command.
3. If the command detail has `auth.required=true`, run `appctl auth status --hostname <host>` before execution. Use `http.default_hostname` when present unless the user provides `--hostname` or `$APPCTL_HOST`.
4. Execute only after flags, body, auth, HTTP path, and output hints are clear from `commands show`.

## General Commands

- `appctl commands --json`: full generated command catalog.
- `appctl commands --include-hidden --json`: include hidden generated commands.
- `appctl commands show <path...> --json`: source of truth for one command.
- `appctl commands schema --json`: catalog schema version for parser compatibility.
- `appctl search "<intent>" --json`: ranked candidate commands.

## Maintenance Commands

- `appctl --version` or `appctl -v`: print CLI build version.

## References

- Read `references/catalog.md` for the command discovery protocol and catalog field meanings.
- Read `references/modules/app.md` for the `app` module command index.

## Rules

- Do not guess flags or request body shape from command names.
- Do not execute directly from search results; confirm with `commands show` first.
- Prefer `-o json` for machine-readable command output unless the user asks for human-readable output.
- With `-o json` or `-o yaml`, branch on `error.code` and process exit status; `error.message` and `error.hint` are safe human guidance, and `error.http.status` is the only optional HTTP context.
- A configured stream pause is successful (`exit 0`); inspect the collected output field mapped from the pause event instead of treating it as an error.
- For collected streams, choose one mode: `-o json` for one stable document, `--stream` in the default output mode when catalog `output.streaming.policy.live` is present, or `-o raw` for wire events.
- Use `--file`, `--set`, or `--set-str` for JSON request bodies according to `commands show` body requirements.
- When `body.runtime_schema` is present, normal execution fetches and validates against that schema before the target request; `--dry-run` stays network-free and skips this preflight.
- For sensitive flags, prefer safe modes from `flags[].input_modes`: `--<flag>-env`, `--<flag>-file`, or `--<flag>-stdin`.
