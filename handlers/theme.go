package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/partials"
)

// ThemeToggleHandler handles the theme toggle requests
func ThemeToggleHandler(w http.ResponseWriter, r *http.Request) error {
	// Get current theme from cookie
	cookie, err := r.Cookie("theme")
	var isLight bool
	
	if err != nil || cookie.Value == "dark" {
		// If no cookie or dark theme (default), switch to light
		isLight = true
		http.SetCookie(w, &http.Cookie{
			Name:     "theme",
			Value:    "light",
			Path:     "/",
			MaxAge:   86400 * 365, // 1 year
			HttpOnly: false,       // Allow JavaScript access for potential future use
			Secure:   false,       // Set to true in production with HTTPS
			SameSite: 2, // SameSiteLax
		})
	} else {
		// If light theme, switch back to dark (default)
		isLight = false
		http.SetCookie(w, &http.Cookie{
			Name:     "theme",
			Value:    "dark",
			Path:     "/",
			MaxAge:   86400 * 365, // 1 year
			HttpOnly: false,
			Secure:   false,
			SameSite: 2, // SameSiteLax
		})
	}

	// Set HTMX header to trigger a page refresh for theme change
	w.Header().Set("HX-Refresh", "true")
	
	// Return the updated theme toggle button
	component := partials.ThemeToggle(isLight)
	return component.Render(r.Context(), w)
}

// GetThemeFromRequest extracts the current theme preference from the request
// Returns true if light mode, false if dark mode (default)
func GetThemeFromRequest(r *http.Request) bool {
	cookie, err := r.Cookie("theme")
	if err != nil {
		// Default to dark theme if no cookie
		return false
	}
	return cookie.Value == "light"
}
