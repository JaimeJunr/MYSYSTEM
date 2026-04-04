# ADR-007: Command pattern para undo/redo

**Data**: 2026-03-14  
**Status**: Proposto  

## Contexto

Permitir reverter instalações.

## Decisão

Usar **command pattern** para operações reversíveis:

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

## Razões

- Histórico de operações.
- Undo inteligente (desinstalar só o que foi instalado).
- Auditoria de ações.
- Operações em batch.

## Implementação (fases futuras)

- **Fase 1**: histórico básico.
- **Fase 2**: undo/redo completo.
- **Fase 3**: transações com rollback em erro.
