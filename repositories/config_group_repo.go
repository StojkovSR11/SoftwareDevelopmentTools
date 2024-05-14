package repositories

import (
	"fmt"
	"projekat/model"
)

// ConfigGroupInMemoryRepository je implementacija ConfigGroupRepository interfejsa u memoriji.
type ConfigGroupInMemoryRepository struct {
	configGroups map[string]map[int]model.ConfigGroup
}

// NewConfigGroupInMemoryRepository kreira novu instancu ConfigGroupInMemoryRepository.
func NewConfigGroupInMemoryRepository() model.ConfigGroupRepository {
	return &ConfigGroupInMemoryRepository{
		configGroups: make(map[string]map[int]model.ConfigGroup),
	}
}

// CreateConfigGroup kreira novu konfiguracionu grupu.
func (repo *ConfigGroupInMemoryRepository) CreateConfigGroup(group model.ConfigGroup) error {
	key := group.Version
	//Proveri da li mapa konf grupa postoji
	if repo.configGroups[group.Name] == nil {
		// Ako ne postoji, napravi je
		repo.configGroups[group.Name] = make(map[int]model.ConfigGroup)
	} else {
		// Proveri da li postoji unutrasnja mapa konfiguracija
		if _, exists := repo.configGroups[group.Name][key]; exists {
			return fmt.Errorf("konfiguraciona grupa sa imenom %s i verzijom %d već postoji", group.Name, group.Version)
		}
	}
	// Dodeli mapu grupi konf
	repo.configGroups[group.Name][key] = group
	return nil
}

// GetConfigGroup dohvata konfiguracionu grupu po imenu i verziji.
func (repo *ConfigGroupInMemoryRepository) GetConfigGroup(name string, version int) (model.ConfigGroup, error) {
	group, exists := repo.configGroups[name][version]
	//Ako ne postoji baca poruku
	if !exists {
		return model.ConfigGroup{}, fmt.Errorf("konfiguraciona grupa sa imenom %s i verzijom %d nije pronađena", name, version)
	}
	return group, nil
}

// AddConfigurationToGroup dodaje konfiguraciju u konfiguracionu grupu po imenu i verziji.
func (repo *ConfigGroupInMemoryRepository) AddConfigurationToGroup(name string, version int, config model.GroupedConfig) error {
    key := version
    group, exists := repo.configGroups[name][key]
    // Provera da li konfiguraciona grupa postoji
    if !exists {
        return fmt.Errorf("konfiguraciona grupa sa imenom %s i verzijom %d nije pronađena", name, version)
    }
    // Provera da li konfiguracija već postoji unutar grupe
    for _, c := range group.Configs {
        if c.Name == config.Name {
            return fmt.Errorf("konfiguracija '%s' već postoji u konfiguracionoj grupi '%s'", config.Name, name)
        }
    }
    // Dodavanje konfiguracije u grupu
    group.Configs = append(group.Configs, config)
    repo.configGroups[name][key] = group
    return nil
}

// RemoveConfigurationFromGroup uklanja konfiguraciju iz konfiguracione grupe po imenu i verziji.
func (repo *ConfigGroupInMemoryRepository) RemoveConfigurationFromGroup(name string, version int, filter string) error {
	key := version
	group, exists := repo.configGroups[name][key]
	if !exists {
		return fmt.Errorf("konfiguraciona grupa sa imenom %s i verzijom %d nije pronađena", name, version)
	}
	var updatedConfigs []model.GroupedConfig
	removed := false
	for _, c := range group.Configs {
		// Use the filter condition here instead of comparing with configName
		if c.Labels[filter] == "" {
			updatedConfigs = append(updatedConfigs, c)
		} else {
			removed = true
		}
	}
	if !removed {
		return fmt.Errorf("konfiguracija sa filterom %s nije pronađena u grupi", filter)
	}
	group.Configs = updatedConfigs
	repo.configGroups[name][key] = group
	return nil
}
func (repo *ConfigGroupInMemoryRepository) GetConfigurationsFromGroup(name string, version int, filter string) ([]model.GroupedConfig, error) {
	key := version
	group, exists := repo.configGroups[name][key]
	if !exists {
		return nil, fmt.Errorf("konfiguraciona grupa sa imenom %s i verzijom %d nije pronađena", name, version)
	}
	var filteredConfigs []model.GroupedConfig
	for _, c := range group.Configs {
		if c.Labels[filter] != "" {
			filteredConfigs = append(filteredConfigs, c)
		}
	}
	if len(filteredConfigs) == 0 {
		return nil, fmt.Errorf("nema konfiguracija sa filterom %s u grupi", filter)
	}
	return filteredConfigs, nil
}

// DeleteConfigGroup briše konfiguracionu grupu po imenu i verziji.
func (repo *ConfigGroupInMemoryRepository) DeleteConfigGroup(name string, version int) error {
	key := version
	if _, exists := repo.configGroups[name][key];
	//ako ne postoji baca poruku
	!exists {
		return fmt.Errorf("konfiguraciona grupa sa imenom %s i verzijom %d nije pronađena", name, version)
	}
	delete(repo.configGroups[name], key)
	return nil
}
