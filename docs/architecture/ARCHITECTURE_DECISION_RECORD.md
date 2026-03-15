# 📋 Architecture Decision Records (ADR)

Registro de decisões arquiteturais importantes do Homestead.

## ADR-001: Arquitetura em Camadas com Clean Architecture

**Data**: 2026-03-14
**Status**: Aceito
**Contexto**: Necessidade de definir estrutura base do projeto

### Decisão

Adotar **Arquitetura em Camadas** (Layered Architecture) com princípios de **Clean Architecture**:

1. **Presentation Layer** (TUI) - `internal/tui/`
2. **Application Layer** (Use Cases) - `internal/app/`
3. **Domain Layer** (Entities + Interfaces) - `internal/domain/`
4. **Infrastructure Layer** (Implementações) - `internal/infrastructure/`

### Razões

✅ **Separação de responsabilidades clara**
- Cada camada tem propósito único
- Fácil entender onde adicionar código

✅ **Testabilidade**
- Domain isolado de frameworks
- Mock de dependências simples
- Testes unitários focados

✅ **Manutenibilidade**
- Mudanças localizadas em uma camada
- Não quebra outras camadas

✅ **Extensibilidade**
- Adicionar features sem modificar base
- Trocar implementações (ex: apt → snap)

### Alternativas Consideradas

❌ **Monolito sem camadas**
- Rápido inicialmente
- Difícil manter quando cresce
- Testes complexos

❌ **Microservices**
- Over-engineering para CLI
- Complexidade desnecessária
- Overhead de comunicação

### Consequências

**Positivas**:
- Código organizado e previsível
- Fácil onboarding de novos devs
- Preparado para crescimento

**Negativas**:
- Mais boilerplate inicial
- Precisa disciplina para seguir
- Curva de aprendizado

---

## ADR-002: Uso de Interfaces para Inversão de Dependência

**Data**: 2026-03-14
**Status**: Aceito
**Contexto**: Como desacoplar camadas e facilitar testes

### Decisão

Definir **interfaces no Domain Layer** e implementar na Infrastructure:

```go
// Domain (interface)
type ScriptExecutor interface {
    Execute(script Script) error
}

// Infrastructure (implementation)
type BashExecutor struct {}
func (b *BashExecutor) Execute(script Script) error { ... }

// Application (uses interface)
type ScriptService struct {
    executor ScriptExecutor // interface, não concrete type
}
```

### Razões

✅ **Testabilidade**: Mock interfaces em testes
✅ **Flexibilidade**: Trocar implementações facilmente
✅ **SOLID**: Dependency Inversion Principle
✅ **Independência**: Domain não depende de infra

### Consequências

- Precisa definir interfaces antes de implementar
- Mais arquivos (interface + implementation)
- Mas vale a pena pela flexibilidade

---

## ADR-003: Repository Pattern para Acesso a Dados

**Data**: 2026-03-14
**Status**: Aceito
**Contexto**: Como gerenciar coleções de scripts/packages

### Decisão

Usar **Repository Pattern**:

```go
type ScriptRepository interface {
    FindAll() ([]Script, error)
    FindByID(id string) (*Script, error)
    FindByCategory(cat Category) ([]Script, error)
}
```

Inicialmente: **In-Memory Repository**
Futuro: File-based ou SQLite se necessário

### Razões

✅ **Abstração de persistência**
✅ **Fácil trocar backend** (memory → file → db)
✅ **Queries centralizadas**
✅ **Testável** (mock repository)

### Alternativas

❌ **Acesso direto** - Acoplamento forte
❌ **DAO** - Mais complexo, não necessário
❌ **ORM** - Overhead para CLI simples

---

## ADR-004: Factory Pattern para Installers

**Data**: 2026-03-14
**Status**: Aceito
**Contexto**: Criar instaladores de diferentes tipos (apt, snap, manual)

### Decisão

Usar **Factory Pattern**:

```go
type InstallerFactory interface {
    Create(packageType string) (Installer, error)
}
```

### Razões

