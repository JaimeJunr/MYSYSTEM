package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/catalog"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/preferences"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/profilestate"
	"github.com/JaimeJunr/Homestead/internal/monitoring"
	"github.com/JaimeJunr/Homestead/internal/tui/cmds"
	"github.com/JaimeJunr/Homestead/internal/tui/items"
	btmsg "github.com/JaimeJunr/Homestead/internal/tui/msg"
	"github.com/JaimeJunr/Homestead/internal/tui/theme"
)

// Model is the main TUI model
type Model struct {
	scriptService    *services.ScriptService
	installerService *services.InstallerService
	configService    *services.ConfigService
	repoService      *services.RepoService
	state            ViewState
	mainMenu         list.Model
	scriptList       list.Model
	installerList    list.Model
	packageList      list.Model
	selectedMenu     int
	selectedItem     interface{} // Can be Script or Package
	confirmYes       bool
	confirmReturn    ViewState // view to restore when confirmation is cancelled
	confirmReturnOK  bool      // if false, cancel returns to main menu
	scriptDryRunNext bool      // next bash run uses HOMESTEAD_DRY_RUN=1 (key d)
	width            int
	height           int
	err              error
	keyboardToast    string // transient open/copy URL feedback

	// Installation progress
	progress       progress.Model
	spinner        spinner.Model
	installStatus  string
	installMessage string
	installPercent float64
	canAbort       bool
	aborted        bool

	zshWizard *ZshWizardModel

	zshRepoWizard *ZshRepoModel

	zshCoreInstalled bool
	zshCoreChecked   bool

	zshApplyPhase string
	zshApplyError error

	scriptOutputView  viewport.Model
	scriptOutputPhase string
	scriptOutputTitle string
	scriptOutputErr   error

	nativeMonitorKind string
	nativeBattery     *monitoring.BatterySnapshot
	nativeBatteryErr  error
	nativeMemory      *monitoring.MemorySnapshot
	nativeMemoryErr   error
	nativeDisk        []monitoring.DiskMount
	nativeDiskErr     error
	nativeLoad        *monitoring.LoadSnapshot
	nativeLoadErr     error
	nativeNetwork     *monitoring.NetworkSnapshot
	nativeNetworkErr  error
	nativeNetworkAt   time.Time
	nativeNetRates    map[string]monitoring.NetRates
	nativeThermal     *monitoring.ThermalSnapshot
	nativeThermalErr  error
	nativeSystemdUser *monitoring.SystemdUserFailedSnapshot
	nativeSystemdErr  error

	scriptListParent ViewState // Esc target from script list (main menu or installer categories)
	scriptListCategory types.Category

	profile                 *profilestate.State
	profilePath             string
	pendingInstallPackageID string

	packageListCategories []types.PackageCategory
	catalogURL            string

	prefs         preferences.Preferences
	prefsPath     string
	catalogEnvSet bool
	settingsModel *SettingsModel

	helpOpen bool
}

// NewModel wires the TUI. Empty catalogURL skips remote catalog fetch (tests).
func NewModel(scriptService *services.ScriptService, installerService *services.InstallerService, configService *services.ConfigService, repoService *services.RepoService, catalogURL string, prefs preferences.Preferences, prefsPath string, catalogEnvSet bool, profile *profilestate.State, profilePath string) Model {
	prefs.Normalize()
	theme.ApplyPreferences(prefs)

	mainItems := getMainMenuItems(false)
	mainList := list.New(mainItems, list.NewDefaultDelegate(), 0, 0)
	mainList.Title = "Homestead - Gerenciador de Sistema"
	mainList.SetShowStatusBar(false)
	mainList.SetFilteringEnabled(false)

	progOpts := []progress.Option{progress.WithWidth(40)}
	if prefs.ReduceMotion {
		progOpts = append(progOpts, progress.WithSolidFill("#888888"))
	} else {
		progOpts = append(progOpts, progress.WithDefaultGradient())
	}
	prog := progress.New(progOpts...)

	spin := spinner.New()
	spin.Spinner = spinner.Dot

	if profile == nil {
		profile = &profilestate.State{}
	}
	return Model{
		scriptService:         scriptService,
		installerService:      installerService,
		configService:         configService,
		repoService:           repoService,
		catalogURL:            catalogURL,
		prefs:                 prefs,
		prefsPath:             prefsPath,
		catalogEnvSet:         catalogEnvSet,
		profile:               profile,
		profilePath:           profilePath,
		state:                 ViewMainMenu,
		mainMenu:              mainList,
		progress:              prog,
		spinner:               spin,
		confirmYes:            false,
		scriptDryRunNext:      false,
		scriptListParent:      ViewMainMenu,
	}
}

