package services

import (
	"testing"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/installer"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/repository"
)

func TestInstallerService_GetAllPackages(t *testing.T) {
	repo := repository.NewInMemoryPackageRepository()
	inst := installer.NewDefaultPackageInstaller()
	service := NewInstallerService(repo, inst)

	packages, err := service.GetAllPackages()
	if err != nil {
		t.Fatalf("GetAllPackages() error = %v", err)
	}

	if len(packages) == 0 {
		t.Error("Expected at least one package")
	}
}

func TestInstallerService_GetPackagesByCategory(t *testing.T) {
	repo := repository.NewInMemoryPackageRepository()
	inst := installer.NewDefaultPackageInstaller()
	service := NewInstallerService(repo, inst)

	packages, err := service.GetPackagesByCategory(types.PackageCategoryIDE)
	if err != nil {
		t.Fatalf("GetPackagesByCategory() error = %v", err)
	}

	if len(packages) == 0 {
		t.Error("Expected at least one IDE package")
	}

	for _, pkg := range packages {
		if pkg.Category != types.PackageCategoryIDE {
			t.Errorf("Expected category IDE, got %s", pkg.Category)
		}
	}
}

func TestInstallerService_GetPackageByID(t *testing.T) {
	repo := repository.NewInMemoryPackageRepository()
	inst := installer.NewDefaultPackageInstaller()
	service := NewInstallerService(repo, inst)

	pkg, err := service.GetPackageByID("claude-code")
	if err != nil {
		t.Fatalf("GetPackageByID() error = %v", err)
	}

	if pkg.ID != "claude-code" {
		t.Errorf("Expected ID claude-code, got %s", pkg.ID)
	}

	// Test non-existent package
	_, err = service.GetPackageByID("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent package")
	}
}

func TestInstallerService_IsPackageInstalled(t *testing.T) {
	repo := repository.NewInMemoryPackageRepository()
	inst := installer.NewDefaultPackageInstaller()
	service := NewInstallerService(repo, inst)

	// This will check if the package is installed on the system
	// We don't know if it actually is, so just verify no error
	_, err := service.IsPackageInstalled("claude-code")
	if err != nil {
		t.Fatalf("IsPackageInstalled() error = %v", err)
	}
}

func TestInstallerService_MergePackages(t *testing.T) {
	repo := repository.NewInMemoryPackageRepository()
	inst := installer.NewDefaultPackageInstaller()
	service := NewInstallerService(repo, inst)

	newPkg := entities.Package{
		ID:          "remote-only",
		Name:        "Remote Only",
		Description: "d",
		Version:     "1",
		Category:    types.PackageCategoryTool,
		InstallCmd:  "true",
		CheckCmd:    "true",
	}
	if err := service.MergePackages([]entities.Package{newPkg}); err != nil {
		t.Fatal(err)
	}
	p, err := service.GetPackageByID("remote-only")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Remote Only" {
		t.Fatalf("name = %s", p.Name)
	}
}
