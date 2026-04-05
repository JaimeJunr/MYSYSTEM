package scripts

import (
	"testing"
)

func TestGetAllScripts(t *testing.T) {
	scripts := GetAllScripts()

	if len(scripts) == 0 {
		t.Error("Expected at least one script, got 0")
	}

	// Verify we have the expected scripts
	expectedCount := 10
	if len(scripts) != expectedCount {
		t.Errorf("Expected %d scripts, got %d", expectedCount, len(scripts))
	}

	// Verify each script has required fields
	for i, script := range scripts {
		if script.ID == "" {
			t.Errorf("Script %d: ID is empty", i)
		}
		if script.Name == "" {
			t.Errorf("Script %d: Name is empty", i)
		}
		if script.Description == "" {
			t.Errorf("Script %d: Description is empty", i)
		}
		if script.Path == "" && script.Native == "" {
			t.Errorf("Script %d: Path is empty without Native", i)
		}
		if script.Category == "" {
			t.Errorf("Script %d: Category is empty", i)
		}
	}
}

func TestGetScriptsByCategory(t *testing.T) {
	tests := []struct {
		category     ScriptCategory
		expectedMin  int
		description  string
	}{
		{CategoryCleanup, 3, "cleanup scripts"},
		{CategoryMonitoring, 7, "monitoring scripts"},
		{CategoryInstall, 0, "install scripts (none yet)"},
	}

	for _, tt := range tests {
		t.Run(string(tt.category), func(t *testing.T) {
			scripts := GetScriptsByCategory(tt.category)

			if len(scripts) < tt.expectedMin {
				t.Errorf("Expected at least %d %s, got %d",
					tt.expectedMin, tt.description, len(scripts))
			}

			// Verify all returned scripts belong to the requested category
			for _, script := range scripts {
				if script.Category != string(tt.category) {
					t.Errorf("Script %s has category %s, expected %s",
						script.Name, script.Category, tt.category)
				}
			}
		})
	}
}

func TestScriptFields(t *testing.T) {
	tests := []struct {
		id           string
		expectedName string
		expectedPath string
		requiresSudo bool
	}{
		{
			id:           "cleanup-full",
			expectedName: "Limpeza Completa (SSD)",
			expectedPath: "scripts/cleanup/limpar_ssd.sh",
			requiresSudo: true,
		},
		{
			id:           "monitor-battery",
			expectedName: "Monitor de Bateria",
			expectedPath: "",
			requiresSudo: false,
		},
		{
			id:           "monitor-memory",
			expectedName: "Uso de Memória",
			expectedPath: "",
			requiresSudo: false,
		},
	}

	allScripts := GetAllScripts()

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			var found *Script
			for _, script := range allScripts {
				if script.ID == tt.id {
					found = &script
					break
				}
			}

			if found == nil {
				t.Fatalf("Script with ID %s not found", tt.id)
			}

			if found.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, found.Name)
			}

			if found.Path != tt.expectedPath {
				t.Errorf("Expected path %s, got %s", tt.expectedPath, found.Path)
			}

			if found.RequiresSudo != tt.requiresSudo {
				t.Errorf("Expected RequiresSudo %v, got %v", tt.requiresSudo, found.RequiresSudo)
			}
		})
	}
}

func TestScriptCategories(t *testing.T) {
	// Verify category constants
	if CategoryCleanup != "cleanup" {
		t.Errorf("CategoryCleanup should be 'cleanup', got %s", CategoryCleanup)
	}
	if CategoryMonitoring != "monitoring" {
		t.Errorf("CategoryMonitoring should be 'monitoring', got %s", CategoryMonitoring)
	}
	if CategoryInstall != "install" {
		t.Errorf("CategoryInstall should be 'install', got %s", CategoryInstall)
	}
}

func TestScriptIDsUnique(t *testing.T) {
	scripts := GetAllScripts()
	seen := make(map[string]bool)

	for _, script := range scripts {
		if seen[script.ID] {
			t.Errorf("Duplicate script ID found: %s", script.ID)
		}
		seen[script.ID] = true
	}
}

// Benchmark tests
func BenchmarkGetAllScripts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllScripts()
	}
}

func BenchmarkGetScriptsByCategory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetScriptsByCategory(CategoryCleanup)
	}
}
