package tui

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/config"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/executor"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/installer"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/repository"
)

// testModel creates a model for testing with mocked dependencies
func testModel() Model {
	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExec := executor.NewBashExecutor()
	scriptService := services.NewScriptService(scriptRepo, scriptExec)

	packageRepo := repository.NewInMemoryPackageRepository()
	packageInstaller := installer.NewDefaultPackageInstaller()
	installerService := services.NewInstallerService(packageRepo, packageInstaller)

	configManager := config.NewFileConfigManager("")
	configService := services.NewConfigService(configManager)

	return NewModel(scriptService, installerService, configService, nil, "")
}

func TestNewModel(t *testing.T) {
	model := testModel()

	if model.state != ViewMainMenu {
		t.Errorf("Expected initial state to be ViewMainMenu, got %d", model.state)
	}

	items := model.mainMenu.Items()
	// 6 items when zsh core not installed: Limpeza, Monitoramento, Instaladores, Configurar Zsh, Configurações, Sair
	if len(items) != 6 {
		t.Errorf("Expected 6 main menu items (zsh core not installed), got %d", len(items))
	}

	if model.scriptService == nil {
		t.Error("Expected scriptService to be initialized")
	}
}

func TestViewStates(t *testing.T) {
	// Verify view state constants
	if ViewMainMenu != 0 {
		t.Errorf("ViewMainMenu should be 0, got %d", ViewMainMenu)
	}
	if ViewScriptList != 1 {
		t.Errorf("ViewScriptList should be 1, got %d", ViewScriptList)
	}
	if ViewInstallerCategories != 2 {
		t.Errorf("ViewInstallerCategories should be 2, got %d", ViewInstallerCategories)
	}
	if ViewPackageList != 3 {
		t.Errorf("ViewPackageList should be 3, got %d", ViewPackageList)
	}
	if ViewConfirmation != 4 {
		t.Errorf("ViewConfirmation should be 4, got %d", ViewConfirmation)
	}
	if ViewScriptOutput != 5 {
		t.Errorf("ViewScriptOutput should be 5, got %d", ViewScriptOutput)
	}
	if ViewNativeMonitor != 6 {
		t.Errorf("ViewNativeMonitor should be 6, got %d", ViewNativeMonitor)
	}
	if ViewInstalling != 7 {
		t.Errorf("ViewInstalling should be 7, got %d", ViewInstalling)
	}
	if ViewZshWizard != 8 {
		t.Errorf("ViewZshWizard should be 8, got %d", ViewZshWizard)
	}
	if ViewZshApplying != 9 {
		t.Errorf("ViewZshApplying should be 9, got %d", ViewZshApplying)
	}
}

func TestModelInit(t *testing.T) {
	model := testModel()
	cmd := model.Init()

	if cmd == nil {
		t.Error("Expected Init() to return spinner tick command")
	}
}

func TestWindowSizeUpdate(t *testing.T) {
	model := testModel()

	msg := tea.WindowSizeMsg{
		Width:  80,
		Height: 24,
	}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.width != 80 {
		t.Errorf("Expected width 80, got %d", m.width)
	}
	if m.height != 24 {
		t.Errorf("Expected height 24, got %d", m.height)
	}
}

func TestQuitOnMainMenu(t *testing.T) {
	model := testModel()
	model.state = ViewMainMenu

	// Test 'q' key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Error("Expected quit command on 'q' key")
	}

	// Test Ctrl+C
	msgCtrlC := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmdCtrlC := model.Update(msgCtrlC)

	if cmdCtrlC == nil {
		t.Error("Expected quit command on Ctrl+C")
	}
}

func TestEscapeFromScriptList(t *testing.T) {
	model := testModel()
	model.state = ViewScriptList

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.state != ViewMainMenu {
		t.Errorf("Expected state to return to ViewMainMenu, got %d", m.state)
	}
}

func TestMenuItemInterface(t *testing.T) {
	item := menuItem{
		title: "Test Title",
		desc:  "Test Description",
	}

	if item.Title() != "Test Title" {
		t.Errorf("Expected title 'Test Title', got %s", item.Title())
	}

	if item.Description() != "Test Description" {
		t.Errorf("Expected description 'Test Description', got %s", item.Description())
	}

	if item.FilterValue() != "Test Title" {
		t.Errorf("Expected filter value 'Test Title', got %s", item.FilterValue())
	}
}

func TestScriptItemInterface(t *testing.T) {
	// Tested via integration
}

func TestViewRendering(t *testing.T) {
	model := testModel()

	// Test initial view (no size set)
	view := model.View()
	if view != "Iniciando..." {
		t.Errorf("Expected 'Iniciando...' for uninitialized view, got %s", view)
	}

	// Set window size
	model.width = 80
	model.height = 24

	// Test main menu view
	model.state = ViewMainMenu
	view = model.View()
	if view == "" {
		t.Error("Expected non-empty view for main menu")
	}
}

func TestModelStateTransitions(t *testing.T) {
	tests := []struct {
		name          string
		initialState  ViewState
		expectedState ViewState
	}{
		{
			name:          "Start at main menu",
			initialState:  ViewMainMenu,
			expectedState: ViewMainMenu,
		},
		{
			name:          "Can be at script list",
			initialState:  ViewScriptList,
			expectedState: ViewScriptList,
		},
		{
			name:          "Can be at script output",
			initialState:  ViewScriptOutput,
			expectedState: ViewScriptOutput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testModel()
			model.state = tt.initialState

			if model.state != tt.expectedState {
				t.Errorf("Expected state %d, got %d", tt.expectedState, model.state)
			}
		})
	}
}

// Benchmark tests
// noOpMsg is a message that falls through Update's switch so the wizard delegate block runs
type noOpMsg struct{}

