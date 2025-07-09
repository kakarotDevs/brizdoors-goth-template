package internal

import (
	"net/http"
)

// SetFlash sets a one-time flash message (short-lived cookie)
func SetFlash(w http.ResponseWriter, message string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "flash",
		Value: message,
		Path:  "/",
		MaxAge: 120, // expire after 2 minutes
		HttpOnly: false,
	})
}

// GetFlash retrieves and deletes a flash message cookie
func GetFlash(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("flash")
	if err == nil && cookie.Value != "" {
		// Clear it
		http.SetCookie(w, &http.Cookie{Name: "flash", Value: "", Path: "/", MaxAge: -1})
		return cookie.Value
	}
	return ""
}

