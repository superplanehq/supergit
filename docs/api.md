# supergit HTTP API

Base URL: `http://<host>:<port>/api`

All JSON responses use `Content-Type: application/json` unless noted otherwise.

## Health

```http
GET /health
```

Returns `200` with body `ok`. This endpoint is outside `/api`.

## Errors

Failed requests return JSON:

```json
{
  "error": "human-readable message"
}
```

Clients should rely on HTTP status codes, not exact error strings.

| Status | Meaning |
|--------|---------|
| `400` | Invalid request (bad path, invalid commit metadata, malformed body) |
| `404` | Repository or Git object not found |
| `409` | Expected head SHA does not match the current branch head |
| `413` | File or commit exceeds configured size limits |
| `500` | Internal server error |

## Repository IDs in URLs

Repository IDs are caller-defined and may contain slashes (for example `acme/widgets`). Encode slashes in URL paths:

```text
acme/widgets
→ acme%2Fwidgets
```

The server URL-decodes the `{id}` path parameter before resolving the repository.

Repository IDs must be relative, non-empty paths without `..`, `.git` segments, or null bytes.

## Git clone (Smart HTTP)

Supergit exposes bare repositories for `git clone`, `git fetch`, and `git push` under:

```text
http://<host>:<port>/git/{repository-id}.git
```

Example:

```bash
git clone http://localhost:8080/git/acme/widgets.git
```

Configure `SUPERGIT_PUBLIC_URL` (for example `http://localhost:8080/git`) so repository metadata includes `clone_url`.

---

## Repositories

### List repositories

```http
GET /api/repos
```

**Response `200`**

```json
{
  "repos": [
    {
      "id": "acme/widgets",
      "default_branch": "main"
    }
  ],
  "next_cursor": "",
  "has_more": false
}
```

Pagination fields are reserved for future use; listing currently returns all repositories.

### Create repository

```http
POST /api/repos
Content-Type: application/json
```

**Request body**

```json
{
  "id": "acme/widgets",
  "default_branch": "main"
}
```

| Field | Required | Description |
|-------|----------|-------------|
| `id` | yes | Repository identifier |
| `default_branch` | no | Branch created on `git init` (defaults to server `SUPERGIT_DEFAULT_BRANCH`) |

If the repository already exists, the call succeeds and returns the existing repository metadata.

**Response `201`**

```json
{
  "id": "acme/widgets",
  "default_branch": "main"
}
```

Creating a repository initializes an empty bare repo. Initial files (for example `README.md`) are added by the client via the commits endpoint.

### Get repository

```http
GET /api/repos/{id}
```

**Response `200`**

```json
{
  "id": "acme/widgets",
  "default_branch": "main"
}
```

### Delete repository

```http
DELETE /api/repos/{id}
```

Deletes the bare repository from disk. Returns `204` with no body. Missing repositories are treated as success.

---

## Files

### List files

```http
GET /api/repos/{id}/files
GET /api/repos/{id}/files?ref={ref}
```

| Query | Description |
|-------|-------------|
| `ref` | Branch, tag, or commit SHA (defaults to the repository default branch) |

**Response `200`**

```json
{
  "paths": ["README.md", "src/main.go"],
  "ref": "main"
}
```

Returns all file paths in the tree at the given ref, recursively.

### Get file

```http
GET /api/repos/{id}/files?path={path}
GET /api/repos/{id}/files?path={path}&ref={ref}
```

| Query | Required | Description |
|-------|----------|-------------|
| `path` | yes | Repository-relative file path |
| `ref` | no | Branch, tag, or commit SHA |

**Response `200`**

- `Content-Type: application/octet-stream`
- Body: raw file bytes

---

## Commits

### List commits

```http
GET /api/repos/{id}/commits
GET /api/repos/{id}/commits?branch={branch}&limit={limit}
```

| Query | Default | Description |
|-------|---------|-------------|
| `branch` | repo default branch | Branch to walk |
| `limit` | `20` | Maximum number of commits |

**Response `200`**

