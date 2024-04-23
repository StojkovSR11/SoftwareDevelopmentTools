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
func (repo *ConfigGroupInMemoryRepository) AddConfigurationToGroup(name string, version int, config model.Config) error {
    key := version
    group, exists := repo.configGroups[name][key]
    //Ako ne postoji baca poruku
    if !exists {
        return fmt.Errorf("konfiguraciona grupa sa imenom %s i verzijom %d nije pronađena", name, version)
    }
    group.Configs = append(group.Configs, config)
    repo.configGroups[name][key] = group
    return nil
}

// RemoveConfigurationFromGroup uklanja konfiguraciju iz konfiguracione grupe po imenu i verziji.
func (repo *ConfigGroupInMemoryRepository) RemoveConfigurationFromGroup(name string, version int, configName string) error {
    key := version
    group, exists := repo.configGroups[name][key]
    //Ako ne postoji baca poruku
    if !exists {
        return fmt.Errorf("konfiguraciona grupa sa imenom %s i verzijom %d nije pronađena", name, version)
    }
    //pravimo listu azuriranih konfiguracija
    var updatedConfigs []model.Config
    removed := false
    //prolazimo kroz mape konfiguracija
    for _, c := range group.Configs {
        //Ako ime nije isto kao ime koje zelimo da izbrisemo dodaje ga u novu listu konfiguracija
        if c.Name != configName {
            updatedConfigs = append(updatedConfigs, c)
        
        }//ako jeste postavi ga na true
         else {
            removed = true
        }
    }
    //ako ga ne pronace baca poruku
    if !removed {
        return fmt.Errorf("konfiguracija sa imenom %s nije pronađena u grupi", configName)
    }
    //azuriramo staru mapu konfiguracije novom
    group.Configs = updatedConfigs
    repo.configGroups[name][key] = group
    return nil
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
