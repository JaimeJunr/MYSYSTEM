# 🏗️ Arquitetura de Software - Homestead

Documento de arquitetura e padrões de projeto para o Homestead.

## 📋 Índice

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

**Localização**: `internal/tui/`

**Componentes**:
- `model.go` - Bubbletea Model (State + Update + View)
- `views/` - Views específicas (menu, script list, installer wizard)
- `components/` - Componentes reutilizáveis (progress bar, input forms)

**Dependências**: Domain Layer (interfaces), Application Layer

### 2. Application Layer

**Responsabilidade**: Casos de uso e orquestração

**Localização**: `internal/app/` (a criar)

**Componentes**:
- `usecases/` - Use cases específicos
  - `execute_script.go`
  - `install_package.go`
  - `export_system.go`
- `services/` - Application services
  - `script_service.go`
  - `installer_service.go`

**Exemplo**:
```go
type ExecuteScriptUseCase struct {
    executor domain.ScriptExecutor
}

func (uc *ExecuteScriptUseCase) Execute(scriptID string) error {
    // 1. Validar input
    // 2. Buscar script
    // 3. Executar
    // 4. Logar resultado
}
```

### 3. Domain Layer

**Responsabilidade**: Lógica de negócio e entidades

**Localização**: `internal/domain/` (a criar)

**Componentes**:
- `entities/` - Entidades do domínio
  - `script.go`
  - `installer.go`
  - `package.go`
  - `system_state.go`
- `interfaces/` - Contratos
  - `executor.go`
  - `repository.go`
  - `installer.go`

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

**Responsabilidade**: Implementações concretas e integrações externas

**Localização**: `internal/infrastructure/` (a criar)

**Componentes**:
- `executor/` - Executores concretos
  - `bash_executor.go`
  - `docker_executor.go`
- `repository/` - Repositórios concretos
  - `script_repository.go` (in-memory ou file-based)
  - `package_repository.go`
- `apt/` - Integração com apt
- `snap/` - Integração com snap
- `fs/` - File system operations

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
│   ├── tui/                        # Presentation Layer
│   │   ├── model.go
│   │   ├── views/
│   │   │   ├── main_menu.go
│   │   │   ├── script_list.go
│   │   │   └── installer_wizard.go
│   │   └── components/
│   │       ├── progress.go
│   │       └── form.go
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

// main.go - wiring
func main() {
    // Infrastructure
    repo := repository.NewInMemoryScriptRepository()
    executor := executor.NewBashExecutor()
    logger := log.New()

    // Application
    scriptService := app.NewScriptService(repo, executor, logger)

    // Presentation
    model := tui.InitialModel(scriptService)

    // Run
    tea.NewProgram(model).Run()
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
