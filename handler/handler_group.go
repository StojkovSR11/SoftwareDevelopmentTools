package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projekat/model"
	"projekat/services"
	"strconv"

	"github.com/gorilla/mux"
)

type ConfigGroupHandler struct {
	service services.ConfigGroupService
}

func NewConfigGroupHandler(service services.ConfigGroupService) ConfigGroupHandler {
	return ConfigGroupHandler{
		service: service,
	}
}

// GET /configs/{name}/{version}
func (c ConfigGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	// dobavi naziv i verziju
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// pozovi servis metodu
	configGroup, err := c.service.GetConfigGroup(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// vrati odgovor
	resp, err := json.Marshal(configGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content−Type", "application/json")
	w.Write(resp)
}

// Post kreira novu konfiguracionu grupu.
func (c ConfigGroupHandler) Post(w http.ResponseWriter, r *http.Request) {
	// Dekodiraj JSON telo zahteva u strukturu ConfigGroup
	var group model.ConfigGroup
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Pozovi servisnu metodu za kreiranje konfiguracione grupe
	err = c.service.CreateConfigGroup(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Vrati odgovor sa statusom 201 Created
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Konfiguraciona grupa uspešno kreirana.")
}

// Delete briše konfiguracionu grupu po imenu i verziji.
func (c ConfigGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Dobavi ime i verziju konfiguracione grupe iz URL parametara
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Nevalidna verzija", http.StatusBadRequest)
		return
	}

	// Pozovi servisnu metodu za brisanje konfiguracione grupe
	err = c.service.DeleteConfigGroup(name, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Vrati odgovor sa statusom 200 OK
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Konfiguraciona grupa uspešno obrisana.")
}

func (c ConfigGroupHandler) AddConfigToGroup(w http.ResponseWriter, r *http.Request) {
	// Dobavi ime grupe i verziju grupe iz URL parametara
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Nevalidna verzija", http.StatusBadRequest)
		return
	}

	// Dekodiraj JSON telo zahteva u strukturu Config
	var config model.GroupedConfig
	err = json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Pozovi servisnu metodu za dodavanje konfiguracije u konfiguracionu grupu
	err = c.service.AddConfigurationToGroup(name, version, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Vrati odgovor sa statusom 201 Created
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Konfiguracija uspesno dodata u konfiguracionu grupu.")
}

// DeleteConfigFromGroup uklanja konfiguracije iz konfiguracione grupe po imenu, verziji i filteru.
func (c ConfigGroupHandler) DeleteConfigFromGroup(w http.ResponseWriter, r *http.Request) {
	// Get name, version, and filter from URL parameters
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	filter := vars["filter"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	// Call service method to delete configurations by filter
	err = c.service.RemoveConfigurationFromGroup(name, version, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Grouped configurations with filter deleted successfully"))
}

// GetConfigurationsFromGroup dobavlja konfiguracije iz konfiguracione grupe po imenu, verziji i filteru.
func (c ConfigGroupHandler) GetConfigurationsFromGroup(w http.ResponseWriter, r *http.Request) {
	// Get name, version, and filter from URL parameters
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	filter := vars["filter"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	// Call service method to get configurations by filter
	configs, err := c.service.GetConfigurationsFromGroup(name, version, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal configurations to JSON
	response, err := json.Marshal(configs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