func listAbsorbsFilterKey(l list.Model, msg tea.KeyMsg) (list.Model, tea.Cmd, bool) {
	if l.SettingFilter() {
		next, cmd := l.Update(msg)
		return next, cmd, true
	}
	if l.IsFiltered() && msg.Type == tea.KeyEsc {
		next, cmd := l.Update(msg)
		return next, cmd, true
	}
	return l, nil, false
}

func forwardListFilterKey(m Model, msg tea.KeyMsg) (Model, tea.Cmd, bool) {
	switch m.state {
	case ViewScriptList:
		if next, cmd, ok := listAbsorbsFilterKey(m.scriptList, msg); ok {
			m.scriptList = next
			return m, cmd, true
		}
	case ViewPackageList:
		if next, cmd, ok := listAbsorbsFilterKey(m.packageList, msg); ok {
			m.packageList = next
			return m, cmd, true
		}
	}
	return m, nil, false
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	batch := []tea.Cmd{cmds.CheckZshCoreInstalled(m.installerService)}
	if !m.prefs.ReduceMotion {
		batch = append([]tea.Cmd{m.spinner.Tick}, batch...)
	}
	if c := cmds.FetchCatalog(m.catalogURL, m.installerService); c != nil {
		batch = append(batch, c)
	}
	return tea.Batch(batch...)
}

