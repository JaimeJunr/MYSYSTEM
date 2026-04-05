# Getting Started

This guide walks you through building and running Homestead for the first time.

## 1. Install Go

### Ubuntu/Debian

```bash
sudo apt update && sudo apt install golang-go
go version   # should print go1.21.x or higher
```

### Install the latest version manually

```bash
# Remove old version if needed
sudo apt remove golang-go

# Download and extract
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$HOME/go/bin
```

## 2. Clone and build

```bash
git clone https://github.com/JaimeJunr/Homestead
cd Homestead
make install
```

This compiles the binary and installs it to `$GOPATH/bin`. After that, `homestead` is available from anywhere in your terminal.

If you just want to try it without installing:

```bash
make run
```

## 3. Zsh: plugins vs repositório

O menu tem duas entradas para Zsh:

- **🔧 Plugins e temas Zsh** (visível depois de instalar Oh My Zsh) — wizard local: escolher plugins, ferramentas (NVM, Bun, etc.) e gerar `.zshrc`. Nada de repositório.
- **⚙️ Configurar Zsh** — repositório de config (dotfiles): criar um novo repo e enviar para a nuvem (GitHub, etc.) ou restaurar/migrar a partir de um repo existente. Útil para migrar de uma máquina para outra: na máquina 1 guarda no repo, na máquina 2 escolhe “já tenho repositório” e restaura.

**Plugins e temas Zsh** (wizard local):

1. **Plugins** — escolher plugins Zsh (git, docker, autosuggestions, etc.)
2. **Ferramentas** — NVM, Bun, SDKMAN, pnpm, etc.
3. **Revisão** — pré-visualizar e aplicar

Navegação: `↑/↓`, `Space` para marcar, `n`/`→` próximo passo, `Esc` voltar. Só escreve em disco quando confirmar na revisão.

**Configurar Zsh** (repo): escolher “Sim” (já tenho repo) ou “Não” (criar novo), indicar URL do repositório; o Homestead faz clone/push e backup/restauro dos ficheiros (por defeito `.zshrc` e `~/.zsh/`).

## 4. Run system scripts

From the main menu:

- **🧹 Limpeza do Sistema** — cleanup scripts (Docker, npm, apt caches, large files)
- **📊 Monitoramento** — native panels (no bash): battery, memory/swap, disk by mount, CPU load, network counters and throughput, temperature (sysfs), failing **systemd --user** units (auto-refresh ~3s)
- **⚙️ Configurações** — installer catalog URL, light/dark theme, script root, dotfiles path, confirmation toggles (`~/.config/homestead/preferences.yaml`; `HOMESTEAD_CATALOG_URL` overrides the saved URL)

For **Monitoramento**, pick an item and confirm to open the panel. For bash scripts under **Limpeza**, confirm runs the script with live output in the TUI.

## 5. Development setup

```bash
make test              # run all tests
make test-coverage     # tests + coverage report
make test-verbose      # verbose output
make build             # build binary (./homestead)
make clean             # remove build artifacts
```

## Troubleshooting

**`make: command not found`**

```bash
sudo apt install make
```

**`go: command not found`**

Make sure `/usr/local/go/bin` is in your PATH. Add it to `~/.bashrc` or `~/.zshrc` and restart your terminal.

**Scripts aren't running**

The scripts in `scripts/` need to be executable:

```bash
chmod +x scripts/cleanup/*.sh scripts/monitoring/*.sh
```

**Build fails with missing module**

```bash
make tidy   # updates go.sum
make build
```

---

See [README.md](README.md) for a full feature overview. Documentation hub: [docs/INDEX.md](docs/INDEX.md) (product context, architecture, testing).
