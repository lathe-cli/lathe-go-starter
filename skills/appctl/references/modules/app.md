# Module `app`

## Source

- Backend: `openapi3`
- Default hostname: `http://127.0.0.1:3000`
- Repository: `unknown`
- Pinned tag: ``unknown``
- Files: `openapi.yaml`

## Health

### `appctl health get`

- Summary: Check application health
- HTTP: `GET /health`
- Auth: public
- Body: none
- Flags: none
- Output: response media `application/json`

## Tasks

### `appctl tasks create`

- Summary: Create a task
- HTTP: `POST /tasks`
- Auth: public
- Body: required; media type `application/json`
- Flags: none

### `appctl tasks delete`

- Summary: Delete a task
- HTTP: `DELETE /tasks/{id}`
- Auth: public
- Body: none
- Flags:
  - `--id` (path, required): id

### `appctl tasks get`

- Summary: Get a task
- HTTP: `GET /tasks/{id}`
- Auth: public
- Body: none
- Flags:
  - `--id` (path, required): id
- Output: response media `application/json`

### `appctl tasks list`

- Summary: List tasks
- HTTP: `GET /tasks`
- Auth: public
- Body: none
- Flags: none
- Output: columns `id`, `completed`, `title`; response media `application/json`

### `appctl tasks update`

- Summary: Update a task
- HTTP: `PATCH /tasks/{id}`
- Auth: public
- Body: required; media type `application/json`
- Flags:
  - `--id` (path, required): id
- Output: response media `application/json`

