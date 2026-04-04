# ADR-010: Erros com wrapping

**Data**: 2026-03-14  
**Status**: Aceito  

## Contexto

Como lidar com erros através das camadas.

## Decisão

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

## Razões

- Contexto preservado na cadeia de erros.
- Suporte a `errors.Is` / `errors.As`.
- Debugging mais simples.

## Convenção

- Envolver com contexto útil em cada limite.
- Usar `%w`, não `%v`, quando o erro deve participar da cadeia.
- Evitar logar e devolver o mesmo erro sem critério (escolher uma estratégia por camada).
