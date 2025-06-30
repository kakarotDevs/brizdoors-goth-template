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
	"github.com/kakarotDevs/brizdoors-goth-template/internal/auth"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load .env file", "error", err)
		os.Exit(1)
	}

	// Initialize Google OAuth config after env vars loaded
	// Using ngrok URL for mobile testing
	auth.InitGoogleOauth(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		"https://6582-220-235-205-48.ngrok-free.app/auth/google/callback",
	)

	r := chi.NewMux()
	r.Use(middleware.Recoverer)

	r.Handle("/public/*", public())

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Public Routes
	r.Get("/", handlers.Make(handlers.HandleHome))
	r.Get("/about", handlers.Make(handlers.HandleAbout))
	r.Get("/contact", handlers.Make(handlers.HandleContact))
	r.Get("/order", handlers.Make(handlers.HandleOrder))
	r.Get("/register", handlers.Make(handlers.HandleRegister))
	r.Get("/login", handlers.Make(handlers.HandleLogin))
	r.Get("/logout", handlers.Make(handlers.HandleLogout))
	r.Post("/login", handlers.Make(handlers.HandleLogin))
	r.Post("/register", handlers.Make(handlers.HandleRegister))

	// Theme
	r.Post("/toggle-theme", handlers.Make(handlers.ThemeToggleHandler))

	// Auth
	r.Get("/auth/google", handlers.Make(handlers.HandleGoogleLogin))
	r.Get("/auth/google/callback", handlers.Make(handlers.HandleGoogleCallback))
	r.Get("/auth/menu", handlers.Make(handlers.HandleAuthMenu))
	r.Get("/auth/menu/toggle", handlers.Make(handlers.HandleAuthMenuToggle))
	r.Get("/auth/menu/content", handlers.Make(handlers.HandleAuthMenuContent))
	

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth)
		r.Get("/lobby", handlers.Make(handlers.HandleLobby))
		r.Get("/profile", handlers.Make(handlers.HandleProfile))
		r.Get("/settings", handlers.Make(handlers.HandleSettings))
		r.Post("/logout", handlers.Make(handlers.HandleLogout))
	})

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":3000"
	}

	slog.Info("Starting server", "localhost", fmt.Sprintf("http://localhost%s", listenAddr), "network", "http://192.168.86.41:3000")

	if err := http.ListenAndServe(listenAddr, r); err != nil {
		slog.Error("Server exited with error", "error", err)
		os.Exit(1)
	}
}