// Update handles messages and updates state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if m.state == ViewScriptOutput && m.scriptOutputPhase == "done" {
			var vcmd tea.Cmd
			m.scriptOutputView, vcmd = m.scriptOutputView.Update(msg)
			return m, vcmd
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.state == ViewSettings && m.settingsModel != nil {
			newS, _ := m.settingsModel.Update(msg)
			sm := newS.(SettingsModel)
			m.settingsModel = &sm
		}
		reserve := theme.ListVerticalReserve()
		m.mainMenu.SetSize(msg.Width, msg.Height-reserve)
		if m.scriptList.Items() != nil {
			m.scriptList.SetSize(msg.Width, msg.Height-reserve)
		}
		if m.installerList.Items() != nil {
			m.installerList.SetSize(msg.Width, msg.Height-reserve)
		}
		if m.packageList.Items() != nil {
			m.packageList.SetSize(msg.Width, msg.Height-reserve)
		}
		if m.state == ViewScriptOutput {
			m.syncScriptOutputViewport()
		}
		return m, nil

	case btmsg.Progress:
		m.installStatus = msg.Status
		m.installMessage = msg.Message
		m.installPercent = float64(msg.Progress) / 100.0
		m.canAbort = msg.CanAbort

		if msg.IsCompleted {
			return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
				return btmsg.InstallComplete{Err: msg.Error}
			})
		}
		return m, nil

	case btmsg.InstallComplete:
		m.state = ViewMainMenu
		m.aborted = false
		if msg.Err == nil && m.profile != nil && strings.TrimSpace(m.pendingInstallPackageID) != "" && strings.TrimSpace(m.profilePath) != "" {
			profilestate.RecordInstalled(m.profile, m.pendingInstallPackageID)
			if err := profilestate.Save(m.profilePath, *m.profile); err != nil {
				m.err = fmt.Errorf("gravar perfil: %w", err)
			}
		}
		m.pendingInstallPackageID = ""
		return m, cmds.CheckZshCoreInstalled(m.installerService)

	case btmsg.ZshCoreInstalled:
		m.zshCoreChecked = true
		m.zshCoreInstalled = msg.Installed
		m.mainMenu.SetItems(getMainMenuItems(m.zshCoreInstalled))
		return m, nil

	case spinner.TickMsg:
		if m.prefs.ReduceMotion {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if m.helpOpen {
			switch msg.String() {
			case "esc", "?", "q":
				m.helpOpen = false
			}
			return m, nil
		}
		if msg.String() == "?" && !m.suppressHelpHotkey() {
			m.helpOpen = true
			return m, nil
		}

		if m.state == ViewSettings && m.settingsModel != nil {
			newS, cmd := m.settingsModel.Update(msg)
			sm := newS.(SettingsModel)
			m.settingsModel = &sm
			return m, cmd
		}
		if m.state == ViewScriptOutput {
			if m.scriptOutputPhase == "done" {
				switch msg.String() {
				case "enter", "esc", "q":
					m.state = m.confirmReturn
					m.scriptOutputPhase = ""
					m.scriptOutputTitle = ""
					m.scriptOutputErr = nil
					return m, nil
				}
				var vcmd tea.Cmd
				m.scriptOutputView, vcmd = m.scriptOutputView.Update(msg)
				return m, vcmd
			}
			return m, nil
		}
		if m.state == ViewNativeMonitor {
			switch msg.String() {
			case "enter", "esc", "q":
				m.state = m.confirmReturn
				m.nativeMonitorKind = ""
				m = m.withClearedNativeMonitors()
				return m, nil
			case "r":
				return m, m.nativeMonitorLoadCmd()
			}
			return m, nil
		}
		if m.state == ViewScriptList && m.err != nil {
			m.err = nil
		}
		if m.state == ViewZshApplying && (m.zshApplyPhase == "success" || m.zshApplyPhase == "error") {
			if msg.String() == "enter" || msg.String() == "esc" {
				m.state = ViewMainMenu
				m.zshApplyPhase = ""
				m.zshApplyError = nil
				return m, nil
			}
		}
		if m.state != ViewZshWizard && m.state != ViewZshRepoWizard {
			if next, cmd, ok := forwardListFilterKey(m, msg); ok {
				return next, cmd
			}
			switch msg.String() {
			case "f":
				if m.state == ViewScriptList {
					return m.handleToggleScriptFavorite()
				}
			case "d", "D":
				if m.state == ViewScriptList {
					return m.handleScriptDryRun()
				}
			case "ctrl+c", "q":
				if m.state == ViewMainMenu {
					return m, tea.Quit
				}
				if m.state == ViewInstalling && m.canAbort {
					m.aborted = true
					m.installMessage = "Instalação abortada pelo usuário"
					m.state = ViewMainMenu
					m.pendingInstallPackageID = ""
					return m, nil
				}
			case "esc":
				switch m.state {
				case ViewScriptList:
					m.state = m.scriptListParent
					m.confirmYes = false
					return m, nil
				case ViewConfirmation:
					if m.confirmReturnOK {
						m.state = m.confirmReturn
					} else {
						m.state = ViewMainMenu
					}
					m.confirmReturnOK = false
					m.confirmYes = false
					return m, nil
				case ViewPackageList:
					m.state = ViewInstallerCategories
					return m, nil
				case ViewInstallerCategories:
					m.state = ViewMainMenu
					return m, nil
				}
			case "left", "h":
				if m.state == ViewConfirmation {
					m.confirmYes = false
					return m, nil
				}
			case "right", "l":
				if m.state == ViewConfirmation {
					m.confirmYes = true
					return m, nil
				}
			case "o", "O":
				return m.handleURLShortcut(false)
			case "c", "C":
				return m.handleURLShortcut(true)
			case "enter":
				return m.handleEnter()
			}
		}

	case btmsg.ZshApplyResult:
		if m.state == ViewZshApplying {
			if msg.Err != nil {
				m.zshApplyPhase = "error"
				m.zshApplyError = msg.Err
			} else {
				m.zshApplyPhase = "success"
				m.zshApplyError = nil
			}
			return m, tea.Tick(time.Second*2, func(time.Time) tea.Msg {
				return btmsg.ZshApplyReturnToMenu{}
			})
		}
		return m, nil

	case btmsg.ZshApplyReturnToMenu:
		if m.state == ViewZshApplying {
			m.state = ViewMainMenu
			m.zshApplyPhase = ""
			m.zshApplyError = nil
		}
		return m, nil

	case btmsg.ScriptCaptured:
		if m.state != ViewScriptOutput {
			return m, nil
		}
		m.scriptOutputPhase = "done"
		m.scriptOutputErr = msg.Err
		text := theme.StripANSI(msg.Output)
		if strings.TrimSpace(text) == "" {
			text = "(sem saída no stdout/stderr)"
		}
		if msg.Err != nil {
			text += "\n\n──\n" + msg.Err.Error()
		}
		m.scriptOutputView.SetContent(text)
		m.scriptOutputView.GotoTop()
		return m, nil

	case btmsg.ScriptExecFinished:
		m.state = m.confirmReturn
		m.scriptOutputPhase = ""
		m.scriptOutputTitle = ""
		m.scriptOutputErr = nil
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.err = nil
		}
		return m, nil

	case btmsg.NativeMonitorReload:
		if m.state != ViewNativeMonitor || msg.Kind != m.nativeMonitorKind {
			return m, nil
		}
		switch msg.Kind {
		case entities.NativeMonitorBattery:
			m.nativeBattery = msg.Battery
			m.nativeBatteryErr = msg.Err
		case entities.NativeMonitorMemory:
			m.nativeMemory = msg.Memory
			m.nativeMemoryErr = msg.Err
		case entities.NativeMonitorDisk:
			m.nativeDisk = msg.Disk
			m.nativeDiskErr = msg.Err
		case entities.NativeMonitorLoad:
			m.nativeLoad = msg.Load
			m.nativeLoadErr = msg.Err
		case entities.NativeMonitorNetwork:
			now := time.Now()
			if msg.Err == nil && msg.Network != nil && m.nativeNetwork != nil && !m.nativeNetworkAt.IsZero() {
				dt := now.Sub(m.nativeNetworkAt).Seconds()
				if dt >= 0.2 {
					m.nativeNetRates = monitoring.ComputeNetRates(m.nativeNetwork, msg.Network, dt)
				}
			}
			if msg.Err != nil {
				m.nativeNetwork = nil
				m.nativeNetworkAt = time.Time{}
				m.nativeNetRates = nil
			} else {
				m.nativeNetwork = msg.Network
				m.nativeNetworkAt = now
			}
			m.nativeNetworkErr = msg.Err
		case entities.NativeMonitorThermal:
			m.nativeThermal = msg.Thermal
			m.nativeThermalErr = msg.Err
		case entities.NativeMonitorSystemdUser:
			m.nativeSystemdUser = msg.SystemdUser
			m.nativeSystemdErr = msg.Err
		}
		return m, m.nativeMonitorScheduleTick()

	case btmsg.NativeMonitorTick:
		if m.state != ViewNativeMonitor {
			return m, nil
		}
		return m, m.nativeMonitorLoadCmd()

	case btmsg.CatalogFetched:
		if msg.Err != nil {
			return m, nil
		}
		var nextCmd tea.Cmd
		if msg.Ok {
			if m.state == ViewPackageList && len(m.packageListCategories) > 0 {
				sel := m.packageList.Index()
				m.loadPackagesFromCategories(m.packageListCategories)
				pkgRows := m.packageList.Items()
				if len(pkgRows) > 0 {
					if sel < 0 {
						sel = 0
					}
					if sel >= len(pkgRows) {
						sel = len(pkgRows) - 1
					}
					m.packageList.Select(sel)
				}
			}
			m.keyboardToast = "Catálogo de instaladores atualizado."
			nextCmd = tea.Tick(2*time.Second, func(time.Time) tea.Msg { return btmsg.ClearKeyboardToast{} })
		}
		return m, nextCmd

	case btmsg.URLActionDone:
		if msg.Err != nil {
			m.keyboardToast = fmt.Sprintf("⚠ %v", msg.Err)
		} else if msg.Verb == "copy" {
			m.keyboardToast = "URL copiada para a área de transferência."
		} else {
			m.keyboardToast = "URL aberta no navegador (app padrão)."
		}
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return btmsg.ClearKeyboardToast{} })

	case btmsg.ClearKeyboardToast:
		m.keyboardToast = ""
		return m, nil

	case btmsg.SettingsSaved:
		m.state = ViewMainMenu
		m.settingsModel = nil
		return m.applySavedPreferences(msg.Prefs)

	case btmsg.SettingsCancelled:
		m.state = ViewMainMenu
		m.settingsModel = nil
		return m, nil
	}

	if m.state == ViewSettings && m.settingsModel != nil {
		newS, cmd := m.settingsModel.Update(msg)
		sm := newS.(SettingsModel)
		m.settingsModel = &sm
		return m, cmd
	}

	if m.state == ViewZshWizard && m.zshWizard != nil {
		newWizard, cmd := m.zshWizard.Update(msg)
		wizard := newWizard.(ZshWizardModel)
		m.zshWizard = &wizard

		if wizard.IsDone() || wizard.IsCancelled() {
			if wizard.IsCancelled() {
				m.state = ViewMainMenu
				m.zshWizard = nil
				return m, cmd
			}
			selections := wizard.GetSelections()
			m.zshWizard = nil
			m.state = ViewZshApplying
			m.zshApplyPhase = "applying"
			m.zshApplyError = nil
			return m, cmds.ApplyZshConfig(m.configService, selections)
		}

		return m, cmd
	}

	if m.state == ViewZshRepoWizard && m.zshRepoWizard != nil {
		newRepo, cmd := m.zshRepoWizard.Update(msg)
		repoWizard := newRepo.(ZshRepoModel)
		m.zshRepoWizard = &repoWizard

		if repoWizard.IsDone() || repoWizard.IsCancelled() {
			m.state = ViewMainMenu
			m.zshRepoWizard = nil
			return m, cmd
		}
		return m, cmd
	}

	var cmd tea.Cmd
	switch m.state {
	case ViewMainMenu:
		m.mainMenu, cmd = m.mainMenu.Update(msg)
	case ViewScriptList:
		m.scriptList, cmd = m.scriptList.Update(msg)
	case ViewInstallerCategories:
		m.installerList, cmd = m.installerList.Update(msg)
	case ViewPackageList:
		m.packageList, cmd = m.packageList.Update(msg)
	}

	return m, cmd
}

