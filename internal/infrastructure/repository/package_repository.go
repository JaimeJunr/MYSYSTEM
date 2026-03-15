package repository

import (
	"sync"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// InMemoryPackageRepository is an in-memory implementation of PackageRepository
type InMemoryPackageRepository struct {
	packages map[string]*entities.Package
	mu       sync.RWMutex
}

// NewInMemoryPackageRepository creates a new in-memory package repository
func NewInMemoryPackageRepository() interfaces.PackageRepository {
	repo := &InMemoryPackageRepository{
		packages: make(map[string]*entities.Package),
	}
	repo.initializeDefaultPackages()
	return repo
}

// initializeDefaultPackages initializes the repository with default packages
func (r *InMemoryPackageRepository) initializeDefaultPackages() {
	defaultPackages := []entities.Package{
		// IDEs
		{
			ID:          "claude-code",
			Name:        "Claude Code CLI",
			Description: "CLI oficial da Anthropic para desenvolvimento com Claude",
			Version:     "latest",
			Category:    types.PackageCategoryIDE,
			DownloadURL: "https://storage.googleapis.com/claude-code/install.sh",
			InstallCmd:  "bash install.sh",
			CheckCmd:    "which claude-code",
		},
		{
			ID:          "cursor",
			Name:        "Cursor AI",
			Description: "Editor de código com IA integrada",
			Version:     "latest",
			Category:    types.PackageCategoryIDE,
			DownloadURL: "https://download.cursor.sh/linux/appImage/x64",
			InstallCmd:  "chmod +x cursor.AppImage && sudo mv cursor.AppImage /usr/local/bin/cursor",
			CheckCmd:    "which cursor",
		},
		{
			ID:          "antigravity",
			Name:        "Antigravity",
			Description: "IDE moderna com recursos avançados",
			Version:     "latest",
			Category:    types.PackageCategoryIDE,
			DownloadURL: "https://antigravity.dev/download/linux",
			InstallCmd:  "sudo dpkg -i antigravity.deb || sudo apt-get install -f -y",
			CheckCmd:    "which antigravity",
		},

		// Shell Core (Zsh, Oh My Zsh, Powerlevel10k) - install via "Instalar componentes core"
		{
			ID:          "zsh",
			Name:        "Zsh",
			Description: "Z Shell - shell poderoso e configurável",
			Version:     "latest",
			Category:    types.PackageCategoryZshCore,
			InstallCmd:  "sudo apt-get install -y zsh",
			CheckCmd:    "which zsh",
		},
		{
			ID:          "oh-my-zsh",
			Name:        "Oh My Zsh",
			Description: "Framework para gerenciar configuração Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryZshCore,
			DownloadURL: "https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh",
			InstallCmd:  "sh -c \"$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\" \"\" --unattended",
			CheckCmd:    "test -d ~/.oh-my-zsh",
		},
		{
			ID:          "powerlevel10k",
			Name:        "Powerlevel10k",
			Description: "Tema Zsh rápido e customizável",
			Version:     "latest",
			Category:    types.PackageCategoryZshCore,
			InstallCmd:  "git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k",
		},

		// Zsh Plugins - Built-in (5)
		{
			ID:          "zsh-plugin-git",
			Name:        "Git Plugin",
			Description: "Plugin built-in do Oh My Zsh para Git",
			Version:     "built-in",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "# Built-in plugin",
			CheckCmd:    "test -f ~/.oh-my-zsh/plugins/git/git.plugin.zsh",
		},
		{
			ID:          "zsh-plugin-docker",
			Name:        "Docker Plugin",
			Description: "Plugin built-in do Oh My Zsh para Docker",
			Version:     "built-in",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "# Built-in plugin",
			CheckCmd:    "test -f ~/.oh-my-zsh/plugins/docker/docker.plugin.zsh",
		},
		{
			ID:          "zsh-plugin-rails",
			Name:        "Rails Plugin",
			Description: "Plugin built-in do Oh My Zsh para Ruby on Rails",
			Version:     "built-in",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "# Built-in plugin",
			CheckCmd:    "test -f ~/.oh-my-zsh/plugins/rails/rails.plugin.zsh",
		},
		{
			ID:          "zsh-plugin-z",
			Name:        "Z Plugin",
			Description: "Plugin built-in para navegação rápida de diretórios",
			Version:     "built-in",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "# Built-in plugin",
			CheckCmd:    "test -f ~/.oh-my-zsh/plugins/z/z.plugin.zsh",
		},
		{
			ID:          "zsh-plugin-sudo",
			Name:        "Sudo Plugin",
			Description: "Plugin built-in para adicionar sudo facilmente",
			Version:     "built-in",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "# Built-in plugin",
			CheckCmd:    "test -f ~/.oh-my-zsh/plugins/sudo/sudo.plugin.zsh",
		},

		// Zsh Plugins - External (10)
		{
			ID:          "zsh-autosuggestions",
			Name:        "Zsh Autosuggestions",
			Description: "Sugestões automáticas baseadas no histórico",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
		},
		{
			ID:          "zsh-syntax-highlighting",
			Name:        "Zsh Syntax Highlighting",
			Description: "Destaque de sintaxe para comandos Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting",
		},
		{
			ID:          "fzf-zsh",
			Name:        "FZF Zsh Integration",
			Description: "Integração do FZF com Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/junegunn/fzf.git ~/.fzf && ~/.fzf/install --all",
			CheckCmd:    "test -d ~/.fzf",
		},
		{
			ID:          "you-should-use",
			Name:        "You Should Use",
			Description: "Lembra aliases existentes ao digitar comandos",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/MichaelAquilina/zsh-you-should-use.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/you-should-use",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/you-should-use",
		},
		{
			ID:          "zsh-completions",
			Name:        "Zsh Completions",
			Description: "Completions adicionais para Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/zsh-users/zsh-completions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-completions",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-completions",
		},
		{
			ID:          "zsh-history-substring-search",
			Name:        "Zsh History Substring Search",
			Description: "Busca no histórico por substring",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/zsh-users/zsh-history-substring-search ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-history-substring-search",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-history-substring-search",
		},
		{
			ID:          "fast-syntax-highlighting",
			Name:        "Fast Syntax Highlighting",
			Description: "Syntax highlighting mais rápido",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/zdharma-continuum/fast-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/fast-syntax-highlighting",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/fast-syntax-highlighting",
		},
		{
			ID:          "zsh-autocomplete",
			Name:        "Zsh Autocomplete",
			Description: "Autocomplete em tempo real para Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone --depth 1 -- https://github.com/marlonrichert/zsh-autocomplete.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autocomplete",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autocomplete",
		},
		{
			ID:          "auto-notify",
			Name:        "Auto Notify",
			Description: "Notificações automáticas para comandos longos",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/MichaelAquilina/zsh-auto-notify.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/auto-notify",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/auto-notify",
		},
		{
			ID:          "zsh-vi-mode",
			Name:        "Zsh Vi Mode",
			Description: "Melhor modo Vi para Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/jeffreytse/zsh-vi-mode ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-vi-mode",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-vi-mode",
		},

		// Development Tools (8)
		{
			ID:          "nvm",
			Name:        "NVM (Node Version Manager)",
			Description: "Gerenciador de versões Node.js",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash",
			CheckCmd:    "test -d ~/.nvm",
		},
		{
			ID:          "bun",
			Name:        "Bun",
			Description: "Runtime JavaScript/TypeScript rápido",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -fsSL https://bun.sh/install | bash",
			CheckCmd:    "test -d ~/.bun",
		},
		{
			ID:          "sdkman",
			Name:        "SDKMAN!",
			Description: "Gerenciador de SDKs para JVM",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -s https://get.sdkman.io | bash",
			CheckCmd:    "test -d ~/.sdkman",
		},
		{
			ID:          "pnpm",
			Name:        "pnpm",
			Description: "Gerenciador de pacotes Node.js eficiente",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -fsSL https://get.pnpm.io/install.sh | sh -",
			CheckCmd:    "which pnpm",
		},
		{
			ID:          "deno",
			Name:        "Deno",
			Description: "Runtime seguro para JavaScript e TypeScript",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -fsSL https://deno.land/install.sh | sh",
			CheckCmd:    "which deno",
		},
		{
			ID:          "angular-cli",
			Name:        "Angular CLI",
			Description: "Interface de linha de comando para Angular",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "npm install -g @angular/cli",
			CheckCmd:    "which ng",
		},
		{
			ID:          "openvpn3",
			Name:        "OpenVPN 3",
			Description: "Cliente VPN moderno",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "sudo apt-get install -y openvpn3",
			CheckCmd:    "which openvpn3",
		},
		{
			ID:          "homebrew",
			Name:        "Homebrew",
			Description: "Gerenciador de pacotes para Linux",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"",
			CheckCmd:    "which brew",
		},
	}

	for _, pkg := range defaultPackages {
		pkgCopy := pkg
		r.packages[pkg.ID] = &pkgCopy
	}
}

// FindAll returns all packages
func (r *InMemoryPackageRepository) FindAll() ([]entities.Package, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	packages := make([]entities.Package, 0, len(r.packages))
	for _, pkg := range r.packages {
		packages = append(packages, *pkg)
	}

	return packages, nil
}

// FindByID finds a package by ID
func (r *InMemoryPackageRepository) FindByID(id string) (*entities.Package, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pkg, exists := r.packages[id]
	if !exists {
		return nil, types.ErrNotFound
	}

	pkgCopy := *pkg
	return &pkgCopy, nil
}

// FindByCategory finds packages by category
func (r *InMemoryPackageRepository) FindByCategory(category types.PackageCategory) ([]entities.Package, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	packages := make([]entities.Package, 0)
	for _, pkg := range r.packages {
		if pkg.Category == category {
			packages = append(packages, *pkg)
		}
	}

	return packages, nil
}

// Save saves a package
func (r *InMemoryPackageRepository) Save(pkg *entities.Package) error {
	if err := pkg.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	pkgCopy := *pkg
	r.packages[pkg.ID] = &pkgCopy

	return nil
}

// Delete deletes a package
func (r *InMemoryPackageRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.packages[id]; !exists {
		return types.ErrNotFound
	}

	delete(r.packages, id)
	return nil
}

// Exists checks if a package exists
func (r *InMemoryPackageRepository) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.packages[id]
	return exists
}
