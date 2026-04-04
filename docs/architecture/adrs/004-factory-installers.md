# ADR-004: Factory pattern para installers

**Data**: 2026-03-14  
**Status**: Aceito  

## Contexto

Criar instaladores de diferentes tipos (apt, snap, manual).

## Decisão

Usar **factory pattern**:

```go
type InstallerFactory interface {
    Create(packageType string) (Installer, error)
}
```

## Razões

- Centraliza a criação de installers.
- Extensível: novos tipos com pouco impacto.
- Type-safe: retorna uma interface comum.
- Configurável: pode usar config para customizar.

## Quando adicionar um novo installer

1. Criar struct que implementa a interface `Installer`.
2. Adicionar o case no factory.
3. Código existente que consome a interface não precisa de alterações em massa.