func (m Model) handleScriptDryRun() (tea.Model, tea.Cmd) {
	sel := m.scriptList.SelectedItem()
	si, ok := sel.(items.ScriptItem)
	if !ok || !si.Script.SupportsDryRun || si.Script.NativeMonitor != "" {
		return m, nil
	}
	m.selectedItem = si.Script
	m.confirmReturn = ViewScriptList
	m.confirmReturnOK = true
	m.scriptDryRunNext = true
	return m.handleConfirmedSelection()
}

func (m Model) handleToggleScriptFavorite() (tea.Model, tea.Cmd) {
	if m.profile == nil {
		return m, nil
	}
	sel := m.scriptList.SelectedItem()
	si, ok := sel.(items.ScriptItem)
	if !ok {
		return m, nil
	}
	profilestate.ToggleFavorite(m.profile, si.Script.ID)
	if strings.TrimSpace(m.profilePath) != "" {
		if err := profilestate.Save(m.profilePath, *m.profile); err != nil {
			m.err = fmt.Errorf("gravar perfil: %w", err)
			return m, nil
		}
	}
	m.reloadScriptList()
	return m, nil
}

func (m Model) handleConfirmedSelection() (tea.Model, tea.Cmd) {
	switch item := m.selectedItem.(type) {
	case entities.Script:
		if item.NativeMonitor != "" {
			m.scriptDryRunNext = false
			m = m.withClearedNativeMonitors()
			m.nativeMonitorKind = item.NativeMonitor
			m.state = ViewNativeMonitor
			return m, m.nativeMonitorLoadCmd()
		}
		dry := m.scriptDryRunNext
		m.scriptDryRunNext = false
		opts := interfaces.ScriptExecOpts{DryRun: dry}
		m.scriptOutputTitle = item.Name
		if dry {
			m.scriptOutputTitle = item.Name + " (simulação)"
		}
		m.scriptOutputPhase = "running"
		m.scriptOutputErr = nil
		m.scriptOutputView = newScriptOutputViewport(m.width, m.height)
		m.state = ViewScriptOutput
		if item.RequiresSudo {
			cmd, err := m.scriptService.ScriptInteractiveCommand(item.ID, opts)
			if err != nil {
				m.state = m.confirmReturn
				m.scriptOutputPhase = ""
				m.scriptOutputTitle = ""
				m.err = err
				return m, nil
			}
			return m, tea.ExecProcess(cmd, func(execErr error) tea.Msg {
				return btmsg.ScriptExecFinished{Err: execErr}
			})
		}
		return m, cmds.RunScriptCapture(m.scriptService, item.ID, opts)
	case entities.Package:
		m.state = ViewInstalling
		m.pendingInstallPackageID = item.ID
		m.installStatus = "preparing"
		m.installMessage = "Preparando instalação..."
		m.installPercent = 0
		m.canAbort = false
		m.aborted = false
		return m, cmds.InstallPackage(m.installerService, item.ID)
	}
	return m, nil
}

