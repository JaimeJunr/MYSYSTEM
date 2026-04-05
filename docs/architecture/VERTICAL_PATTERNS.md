# Padrões do sistema (de cima para baixo)

Este documento define **a pilha de padrões** do Homestead na ordem em que o pedido atravessa o sistema: da interação no terminal até a execução no SO. Use-o como mapa mental antes do [ARCHITECTURE.md](ARCHITECTURE.md) (visão ampla) ou do [PATTERNS_GUIDE.md](PATTERNS_GUIDE.md) (exemplos longos de código).

## 1. Visão em quatro níveis

```
Utilizador (teclado)
        │
        ▼
┌───────────────────────────────────────┐
│ 1. Apresentação — Bubble Tea (TUI)    │  Model / Update / View; msg; cmds; items
├───────────────────────────────────────┤
│ 2. Aplicação — serviços               │  Orquestração, casos de uso finos
├───────────────────────────────────────┤
│ 3. Domínio — entidades e portas       │  Regras e contratos (interfaces)
├───────────────────────────────────────┤
│ 4. Infraestrutura — adaptadores       │  Bash, ficheiros, git, catálogo, instaladores
└───────────────────────────────────────┘
        │
        ▼
Sistema operativo (bash, apt, ficheiros, …)
```

**Regra de dependência:** camadas externas dependem das internas; `internal/domain` não importa `tui` nem `infrastructure`.

---

## 2. Nível 1 — Apresentação (`internal/tui`)

| Padrão / ideia | Papel no Homestead |
| -------------- | ------------------ |
| **ELM / M.V.U.** | `Model` + `Update` + `View`: estado imutável, mensagens como reducers. |
| **Message taxonomy** | Tipos em `internal/tui/msg` — eventos de UI, resultados assíncronos, erros. |
| **Command factory** | `internal/tui/cmds` — `tea.Cmd` que chamam serviços ou IO sem bloquear o loop. |
| **List item mapping** | `internal/tui/items` — cada linha de menu/lista como `list.Item` com título/descr. |
| **Theming** | `internal/tui/theme` — estilos Lipgloss partilhados. |
| **Separation of concerns** | `view_render.go`, `lists.go`, `menu.go` concentram renderização e navegação; o `model` coordena. |

O TUI **não** implementa acesso direto ao disco ou execução de bash como regra: delega em serviços da camada 2 (interfaces definidas no domínio).

*Detalhe de ficheiros:* [TUI_LAYOUT.md](TUI_LAYOUT.md).

---

## 3. Nível 2 — Aplicação (`internal/app/services`)

| Padrão / ideia | Papel no Homestead |
| -------------- | ------------------ |
| **Application service** | `ScriptService`, `InstallerService`, serviços de config/wizard — um método por fluxo de utilizador relevante. |
| **Orquestração** | Valida pré-condições, pede entidades aos repositórios, chama executor/instalador, compõe erros. |
| **Sem regra de negócio pesada** | Regras estáveis e tipos vivem no domínio; o serviço **coordena**. |

A composição (quem implementa cada interface) é feita em `cmd/homestead/main.go` (ver nível transversal).

---

## 4. Nível 3 — Domínio (`internal/domain`)

| Padrão / ideia | Papel no Homestead |
| -------------- | ------------------ |
| **Entidades** | `Script`, `Package`, categorias, tipos de instalação — dados e invariantes do problema. |
| **Ports (interfaces)** | `ScriptRepository`, `ScriptExecutor`, contratos de instalador — **o que** o sistema precisa, não **como**. |
| **Tipos e erros de domínio** | `internal/domain/types` — categorias, erros sentinelas quando aplicável. |

Nada aqui importa Bubble Tea, YAML concreto ou `exec.Command`.

---

## 5. Nível 4 — Infraestrutura (`internal/infrastructure`)

| Padrão / ideia | Papel no Homestead |
| -------------- | ------------------ |
| **Repository (concreto)** | Implementações em memória ou ficheiro dos catálogos de scripts/pacotes. |
| **Executor** | `BashExecutor` (e variantes) — adapta `Script` do domínio para processo bash com política sudo/TTY. |
| **Strategy** | Estratégias de instalação por tipo de pacote / origem (ver `installer` e ADRs de instalação). |
| **Adapter** | Catálogo JSON embutido ou remoto, templates, preferências em disco — detalhes de formato ficam aqui. |
| **Catalog / parsing** | Leitura e validação de metadados de instaladores fora do domínio puro quando é IO. |

---

## 6. Transversal — composição e limites

| Tema | Onde |
| ---- | ---- |
| **Dependency injection manual** | `cmd/homestead/main.go` — instancia repositórios, executor, installer, passa referências ao `tui.NewModel`. |
| **Observabilidade de progresso** | Mensagens para o TUI a partir de operações longas (alinhado a ADRs de progresso). |
| **Testes** | Domínio e serviços com mocks das interfaces; integração no topo (`integration_test.go`). |

---

## 7. Leitura sugerida

1. Este ficheiro (ordem vertical).
2. [ARCHITECTURE.md](ARCHITECTURE.md) — princípios SOLID, diagramas de camada, convenções.
3. [PATTERNS_GUIDE.md](PATTERNS_GUIDE.md) — código de exemplo por padrão.
4. [adrs/README.md](adrs/README.md) — decisões que fixam variantes (ex.: catálogo remoto, estratégias).

---

**Última atualização:** 2026-04-04
