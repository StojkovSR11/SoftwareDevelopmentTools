package services

import (
	"projekat/model"
)

//struktura koja omogućava rad sa konfiguracijama preko repozitorijuma.
type ConfigService struct {
	repo model.ConfigRepository
}

// NewConfigInService kreira novu instancu servisa za konfiguracije.
func NewConfigInService(repo model.ConfigRepository) ConfigService {
	return ConfigService{
		repo: repo,
	}
}

// CreateConfig kreira novu konfiguraciju.
func (s ConfigService) CreateConfig(config model.Config) error {
	return s.repo.CreateConfig(config)
}

// GetConfig dobavlja konfiguraciju po imenu i verziji.
func (s ConfigService) GetConfig(name string, version int) (model.Config, error) {
	return s.repo.GetConfig(name, version)
}

// UpdateConfig ažurira postojeću konfiguraciju.
func (s ConfigService) UpdateConfig(name string, version int, newConfig model.Config) error {
	return s.repo.UpdateConfig(name, version, newConfig)
}

// DeleteConfig briše konfiguraciju po imenu i verziji.
func (s ConfigService) DeleteConfig(name string, version int) error {
	return s.repo.DeleteConfig(name, version)
}
