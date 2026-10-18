# Security Policy

## Supported Versions

| Version | Supported |
| ------- | --------- |
| 2.x     | ✅        |

## Built-in Protections

- **Brute-force limiting** — 10 failed logins per address within 15 minutes, then `429` until the window passes. Password-protected share links are limited the same way. See [docs/deployment.md](docs/deployment.md#fail2ban).
- **Security headers** — every response carries `X-Content-Type-Options: nosniff`, `X-Frame-Options: SAMEORIGIN` (no framing by other sites) and `Referrer-Policy: same-origin`. Raw files are served with scripts disabled.
- **Scope confinement** — users can only reach files inside their scope; symbolic links pointing outside it are not followed unless `followExternalSymlinks` is enabled.
- **Hook commands** — file names are passed to hooks as environment variables, never pasted into a shell script. See [docs/command-execution.md](docs/command-execution.md).

> **Note:** The terminal is powerful. With a `shell` configured, the allowed-commands list is not a security boundary: see [docs/command-execution.md](docs/command-execution.md#interactive-shell).

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly.

Please include, where possible:

- The commit or version the issue was found at
- A plaintext proof of concept (no binaries)
- Steps to reproduce
- Recommended remediation, if any

We take security seriously and will respond to valid reports promptly.
