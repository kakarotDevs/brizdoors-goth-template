package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/order"
)

func HandleOrder(w http.ResponseWriter, r *http.Request) error {
	isDarkMode := GetThemeFromRequest(r)
	return Render(w, r, order.Index(isDarkMode))
}

// Authenticated order page — user must be logged in
func HandleAuthOrder(w http.ResponseWriter, r *http.Request) error {

	// Pass user info to the order page if needed (or just render)
	isDarkMode := GetThemeFromRequest(r)
	return Render(w, r, order.Index(isDarkMode))
}
