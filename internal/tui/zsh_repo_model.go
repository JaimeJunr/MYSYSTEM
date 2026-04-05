package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/JaimeJunr/Homestead/internal/app/services"
	btmsg "github.com/JaimeJunr/Homestead/internal/tui/msg"
	"github.com/JaimeJunr/Homestead/internal/tui/sysurl"
)

// zshRepoResultMsg is sent when a repo operation (push/clone/restore) finishes
type zshRepoResultMsg struct {
	Err error
}

// ZshRepoView represents each step of the Configurar Zsh (repo) wizard
type ZshRepoView int

const (
	ZshRepoViewAlreadyConfigured ZshRepoView = iota // Já tem repo configurado: backup, restaurar, mudar
	ZshRepoViewChoice                               // Já tem repo? Sim / Não
	ZshRepoViewNewRepoMethod                        // Criar novo: [g] GitHub automático ou [u] colar URL
	ZshRepoViewNewRepo                              // Colar URL do repo (quando não usa gh)
	ZshRepoViewNewRepoGhName                        // Nome do repo no GitHub (quando usa gh)
	ZshRepoViewExistingRepo                         // Repo existente: URL, clone, aplicar
	ZshRepoViewRunning                              // Async operation in progress
	ZshRepoViewSuccess
	ZshRepoViewError
)

// ZshRepoModel is the Bubbletea model for the Configurar Zsh (repo) flow
type ZshRepoModel struct {
	repoService   *services.RepoService
	configService *services.ConfigService

	currentView ZshRepoView
	width       int
	height      int

	urlInput     textinput.Model
	repoNameInput textinput.Model

	// User choice: true = already has repo, false = create new
	hasExistingRepo *bool
	// true = create via gh, false = paste URL
	useGhForNew *bool

	// Result
	lastError error

	// Quando true, após success/error o Esc volta ao painel "já configurado" em vez de sair
	returnToDashboard bool

	// Done/cancelled (cancel returns to menu; success stays on success view until Esc)
	done      bool
	cancelled bool
}

// NewZshRepoModel creates a new repo wizard model
func (m *ZshRepoModel) isTextInputView() bool {
	switch m.currentView {
	case ZshRepoViewNewRepo, ZshRepoViewExistingRepo, ZshRepoViewNewRepoGhName:
		return true
	default:
		return false
	}
}

func NewZshRepoModel(repoService *services.RepoService, configService *services.ConfigService) ZshRepoModel {
	ti := textinput.New()
	ti.Placeholder = "https://github.com/user/dotfiles.git"
	ti.Width = 60

	nameTi := textinput.New()
	nameTi.Placeholder = "dotfiles"
	nameTi.Width = 40

	initialView := ZshRepoViewChoice
	if repoService != nil && repoService.IsRepo() {
		initialView = ZshRepoViewAlreadyConfigured
	}

	return ZshRepoModel{
		repoService:     repoService,
		configService:   configService,
		currentView:     initialView,
		hasExistingRepo: nil,
		useGhForNew:     nil,
		urlInput:        ti,
		repoNameInput:   nameTi,
	}
}

