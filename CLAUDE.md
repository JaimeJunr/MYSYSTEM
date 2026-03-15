# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Homestead** is a CLI tool built with Go and the Charm Ecosystem for system maintenance, configuration, and migration. It helps users manage their Linux systems, run maintenance scripts, install tools, and migrate between systems with ease.

## Tech Stack

- **Language**: Go
- **UI Framework**: [Charm Ecosystem](https://github.com/charmbracelet)
  - `bubbletea` - TUI framework (The Elm Architecture)
  - `lipgloss` - Style definitions and layout
  - `bubbles` - Common TUI components
  - `huh` - Forms and interactive prompts
  - `log` - Charm logging

## Project Goals

1. **Script Management** - Organize and run maintenance scripts (existing bash scripts will be integrated)
2. **System Installation** - Install and configure common development tools and applications
3. **Migration Assistant** - Help migrate configurations and tools between systems
4. **Interactive TUI** - Beautiful, user-friendly terminal interface
5. **Extensible** - Easy to add new scripts, installers, and features

## Current Assets (To Be Integrated)

The repository contains existing bash scripts that should be integrated:
- `limpar_ssd.sh` - System cleanup orchestrator
- `limpar_geral.sh` - Cache and system cleanup (Docker, Poetry, npm, apt, etc.)
- `limpar_grandes.sh` - Large file/directory scanner
- `teste_bateria.sh` - Battery monitoring
- `memoria.sh` - Memory usage display

These scripts should be callable from the Go CLI tool.

## Architecture (Planned)

```
homestead/
├── cmd/
│   └── homestead/          # Main CLI entry point
│       └── main.go
├── internal/
│   ├── tui/               # Bubbletea UI components
│   │   ├── model.go       # Root model
│   │   ├── menu.go        # Main menu
│   │   └── views/         # Different view states
│   ├── scripts/           # Script execution and management
│   │   ├── runner.go      # Execute bash scripts
│   │   └── registry.go    # Script catalog
│   ├── installers/        # System/tool installers
│   │   ├── installer.go   # Interface
│   │   └── packages/      # Individual package installers
│   ├── migration/         # Migration tools
│   │   ├── export.go      # Export system config
│   │   └── import.go      # Import/restore config
│   └── config/            # Configuration management
│       └── config.go
├── scripts/               # Bash scripts (existing ones)
│   ├── cleanup/
│   │   ├── limpar_geral.sh
│   │   ├── limpar_grandes.sh
│   │   └── limpar_ssd.sh
│   ├── monitoring/
│   │   ├── teste_bateria.sh
│   │   └── memoria.sh
│   └── install/           # Future installation scripts
├── configs/               # Configuration files and templates
├── go.mod
├── go.sum
└── README.md
```

## Development Commands

### Setup
```bash
# Initialize Go module (if not done)
go mod init github.com/JaimeJunr/Homestead

# Install Charm dependencies
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/huh
go get github.com/charmbracelet/log
```

### Build and Run
```bash
# Build the CLI
go build -o homestead ./cmd/homestead

# Run directly
go run ./cmd/homestead

# Install to $GOPATH/bin
go install ./cmd/homestead
```

### Testing
```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/scripts
```

## Key Design Patterns

### Bubbletea (The Elm Architecture)

All UI components follow the Elm Architecture:
```go
type Model struct {
    // state
}

func (m Model) Init() tea.Cmd {
    // initialization
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // handle messages and update state
}

func (m Model) View() string {
    // render UI
}
```

### Script Execution

Scripts should be executed with proper user context:
```go
// Maintain the existing user detection pattern
cmd := exec.Command("bash", scriptPath)
cmd.Env = append(os.Environ(),
    fmt.Sprintf("REAL_USER=%s", realUser),
    fmt.Sprintf("REAL_HOME=%s", realHome),
)
```

### Navigation Pattern

Use a state machine for navigation:
- Menu state → shows main menu
- Script execution state → runs script with live output
- Installer state → interactive installation wizard
- Migration state → export/import flows

## Feature Roadmap

### Phase 1: Core CLI
- [x] Existing bash scripts (cleanup, battery, memory)
- [ ] Go project setup with Charm
- [ ] Main menu TUI
- [ ] Script runner with live output
- [ ] Script registry/catalog

### Phase 2: Installers
- [ ] Package installer framework
- [ ] Common dev tools (git, docker, node, python, rust, etc.)
- [ ] IDE installers (VSCode, Cursor)
- [ ] Configuration templates

### Phase 3: Migration
- [ ] Export system state (installed packages, configs)
- [ ] Import/restore on new system
- [ ] Dotfile management
- [ ] Backup/restore scripts

### Phase 4: Advanced
- [ ] Plugin system
- [ ] Remote execution (SSH)
- [ ] Scheduled maintenance
- [ ] System health dashboard

## Coding Guidelines

- Use Go standard project layout
- Keep UI logic separate from business logic
- Make scripts idempotent when possible
- Handle errors gracefully with user-friendly messages
- Use Charm's logging for debug output
- Style with lipgloss for consistent UI
- Prefer composition over inheritance

## Charm Component Usage

- **bubbletea**: Main TUI framework, state management
- **lipgloss**: All styling (colors, borders, padding, alignment)
- **bubbles**: Use existing components (list, spinner, progress, viewport, etc.)
- **huh**: Forms for user input (installation configs, migration options)
- **log**: Development logging and debug output

## Integration with Existing Scripts

The existing bash scripts use an interactive pattern with `confirm_action()`. When integrating:
1. Option 1: Run scripts as-is with PTY for interactivity
2. Option 2: Parse scripts and recreate logic in Go with huh forms
3. Option 3: Hybrid - use Go for UI, bash for execution

Prefer Option 1 for quick integration, then refactor to Option 2 for better UX.
