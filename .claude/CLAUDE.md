# CLAUDE.md

Guidance for Claude when working in this repository (File Browser).

## Repo Orientation

- Go backend + Vue 3 frontend under `frontend/`.
- Vue components use both Options API and Composition API (`<script setup>`).
- State management via Pinia stores (`frontend/src/stores/`).
- Internationalization via vue-i18n (`frontend/src/i18n/en.json`).
- Design system tokens: `--primary`, `--textPrimary`, `--textSecondary`, `--surfaceSecondary`, `--borderPrimary`.
- Bootstrap 5 utilities preferred, minimal custom CSS class files.
- Conventional-commit messages (`fix(scope): …`, `feat: …`, `chore: …`).

## Build & Test

```bash
# Backend
go build ./...
go vet ./...
go test ./...

# Frontend
cd frontend
npm install
npm run build    # production build
npm run serve    # dev server with hot-reload
```

## Code Conventions

- **Go**: Standard Go conventions. Use `go vet` and `go fmt`.
- **Vue/TypeScript**: Use Composition API with `<script setup>` for new components.
- **Styling**: Inline styles or scoped `<style>` blocks using design tokens. Light-themed, minimal, no heavy custom CSS frameworks.
- **i18n**: All user-facing strings go in `frontend/src/i18n/en.json`. Use `$t("key")` in templates, `t("key")` in setup scripts.

## License

This project uses a **dual-license structure**:
- **LICENSE** — Apache License 2.0 for the original File Browser code (must be preserved)
- **LICENSE-MODIFICATIONS** — Non-Commercial Personal Use License for Gab's modifications

When modifying license-related content:
- Never remove or alter the Apache 2.0 LICENSE file — it covers the original codebase
- New copyright notices for modifications should reference "Gab Huerte"
- The README must reference both licenses
- Any new code added to this fork falls under LICENSE-MODIFICATIONS

---

# Documentation Maintenance Rules

These rules ensure documentation stays in sync with code changes. Follow them whenever making changes to the codebase.

## When to Update Documentation

### Always update docs when:

1. **Adding a new feature or UI component**
   - Add the feature to `README.md` under the appropriate Features section
   - If it has user-facing interactions, add it to `Help.vue` (the help modal)
   - Update `docs/README.md` Recent Updates section
   - If it involves configuration, update the relevant doc in `docs/`

2. **Changing keyboard shortcuts or hotkeys**
   - Update `frontend/src/components/prompts/Help.vue` (Shortcuts tab)
   - Update `frontend/src/i18n/en.json` under the `help` key

3. **Changing authentication or security behavior**
   - Update `docs/authentication.md`
   - Update `SECURITY.md` if security-relevant
   - Update Help.vue troubleshooting section if user-facing

4. **Changing CLI commands or flags**
   - Update the relevant file in `docs/cli/`
   - If it's a new command, create a new CLI doc following the existing pattern

5. **Changing Docker configuration or deployment**
   - Update `docs/deployment.md` and/or `docs/installation.md`
   - Update README.md Quick Start section if the basic usage changes
   - Update Help.vue troubleshooting Docker section

6. **Changing branding or theming**
   - Update `docs/customization.md`
   - Update any screenshots in `docs/static/` if the UI has changed

7. **Changing command execution / hooks**
   - Update `docs/command-execution.md`

8. **Adding or modifying i18n strings**
   - Always add new keys to `frontend/src/i18n/en.json`
   - Follow the existing key hierarchy pattern

### Files to check for each change type:

| Change Type | Files to Update |
|---|---|
| New feature | `README.md`, `docs/README.md`, `Help.vue` |
| UI change | `Help.vue`, screenshots in `docs/static/` |
| Config change | Relevant `docs/*.md`, `docs/cli/*.md` |
| API change | `docs/README.md` |
| Security change | `SECURITY.md`, `docs/authentication.md` |
| Docker change | `docs/installation.md`, `docs/deployment.md`, `README.md` |
| Keyboard shortcut | `Help.vue`, `en.json` |
| Password/auth flow | `docs/authentication.md`, `Help.vue` |

## Documentation Style

- Use clear, concise language
- Include code examples for configuration options
- Use `> **Note:**` for important callouts (not `> [!WARNING]` or deprecated notices)
- Keep troubleshooting sections practical — describe the symptom, then the fix
- Reference design tokens by name when documenting UI components

## Project Structure Reference

```
├── .claude/CLAUDE.md          # This file — AI development rules
├── LICENSE                     # Apache License 2.0 (ORIGINAL — do not alter)
├── README.md                   # Project overview and features
├── SECURITY.md                 # Security policy
├── CONTRIBUTING.md             # Build instructions
├── branding/                   # Logo and banner assets
├── docs/                       # Full documentation
│   ├── README.md              # Docs index with recent updates
│   ├── installation.md        # Installation methods
│   ├── deployment.md          # Reverse proxy and production setup
│   ├── authentication.md      # Auth providers
│   ├── customization.md       # Branding and themes
│   ├── command-execution.md   # Hooks and shell
│   ├── troubleshooting.md     # Common issues
│   ├── cli/                   # CLI command reference
│   └── static/                # Screenshots and media
├── frontend/                   # Vue 3 frontend
│   └── src/
│       ├── components/        # Vue components
│       │   └── prompts/       # Modal components (Help.vue, Info.vue)
│       ├── views/             # Page-level components
│       ├── stores/            # Pinia stores
│       ├── i18n/              # Translations (en.json)
│       └── api/               # API client
├── http/                       # Go HTTP handlers
├── users/                      # User management (password validation)
├── auth/                       # Authentication providers
├── storage/                    # Database layer (BoltDB)
└── cmd/                        # CLI entry point
```