// Init initializes the model
func (m ZshRepoModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m ZshRepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case zshRepoResultMsg:
		if msg.Err != nil {
			m.lastError = msg.Err
			m.currentView = ZshRepoViewError
		} else {
			m.currentView = ZshRepoViewSuccess
		}
		return m, nil

	case tea.KeyMsg:
		keyStr := msg.String()
		isEnter := msg.Type == tea.KeyEnter || keyStr == "enter" || keyStr == "ctrl+m"

		// Running view: only allow Esc
		if m.currentView == ZshRepoViewRunning {
			if msg.String() == "esc" || msg.String() == "ctrl+c" {
				m.cancelled = true
				m.done = true
				return m, nil
			}
			return m, nil
		}
		// Success/Error: Esc or Enter — volta ao menu ou ao painel "já configurado"
		if m.currentView == ZshRepoViewSuccess || m.currentView == ZshRepoViewError {
			if keyStr == "o" || keyStr == "O" {
				if m.currentView == ZshRepoViewSuccess {
					web := m.browserOriginURL()
					if web != "" {
						return m, zshRepoOpenURLCmd(web)
					}
					return m, zshRepoOpenURLCmd("") // dispara toast de erro via OpenURL
				}
				return m, nil
			}
			if msg.String() == "esc" || isEnter {
				if m.returnToDashboard {
					m.returnToDashboard = false
					m.currentView = ZshRepoViewAlreadyConfigured
					m.lastError = nil
				} else {
					m.done = true
					return m, tea.Quit
				}
				return m, nil
			}
			return m, nil
		}

		// Enter/ctrl+m para enviar URL ou nome do repo (tratado antes do switch para não ser consumido pelo textinput)
		if isEnter {
			if m.currentView == ZshRepoViewNewRepoGhName {
				name := m.repoNameInput.Value()
				if name == "" {
					name = "dotfiles"
				}
				m.currentView = ZshRepoViewRunning
				return m, runCreateGitHubRepoCmd(m.repoService, name)
			}
			if m.currentView == ZshRepoViewNewRepo || m.currentView == ZshRepoViewExistingRepo {
				url := m.urlInput.Value()
				if url != "" {
					m.currentView = ZshRepoViewRunning
					if m.hasExistingRepo != nil && *m.hasExistingRepo {
						return m, runRestoreRepoCmd(m.repoService, m.configService, url)
					}
					return m, runPushNewRepoCmd(m.repoService, url)
				}
			}
		}

		switch keyStr {
		case "ctrl+c":
			m.cancelled = true
			m.done = true
			return m, nil
		case "esc":
			if m.currentView == ZshRepoViewAlreadyConfigured {
				m.cancelled = true
				m.done = true
				return m, nil
			}
			if m.currentView > ZshRepoViewChoice {
				switch m.currentView {
				case ZshRepoViewNewRepoMethod:
					m.currentView = ZshRepoViewChoice
					m.hasExistingRepo = nil
				case ZshRepoViewNewRepo, ZshRepoViewNewRepoGhName:
					m.currentView = ZshRepoViewNewRepoMethod
					m.useGhForNew = nil
					m.urlInput.SetValue("")
					m.urlInput.Blur()
					m.repoNameInput.SetValue("")
					m.repoNameInput.Blur()
				case ZshRepoViewExistingRepo:
					m.currentView = ZshRepoViewChoice
					m.hasExistingRepo = nil
					m.urlInput.SetValue("")
					m.urlInput.Blur()
				default:
					m.currentView--
				}
			} else {
				// Estava no Choice: se já tem repo (veio de "Mudar"), volta ao painel
				if m.repoService != nil && m.repoService.IsRepo() {
					m.currentView = ZshRepoViewAlreadyConfigured
				} else {
					m.cancelled = true
					m.done = true
				}
			}
			return m, nil
		case "b":
			if m.currentView == ZshRepoViewAlreadyConfigured {
				m.currentView = ZshRepoViewRunning
				m.returnToDashboard = true
				return m, runBackupOnlyCmd(m.repoService)
			}
			return m, nil
		case "r":
			if m.currentView == ZshRepoViewAlreadyConfigured {
				m.currentView = ZshRepoViewRunning
				m.returnToDashboard = true
				return m, runRestoreOnlyCmd(m.repoService, m.configService)
			}
			return m, nil
		case "m":
			if m.currentView == ZshRepoViewAlreadyConfigured {
				m.currentView = ZshRepoViewChoice
			}
			return m, nil
		case "q":
			if m.currentView == ZshRepoViewAlreadyConfigured {
				m.cancelled = true
				m.done = true
				return m, nil
			}
			return m, nil
		case "y", "s":
			if m.currentView == ZshRepoViewChoice {
				yes := true
				m.hasExistingRepo = &yes
				m.currentView = ZshRepoViewExistingRepo
				m.urlInput.Focus()
				m.urlInput.SetValue("")
			}
			return m, nil
		case "n":
			if m.currentView == ZshRepoViewChoice {
				no := false
				m.hasExistingRepo = &no
				m.currentView = ZshRepoViewNewRepoMethod
			}
			return m, nil
		case "g":
			if m.currentView == ZshRepoViewNewRepoMethod {
				yes := true
				m.useGhForNew = &yes
				m.currentView = ZshRepoViewNewRepoGhName
				m.repoNameInput.Focus()
				m.repoNameInput.SetValue("")
			}
			return m, nil
		case "u":
			if m.currentView == ZshRepoViewNewRepoMethod {
				no := false
				m.useGhForNew = &no
				m.currentView = ZshRepoViewNewRepo
				m.urlInput.Focus()
				m.urlInput.SetValue("")
			}
			return m, nil
		}

		// Forward to textinput when on URL step or repo name step (Enter já foi tratado acima)
		if !isEnter {
			if m.currentView == ZshRepoViewNewRepo || m.currentView == ZshRepoViewExistingRepo {
				var cmd tea.Cmd
				m.urlInput, cmd = m.urlInput.Update(msg)
				return m, cmd
			}
			if m.currentView == ZshRepoViewNewRepoGhName {
				var cmd tea.Cmd
				m.repoNameInput, cmd = m.repoNameInput.Update(msg)
				return m, cmd
			}
		}

		return m, nil
	}

	return m, nil
}

