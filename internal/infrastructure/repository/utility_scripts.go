package repository

import (
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// defaultUtilityScripts aligns catalog utility installers with ScriptRepository entries.
func defaultUtilityScripts() []entities.Script {
	defs := utilityInstallerDefinitions()
	out := make([]entities.Script, len(defs))
	for i, d := range defs {
		out[i] = entities.Script{
			ID:           d.ID,
			Name:         d.Name,
			Description:  d.Description,
			Path:         d.Path,
			Category:     types.CategoryUtilities,
			RequiresSudo: d.RequiresSudo,
		}
	}
	return out
}

func defaultUtilityPackages() []entities.Package {
	defs := utilityInstallerDefinitions()
	out := make([]entities.Package, len(defs))
	for i, d := range defs {
		out[i] = entities.Package{
			ID:                d.ID,
			Name:              d.Name,
			Description:       d.Description,
			Version:           "latest",
			Category:          types.PackageCategoryUtilities,
			UtilityScriptPath: d.Path,
			RequiresSudo:      d.RequiresSudo,
			ProjectURL:        d.ProjectURL,
		}
	}
	return out
}
