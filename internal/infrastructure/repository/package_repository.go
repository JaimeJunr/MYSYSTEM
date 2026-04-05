package repository

import (
	"sync"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/catalog"
)

// InMemoryPackageRepository is an in-memory implementation of PackageRepository
type InMemoryPackageRepository struct {
	packages map[string]*entities.Package
	mu       sync.RWMutex
}

// NewInMemoryPackageRepository creates a new in-memory package repository
// pre-filled from the embedded installer catalog JSON (see internal/infrastructure/catalog).
func NewInMemoryPackageRepository() interfaces.PackageRepository {
	repo := &InMemoryPackageRepository{
		packages: make(map[string]*entities.Package),
	}
	repo.bootstrapFromEmbeddedCatalog()
	return repo
}

func (r *InMemoryPackageRepository) bootstrapFromEmbeddedCatalog() {
	raw := catalog.EmbeddedCatalogBytes()
	if len(raw) == 0 {
		return
	}
	pkgs, _, err := catalog.ParseManifest(raw)
	if err != nil {
		return
	}
	for _, p := range pkgs {
		pkg := p
		_ = r.Save(&pkg)
	}
	for _, p := range defaultUtilityPackages() {
		pkg := p
		_ = r.Save(&pkg)
	}
}

// FindAll returns all packages
func (r *InMemoryPackageRepository) FindAll() ([]entities.Package, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	packages := make([]entities.Package, 0, len(r.packages))
	for _, pkg := range r.packages {
		packages = append(packages, *pkg)
	}

	return packages, nil
}

// FindByID finds a package by ID
func (r *InMemoryPackageRepository) FindByID(id string) (*entities.Package, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pkg, exists := r.packages[id]
	if !exists {
		return nil, types.ErrNotFound
	}

	pkgCopy := *pkg
	return &pkgCopy, nil
}

// FindByCategory finds packages by category
func (r *InMemoryPackageRepository) FindByCategory(category types.PackageCategory) ([]entities.Package, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	packages := make([]entities.Package, 0)
	for _, pkg := range r.packages {
		if pkg.Category == category {
			packages = append(packages, *pkg)
		}
	}

	return packages, nil
}

// Save saves a package
func (r *InMemoryPackageRepository) Save(pkg *entities.Package) error {
	if err := pkg.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	pkgCopy := *pkg
	r.packages[pkg.ID] = &pkgCopy

	return nil
}

// Delete deletes a package
func (r *InMemoryPackageRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.packages[id]; !exists {
		return types.ErrNotFound
	}

	delete(r.packages, id)
	return nil
}

// Exists checks if a package exists
func (r *InMemoryPackageRepository) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.packages[id]
	return exists
}
