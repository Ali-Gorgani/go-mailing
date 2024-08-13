package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ali-Gorgani/go-mailing/internal/config"
	"github.com/Ali-Gorgani/go-mailing/internal/handler"
	"github.com/Ali-Gorgani/go-mailing/internal/loadbalancer"
	"github.com/Ali-Gorgani/go-mailing/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Load configuration
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Initialize services
	mailService := service.NewMailService()

	// Set up HTTP handlers with Chi router
	r := chi.NewRouter()
	r.Get("/get_providers", handler.GetProviders(mailService))
	r.Post("/send_mail", handler.SendMail(mailService))
	r.Handle("/swagger/*", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))))

	r.Get("/test", Test)

	// Set up Load Balancer
	serverPool := loadbalancer.NewServerPool(config.Server.LoadbalancerAddresses)
	r.HandleFunc("/*", serverPool.ProxyHandler)

	// Start a server
	fmt.Printf("Starting server on %s...\n", config.Server.Address)
	log.Fatal(http.ListenAndServe(config.Server.Address, r))
}

func Test(w http.ResponseWriter, r *http.Request) {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	fmt.Fprintf(w, "Handled by server on port %s\n", config.Server.Address)
	fmt.Printf("Handled by server on port %s\n", config.Server.Address)
	time.Sleep(50 * time.Millisecond)
}
