# ADR-005: Strategy pattern para métodos de instalação

**Data**: 2026-03-14  
**Status**: Aceito  

## Contexto

O mesmo package pode ser instalado de múltiplas formas.

## Decisão

Usar **strategy pattern**:

```go
type InstallStrategy interface {
    Install(pkg Package) error
    CanInstall() bool
}

// Strategies: AptStrategy, SnapStrategy, FromSourceStrategy
```

## Razões

- Algoritmos intercambiáveis.
- Fallback automático (tentar apt, depois snap, etc.).
- Extensível sem modificar código existente (Open/Closed).

## Exemplo

Docker pode ser instalado via:

1. Apt (preferido em Ubuntu).
2. Snap (fallback).
3. Script oficial (manual).

O serviço escolhe automaticamente a melhor strategy disponível.
