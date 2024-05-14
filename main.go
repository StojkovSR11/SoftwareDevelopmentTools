package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlers "projekat/handler"
	"projekat/repositories"
	"projekat/services"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func main() {
	repo := repositories.NewConfigInMemory()
	service := services.NewConfigInService(repo)

	repogroup := repositories.NewConfigGroupInMemoryRepository()
	servicegroup := services.NewConfigGroupInService(repogroup)

	configHandler := handlers.NewConfigHandler(service)
	configGroupHandler := handlers.NewConfigGroupHandler(servicegroup)

	router := mux.NewRouter()

	// Definisanje rate limitera sa limitom od 10 zahteva po min
	limiter := rate.NewLimiter(rate.Limit(0.167), 3)

	// Middleware koji implementira rate limiter
	rateLimitMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// Dodavanje middleware-a na sver route
	router.Use(rateLimitMiddleware)

	router.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")
	router.HandleFunc("/configs", configHandler.Post).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Put).Methods("PUT")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE")

	router.HandleFunc("/configGroups/{name}/{version}", configGroupHandler.Get).Methods("GET")
	router.HandleFunc("/configGroups", configGroupHandler.Post).Methods("POST")
	router.HandleFunc("/configGroups/{name}/{version}", configGroupHandler.Delete).Methods("DELETE")
	//router.HandleFunc("/configGroups/{name}/{version}/{configName}", configGroupHandler.DeleteConfigFromGroup).Methods("DELETE")
	router.HandleFunc("/configGroups/{name}/{version}/addConfig", configGroupHandler.AddConfigToGroup).Methods("POST")

	router.HandleFunc("/configGroups/{name}/{version}/{filter}", configGroupHandler.GetConfigurationsFromGroup).Methods("GET")
	router.HandleFunc("/configGroups/{name}/{version}/{filter}", configGroupHandler.DeleteConfigFromGroup).Methods("DELETE")
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}

	// Kanal za hvatanje signala za zaustavljanje
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Pokretanje servera u posebnoj gorutini
	go func() {
		log.Println("Server se pokrece...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Greška prilikom pokretanja servera: %v", err)
		}
	}()
	// Čekanje na signal zaustavljanja
	<-shutdown
	// Logovanje početka procesa graceful shutdown-a
	log.Println("Zatvaranje servera...")

	// Pravljenje kanala za oznaku zatvaranja servera
	stop := make(chan struct{})
	go func() {
		// Postavljanje timeout-a za graceful shutdown
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Zatvaranje HTTP servera
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Greška prilikom graceful shutdown-a servera: %v", err)
		}
		close(stop)
	}()

	// Čekanje na zatvaranje servera ili prekid izvršavanja
	<-stop
	// Logovanje završetka procesa graceful shutdown-a
	log.Println("Završeno zatvaranje servera")
}
