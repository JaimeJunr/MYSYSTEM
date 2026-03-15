package installer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
)

// DefaultPackageInstaller is the default implementation of PackageInstaller
type DefaultPackageInstaller struct {
	tempDir string
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewDefaultPackageInstaller creates a new default package installer
func NewDefaultPackageInstaller() interfaces.PackageInstaller {
	ctx, cancel := context.WithCancel(context.Background())
	return &DefaultPackageInstaller{
		tempDir: os.TempDir(),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Install installs a package with progress reporting
func (i *DefaultPackageInstaller) Install(pkg *entities.Package, progressCallback interfaces.ProgressCallback) error {
	// Check if already installed
	installed, err := i.IsInstalled(pkg)
	if err == nil && installed {
		progressCallback(interfaces.InstallProgress{
			Package:     pkg,
			Status:      "complete",
			Progress:    100,
			Message:     "Já instalado",
			IsCompleted: true,
		})
		return nil
	}

	// Download phase
	progressCallback(interfaces.InstallProgress{
		Package:  pkg,
		Status:   "downloading",
		Progress: 0,
		Message:  "Iniciando download...",
		CanAbort: true,
	})

	downloadPath, err := i.downloadPackage(pkg, progressCallback)
	if err != nil {
		progressCallback(interfaces.InstallProgress{
			Package:     pkg,
			Status:      "failed",
			Progress:    0,
			Message:     "Erro ao baixar",
			Error:       err,
			IsCompleted: true,
		})
		return fmt.Errorf("download failed: %w", err)
	}

	// Installation phase
	progressCallback(interfaces.InstallProgress{
		Package:  pkg,
		Status:   "installing",
		Progress: 70,
		Message:  "Instalando...",
		CanAbort: false,
	})

	if err := i.installPackage(pkg, downloadPath); err != nil {
		progressCallback(interfaces.InstallProgress{
			Package:     pkg,
			Status:      "failed",
			Progress:    70,
			Message:     "Erro ao instalar",
			Error:       err,
			IsCompleted: true,
		})
		return fmt.Errorf("installation failed: %w", err)
	}

	// Complete
	progressCallback(interfaces.InstallProgress{
		Package:     pkg,
		Status:      "complete",
		Progress:    100,
		Message:     "Instalação concluída com sucesso!",
		IsCompleted: true,
	})

	return nil
}

// downloadPackage downloads a package and reports progress
func (i *DefaultPackageInstaller) downloadPackage(pkg *entities.Package, progressCallback interfaces.ProgressCallback) (string, error) {
	// Create HTTP request with context for cancellation
	req, err := http.NewRequestWithContext(i.ctx, "GET", pkg.DownloadURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	// Determine filename
	filename := filepath.Base(pkg.DownloadURL)
	if filename == "" || filename == "/" {
		filename = pkg.ID + ".download"
	}

	// Create temp file
	tempFile := filepath.Join(i.tempDir, filename)
	out, err := os.Create(tempFile)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Download with progress
	totalSize := resp.ContentLength
	downloaded := int64(0)
	buffer := make([]byte, 32*1024) // 32KB buffer

	lastUpdate := time.Now()
	for {
		// Check for context cancellation
		select {
		case <-i.ctx.Done():
			return "", fmt.Errorf("download cancelled")
		default:
		}

		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, writeErr := out.Write(buffer[:n])
			if writeErr != nil {
				return "", writeErr
			}
			downloaded += int64(n)

			// Update progress every 100ms
			if time.Since(lastUpdate) > 100*time.Millisecond {
				progress := 0
				if totalSize > 0 {
					progress = int((float64(downloaded) / float64(totalSize)) * 60) // 0-60% for download
				}

				progressCallback(interfaces.InstallProgress{
					Package:  pkg,
					Status:   "downloading",
					Progress: progress,
					Message:  fmt.Sprintf("Baixando... %d/%d bytes", downloaded, totalSize),
					CanAbort: true,
				})
				lastUpdate = time.Now()
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}

	progressCallback(interfaces.InstallProgress{
		Package:  pkg,
		Status:   "downloading",
		Progress: 60,
		Message:  "Download concluído",
		CanAbort: false,
	})

	return tempFile, nil
}

// installPackage installs a downloaded package
func (i *DefaultPackageInstaller) installPackage(pkg *entities.Package, downloadPath string) error {
	if pkg.InstallCmd == "" {
		return nil
	}

	// Replace placeholders in install command
	installCmd := strings.ReplaceAll(pkg.InstallCmd, "install.sh", downloadPath)
	installCmd = strings.ReplaceAll(installCmd, "cursor.AppImage", downloadPath)
	installCmd = strings.ReplaceAll(installCmd, "antigravity.deb", downloadPath)

	// Execute install command
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Dir(downloadPath)

	return cmd.Run()
}

// IsInstalled checks if a package is already installed
func (i *DefaultPackageInstaller) IsInstalled(pkg *entities.Package) (bool, error) {
	if pkg.CheckCmd == "" {
		return false, nil
	}

	cmd := exec.Command("bash", "-c", pkg.CheckCmd)
	err := cmd.Run()
	return err == nil, nil
}

// Uninstall removes a package
func (i *DefaultPackageInstaller) Uninstall(pkg *entities.Package) error {
	// For now, just a placeholder
	// Real implementation would depend on package type
	return fmt.Errorf("uninstall not implemented for %s", pkg.ID)
}

// CanInstall checks if the system can install this package
func (i *DefaultPackageInstaller) CanInstall(pkg *entities.Package) bool {
	// Check if we have wget or curl for downloading
	if _, err := exec.LookPath("wget"); err == nil {
		return true
	}
	if _, err := exec.LookPath("curl"); err == nil {
		return true
	}
	return true // We use Go's http.Get, so always true
}
