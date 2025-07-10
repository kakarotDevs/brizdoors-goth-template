package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/about"
)

func HandleAbout(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, about.Index())
}
