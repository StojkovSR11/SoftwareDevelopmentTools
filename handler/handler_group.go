package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projekat/model"
	"projekat/services"
	"strconv"
	"strings"
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
	// Get name and version from URL parameters
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call service method
	configGroup, err := c.service.GetConfigGroup(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return response
	resp, err := json.Marshal(configGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// POST /configs
// Post kreira novu konfiguracionu grupu.
func (c ConfigGroupHandler) Post(w http.ResponseWriter, r *http.Request) {
	// Decode JSON request body into ConfigGroup struct
	var group model.ConfigGroup
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call service method to create config group
	err = c.service.CreateConfigGroup(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response with status 201 Created
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Konfiguraciona grupa uspešno kreirana.")
}

// DELETE /configs/{name}/{version}
// Delete briše konfiguracionu grupu po imenu i verziji.
func (c ConfigGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Get name and version from URL parameters
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Nevalidna verzija", http.StatusBadRequest)
		return
	}

	// Call service method to delete config group
	err = c.service.DeleteConfigGroup(name, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return response with status 200 OK
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Konfiguraciona grupa uspešno obrisana.")
}

// POST /configs/{name}/{version}/config
// AddConfigToGroup dodaje konfiguraciju u konfiguracionu grupu po imenu i verziji.
func (c ConfigGroupHandler) AddConfigToGroup(w http.ResponseWriter, r *http.Request) {
	// Get group name and version from URL parameters
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Nevalidna verzija", http.StatusBadRequest)
		return
	}

	// Decode JSON request body into GroupedConfig struct
	var config model.GroupedConfig
	err = json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call service method to add configuration to group
	err = c.service.AddConfigurationToGroup(name, version, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response with status 201 Created
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Konfiguracija uspešno dodata u konfiguracionu grupu.")
}

// DELETE /configs/{name}/{version}/config/{filter}
// DeleteConfigFromGroup uklanja konfiguracije iz konfiguracione grupe po imenu, verziji i filteru.
func (c ConfigGroupHandler) DeleteConfigFromGroup(w http.ResponseWriter, r *http.Request) {
	// Get name, version, and filter from URL parameters
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]
	filter := vars["filter"]

	// Split filter into key and value
	filterParts := strings.SplitN(filter, ":", 2)
	if len(filterParts) != 2 {
		http.Error(w, "Invalid filter format. Expected key:value", http.StatusBadRequest)
		return
	}
	filterKey := filterParts[0]
	filterValue := filterParts[1]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	// Call service method to delete configurations by filter
	err = c.service.RemoveConfigurationFromGroup(name, version, filterKey, filterValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Konfiguracije sa filterom uspešno obrisane."))
}


func (c ConfigGroupHandler) GetConfigurationsFromGroup(w http.ResponseWriter, r *http.Request) {
    // Get name, version, and filter from URL parameters
    vars := mux.Vars(r)
    name := vars["name"]
    versionStr := vars["version"]
    filter := vars["filter"]

    // Split filter into key and value
    filterParts := strings.SplitN(filter, ":", 2)
    if len(filterParts) != 2 {
        http.Error(w, "Invalid filter format. Expected key:value", http.StatusBadRequest)
        return
    }
    filterKey := filterParts[0]
    filterValue := filterParts[1]

    version, err := strconv.Atoi(versionStr)
    if err != nil {
        http.Error(w, "Invalid version", http.StatusBadRequest)
        return
    }

    // Call service method to get configurations by filter
    configs, err := c.service.GetConfigurationsFromGroup(name, version, filterKey, filterValue)
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
