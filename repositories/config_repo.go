package repositories

import (
	"encoding/json"
	"fmt"
	"os"
	"projekat/model"

	"github.com/hashicorp/consul/api"
)

type ConfigConsulRepository struct {
	cli *api.Client
}

// Kreiranje nove instance repozitorijuma za konfiguracije sa Consul
func NewConfigConsulRepository() (*ConfigConsulRepository, error) {
    consulAddress := fmt.Sprintf("%s:%s", os.Getenv("DB"), os.Getenv("DBPORT"))
    config := api.DefaultConfig()
    config.Address = consulAddress

    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }
    return &ConfigConsulRepository{cli: client}, nil
}

func (c *ConfigConsulRepository) CreateConfig(config model.Config) error {
	kv := c.cli.KV()
	key := fmt.Sprintf("configs/%s/%d", config.Name, config.Version)
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	pair := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(pair, nil)
	return err
}

func (c *ConfigConsulRepository) GetConfig(name string, version int) (model.Config, error) {
	kv := c.cli.KV()
	key := fmt.Sprintf("configs/%s/%d", name, version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return model.Config{}, err
	}
	if pair == nil {
		return model.Config{}, fmt.Errorf("config not found")
	}

	var config model.Config
	err = json.Unmarshal(pair.Value, &config)
	return config, err
}

func (c *ConfigConsulRepository) UpdateConfig(name string, version int, newConfig model.Config) error {
	// Ensure the original config exists
	_, err := c.GetConfig(name, version)
	if err != nil {
		return err
	}

	// Create new config
	return c.CreateConfig(newConfig)
}

func (c *ConfigConsulRepository) DeleteConfig(name string, version int) error {
	kv := c.cli.KV()
	key := fmt.Sprintf("configs/%s/%d", name, version)
	_, err := kv.Delete(key, nil)
	return err
}
