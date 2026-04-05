package profilestate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// State is persisted under the Homestead XDG config dir (profile.yaml).
// It records installs triggered from the TUI and user-marked script favorites for migration.
type State struct {
	InstalledPackageIDs []string `yaml:"installed_packages,omitempty"`
	FavoriteScriptIDs   []string `yaml:"favorite_scripts,omitempty"`
}

func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user config dir: %w", err)
	}
	return filepath.Join(dir, "homestead", "profile.yaml"), nil
}

func Load(path string) (State, error) {
	var s State
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return State{}, nil
		}
		return State{}, fmt.Errorf("read profile state: %w", err)
	}
	if err := yaml.Unmarshal(data, &s); err != nil {
		return State{}, fmt.Errorf("parse profile state: %w", err)
	}
	s.Normalize()
	return s, nil
}

func Save(path string, s State) error {
	s.Normalize()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("profile state dir: %w", err)
	}
	data, err := yaml.Marshal(&s)
	if err != nil {
		return fmt.Errorf("marshal profile state: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write profile state: %w", err)
	}
	return nil
}

func (s *State) Normalize() {
	s.InstalledPackageIDs = dedupeSorted(s.InstalledPackageIDs)
	s.FavoriteScriptIDs = dedupeSorted(s.FavoriteScriptIDs)
}

// RecordInstalled adds a package id after a successful TUI install (idempotent).
func RecordInstalled(s *State, packageID string) {
	id := strings.TrimSpace(packageID)
	if id == "" || s == nil {
		return
	}
	s.InstalledPackageIDs = append(s.InstalledPackageIDs, id)
	s.Normalize()
}

// ToggleFavorite returns true if id is favorite after the call.
func ToggleFavorite(s *State, scriptID string) bool {
	id := strings.TrimSpace(scriptID)
	if id == "" || s == nil {
		return false
	}
	for i, x := range s.FavoriteScriptIDs {
		if x == id {
			s.FavoriteScriptIDs = append(s.FavoriteScriptIDs[:i], s.FavoriteScriptIDs[i+1:]...)
			s.Normalize()
			return false
		}
	}
	s.FavoriteScriptIDs = append(s.FavoriteScriptIDs, id)
	s.Normalize()
	return true
}

func IsFavorite(s *State, scriptID string) bool {
	if s == nil {
		return false
	}
	id := strings.TrimSpace(scriptID)
	for _, x := range s.FavoriteScriptIDs {
		if x == id {
			return true
		}
	}
	return false
}

func dedupeSorted(in []string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, raw := range in {
		x := strings.TrimSpace(raw)
		if x == "" {
			continue
		}
		if _, ok := seen[x]; ok {
			continue
		}
		seen[x] = struct{}{}
		out = append(out, x)
	}
	sort.Strings(out)
	return out
}
