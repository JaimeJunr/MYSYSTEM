package interfaces

import (
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// PackageRepository defines the interface for managing packages
type PackageRepository interface {
	FindAll() ([]entities.Package, error)
	FindByID(id string) (*entities.Package, error)
	FindByCategory(category types.PackageCategory) ([]entities.Package, error)
	Save(pkg *entities.Package) error
	Delete(id string) error
	Exists(id string) bool
}