✅ **Centraliza criação** de installers
✅ **Extensível** - adicionar novos tipos fácil
✅ **Type-safe** - retorna interface comum
✅ **Configurável** - pode usar config para customizar

### Quando Adicionar Novo Installer

1. Criar struct que implementa `Installer` interface
2. Adicionar case no Factory
3. Pronto! Nenhum código existente quebra

---

## ADR-005: Strategy Pattern para Métodos de Instalação

**Data**: 2026-03-14
**Status**: Aceito
**Contexto**: Mesmo package pode ser instalado de múltiplas formas

### Decisão

Usar **Strategy Pattern**:

```go
type InstallStrategy interface {
    Install(pkg Package) error
    CanInstall() bool
}

// Strategies: AptStrategy, SnapStrategy, FromSourceStrategy
```

### Razões

✅ **Algoritmos intercambiáveis**
✅ **Fallback automático** (tentar apt, depois snap, etc.)
✅ **Extensível** sem modificar código existente

### Exemplo

Docker pode ser instalado via:
1. Apt (preferido para Ubuntu)
2. Snap (fallback)
3. Script oficial (manual)

Service escolhe automaticamente melhor strategy disponível.

---

## ADR-006: Observer Pattern para Progresso

**Data**: 2026-03-14
**Status**: Aceito
**Contexto**: Notificar TUI sobre progresso de instalações longas

### Decisão

Usar **Observer Pattern**:

```go
type ProgressObserver interface {
    OnProgress(current, total int, msg string)
    OnComplete(success bool, msg string)
}
```

### Razões

✅ **Desacoplamento** - Installer não conhece TUI
✅ **Múltiplos observers** - TUI + Logger + Metrics
✅ **Real-time updates** - Progress bar responsivo

### Observers

1. **TUIObserver** - Atualiza progress bar
2. **LogObserver** - Escreve em arquivo
3. **MetricsObserver** - Envia para analytics (futuro)

---

## ADR-007: Command Pattern para Undo/Redo

**Data**: 2026-03-14
**Status**: Proposto
**Contexto**: Permitir reverter instalações

### Decisão

Usar **Command Pattern** para operações reversíveis:

```go
type Command interface {
    Execute() error
    Undo() error
}

type InstallCommand struct {
    installer Installer
    pkg Package
}
```

### Razões

✅ **Histórico** de operações
✅ **Undo inteligente** - só desinstala se instalou
✅ **Auditoria** - log de todas ações
✅ **Batch operations** - executar múltiplos comandos

### Implementação

**Fase 1** (Futuro): Histórico básico
**Fase 2** (Futuro): Undo/Redo completo
**Fase 3** (Futuro): Transações (rollback em erro)

---

## ADR-008: Builder Pattern para Configurações

**Data**: 2026-03-14
**Status**: Aceito
**Contexto**: Installers precisam configurações complexas

### Decisão

Usar **Builder Pattern**:

```go
config := NewInstallerConfigBuilder().
    WithPackageName("docker").
    WithVersion("latest").
    AddDependency("ca-certificates").
    EnableAutoStart().
    Build()
```

### Razões

✅ **Fluent interface** - legível
✅ **Validação** no Build()
✅ **Defaults** automáticos
✅ **Imutabilidade** - não modifica após Build()

### Uso

- Wizards no TUI
- Configs de YAML
- Presets (via Director)

---

## ADR-009: Dependency Injection Manual

**Data**: 2026-03-14
**Status**: Aceito
**Contexto**: Como conectar dependencies entre camadas

### Decisão

**Dependency Injection Manual** via constructores:

```go
// main.go
executor := executor.NewBashExecutor()
repo := repository.NewInMemoryScriptRepository()
service := services.NewScriptService(repo, executor)
model := tui.NewModel(service)
```

### Razões

✅ **Simplicidade** - Sem frameworks
✅ **Explícito** - Fácil de entender
✅ **Type-safe** - Compilador valida
✅ **Testável** - Injetar mocks

### Alternativas

