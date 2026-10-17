<p align="center">
  <img src="./branding/banner.png" width="550"/>
</p>

<p align="center">
  <strong>A modern, self-hosted file management interface for your server.</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/license-Non--Commercial-blue" alt="License" />
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-00ADD8?logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/vue-3-4FC08D?logo=vue.js" alt="Vue 3" />
  <img src="https://img.shields.io/badge/docker-ready-2496ED?logo=docker" alt="Docker" />
</p>

---

## Overview

File Browser provides a clean, intuitive web interface for managing files on your server. Upload, download, preview, edit, and organize your files from any browser — no FTP client needed. Think of it as your own personal cloud storage.

## Features

### File Management
- **Upload & Download** — Drag-and-drop file uploads, single and bulk downloads
- **File Preview** — Built-in preview for images, videos, audio, PDFs, and text files
- **Media Viewer** — Image toolbar (zoom, rotate, original quality), themed video player with FLV support, and a music player with album art, repeat, speed, and auto-advance
- **Download Progress** — Images, video, and audio show how much has loaded (percentage and size) while they download or buffer
- **File Editing** — Edit text files and code directly in the browser with syntax highlighting
- **Bulk Operations** — Select multiple files for move, copy, delete, or download as zip
- **Search** — Full-text search across your file system
- **Favorites** — Star files and folders for quick access from the sidebar

### User Interface
- **Modern Design** — Clean, responsive interface that works on desktop and mobile
- **Dark & Light Themes** — System-aware theme with manual override
- **File Info Modal** — Detailed file information with collapsible advanced section (permissions, checksums)
- **Keyboard Shortcuts** — Full keyboard navigation (F1 for help, F2 rename, Ctrl+S save/download)
- **Customizable Branding** — Set instance name, custom logos, and theme colors

### Security & Users
- **Multi-User Support** — Create and manage multiple user accounts with individual permissions
- **Role-Based Access** — Admin and regular user roles with configurable scopes
- **Flexible Passwords** — Strong password recommendations with option to use simpler passwords
- **Share Links** — Generate shareable links for files and directories with optional expiration

### Administration
- **User Management** — Create, edit, and delete users from the admin panel
- **Default Settings** — Configure default permissions, scope, and locale for new users
- **Storage Overview** — Monitor disk usage at a glance from the sidebar
- **Home Dashboard Auto-Refresh** — Rescan file statistics in the background every 5–60 minutes
- **Authentication Providers** — Support for proxy-based authentication, JSON auth, and more

### Deployment
- **Docker Ready** — Official Docker images with PUID/PGID support for proper file ownership
- **Single Binary** — Self-contained Go binary, easy to deploy anywhere
- **Cross-Platform** — Runs on Linux, macOS, Windows, ARM, and more
- **Reverse Proxy Friendly** — Works behind nginx, Caddy, Traefik, and others

### Developer Features
- **RESTful API** — Full API for programmatic file management
- **Command Hooks** — Execute custom scripts on file events (upload, delete, rename, etc.)
- **Shell Access** — Built-in terminal for server-side command execution (when enabled)
- **Customization** — Custom CSS, branding directory, and logo overrides

## Quick Start

### Docker (Recommended)

```bash
docker run -d \
  --name filebrowser \
  -p 8080:80 \
  -v /path/to/your/files:/srv \
  -v /path/to/config:/config \
  -e PUID=$(id -u) \
  -e PGID=$(id -g) \
  filebrowser
```

### Binary

```bash
# Download and install
curl -fsSL https://raw.githubusercontent.com/filebrowser/filebrowser/master/get.sh | bash

# Run with default settings
filebrowser -r /path/to/your/files
```

Default credentials: `admin` / `admin` — change the password immediately after first login.

## Configuration

File Browser can be configured via CLI flags, environment variables, or a config file.

```bash
# Basic usage
filebrowser --address 0.0.0.0 --port 8080 --root /srv

# With a database file
filebrowser --database /config/filebrowser.db

# With a config file
filebrowser --config /config/.filebrowser.json
```

See the [`docs`](docs) directory for full configuration documentation:

- [Installation](docs/installation.md) — All installation methods
- [Deployment](docs/deployment.md) — Reverse proxy setup and production tips
- [Authentication](docs/authentication.md) — Auth providers and configuration
- [Customization](docs/customization.md) — Branding, themes, and custom CSS
- [Command Execution](docs/command-execution.md) — Hooks and shell access
- [Troubleshooting](docs/troubleshooting.md) — Common issues and solutions
- [CLI Reference](docs/cli) — All CLI commands and flags

## Screenshots

<p align="center">
  <em>File listing with sidebar navigation</em>
</p>

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for build instructions and development setup.

## License
## License

The original File Browser project is licensed under [Apache License 2.0](LICENSE) © File Browser Contributors.

All modifications and additions in this fork are licensed under a [Non-Commercial Personal Use License](LICENSE-MODIFICATIONS) © Gab Huerte. Commercial use of the modifications is not permitted.
