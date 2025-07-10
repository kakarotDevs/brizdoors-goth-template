package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/order"
)

func HandleOrder(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, order.Index())
}

// Authenticated order page — user must be logged in
func HandleAuthOrder(w http.ResponseWriter, r *http.Request) error {

	// Pass user info to the order page if needed (or just render)
	return Render(w, r, order.Index())
}
