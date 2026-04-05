package types

// Category represents different categories of scripts and packages
type Category string

const (
	CategoryCleanup    Category = "cleanup"
	CategoryMonitoring Category = "monitoring"
	CategoryCheckup    Category = "checkup"
	CategoryInstall    Category = "install"
	CategoryUtilities  Category = "utilities"
)

// String returns the string representation of the category
func (c Category) String() string {
	return string(c)
}

// IsValid checks if the category is valid
func (c Category) IsValid() bool {
	switch c {
	case CategoryCleanup, CategoryMonitoring, CategoryCheckup, CategoryInstall, CategoryUtilities:
		return true
	default:
		return false
	}
}
