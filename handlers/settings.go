package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/settings"
)

func HandleSettings(w http.ResponseWriter, r *http.Request) error {
	isDarkMode := GetThemeFromRequest(r)
	return Render(w, r, settings.Index(isDarkMode))
}
