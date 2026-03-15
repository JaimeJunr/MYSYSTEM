package interfaces

import "github.com/JaimeJunr/Homestead/internal/domain/entities"

// InstallProgress represents the progress of an installation
type InstallProgress struct {
	Package     *entities.Package
	Status      string // "downloading", "installing", "complete", "failed"
	Progress    int    // 0-100
	Message     string
	Error       error
	CanAbort    bool
	IsAborted   bool
	IsCompleted bool
}

// ProgressCallback is called to report installation progress
type ProgressCallback func(progress InstallProgress)

// PackageInstaller defines the interface for installing packages
type PackageInstaller interface {
	// Install installs a package with progress reporting
	Install(pkg *entities.Package, progressCallback ProgressCallback) error

	// IsInstalled checks if a package is already installed
	IsInstalled(pkg *entities.Package) (bool, error)

	// Uninstall removes a package
	Uninstall(pkg *entities.Package) error

	// CanInstall checks if the system can install this package
	CanInstall(pkg *entities.Package) bool
}
