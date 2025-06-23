package handlers

import (
	"log/slog"
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/order"
)

func HandleOrder(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Serving orderpage")
	return Render(w, r, order.Index())
}
