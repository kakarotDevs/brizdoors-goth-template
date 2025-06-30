package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/internal"
	"github.com/kakarotDevs/brizdoors-goth-template/models"
	"github.com/kakarotDevs/brizdoors-goth-template/views/order"
)

// Public order page — no login required
func HandleOrder(w http.ResponseWriter, r *http.Request) error {
	isDarkMode := GetThemeFromRequest(r)
	return Render(w, r, order.Index(isDarkMode))
}

// Authenticated order page — user must be logged in
func HandleAuthOrder(w http.ResponseWriter, r *http.Request) error {
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		// Not logged in, redirect to sign-in page
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		internal.ClearUserSession(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	_ = user // Placeholder to prevent unused variable error

	// Pass user info to the order page if needed (or just render)
	isDarkMode := GetThemeFromRequest(r)
	return Render(w, r, order.Index(isDarkMode))
}
