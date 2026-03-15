package services

import (
	"fmt"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// InstallerService orchestrates package installation operations
type InstallerService struct {
	repo      interfaces.PackageRepository
	installer interfaces.PackageInstaller
}

// NewInstallerService creates a new installer service
func NewInstallerService(
	repo interfaces.PackageRepository,
	installer interfaces.PackageInstaller,
) *InstallerService {
	return &InstallerService{
		repo:      repo,
		installer: installer,
	}
}

// GetAllPackages returns all available packages
func (s *InstallerService) GetAllPackages() ([]entities.Package, error) {
	packages, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("get all packages: %w", err)
	}
	return packages, nil
}

// GetPackagesByCategory returns packages filtered by category
func (s *InstallerService) GetPackagesByCategory(category types.PackageCategory) ([]entities.Package, error) {
	packages, err := s.repo.FindByCategory(category)
	if err != nil {
		return nil, fmt.Errorf("get packages by category %s: %w", category, err)
	}
	return packages, nil
}

// GetPackagesByCategories returns packages from multiple categories (combined, no duplicates by ID)
func (s *InstallerService) GetPackagesByCategories(categories []types.PackageCategory) ([]entities.Package, error) {
	seen := make(map[string]bool)
	var result []entities.Package
	for _, cat := range categories {
		packages, err := s.repo.FindByCategory(cat)
		if err != nil {
			return nil, fmt.Errorf("get packages by category %s: %w", cat, err)
		}
		for _, pkg := range packages {
			if !seen[pkg.ID] {
				seen[pkg.ID] = true
				result = append(result, pkg)
			}
		}
	}
	return result, nil
}

// GetPackageByID returns a package by ID
func (s *InstallerService) GetPackageByID(id string) (*entities.Package, error) {
	pkg, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("get package %s: %w", id, err)
	}
	return pkg, nil
}

// InstallPackage installs a package with progress reporting
func (s *InstallerService) InstallPackage(id string, progressCallback interfaces.ProgressCallback) error {
	// Get package
	pkg, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("install package %s: %w", id, err)
	}

	// Validate package
	if err := pkg.Validate(); err != nil {
		return fmt.Errorf("install package %s: invalid package: %w", id, err)
	}

	// Check if can install
	if !s.installer.CanInstall(pkg) {
		return fmt.Errorf("install package %s: system cannot install this package", id)
	}

	// Install with progress
	if err := s.installer.Install(pkg, progressCallback); err != nil {
		return fmt.Errorf("install package %s: %w", id, err)
	}

	return nil
}

// IsPackageInstalled checks if a package is already installed
func (s *InstallerService) IsPackageInstalled(id string) (bool, error) {
	pkg, err := s.repo.FindByID(id)
	if err != nil {
		return false, fmt.Errorf("check if package %s installed: %w", id, err)
	}

	installed, err := s.installer.IsInstalled(pkg)
	if err != nil {
		return false, fmt.Errorf("check if package %s installed: %w", id, err)
	}

	return installed, nil
}

// UninstallPackage uninstalls a package
func (s *InstallerService) UninstallPackage(id string) error {
	pkg, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("uninstall package %s: %w", id, err)
	}

	if err := s.installer.Uninstall(pkg); err != nil {
		return fmt.Errorf("uninstall package %s: %w", id, err)
	}

	return nil
}
