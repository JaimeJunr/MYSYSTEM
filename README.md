# Homestead

> A CLI tool for managing, maintaining, and migrating Linux systems — built with Go and the [Charm ecosystem](https://github.com/charmbracelet).

Homestead brings together system cleanup scripts, dev tool installers, and shell configuration into a single interactive terminal interface. Whether you're setting up a new machine or keeping an existing one tidy, it's designed to make the tedious stuff fast.

---

## Features

- **System Cleanup** — Run maintenance scripts for Docker, npm, apt caches, and large file detection
- **System Monitoring** — Native in-TUI panels (Go, ~3s refresh): battery, RAM/swap, disk usage per mount, CPU load average, network RX/TX and per-interval throughput, thermal sensors (sysfs), and failing **systemd --user** units
- **Package Installers** — Install IDEs (Cursor, Claude Code CLI) and dev tools from a curated list (embedded catalog + optional remote manifest; see below)
- **Plugins e temas Zsh** — Wizard local para escolher plugins, ferramentas e gerar `.zshrc` (requer Oh My Zsh instalado)
- **Configurar Zsh** — Repositório de config (dotfiles): criar novo repo e enviar para a nuvem ou restaurar a partir de um repo existente; ideal para migração entre máquinas
- **Beautiful TUI** — Keyboard-driven interface built with Bubbletea and Lipgloss

### Preferences, theme, and installer catalog

On startup, `main` loads **`~/.config/homestead/preferences.yaml`** (see `internal/infrastructure/preferences`). You can change the same options from the TUI under **⚙️ Configurações**:

| Area | What it does |
|------|----------------|
| **Catalog URL** | HTTPS URL of the installer JSON manifest. Empty means the [default raw GitHub URL](https://raw.githubusercontent.com/JaimeJunr/Homestead/main/internal/infrastructure/catalog/installer-catalog.json) baked into `catalog.EffectiveCatalogURL`. If **`HOMESTEAD_CATALOG_URL`** is set in the environment, it overrides both the file and the default until you unset it. |
| **Theme** | Toggles **light** vs **dark** Lipgloss palettes (`internal/tui/theme/variants.go`); applied when the TUI starts and again after you save settings. |
| **Script root** | Directory that must contain a `scripts/` folder (bash assets). Empty uses the current working directory (see `preferences.ValidateScriptRoot`). |
| **Dotfiles repo** | Default path used by **Configurar Zsh** when no custom repo is set. |
| **Confirmations** | Optional prompts before running scripts or installing packages. |

The binary merges the **on-disk cache** of the catalog at startup (`main.go`), then the TUI **fetches** the effective URL in the background (`tui.Init` → `cmds.FetchCatalog`) and merges packages when the request succeeds.

---

## Quick Start

### Prerequisites

- Go 1.21+
- Linux (Ubuntu/Debian recommended)

### Install

```bash
git clone https://github.com/JaimeJunr/Homestead
cd Homestead
make install   # builds and installs to $GOPATH/bin
```

Or just run it directly:

```bash
make run
```

### Usage

```
↑/↓ or j/k    Navigate menus
Space          Toggle selection (in wizard)
Enter          Confirm
n or →         Next step (in wizard)
Esc            Go back
q or Ctrl+C    Quit
```

---

## Zsh: two flows

- **Plugins e temas Zsh** (menu, when Oh My Zsh is installed): local wizard — Plugins, Dev tools, Review & apply. Generates `.zshrc` and `~/.zsh/general/`; no git.
- **Configurar Zsh**: repo-based backup and migration. Choose “create new repo” (init, copy dotfiles, push to GitHub/etc.) or “already have repo” (clone/pull, backup, restore). Default dotfiles: `.zshrc`, `~/.zsh/`.
---

## Project Structure

```
homestead/
├── cmd/homestead/             # CLI entry point
├── internal/
│   ├── domain/               # Core business logic (entities, interfaces, types)
│   ├── app/services/         # Orchestration layer (ConfigService, PluginService, WizardService)
│   ├── infrastructure/       # File storage, plugin installer, template engine
│   │   ├── config/           # YAML-based config persistence
│   │   ├── plugins/          # Zsh plugin install/update/remove
│   │   ├── templates/        # Go template renderer
│   │   └── repository/       # In-memory package catalog (29 packages)
│   └── tui/                  # Bubbletea models and views
├── scripts/
│   ├── cleanup/              # System maintenance scripts
│   └── monitoring/           # Optional bash helpers (native monitors are Go under internal/monitoring/)
└── docs/                     # INDEX.md, product context, architecture, ADRs
```

The codebase follows [Clean Architecture](docs/architecture/ARCHITECTURE.md) — domain logic has zero external dependencies, infrastructure is pluggable, and the TUI layer only talks to services.

---

## Development

```bash
make test              # Run all tests
make test-coverage     # Tests with coverage report
make test-verbose      # Verbose output
make benchmark         # Performance benchmarks
make build             # Build binary
make clean             # Remove build artifacts
```

The project uses strict TDD. Each feature is developed Red → Green → Refactor, with tests written before implementation.

**Current test count: 97+ tests across 9 packages.**

### Adding a Script

1. Drop the script in `scripts/<category>/`
2. Register it in `internal/scripts/script.go` inside `GetAllScripts()`

```go
{
    ID:           "my-script",
    Name:         "My Script",
    Description:  "What it does",
    Path:         "scripts/cleanup/my_script.sh",
    Category:     string(CategoryCleanup),
    RequiresSudo: true,
}
```

### Adding a Package to the Installer

Add an entry to `initializeDefaultPackages()` in `internal/infrastructure/repository/package_repository.go`:

```go
{
    ID:          "my-tool",
    Name:        "My Tool",
    Description: "What it does",
    Version:     "latest",
    Category:    types.PackageCategoryTool,
    InstallCmd:  "curl -fsSL https://example.com/install.sh | bash",
    CheckCmd:    "which my-tool",
},
```

---

## Architecture

The project is structured in 4 layers, strictly following Clean Architecture:

| Layer | Package | Responsibility |
|-------|---------|----------------|
| Domain | `internal/domain` | Entities, interfaces, types — no external deps |
| Application | `internal/app/services` | Business orchestration |
| Infrastructure | `internal/infrastructure` | File I/O, git, system commands |
| Presentation | `internal/tui` | Bubbletea models and views |

Dependencies only flow inward. The TUI never touches the filesystem directly.

Documentation map: [docs/INDEX.md](docs/INDEX.md). For diagrams, ADRs, and patterns, see [docs/architecture/](docs/architecture/).

---

## Roadmap

- [x] System cleanup and monitoring scripts
- [x] Clean Architecture with 4 layers
- [x] Full TDD suite (97+ tests)
- [x] Package installer framework
- [x] IDEs: Cursor AI, Claude Code CLI, Antigravity
- [x] Zsh configuration wizard
- [x] 15 Zsh plugins, 8 dev tools
- [ ] Real `.zshrc` template generation (Go templates)
- [ ] End-to-end integration tests
- [ ] Migration: export/import dotfiles and configs
- [ ] Remote execution via SSH

---

## Contributing

Pull requests are welcome.

Before submitting:

```bash
make test           # all tests must pass
make test-coverage  # check coverage
```

Follow the existing patterns in `docs/architecture/PATTERNS_GUIDE.md`. If you're adding a new layer or changing how things are wired, consider adding an ADR under `docs/architecture/adrs/` (see `docs/architecture/adrs/README.md`).

---

## Tech Stack

- [Go](https://go.dev) 1.21+
- [Bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — UI components
- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) — config serialization

---

## License

MIT
