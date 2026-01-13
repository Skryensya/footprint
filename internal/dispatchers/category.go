package dispatchers

type CommandCategory int

const (
	CategoryUncategorized CommandCategory = iota
	CategoryInfo
	CategoryConfig
	CategoryRepo
)

func (c CommandCategory) String() string {
	switch c {
	case CategoryInfo:
		return "Information and diagnostics"
	case CategoryConfig:
		return "Configuration and preferences"
	case CategoryRepo:
		return "Handle repository tracking"
	default:
		return "other commands"
	}
}

var categoryOrder = []CommandCategory{
	CategoryInfo,
	CategoryRepo,
	CategoryConfig,
	CategoryUncategorized,
}
