package cmds

import (
	"context"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/catalog"
	"github.com/JaimeJunr/Homestead/internal/tui/msg"
	"github.com/JaimeJunr/Homestead/internal/tui/sysurl"
)

func CheckZshCoreInstalled(installerService *services.InstallerService) tea.Cmd {
	return func() tea.Msg {
		installed, _ := installerService.IsPackageInstalled("oh-my-zsh")
		return msg.ZshCoreInstalled{Installed: installed}
	}
}

func FetchCatalog(url string, svc *services.InstallerService) tea.Cmd {
	if strings.TrimSpace(url) == "" {
		return nil
	}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		body, err := catalog.Fetch(ctx, url)
		if err != nil {
			return msg.CatalogFetched{Err: err}
		}
		pkgs, _, err := catalog.ParseManifest(body)
		if err != nil {
			return msg.CatalogFetched{Err: err}
		}
		if err := svc.MergePackages(pkgs); err != nil {
			return msg.CatalogFetched{Err: err}
		}
		_ = catalog.WriteCache(body)
		return msg.CatalogFetched{Ok: true}
	}
}

func RunScriptCapture(service *services.ScriptService, scriptID string, opts interfaces.ScriptExecOpts) tea.Cmd {
	return func() tea.Msg {
		out, err := service.ExecuteScriptCapture(scriptID, opts)
		return msg.ScriptCaptured{Output: out, Err: err}
	}
}

func InstallPackage(service *services.InstallerService, packageID string) tea.Cmd {
	return func() tea.Msg {
		progressChan := make(chan interfaces.InstallProgress, 10)

		go func() {
			err := service.InstallPackage(packageID, func(progress interfaces.InstallProgress) {
				progressChan <- progress
			})
			if err != nil {
				progressChan <- interfaces.InstallProgress{
					Status:      "failed",
					Message:     err.Error(),
					IsCompleted: true,
					Error:       err,
				}
			}
			close(progressChan)
		}()

		for progress := range progressChan {
			return msg.Progress(progress)
		}

		return msg.InstallComplete{Err: nil}
	}
}

func ApplyZshConfig(configService *services.ConfigService, selections interfaces.ConfigSelections) tea.Cmd {
	return func() tea.Msg {
		err := configService.ApplyConfig(selections)
		return msg.ZshApplyResult{Err: err}
	}
}

func OpenURL(url string) tea.Cmd {
	return func() tea.Msg {
		err := sysurl.Open(url)
		return msg.URLActionDone{Verb: "open", Err: err}
	}
}

func CopyURL(url string) tea.Cmd {
	return func() tea.Msg {
		err := sysurl.CopyToClipboard(url)
		return msg.URLActionDone{Verb: "copy", Err: err}
	}
}
