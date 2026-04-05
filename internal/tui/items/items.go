package items

import (
	"strings"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

type MenuItem struct {
	Label  string
	Desc   string
	Action string
}

func (i MenuItem) Title() string       { return i.Label }
func (i MenuItem) Description() string { return i.Desc }
func (i MenuItem) FilterValue() string { return i.Label }

type ScriptItem struct {
	Script   entities.Script
	Favorite bool
}

func (i ScriptItem) Title() string { return i.Script.Name }
func (i ScriptItem) Description() string {
	if i.Favorite {
		return "★ " + i.Script.Description
	}
	return i.Script.Description
}
func (i ScriptItem) FilterValue() string {
	return strings.ToLower(strings.TrimSpace(i.Script.Name + " " + i.Script.Description + " " + i.Script.ID))
}

type PackageItem struct {
	Pkg entities.Package
}

func (i PackageItem) Title() string       { return i.Pkg.Name }
func (i PackageItem) Description() string { return i.Pkg.Description }
func (i PackageItem) FilterValue() string {
	return strings.ToLower(strings.TrimSpace(i.Pkg.Name + " " + i.Pkg.Description + " " + i.Pkg.ID))
}

type InstallerCategoryItem struct {
	Heading    string
	Desc       string
	Categories []types.PackageCategory
}

func (i InstallerCategoryItem) Title() string       { return i.Heading }
func (i InstallerCategoryItem) Description() string { return i.Desc }
func (i InstallerCategoryItem) FilterValue() string { return i.Heading }
