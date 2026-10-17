# Contributing

Thank you for your interest in contributing to File Browser!

## Building from Source

### Prerequisites

- [Go](https://golang.org/) (see `go.mod` for the required version)
- [Node.js](https://nodejs.org/) (latest LTS recommended)
- [pnpm](https://pnpm.io/) or npm

### Getting Started

```bash
git clone <your-repo-url>
cd filebrowser
```

### Frontend

The frontend lives under `frontend/` and uses Vue 3.

```bash
cd frontend
npm install
npm run build
```

For development with hot-reload:

```bash
npm run serve
```

### Backend

From the project root:

```bash
go build -o filebrowser
```

To run tests:

```bash
go test ./...
go vet ./...
```

### Building the Docker Image

```bash
docker build -t filebrowser .
```

## Code Style

- **Go**: Follow standard Go conventions. Run `go vet` and `go fmt` before committing.
- **Vue/TypeScript**: Follow the existing patterns in the codebase. Use Composition API with `<script setup>` for new components.
- **Commits**: Use conventional commit messages (`fix(scope): …`, `feat: …`, `chore: …`).

## Authentication Providers

To build a new authentication provider, implement the `Auther` interface in `auth/auth.go`. The method will be called on the login page after the user has submitted their login data.

Steps:
1. Add your provider to the `auth` directory
2. Add it to the configuration parser in `cmd/config.go`
3. Register it in `storage/bolt/auth.go`

## Pull Request Guidelines

- Keep changes focused and surgical
- Include tests for new functionality
- Make sure `go build ./...`, `go vet ./...`, and `go test ./...` pass
- Follow the existing code conventions
