package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/handler"
	"github.com/stwalsh4118/hephaestus/backend/internal/middleware"
	"github.com/stwalsh4118/hephaestus/backend/internal/storage"
)

const (
	defaultPort     = "8080"
	shutdownTimeout = 5 * time.Second
	readTimeout     = 15 * time.Second
	writeTimeout    = 15 * time.Second
	idleTimeout     = 60 * time.Second
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	store, err := storage.NewFileStore("")
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}

	// Initialize Docker orchestrator (non-fatal if Docker is unavailable).
	// pollingCtx controls any background health polling; cancel it before teardown.
	pollingCtx, cancelPolling := context.WithCancel(context.Background())
	_ = pollingCtx // will be passed to StartHealthPolling when deploy flow is wired

	dockerClient, dockerErr := docker.NewClient()
	var orchestrator *docker.DockerOrchestrator
	if dockerErr != nil {
		log.Printf("docker client unavailable: %v (orchestration features disabled)", dockerErr)
	} else {
		orchestrator = docker.NewDockerOrchestrator(dockerClient)
		log.Println("docker orchestrator initialized")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	diagramHandler := handler.NewDiagramHandler(store)
	diagramHandler.RegisterRoutes(mux)

	wsHandler := handler.NewWebSocketHandler()
	wsHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      middleware.CORS()(mux),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	errs := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", server.Addr)
		errs <- server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("received signal %s; shutting down", sig)
	case err := <-errs:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	// Cancel health polling before teardown to stop background goroutines.
	cancelPolling()

	// Teardown Docker resources after HTTP server has drained.
	if orchestrator != nil {
		log.Println("tearing down docker resources...")
		if err := orchestrator.TeardownAll(ctx); err != nil {
			log.Printf("docker teardown errors: %v", err)
		} else {
			log.Println("docker teardown complete")
		}
		if dockerClient != nil {
			if err := dockerClient.Close(); err != nil {
				log.Printf("failed to close docker client: %v", err)
			}
		}
	}

	log.Println("server stopped")
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	response := healthResponse{Status: "ok"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to write health response: %v", err)
	}
}
