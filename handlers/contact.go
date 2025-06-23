package handlers

import (
	"log/slog"
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/contact"
)

func HandleContact(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Serving contactpage")
	return Render(w, r, contact.Index())
}
