package handlers

import (
	"log/slog"
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/home"
)

func HandleHome(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Serving homepage")
	return Render(w, r, home.Index())
}
