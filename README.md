# supergit

Git-backed file storage for SuperPlane, exposed over HTTP. supergit manages bare Git repositories on disk and lets callers create repos, list and read files, and write commits without running Git in the client.

## Requirements

- Go 1.26+
- `git` on `PATH` (used for all repository operations)

## API reference

See [docs/api.md](docs/api.md) for endpoint details, request/response shapes, and the NDJSON commit format.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `SUPERGIT_ROOT` | `/var/lib/supergit/repos` | Directory for bare Git repositories |
| `SUPERGIT_PORT` | `8080` | HTTP listen port |
| `SUPERGIT_PUBLIC_URL` | _(empty)_ | Base URL for `git clone` (for example `http://localhost:8080/git`) |
| `SUPERGIT_DEFAULT_BRANCH` | `main` | Default branch when a repo is created without one |
| `SUPERGIT_MAX_FILE_BYTES` | `10485760` (10 MiB) | Maximum size of a single file in a commit |
| `SUPERGIT_MAX_COMMIT_BYTES` | `26214400` (25 MiB) | Maximum total blob size per commit |
| `SUPERGIT_RESERVED_PATHS` | _(empty)_ | Comma-separated top-level paths callers cannot write (for example `.superplane`) |

## Repository IDs

Repository IDs are opaque strings chosen by the caller. They may contain slashes (for example `acme/widgets`) and are mapped to bare Git repositories under `SUPERGIT_ROOT`. IDs must be relative paths: non-empty, without `..`, `.git` segments, or null bytes.

When an ID contains slashes, URL-encode it in request paths (see [docs/api.md](docs/api.md)).

Paths listed in `SUPERGIT_RESERVED_PATHS` and anything under them cannot be written by callers.

## Git clone

Supergit serves repositories over Git Smart HTTP at `/git/{repository-id}.git`.

Set `SUPERGIT_PUBLIC_URL` to the externally reachable git base URL (for example `http://localhost:8080/git`). Repository metadata responses include `clone_url` when this is set.

```bash
git clone http://localhost:8080/git/acme/widgets.git
```

Repository IDs with slashes use normal URL path segments (no encoding) in git clone URLs.
