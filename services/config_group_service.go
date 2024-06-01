package services

import "projekat/model"

type ConfigGroupService struct {
	repo model.ConfigGroupRepository
}

func NewConfigGroupInService(repo model.ConfigGroupRepository) ConfigGroupService {
	return ConfigGroupService{
		repo: repo,
	}
}

// CreateConfigGroup creates a new configuration group.
func (s ConfigGroupService) CreateConfigGroup(group model.ConfigGroup) error {
	return s.repo.CreateConfigGroup(group)
}

// GetConfigGroup retrieves a configuration group by its name and version.
func (s ConfigGroupService) GetConfigGroup(name string, version int) (model.ConfigGroup, error) {
	return s.repo.GetConfigGroup(name, version)
}

// AddConfigurationToGroup adds a configuration to a configuration group by its name and version.
func (s ConfigGroupService) AddConfigurationToGroup(name string, version int, config model.GroupedConfig) error {
	return s.repo.AddConfigurationToGroup(name, version, config)
}

func (s *ConfigGroupService) RemoveConfigurationFromGroup(name string, version int, filterKey, filterValue string) error {
    return s.repo.RemoveConfigurationFromGroup(name, version, filterKey, filterValue)
}

// DeleteConfigGroup deletes a configuration group by its name and version.
func (s ConfigGroupService) DeleteConfigGroup(name string, version int) error {
	return s.repo.DeleteConfigGroup(name, version)
}
func (s *ConfigGroupService) GetConfigurationsFromGroup(name string, version int, filterKey, filterValue string) ([]model.GroupedConfig, error) {
    return s.repo.GetConfigurationsFromGroup(name, version, filterKey, filterValue)
}

