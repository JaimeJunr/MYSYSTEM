# Produto e contexto

Documento para alinhar **o que é o Homestead**, **para quem**, e **que problemas resolve** — antes de mergulhar em arquitetura ou código.

## O produto

Homestead é uma **ferramenta de linha de comando com TUI** (interface rica no terminal) para **Linux**, focada em:

- **Manutenção do sistema** — scripts de limpeza (caches, Docker, espaço em disco, etc.).
- **Monitorização simples** — painéis nativos no TUI (Go, atualização periódica): bateria, RAM/swap, disco por mount, carga de CPU, rede, temperatura (sysfs), unidades **systemd --user** em falha.
- **Instalação de ferramentas** — catálogo embutido + manifesto JSON remoto opcional; preferências (`preferences.yaml`) e variável `HOMESTEAD_CATALOG_URL` definem a URL efetiva; estratégias distintas por pacote (script, pacote, URL).
- **Preferências no TUI** — tema claro/escuro, URL do catálogo, raiz dos scripts, repo de dotfiles e confirmações antes de executar.
- **Configuração de shell** — fluxos Zsh (wizard de plugins/temas com Oh My Zsh; backup/restauro de dotfiles via Git).

O binário é o ponto único de entrada: o utilizador navega por menus e confirma ações; a aplicação orquestra scripts bash, comandos do sistema e ficheiros de configuração.

## Contexto de negócio (no sentido de “porque existe”)

### Problema

Quem desenvolve ou administra Linux acumula tarefas repetitivas: limpar caches, instalar o mesmo conjunto de ferramentas em máquinas novas, alinhar `.zshrc` e dotfiles entre computadores. Essas tarefas estão espalhadas por scripts soltos, documentação pessoal e memória.

### Proposta de valor

- **Um só lugar** para tarefas de “casa arrumada” e setup, sem memorizar caminhos de scripts.
- **Curadoria** — lista de pacotes/scripts mantida no repositório, em vez de cada utilizador reinventar.
- **Previsibilidade** — camadas e testes reduzem regressões quando se adiciona um script ou um instalador.
- **Migração** — fluxo explícito para dotfiles/Zsh entre máquinas (Git).

### Utilizadores-alvo

- Desenvolvedores em **Linux** (Ubuntu/Debian como referência).
- Quem já está confortável com terminal e, em alguns fluxos, com `sudo` ou credenciais Git remotas.

### Fora de escopo (explícito)

- Não é um gestor de pacotes global do SO (substitui parcialmente “coleções de comandos”, não o `apt` inteiro).
- Não é uma aplicação web nem API remota obrigatória: execução **local**, opcionalmente com catálogo remoto de instaladores conforme ADRs.
- Não promete suporte a todos os ambientes Linux: distros e permissões variam; o código assume convenções documentadas no repositório.

## Relação com a documentação técnica

| Necessidade              | Onde ir                                                          |
| ------------------------ | ---------------------------------------------------------------- |
| Instalar e correr        | [README.md](../../README.md), [GETTING_STARTED.md](../../GETTING_STARTED.md) |
| Camadas e padrões        | [../architecture/VERTICAL_PATTERNS.md](../architecture/VERTICAL_PATTERNS.md), [../architecture/ARCHITECTURE.md](../architecture/ARCHITECTURE.md) |
| Decisões passo a passo   | [../architecture/adrs/README.md](../architecture/adrs/README.md) |

---

**Última atualização:** 2026-04-04
