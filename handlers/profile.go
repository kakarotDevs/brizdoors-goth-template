package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/profile"
)

func HandleProfile(w http.ResponseWriter, r *http.Request) error {
	isDarkMode := GetThemeFromRequest(r)
	return Render(w, r, profile.Index(isDarkMode))
}
