package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/db"
	"github.com/kakarotDevs/brizdoors-goth-template/internal"
	"github.com/kakarotDevs/brizdoors-goth-template/models"
	"github.com/kakarotDevs/brizdoors-goth-template/utils"
	"github.com/kakarotDevs/brizdoors-goth-template/views/profile"
)

func HandleProfile(w http.ResponseWriter, r *http.Request) error {
	isDarkMode := GetThemeFromRequest(r)
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	user, err := models.GetUserByID(r.Context(), db.DB, userID)
	if err != nil {
		internal.ClearUserSession(w, r)
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	return Render(w, r, profile.Index(*user, isDarkMode))
}

func HandleProfilePassword(w http.ResponseWriter, r *http.Request) error {
	// htmx handler for in-place password change
	isDarkMode := GetThemeFromRequest(r)
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}
	user, err := models.GetUserByID(r.Context(), db.DB, userID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	if r.Method == http.MethodPost {
		old := r.FormValue("old_password")
		newpw := r.FormValue("new_password")
		conf := r.FormValue("confirm_password")
		if newpw == "" || newpw != conf {
			return Render(w, r, profile.PasswordForm(true, "Passwords do not match", isDarkMode))
		}
		if user.Password != nil && !utils.VerifyPassword(*user.Password, old) {
			return Render(w, r, profile.PasswordForm(true, "Current password incorrect", isDarkMode))
		}
		if err := models.SetUserPassword(r.Context(), db.DB, user.ID, newpw); err != nil {
			return Render(w, r, profile.PasswordForm(true, "Failed to set new password", isDarkMode))
		}
		return Render(w, r, profile.PasswordForm(false, "Password changed successfully!", isDarkMode))
	}
	return Render(w, r, profile.PasswordForm(false, "", isDarkMode))
}

func HandleProfileDelete(w http.ResponseWriter, r *http.Request) error {
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
	if err := models.DeleteUser(r.Context(), db.DB, userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	internal.ClearUserSession(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}
