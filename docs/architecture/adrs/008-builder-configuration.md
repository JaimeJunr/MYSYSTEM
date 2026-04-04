# ADR-008: Builder pattern para configurações

**Data**: 2026-03-14  
**Status**: Aceito  

## Contexto

Installers precisam de configurações complexas.

## Decisão

Usar **builder pattern**:

```go
config := NewInstallerConfigBuilder().
    WithPackageName("docker").
    WithVersion("latest").
    AddDependency("ca-certificates").
    EnableAutoStart().
    Build()
```

## Razões

- Interface fluente e legível.
- Validação no `Build()`.
- Defaults automáticos.
- Imutabilidade após `Build()`.

## Uso

- Wizards na TUI.
- Configs YAML.
- Presets (via Director, quando aplicável).
