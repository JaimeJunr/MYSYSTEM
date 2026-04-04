package types

// PackageCategory represents the category of a package
type PackageCategory string

const (
	PackageCategoryIDE       PackageCategory = "ide"
	PackageCategoryTool      PackageCategory = "tool"
	PackageCategoryApp       PackageCategory = "app"
	PackageCategoryZshCore   PackageCategory = "zsh_core"
	PackageCategoryTerminal  PackageCategory = "terminal"
	PackageCategoryShell     PackageCategory = "shell"
	PackageCategoryAI        PackageCategory = "ai"
	PackageCategoryGames     PackageCategory = "games"
	PackageCategorySysAdmin  PackageCategory = "sys_admin"
	PackageCategoryOther     PackageCategory = "other"
)

// IsValid checks if the category is valid
func (c PackageCategory) IsValid() bool {
	switch c {
	case PackageCategoryIDE,
		PackageCategoryTool,
		PackageCategoryApp,
		PackageCategoryZshCore,
		PackageCategoryTerminal,
		PackageCategoryShell,
		PackageCategoryAI,
		PackageCategoryGames,
		PackageCategorySysAdmin,
		PackageCategoryOther:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (c PackageCategory) String() string {
	return string(c)
}