```json
{
  "commits": [
    {
      "commit_sha": "9e6a46e7f5affc39101a2bbe03b85ea0c934cdef",
      "tree_sha": "abc123...",
      "message": "Initialize repository",
      "author": {
        "name": "SuperPlane",
        "email": "bot@superplane.local"
      }
    }
  ],
  "next_cursor": "",
  "has_more": false
}
```

Commits are ordered newest first.

### Create commit

```http
POST /api/repos/{id}/commits
Content-Type: application/x-ndjson
```

Creates a commit by streaming file operations in [NDJSON](https://github.com/ndjson/ndjson-spec) format.

The body is a newline-delimited stream:

1. First line: `metadata` object describing the commit and file operations
2. Following lines: `blob_chunk` objects with base64-encoded file content

#### Metadata line

```json
{
  "metadata": {
    "target_branch": "main",
    "base_branch": "",
    "expected_head_sha": "",
    "commit_message": "Update files",
    "author": {
      "name": "Jane Doe",
      "email": "jane@example.com"
    },
    "files": [
      {
        "path": "README.md",
        "operation": "upsert",
        "content_id": "blob-1",
        "mode": "100644"
      },
      {
        "path": "old.txt",
        "operation": "delete",
        "content_id": "blob-2"
      }
    ]
  }
}
```

| Field | Description |
|-------|-------------|
| `target_branch` | Branch to update |
| `base_branch` | Optional base branch when creating `target_branch` |
| `expected_head_sha` | Optional optimistic concurrency check; commit fails with `409` on mismatch |
| `commit_message` | Commit message (required) |
| `author.name` / `author.email` | Commit author (required) |
| `files[].operation` | `upsert`, `add`, `update`, or `delete` |
| `files[].content_id` | References a `blob_chunk` stream (required for every file entry) |
| `files[].mode` | Git file mode for upserts (e.g. `100644`); optional |

#### Blob chunk lines

```json
{"blob_chunk":{"content_id":"blob-1","data":"SGVsbG8=","eof":true}}
```

| Field | Description |
|-------|-------------|
| `content_id` | Must match a `content_id` from the metadata `files` list |
| `data` | Base64-encoded chunk payload |
| `eof` | Set to `true` on the final chunk for this `content_id` |

Each blob may be sent in multiple chunks; decoded content is concatenated until `eof: true`.

#### Example request

```bash
curl -X POST "http://localhost:8080/api/repos/acme%2Fwidgets/commits" \
  -H "Content-Type: application/x-ndjson" \
  --data-binary @- <<'EOF'
{"metadata":{"target_branch":"main","commit_message":"Add README","author":{"name":"SuperPlane","email":"bot@superplane.local"},"files":[{"path":"README.md","operation":"upsert","content_id":"blob-1","mode":"100644"}]}}
{"blob_chunk":{"content_id":"blob-1","data":"","eof":true}}
EOF
```

**Response `201`**

```json
{
  "commit": {
    "commit_sha": "9e6a46e7f5affc39101a2bbe03b85ea0c934cdef"
  },
  "result": {
    "branch": "main",
    "new_sha": "9e6a46e7f5affc39101a2bbe03b85ea0c934cdef",
    "old_sha": "",
    "success": true
  }
}
```

If the commit produces no tree changes, `new_sha` may equal `old_sha`.

### Get commit

```http
GET /api/repos/{id}/commit?sha={sha}
```

| Query | Required | Description |
|-------|----------|-------------|
| `sha` | yes | Branch name, tag, or commit SHA |

**Response `200`**

```json
{
  "commit_sha": "9e6a46e7f5affc39101a2bbe03b85ea0c934cdef",
  "tree_sha": "abc123...",
  "message": "Add README",
  "author": {
    "name": "SuperPlane",
    "email": "bot@superplane.local"
  }
}
```

Branch names (for example `main`) are resolved to the branch head commit.

---

## Path rules

- File paths must be relative, non-empty, and must not contain `..`, `.git`, or null bytes.
- Paths configured in `SUPERGIT_RESERVED_PATHS` and anything under them are rejected.

## Size limits

Configured via `SUPERGIT_MAX_FILE_BYTES` and `SUPERGIT_MAX_COMMIT_BYTES`. Oversized requests return `413`.
