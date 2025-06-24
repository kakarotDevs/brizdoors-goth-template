package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/internal"
	"github.com/kakarotDevs/brizdoors-goth-template/models"
	"github.com/kakarotDevs/brizdoors-goth-template/views/lobby"
)

func HandleLobby(w http.ResponseWriter, r *http.Request) error {
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		internal.ClearUserSession(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	return Render(w, r, lobby.Index(user.Name))
}
