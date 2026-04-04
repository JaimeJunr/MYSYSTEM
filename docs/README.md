# Documentation

This folder contains architecture documentation, design decisions, and development history for Homestead.

## Structure

```
docs/
├── README.md                              # You are here
├── TESTING.md                             # Testing guide and conventions
├── architecture/
│   ├── ARCHITECTURE.md                    # Layer overview and conventions
│   ├── PATTERNS_GUIDE.md                  # Design patterns with code examples
│   ├── adrs/                              # Architecture Decision Records (per-file)
│   └── DIAGRAMS.md                        # Data flow and module diagrams
└── development/
    ├── REFACTORING_SUMMARY.md             # History of architectural changes
    └── INSTALLER_IMPLEMENTATION.md        # Notes on the installer system
```

## Where to start

**New to the project?**
→ Read [../README.md](../README.md), then [../GETTING_STARTED.md](../GETTING_STARTED.md)

**Want to understand the architecture?**
→ [architecture/ARCHITECTURE.md](architecture/ARCHITECTURE.md) — the 4-layer overview
→ [architecture/DIAGRAMS.md](architecture/DIAGRAMS.md) — visual data flows

**Adding a feature?**
→ [architecture/ARCHITECTURE.md](architecture/ARCHITECTURE.md) — which layer it belongs to
→ [architecture/PATTERNS_GUIDE.md](architecture/PATTERNS_GUIDE.md) — code patterns to follow
→ [TESTING.md](TESTING.md) — how to write tests

**Wondering why something was built a certain way?**
→ [architecture/adrs/README.md](architecture/adrs/README.md)

---

## Architecture overview

Homestead follows Clean Architecture with 4 layers:

```
┌─────────────────────────────────────┐
│  Presentation  (internal/tui)       │  Bubbletea models, views
├─────────────────────────────────────┤
│  Application   (internal/app)       │  ConfigService, PluginService, WizardService
├─────────────────────────────────────┤
│  Infrastructure (internal/infra)    │  Files, git, system commands
├─────────────────────────────────────┤
│  Domain        (internal/domain)    │  Entities, interfaces, types — no deps
└─────────────────────────────────────┘
```

Dependencies only flow inward. Nothing in `domain` imports from other layers. `tui` never touches the filesystem directly — it goes through services.

---

## Project metrics

| | |
|---|---|
| Test packages | 9 |
| Tests | 97+ |
| Architecture layers | 4 |
| Packages in installer catalog | 29 |
| Zsh plugins available | 15 |
| ADRs documented | 12 |

---

**Last updated:** 2026-03-15
