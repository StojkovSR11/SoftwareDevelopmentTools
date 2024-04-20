package model

type ConfigGroup struct {
    Name       string    `json:"name"`
    Version    int       `json:"version"`
    Configs    []Config  `json:"configs"`
}

type ConfigGroupRepository interface {

    CreateConfigGroup(group ConfigGroup) error

    GetConfigGroup(name string, version int) (ConfigGroup, error)

    AddConfigurationToGroup(name string, version int, config Config) error

    RemoveConfigurationFromGroup(name string, version int, configName string) error

    DeleteConfigGroup(name string, version int) error
	
}
