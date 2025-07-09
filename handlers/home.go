package handlers

import (
	"log/slog"
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/home"
)

func HandleHome(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Serving homepage")

	// Set no-cache headers for login page to prevent form caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Get theme preference
	isLightMode := GetThemeFromRequest(r)

	if isUserLoggedIn(r) {
		// Logged-in users should be redirected to lobby
		http.Redirect(w, r, "/lobby", http.StatusFound)
		return nil
	}

	// Only show login form for non-authenticated users
	return Render(w, r, home.Index(isLightMode))
}
