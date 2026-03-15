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

## 3. Run the Zsh wizard

From the main menu, select **🐚 Configurar Zsh**. The wizard walks you through:

1. **Core components** — choose Zsh, Oh My Zsh, and/or Powerlevel10k
2. **Plugins** — pick from 15 Zsh plugins (git, docker, autosuggestions, etc.)
3. **Dev tools** — NVM, Bun, SDKMAN, pnpm, and others
4. **Project config** — optionally include project-specific aliases/functions
5. **Review** — preview your future `.zshrc` before applying anything

Navigate with `↑/↓`, toggle items with `Space`, move between steps with `n` or `→`, go back with `Esc`.

Nothing is written to disk until you confirm at the review step.

## 4. Run system scripts

From the main menu:

- **🧹 Limpeza do Sistema** — cleanup scripts (Docker, npm, apt caches, large files)
- **📊 Monitoramento** — battery and memory info

Select a script, confirm, and it runs with live output in your terminal.

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

See [README.md](README.md) for a full feature overview, and [docs/](docs/) for architecture documentation.
