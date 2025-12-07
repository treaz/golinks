# AGENTS.md

## Project Overview

Golinks is an internal URL shortener service that redirects short `go/keyword` links to full URLs.
- **Backend**: Go (Golang) using Chi router and GORM.
- **Frontend**: Single-page Vue.js application (Vue 3, CDN-hosted) embedded in the Go binary.
- **Database**: SQLite (default), supports MySQL and PostgreSQL.

## Setup Commands

- **Prerequisites**: Go 1.23+, Docker (optional).
- **Run Locally**:
  ```bash
  go run main.go serve
  ```
  Default port: 8998. Access at `http://localhost:8998`.

- **Live Reload (Recommended)**:
  This project uses [air](https://github.com/cosmtrek/air) for live reloading.
  ```bash
  air
  ```

- **Database Setup**:
  The application automatically runs migrations on startup.
  Default DB location: `~/.golinks/golinks.db` (SQLite).
  Override with environment variable `GOLINKS_DB`.

## Development Workflow

### Backend
- Entry point: `cmd/serve.go` (and `main.go`).
- Routes defined in: `pkg/router/routes.go`.
- Controllers: `pkg/controllers/`.
- Models: `pkg/models/`.

### Frontend
- The frontend is contained entirely in `pkg/router/static/index.html`.
- It uses Vue 3 and Picnic CSS via CDN.
- No `npm` or `package.json` build step required. Modify `index.html` directly.

### Configuration
- Config via flags, environment variables (`GOLINKS_...`), or `viper`.
- See `.air.toml` for live reload configuration.

## Testing Instructions

- **Run Tests**:
  The project has unit tests for models (`pkg/models`) and integration tests for controllers (`pkg/controllers`).
  ```bash
  go test ./...
  ```
- **Linting**:
  Use `go vet` or `staticcheck`.
  ```bash
  go vet ./...
  ```

## Build and Deployment

- **Build Binary**:
  ```bash
  go build -o tmp/main .
  ```
- **Docker Build**:
  ```bash
  docker build . -t golinks:latest
  ```
- **Releases**:
  Project uses `goreleaser` (configs: `.goreleaser-darwin.yaml`, `.goreleaser-linux.yaml`).

## Code Style

- **Go**: Follow standard Go conventions (`gofmt`, `goimports`).
- **Commits**: Use descriptive commit messages.

## Common Tasks

### Adding a new Route
1. Define the handler in a controller (e.g., `pkg/controllers/controller.go`).
2. Register the route in `pkg/router/routes.go`.

### modifying the UI
1. Edit `pkg/router/static/index.html`.
2. App Logic is in the `<script>` section at the bottom of the file.
