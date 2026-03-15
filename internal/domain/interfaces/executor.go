package interfaces

import "github.com/JaimeJunr/Homestead/internal/domain/entities"

// ScriptExecutor defines the interface for script execution
type ScriptExecutor interface {
	// Execute executes a script
	Execute(script *entities.Script) error

	// CanExecute checks if a script can be executed
	CanExecute(script *entities.Script) bool

	// Validate validates a script before execution
	Validate(script *entities.Script) error
}
