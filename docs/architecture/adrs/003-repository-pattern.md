# ADR-003: Repository pattern para acesso a dados

**Data**: 2026-03-14  
**Status**: Aceito  

## Contexto

Como gerir coleções de scripts e packages.

## Decisão

Usar **repository pattern**:

```go
type ScriptRepository interface {
    FindAll() ([]Script, error)
    FindByID(id string) (*Script, error)
    FindByCategory(cat Category) ([]Script, error)
}
```

Inicialmente: **repositório em memória**. Futuro: baseado em ficheiros ou SQLite se necessário.

## Razões

- Abstração da persistência.
- Facilita trocar o backend (memória → ficheiro → BD).
- Queries centralizadas.
- Testável com mock do repositório.

## Alternativas consideradas

- **Acesso direto**: acoplamento forte.
- **DAO**: mais complexo do que o necessário aqui.
- **ORM**: overhead para um CLI simples.
