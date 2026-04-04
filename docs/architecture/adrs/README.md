# Architecture Decision Records (ADR)

Registro de decisões arquiteturais importantes do Homestead. Cada ADR está num ficheiro dedicado nesta pasta.

## Índice

| ADR | Título | Status |
| --- | --- | --- |
| [001](001-layered-clean-architecture.md) | Arquitetura em camadas com Clean Architecture | Aceito |
| [002](002-interfaces-dependency-inversion.md) | Interfaces para inversão de dependência | Aceito |
| [003](003-repository-pattern.md) | Repository pattern para acesso a dados | Aceito |
| [004](004-factory-installers.md) | Factory pattern para installers | Aceito |
| [005](005-strategy-installation.md) | Strategy pattern para métodos de instalação | Aceito |
| [006](006-observer-progress.md) | Observer pattern para progresso | Aceito |
| [007](007-command-undo-redo.md) | Command pattern para undo/redo | Proposto |
| [008](008-builder-configuration.md) | Builder pattern para configurações | Aceito |
| [009](009-dependency-injection-manual.md) | Dependency injection manual | Aceito |
| [010](010-error-wrapping.md) | Erros com wrapping | Aceito |
| [011](011-structured-logging.md) | Logs estruturados | Proposto (futuro) |
| [012](012-yaml-configuration.md) | Configuração via YAML | Proposto (futuro) |
| [013](013-remote-installer-catalog.md) | Catálogo remoto de instaladores (JSON + cache) | Aceito |

## Resumo por prioridade (legado)

| ADR | Decisão | Status | Prioridade |
| --- | --- | --- | --- |
| 001 | Layered Architecture | Aceito | Alta |
| 002 | Interfaces | Aceito | Alta |
| 003 | Repository Pattern | Aceito | Alta |
| 004 | Factory Pattern | Aceito | Média |
| 005 | Strategy Pattern | Aceito | Média |
| 006 | Observer Pattern | Aceito | Média |
| 007 | Command Pattern | Proposto | Baixa |
| 008 | Builder Pattern | Aceito | Média |
| 009 | DI Manual | Aceito | Alta |
| 010 | Error Wrapping | Aceito | Alta |
| 011 | Structured Logging | Proposto | Baixa |
| 012 | YAML Config | Proposto | Baixa |
| 013 | Catálogo remoto | Aceito | Alta |

## Próximas decisões (ideias)

1. **ADR-014**: Sistema de plugins
2. **ADR-015**: Migração/export de sistema
3. **ADR-016**: Testes de integração (expansão)
4. **ADR-017**: Pipeline CI/CD

## Manutenção

Adicionar um novo ficheiro `NNN-kebab-title.md` quando tomar uma decisão arquitetural importante e atualizar a tabela acima.

## Template

```markdown
# ADR-XXX: [Título]

**Data**: YYYY-MM-DD
**Status**: [Proposto/Aceito/Rejeitado/Obsoleto]

## Contexto

[Problema/situação]

## Decisão

[O que decidimos]

## Razões

[Por que decidimos assim]

## Alternativas consideradas

[O que mais consideramos]

## Consequências

[Impactos da decisão]
```
