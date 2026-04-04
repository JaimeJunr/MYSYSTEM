# ADR-001: Arquitetura em camadas com Clean Architecture

**Data**: 2026-03-14  
**Status**: Aceito  

## Contexto

Necessidade de definir estrutura base do projeto.

## Decisão

Adotar **arquitetura em camadas** (Layered Architecture) com princípios de **Clean Architecture**:

1. **Presentation Layer** (TUI) — `internal/tui/`
2. **Application Layer** (use cases) — `internal/app/`
3. **Domain Layer** (entidades + interfaces) — `internal/domain/`
4. **Infrastructure Layer** (implementações) — `internal/infrastructure/`

## Razões

- **Separação de responsabilidades**: cada camada tem propósito único; é fácil saber onde adicionar código.
- **Testabilidade**: domínio isolado de frameworks; mock de dependências simples; testes unitários focados.
- **Manutenibilidade**: mudanças localizadas; não rebenta outras camadas por acoplamento acidental.
- **Extensibilidade**: novas funcionalidades sem reescrever a base; trocar implementações (ex.: apt → snap).

## Alternativas consideradas

- **Monólito sem camadas**: rápido no início, difícil de manter quando cresce, testes mais complexos.
- **Microserviços**: over-engineering para um CLI, complexidade e overhead de comunicação desnecessários.

## Consequências

**Positivas**: código organizado e previsível; onboarding mais simples; preparado para crescimento.

**Negativas**: mais boilerplate inicial; exige disciplina; curva de aprendizado.
