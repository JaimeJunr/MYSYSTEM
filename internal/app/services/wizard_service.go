package services

import (
	"fmt"
	"strings"

	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
)

// WizardService provides wizard flow management
type WizardService struct {
	steps []WizardStep
}

// WizardState represents the state of a wizard session
type WizardState struct {
	CurrentStep int
	Selections  interfaces.ConfigSelections
	Completed   bool
}

// WizardStep represents a step in the wizard
type WizardStep struct {
	Name        string
	Description string
	Required    bool
}

// NewWizardService creates a new wizard service (core components are installed separately)
func NewWizardService() *WizardService {
	return &WizardService{
		steps: []WizardStep{
			{
				Name:        "Plugins",
				Description: "Select Zsh plugins to install",
				Required:    false,
			},
			{
				Name:        "Development Tools",
				Description: "Select development tools (NVM, Bun, etc)",
				Required:    false,
			},
			{
				Name:        "Project Configuration",
				Description: "Include project-specific configurations",
				Required:    false,
			},
			{
				Name:        "Review & Confirm",
				Description: "Review your selections and apply configuration",
				Required:    true,
			},
		},
	}
}

// CreateNewWizard creates a new wizard session
func (ws *WizardService) CreateNewWizard() *WizardState {
	return &WizardState{
		CurrentStep: 0,
		Selections: interfaces.ConfigSelections{
			CoreComponents:       make([]string, 0),
			Plugins:              make([]string, 0),
			Tools:                make([]string, 0),
			IncludeProjectConfig: false,
			CustomAliases:        make(map[string]string),
			CustomFunctions:      make(map[string]string),
			CustomEnvVars:        make(map[string]string),
		},
		Completed: false,
	}
}

// GetCurrentStep returns the current wizard step
func (ws *WizardService) GetCurrentStep(state *WizardState) *WizardStep {
	if state.CurrentStep < 0 || state.CurrentStep >= len(ws.steps) {
		return nil
	}
	return &ws.steps[state.CurrentStep]
}

// NextStep advances to the next step
func (ws *WizardService) NextStep(state *WizardState) error {
	if state.CurrentStep >= len(ws.steps)-1 {
		return fmt.Errorf("already at last step")
	}

	state.CurrentStep++
	return nil
}

// PreviousStep goes back to the previous step
func (ws *WizardService) PreviousStep(state *WizardState) error {
	if state.CurrentStep <= 0 {
		return fmt.Errorf("already at first step")
	}

	state.CurrentStep--
	return nil
}

// IsFirstStep checks if at first step
func (ws *WizardService) IsFirstStep(state *WizardState) bool {
	return state.CurrentStep == 0
}

// IsLastStep checks if at last step
func (ws *WizardService) IsLastStep(state *WizardState) bool {
	return state.CurrentStep == len(ws.steps)-1
}

// AddCoreComponent adds a core component to selections
func (ws *WizardService) AddCoreComponent(state *WizardState, component string) {
	if !sliceContains(state.Selections.CoreComponents, component) {
		state.Selections.CoreComponents = append(state.Selections.CoreComponents, component)
	}
}

// RemoveCoreComponent removes a core component from selections
func (ws *WizardService) RemoveCoreComponent(state *WizardState, component string) {
	state.Selections.CoreComponents = removeFromSlice(state.Selections.CoreComponents, component)
}

// AddPlugin adds a plugin to selections
func (ws *WizardService) AddPlugin(state *WizardState, plugin string) {
	if !sliceContains(state.Selections.Plugins, plugin) {
		state.Selections.Plugins = append(state.Selections.Plugins, plugin)
	}
}

// RemovePlugin removes a plugin from selections
func (ws *WizardService) RemovePlugin(state *WizardState, plugin string) {
	state.Selections.Plugins = removeFromSlice(state.Selections.Plugins, plugin)
}

// AddTool adds a tool to selections
func (ws *WizardService) AddTool(state *WizardState, tool string) {
	if !sliceContains(state.Selections.Tools, tool) {
		state.Selections.Tools = append(state.Selections.Tools, tool)
	}
}

// RemoveTool removes a tool from selections
func (ws *WizardService) RemoveTool(state *WizardState, tool string) {
	state.Selections.Tools = removeFromSlice(state.Selections.Tools, tool)
}

// SetIncludeProjectConfig sets whether to include project configs
func (ws *WizardService) SetIncludeProjectConfig(state *WizardState, include bool) {
	state.Selections.IncludeProjectConfig = include
}

// GeneratePreview generates a preview of the configuration
func (ws *WizardService) GeneratePreview(state *WizardState) string {
	var builder strings.Builder

	builder.WriteString("=== Configuration Preview ===\n\n")

	// Core Components
	builder.WriteString("Core Components:\n")
	if len(state.Selections.CoreComponents) > 0 {
		for _, component := range state.Selections.CoreComponents {
			builder.WriteString(fmt.Sprintf("  - %s\n", component))
		}
	} else {
		builder.WriteString("  (none selected)\n")
	}
	builder.WriteString("\n")

	// Plugins
	builder.WriteString("Plugins:\n")
	if len(state.Selections.Plugins) > 0 {
		for _, plugin := range state.Selections.Plugins {
			builder.WriteString(fmt.Sprintf("  - %s\n", plugin))
		}
	} else {
		builder.WriteString("  (none selected)\n")
	}
	builder.WriteString("\n")

	// Tools
	builder.WriteString("Development Tools:\n")
	if len(state.Selections.Tools) > 0 {
		for _, tool := range state.Selections.Tools {
			builder.WriteString(fmt.Sprintf("  - %s\n", tool))
		}
	} else {
		builder.WriteString("  (none selected)\n")
	}
	builder.WriteString("\n")

	// Project Config
	builder.WriteString(fmt.Sprintf("Include Project Config: %t\n", state.Selections.IncludeProjectConfig))

	return builder.String()
}

// ValidateSelections validates the wizard selections (core is assumed installed when entering config wizard)
func (ws *WizardService) ValidateSelections(state *WizardState) error {
	return nil
}

// GetTotalSteps returns the total number of steps
func (ws *WizardService) GetTotalSteps() int {
	return len(ws.steps)
}

// GetProgress returns the wizard progress percentage (0-100)
func (ws *WizardService) GetProgress(state *WizardState) int {
	if len(ws.steps) == 0 {
		return 0
	}

	return (state.CurrentStep * 100) / len(ws.steps)
}

// Reset resets the wizard to initial state
func (ws *WizardService) Reset(state *WizardState) {
	state.CurrentStep = 0
	state.Selections = interfaces.ConfigSelections{
		CoreComponents:       make([]string, 0),
		Plugins:              make([]string, 0),
		Tools:                make([]string, 0),
		IncludeProjectConfig: false,
		CustomAliases:        make(map[string]string),
		CustomFunctions:      make(map[string]string),
		CustomEnvVars:        make(map[string]string),
	}
	state.Completed = false
}

// CanProceed checks if can proceed to next step
func (ws *WizardService) CanProceed(state *WizardState) bool {
	step := ws.GetCurrentStep(state)
	if step == nil {
		return false
	}

	return true
}

// Complete marks the wizard as completed
func (ws *WizardService) Complete(state *WizardState) {
	state.Completed = true
}

// Helper functions
func sliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeFromSlice(slice []string, item string) []string {
	var result []string
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}
