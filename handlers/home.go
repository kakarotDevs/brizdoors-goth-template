package handlers

import (
	"log/slog"
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/internal"
	"github.com/kakarotDevs/brizdoors-goth-template/models"
	"github.com/kakarotDevs/brizdoors-goth-template/views/home"
)

func HandleHome(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Serving homepage")

	if isUserLoggedIn(r) {
		userID, _ := internal.GetUserFromSession(r)
		user, err := models.GetUserByID(userID)
		if err != nil {
			internal.ClearUserSession(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		}
		return Render(w, r, home.Index(true, user.Name, user.Email, user.Picture))
	} else {
		return Render(w, r, home.Index(false, "", "", ""))
	}

}
