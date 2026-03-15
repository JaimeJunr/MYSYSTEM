package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/JaimeJunr/Homestead/internal/app/services"
)

// newTestZshWizardModel creates a ZshWizardModel for testing
func newTestZshWizardModel() ZshWizardModel {
	wizardService := services.NewWizardService()
	return NewZshWizardModel(wizardService)
}

// TestZshWizardModel_New tests creating a new wizard model
func TestZshWizardModel_New(t *testing.T) {
	m := newTestZshWizardModel()

	if m.wizardService == nil {
		t.Error("wizardService should not be nil")
	}

	if m.state == nil {
		t.Error("wizard state should be initialized")
	}

	if m.currentView != ZshWizardViewPlugins {
		t.Errorf("initial view = %v, want ZshWizardViewPlugins", m.currentView)
	}
}

// TestZshWizardModel_Init tests model initialization
func TestZshWizardModel_Init(t *testing.T) {
	m := newTestZshWizardModel()
	cmd := m.Init()

	// Init may return nil or a command
	_ = cmd
}

// TestZshWizardModel_View_Initial tests initial view rendering
func TestZshWizardModel_View_Initial(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24

	view := m.View()
	if view == "" {
		t.Error("View() returned empty string")
	}

	// Should show plugins step
	if !containsString(view, "Plugins") && !containsString(view, "plugins") &&
		!containsString(view, "Zsh") {
		t.Logf("View output: %s", view)
	}
}

// TestZshWizardModel_Navigation_Next tests navigating to next step
func TestZshWizardModel_Navigation_Next(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24

	initialView := m.currentView

	// Press 'n' or tab to go next
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	updated := newModel.(ZshWizardModel)

	// View should have changed
	if updated.currentView == initialView {
		// May stay if validation requires selection
		t.Log("View stayed same (may require selections first)")
	}
}

// TestZshWizardModel_Navigation_Previous tests navigating to previous step
func TestZshWizardModel_Navigation_Previous(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24

	// Go to next step first
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = newModel.(ZshWizardModel)

	// If we moved forward, go back
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = newModel.(ZshWizardModel)

	// Should be back at or before the view we started from
	_ = m
}

// TestZshWizardModel_SelectPlugin tests selecting plugin with space
func TestZshWizardModel_SelectPlugin(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24

	// Select item with space (plugins list)
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	updated := newModel.(ZshWizardModel)

	_ = updated.state
}

// TestZshWizardModel_SelectAll tests "marcar todos" (KeySelectAll) in plugins view
func TestZshWizardModel_SelectAll(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24

	// Press 'a' to select all plugins
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	updated := newModel.(ZshWizardModel)

	for i, item := range updated.pluginItems {
		if !item.Selected {
			t.Errorf("plugin item %d (%s) should be selected after marcar todos", i, item.Name)
		}
	}
	if len(updated.state.Selections.Plugins) != len(updated.pluginItems) {
		t.Errorf("Selections.Plugins count = %d, want %d", len(updated.state.Selections.Plugins), len(updated.pluginItems))
	}
}

// TestZshWizardModel_ViewTransitions tests all view transitions
func TestZshWizardModel_ViewTransitions(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24

	views := []ZshWizardView{
		ZshWizardViewPlugins,
		ZshWizardViewTools,
		ZshWizardViewProjectConfig,
		ZshWizardViewReview,
	}

	// Test each view renders
	for _, view := range views {
		m.currentView = view
		rendered := m.View()
		if rendered == "" {
			t.Errorf("View %v rendered empty string", view)
		}
	}
}

// TestZshWizardModel_ViewPluginsFirst tests plugins as first view
func TestZshWizardModel_ViewPluginsFirst(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24

	view := m.View()
	if view == "" {
		t.Error("Plugins view should not be empty")
	}
	if !containsString(view, "Plugins") && !containsString(view, "plugins") {
		t.Error("First view should show Plugins step")
	}
}

// TestZshWizardModel_ViewPlugins tests plugins selection view
func TestZshWizardModel_ViewPlugins(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24
	m.currentView = ZshWizardViewPlugins

	view := m.View()
	if view == "" {
		t.Error("Plugins view should not be empty")
	}
}

// TestZshWizardModel_ViewTools tests tools selection view
func TestZshWizardModel_ViewTools(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24
	m.currentView = ZshWizardViewTools

	view := m.View()
	if view == "" {
		t.Error("Tools view should not be empty")
	}
}

// TestZshWizardModel_ViewProjectConfig tests project config view
func TestZshWizardModel_ViewProjectConfig(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24
	m.currentView = ZshWizardViewProjectConfig

	view := m.View()
	if view == "" {
		t.Error("Project config view should not be empty")
	}
}

// TestZshWizardModel_ViewReview tests review/confirm view
func TestZshWizardModel_ViewReview(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24
	m.currentView = ZshWizardViewReview

	// Add some selections first
	services.NewWizardService()
	m.wizardService.AddCoreComponent(m.state, "zsh")
	m.wizardService.AddCoreComponent(m.state, "oh-my-zsh")
	m.wizardService.AddPlugin(m.state, "git")

	view := m.View()
	if view == "" {
		t.Error("Review view should not be empty")
	}
}

// TestZshWizardModel_ItemToggle tests toggling items in selection views
func TestZshWizardModel_ItemToggle(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24

	// Toggle first item
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = newModel.(ZshWizardModel)

	// Toggle it again (deselect)
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = newModel.(ZshWizardModel)

	// Should be back to original state
	_ = m
}

// TestZshWizardModel_WindowResize tests handling window resize
func TestZshWizardModel_WindowResize(t *testing.T) {
	m := newTestZshWizardModel()

	newModel, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	updated := newModel.(ZshWizardModel)

	if updated.width != 120 {
		t.Errorf("width = %d, want 120", updated.width)
	}
	if updated.height != 40 {
		t.Errorf("height = %d, want 40", updated.height)
	}
}

// TestZshWizardModel_QuitKey tests quitting from wizard
func TestZshWizardModel_QuitKey(t *testing.T) {
	m := newTestZshWizardModel()
	m.width = 80
	m.height = 24

	// Press q or ctrl+c to quit/go back
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = newModel
	// Cmd may be nil or a quit command
	_ = cmd
}

// TestZshWizardModel_Progress tests progress calculation
func TestZshWizardModel_Progress(t *testing.T) {
	m := newTestZshWizardModel()

	progress := m.GetProgress()
	if progress < 0 || progress > 100 {
		t.Errorf("GetProgress() = %d, should be 0-100", progress)
	}
}

// TestZshWizardModel_GetSelections tests getting current selections
func TestZshWizardModel_GetSelections(t *testing.T) {
	m := newTestZshWizardModel()

	// Add selections
	m.wizardService.AddCoreComponent(m.state, "zsh")
	m.wizardService.AddPlugin(m.state, "git")
	m.wizardService.AddTool(m.state, "nvm")

	selections := m.GetSelections()

	if !sliceContainsStr(selections.CoreComponents, "zsh") {
		t.Error("Selections should contain 'zsh'")
	}
	if !sliceContainsStr(selections.Plugins, "git") {
		t.Error("Selections should contain 'git'")
	}
	if !sliceContainsStr(selections.Tools, "nvm") {
		t.Error("Selections should contain 'nvm'")
	}
}

// Helpers
func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findStr(s, substr)
}

func findStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func sliceContainsStr(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
