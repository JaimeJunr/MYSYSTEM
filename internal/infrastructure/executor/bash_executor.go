package executor

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// BashExecutor implements ScriptExecutor for bash scripts
type BashExecutor struct {
	rootDir string
}

// NewBashExecutor creates a new bash executor
func NewBashExecutor() interfaces.ScriptExecutor {
	rootDir, err := os.Getwd()
	if err != nil {
		// Fallback to empty string, will be relative paths
		rootDir = ""
	}

	return &BashExecutor{
		rootDir: rootDir,
	}
}

// Execute executes a bash script
func (e *BashExecutor) Execute(script *entities.Script) error {
	if err := e.Validate(script); err != nil {
		return fmt.Errorf("execute script %s: %w", script.ID, err)
	}

	// Construct full path to script
	scriptPath := filepath.Join(e.rootDir, script.Path)

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("execute script %s: script file not found at %s: %w",
			script.ID, scriptPath, types.ErrNotFound)
	}

	// Get current user information
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("execute script %s: get current user: %w", script.ID, err)
	}

	// Prepare command
	var cmd *exec.Cmd
	if script.RequiresSudo {
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execute script %s: %w", script.ID, types.ErrExecutionFailed)
	}

	return nil
}

// CanExecute checks if a script can be executed
func (e *BashExecutor) CanExecute(script *entities.Script) bool {
	if script == nil {
		return false
	}

	// Check if script file exists
	scriptPath := filepath.Join(e.rootDir, script.Path)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return false
	}

	// Check if bash is available
	if _, err := exec.LookPath("bash"); err != nil {
		return false
	}

	// If requires sudo, check if sudo is available
	if script.RequiresSudo {
		if _, err := exec.LookPath("sudo"); err != nil {
			return false
		}
	}

	return true
}

// Validate validates a script before execution
func (e *BashExecutor) Validate(script *entities.Script) error {
	if script == nil {
		return fmt.Errorf("validate script: %w", types.ErrInvalidInput)
	}

	// Validate entity
	if err := script.Validate(); err != nil {
		return fmt.Errorf("validate script: %w", err)
	}

	// Check if can execute
	if !e.CanExecute(script) {
		return fmt.Errorf("validate script: cannot execute script %s: %w",
			script.ID, types.ErrExecutionFailed)
	}

	return nil
}
