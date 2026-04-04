# ADR-012: Configuração via YAML

**Data**: 2026-03-14  
**Status**: Proposto (futuro)  

## Contexto

Como utilizadores customizam o Homestead.

## Decisão

Configuração em `~/.config/homestead/config.yaml`:

```yaml
installer:
  preferred_method: apt
  auto_update: true
  backup_before_install: true

scripts:
  custom_directory: ~/my-scripts

ui:
  theme: dark
  confirm_destructive: true
```

## Razões

- YAML legível para humanos.
- Versionável em dotfiles.
- Defaults em código para funcionar sem ficheiro.

## Implementação (direção)

1. Biblioteca de parsing (ex.: Viper) se adotado.
2. Defaults no código.
3. Overrides via variáveis de ambiente quando fizer sentido.
