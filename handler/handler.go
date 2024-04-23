package handlers

import (
	"projekat/model"
	"fmt"
	"encoding/json"
	"net/http"
	"projekat/services"
	"strconv"
	"github.com/gorilla/mux"
    "time"
)

type ConfigHandler struct {
	service services.ConfigService
}

func NewConfigHandler(service services.ConfigService) ConfigHandler {
	return ConfigHandler{
		service: service,
	}
}


// GET /configs/{name}/{version}
func (c ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
    // dobavi naziv i verziju
    name := mux.Vars(r)["name"]
    version := mux.Vars(r)["version"]
    versionInt, err := strconv.Atoi(version)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // pozovi servis metodu
    config, err := c.service.GetConfig(name, versionInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    // čekaj 5 sekundi
    time.Sleep(5 * time.Second)

    // vrati odgovor
    resp, err := json.Marshal(config)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(resp)
}


func (c ConfigHandler) Post(w http.ResponseWriter, r *http.Request) {
    // Dekodiraj JSON telo zahteva u strukturu Config
    var config model.Config
    err := json.NewDecoder(r.Body).Decode(&config)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Pozovi servis metodu za kreiranje konfiguracije
    err = c.service.CreateConfig(config)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Vrati odgovor sa statusom 201 Created
    w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Konfiguracija uspešno kreirana.")
}


// UpdateConfig ažurira postojeću konfiguraciju po imenu i verziji.
// Ako konfiguracija ne postoji, vraća grešku. Ako nova verzija već postoji, takođe vraća grešku.
func (c ConfigHandler) Put(w http.ResponseWriter, r *http.Request) {
    // Dobavi ime i verziju konfiguracije iz URL parametara
    vars := mux.Vars(r)
    name := vars["name"]
    versionStr := vars["version"]
    version, err := strconv.Atoi(versionStr)
    if err != nil {
        http.Error(w, "Nevalidna verzija", http.StatusBadRequest)
        return
    }

    // Dekodiraj JSON telo zahteva u strukturu Config
    var newConfig model.Config
    err = json.NewDecoder(r.Body).Decode(&newConfig)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Proveri da li nova verzija već postoji
    _, err = c.service.GetConfig(name, newConfig.Version)
    if err == nil {
        http.Error(w, fmt.Sprintf("Konfiguracija sa imenom %s i verzijom %d već postoji", name, newConfig.Version), http.StatusBadRequest)
        return
    }


    // Pozovi servisnu metodu za ažuriranje konfiguracije
    err = c.service.UpdateConfig(name, version, newConfig)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    // Vrati odgovor sa statusom 200 OK
    w.WriteHeader(http.StatusOK)
    // Pošalji poruku da je ažuriranje uspešno
    fmt.Fprintf(w, "Konfiguracija uspešno ažurirana.")
}



func (c ConfigHandler) Delete(w http.ResponseWriter, r *http.Request) {
    // Dobavi ime i verziju konfiguracije iz URL parametara
    vars := mux.Vars(r)
    name := vars["name"]
    versionStr := vars["version"]
    version, err := strconv.Atoi(versionStr)
    if err != nil {
        http.Error(w, "Nevalidna verzija", http.StatusBadRequest)
        return
    }

    // Pozovi servisnu metodu za brisanje konfiguracije
    err = c.service.DeleteConfig(name, version)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    //status odgovora na 200 OK
    w.WriteHeader(http.StatusOK)
	// Pošalji poruku da je brisanje uspešno
	fmt.Fprintf(w, "Konfiguracija uspešno obrisana.")
}
