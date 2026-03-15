# 📐 Diagramas de Arquitetura - Homestead

Diagramas visuais da arquitetura do Homestead.

## 🏗️ Diagrama de Camadas

```
┌─────────────────────────────────────────────────────────────┐
│                     PRESENTATION LAYER                       │
│                      (internal/tui/)                         │
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  Model   │  │  Views   │  │Components│  │ Observers│   │
│  │  (State) │  │  (Render)│  │(Reusable)│  │(Progress)│   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│       │             │              │             │          │
└───────┼─────────────┼──────────────┼─────────────┼──────────┘
        │             │              │             │
        ▼             ▼              ▼             ▼
┌─────────────────────────────────────────────────────────────┐
│                    APPLICATION LAYER                         │
│                      (internal/app/)                         │
│                                                              │
│  ┌──────────────┐          ┌──────────────┐                │
│  │  Use Cases   │          │   Services   │                │
│  ├──────────────┤          ├──────────────┤                │
│  │ • Execute    │──────────│ • Script     │                │
│  │   Script     │          │   Service    │                │
│  │ • Install    │──────────│ • Installer  │                │
│  │   Package    │          │   Service    │                │
│  │ • Export     │          │ • Config     │                │
│  │   System     │          │   Service    │                │
│  └──────────────┘          └──────────────┘                │
│         │                         │                          │
└─────────┼─────────────────────────┼──────────────────────────┘
          │                         │
          ▼                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      DOMAIN LAYER                            │
│                    (internal/domain/)                        │
│                                                              │
│  ┌─────────────┐         ┌──────────────────┐              │
│  │  Entities   │         │   Interfaces     │              │
│  ├─────────────┤         ├──────────────────┤              │
│  │ • Script    │         │ • Repository     │              │
│  │ • Package   │         │ • Executor       │              │
│  │ • Installer │         │ • Installer      │              │
│  │ • Config    │         │ • Observer       │              │
│  └─────────────┘         └──────────────────┘              │
│         ▲                         ▲                          │
└─────────┼─────────────────────────┼──────────────────────────┘
          │                         │
          │ implements              │ implements
          │                         │
┌─────────┴─────────────────────────┴──────────────────────────┐
│                  INFRASTRUCTURE LAYER                         │
│                 (internal/infrastructure/)                    │
│                                                              │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │ Executors  │  │Repositories│  │  Adapters  │            │
│  ├────────────┤  ├────────────┤  ├────────────┤            │
│  │• Bash      │  │• InMemory  │  │• Apt       │            │
│  │• Docker    │  │• File      │  │• Snap      │            │
│  └────────────┘  └────────────┘  │• Flatpak   │            │
│                                   └────────────┘            │
│                                                              │
│  ┌────────────────────────────────────────────────┐         │
│  │          External Systems                      │         │
│  │  • Bash Scripts  • apt  • snap  • docker      │         │
│  └────────────────────────────────────────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

## 🔄 Fluxo de Dados - Executar Script

```
     USER
       │
       │ 1. Seleciona "Limpeza" → "Limpeza Geral"
       ▼
  ┌─────────┐
  │   TUI   │
  │  Model  │
  └────┬────┘
       │ 2. handleEnter()
       │    scriptID = "cleanup-general"
       ▼
  ┌─────────────────┐
  │ ScriptService   │
  │ (Application)   │
  └────┬────────────┘
       │ 3. Execute(scriptID)
       │
       ├──────────────────┐
       │                  │
       ▼                  ▼
  ┌──────────┐      ┌──────────┐
  │  Script  │      │  Bash    │
  │Repository│      │ Executor │
  │  (Infra) │      │  (Infra) │
  └────┬─────┘      └────┬─────┘
       │ 4. FindByID()    │
       │    returns       │
       │    Script        │
       ▼                  │
       Script Entity      │
       │                  │
       └──────────────────┘
                │ 5. Execute(script)
                ▼
           ┌─────────┐
           │  bash   │
           │ process │
           └────┬────┘
                │ 6. stdout/stderr
                ▼
              Result
                │
                └─────► TUI (mostra sucesso/erro)
```

## 🏭 Fluxo de Dados - Instalar Package

```
USER seleciona "Instalar Git"
  │
  ▼
┌──────────────────┐
│  TUI Wizard      │
│ 1. Coleta config │
└────────┬─────────┘
         │ config: {name: "git", type: "apt"}
         ▼
