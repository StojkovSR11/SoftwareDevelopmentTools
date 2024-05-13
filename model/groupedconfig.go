package model

// Definicija strukture grupisane konfiguracije
type GroupedConfig struct {
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"`
	Labels     map[string]string `json:"labels"`
}
