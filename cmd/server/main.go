package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dharavath-agency/configs"
	deliveryHttp "dharavath-agency/internal/delivery/http"
	"dharavath-agency/internal/delivery/http/handler"
	"dharavath-agency/internal/repository/memory"
	"dharavath-agency/internal/service"
	"dharavath-agency/internal/view"
)

func main() {
	cfg := configs.LoadConfig()

	propRepo, err := memory.NewPropertyRepo()
	if err != nil {
		log.Fatalf("Fatal: failed to initialize property repository: %v", err)
	}
	leadRepo := memory.NewLeadRepo()

	log.Printf("Repository initialized: %d properties, %d agents, %d projects",
		propRepo.Count(), len(propRepo.FindAgents()), len(propRepo.FindProjects()))

	propertyService := service.NewPropertyService(propRepo, propRepo)
	leadService := service.NewLeadService(leadRepo)
	catalogService := service.NewCatalogService(propRepo)
	adminService := service.NewAdminService(propRepo, leadRepo)

	viewEngine, err := view.NewEngine()
	if err != nil {
		log.Fatalf("Fatal: failed to initialize view engine: %v", err)
	}

	pageHandler := handler.NewPageHandler(propertyService, catalogService, viewEngine, cfg.Company)
	propertyHandler := handler.NewPropertyHandler(propertyService, viewEngine, cfg.Company)
	leadHandler := handler.NewLeadHandler(leadService)
	adminHandler := handler.NewAdminHandler(propertyService, adminService, viewEngine, cfg.Company)
	seoHandler := handler.NewSEOHandler(propertyService, catalogService, "https://dharavathagency.in")
	healthHandler := handler.NewHealthHandler(cfg.Env, "1.0.0", propRepo, leadRepo)

	router := deliveryHttp.NewRouter(deliveryHttp.RouterConfig{
		PageHandler:     pageHandler,
		PropertyHandler: propertyHandler,
		LeadHandler:     leadHandler,
		AdminHandler:    adminHandler,
		SEOHandler:      seoHandler,
		HealthHandler:   healthHandler,
	})

	server := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Printf("Dharavath Agency enterprise server running on http://%s (Env: %s)",
			cfg.Addr(), cfg.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Received shutdown signal. Stopping server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server exited successfully.")
}
