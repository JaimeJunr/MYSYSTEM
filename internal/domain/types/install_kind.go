package types

// PackageInstallKind identifies how DefaultPackageInstaller runs a package.
// New kinds must be registered in internal/infrastructure/installer (see install_strategies.go).
type PackageInstallKind string

const (
	// InstallKindUtilityScript runs scripts/utilities/… via bash from HOMESTEAD_ROOT (category utilities).
	InstallKindUtilityScript PackageInstallKind = "utility_script"
	// InstallKindShellWithDownload downloads DownloadURL then runs InstallCmd (use {{download_path}} for the file).
	InstallKindShellWithDownload PackageInstallKind = "shell_with_download"
	// InstallKindShellLocal runs InstallCmd only (apt, curl|bash, etc.); working dir is a per-package temp folder.
	InstallKindShellLocal PackageInstallKind = "shell_local"
)

// AllPackageInstallKinds is the exhaustive list for tests and registry validation.
func AllPackageInstallKinds() []PackageInstallKind {
	return []PackageInstallKind{
		InstallKindUtilityScript,
		InstallKindShellWithDownload,
		InstallKindShellLocal,
	}
}

// IsValid reports whether k is a known install kind.
func (k PackageInstallKind) IsValid() bool {
	switch k {
	case InstallKindUtilityScript, InstallKindShellWithDownload, InstallKindShellLocal:
		return true
	default:
		return false
	}
}
