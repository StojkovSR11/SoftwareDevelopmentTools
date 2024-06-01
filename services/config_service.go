package services

import "projekat/model"

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(repo model.ConfigRepository) *ConfigService {
	return &ConfigService{repo: repo}
}

func (c *ConfigService) CreateConfig(config model.Config) error {
	return c.repo.CreateConfig(config)
}

func (c *ConfigService) GetConfig(name string, version int) (model.Config, error) {
	return c.repo.GetConfig(name, version)
}

func (c *ConfigService) UpdateConfig(name string, version int, newConfig model.Config) error {
	return c.repo.UpdateConfig(name, version, newConfig)
}

func (c *ConfigService) DeleteConfig(name string, version int) error {
	return c.repo.DeleteConfig(name, version)
}
