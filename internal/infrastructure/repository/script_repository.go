package repository

import (
	"fmt"
	"sort"
	"sync"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// InMemoryScriptRepository is an in-memory implementation of ScriptRepository
type InMemoryScriptRepository struct {
	scripts map[string]*entities.Script
	mu      sync.RWMutex
}

// NewInMemoryScriptRepository creates a new in-memory script repository
// and initializes it with the default scripts
func NewInMemoryScriptRepository() interfaces.ScriptRepository {
	repo := &InMemoryScriptRepository{
		scripts: make(map[string]*entities.Script),
	}

	// Initialize with default scripts
	repo.initializeDefaultScripts()

	return repo
}

// initializeDefaultScripts populates the repository with default scripts
func (r *InMemoryScriptRepository) initializeDefaultScripts() {
	defaultScripts := []entities.Script{
		{
			ID:           "cleanup-full",
			Name:         "Limpeza Completa (SSD)",
			Description:  "Orquestrador completo de limpeza do sistema",
			Path:         "scripts/cleanup/limpar_ssd.sh",
			Category:     types.CategoryCleanup,
			RequiresSudo: true,
		},
		{
			ID:           "cleanup-general",
			Name:         "Limpeza Geral (Caches)",
			Description:  "Limpa caches de Docker, Poetry, npm, apt, etc.",
			Path:         "scripts/cleanup/limpar_geral.sh",
			Category:     types.CategoryCleanup,
			RequiresSudo: true,
		},
		{
			ID:           "cleanup-large",
			Name:         "Buscar Arquivos Grandes",
			Description:  "Encontra e remove arquivos/pastas grandes (>100MB)",
			Path:         "scripts/cleanup/limpar_grandes.sh",
			Category:     types.CategoryCleanup,
			RequiresSudo: true,
		},
		{
			ID:            "monitor-battery",
			Name:          "Monitor de Bateria",
			Description:   "Carga, carregador e detalhes da bateria (atualiza automaticamente)",
			Path:          "",
			Category:      types.CategoryMonitoring,
			RequiresSudo:  false,
			NativeMonitor: entities.NativeMonitorBattery,
		},
		{
			ID:            "monitor-memory",
			Name:          "Uso de Memória",
			Description:   "Uso de memória RAM e swap",
			Path:          "",
			Category:      types.CategoryMonitoring,
			RequiresSudo:  false,
			NativeMonitor: entities.NativeMonitorMemory,
		},
	}

	for i := range defaultScripts {
		r.scripts[defaultScripts[i].ID] = &defaultScripts[i]
	}
	for _, s := range defaultUtilityScripts() {
		sc := s
		r.scripts[sc.ID] = &sc
	}
}

// FindAll returns all scripts
func (r *InMemoryScriptRepository) FindAll() ([]entities.Script, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	scripts := make([]entities.Script, 0, len(r.scripts))
	for _, script := range r.scripts {
		scripts = append(scripts, *script)
	}

	return scripts, nil
}

// FindByID returns a script by its ID
func (r *InMemoryScriptRepository) FindByID(id string) (*entities.Script, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	script, ok := r.scripts[id]
	if !ok {
		return nil, fmt.Errorf("script %s: %w", id, types.ErrNotFound)
	}

	// Return a copy to prevent external modification
	scriptCopy := *script
	return &scriptCopy, nil
}

// FindByCategory returns all scripts in a category
func (r *InMemoryScriptRepository) FindByCategory(category types.Category) ([]entities.Script, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	scripts := make([]entities.Script, 0)
	for _, script := range r.scripts {
		if script.Category == category {
			scripts = append(scripts, *script)
		}
	}

	sort.Slice(scripts, func(i, j int) bool {
		return scripts[i].Name < scripts[j].Name
	})

	return scripts, nil
}

// Save saves a script
func (r *InMemoryScriptRepository) Save(script *entities.Script) error {
	if err := script.Validate(); err != nil {
		return fmt.Errorf("save script: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Make a copy to store
	scriptCopy := *script
	r.scripts[script.ID] = &scriptCopy

	return nil
}

// Delete deletes a script by ID
func (r *InMemoryScriptRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.scripts[id]; !ok {
		return fmt.Errorf("delete script %s: %w", id, types.ErrNotFound)
	}

	delete(r.scripts, id)
	return nil
}

// Exists checks if a script exists
func (r *InMemoryScriptRepository) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.scripts[id]
	return ok
}
