package model

type ConfigGroup struct {
	Name    string          `json:"name"`
	Version int             `json:"version"`
	Configs []GroupedConfig `json:"configs"`
}

type ConfigGroupRepository interface {
	CreateConfigGroup(group ConfigGroup) error

	GetConfigGroup(name string, version int) (ConfigGroup, error)

	AddConfigurationToGroup(name string, version int, config GroupedConfig) error

	RemoveConfigurationFromGroup(name string, version int, configName string) error

	DeleteConfigGroup(name string, version int) error

	GetConfigurationsFromGroup(name string, version int, filter string) ([]GroupedConfig, error)
}
