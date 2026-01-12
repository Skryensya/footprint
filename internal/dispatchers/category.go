package dispatchers

type CommandCategory int

const (
	CategoryUncategorized CommandCategory = iota
	CategoryInfo
	CategoryConfig
)

func (c CommandCategory) String() string {
	switch c {
	case CategoryInfo:
		return "information and diagnostics"
	case CategoryConfig:
		return "configuration and preferences"
	default:
		return "other commands"
	}
}

var categoryOrder = []CommandCategory{
	CategoryInfo,
	CategoryConfig,
	CategoryUncategorized,
}
