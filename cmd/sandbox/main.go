package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "main/internal"
	"main/internal/database"
	"main/internal/sandbox/core"
	jwtutil "main/internal/security/jwt"
)

func main() {
	// docker client
	// client.FromEnv == reads docker connection string from environment

	ctx := context.Background()

	apiClient, err := core.NewSandboxClient()
	if err != nil {
		panic(err)
	}

	defer apiClient.Close()

	db, err := database.ConnectFromEnv(ctx)
	if err != nil {
		panic(err)
	}

	if err := jwtutil.InitFromEnv(); err != nil {
		panic(err)
	}

	application, err := app.New(db, apiClient)
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Addr:         ":8080",
		Handler:      application.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	//run seperately from the server
	go func() {
		log.Println("server listening on port: 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("error starting server: ", err)
		}
	}()

	//waits for ctr+c macro to shutdown the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	<-sigChan

	log.Println("Shutting down the server")

	shutDownCtr, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutDownCtr); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	apiClient.CleanUp(ctx)

	log.Println("Server stopped")
}
