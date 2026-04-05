package services

import (
	"fmt"
	"os/exec"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// ScriptService handles script-related operations
type ScriptService struct {
	repo     interfaces.ScriptRepository
	executor interfaces.ScriptExecutor
}

// NewScriptService creates a new script service
func NewScriptService(
	repo interfaces.ScriptRepository,
	executor interfaces.ScriptExecutor,
) *ScriptService {
	return &ScriptService{
		repo:     repo,
		executor: executor,
	}
}

// GetAllScripts returns all scripts
func (s *ScriptService) GetAllScripts() ([]entities.Script, error) {
	return s.repo.FindAll()
}

// GetScriptByID returns a script by ID
func (s *ScriptService) GetScriptByID(id string) (*entities.Script, error) {
	if id == "" {
		return nil, fmt.Errorf("get script: %w", types.ErrInvalidInput)
	}

	return s.repo.FindByID(id)
}

// GetScriptsByCategory returns scripts by category
func (s *ScriptService) GetScriptsByCategory(category types.Category) ([]entities.Script, error) {
	if !category.IsValid() {
		return nil, fmt.Errorf("get scripts by category: invalid category %s: %w",
			category, types.ErrInvalidInput)
	}

	return s.repo.FindByCategory(category)
}

// ExecuteScript executes a script by ID
func (s *ScriptService) ExecuteScript(id string) error {
	if id == "" {
		return fmt.Errorf("execute script: %w", types.ErrInvalidInput)
	}

	// Find script
	script, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("execute script %s: %w", id, err)
	}

	// Validate before execution
	if err := s.executor.Validate(script); err != nil {
		return fmt.Errorf("execute script %s: %w", id, err)
	}

	// Execute
	if err := s.executor.Execute(script, interfaces.ScriptExecOpts{}); err != nil {
		return fmt.Errorf("execute script %s: %w", id, err)
	}

	return nil
}

// ExecuteScriptCapture runs the script and returns combined stdout/stderr (for in-app TUI).
func (s *ScriptService) ExecuteScriptCapture(id string, opts interfaces.ScriptExecOpts) (output string, err error) {
	if id == "" {
		return "", fmt.Errorf("execute script: %w", types.ErrInvalidInput)
	}

	script, err := s.repo.FindByID(id)
	if err != nil {
		return "", fmt.Errorf("execute script %s: %w", id, err)
	}

	out, err := s.executor.ExecuteCapture(script, opts)
	if err != nil {
		return out, fmt.Errorf("execute script %s: %w", id, err)
	}
	return out, nil
}

// ScriptInteractiveCommand builds an exec.Cmd for tea.ExecProcess (sudo / TTY).
func (s *ScriptService) ScriptInteractiveCommand(id string, opts interfaces.ScriptExecOpts) (*exec.Cmd, error) {
	if id == "" {
		return nil, fmt.Errorf("execute script: %w", types.ErrInvalidInput)
	}

	script, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("execute script %s: %w", id, err)
	}

	cmd, err := s.executor.InteractiveCommand(script, opts)
	if err != nil {
		return nil, fmt.Errorf("execute script %s: %w", id, err)
	}
	return cmd, nil
}

// CanExecuteScript checks if a script can be executed
func (s *ScriptService) CanExecuteScript(id string) bool {
	script, err := s.repo.FindByID(id)
	if err != nil {
		return false
	}

	return s.executor.CanExecute(script)
}

// ScriptExists checks if a script exists
func (s *ScriptService) ScriptExists(id string) bool {
	return s.repo.Exists(id)
}

func (s *ScriptService) ConfigureScriptRoot(dir string) error {
	return s.executor.SetScriptRoot(dir)
}
