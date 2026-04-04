# ADR-009: Dependency injection manual

**Data**: 2026-03-14  
**Status**: Aceito  

## Contexto

Como ligar dependências entre camadas.

## Decisão

**Dependency injection manual** via construtores:

```go
// main.go
executor := executor.NewBashExecutor()
repo := repository.NewInMemoryScriptRepository()
service := services.NewScriptService(repo, executor)
model := tui.NewModel(service)
```

## Razões

- Simplicidade: sem frameworks de DI.
- Explícito: fácil de seguir no código.
- Type-safe: o compilador valida.
- Testável: injetar mocks nos testes.

## Alternativas consideradas

- **Wire/Dig**: over-engineering para o tamanho atual.
- **Service locator**: anti-pattern neste contexto.
- **Variáveis globais**: dificultam testes e razão sobre o fluxo.

## Consequências

O wiring concentra-se em `main.go`; se crescer muito, reconsiderar ferramentas como Wire.
