# Layout do pacote TUI (`internal/tui`)

A camada de apresentação usa [Bubble Tea](https://github.com/charmbracelet/bubbletea) com o pacote raiz `**tui**` como ponto de entrada (`NewModel`, `Model`, wizards). Subpacotes isolam tipos e comandos assíncronos **sem ciclos de import** com o raiz.

## Pacote raiz (`package tui`)

Ficheiros principais:


| Ficheiro                                    | Conteúdo                                                              |
| ------------------------------------------- | --------------------------------------------------------------------- |
| `model.go`                                  | `Model`, `NewModel`, `Init`, `Update`, `handleEnter`                  |
| `view_state.go`                             | `ViewState`, constantes de ações do menu                              |
| `menu.go`                                   | Itens do menu principal (`getMainMenuItems`)                          |
| `lists.go`                                  | Carregar listas (scripts, pacotes, categorias de instaladores)        |
| `view_render.go`                            | `View()` e ecrãs (confirmação, progresso, saída de script, Zsh apply) |
| `native_monitor.go`                         | Painéis nativos (bateria, RAM, disco, load, rede, térmico, systemd user) |
| `settings_model.go`                         | Ecrã **Configurações** (YAML, catálogo, tema, caminhos, confirmações)   |
| `zsh_wizard_model.go`                       | Wizard “Plugins e temas Zsh”                                          |
| `zsh_repo_model.go`                         | Wizard “Configurar Zsh” (repositório)                                 |
| `model_test.go`, `zsh_wizard_model_test.go` | Testes                                                                |


**Import público** (inalterado para `cmd/homestead` e `integration_test.go`):

```go
import "github.com/JaimeJunr/Homestead/internal/tui"

model := tui.NewModel(scriptService, installerService, configService, repoService, catalogURL, prefs, prefsPath, catalogEnvSet)
```

## Subpacotes


| Pacote   | Caminho                | Responsabilidade                                                                                    |
| -------- | ---------------------- | --------------------------------------------------------------------------------------------------- |
| `cmds`   | `internal/tui/cmds/`   | Fábricas `tea.Cmd`: catálogo remoto, instalação, captura de script, deteção Oh My Zsh, URLs         |
| `items`  | `internal/tui/items/`  | Implementações `list.Item` (menu, script, pacote, grupo de instalador)                              |
| `msg`    | `internal/tui/msg/`    | Tipos de mensagens Bubble Tea (`Progress`, `CatalogFetched`, `ScriptCaptured`, …)                   |
| `sysurl` | `internal/tui/sysurl/` | Abrir URL no SO, clipboard, `PackageKeyboardURL` — separado para `cmds` não importar o pacote `tui` |
| `theme`  | `internal/tui/theme/`  | Estilos Lipgloss partilhados, `StripANSI`, títulos de secções de instaladores                       |


## Fluxo de dependências

```
cmd/homestead  →  tui (raiz)
                    ├── app/services, domain, infrastructure (catalog), monitoring
                    ├── tui/cmds  →  tui/msg, tui/sysurl, services, catalog
                    ├── tui/items →  domain (entities, types)
                    ├── tui/theme →  domain/types (títulos de categoria)
                    └── tui/msg   →  domain/interfaces, monitoring (snapshots)
```

O raiz **não** é importado por `cmds`, `sysurl`, `items`, `msg` ou `theme`, evitando ciclos.

## Onde adicionar código novo

- **Novo tipo de mensagem assíncrona** → `internal/tui/msg` e `case` em `model.Update`.
- **Novo `tea.Cmd` que chama serviços ou rede** → `internal/tui/cmds`.
- **Nova linha de lista reutilizável** → `internal/tui/items`.
- **Cores / estilos globais da TUI** → `internal/tui/theme`.
- **Novo ecrã ou ramo grande de `Update`** → preferir ficheiro no raiz `tui` (ex. `view_render.go`, `lists.go`) ou extrair funções auxiliares no mesmo pacote.

---

**Última atualização:** 2026-04-04