❌ **Wire/Dig** - Over-engineering
❌ **Service Locator** - Anti-pattern
❌ **Global vars** - Dificulta testes

### Consequências

- Wiring em `main.go`
- Se crescer muito, considerar Wire

---

## ADR-010: Erros com Wrapping

**Data**: 2026-03-14
**Status**: Aceito
**Contexto**: Como lidar com erros em camadas

### Decisão

Usar **error wrapping** com `fmt.Errorf("%w")`:

```go
func (s *Service) Execute(id string) error {
    script, err := s.repo.FindByID(id)
    if err != nil {
        return fmt.Errorf("execute: find script %s: %w", id, err)
    }

    if err := s.executor.Execute(script); err != nil {
        return fmt.Errorf("execute: run script %s: %w", id, err)
    }

    return nil
}
```

### Razões

✅ **Contexto** preservado
✅ **Stack trace** via wrapping
✅ **errors.Is/As** funciona
✅ **Debugging** mais fácil

### Convenção

- Sempre wrap com contexto
- Use `%w` não `%v`
- Não log + return (escolha um)

---

## ADR-011: Logs Estruturados

**Data**: 2026-03-14
**Status**: Proposto (Futuro)
**Contexto**: Como fazer logging

### Decisão

Usar **structured logging** (ex: Charm Log):

```go
logger.Info("installing package",
    "package", pkgName,
    "version", version,
    "method", "apt",
)
```

### Razões

✅ **Parseável** - Fácil analisar
✅ **Consistente** - Formato fixo
✅ **Filtrável** - Por campos

### Níveis

- **Debug**: Detalhes internos
- **Info**: Operações importantes
- **Warn**: Problemas não-críticos
- **Error**: Erros que precisam atenção

---

## ADR-012: Configuração via YAML

**Data**: 2026-03-14
**Status**: Proposto (Futuro)
**Contexto**: Como users customizam Homestead

### Decisão

Configuração em `~/.config/homestead/config.yaml`:

```yaml
installer:
  preferred_method: apt
  auto_update: true
  backup_before_install: true

scripts:
  custom_directory: ~/my-scripts

ui:
  theme: dark
  confirm_destructive: true
```

### Razões

✅ **User-friendly** - YAML legível
✅ **Versionável** - Commit em dotfiles
✅ **Defaults** - Funciona sem config

### Implementação

1. Viper para parsing
2. Defaults em código
3. Override via env vars

---

## 📊 Resumo de Decisões

| ADR | Decisão | Status | Prioridade |
|-----|---------|--------|-----------|
| 001 | Layered Architecture | ✅ Aceito | Alta |
| 002 | Interfaces | ✅ Aceito | Alta |
| 003 | Repository Pattern | ✅ Aceito | Alta |
| 004 | Factory Pattern | ✅ Aceito | Média |
| 005 | Strategy Pattern | ✅ Aceito | Média |
| 006 | Observer Pattern | ✅ Aceito | Média |
| 007 | Command Pattern | 🔄 Proposto | Baixa |
| 008 | Builder Pattern | ✅ Aceito | Média |
| 009 | DI Manual | ✅ Aceito | Alta |
| 010 | Error Wrapping | ✅ Aceito | Alta |
| 011 | Structured Logging | 🔄 Proposto | Baixa |
| 012 | YAML Config | 🔄 Proposto | Baixa |

---

## 🎯 Próximas Decisões

1. **ADR-013**: Sistema de Plugins
2. **ADR-014**: Migração/Export de Sistema
3. **ADR-015**: Testes de Integração
4. **ADR-016**: CI/CD Pipeline

---

**Manutenção**: Adicionar novo ADR quando tomar decisão arquitetural importante

**Template**:
```markdown
## ADR-XXX: [Título]

**Data**: YYYY-MM-DD
**Status**: [Proposto/Aceito/Rejeitado/Obsoleto]
**Contexto**: [Problema/situação]

### Decisão
[O que decidimos]

### Razões
[Por que decidimos assim]

### Alternativas Consideradas
[O que mais consideramos]

### Consequências
[Impactos da decisão]
```
