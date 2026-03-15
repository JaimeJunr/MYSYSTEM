package repository

import (
	"testing"

	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

func TestPackageRepository_FindAll(t *testing.T) {
	repo := NewInMemoryPackageRepository()

	packages, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if len(packages) == 0 {
		t.Error("Expected at least one package")
	}

	// Check for expected packages
	expectedPackages := []string{"claude-code", "cursor", "antigravity"}
	found := make(map[string]bool)

	for _, pkg := range packages {
		found[pkg.ID] = true
	}

	for _, expected := range expectedPackages {
		if !found[expected] {
			t.Errorf("Expected package %s not found", expected)
		}
	}
}

func TestPackageRepository_FindByID(t *testing.T) {
	repo := NewInMemoryPackageRepository()

	// Test finding existing package
	pkg, err := repo.FindByID("claude-code")
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if pkg.ID != "claude-code" {
		t.Errorf("Expected ID claude-code, got %s", pkg.ID)
	}

	if pkg.Name != "Claude Code CLI" {
		t.Errorf("Expected name 'Claude Code CLI', got %s", pkg.Name)
	}

	// Test finding non-existent package
	_, err = repo.FindByID("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent package")
	}
}

func TestPackageRepository_FindByCategory(t *testing.T) {
	repo := NewInMemoryPackageRepository()

	packages, err := repo.FindByCategory(types.PackageCategoryIDE)
	if err != nil {
		t.Fatalf("FindByCategory() error = %v", err)
	}

	if len(packages) == 0 {
		t.Error("Expected at least one IDE package")
	}

	// All packages should be IDEs
	for _, pkg := range packages {
		if pkg.Category != types.PackageCategoryIDE {
			t.Errorf("Expected category IDE, got %s", pkg.Category)
		}
	}
}

func TestPackageRepository_Exists(t *testing.T) {
	repo := NewInMemoryPackageRepository()

	if !repo.Exists("claude-code") {
		t.Error("Expected claude-code to exist")
	}

	if repo.Exists("non-existent") {
		t.Error("Expected non-existent package to not exist")
	}
}