// TestZshWizardDoneTriggersApply tests that when the wizard completes (done, not cancelled),
// the model transitions to ViewZshApplying and returns the apply Cmd.
func TestZshWizardDoneTriggersApply(t *testing.T) {
	model := testModel()
	wizardService := services.NewWizardService()
	wizard := NewZshWizardModel(wizardService)
	wizard.width = 80
	wizard.height = 24

	// Advance wizard to Review (Plugins -> Tools -> ProjectConfig -> Review): 3x 'n', then Enter
	for i := 0; i < 3; i++ {
		var next tea.Model
		next, _ = wizard.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
		wizard = next.(ZshWizardModel)
	}
	// Confirm on Review (Enter) so wizard is done
	{
		var next tea.Model
		next, _ = wizard.Update(tea.KeyMsg{Type: tea.KeyEnter})
		wizard = next.(ZshWizardModel)
	}
	if !wizard.IsDone() {
		t.Fatal("wizard should be done after Enter on Review")
	}

	model.state = ViewZshWizard
	model.zshWizard = &wizard
	model.width = 80
	model.height = 24

	// Use a message that falls through the switch so we reach the wizard delegate block
	updated, cmd := model.Update(noOpMsg{})
	m := updated.(Model)

	if m.state != ViewZshApplying {
		t.Errorf("state = %d, want ViewZshApplying (%d)", m.state, ViewZshApplying)
	}
	if m.zshApplyPhase != "applying" {
		t.Errorf("zshApplyPhase = %q, want \"applying\"", m.zshApplyPhase)
	}
	if cmd == nil {
		t.Error("expected non-nil Cmd (apply config)")
	}
}

// TestZshWizardCancelledDoesNotApply tests that when the wizard is cancelled, we return to main menu.
func TestZshWizardCancelledDoesNotApply(t *testing.T) {
	model := testModel()
	wizardService := services.NewWizardService()
	wizard := NewZshWizardModel(wizardService)
	wizard.width = 80
	wizard.height = 24
	wizard.cancelled = true
	wizard.done = true

	model.state = ViewZshWizard
	model.zshWizard = &wizard

	updated, _ := model.Update(noOpMsg{})
	m := updated.(Model)

	if m.state != ViewMainMenu {
		t.Errorf("state = %d, want ViewMainMenu", m.state)
	}
	if m.zshWizard != nil {
		t.Error("zshWizard should be nil after cancel")
	}
}

// TestZshApplyResultSuccess tests that zshApplyResultMsg with nil error sets success and schedules return.
func TestZshApplyResultSuccess(t *testing.T) {
	model := testModel()
	model.state = ViewZshApplying
	model.zshApplyPhase = "applying"

	updated, cmd := model.Update(zshApplyResultMsg{Err: nil})
	m := updated.(Model)

	if m.zshApplyPhase != "success" {
		t.Errorf("zshApplyPhase = %q, want \"success\"", m.zshApplyPhase)
	}
	if m.zshApplyError != nil {
		t.Errorf("zshApplyError = %v, want nil", m.zshApplyError)
	}
	if cmd == nil {
		t.Error("expected non-nil Cmd (tick to return to menu)")
	}
}

// TestZshApplyResultError tests that zshApplyResultMsg with error sets error phase.
func TestZshApplyResultError(t *testing.T) {
	model := testModel()
	model.state = ViewZshApplying
	model.zshApplyPhase = "applying"
	err := fmt.Errorf("test apply error")

	updated, _ := model.Update(zshApplyResultMsg{Err: err})
	m := updated.(Model)

	if m.zshApplyPhase != "error" {
		t.Errorf("zshApplyPhase = %q, want \"error\"", m.zshApplyPhase)
	}
	if m.zshApplyError != err {
		t.Errorf("zshApplyError = %v, want %v", m.zshApplyError, err)
	}
}

// TestZshApplyReturnToMenuMsg tests that zshApplyReturnToMenuMsg clears state and returns to main menu.
func TestZshApplyReturnToMenuMsg(t *testing.T) {
	model := testModel()
	model.state = ViewZshApplying
	model.zshApplyPhase = "success"

	updated, _ := model.Update(zshApplyReturnToMenuMsg{})
	m := updated.(Model)

	if m.state != ViewMainMenu {
		t.Errorf("state = %d, want ViewMainMenu", m.state)
	}
	if m.zshApplyPhase != "" {
		t.Errorf("zshApplyPhase = %q, want empty", m.zshApplyPhase)
	}
	if m.zshApplyError != nil {
		t.Errorf("zshApplyError should be nil after return")
	}
}

// TestZshApplyResultStateEnterReturnsToMenu tests that Enter in success/error phase returns to menu.
func TestZshApplyResultStateEnterReturnsToMenu(t *testing.T) {
	model := testModel()
	model.state = ViewZshApplying
	model.zshApplyPhase = "success"

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updated.(Model)

	if m.state != ViewMainMenu {
		t.Errorf("state = %d, want ViewMainMenu after Enter", m.state)
	}
}

func BenchmarkNewModel(b *testing.B) {
	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExec := executor.NewBashExecutor()
	scriptService := services.NewScriptService(scriptRepo, scriptExec)

	packageRepo := repository.NewInMemoryPackageRepository()
	packageInstaller := installer.NewDefaultPackageInstaller()
	installerService := services.NewInstallerService(packageRepo, packageInstaller)

	configManager := config.NewFileConfigManager("")
	configService := services.NewConfigService(configManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewModel(scriptService, installerService, configService, nil, "")
	}
}

func BenchmarkModelUpdate(b *testing.B) {
	model := testModel()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.Update(msg)
	}
}

func BenchmarkModelView(b *testing.B) {
	model := testModel()
	model.width = 80
	model.height = 24

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.View()
	}
}
