package entities

import (
	"testing"

	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

func TestPackageValidate(t *testing.T) {
	tests := []struct {
		name    string
		pkg     Package
		wantErr bool
	}{
		{
			name: "Valid package",
			pkg: Package{
				ID:          "test",
				Name:        "Test Package",
				Description: "Test",
				Version:     "1.0",
				Category:    types.PackageCategoryIDE,
				DownloadURL: "https://example.com/pkg",
			},
			wantErr: false,
		},
		{
			name: "Missing ID",
			pkg: Package{
				Name:        "Test",
				Description: "Test",
				Category:    types.PackageCategoryIDE,
				DownloadURL: "https://example.com/pkg",
			},
			wantErr: true,
		},
		{
			name: "Missing Name",
			pkg: Package{
				ID:          "test",
				Description: "Test",
				Category:    types.PackageCategoryIDE,
				DownloadURL: "https://example.com/pkg",
			},
			wantErr: true,
		},
		{
			name: "Invalid category",
			pkg: Package{
				ID:          "test",
				Name:        "Test",
				Description: "Test",
				Category:    types.PackageCategory("invalid"),
				DownloadURL: "https://example.com/pkg",
			},
			wantErr: true,
		},
		{
			name: "Package with InstallCmd only",
			pkg: Package{
				ID:          "test",
				Name:        "Test",
				Description: "Test",
				Category:    types.PackageCategoryIDE,
				InstallCmd:  "apt install test",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pkg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Package.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPackageIsIDE(t *testing.T) {
	pkg := Package{
		Category: types.PackageCategoryIDE,
	}

	if !pkg.IsIDE() {
		t.Error("Expected IsIDE() to return true")
	}
}

func TestPackageIsTool(t *testing.T) {
	pkg := Package{
		Category: types.PackageCategoryTool,
	}

	if !pkg.IsTool() {
		t.Error("Expected IsTool() to return true")
	}
}
