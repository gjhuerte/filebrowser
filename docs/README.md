<p align="center">
  <img src="../branding/banner.png" width="550"/>
</p>

File Browser provides a file managing interface within a specified directory and it can be used to upload, delete, preview and edit your files. It is a **create-your-own-cloud**-kind of software where you can just install it on your server, direct it to a path and access your files through a nice web interface.

![Preview](static/example.gif)

## Contents

- [Installation](installation.md) — How to install and run File Browser
- [Deployment](deployment.md) — Reverse proxy setup and production configuration
- [Authentication](authentication.md) — Auth providers and SSO configuration
- [Customization](customization.md) — Branding, themes, and custom CSS
- [Command Execution](command-execution.md) — Hooks and shell configuration
- [Troubleshooting](troubleshooting.md) — Common issues and solutions
- [CLI Reference](cli/filebrowser.md) — All CLI commands and flags

## Recent Updates

- **Security Hardening** — Login and share-password rate limiting, safer hook commands, clickjacking and MIME-sniffing protection, and several bug fixes (see [SECURITY.md](../SECURITY.md))
- **Terminal** — Closing the terminal now stops the running command; long or stderr-heavy output no longer hangs it
- **Preview & Editor Top Bar** — Back button, file type icon, folder/size/position details, and a primary Download or Save button
- **Media Viewer** — Redesigned image, video, and audio previews that follow the light/dark theme
- **FLV Playback** — Flash video files play in the browser (remuxed on the fly; not supported on iPhone)
- **Download Progress** — Loading and buffering indicators with percentage and size
- **Home Auto-Refresh** — Background rescans of the Home statistics on a chosen interval
- **File Info Modal** — Redesigned with basic info and collapsible technical details (permissions, checksums)
- **Favorites** — Star files and folders for quick sidebar access
- **Flexible Passwords** — Allow simpler passwords with explicit confirmation
- **Docker PUID/PGID** — Proper file ownership support in Docker containers
- **Improved Help** — Comprehensive user guide with keyboard shortcuts, troubleshooting, and tips
