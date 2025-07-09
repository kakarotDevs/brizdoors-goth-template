package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/kakarotDevs/brizdoors-goth-template/db"
	"github.com/kakarotDevs/brizdoors-goth-template/handlers"
	"github.com/kakarotDevs/brizdoors-goth-template/internal/auth"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load .env file", "error", err)
		os.Exit(1)
	}

	// Run DB migrations before connecting ORM or starting app
	m, err := migrate.New(
		"file://migrations",
		os.Getenv("DEV_DATABASE_URL"),
	)
	if err != nil {
		log.Fatalf("migration setup failed: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	// Connect to Postgres using Bun+pgxpool
	if err := db.InitDB(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.CloseDB()
	log.Println("Connected to database!")

	// Initialize Google OAuth config after env vars loaded
	// Smart callback URL detection
	baseURL := determineBaseURL()
	callbackURL := baseURL + "/auth/google/callback"

	auth.InitGoogleOauth(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		callbackURL,
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
	r.Post("/chat", handlers.Make(handlers.ChatHandler))
	r.Post("/chat/message", handlers.Make(handlers.HandleChatMessage))
	r.Post("/chat/new", handlers.Make(handlers.HandleNewChat))

	// Theme
	r.Post("/toggle-theme", handlers.Make(handlers.ThemeToggleHandler))

	// Navigation
	r.Get("/nav/mobile", handlers.Make(handlers.HandleMobileNav))
	r.Get("/nav/mobile-close", handlers.Make(handlers.HandleMobileNavClose))

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
		r.Get("/chat", handlers.Make(handlers.HandleChatPage))
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

// determineBaseURL intelligently determines the base URL for OAuth callbacks
func determineBaseURL() string {
	// Priority 1: Check if FORCE_LOCALHOST is set (bypasses ngrok detection)
	if forceLocal := os.Getenv("FORCE_LOCALHOST"); forceLocal == "true" {
		slog.Info("FORCE_LOCALHOST enabled, using localhost for OAuth callbacks")
		return "http://localhost:3000"
	}

	// Priority 2: Check if NGROK_URL is explicitly set
	if ngrokURL := os.Getenv("NGROK_URL"); ngrokURL != "" {
		slog.Info("Using explicitly set NGROK_URL for OAuth callbacks", "url", ngrokURL)
		return ngrokURL
	}

	// Priority 3: Check for production domain
	if prodDomain := os.Getenv("PRODUCTION_DOMAIN"); prodDomain != "" {
		slog.Info("Using production domain for OAuth callbacks", "domain", prodDomain)
		return prodDomain
	}

	// Default: Use localhost for development
	slog.Info("No specific URL configured, defaulting to localhost for OAuth callbacks")
	return "http://localhost:3000"
}