┌──────────────────────┐
│ InstallPackageUseCase│
│   (Application)      │
└────────┬─────────────┘
         │ 2. Execute(config)
         │
         ├───────────────┬────────────────┐
         │               │                │
         ▼               ▼                ▼
    ┌─────────┐    ┌─────────┐    ┌──────────┐
    │ Factory │    │Observer │    │Repository│
    │         │    │ (TUI)   │    │          │
    └────┬────┘    └────┬────┘    └──────────┘
         │              │
         │ 3. Create("apt")
         ▼              │
    ┌─────────┐         │
    │   Apt   │         │
    │Installer│         │
    └────┬────┘         │
         │              │
         │ 4. Install("git")
         ├──────────────┤
         │ OnStart()    ├──► TUI: "Instalando..."
         │              │
         │ OnProgress() ├──► TUI: Progress bar 50%
         │              │
         │ OnComplete() ├──► TUI: "✓ Instalado!"
         ▼              │
    ┌─────────┐         │
    │apt install│        │
    │   git    │        │
    └──────────┘        │
         │              │
         └──────────────┘
```

## 🎯 Diagrama de Padrões

```
┌─────────────────────────────────────────────────┐
│              DESIGN PATTERNS                     │
└─────────────────────────────────────────────────┘

Repository Pattern
═════════════════
┌──────────────┐
│ScriptService │
└──────┬───────┘
       │ depends on
       ▼
┌──────────────────┐  interface
│ScriptRepository  │◄─────────────┐
└──────────────────┘              │
       ▲                           │ implements
       │                           │
       └───────────────────────────┘
              ┌──────────────────────┐
              │InMemoryScriptRepo    │
              └──────────────────────┘

Factory Pattern
═══════════════
┌─────────────┐
│   Service   │
└──────┬──────┘
       │ uses
       ▼
┌──────────────────┐
│InstallerFactory  │
└──────┬───────────┘
       │ creates
       ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│AptInstaller  │   │SnapInstaller │   │GitInstaller  │
└──────────────┘   └──────────────┘   └──────────────┘

Strategy Pattern
════════════════
┌─────────────────┐
│InstallerService │
└────────┬────────┘
         │ uses
         ▼
  ┌──────────────┐
  │  Strategy    │  interface
  └──────────────┘
         ▲
         │ implements
    ┌────┴────┬────────┬────────┐
    │         │        │        │
┌───────┐ ┌───────┐ ┌──────┐ ┌──────┐
│  Apt  │ │ Snap  │ │Manual│ │Source│
│Strategy│ │Strategy│ │Strategy│ │Strategy│
└───────┘ └───────┘ └──────┘ └──────┘

Observer Pattern
════════════════
┌──────────────────┐
│ObservableInstaller│
└────────┬─────────┘
         │ notifies
         ▼
  ┌──────────────┐
  │   Observer   │  interface
  └──────────────┘
         ▲
         │ implements
    ┌────┴─────┬──────────┐
    │          │          │
┌────────┐ ┌────────┐ ┌─────────┐
│  TUI   │ │  Log   │ │ Metrics │
│Observer│ │Observer│ │ Observer│
└────────┘ └────────┘ └─────────┘

Command Pattern
═══════════════
┌──────────────┐
│   Executor   │
└──────┬───────┘
       │ executes
       ▼
┌──────────────┐
│   Command    │  interface
└──────────────┘
       ▲
       │ implements
  ┌────┴────┬────────┬────────┐
  │         │        │        │
┌─────┐ ┌──────┐ ┌──────┐ ┌──────┐
│Install│ │Uninstall│ │Update│ │Export│
│Command│ │Command │ │Command│ │Command│
└─────┘ └──────┘ └──────┘ └──────┘

Builder Pattern
═══════════════
┌──────────────────────┐
│InstallerConfigBuilder│
└──────────┬───────────┘
           │ builds
           ▼
    ┌──────────────┐
    │InstallerConfig│
    └──────────────┘

Adapter Pattern
═══════════════
┌─────────────┐
│   Service   │
└──────┬──────┘
       │ uses
       ▼
┌──────────────────┐
│ PackageManager   │  interface
└──────────────────┘
       ▲
       │ implements
  ┌────┴────┬────────┐
  │         │        │
┌─────┐ ┌──────┐ ┌────────┐
│ Apt │ │ Snap │ │Flatpak │
│Adapter│ │Adapter│ │ Adapter│
└─────┘ └──────┘ └────────┘
  │       │        │
  │       │        │ adapts
  ▼       ▼        ▼
┌─────┐ ┌──────┐ ┌────────┐
│ apt │ │ snap │ │flatpak │
│ CLI │ │  CLI │ │  CLI   │
└─────┘ └──────┘ └────────┘
```

## 📦 Diagrama de Dependências

```
main.go
  │
  ├──► TUI (Presentation)
  │     │
  │     └──► Application Services
  │           │
  │           └──► Domain Interfaces
  │
  ├──► Application Services
  │     │
  │     ├──► Domain Interfaces
  │     └──► Use Cases
  │           │
  │           └──► Domain Interfaces
  │
  └──► Infrastructure
        │
        ├──► Domain Interfaces (implements)
        └──► External Systems (apt, snap, bash)

