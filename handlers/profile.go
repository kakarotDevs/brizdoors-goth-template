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

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX requests, return just the content without the layout
		return profile.IndexContent(*user).Render(r.Context(), w)
	}

	// Full page load - use auth base layout
	return Render(w, r, profile.Index(*user))
}

// HTMX Tab Handlers
func HandleProfileOverview(w http.ResponseWriter, r *http.Request) error {
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
	return profile.OverviewContent(*user).Render(r.Context(), w)
}

func HandleProfileAccount(w http.ResponseWriter, r *http.Request) error {
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
	return profile.AccountContent(*user).Render(r.Context(), w)
}

func HandleProfileSecurity(w http.ResponseWriter, r *http.Request) error {
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
	return profile.SecurityContent(*user).Render(r.Context(), w)
}

func HandleProfilePreferences(w http.ResponseWriter, r *http.Request) error {
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
	return profile.PreferencesContent(*user).Render(r.Context(), w)
}

// Form Submission Handlers
func HandleProfileUpdate(w http.ResponseWriter, r *http.Request) error {
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
		// Parse form data
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		// Update user fields
		if firstName != "" {
			user.FirstName = firstName
		}
		if lastName != "" {
			user.LastName = lastName
		}
		if email != "" {
			user.Email = email
		}
		if phone != "" {
			user.Phone = &phone
		} else {
			user.Phone = nil
		}

		// Save to database
		if err := models.UpdateUser(r.Context(), db.DB, user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return nil
		}

		// Return updated account content
		return profile.AccountContent(*user).Render(r.Context(), w)
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}

func HandleProfilePassword(w http.ResponseWriter, r *http.Request) error {
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
			// Return error message in the security content
			return profile.SecurityContent(*user).Render(r.Context(), w)
		}
		if user.Password != nil && !utils.VerifyPassword(*user.Password, old) {
			// Return error message in the security content
			return profile.SecurityContent(*user).Render(r.Context(), w)
		}
		if err := models.SetUserPassword(r.Context(), db.DB, user.ID, newpw); err != nil {
			// Return error message in the security content
			return profile.SecurityContent(*user).Render(r.Context(), w)
		}
		// Return success message in the security content
		return profile.SecurityContent(*user).Render(r.Context(), w)
	}

	return profile.SecurityContent(*user).Render(r.Context(), w)
}

func HandleProfilePreferencesUpdate(w http.ResponseWriter, r *http.Request) error {
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
		// Parse form data
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		// Handle preferences update logic here
		// For now, just return the updated preferences content
		return profile.PreferencesContent(*user).Render(r.Context(), w)
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}

func HandleProfilePrivacyUpdate(w http.ResponseWriter, r *http.Request) error {
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
		// Parse form data
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		// Handle privacy settings update logic here
		// For now, just return the updated preferences content
		return profile.PreferencesContent(*user).Render(r.Context(), w)
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}

func HandleProfileLanguageUpdate(w http.ResponseWriter, r *http.Request) error {
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
		// Parse form data
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		// Handle language settings update logic here
		// For now, just return the updated preferences content
		return profile.PreferencesContent(*user).Render(r.Context(), w)
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
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
