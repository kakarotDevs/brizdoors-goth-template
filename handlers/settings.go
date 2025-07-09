package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/db"
	"github.com/kakarotDevs/brizdoors-goth-template/internal"
	"github.com/kakarotDevs/brizdoors-goth-template/models"
	"github.com/kakarotDevs/brizdoors-goth-template/views/settings"
)

func HandleSettings(w http.ResponseWriter, r *http.Request) error {
	isDarkMode := GetThemeFromRequest(r)
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}
	user, err := models.GetUserByID(r.Context(), db.DB, userID)
	if err != nil {
		internal.ClearUserSession(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}
	_ = user // Placeholder to prevent unused variable error
	return Render(w, r, settings.Index(isDarkMode))
}
