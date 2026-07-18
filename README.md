# Omega Learn

Omega Learn is a small Go website for browsing the Markdown knowledge base in `docs/`.
The Markdown files remain the source of truth and are never rewritten by the application.

## Run locally

```bash
go run .
```

Open <http://localhost:8080>. The server discovers every `.md` file below `docs/` when it starts.

You can change the listener or content directory with flags:

```bash
go run . -addr :3000 -docs ./docs
```

The equivalent environment variables are `OMEGA_ADDR` and `OMEGA_DOCS`.

## Verify

```bash
go test ./...
go vet ./...
```
