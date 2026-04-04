# ADR-002: Uso de interfaces para inversão de dependência

**Data**: 2026-03-14  
**Status**: Aceito  

## Contexto

Como desacoplar camadas e facilitar testes.

## Decisão

Definir **interfaces no domain** e implementar na infraestrutura:

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

## Razões

- **Testabilidade**: mock de interfaces em testes.
- **Flexibilidade**: trocar implementações sem alterar o serviço.
- **SOLID**: princípio de inversão de dependências.
- **Independência**: o domínio não depende de infraestrutura.

## Consequências

É preciso definir interfaces antes ou em paralelo às implementações; mais ficheiros (interface + implementação), com ganho claro em flexibilidade.
