package model

// Definicija strukture konfiguracije
type Config struct {
	Name       string             `json:"name"`
	Version    int                `json:"version"`
	Parameters map[string]string `json:"parameters"`
}

// Interfejs za rad sa konfiguracijama
type ConfigRepository interface {

	CreateConfig(config Config) error

	GetConfig(name string, version int) (Config, error)

	UpdateConfig(name string, version int, newConfig Config) error

	DeleteConfig(name string, version int) error
    
}
