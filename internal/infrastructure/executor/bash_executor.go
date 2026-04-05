package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// BashExecutor implements ScriptExecutor for bash scripts
type BashExecutor struct {
	rootDir string
}

func NewBashExecutor() interfaces.ScriptExecutor {
	return NewBashExecutorWithRoot("")
}

func NewBashExecutorWithRoot(root string) interfaces.ScriptExecutor {
	e := &BashExecutor{}
	_ = e.SetScriptRoot(root)
	return e
}

// ResolveScriptRoot returns the absolute script/repo root; empty dir means cwd.
func ResolveScriptRoot(dir string) (string, error) {
	return expandScriptRoot(dir)
}

func expandScriptRoot(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get working directory: %w", err)
		}
		return filepath.Clean(cwd), nil
	}
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("home dir: %w", err)
		}
		if path == "~" {
			return filepath.Clean(home), nil
		}
		return filepath.Clean(filepath.Join(home, path[2:])), nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

func (e *BashExecutor) SetScriptRoot(dir string) error {
	root, err := expandScriptRoot(dir)
	if err != nil {
		return err
	}
	e.rootDir = root
	return nil
}

func (e *BashExecutor) newScriptCmd(script *entities.Script, opts interfaces.ScriptExecOpts) (*exec.Cmd, error) {
	if script.NativeMonitor != "" {
		return nil, fmt.Errorf("native monitor script %s", script.ID)
	}
	scriptPath := filepath.Join(e.rootDir, script.Path)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("execute script %s: script file not found at %s: %w",
			script.ID, scriptPath, types.ErrNotFound)
	}

	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("execute script %s: get current user: %w", script.ID, err)
	}

	var cmd *exec.Cmd
	if script.RequiresSudo {
		cmd = exec.Command("sudo", "-E", "bash", scriptPath)
	} else {
		cmd = exec.Command("bash", scriptPath)
	}

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("REAL_USER=%s", currentUser.Username),
		fmt.Sprintf("REAL_HOME=%s", currentUser.HomeDir),
		fmt.Sprintf("HOMESTEAD_ROOT=%s", e.rootDir),
	)
	if opts.DryRun {
		cmd.Env = append(cmd.Env, "HOMESTEAD_DRY_RUN=1")
	}

	return cmd, nil
}

// Execute executes a bash script attached to the current terminal
func (e *BashExecutor) Execute(script *entities.Script, opts interfaces.ScriptExecOpts) error {
	if err := e.Validate(script); err != nil {
		return fmt.Errorf("execute script %s: %w", script.ID, err)
	}
	if script.NativeMonitor != "" {
		return fmt.Errorf("execute script %s: native monitor (use o TUI): %w", script.ID, types.ErrInvalidInput)
	}

	cmd, err := e.newScriptCmd(script, opts)
	if err != nil {
		return err
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execute script %s: %w", script.ID, types.ErrExecutionFailed)
	}

	return nil
}

// ExecuteCapture runs the script and returns combined stdout/stderr (no TTY).
func (e *BashExecutor) ExecuteCapture(script *entities.Script, opts interfaces.ScriptExecOpts) (string, error) {
	if err := e.Validate(script); err != nil {
		return "", fmt.Errorf("execute script %s: %w", script.ID, err)
	}
	if script.NativeMonitor != "" {
		return "", fmt.Errorf("execute script %s: native monitor: %w", script.ID, types.ErrInvalidInput)
	}

	cmd, err := e.newScriptCmd(script, opts)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	cmd.Stdin = nil

	runErr := cmd.Run()
	out := buf.String()
	if runErr != nil {
		return out, fmt.Errorf("execute script %s: %w", script.ID, types.ErrExecutionFailed)
	}
	return out, nil
}

// InteractiveCommand returns a command for tea.ExecProcess (sudo / password prompts).
func (e *BashExecutor) InteractiveCommand(script *entities.Script, opts interfaces.ScriptExecOpts) (*exec.Cmd, error) {
	if err := e.Validate(script); err != nil {
		return nil, fmt.Errorf("execute script %s: %w", script.ID, err)
	}
	if script.NativeMonitor != "" {
		return nil, fmt.Errorf("execute script %s: native monitor: %w", script.ID, types.ErrInvalidInput)
	}
	return e.newScriptCmd(script, opts)
}

// CanExecute checks if a script can be executed
func (e *BashExecutor) CanExecute(script *entities.Script) bool {
	if script == nil {
		return false
	}
	if script.NativeMonitor != "" {
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
	if script.NativeMonitor != "" {
		return nil
	}

	// Check if can execute
	if !e.CanExecute(script) {
		return fmt.Errorf("validate script: cannot execute script %s: %w",
			script.ID, types.ErrExecutionFailed)
	}

	return nil
}