func (m Model) applySavedPreferences(p preferences.Preferences) (tea.Model, tea.Cmd) {
	p.Normalize()
	theme.ApplyPreferences(p)
	m.prefs = p
	if err := m.scriptService.ConfigureScriptRoot(p.ScriptRoot); err != nil {
		m.err = err
	}
	if err := m.installerService.ConfigureHomesteadRoot(p.ScriptRoot); err != nil {
		m.err = err
	}
	dotfiles := p.DotfilesRepo
	if strings.TrimSpace(dotfiles) == "" {
		dotfiles = preferences.DefaultDotfilesRepo()
	}
	newRepo, err := services.NewRepoService(dotfiles)
	if err != nil {
		m.err = err
	} else {
		m.repoService = newRepo
		if m.zshRepoWizard != nil {
			m.zshRepoWizard.repoService = newRepo
		}
	}
	m.catalogURL = catalog.EffectiveCatalogURL(p.CatalogURL)
	if strings.TrimSpace(m.prefsPath) != "" {
		if err := preferences.Save(m.prefsPath, p); err != nil {
			m.err = err
		}
	}
	if m.width > 0 && m.height > 0 {
		r := theme.ListVerticalReserve()
		m.mainMenu.SetSize(m.width, m.height-r)
		if m.scriptList.Items() != nil {
			m.scriptList.SetSize(m.width, m.height-r)
		}
		if m.installerList.Items() != nil {
			m.installerList.SetSize(m.width, m.height-r)
		}
		if m.packageList.Items() != nil {
			m.packageList.SetSize(m.width, m.height-r)
		}
		m.syncScriptOutputViewport()
	}
	var batch []tea.Cmd
	if c := cmds.FetchCatalog(m.catalogURL, m.installerService); c != nil {
		batch = append(batch, c)
	}
	return m, tea.Batch(batch...)
}

