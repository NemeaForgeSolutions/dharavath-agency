# Dharavath Agency

A high-performance real estate discovery and advisory platform built with **Go**, **HTMX**, **Alpine.js**, and **Tailwind CSS v4**.

Features server-side rendering with Go's standard library `net/http` and `html/template`, reactive partial updates via HTMX, lightweight client state with Alpine.js, and self-contained static assets compiled directly into the binary using Go's `embed.FS`.

---

## Tech Stack

- **Backend**: Go 1.24+ (`net/http`, `html/template`, `embed.FS`)
- **Frontend**: HTMX (dynamic swaps), Alpine.js (modals, dark mode, calculators)
- **Styling**: Tailwind CSS v4 CLI
- **Data**: In-memory repository seeded from `seed.json`
- **Deployment**: Multi-stage Docker (Alpine Linux runtime)

---

## Getting Started

### Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- [Node.js 20+](https://nodejs.org/) & `npm` (for Tailwind CSS)
- [Docker](https://www.docker.com/) *(optional)*

### Quick Start

1. **Install frontend dependencies**:
   ```bash
   npm install
   ```

2. **Build stylesheets**:
   ```bash
   npm run build:css
   # or: make css
   ```

3. **Start the server**:
   ```bash
   make dev
   # or: go run ./cmd/server
   ```

4. Open [http://localhost:8080](http://localhost:8080) in your browser.

---

## Development Commands

| Command | Description |
|---|---|
| `make dev` | Run server directly with `go run` |
| `make build` | Compile static binary to `bin/server` |
| `make run` | Build and run the compiled binary |
| `make css` | Compile and minify production Tailwind CSS |
| `npm run watch:css` | Watch and recompile CSS during development |
| `make test` | Run test suite (`go test -v ./...`) |
| `make clean` | Remove compiled binaries and Go cache |
| `make docker-build` | Build the multi-stage Docker image |
| `make docker-run` | Run container locally on port 8080 |

---

## Configuration

Configure the server via environment variables or a `.env` file:

| Variable | Default | Description |
|---|---|---|
| `HOST` | `0.0.0.0` | Bind IP address |
| `PORT` | `8080` | HTTP listening port |
| `APP_ENV` | `development` | Environment mode (`development` or `production`) |
| `READ_TIMEOUT_SEC` | `10` | Maximum request read duration in seconds |
| `WRITE_TIMEOUT_SEC` | `15` | Maximum response write duration in seconds |
| `IDLE_TIMEOUT_SEC` | `60` | Maximum keep-alive wait time in seconds |

---

## Project Structure

```text
├── cmd/server/          # Application entrypoint and server lifecycle
├── configs/             # Environment configuration & company metadata
├── internal/
│   ├── delivery/http/   # HTTP handlers, middleware, and router
│   ├── domain/          # Core entities, models, and interfaces
│   ├── repository/      # In-memory data store and seed.json
│   ├── service/         # Business logic services
│   └── view/            # Template engine and view models
├── web/
│   ├── static/          # Source & compiled CSS, icons, manifest
│   ├── template/        # HTML templates (layouts, pages, partials, admin)
│   └── web.go           # embed.FS declarations for binary embedding
├── Dockerfile           # Multi-stage production container
└── Makefile             # Development and build shortcuts
```

---

## License

Proprietary. Developed for **Dharavath Agency** by **NemeaForge Solutions**.