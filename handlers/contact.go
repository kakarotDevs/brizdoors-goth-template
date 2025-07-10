package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/contact"
)

func HandleContact(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, contact.Index())
}
