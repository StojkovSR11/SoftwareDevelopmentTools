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
func (s ConfigGroupService) AddConfigurationToGroup(name string, version int, config model.Config) error {
	return s.repo.AddConfigurationToGroup(name, version, config)
}

// RemoveConfigurationFromGroup removes a configuration from a configuration group by its name and version.
func (s ConfigGroupService) RemoveConfigurationFromGroup(name string, version int, configName string) error {
	return s.repo.RemoveConfigurationFromGroup(name, version, configName)
}

// DeleteConfigGroup deletes a configuration group by its name and version.
func (s ConfigGroupService) DeleteConfigGroup(name string, version int) error {
	return s.repo.DeleteConfigGroup(name, version)
}
