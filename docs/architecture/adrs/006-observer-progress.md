# ADR-006: Observer pattern para progresso

**Data**: 2026-03-14  
**Status**: Aceito  

## Contexto

Notificar a TUI sobre progresso de instalações longas.

## Decisão

Usar **observer pattern**:

```go
type ProgressObserver interface {
    OnProgress(current, total int, msg string)
    OnComplete(success bool, msg string)
}
```

## Razões

- **Desacoplamento**: o installer não conhece a TUI.
- **Múltiplos observers**: TUI, logger, métricas.
- **Atualizações em tempo real**: barra de progresso responsiva.

## Observers (exemplos)

1. **TUIObserver** — atualiza a barra de progresso.
2. **LogObserver** — escreve em ficheiro.
3. **MetricsObserver** — analytics (futuro).
