package installer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/executor"
)

// DefaultPackageInstaller is the default implementation of PackageInstaller
type DefaultPackageInstaller struct {
	tempDir    string
	rootDir    string
	ctx        context.Context
	cancel     context.CancelFunc
	strategies map[types.PackageInstallKind]installStrategy
}

// NewDefaultPackageInstaller creates a new default package installer (Homestead root = cwd).
func NewDefaultPackageInstaller() interfaces.PackageInstaller {
	return NewDefaultPackageInstallerWithRoot("")
}

// NewDefaultPackageInstallerWithRoot resolves HOMESTEAD_ROOT from scriptRoot like BashExecutor.
func NewDefaultPackageInstallerWithRoot(scriptRoot string) interfaces.PackageInstaller {
	ctx, cancel := context.WithCancel(context.Background())
	rootDir, err := executor.ResolveScriptRoot(scriptRoot)
	if err != nil || rootDir == "" {
		rootDir, _ = os.Getwd()
		if rootDir == "" {
			rootDir = "."
		}
	}
	return &DefaultPackageInstaller{
		tempDir:    os.TempDir(),
		rootDir:    rootDir,
		ctx:        ctx,
		cancel:     cancel,
		strategies: newStrategyMap(defaultStrategies()),
	}
}

// SetHomesteadRoot updates HOMESTEAD_ROOT for utility script installs.
func (i *DefaultPackageInstaller) SetHomesteadRoot(dir string) error {
	r, err := executor.ResolveScriptRoot(dir)
	if err != nil {
		return err
	}
	i.rootDir = r
	return nil
}

// Install installs a package with progress reporting
func (i *DefaultPackageInstaller) Install(pkg *entities.Package, progressCallback interfaces.ProgressCallback) error {
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

	kind := pkg.ResolveInstallKind()
	s, ok := i.strategies[kind]
	if !ok || s == nil {
		return fmt.Errorf("no installer strategy for kind %q", kind)
	}
	return s.install(i, pkg, progressCallback)
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

func (i *DefaultPackageInstaller) installHomesteadUtility(pkg *entities.Package) error {
	rel := strings.TrimSpace(pkg.UtilityScriptPath)
	scriptPath := filepath.Join(i.rootDir, rel)
	if _, err := os.Stat(scriptPath); err != nil {
		return fmt.Errorf("script não encontrado em %s: %w", scriptPath, err)
	}
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("utilizador atual: %w", err)
	}
	var cmd *exec.Cmd
	if pkg.RequiresSudo {
		cmd = exec.Command("sudo", "-E", "bash", scriptPath)
	} else {
		cmd = exec.Command("bash", scriptPath)
	}
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("REAL_USER=%s", currentUser.Username),
		fmt.Sprintf("REAL_HOME=%s", currentUser.HomeDir),
		fmt.Sprintf("HOMESTEAD_ROOT=%s", i.rootDir),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
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
	kind := pkg.ResolveInstallKind()
	s, ok := i.strategies[kind]
	if !ok || s == nil {
		return false
	}
	return s.canInstall(i, pkg)
}
