package main

import (
	"chat-application/users/api/routes"
	"chat-application/users/config"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	env, err := config.NewEnvDev()
	if err != nil {
		log.Fatal("[server] - failed to load env")
	}

	router := routes.NewRoute()

	server := &http.Server{
		Addr:           ":" + env.APP_PORT,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// OS interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Printf("[server] - Server is starting on port %s", env.APP_PORT)
		log.Printf("[server] - Server is running on http://localhost:%s", env.APP_PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[server] - Failed to run server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop

	log.Println("[server] - Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[server] - Server forced to shutdown: %v", err)
	}

	log.Println("[server] - Server gracefully stopped")
}
