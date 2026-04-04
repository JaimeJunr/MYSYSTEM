# Agent Instructions: Homestead

## Overview

Homestead is a Go CLI with a Bubble Tea TUI for Linux maintenance: run bundled bash scripts (cleanup, utilities), view native monitors (battery, memory), and manage installers/config (packages, Zsh wizard, file config). Entry point wires dependencies manually and starts the full-screen TUI.

## Tech Stack

| Layer    | Technology                   | Notes                                                             |
| -------- | ---------------------------- | ----------------------------------------------------------------- |
| Language | Go                           | `go 1.21` (see `go.mod`, CI release workflow)                     |
| TUI      | bubbletea, lipgloss, bubbles | Primary UI stack                                                  |
| Config   | `gopkg.in/yaml.v3`           | File-based config                                                 |
| Scripts  | Bash under `scripts/`        | Invoked via `BashExecutor`; `REAL_USER` / `REAL_HOME` env pattern |

`huh` and Charm `log` are mentioned in older docs but are **not** current `require` entries in `go.mod`—verify before importing.

## Project structure

| Path                       | Role                                                                                                      |
| -------------------------- | --------------------------------------------------------------------------------------------------------- |
| `cmd/homestead/main.go`    | DI wiring: repos, executor, installer, config, `tui.NewModel`, `tea.NewProgram`                           |
| `internal/tui/`            | Bubble Tea models, menus, script output, Zsh wizard                                                       |
| `internal/app/services/`   | Application services (`ScriptService`, installers, config, repo, wizard)                                  |
| `internal/domain/`         | Entities (`Script`, packages, shell config), `types` (e.g. `Category`), interfaces                        |
| `internal/infrastructure/` | `repository/` (in-memory catalogs), `executor/` (bash), `installer/`, `config/`, `templates/`, `plugins/` |
| `internal/scripts/`        | Legacy script helpers (parallel to domain-driven path; prefer domain + infra for new work)                |
| `scripts/`                 | Bash assets: `cleanup/`, `monitoring/`, `install/`, `utilities/`, `lib/`                                  |
| `configs/`                 | Templates / static config                                                                                 |
| `integration_test.go`      | Top-level integration test; mirrors `main` wiring                                                         |
| `Makefile`                 | `build`, `run`, `test`, `test-short`, `test-integration`, `test-coverage`                                 |

## Data flow (run a script)

1. TUI selects a script → `ScriptService` loads `entities.Script` from `ScriptRepository`.
2. `BashExecutor` validates paths / permissions; executes bash (or builds interactive `exec.Cmd` for sudo/TTY).
3. Native monitors use `NativeMonitor` on `entities.Script` (`battery` / `memory`) with TUI-side handling (no bash path).

Categories include `cleanup`, `monitoring`, `install`, `utilities` (`internal/domain/types/category.go`).

## Code style & errors

- Go std layout; layered **domain → services → infrastructure → tui**.
- Errors: `fmt.Errorf` with `%w`; domain sentinel-style errors in `internal/domain/types` where used.
- File naming: `snake_case` Go files (e.g. `script_service.go`); tests `*_test.go` alongside or in same package.

## Testing

- All packages: `go test ./...` or `make test`
- Skip slow/integration: `go test -short ./...` or `make test-short`
- Integration: `make test-integration` (`go test -v -run Integration ./...`)
- Coverage: `make test-coverage`

## Build & run

- `make run` / `go run ./cmd/homestead`
- `make build` → `./homestead`
- `make install` → `$GOPATH/bin/homestead`

## CI / release

- `.github/workflows/release.yml`: on GitHub **release published**, builds `linux/amd64` and `linux/arm64` and uploads artifacts.

## Conventions (git)

- Recent commits use prefixes such as `feat(scope):`, `Refactor`, `Update` (mixed English/Portuguese).

## Do not

- Assume `go run` cwd: bash paths in the catalog are relative to repo root; running the binary from elsewhere may break script resolution unless code resolves embed/root explicitly.
- List `huh` / `log` as installed without checking `go.mod`.
- Expand scope into `linuxtoys/` or `.agents/` unless the task explicitly includes them.