// handleEnter handles the enter key based on current state
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case ViewMainMenu:
		selected := m.mainMenu.SelectedItem()
		item, ok := selected.(items.MenuItem)
		if !ok {
			return m, nil
		}
		switch item.Action {
		case menuActionCleanup:
			m.state = ViewScriptList
			m.selectedMenu = 0
			m.loadScripts(types.CategoryCleanup)
		case menuActionMonitoring:
			m.state = ViewScriptList
			m.selectedMenu = 1
			m.loadScripts(types.CategoryMonitoring)
		case menuActionCheckup:
			m.state = ViewScriptList
			m.selectedMenu = 2
			m.loadScripts(types.CategoryCheckup)
		case menuActionInstallers:
			m.state = ViewInstallerCategories
			m.selectedMenu = 3
			m.loadInstallerCategories()
		case menuActionZshPlugins:
			m.state = ViewZshWizard
			wizardService := services.NewWizardService()
			wizard := NewZshWizardModel(wizardService)
			wizard.width = m.width
			wizard.height = m.height
			m.zshWizard = &wizard
		case menuActionZshRepo:
			m.state = ViewZshRepoWizard
			repoWizard := NewZshRepoModel(m.repoService, m.configService)
			repoWizard.width = m.width
			repoWizard.height = m.height
			m.zshRepoWizard = &repoWizard
		case menuActionQuit:
			return m, tea.Quit
		case menuActionSettings:
			sm := NewSettingsModel(m.prefs, m.prefsPath, m.catalogEnvSet)
			sm.width = m.width
			sm.height = m.height
			m.settingsModel = &sm
			m.state = ViewSettings
			return m, m.settingsModel.Init()
		default:
			return m, nil
		}

	case ViewScriptList:
		selected := m.scriptList.SelectedItem()
		if scriptItem, ok := selected.(items.ScriptItem); ok {
			m.scriptDryRunNext = false
			m.selectedItem = scriptItem.Script
			m.confirmReturn = ViewScriptList
			m.confirmReturnOK = true
			if m.prefs.ConfirmBeforeScript {
				m.state = ViewConfirmation
				m.confirmYes = false
			} else {
				return m.handleConfirmedSelection()
			}
		}

	case ViewPackageList:
		selected := m.packageList.SelectedItem()
		if pkgItem, ok := selected.(items.PackageItem); ok {
			m.selectedItem = pkgItem.Pkg
			m.confirmReturn = ViewPackageList
			m.confirmReturnOK = true
			if m.prefs.ConfirmBeforePackage {
				m.state = ViewConfirmation
				m.confirmYes = false
			} else {
				return m.handleConfirmedSelection()
			}
		}

	case ViewInstallerCategories:
		selected := m.installerList.SelectedItem()
		catItem, ok := selected.(items.InstallerCategoryItem)
		if !ok {
			break
		}
		if len(catItem.Categories) > 0 {
			m.state = ViewPackageList
			m.packageListCategories = append([]types.PackageCategory(nil), catItem.Categories...)
			m.loadPackagesFromCategories(catItem.Categories)
		}

	case ViewConfirmation:
		if m.confirmYes {
			return m.handleConfirmedSelection()
		} else {
			if m.confirmReturnOK {
				m.state = m.confirmReturn
			} else {
				m.state = ViewMainMenu
			}
			m.confirmReturnOK = false
		}
	}

	return m, nil
}
