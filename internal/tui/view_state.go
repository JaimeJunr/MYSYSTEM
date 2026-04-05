package tui

type ViewState int

const (
	ViewMainMenu ViewState = iota
	ViewScriptList
	ViewInstallerCategories
	ViewPackageList
	ViewConfirmation
	ViewScriptOutput
	ViewNativeMonitor
	ViewInstalling
	ViewZshWizard
	ViewZshApplying
	ViewZshRepoWizard
	ViewSettings
)

const (
	menuActionCleanup    = "cleanup"
	menuActionMonitoring = "monitoring"
	menuActionCheckup    = "checkup"
	menuActionInstallers = "installers"
	menuActionZshPlugins = "zsh_plugins"
	menuActionZshRepo    = "zsh_repo"
	menuActionSettings   = "settings"
	menuActionQuit       = "quit"
)
