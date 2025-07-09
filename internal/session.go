package internal

import (
	"net/http"
	"time"

	"github.com/kakarotDevs/brizdoors-goth-template/db"
	"github.com/kakarotDevs/brizdoors-goth-template/models"
)

const sessionCookieName = "session_user_id"

// GetUserFromSession reads the user ID from the session cookie.
func GetUserFromSession(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return "", false
	}
	id := cookie.Value
	if _, err := models.GetUserByID(r.Context(), db.DB, id); err != nil {
		return "", false
	}

	return id, true

}

// SetUserSession sets the user ID in a secure cookie.
func SetUserSession(w http.ResponseWriter, userID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    userID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // set to false if developing locally without HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

// ClearUserSession clears the session cookie.
func ClearUserSession(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // same as SetUserSession
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}
