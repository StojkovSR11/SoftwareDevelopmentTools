package repositories

import (
	"fmt"
	"projekat/model"
)

// Struktura za čuvanje konfiguracija u memoriji
type ConfigInRepositoryMemory struct {
	configs map[string]map[int]model.Config // Mapa sa konfiguracijama, gde ključ predstavlja naziv, a vrednost mapu verzija
}

// Kreiranje nove instance repozitorijuma za konfiguracije u memoriji
func NewConfigInMemory() model.ConfigRepository {
	return &ConfigInRepositoryMemory{
		configs: make(map[string]map[int]model.Config),
	}
}

// Kreiranje konfiguracije.
func (cim *ConfigInRepositoryMemory) CreateConfig(config model.Config) error {
	if cim.configs[config.Name] == nil {
		cim.configs[config.Name] = make(map[int]model.Config)
	}
	if _, exists := cim.configs[config.Name][config.Version];
	// Ako konfiguracija već postoji, vraća grešku.
	exists {
		return fmt.Errorf("konfiguracija sa imenom %s i verzijom %d već postoji", config.Name, config.Version)
	}
	cim.configs[config.Name][config.Version] = config
	return nil
}

// Dobavljanje konfiguracije po imenu i verziji.
func (cim *ConfigInRepositoryMemory) GetConfig(name string, version int) (model.Config, error) {
	config, exists := cim.configs[name][version]
	// Ako konfiguracija već postoji, vraća grešku.
	if !exists {
		return model.Config{}, fmt.Errorf("konfiguracija sa imenom %s i verzijom %d nije pronađena", name, version)
	}
	return config, nil
}

// Ažuriranje konfiguracije po imenu i verziji.
func (cim *ConfigInRepositoryMemory) UpdateConfig(name string, version int, newConfig model.Config) error {
	// Proveri da li postoji konfiguracija sa starim imenom i verzijom
	if _, exists := cim.configs[name][version]; !exists {
		return fmt.Errorf("konfiguracija sa imenom %s i verzijom %d nije pronađena", name, version)
	}

	// Proveri da li postoji konfiguracija sa novim imenom i verzijom
	if _, exists := cim.configs[newConfig.Name][newConfig.Version]; exists {
		return fmt.Errorf("konfiguracija sa imenom %s i verzijom %d već postoji", newConfig.Name, newConfig.Version)
	}

	// Kreiraj novu mapu za novu konfiguraciju
	if cim.configs[newConfig.Name] == nil {
		cim.configs[newConfig.Name] = make(map[int]model.Config)
	}

	// Kopiraj konfiguraciju sa novim imenom i verzijom u novu mapu
	cim.configs[newConfig.Name][newConfig.Version] = newConfig

	// Obriši staru konfiguraciju
	delete(cim.configs[name], version)

	return nil
}

// Brisanje konfiguracije po imenu i verziji.
func (cim *ConfigInRepositoryMemory) DeleteConfig(name string, version int) error {
	if _, exists := cim.configs[name][version];
	// Ako konfiguracija već postoji, vraća grešku.
	!exists {
		return fmt.Errorf("konfiguracija sa imenom %s i verzijom %d nije pronađena", name, version)
	}
	delete(cim.configs[name], version)
	// Ako su sve verzije obrisane, obriši i samu mapu konfiguracija
	if len(cim.configs[name]) == 0 {
		delete(cim.configs, name)
	}
	return nil
}
