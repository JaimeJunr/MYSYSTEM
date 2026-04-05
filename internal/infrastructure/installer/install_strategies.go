package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// installStrategy is one PackageInstaller execution path (strategy registry).
type installStrategy interface {
	kind() types.PackageInstallKind
	install(i *DefaultPackageInstaller, pkg *entities.Package, progressCallback interfaces.ProgressCallback) error
	canInstall(i *DefaultPackageInstaller, pkg *entities.Package) bool
}

func defaultStrategies() []installStrategy {
	return []installStrategy{
		utilityScriptStrategy{},
		shellWithDownloadStrategy{},
		shellLocalStrategy{},
	}
}

func newStrategyMap(strategies []installStrategy) map[types.PackageInstallKind]installStrategy {
	m := make(map[types.PackageInstallKind]installStrategy, len(strategies))
	for _, s := range strategies {
		k := s.kind()
		if _, dup := m[k]; dup {
			panic(fmt.Sprintf("installer: duplicate strategy for kind %q", k))
		}
		m[k] = s
	}
	for _, k := range types.AllPackageInstallKinds() {
		if m[k] == nil {
			panic(fmt.Sprintf("installer: missing strategy for kind %q", k))
		}
	}
	return m
}

// expandInstallCmd replaces documented placeholders. Catalog entries must use {{download_path}} for downloaded artifacts.
func expandInstallCmd(installCmd, downloadPath string) string {
	return strings.ReplaceAll(installCmd, "{{download_path}}", downloadPath)
}

func runShellInstallCmd(pkg *entities.Package, downloadPath string) error {
	if pkg.InstallCmd == "" {
		return nil
	}
	cmdLine := expandInstallCmd(pkg.InstallCmd, downloadPath)
	cmd := exec.Command("bash", "-c", cmdLine)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if fi, err := os.Stat(downloadPath); err == nil && fi.IsDir() {
		cmd.Dir = downloadPath
	} else {
		cmd.Dir = filepath.Dir(downloadPath)
	}
	return cmd.Run()
}

type utilityScriptStrategy struct{}

func (utilityScriptStrategy) kind() types.PackageInstallKind { return types.InstallKindUtilityScript }

func (utilityScriptStrategy) install(i *DefaultPackageInstaller, pkg *entities.Package, cb interfaces.ProgressCallback) error {
	cb(interfaces.InstallProgress{
		Package:  pkg,
		Status:   "preparing",
		Progress: 20,
		Message:  "Preparando script Homestead...",
		CanAbort: false,
	})
	cb(interfaces.InstallProgress{
		Package:  pkg,
		Status:   "installing",
		Progress: 70,
		Message:  "Instalando...",
		CanAbort: false,
	})
	if err := i.installHomesteadUtility(pkg); err != nil {
		cb(interfaces.InstallProgress{
			Package:     pkg,
			Status:      "failed",
			Progress:    70,
			Message:     "Erro ao instalar",
			Error:       err,
			IsCompleted: true,
		})
		return fmt.Errorf("installation failed: %w", err)
	}
	cb(interfaces.InstallProgress{
		Package:     pkg,
		Status:      "complete",
		Progress:    100,
		Message:     "Instalação concluída com sucesso!",
		IsCompleted: true,
	})
	return nil
}

func (utilityScriptStrategy) canInstall(i *DefaultPackageInstaller, pkg *entities.Package) bool {
	if _, err := exec.LookPath("bash"); err != nil {
		return false
	}
	if pkg.RequiresSudo {
		if _, err := exec.LookPath("sudo"); err != nil {
			return false
		}
	}
	p := filepath.Join(i.rootDir, strings.TrimSpace(pkg.UtilityScriptPath))
	_, err := os.Stat(p)
	return err == nil
}

type shellWithDownloadStrategy struct{}

func (shellWithDownloadStrategy) kind() types.PackageInstallKind { return types.InstallKindShellWithDownload }

func (shellWithDownloadStrategy) install(i *DefaultPackageInstaller, pkg *entities.Package, cb interfaces.ProgressCallback) error {
	cb(interfaces.InstallProgress{
		Package:  pkg,
		Status:   "downloading",
		Progress: 0,
		Message:  "Iniciando download...",
		CanAbort: true,
	})
	downloadPath, err := i.downloadPackage(pkg, cb)
	if err != nil {
		cb(interfaces.InstallProgress{
			Package:     pkg,
			Status:      "failed",
			Progress:    0,
			Message:     "Erro ao baixar",
			Error:       err,
			IsCompleted: true,
		})
		return fmt.Errorf("download failed: %w", err)
	}
	cb(interfaces.InstallProgress{
		Package:  pkg,
		Status:   "installing",
		Progress: 70,
		Message:  "Instalando...",
		CanAbort: false,
	})
	if err := runShellInstallCmd(pkg, downloadPath); err != nil {
		cb(interfaces.InstallProgress{
			Package:     pkg,
			Status:      "failed",
			Progress:    70,
			Message:     "Erro ao instalar",
			Error:       err,
			IsCompleted: true,
		})
		return fmt.Errorf("installation failed: %w", err)
	}
	cb(interfaces.InstallProgress{
		Package:     pkg,
		Status:      "complete",
		Progress:    100,
		Message:     "Instalação concluída com sucesso!",
		IsCompleted: true,
	})
	return nil
}

func (shellWithDownloadStrategy) canInstall(*DefaultPackageInstaller, *entities.Package) bool {
	if _, err := exec.LookPath("wget"); err == nil {
		return true
	}
	if _, err := exec.LookPath("curl"); err == nil {
		return true
	}
	return true
}

type shellLocalStrategy struct{}

func (shellLocalStrategy) kind() types.PackageInstallKind { return types.InstallKindShellLocal }

func (shellLocalStrategy) install(i *DefaultPackageInstaller, pkg *entities.Package, cb interfaces.ProgressCallback) error {
	downloadPath := filepath.Join(i.tempDir, pkg.ID)
	if err := os.MkdirAll(downloadPath, 0750); err != nil {
		cb(interfaces.InstallProgress{
			Package:     pkg,
			Status:      "failed",
			Progress:    0,
			Message:     "Erro ao preparar instalação",
			Error:       err,
			IsCompleted: true,
		})
		return fmt.Errorf("prepare install dir: %w", err)
	}
	cb(interfaces.InstallProgress{
		Package:  pkg,
		Status:   "downloading",
		Progress: 60,
		Message:  "Preparando instalação...",
		CanAbort: false,
	})
	cb(interfaces.InstallProgress{
		Package:  pkg,
		Status:   "installing",
		Progress: 70,
		Message:  "Instalando...",
		CanAbort: false,
	})
	if err := runShellInstallCmd(pkg, downloadPath); err != nil {
		cb(interfaces.InstallProgress{
			Package:     pkg,
			Status:      "failed",
			Progress:    70,
			Message:     "Erro ao instalar",
			Error:       err,
			IsCompleted: true,
		})
		return fmt.Errorf("installation failed: %w", err)
	}
	cb(interfaces.InstallProgress{
		Package:     pkg,
		Status:      "complete",
		Progress:    100,
		Message:     "Instalação concluída com sucesso!",
		IsCompleted: true,
	})
	return nil
}

func (shellLocalStrategy) canInstall(*DefaultPackageInstaller, *entities.Package) bool {
	return true
}
