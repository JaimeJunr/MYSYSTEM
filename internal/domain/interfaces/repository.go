package interfaces

import (
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// ScriptRepository defines the interface for script data access
type ScriptRepository interface {
	// FindAll returns all scripts
	FindAll() ([]entities.Script, error)

	// FindByID returns a script by its ID
	FindByID(id string) (*entities.Script, error)

	// FindByCategory returns all scripts in a category
	FindByCategory(category types.Category) ([]entities.Script, error)

	// Save saves a script
	Save(script *entities.Script) error

	// Delete deletes a script by ID
	Delete(id string) error

	// Exists checks if a script exists
	Exists(id string) bool
}
