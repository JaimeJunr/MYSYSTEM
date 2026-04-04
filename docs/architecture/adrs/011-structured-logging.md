# ADR-011: Logs estruturados

**Data**: 2026-03-14  
**Status**: Proposto (futuro)  

## Contexto

Como fazer logging de forma consistente.

## Decisão

Usar **structured logging** (ex.: Charm Log):

```go
logger.Info("installing package",
    "package", pkgName,
    "version", version,
    "method", "apt",
)
```

## Razões

- Parseável para análise automatizada.
- Formato consistente.
- Filtragem por campos.

## Níveis sugeridos

- **Debug**: detalhes internos.
- **Info**: operações importantes.
- **Warn**: problemas não críticos.
- **Error**: erros que exigem atenção.
