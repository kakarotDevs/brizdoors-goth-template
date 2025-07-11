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

	// Initialize form handler
	formHandler := handlers.NewFormHandler()

	// Public Routes
	r.Get("/", handlers.Make(handlers.HandleHome))
	r.Get("/login", handlers.Make(handlers.HandleLogin))
	r.Get("/logout", handlers.Make(handlers.HandleLogout))

	// Form handling routes
	r.Post("/form/login", handlers.Make(formHandler.HandleLoginForm))
	r.Post("/form/forgot-password", handlers.Make(formHandler.HandleForgotPasswordForm))
	r.Post("/form/reset-password", handlers.Make(formHandler.HandleResetPasswordForm))

	// Demo Routes (for testing spinners)
	r.Get("/demo", handlers.Make(handlers.HandleSpinnerDemo))
	r.Post("/demo/loading", handlers.Make(handlers.HandleDemoLoading))
	r.Post("/demo/form", handlers.Make(handlers.HandleDemoForm))
	r.Get("/demo/content", handlers.Make(handlers.HandleDemoContent))
	r.Get("/demo/clear", handlers.Make(handlers.HandleDemoClear))

	// Password Reset Routes (page display only)
	r.Get("/forgot-password", handlers.Make(handlers.HandleForgotPassword))
	r.Get("/reset-password", handlers.Make(handlers.HandleResetPassword))

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

		// Profile HTMX Routes
		r.Get("/profile/overview", handlers.Make(handlers.HandleProfileOverview))
		r.Get("/profile/account", handlers.Make(handlers.HandleProfileAccount))
		r.Get("/profile/security", handlers.Make(handlers.HandleProfileSecurity))
		r.Get("/profile/preferences", handlers.Make(handlers.HandleProfilePreferences))

		// Profile Form Submission Routes
		r.Post("/profile/update", handlers.Make(handlers.HandleProfileUpdate))
		r.Post("/profile/password", handlers.Make(handlers.HandleProfilePassword))
		r.Post("/profile/preferences", handlers.Make(handlers.HandleProfilePreferencesUpdate))
		r.Post("/profile/privacy", handlers.Make(handlers.HandleProfilePrivacyUpdate))
		r.Post("/profile/language", handlers.Make(handlers.HandleProfileLanguageUpdate))
		r.Post("/profile/delete", handlers.Make(handlers.HandleProfileDelete))
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