// runCreateGitHubRepoCmd creates repo locally, copies dotfiles, commits, then creates GitHub repo via gh and pushes
func runCreateGitHubRepoCmd(repoService *services.RepoService, repoName string) tea.Cmd {
	return func() tea.Msg {
		if repoService == nil {
			return zshRepoResultMsg{Err: fmt.Errorf("serviço de repositório indisponível")}
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return zshRepoResultMsg{Err: err}
		}
		if !repoService.IsRepo() {
			if err := repoService.InitRepo(); err != nil {
				return zshRepoResultMsg{Err: err}
			}
		}
		if err := repoService.CopyToRepo(homeDir, services.DefaultDotfilesPaths); err != nil {
			return zshRepoResultMsg{Err: err}
		}
		if err := repoService.CommitAll("homestead: backup dotfiles"); err != nil {
			return zshRepoResultMsg{Err: err}
		}
		// gh repo create <name> --private --source=. --remote=origin --push
		if err := services.CreateGitHubRepoWithGh(repoService.RepoDir(), repoName, true); err != nil {
			return zshRepoResultMsg{Err: err}
		}
		return zshRepoResultMsg{}
	}
}

// runPushNewRepoCmd runs init, copy, commit, add remote, push and sends zshRepoResultMsg
func runPushNewRepoCmd(repoService *services.RepoService, remoteURL string) tea.Cmd {
	return func() tea.Msg {
		if repoService == nil {
			return zshRepoResultMsg{Err: fmt.Errorf("serviço de repositório indisponível")}
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return zshRepoResultMsg{Err: err}
		}
		if !repoService.IsRepo() {
			if err := repoService.InitRepo(); err != nil {
				return zshRepoResultMsg{Err: err}
			}
		}
		if err := repoService.CopyToRepo(homeDir, services.DefaultDotfilesPaths); err != nil {
			return zshRepoResultMsg{Err: err}
		}
		if err := repoService.CommitAll("homestead: backup dotfiles"); err != nil {
			return zshRepoResultMsg{Err: err}
		}
		if !repoService.HasRemote("origin") {
			if err := repoService.AddRemote("origin", remoteURL); err != nil {
				return zshRepoResultMsg{Err: err}
			}
		}
		if err := repoService.Push("origin", "main"); err != nil {
			// Try "master" if "main" fails (older repos)
			if err2 := repoService.Push("origin", "master"); err2 != nil {
				return zshRepoResultMsg{Err: err}
			}
		}
		return zshRepoResultMsg{}
	}
}

// runBackupOnlyCmd copies from home, commits and pushes (repo already configured)
func runBackupOnlyCmd(repoService *services.RepoService) tea.Cmd {
	return func() tea.Msg {
		if repoService == nil {
			return zshRepoResultMsg{Err: fmt.Errorf("serviço de repositório indisponível")}
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return zshRepoResultMsg{Err: err}
		}
		if err := repoService.CopyToRepo(homeDir, services.DefaultDotfilesPaths); err != nil {
			return zshRepoResultMsg{Err: err}
		}
		if err := repoService.CommitAll("homestead: backup dotfiles"); err != nil {
			return zshRepoResultMsg{Err: err}
		}
		if err := repoService.Push("origin", "main"); err != nil {
			if err2 := repoService.Push("origin", "master"); err2 != nil {
				return zshRepoResultMsg{Err: err}
			}
		}
		return zshRepoResultMsg{}
	}
}

