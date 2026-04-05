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
		{
			name: "Utilities package with script path and project URL",
			pkg: Package{
				ID:                "util-x",
				Name:              "X",
				Description:       "D",
				Version:           "latest",
				Category:          types.PackageCategoryUtilities,
				UtilityScriptPath: "scripts/utilities/x.sh",
				ProjectURL:        "https://example.com",
				RequiresSudo:      true,
			},
			wantErr: false,
		},
		{
			name: "Utilities missing project URL",
			pkg: Package{
				ID:                "util-x",
				Name:              "X",
				Description:       "D",
				Category:          types.PackageCategoryUtilities,
				UtilityScriptPath: "scripts/utilities/x.sh",
			},
			wantErr: true,
		},
		{
			name: "Utility script path on non-utilities category",
			pkg: Package{
				ID:                "bad",
				Name:              "B",
				Description:       "D",
				Category:          types.PackageCategoryApp,
				UtilityScriptPath: "scripts/utilities/x.sh",
				ProjectURL:        "https://example.com",
			},
			wantErr: true,
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

func TestPackageResolveInstallKind(t *testing.T) {
	tests := []struct {
		name string
		pkg  Package
		want types.PackageInstallKind
	}{
		{
			name: "utility script",
			pkg: Package{
				UtilityScriptPath: "scripts/x.sh",
			},
			want: types.InstallKindUtilityScript,
		},
		{
			name: "download URL",
			pkg: Package{
				DownloadURL: "https://example.com/a",
				InstallCmd:  "bash {{download_path}}",
			},
			want: types.InstallKindShellWithDownload,
		},
		{
			name: "local shell only",
			pkg: Package{
				InstallCmd: "apt install x",
			},
			want: types.InstallKindShellLocal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pkg.ResolveInstallKind()
			if got != tt.want {
				t.Errorf("ResolveInstallKind() = %q, want %q", got, tt.want)
			}
			if !got.IsValid() {
				t.Errorf("ResolveInstallKind() = %q not valid", got)
			}
		})
	}
}
