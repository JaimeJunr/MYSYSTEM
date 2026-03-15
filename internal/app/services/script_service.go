package services

import (
	"fmt"

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
	if err := s.executor.Execute(script); err != nil {
		return fmt.Errorf("execute script %s: %w", id, err)
	}

	return nil
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