// runRestoreOnlyCmd pulls and restores to home (repo already exists)
func runRestoreOnlyCmd(repoService *services.RepoService, configService *services.ConfigService) tea.Cmd {
	return func() tea.Msg {
		if repoService == nil {
			return zshRepoResultMsg{Err: fmt.Errorf("serviço de repositório indisponível")}
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return zshRepoResultMsg{Err: err}
		}
		if err := repoService.Pull(); err != nil {
			return zshRepoResultMsg{Err: err}
		}
		if configService != nil {
			_ = configService.BackupCurrentConfig()
		}
		if err := repoService.RestoreToHome(homeDir, services.DefaultDotfilesPaths); err != nil {
			return zshRepoResultMsg{Err: err}
		}
		return zshRepoResultMsg{}
	}
}

// runRestoreRepoCmd clones (or pulls), backs up, restores and sends zshRepoResultMsg
func runRestoreRepoCmd(repoService *services.RepoService, configService *services.ConfigService, repoURL string) tea.Cmd {
	return func() tea.Msg {
		if repoService == nil {
			return zshRepoResultMsg{Err: fmt.Errorf("serviço de repositório indisponível")}
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return zshRepoResultMsg{Err: err}
		}
		if !repoService.IsRepo() {
			if err := repoService.Clone(repoURL); err != nil {
				return zshRepoResultMsg{Err: err}
			}
		} else {
			if err := repoService.Pull(); err != nil {
				return zshRepoResultMsg{Err: err}
			}
		}
		if configService != nil {
			_ = configService.BackupCurrentConfig()
		}
		if err := repoService.RestoreToHome(homeDir, services.DefaultDotfilesPaths); err != nil {
			return zshRepoResultMsg{Err: err}
		}
		return zshRepoResultMsg{}
	}
}

