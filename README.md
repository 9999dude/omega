# Omega Learn

Omega Learn is a small, self-contained Go website for browsing a personal Markdown knowledge base. It turns the files under [`docs/`](docs/) into a searchable documentation library while keeping those Markdown files as the source of truth.

The site uses the Gruvbox Dark Soft theme and includes responsive navigation, generated tables of contents, Mermaid diagrams, and syntax highlighting for Go, YAML, and JSON.

## Features

- Discovers every `.md` file below `docs/` at startup
- Organizes articles by their top-level directory
- Generates article titles, summaries, reading times, breadcrumbs, and navigation
- Builds an “On this page” menu from level-two and level-three headings
- Rewrites relative links between local Markdown files to website routes
- Provides full-text search with `⌘ K` or `Ctrl K`
- Renders GitHub Flavored Markdown, tables, footnotes, and typographic substitutions
- Renders Mermaid diagrams with a pop-up viewer and zoom controls
- Highlights fenced Go, YAML/YML, and JSON code blocks
- Serves templates, styles, scripts, Mermaid, and Highlight.js from the compiled binary
- Adds a health-check endpoint and basic browser security headers

## Requirements

- Go version compatible with [`go.mod`](go.mod)

No Node.js installation or frontend build step is required.

## Run locally

```bash
go run .
```

Open [http://localhost:8080](http://localhost:8080). The server reads the Markdown library when it starts, so restart it after adding or editing a document.

To use a different address or content directory:

```bash
go run . -addr :3000 -docs ./docs
```

The same settings can be supplied with environment variables:

```bash
OMEGA_ADDR=:3000 OMEGA_DOCS=./docs go run .
```

| Setting | Flag | Environment variable | Default |
| --- | --- | --- | --- |
| Listen address | `-addr` | `OMEGA_ADDR` | `:8080` |
| Markdown directory | `-docs` | `OMEGA_DOCS` | `docs` |

Command-line flags take precedence over environment-variable defaults.

## Build a binary

```bash
go build -o omega-learn .
./omega-learn
```

The web templates and static assets are embedded in the binary. The `docs/` directory remains external so the same binary can serve a different knowledge base with `-docs`.

## Add content

Create Markdown files anywhere below `docs/`:

```text
docs/
├── algorithms/
│   └── recursion.md
├── k8s/
│   └── storage/
│       └── storage-concepts.md
└── sre/
    └── role.md
```

The first path segment becomes the library section. For example, `docs/k8s/storage/storage-concepts.md` is served at:

```text
/docs/k8s/storage/storage-concepts
```

A file named `README.md` maps to its containing directory. For example, `docs/k8s/README.md` is available at `/docs/k8s`.

For the best result, begin each document with a level-one heading and an introductory paragraph:

```markdown
# Storage concepts

A practical guide to storage lifetimes and Kubernetes storage objects.

## Container storage

Article content starts here.
```

The level-one heading is displayed as the article title and omitted from the rendered article body to avoid duplication. The first suitable paragraph becomes the article summary.

### Link between documents

Use ordinary relative Markdown links:

```markdown
[Read about workload configuration](../config/config-concepts.md)
[Jump to a section](../config/config-concepts.md#configmaps)
```

Links to existing local `.md` files are automatically converted to their `/docs/...` routes.

### Render Mermaid diagrams

Use a fenced block labeled `mermaid`:

````markdown
```mermaid
flowchart LR
    A[Markdown] --> B[Goldmark]
    B --> C[Omega Learn]
```
````

Rendered diagrams can be clicked to open the diagram viewer. The viewer supports zoom buttons, reset, scrolling, keyboard zoom shortcuts, and `Esc` to close.

### Highlight source code

Use fenced blocks labeled `go`, `yaml`, `yml`, or `json`:

````markdown
```go
// Exact question: How do you add a syntax-highlighted Go code example to Markdown?
//
// Possible answer: Use `main` to format and write the demonstrated value to standard output.
//
// Output format: `main` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
func main() {
    fmt.Println("learn deeply")
}

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

```yaml
apiVersion: v1
kind: ConfigMap
```

```json
{
  "enabled": true
}
```
````

## Project structure

```text
.
├── docs/                   Markdown knowledge base
├── web/
│   ├── static/             Gruvbox theme, interactions, and vendored libraries
│   └── templates/          Home and article templates
├── main.go                 Flags and HTTP server startup
├── site.go                 Content loading, Markdown rendering, search, and routes
├── site_test.go            Rendering, routing, search, and error-path tests
├── go.mod
└── README.md
```

## HTTP routes

| Route | Purpose |
| --- | --- |
| `/` | Library home page |
| `/docs/{path}` | Rendered Markdown article |
| `/api/search?q={query}` | JSON search results |
| `/assets/{path}` | Embedded frontend assets |
| `/healthz` | Health check; returns `ok` |

## Verify changes

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

These checks exercise document discovery, Markdown rendering, local-link rewriting, routes, search, and empty-library handling.
