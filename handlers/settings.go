package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/db"
	"github.com/kakarotDevs/brizdoors-goth-template/internal"
	"github.com/kakarotDevs/brizdoors-goth-template/models"
	"github.com/kakarotDevs/brizdoors-goth-template/views/settings"
)

func HandleSettings(w http.ResponseWriter, r *http.Request) error {
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}
	_, err := models.GetUserByID(r.Context(), db.DB, userID)
	if err != nil {
		internal.ClearUserSession(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX requests, return just the content without the layout
		return settings.IndexContent().Render(r.Context(), w)
	}

	// Full page load - use auth base layout
	return Render(w, r, settings.Index())
}