// View renders the UI
func (m ZshRepoModel) View() string {
	if m.width == 0 {
		m.width = 80
	}
	if m.height == 0 {
		m.height = 24
	}

	title := wizardTitleStyle.Render("⚙️  Configurar Zsh - Repositório")
	help := wizardHelpStyle.Render("esc: voltar • ctrl+c: cancelar")

	switch m.currentView {
	case ZshRepoViewAlreadyConfigured:
		remoteURL := ""
		if m.repoService != nil {
			remoteURL = m.repoService.GetRemoteURL("origin")
		}
		if remoteURL == "" {
			remoteURL = "(sem remote origin)"
		}
		body := lipgloss.NewStyle().Padding(1, 2).Render(
			"Repositório já configurado.\n\n"+
				"  [b] Fazer backup agora — copia .zshrc e ~/.zsh para o repo e envia (push)\n"+
				"  [r] Restaurar do remoto — puxa alterações e aplica no sistema (pull + aplicar)\n"+
				"  [m] Mudar repositório — usar outro repo ou criar um novo\n"+
				"  [q] Sair\n\n"+
				"Remote: "+remoteURL,
		)
		box := wizardPreviewStyle.Render(body)
		return lipgloss.JoinVertical(lipgloss.Left, title, box, help)
	case ZshRepoViewChoice:
		body := lipgloss.NewStyle().Padding(1, 2).Render(
			"Já tem um repositório de config (dotfiles)?\n\n" +
				"  [s] Sim — quero restaurar/migrar a partir de um repo\n" +
				"  [n] Não — quero criar um novo e enviar para a nuvem (GitHub, etc.)",
		)
		box := wizardPreviewStyle.Render(body)
		return lipgloss.JoinVertical(lipgloss.Left, title, box, help)
	case ZshRepoViewNewRepoMethod:
		body := lipgloss.NewStyle().Padding(1, 2).Render(
			"Como quer criar o repositório?\n\n" +
				"  [g] Criar automaticamente no GitHub (requer gh instalado e autenticado)\n" +
				"  [u] Já criei o repo — colar URL (GitHub, GitLab ou outro)",
		)
		box := wizardPreviewStyle.Render(body)
		return lipgloss.JoinVertical(lipgloss.Left, title, box, help)
	case ZshRepoViewNewRepoGhName:
		body := lipgloss.NewStyle().Padding(1, 2).Render(
			"Nome do repositório no GitHub (será criado como privado):\n\n" + m.repoNameInput.View(),
		)
		box := wizardPreviewStyle.Render(body)
		return lipgloss.JoinVertical(lipgloss.Left, title, box, help+" • enter: criar e enviar")
	case ZshRepoViewNewRepo:
		body := lipgloss.NewStyle().Padding(1, 2).Render(
			"URL do repositório remoto (crie um repo vazio no GitHub/GitLab e cole o URL):\n\n" + m.urlInput.View() +
				"\n\nRepos públicos e privados funcionam (use SSH ou token para autenticação).",
		)
		box := wizardPreviewStyle.Render(body)
		return lipgloss.JoinVertical(lipgloss.Left, title, box, help+" • enter: enviar")
	case ZshRepoViewExistingRepo:
		body := lipgloss.NewStyle().Padding(1, 2).Render(
			"URL do repositório existente:\n\n" + m.urlInput.View() +
				"\n\nRepos públicos e privados funcionam (use SSH ou token para autenticação).",
		)
		box := wizardPreviewStyle.Render(body)
		return lipgloss.JoinVertical(lipgloss.Left, title, box, help+" • enter: restaurar")
	case ZshRepoViewRunning:
		body := lipgloss.NewStyle().Padding(1, 2).Render("A aguardar... (clone/push em curso)")
		box := wizardPreviewStyle.Render(body)
		return lipgloss.JoinVertical(lipgloss.Left, title, box, help)
	case ZshRepoViewSuccess:
		backHint := "voltar ao menu"
		if m.returnToDashboard {
			backHint = "voltar ao painel"
		}
		body := lipgloss.NewStyle().Padding(1, 2).Render(
			"✅ Operação concluída com sucesso.\n\n"+
				"[o] Abrir o repositório no navegador\n\n"+
				"Pressione Enter ou Esc para "+backHint+".",
		)
		box := wizardPreviewStyle.Render(body)
		return lipgloss.JoinVertical(lipgloss.Left, title, box, help)
	case ZshRepoViewError:
		errStr := ""
		if m.lastError != nil {
			errStr = m.lastError.Error()
		}
		backHint := "voltar"
		if m.returnToDashboard {
			backHint = "voltar ao painel"
		}
		body := lipgloss.NewStyle().Padding(1, 2).Render(
			"❌ Erro:\n\n" + errStr + "\n\nPressione Enter ou Esc para " + backHint + ".",
		)
		box := wizardPreviewStyle.Render(body)
		return lipgloss.JoinVertical(lipgloss.Left, title, box, help)
	default:
		return fmt.Sprintf("%s\n\n%s", title, help)
	}
}

func (m ZshRepoModel) browserOriginURL() string {
	if m.repoService == nil {
		return ""
	}
	return gitRemoteToWebURL(m.repoService.GetRemoteURL("origin"))
}

func zshRepoOpenURLCmd(webURL string) tea.Cmd {
	return func() tea.Msg {
		err := sysurl.Open(webURL)
		return btmsg.URLActionDone{Verb: "open", Err: err}
	}
}

// gitRemoteToWebURL converts a git remote URL to HTTPS for opening in a browser.
func gitRemoteToWebURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return strings.TrimSuffix(raw, ".git")
	}
	if strings.HasPrefix(raw, "git@") {
		rest := strings.TrimPrefix(raw, "git@")
		colon := strings.Index(rest, ":")
		if colon <= 0 {
			return ""
		}
		host := rest[:colon]
		path := strings.TrimSuffix(rest[colon+1:], ".git")
		return "https://" + host + "/" + path
	}
	if strings.HasPrefix(raw, "ssh://") {
		u := strings.TrimPrefix(raw, "ssh://")
		u = strings.TrimSuffix(u, ".git")
		if strings.HasPrefix(u, "git@") {
			u = strings.TrimPrefix(u, "git@")
		}
		slash := strings.Index(u, "/")
		if slash <= 0 {
			return ""
		}
		return "https://" + u[:slash] + "/" + u[slash+1:]
	}
	return ""
}

// IsDone returns true when the wizard completed (user left success/error or cancelled)
func (m ZshRepoModel) IsDone() bool {
	return m.done
}

// IsCancelled returns true when the user cancelled
func (m ZshRepoModel) IsCancelled() bool {
	return m.cancelled
}
