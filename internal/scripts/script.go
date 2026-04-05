package scripts

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// Script represents a system maintenance script
type Script struct {
	ID           string
	Name         string
	Description  string
	Path         string
	Category     string
	RequiresSudo bool
	Native string // monitores nativos do TUI; Path pode ficar vazio
}

// ScriptCategory represents different categories of scripts
type ScriptCategory string

const (
	CategoryCleanup    ScriptCategory = "cleanup"
	CategoryMonitoring ScriptCategory = "monitoring"
	CategoryInstall    ScriptCategory = "install"
)

// GetAllScripts returns all available scripts
func GetAllScripts() []Script {
	return []Script{
		{
			ID:          "cleanup-full",
			Name:        "Limpeza Completa (SSD)",
			Description: "Orquestrador completo de limpeza do sistema",
			Path:        "scripts/cleanup/limpar_ssd.sh",
			Category:    string(CategoryCleanup),
			RequiresSudo: true,
		},
		{
			ID:          "cleanup-general",
			Name:        "Limpeza Geral (Caches)",
			Description: "Limpa caches de Docker, Poetry, npm, apt, etc.",
			Path:        "scripts/cleanup/limpar_geral.sh",
			Category:    string(CategoryCleanup),
			RequiresSudo: true,
		},
		{
			ID:          "cleanup-large",
			Name:        "Buscar Arquivos Grandes",
			Description: "Encontra e remove arquivos/pastas grandes (>100MB)",
			Path:        "scripts/cleanup/limpar_grandes.sh",
			Category:    string(CategoryCleanup),
			RequiresSudo: true,
		},
		{
			ID:           "monitor-battery",
			Name:         "Monitor de Bateria",
			Description:  "Carga, carregador e detalhes da bateria",
			Path:         "",
			Category:     string(CategoryMonitoring),
			RequiresSudo: false,
			Native:       "battery",
		},
		{
			ID:           "monitor-memory",
			Name:         "Uso de Memória",
			Description:  "Uso de memória RAM e swap",
			Path:         "",
			Category:     string(CategoryMonitoring),
			RequiresSudo: false,
			Native:       "memory",
		},
		{
			ID:           "monitor-disk",
			Name:         "Espaço em disco",
			Description:  "Uso por ponto de montagem",
			Path:         "",
			Category:     string(CategoryMonitoring),
			RequiresSudo: false,
			Native:       "disk",
		},
		{
			ID:           "monitor-load",
			Name:         "Carga da CPU",
			Description:  "Load average",
			Path:         "",
			Category:     string(CategoryMonitoring),
			RequiresSudo: false,
			Native:       "load",
		},
		{
			ID:           "monitor-network",
			Name:         "Rede (interfaces)",
			Description:  "RX/TX e vazão simples",
			Path:         "",
			Category:     string(CategoryMonitoring),
			RequiresSudo: false,
			Native:       "network",
		},
		{
			ID:           "monitor-thermal",
			Name:         "Temperatura",
			Description:  "Sensores em sysfs",
			Path:         "",
			Category:     string(CategoryMonitoring),
			RequiresSudo: false,
			Native:       "thermal",
		},
		{
			ID:           "monitor-systemd-user",
			Name:         "systemd usuário (falhas)",
			Description:  "Unidades --user falhando",
			Path:         "",
			Category:     string(CategoryMonitoring),
			RequiresSudo: false,
			Native:       "systemd-user",
		},
	}
}

// GetScriptsByCategory returns scripts filtered by category
func GetScriptsByCategory(category ScriptCategory) []Script {
	all := GetAllScripts()
	filtered := []Script{}

	for _, script := range all {
		if script.Category == string(category) {
			filtered = append(filtered, script)
		}
	}

	return filtered
}

// Execute runs a script with proper environment setup
func (s *Script) Execute() error {
	if s.Native != "" {
		return fmt.Errorf("este item só funciona no Homestead (menu Monitoramento)")
	}
	// Get project root directory
	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Construct full path to script
	scriptPath := filepath.Join(rootDir, s.Path)

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script not found: %s", scriptPath)
	}

	// Get current user information
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Prepare command
	var cmd *exec.Cmd
	if s.RequiresSudo {
		cmd = exec.Command("sudo", "-E", "bash", scriptPath)
	} else {
		cmd = exec.Command("bash", scriptPath)
	}

	// Set environment variables (preserve user context for sudo)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("REAL_USER=%s", currentUser.Username),
		fmt.Sprintf("REAL_HOME=%s", currentUser.HomeDir),
	)

	// Connect to terminal for interactive scripts
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute
	return cmd.Run()
}
