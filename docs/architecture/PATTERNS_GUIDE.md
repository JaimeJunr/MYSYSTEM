# 📚 Guia Prático de Padrões - Homestead

Exemplos práticos de como implementar cada padrão de projeto no Homestead.

Para **onde cada padrão encaixa na pilha** (TUI → serviços → domínio → infra), veja primeiro [VERTICAL_PATTERNS.md](VERTICAL_PATTERNS.md).

## 📋 Índice

- [Repository Pattern](#repository-pattern)
- [Factory Pattern](#factory-pattern)
- [Strategy Pattern](#strategy-pattern)
- [Command Pattern](#command-pattern)
- [Observer Pattern](#observer-pattern)
- [Builder Pattern](#builder-pattern)
- [Adapter Pattern](#adapter-pattern)

---

## Repository Pattern

### Quando Usar
- Precisar acessar/gerenciar coleções de entidades
- Separar lógica de acesso a dados da lógica de negócio
- Facilitar testes com mock repositories

### Implementação Completa

#### 1. Domain Interface
```go
// internal/domain/interfaces/repository.go
package interfaces

import "github.com/JaimeJunr/Homestead/internal/domain/entities"

type ScriptRepository interface {
    // Queries
    FindAll() ([]entities.Script, error)
    FindByID(id string) (*entities.Script, error)
    FindByCategory(category string) ([]entities.Script, error)

    // Commands
    Save(script *entities.Script) error
    Delete(id string) error

    // Checks
    Exists(id string) bool
}
```

#### 2. Infrastructure Implementation
```go
// internal/infrastructure/repository/script_repository.go
package repository

import (
    "fmt"
    "sync"

    "github.com/JaimeJunr/Homestead/internal/domain/entities"
    "github.com/JaimeJunr/Homestead/internal/domain/interfaces"
)

type InMemoryScriptRepository struct {
    scripts map[string]*entities.Script
    mu      sync.RWMutex
}

func NewInMemoryScriptRepository() interfaces.ScriptRepository {
    return &InMemoryScriptRepository{
        scripts: make(map[string]*entities.Script),
    }
}

func (r *InMemoryScriptRepository) FindAll() ([]entities.Script, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    scripts := make([]entities.Script, 0, len(r.scripts))
    for _, script := range r.scripts {
        scripts = append(scripts, *script)
    }

    return scripts, nil
}

func (r *InMemoryScriptRepository) FindByID(id string) (*entities.Script, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    script, ok := r.scripts[id]
    if !ok {
        return nil, fmt.Errorf("script not found: %s", id)
    }

    return script, nil
}

func (r *InMemoryScriptRepository) FindByCategory(category string) ([]entities.Script, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    var scripts []entities.Script
    for _, script := range r.scripts {
        if script.Category == category {
            scripts = append(scripts, *script)
        }
    }

    return scripts, nil
}

func (r *InMemoryScriptRepository) Save(script *entities.Script) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.scripts[script.ID] = script
    return nil
}

func (r *InMemoryScriptRepository) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    delete(r.scripts, id)
    return nil
}

func (r *InMemoryScriptRepository) Exists(id string) bool {
    r.mu.RLock()
    defer r.mu.RUnlock()

    _, ok := r.scripts[id]
    return ok
}
```

#### 3. Uso no Application Service
```go
// internal/app/services/script_service.go
package services

import (
    "github.com/JaimeJunr/Homestead/internal/domain/interfaces"
)

type ScriptService struct {
    repo interfaces.ScriptRepository
}

func NewScriptService(repo interfaces.ScriptRepository) *ScriptService {
    return &ScriptService{repo: repo}
}

func (s *ScriptService) GetScriptsByCategory(category string) ([]entities.Script, error) {
    return s.repo.FindByCategory(category)
}
```

#### 4. Mock para Testes
```go
// internal/app/services/script_service_test.go
package services_test

type MockScriptRepository struct {
    FindAllFunc        func() ([]entities.Script, error)
    FindByIDFunc       func(id string) (*entities.Script, error)
    FindByCategoryFunc func(category string) ([]entities.Script, error)
}

func (m *MockScriptRepository) FindAll() ([]entities.Script, error) {
    if m.FindAllFunc != nil {
        return m.FindAllFunc()
    }
    return nil, nil
}

// Test
func TestScriptService_GetScriptsByCategory(t *testing.T) {
    mockRepo := &MockScriptRepository{
        FindByCategoryFunc: func(cat string) ([]entities.Script, error) {
            return []entities.Script{{ID: "test"}}, nil
        },
    }

    service := NewScriptService(mockRepo)
    scripts, err := service.GetScriptsByCategory("cleanup")

    assert.NoError(t, err)
    assert.Len(t, scripts, 1)
}
```

---

## Factory Pattern

### Quando Usar
- Criar objetos complexos
- Tipo de objeto definido em runtime
- Centralizar lógica de criação

### Implementação Completa

#### 1. Domain Interface
```go
// internal/domain/interfaces/installer.go
package interfaces

type Installer interface {
    Install(packageName string) error
    Uninstall(packageName string) error
    IsInstalled(packageName string) (bool, error)
    GetInfo(packageName string) (*PackageInfo, error)
}

type PackageInfo struct {
    Name        string
    Version     string
    Description string
    Installed   bool
}
```

#### 2. Factory Interface
```go
// internal/domain/interfaces/factory.go
package interfaces

type InstallerFactory interface {
    Create(packageType string) (Installer, error)
    Supports(packageType string) bool
    ListSupported() []string
}
```

#### 3. Concrete Implementations
```go
// internal/infrastructure/installer/apt_installer.go
package installer

type AptInstaller struct {
    executor CommandExecutor
}

func NewAptInstaller(executor CommandExecutor) *AptInstaller {
    return &AptInstaller{executor: executor}
}

func (a *AptInstaller) Install(packageName string) error {
    return a.executor.Run("sudo", "apt", "install", "-y", packageName)
}

// internal/infrastructure/installer/snap_installer.go
type SnapInstaller struct {
    executor CommandExecutor
}

func NewSnapInstaller(executor CommandExecutor) *SnapInstaller {
    return &SnapInstaller{executor: executor}
}

func (s *SnapInstaller) Install(packageName string) error {
    return s.executor.Run("sudo", "snap", "install", packageName)
}
```

#### 4. Factory Implementation
```go
// internal/infrastructure/factory/installer_factory.go
package factory

import (
    "fmt"

    "github.com/JaimeJunr/Homestead/internal/domain/interfaces"
    "github.com/JaimeJunr/Homestead/internal/infrastructure/installer"
)

type DefaultInstallerFactory struct {
    executor CommandExecutor
}

func NewInstallerFactory(executor CommandExecutor) interfaces.InstallerFactory {
    return &DefaultInstallerFactory{executor: executor}
}

func (f *DefaultInstallerFactory) Create(packageType string) (interfaces.Installer, error) {
    switch packageType {
    case "apt":
        return installer.NewAptInstaller(f.executor), nil
    case "snap":
        return installer.NewSnapInstaller(f.executor), nil
    case "flatpak":
        return installer.NewFlatpakInstaller(f.executor), nil
    case "manual":
        return installer.NewManualInstaller(f.executor), nil
    default:
        return nil, fmt.Errorf("unsupported package type: %s", packageType)
    }
}

func (f *DefaultInstallerFactory) Supports(packageType string) bool {
    supported := []string{"apt", "snap", "flatpak", "manual"}
    for _, s := range supported {
        if s == packageType {
            return true
        }
    }
    return false
}

func (f *DefaultInstallerFactory) ListSupported() []string {
    return []string{"apt", "snap", "flatpak", "manual"}
}
```

#### 5. Uso
```go
// internal/app/usecases/install_package.go
package usecases

type InstallPackageUseCase struct {
    factory interfaces.InstallerFactory
}

func (uc *InstallPackageUseCase) Execute(packageType, packageName string) error {
    installer, err := uc.factory.Create(packageType)
    if err != nil {
        return err
    }

    return installer.Install(packageName)
}
```

---

## Strategy Pattern

### Quando Usar
- Múltiplos algoritmos para mesma operação
- Trocar comportamento em runtime
- Evitar múltiplos if/switch

### Implementação Completa

#### 1. Strategy Interface
```go
// internal/domain/interfaces/install_strategy.go
package interfaces

type InstallStrategy interface {
    CanInstall() bool
    Install(pkg *Package) error
    Uninstall(pkg *Package) error
    Upgrade(pkg *Package) error
}
```

#### 2. Concrete Strategies
```go
// internal/infrastructure/strategy/apt_strategy.go
package strategy

type AptStrategy struct {
    executor CommandExecutor
}

func (s *AptStrategy) CanInstall() bool {
    // Check if apt is available
    return s.executor.CommandExists("apt")
}

func (s *AptStrategy) Install(pkg *Package) error {
    return s.executor.Run("sudo", "apt", "install", "-y", pkg.Name)
}

// internal/infrastructure/strategy/snap_strategy.go
type SnapStrategy struct {
    executor CommandExecutor
}

func (s *SnapStrategy) CanInstall() bool {
    return s.executor.CommandExists("snap")
}

func (s *SnapStrategy) Install(pkg *Package) error {
    return s.executor.Run("sudo", "snap", "install", pkg.Name)
}

// internal/infrastructure/strategy/from_source_strategy.go
type FromSourceStrategy struct {
    executor CommandExecutor
}

func (s *FromSourceStrategy) Install(pkg *Package) error {
    // Complex multi-step installation
    steps := [][]string{
        {"git", "clone", pkg.SourceURL},
        {"cd", pkg.Name},
        {"./configure"},
        {"make"},
        {"sudo", "make", "install"},
    }

    for _, step := range steps {
        if err := s.executor.Run(step[0], step[1:]...); err != nil {
            return err
        }
    }

    return nil
}
```

#### 3. Context
```go
// internal/app/services/installer_service.go
package services

type InstallerService struct {
    strategies map[string]interfaces.InstallStrategy
}

func NewInstallerService() *InstallerService {
    return &InstallerService{
        strategies: make(map[string]interfaces.InstallStrategy),
    }
}

func (s *InstallerService) RegisterStrategy(name string, strategy interfaces.InstallStrategy) {
    s.strategies[name] = strategy
}

func (s *InstallerService) Install(pkg *Package) error {
    // Select best strategy
    var strategy interfaces.InstallStrategy

    // Try preferred strategy first
    if s, ok := s.strategies[pkg.PreferredMethod]; ok && s.CanInstall() {
        strategy = s
    } else {
        // Fallback to first available
        for _, s := range s.strategies {
            if s.CanInstall() {
                strategy = s
                break
            }
        }
    }

    if strategy == nil {
        return fmt.Errorf("no available installation strategy")
    }

    return strategy.Install(pkg)
}
```

#### 4. Configuration
```go
// cmd/homestead/main.go
func setupInstallerService() *services.InstallerService {
    executor := executor.NewBashExecutor()

    service := services.NewInstallerService()
    service.RegisterStrategy("apt", strategy.NewAptStrategy(executor))
    service.RegisterStrategy("snap", strategy.NewSnapStrategy(executor))
    service.RegisterStrategy("source", strategy.NewFromSourceStrategy(executor))

    return service
}
```

---

## Command Pattern

### Quando Usar
- Precisar de undo/redo
- Enfileirar operações
- Logging de ações
- Histórico de operações

### Implementação Completa

#### 1. Command Interface
```go
// internal/domain/interfaces/command.go
package interfaces

type Command interface {
    Execute() error
    Undo() error
    Name() string
    Description() string
}
```

#### 2. Concrete Commands
```go
// internal/app/commands/install_package_command.go
package commands

type InstallPackageCommand struct {
    installer  interfaces.Installer
    packageName string
    wasInstalled bool // Para undo inteligente
}

func NewInstallPackageCommand(installer interfaces.Installer, pkg string) Command {
    return &InstallPackageCommand{
        installer:   installer,
        packageName: pkg,
    }
}

func (c *InstallPackageCommand) Execute() error {
    // Check if already installed
    installed, err := c.installer.IsInstalled(c.packageName)
    if err != nil {
        return err
    }

    c.wasInstalled = installed

    if installed {
        return nil // Already installed, skip
    }

    return c.installer.Install(c.packageName)
}

func (c *InstallPackageCommand) Undo() error {
    // Only uninstall if we installed it
    if !c.wasInstalled {
        return c.installer.Uninstall(c.packageName)
    }
    return nil
}

func (c *InstallPackageCommand) Name() string {
    return fmt.Sprintf("install-%s", c.packageName)
}

func (c *InstallPackageCommand) Description() string {
    return fmt.Sprintf("Install package: %s", c.packageName)
}
```

#### 3. Command Executor/Invoker
```go
// internal/app/services/command_executor.go
package services

type CommandExecutor struct {
    history []interfaces.Command
    current int
}

func NewCommandExecutor() *CommandExecutor {
    return &CommandExecutor{
        history: make([]interfaces.Command, 0),
        current: -1,
    }
}

func (e *CommandExecutor) Execute(cmd interfaces.Command) error {
    // Execute command
    if err := cmd.Execute(); err != nil {
        return err
    }

    // Add to history (remove redo stack if exists)
    e.history = e.history[:e.current+1]
    e.history = append(e.history, cmd)
    e.current++

    return nil
}

func (e *CommandExecutor) Undo() error {
    if e.current < 0 {
        return fmt.Errorf("nothing to undo")
    }

    cmd := e.history[e.current]
    if err := cmd.Undo(); err != nil {
        return err
    }

    e.current--
    return nil
}

func (e *CommandExecutor) Redo() error {
    if e.current >= len(e.history)-1 {
        return fmt.Errorf("nothing to redo")
    }

    e.current++
    cmd := e.history[e.current]
    return cmd.Execute()
}

func (e *CommandExecutor) History() []interfaces.Command {
    return e.history[:e.current+1]
}
```

#### 4. Uso
```go
// Example usage
executor := NewCommandExecutor()

// Install git
gitCmd := NewInstallPackageCommand(aptInstaller, "git")
executor.Execute(gitCmd)

// Install docker
dockerCmd := NewInstallPackageCommand(aptInstaller, "docker")
executor.Execute(dockerCmd)

// Oops, undo docker
executor.Undo()

// Changed my mind, redo
executor.Redo()

// View history
for _, cmd := range executor.History() {
    fmt.Println(cmd.Description())
}
```

---

## Observer Pattern

### Quando Usar
- Notificar progresso de operações longas
- Múltiplos componentes interessados em eventos
- Logging e metrics

### Implementação Completa

#### 1. Observer Interface
```go
// internal/domain/interfaces/observer.go
package interfaces

type ProgressObserver interface {
    OnStart(operation string)
    OnProgress(current, total int, message string)
    OnComplete(success bool, message string)
    OnError(err error)
}
```

#### 2. Observable/Subject
```go
// internal/app/services/observable_installer.go
package services

type ObservableInstaller struct {
    installer interfaces.Installer
    observers []interfaces.ProgressObserver
}

func NewObservableInstaller(installer interfaces.Installer) *ObservableInstaller {
    return &ObservableInstaller{
        installer: installer,
        observers: make([]interfaces.ProgressObserver, 0),
    }
}

func (o *ObservableInstaller) AddObserver(observer interfaces.ProgressObserver) {
    o.observers = append(o.observers, observer)
}

func (o *ObservableInstaller) RemoveObserver(observer interfaces.ProgressObserver) {
    for i, obs := range o.observers {
        if obs == observer {
            o.observers = append(o.observers[:i], o.observers[i+1:]...)
            break
        }
    }
}

func (o *ObservableInstaller) Install(packageName string) error {
    // Notify start
    o.notifyStart(fmt.Sprintf("Installing %s", packageName))

    // Simulate steps with progress
    steps := []string{
        "Downloading package...",
        "Verifying dependencies...",
        "Installing files...",
        "Configuring package...",
    }

    for i, step := range steps {
        o.notifyProgress(i+1, len(steps), step)
        time.Sleep(time.Second) // Simulate work
    }

    // Actual installation
    err := o.installer.Install(packageName)

    if err != nil {
        o.notifyError(err)
        o.notifyComplete(false, fmt.Sprintf("Failed to install %s", packageName))
        return err
    }

    o.notifyComplete(true, fmt.Sprintf("Successfully installed %s", packageName))
    return nil
}

func (o *ObservableInstaller) notifyStart(operation string) {
    for _, observer := range o.observers {
        observer.OnStart(operation)
    }
}

func (o *ObservableInstaller) notifyProgress(current, total int, msg string) {
    for _, observer := range o.observers {
        observer.OnProgress(current, total, msg)
    }
}

func (o *ObservableInstaller) notifyComplete(success bool, msg string) {
    for _, observer := range o.observers {
        observer.OnComplete(success, msg)
    }
}

func (o *ObservableInstaller) notifyError(err error) {
    for _, observer := range o.observers {
        observer.OnError(err)
    }
}
```

#### 3. Concrete Observers
```go
// internal/tui/observers/progress_observer.go
package observers

type TUIProgressObserver struct {
    model *tui.Model
}

func NewTUIProgressObserver(model *tui.Model) *TUIProgressObserver {
    return &TUIProgressObserver{model: model}
}

func (o *TUIProgressObserver) OnStart(operation string) {
    o.model.SetStatus(operation)
    o.model.ShowProgressBar(true)
}

func (o *TUIProgressObserver) OnProgress(current, total int, message string) {
    percent := float64(current) / float64(total)
    o.model.UpdateProgress(percent, message)
}

func (o *TUIProgressObserver) OnComplete(success bool, message string) {
    o.model.ShowProgressBar(false)
    o.model.SetStatus(message)
}

func (o *TUIProgressObserver) OnError(err error) {
    o.model.ShowError(err.Error())
}

// internal/infrastructure/observers/log_observer.go
package observers

type LogObserver struct {
    logger Logger
}

func NewLogObserver(logger Logger) *LogObserver {
    return &LogObserver{logger: logger}
}

func (o *LogObserver) OnStart(operation string) {
    o.logger.Info("Started: " + operation)
}

func (o *LogObserver) OnProgress(current, total int, message string) {
    o.logger.Debug(fmt.Sprintf("[%d/%d] %s", current, total, message))
}

func (o *LogObserver) OnComplete(success bool, message string) {
    if success {
        o.logger.Info("Completed: " + message)
    } else {
        o.logger.Error("Failed: " + message)
    }
}

func (o *LogObserver) OnError(err error) {
    o.logger.Error("Error: " + err.Error())
}
```

#### 4. Uso
```go
// Setup
installer := NewAptInstaller(executor)
observable := NewObservableInstaller(installer)

// Add observers
tuiObserver := NewTUIProgressObserver(model)
logObserver := NewLogObserver(logger)

observable.AddObserver(tuiObserver)
observable.AddObserver(logObserver)

// Execute - both observers will be notified
observable.Install("git")
```

---

## Builder Pattern

### Quando Usar
- Objetos complexos com muitos parâmetros
- Wizards e forms
- Configurações complexas

### Implementação Completa

#### 1. Domain Entity
```go
// internal/domain/entities/installer_config.go
package entities

type InstallerConfig struct {
    PackageName    string
    PackageType    string
    Version        string
    CustomOptions  map[string]string
    Dependencies   []string
    PostInstall    []string
    Backup         bool
    AutoStart      bool
}
```

#### 2. Builder
```go
// internal/app/builders/installer_config_builder.go
package builders

type InstallerConfigBuilder struct {
    config *entities.InstallerConfig
}

func NewInstallerConfigBuilder() *InstallerConfigBuilder {
    return &InstallerConfigBuilder{
        config: &entities.InstallerConfig{
            CustomOptions: make(map[string]string),
            Dependencies:  make([]string, 0),
            PostInstall:   make([]string, 0),
        },
    }
}

func (b *InstallerConfigBuilder) WithPackageName(name string) *InstallerConfigBuilder {
    b.config.PackageName = name
    return b
}

func (b *InstallerConfigBuilder) WithPackageType(pkgType string) *InstallerConfigBuilder {
    b.config.PackageType = pkgType
    return b
}

func (b *InstallerConfigBuilder) WithVersion(version string) *InstallerConfigBuilder {
    b.config.Version = version
    return b
}

func (b *InstallerConfigBuilder) WithOption(key, value string) *InstallerConfigBuilder {
    b.config.CustomOptions[key] = value
    return b
}

func (b *InstallerConfigBuilder) AddDependency(dep string) *InstallerConfigBuilder {
    b.config.Dependencies = append(b.config.Dependencies, dep)
    return b
}

func (b *InstallerConfigBuilder) AddPostInstallScript(script string) *InstallerConfigBuilder {
    b.config.PostInstall = append(b.config.PostInstall, script)
    return b
}

func (b *InstallerConfigBuilder) EnableBackup() *InstallerConfigBuilder {
    b.config.Backup = true
    return b
}

func (b *InstallerConfigBuilder) EnableAutoStart() *InstallerConfigBuilder {
    b.config.AutoStart = true
    return b
}

func (b *InstallerConfigBuilder) Build() (*entities.InstallerConfig, error) {
    // Validation
    if b.config.PackageName == "" {
        return nil, fmt.Errorf("package name is required")
    }
    if b.config.PackageType == "" {
        return nil, fmt.Errorf("package type is required")
    }

    return b.config, nil
}
```

#### 3. Director (Optional)
```go
// internal/app/builders/config_director.go
package builders

type ConfigDirector struct {
    builder *InstallerConfigBuilder
}

func NewConfigDirector() *ConfigDirector {
    return &ConfigDirector{
        builder: NewInstallerConfigBuilder(),
    }
}

// Preset configurations
func (d *ConfigDirector) BuildDockerConfig() (*entities.InstallerConfig, error) {
    return d.builder.
        WithPackageName("docker").
        WithPackageType("apt").
        WithVersion("latest").
        AddDependency("ca-certificates").
        AddDependency("curl").
        AddDependency("gnupg").
        AddPostInstallScript("sudo usermod -aG docker $USER").
        EnableAutoStart().
        Build()
}

func (d *ConfigDirector) BuildNodeJSConfig(version string) (*entities.InstallerConfig, error) {
    return d.builder.
        WithPackageName("nodejs").
        WithPackageType("manual").
        WithVersion(version).
        WithOption("source", "https://nodejs.org/dist").
        AddDependency("build-essential").
        EnableBackup().
        Build()
}
```

#### 4. Uso
```go
// Simple usage
config, err := NewInstallerConfigBuilder().
    WithPackageName("git").
    WithPackageType("apt").
    WithVersion("latest").
    Build()

// With director (preset)
director := NewConfigDirector()
dockerConfig, err := director.BuildDockerConfig()

// Complex usage
config, err := NewInstallerConfigBuilder().
    WithPackageName("postgresql").
    WithPackageType("apt").
    WithVersion("14").
    WithOption("locale", "en_US.UTF-8").
    AddDependency("postgresql-contrib").
    AddDependency("postgresql-client").
    AddPostInstallScript("sudo systemctl enable postgresql").
    AddPostInstallScript("sudo -u postgres createuser myapp").
    EnableBackup().
    EnableAutoStart().
    Build()
```

---

## Adapter Pattern

### Quando Usar
- Integrar sistemas externos
- Unificar interfaces diferentes
- Wrapper de bibliotecas terceiras

### Implementação Completa

#### 1. Target Interface
```go
// internal/domain/interfaces/package_manager.go
package interfaces

type PackageManager interface {
    // Queries
    Search(query string) ([]PackageInfo, error)
    Info(packageName string) (*PackageInfo, error)
    ListInstalled() ([]PackageInfo, error)
    IsInstalled(packageName string) (bool, error)

    // Commands
    Install(packageName string, options ...InstallOption) error
    Uninstall(packageName string) error
    Update(packageName string) error
    UpdateAll() error
}

type PackageInfo struct {
    Name        string
    Version     string
    Description string
    Size        int64
    Installed   bool
}

type InstallOption func(*InstallConfig)

type InstallConfig struct {
    AssumeYes  bool
    NoRecommends bool
    Version    string
}
```

#### 2. Adaptees (External systems)
```go
// Apt command-line tool
// Snap command-line tool
// Flatpak command-line tool
// Each has different CLI interface
```

#### 3. Adapters
```go
// internal/infrastructure/adapter/apt_adapter.go
package adapter

type AptAdapter struct {
    executor CommandExecutor
}

func NewAptAdapter(executor CommandExecutor) *AptAdapter {
    return &AptAdapter{executor: executor}
}

func (a *AptAdapter) Search(query string) ([]interfaces.PackageInfo, error) {
    output, err := a.executor.Output("apt", "search", query)
    if err != nil {
        return nil, err
    }

    return a.parseSearchOutput(output), nil
}

func (a *AptAdapter) Install(packageName string, options ...interfaces.InstallOption) error {
    config := &interfaces.InstallConfig{}
    for _, opt := range options {
        opt(config)
    }

    args := []string{"install"}
    if config.AssumeYes {
        args = append(args, "-y")
    }
    if config.NoRecommends {
        args = append(args, "--no-install-recommends")
    }
    if config.Version != "" {
        packageName = fmt.Sprintf("%s=%s", packageName, config.Version)
    }
    args = append(args, packageName)

    return a.executor.Run("sudo", append([]string{"apt"}, args...)...)
}

func (a *AptAdapter) parseSearchOutput(output string) []interfaces.PackageInfo {
    // Parse apt search output format
    var packages []interfaces.PackageInfo
    // Implementation...
    return packages
}

// internal/infrastructure/adapter/snap_adapter.go
package adapter

type SnapAdapter struct {
    executor CommandExecutor
}

func NewSnapAdapter(executor CommandExecutor) *SnapAdapter {
    return &SnapAdapter{executor: executor}
}

func (s *SnapAdapter) Search(query string) ([]interfaces.PackageInfo, error) {
    output, err := s.executor.Output("snap", "find", query)
    if err != nil {
        return nil, err
    }

    return s.parseSearchOutput(output), nil
}

func (s *SnapAdapter) Install(packageName string, options ...interfaces.InstallOption) error {
    // Snap has different options, adapt them
    args := []string{"install", packageName}

    return s.executor.Run("sudo", append([]string{"snap"}, args...)...)
}

func (s *SnapAdapter) parseSearchOutput(output string) []interfaces.PackageInfo {
    // Parse snap find output format (different from apt!)
    var packages []interfaces.PackageInfo
    // Implementation...
    return packages
}
```

#### 4. Uso
```go
// Create adapters
aptAdapter := adapter.NewAptAdapter(executor)
snapAdapter := adapter.NewSnapAdapter(executor)

// Use unified interface
var pm interfaces.PackageManager

// Use apt
pm = aptAdapter
results, _ := pm.Search("docker")
pm.Install("docker.io",
    func(c *InstallConfig) { c.AssumeYes = true },
)

// Switch to snap (same interface!)
pm = snapAdapter
results, _ = pm.Search("docker")
pm.Install("docker")

// Can abstract further
type MultiPackageManager struct {
    managers []interfaces.PackageManager
}

func (m *MultiPackageManager) Install(pkg string) error {
    // Try each manager until one succeeds
    for _, pm := range m.managers {
        if err := pm.Install(pkg); err == nil {
            return nil
        }
    }
    return fmt.Errorf("all managers failed")
}
```

---

## 🎯 Resumo de Quando Usar

| Padrão | Quando Usar | Exemplo no Homestead |
|--------|-------------|---------------------|
| **Repository** | Acesso a dados | Scripts, Packages, Configs |
| **Factory** | Criar objetos variados | Criar installers por tipo |
| **Strategy** | Algoritmos intercambiáveis | Métodos de instalação |
| **Command** | Undo/redo, histórico | Instalações reversíveis |
| **Observer** | Notificar progresso | TUI + Logger + Metrics |
| **Builder** | Objetos complexos | Configurações de instalação |
| **Adapter** | Unificar interfaces externas | apt, snap, flatpak |

---

**Última atualização**: 2026-03-14

Para mais detalhes, veja [ARCHITECTURE.md](../ARCHITECTURE.md)