Regras:
• Presentation → Application → Domain ← Infrastructure
• Domain não depende de nada (core isolado)
• Infrastructure implementa interfaces do Domain
• Application orquestra Domain + Infrastructure
```

## 🔌 Diagrama de Módulos

```
github.com/JaimeJunr/Homestead
│
├── cmd/
│   └── homestead/
│       └── main.go  ────┐
│                        │ imports
├── internal/            │
│   │                    │
│   ├── domain/          │
│   │   ├── entities/ ◄──┼──────────┐
│   │   ├── interfaces/◄─┼────┐     │
│   │   └── types/    ◄──┼──┐ │     │
│   │                     │  │ │     │
│   ├── app/              │  │ │     │
│   │   ├── usecases/  ──┼──┘ │     │
│   │   └── services/  ──┼────┘     │
│   │                     │          │
│   ├── infrastructure/   │          │
│   │   ├── executor/  ──┼──────────┘
│   │   ├── repository/──┼──────────┐
│   │   └── adapter/  ───┼────┐     │
│   │                     │    │     │
│   └── tui/              │    │     │
│       ├── model.go ─────┤    │     │
│       ├── views/    ────┤    │     │
│       └── components/───┤    │     │
│                         │    │     │
└── scripts/              │    │     │
    ├── cleanup/          │    │     │
    └── monitoring/       │    │     │
                          │    │     │
                          ▼    ▼     ▼
                    Dependency Flow
```

## 🎨 Diagrama de Instalação de Package

```
┌─────────────────────────────────────────────────┐
│            INSTALL PACKAGE FLOW                  │
└─────────────────────────────────────────────────┘

1. User Input
   ┌──────────┐
   │   User   │
   └────┬─────┘
        │ "Instalar Docker"
        ▼
   ┌──────────┐
   │   TUI    │
   │  Wizard  │
   └────┬─────┘
        │ collect config
        ▼

2. Application Layer
   ┌─────────────────────┐
   │InstallPackageUseCase│
   └──────────┬──────────┘
              │
              ├─── Validate input
              ├─── Check if installed
              ├─── Resolve dependencies
              └─── Execute installation
                   │
                   ▼

3. Domain Layer
   ┌──────────────────┐
   │  Package Entity  │
   │ ┌──────────────┐ │
   │ │ name: docker │ │
   │ │ type: apt    │ │
   │ │ version: *   │ │
   │ └──────────────┘ │
   └──────────────────┘

4. Infrastructure Layer
   ┌──────────────┐
   │   Factory    │
   └──────┬───────┘
          │ Create("apt")
          ▼
   ┌──────────────┐     ┌──────────────┐
   │AptInstaller  │────►│   Strategy   │
   └──────┬───────┘     └──────────────┘
          │
          │ Install()
          ▼
   ┌──────────────┐
   │   Observer   │
   │  (notifies)  │
   └──────┬───────┘
          │
          ├─── TUIObserver    → Update UI
          ├─── LogObserver    → Write log
          └─── MetricObserver → Send metrics
          │
          ▼
   ┌──────────────┐
   │  Apt Adapter │
   └──────┬───────┘
          │ Run command
          ▼
   ┌──────────────┐
   │apt install   │
   │   docker     │
   └──────────────┘

5. Result
   ┌──────────────┐
   │   Success    │
   └──────┬───────┘
          │
          └─── Back to TUI
               │
               ▼
          ┌──────────┐
          │Show: ✓   │
          │Installed!│
          └──────────┘
```

## 📊 Resumo Visual

```
┌────────────────────────────────────────┐
│      ARCHITECTURE SUMMARY              │
├────────────────────────────────────────┤
│                                        │
│  Layers:        4 (Presentation,      │
│                    Application,        │
│                    Domain,             │
│                    Infrastructure)     │
│                                        │
│  Patterns:      7 (Repository,        │
│                    Factory,            │
│                    Strategy,           │
│                    Observer,           │
│                    Command,            │
│                    Builder,            │
│                    Adapter)            │
│                                        │
│  Dependencies:  Always inward          │
│                 (Dependency Inversion) │
│                                        │
│  Testing:       Mock interfaces        │
│                 Unit + Integration     │
│                                        │
└────────────────────────────────────────┘
```

---

**Última atualização**: 2026-03-14

Para mais detalhes, veja:
- [ARCHITECTURE.md](../ARCHITECTURE.md)
- [PATTERNS_GUIDE.md](PATTERNS_GUIDE.md)
- [ARCHITECTURE_DECISION_RECORD.md](ARCHITECTURE_DECISION_RECORD.md)
