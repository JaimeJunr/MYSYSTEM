package homesteadcli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/JaimeJunr/Homestead/internal/infrastructure/preferences"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/profilestate"
)

// RunExportProfile parses argv like: export-profile [-format json|text] [-o path]
func RunExportProfile(argv []string, appVersion string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("export-profile", flag.ContinueOnError)
	fs.SetOutput(stderr)
	format := fs.String("format", "json", "json ou text")
	outPath := fs.String("o", "", "ficheiro de saída (vazio = stdout)")
	if err := fs.Parse(argv); err != nil {
		return 2
	}

	prefsPath, err := preferences.DefaultPath()
	if err != nil {
		fmt.Fprintf(stderr, "Homestead: preferências: %v\n", err)
		return 1
	}
	prefs, err := preferences.Load(prefsPath)
	if err != nil {
		fmt.Fprintf(stderr, "Homestead: carregar preferências: %v\n", err)
		return 1
	}
	profPath, err := profilestate.DefaultPath()
	if err != nil {
		fmt.Fprintf(stderr, "Homestead: perfil: %v\n", err)
		return 1
	}
	state, err := profilestate.Load(profPath)
	if err != nil {
		fmt.Fprintf(stderr, "Homestead: carregar perfil: %v\n", err)
		return 1
	}

	var out io.Writer = stdout
	if strings.TrimSpace(*outPath) != "" {
		f, err := os.Create(*outPath)
		if err != nil {
			fmt.Fprintf(stderr, "Homestead: criar ficheiro: %v\n", err)
			return 1
		}
		defer f.Close()
		out = f
	}

	if err := profilestate.WriteExport(out, *format, state, prefs, appVersion); err != nil {
		fmt.Fprintf(stderr, "Homestead: exportar: %v\n", err)
		return 1
	}
	return 0
}

// RunShellInit prints minimal env snippets for bash or fish (paths under XDG).
func RunShellInit(argv []string, stderr io.Writer) int {
	if len(argv) != 1 {
		fmt.Fprintf(stderr, "uso: homestead shell-init bash|fish\n")
		return 2
	}
	shell := strings.ToLower(strings.TrimSpace(argv[0]))
	switch shell {
	case "bash":
		fmt.Print(`# Homestead — variáveis úteis para migração (gerado por: homestead shell-init bash)
: "${XDG_CONFIG_HOME:=$HOME/.config}"
export HOMESTEAD_CONFIG_DIR="$XDG_CONFIG_HOME/homestead"
# Preferências: $HOMESTEAD_CONFIG_DIR/preferences.yaml
# Perfil (instalados/favoritos): $HOMESTEAD_CONFIG_DIR/profile.yaml
`)
	case "fish":
		fmt.Print(`# Homestead — variáveis úteis para migração (gerado por: homestead shell-init fish)
if test -z "$XDG_CONFIG_HOME"
    set -gx XDG_CONFIG_HOME "$HOME/.config"
end
set -gx HOMESTEAD_CONFIG_DIR "$XDG_CONFIG_HOME/homestead"
# Preferências: $HOMESTEAD_CONFIG_DIR/preferences.yaml
# Perfil (instalados/favoritos): $HOMESTEAD_CONFIG_DIR/profile.yaml
`)
	default:
		fmt.Fprintf(stderr, "shell desconhecido: %q (use bash ou fish)\n", argv[0])
		return 2
	}
	return 0
}

func PrintHelp(stdout io.Writer) {
	fmt.Fprintln(stdout, `uso: homestead [opções]
       homestead export-profile [-format json|text] [-o caminho]
       homestead shell-init bash|fish

opções:
  -version   mostra a versão e termina

export-profile grava o que o Homestead registou no TUI (pacotes instalados com sucesso
e scripts marcados com f como favoritos), mais preferências relevantes para clonar
setup noutra máquina.

shell-init imprime linhas para colar em ~/.bashrc ou ~/.config/fish/config.fish.`)
}
