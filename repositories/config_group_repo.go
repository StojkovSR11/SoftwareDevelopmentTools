package repositories

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/hashicorp/consul/api"
	"projekat/model"
)

// ConfigGroupConsulRepository is an implementation of ConfigGroupRepository using Consul.
type ConfigGroupConsulRepository struct {
	cli *api.Client
}

func NewConfigGroupConsulRepository() (*ConfigGroupConsulRepository, error) {
	consulAddress := fmt.Sprintf("%s:%s", os.Getenv("DB"), os.Getenv("DBPORT"))
	config := api.DefaultConfig()
	config.Address = consulAddress

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConfigGroupConsulRepository{cli: client}, nil
}

// CreateConfigGroup creates a new configuration group.
func (repo *ConfigGroupConsulRepository) CreateConfigGroup(group model.ConfigGroup) error {
	key := fmt.Sprintf("configGroup/%s/%d", group.Name, group.Version)

	// Check if the configuration group already exists
	kv := repo.cli.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return err
	}
	if pair != nil {
		return fmt.Errorf("konfiguraciona grupa sa imenom %s i verzijom %d već postoji", group.Name, group.Version)
	}

	// Save the configuration group to Consul
	data, err := json.Marshal(group)
	if err != nil {
		return err
	}

	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	return err
}

// GetConfigGroup retrieves a configuration group by name and version.
func (repo *ConfigGroupConsulRepository) GetConfigGroup(name string, version int) (model.ConfigGroup, error) {
	key := fmt.Sprintf("configGroup/%s/%d", name, version)

	kv := repo.cli.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	if pair == nil {
		return model.ConfigGroup{}, fmt.Errorf("konfiguraciona grupa sa imenom %s i verzijom %d nije pronađena", name, version)
	}

	var group model.ConfigGroup
	err = json.Unmarshal(pair.Value, &group)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	return group, nil
}

// AddConfigurationToGroup adds a configuration to a configuration group by name and version.
func (repo *ConfigGroupConsulRepository) AddConfigurationToGroup(name string, version int, config model.GroupedConfig) error {
    group, err := repo.GetConfigGroup(name, version)
    if err != nil {
        return err
    }

    for _, c := range group.Configs {
        if c.Name == config.Name {
            return fmt.Errorf("konfiguracija '%s' već postoji u konfiguracionoj grupi '%s'", config.Name, name)
        }
    }

    group.Configs = append(group.Configs, config)

    return repo.UpdateConfigGroup(group)
}

// UpdateConfigGroup updates an existing configuration group.
func (repo *ConfigGroupConsulRepository) UpdateConfigGroup(group model.ConfigGroup) error {
    key := fmt.Sprintf("configGroup/%s/%d", group.Name, group.Version)

    // Save the updated configuration group to Consul
    data, err := json.Marshal(group)
    if err != nil {
        return err
    }

    p := &api.KVPair{Key: key, Value: data}
    _, err = repo.cli.KV().Put(p, nil)
    return err
}


// RemoveConfigurationFromGroup removes a configuration from a configuration group by name and version.
func (repo *ConfigGroupConsulRepository) RemoveConfigurationFromGroup(name string, version int, filterKey, filterValue string) error {
	group, err := repo.GetConfigGroup(name, version)
	if err != nil {
		return err
	}

	var updatedConfigs []model.GroupedConfig
	removed := false
	for _, c := range group.Configs {
		if val, ok := c.Labels[filterKey]; !ok || val != filterValue {
			updatedConfigs = append(updatedConfigs, c)
		} else {
			removed = true
		}
	}
	if !removed {
		return fmt.Errorf("konfiguracija sa filterom %s:%s nije pronađena u grupi", filterKey, filterValue)
	}
	group.Configs = updatedConfigs

	return repo.UpdateConfigGroup(group)
}


func (repo *ConfigGroupConsulRepository) GetConfigurationsFromGroup(name string, version int, filterKey, filterValue string) ([]model.GroupedConfig, error) {
    group, err := repo.GetConfigGroup(name, version)
    if err != nil {
        return nil, err
    }

    var filteredConfigs []model.GroupedConfig
    for _, c := range group.Configs {
        if val, ok := c.Labels[filterKey]; ok && val == filterValue {
            filteredConfigs = append(filteredConfigs, c)
        }
    }
    if len(filteredConfigs) == 0 {
        return nil, fmt.Errorf("nema konfiguracija sa filterom %s:%s u grupi", filterKey, filterValue)
    }
    return filteredConfigs, nil
}




// DeleteConfigGroup deletes a configuration group by name and version.
func (repo *ConfigGroupConsulRepository) DeleteConfigGroup(name string, version int) error {
	key := fmt.Sprintf("configGroup/%s/%d", name, version)

	kv := repo.cli.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return err
	}

	if pair == nil {
		return fmt.Errorf("konfiguraciona grupa sa imenom %s i verzijom %d nije pronađena", name, version)
	}

	_, err = kv.Delete(key, nil)
	return err
}

