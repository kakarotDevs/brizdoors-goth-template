package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/kakarotDevs/brizdoors-goth-template/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load .env file", "error", err)
		os.Exit(1)
	}

	r := chi.NewMux()

	// Middleware should come before route declarations
	r.Use(middleware.Recoverer)

	r.Handle("/public/*", public())

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	r.Get("/", handlers.Make(handlers.HandleHome))
	r.Get("/about", handlers.Make(handlers.HandleAbout))
	r.Get("/contact", handlers.Make(handlers.HandleContact))
	r.Get("/order", handlers.Make(handlers.HandleOrder))

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":3000"
	}

	slog.Info("Starting server", "url", fmt.Sprintf("http://localhost%s", listenAddr))

	if err := http.ListenAndServe(listenAddr, r); err != nil {
		slog.Error("Server exited with error", "error", err)
		os.Exit(1)
	}
}
