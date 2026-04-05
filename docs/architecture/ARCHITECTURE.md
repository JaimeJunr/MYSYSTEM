# 🏗️ Arquitetura de Software - Homestead

Documento de arquitetura e padrões de projeto para o Homestead.

**Ordem de leitura:** [VERTICAL_PATTERNS.md](VERTICAL_PATTERNS.md) define a pilha de padrões **de cima para baixo** (TUI → aplicação → domínio → infra). Este ficheiro aprofunda princípios, camadas e padrões com mais detalhe.

## 📋 Índice

- [Padrões (visão vertical)](VERTICAL_PATTERNS.md) — mapa top-down do sistema
- [Visão Geral](#visão-geral)
- [Princípios Arquiteturais](#princípios-arquiteturais)
- [Camadas da Aplicação](#camadas-da-aplicação)
- [Padrões de Projeto](#padrões-de-projeto)
- [Estrutura de Diretórios](#estrutura-de-diretórios)
- [Convenções de Código](#convenções-de-código)
- [Casos de Uso](#casos-de-uso)

## 🎯 Visão Geral

Homestead segue uma **arquitetura em camadas** (Layered Architecture) com princípios de **Clean Architecture** e **Domain-Driven Design (DDD) simplificado**.

### Objetivos Arquiteturais

1. **Manutenibilidade** - Código fácil de entender e modificar
2. **Testabilidade** - Testes unitários e de integração simples
3. **Extensibilidade** - Fácil adicionar novos instaladores e features
4. **Separação de Responsabilidades** - Cada camada com propósito claro
5. **Independência de Framework** - Lógica de negócio não depende do TUI

## 🧱 Princípios Arquiteturais

### SOLID Principles

✅ **S - Single Responsibility Principle**

- Cada package/struct tem uma responsabilidade única
- Exemplo: `ScriptRunner` só executa, `ScriptRegistry` só registra

✅ **O - Open/Closed Principle**

- Aberto para extensão, fechado para modificação
- Novos installers via interface, sem modificar código existente

✅ **L - Liskov Substitution Principle**

- Implementações de interfaces são substituíveis
- Qualquer `Installer` pode ser usado onde a interface é esperada

✅ **I - Interface Segregation Principle**

- Interfaces pequenas e focadas
- `Installer`, `Executor`, `Registry` separadas

✅ **D - Dependency Inversion Principle**

- Dependa de abstrações, não de implementações
- TUI depende de interface `ScriptExecutor`, não da implementação

### Outros Princípios

- **DRY** (Don't Repeat Yourself) - Evitar duplicação
- **KISS** (Keep It Simple, Stupid) - Simplicidade sobre complexidade
- **YAGNI** (You Aren't Gonna Need It) - Implementar apenas o necessário

## 🏛️ Camadas da Aplicação

```
┌─────────────────────────────────────────────┐
│          Presentation Layer (TUI)           │
│         internal/tui/                       │
│  - Bubbletea Models                         │
│  - View Rendering                           │
│  - User Input Handling                      │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│          Application Layer                  │
│         internal/app/                       │
│  - Use Cases                                │
│  - Application Services                     │
│  - Orchestration                            │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│          Domain Layer                       │
│         internal/domain/                    │
│  - Entities (Script, Installer, Package)    │
│  - Interfaces (Repository, Executor)        │
│  - Domain Logic                             │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│          Infrastructure Layer               │
│    internal/infrastructure/                 │
│  - Script Execution (bash)                  │
│  - File System                              │
│  - External Commands (apt, snap, etc.)      │
└─────────────────────────────────────────────┘
```

### 1. Presentation Layer (TUI)

**Responsabilidade**: Interface com o usuário

**Localização**: `internal/tui/` (pacote raiz `tui`) e subpacotes auxiliares.

**Componentes** (ver [TUI_LAYOUT.md](TUI_LAYOUT.md)):

- **Raiz** — `model.go` (`Model`, `Update`, `Init`, `handleEnter`), `view_render.go`, `lists.go`, `menu.go`, `native_monitor.go`, wizards Zsh (`zsh_*_model.go`).
- `**internal/tui/cmds`** — `tea.Cmd` (catálogo, instalação, scripts, URLs).
- `**internal/tui/items**` — linhas `list.Item` (menu, scripts, pacotes).
- `**internal/tui/msg**` — tipos de mensagem Bubble Tea.
- `**internal/tui/theme**` — Lipgloss partilhado; `**internal/tui/sysurl**` — abrir/copiar URL sem ciclo com o raiz.

**Dependências**: Application Layer (services), tipos do Domain; o raiz também usa `internal/infrastructure/catalog` e `internal/monitoring` onde necessário.

### 2. Application Layer

**Responsabilidade**: Orquestração de fluxos que o utilizador dispara no TUI — validar inputs, compor chamadas ao domínio e à infraestrutura via interfaces.

**Localização**: `internal/app/services/`

**Componentes** (serviços de aplicação):

- `script_service.go` — carregar e executar scripts catalogados
- `installer_service.go` — instalação de pacotes do catálogo
- `config_service.go`, `plugin_service.go`, `wizard_service.go`, `repo_service.go` — configuração, Zsh e repositório de dotfiles

Não existe hoje uma pasta separada `usecases/`; a orquestração vive nestes serviços. Para novos fluxos grandes, pode-se extrair um use case explícito sem mudar a regra de dependência.

### 3. Domain Layer

**Responsabilidade**: Lógica de negócio, entidades e contratos (portas) independentes de UI e de IO.

**Localização**: `internal/domain/`

**Componentes**:

- `entities/` — `Script`, `Package`, e outras entidades do problema
- `interfaces/` — `ScriptExecutor`, repositórios, contratos de instalação
- `types/` — categorias, tipos de instalação, erros de domínio quando aplicável

**Exemplo**:

```go
// Domain Entity
type Script struct {
    ID          string
    Name        string
    Description string
    Category    Category
    RequiresSudo bool
}

// Domain Interface
type ScriptExecutor interface {
    Execute(script Script) error
    CanExecute(script Script) bool
}

// Domain Interface
type ScriptRepository interface {
    FindAll() ([]Script, error)
    FindByID(id string) (Script, error)
    FindByCategory(cat Category) ([]Script, error)
}
```

### 4. Infrastructure Layer

**Responsabilidade**: Implementações concretas e integrações externas (processos, ficheiros, rede, formatos).

**Localização**: `internal/infrastructure/`

**Componentes** (evolutivo; ver árvore do repositório):

- `executor/` — execução bash (`bash_executor.go`) e política sudo/TTY
- `repository/` — catálogos em memória, scripts utilitários, definições de instaladores
- `installer/` — estratégias de instalação e orquestração de pacotes
- `catalog/` — parse e carga de metadados (ex. JSON de instaladores)
- `config/`, `templates/`, `plugins/`, `preferences/` — persistência e extensão conforme o código atual

## 🎨 Padrões de Projeto

### 1. Repository Pattern

**Uso**: Acesso a dados (scripts, instaladores, configurações)

**Implementação**:

```go
// Domain interface
type ScriptRepository interface {
    FindAll() ([]Script, error)
    FindByID(id string) (Script, error)
    FindByCategory(category Category) ([]Script, error)
}

// Infrastructure implementation
type InMemoryScriptRepository struct {
    scripts []domain.Script
}

func (r *InMemoryScriptRepository) FindAll() ([]domain.Script, error) {
    return r.scripts, nil
}
```

**Vantagens**:

- Separa lógica de acesso a dados
- Facilita testes (mock repositories)
- Permite trocar implementação (in-memory → file → database)

### 2. Factory Pattern

**Uso**: Criar instaladores baseado no tipo

**Implementação**:

```go
type InstallerFactory interface {
    Create(packageType string) (Installer, error)
}

type DefaultInstallerFactory struct{}

func (f *DefaultInstallerFactory) Create(pkgType string) (Installer, error) {
    switch pkgType {
    case "apt":
        return &AptInstaller{}, nil
    case "snap":
        return &SnapInstaller{}, nil
    case "git":
        return &GitInstaller{}, nil
    default:
        return nil, fmt.Errorf("unknown package type: %s", pkgType)
    }
}
```

**Vantagens**:

- Centraliza criação de objetos
- Fácil adicionar novos tipos
- Esconde complexidade de criação

### 3. Strategy Pattern

**Uso**: Diferentes estratégias de instalação

**Implementação**:

```go
// Domain interface
type InstallStrategy interface {
    Install(pkg Package) error
    Uninstall(pkg Package) error
    IsInstalled(pkg Package) (bool, error)
}

// Concrete strategies
type AptStrategy struct{}
type SnapStrategy struct{}
type ManualStrategy struct{}

// Context
type Installer struct {
    strategy InstallStrategy
}

func (i *Installer) Install(pkg Package) error {
    return i.strategy.Install(pkg)
}
```

**Vantagens**:

- Algoritmos intercambiáveis
- Adicionar novas estratégias sem modificar código existente
- Testável individualmente

### 4. Command Pattern

**Uso**: Encapsular operações (útil para undo/redo, logging)

**Implementação**:

```go
type Command interface {
    Execute() error
    Undo() error
    Name() string
}

type InstallPackageCommand struct {
    installer Installer
    pkg       Package
}

func (c *InstallPackageCommand) Execute() error {
    return c.installer.Install(c.pkg)
}

func (c *InstallPackageCommand) Undo() error {
    return c.installer.Uninstall(c.pkg)
}

// Executor
type CommandExecutor struct {
    history []Command
}

func (e *CommandExecutor) Execute(cmd Command) error {
    err := cmd.Execute()
    if err == nil {
        e.history = append(e.history, cmd)
    }
    return err
}
```

**Vantagens**:

- Histórico de operações
- Suporte a undo
- Logging e auditoria

### 5. Observer Pattern

**Uso**: Notificar progresso de operações longas

**Implementação**:

```go
type ProgressObserver interface {
    OnProgress(current, total int, message string)
    OnComplete(success bool, message string)
}

type InstallOperation struct {
    observers []ProgressObserver
}

func (o *InstallOperation) NotifyProgress(current, total int, msg string) {
    for _, observer := range o.observers {
        observer.OnProgress(current, total, msg)
    }
}

// TUI implementa ProgressObserver
type TUIProgressObserver struct {
    model *Model
}

func (t *TUIProgressObserver) OnProgress(current, total int, msg string) {
    // Atualizar progress bar no TUI
}
```

**Vantagens**:

- Desacoplamento
- Múltiplos observers (TUI, logger, metrics)
- Fácil adicionar novos observers

### 6. Builder Pattern

**Uso**: Construir objetos complexos (configurações, wizard)

**Implementação**:

```go
type InstallerConfigBuilder struct {
    config InstallerConfig
}

func NewInstallerConfigBuilder() *InstallerConfigBuilder {
    return &InstallerConfigBuilder{
        config: InstallerConfig{},
    }
}

func (b *InstallerConfigBuilder) WithPackageManager(pm string) *InstallerConfigBuilder {
    b.config.PackageManager = pm
    return b
}

func (b *InstallerConfigBuilder) WithVersion(v string) *InstallerConfigBuilder {
    b.config.Version = v
    return b
}

func (b *InstallerConfigBuilder) Build() InstallerConfig {
    return b.config
}

// Uso
config := NewInstallerConfigBuilder().
    WithPackageManager("apt").
    WithVersion("latest").
    Build()
```

### 7. Adapter Pattern

**Uso**: Adaptar interfaces externas (apt, snap, docker CLI)

**Implementação**:

```go
// Interface que queremos
type PackageManager interface {
    Install(name string) error
    Remove(name string) error
    IsInstalled(name string) bool
}

// Adapter para apt
type AptAdapter struct {
    executor CommandExecutor
}

func (a *AptAdapter) Install(name string) error {
    return a.executor.Run("apt", "install", "-y", name)
}

// Adapter para snap
type SnapAdapter struct {
    executor CommandExecutor
}

func (s *SnapAdapter) Install(name string) error {
    return s.executor.Run("snap", "install", name)
}
```

## 📁 Estrutura de Diretórios Proposta

```
Homestead/
├── cmd/
│   └── homestead/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── domain/                     # Domain Layer
│   │   ├── entities/
│   │   │   ├── script.go
│   │   │   ├── installer.go
│   │   │   ├── package.go
│   │   │   └── system_state.go
│   │   ├── interfaces/
│   │   │   ├── executor.go
│   │   │   ├── repository.go
│   │   │   └── installer.go
│   │   └── types/
│   │       ├── category.go
│   │       └── errors.go
│   │
│   ├── app/                        # Application Layer
│   │   ├── usecases/
│   │   │   ├── execute_script.go
│   │   │   ├── install_package.go
│   │   │   ├── export_system.go
│   │   │   └── import_system.go
│   │   └── services/
│   │       ├── script_service.go
│   │       └── installer_service.go
│   │
│   ├── infrastructure/             # Infrastructure Layer
│   │   ├── executor/
│   │   │   ├── bash_executor.go
│   │   │   └── docker_executor.go
│   │   ├── repository/
│   │   │   ├── script_repository.go
│   │   │   └── package_repository.go
│   │   ├── packagemanager/
│   │   │   ├── apt.go
│   │   │   ├── snap.go
│   │   │   └── flatpak.go
│   │   └── fs/
│   │       └── file_system.go
│   │
│   ├── tui/                        # Presentation Layer (Bubble Tea)
│   │   ├── model.go
│   │   ├── view_render.go
│   │   ├── lists.go
│   │   ├── cmds/                   # tea.Cmd factories
│   │   ├── items/                  # list.Item implementations
│   │   ├── msg/                    # Bubble Tea message types
│   │   ├── theme/                  # Lipgloss styles
│   │   └── sysurl/                 # Open URL / clipboard helpers
│   │
│   ├── config/                     # Configuration
│   │   ├── config.go
│   │   └── loader.go
│   │
│   └── testutil/                   # Test utilities
│       └── testutil.go
│
├── scripts/                        # Bash scripts
│   ├── cleanup/
│   ├── monitoring/
│   └── install/
│
├── configs/                        # Configuration files
│   └── default.yaml
│
└── docs/                           # Documentation
    ├── ARCHITECTURE.md
    ├── PATTERNS.md
    └── API.md
```

## 📐 Convenções de Código

### Naming Conventions

**Packages**:

- Lowercase, singular
- Exemplo: `domain`, `executor`, `repository`

**Interfaces**:

- Substantivo ou adjetivo + "er"
- Exemplo: `Executor`, `Repository`, `Installer`

**Structs**:

- PascalCase
- Exemplo: `Script`, `AptInstaller`, `ScriptService`

**Methods**:

- PascalCase (públicos), camelCase (privados)
- Verbos no início
- Exemplo: `Execute()`, `Install()`, `loadConfig()`

### Error Handling

```go
// Sempre retornar erros, não usar panic
func (s *ScriptService) Execute(id string) error {
    script, err := s.repo.FindByID(id)
    if err != nil {
        return fmt.Errorf("failed to find script %s: %w", id, err)
    }

    if err := s.executor.Execute(script); err != nil {
        return fmt.Errorf("failed to execute script %s: %w", id, err)
    }

    return nil
}

// Definir erros customizados
var (
    ErrScriptNotFound = errors.New("script not found")
    ErrPermissionDenied = errors.New("permission denied")
)
```

### Dependency Injection

```go
// Constructor injection (preferido)
type ScriptService struct {
    repo     domain.ScriptRepository
    executor domain.ScriptExecutor
    logger   Logger
}

func NewScriptService(
    repo domain.ScriptRepository,
    executor domain.ScriptExecutor,
    logger Logger,
) *ScriptService {
    return &ScriptService{
        repo:     repo,
        executor: executor,
        logger:   logger,
    }
}

// main.go - wiring (detalhes em cmd/homestead/main.go)
func main() {
    // … preferences.Load, catalog.EffectiveCatalogURL, infra + services …
    model := tui.NewModel(scriptService, installerService, configService, repoService, catalogURL, prefs, prefsPath, catalogEnvSet)
    tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
}
```

### Testing Conventions

```go
// Mocks em _test.go
type MockExecutor struct {
    ExecuteFunc func(script domain.Script) error
}

func (m *MockExecutor) Execute(script domain.Script) error {
    if m.ExecuteFunc != nil {
        return m.ExecuteFunc(script)
    }
    return nil
}

// Table-driven tests
func TestScriptService_Execute(t *testing.T) {
    tests := []struct {
        name    string
        scriptID string
        wantErr bool
    }{
        {"success", "valid-id", false},
        {"not found", "invalid-id", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 🎯 Casos de Uso

### Use Case 1: Executar Script

```
1. Usuário seleciona script no menu
2. TUI chama ScriptService.Execute(scriptID)
3. ScriptService:
   a. Busca script no repository
   b. Valida permissões
   c. Chama executor.Execute(script)
4. Executor executa bash script
5. Resultado retorna ao TUI
6. TUI mostra sucesso/erro
```

### Use Case 2: Instalar Pacote

```
1. Usuário seleciona "Instaladores" → "Git"
2. TUI mostra wizard de configuração
3. Usuário confirma instalação
4. TUI chama InstallerService.Install(packageID, config)
5. InstallerService:
   a. Busca installer via Factory
   b. Verifica se já instalado
   c. Executa instalação
   d. Notifica progresso (Observer)
6. TUI atualiza progress bar
7. Instalação completa, TUI mostra sucesso
```

## 🔄 Próximos Passos

1. **Refatorar código atual** para seguir esta arquitetura
2. **Criar domain layer** com entities e interfaces
3. **Implementar application layer** com use cases
4. **Criar infrastructure** para installers
5. **Atualizar testes** para nova estrutura

---

**Última atualização**: 2026-03-14

Este documento evolui conforme o projeto cresce.