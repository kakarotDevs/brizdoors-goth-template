package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/internal/auth"
	authviews "github.com/kakarotDevs/brizdoors-goth-template/views/auth"
	"github.com/kakarotDevs/brizdoors-goth-template/views/partials"
)

// HandleResetPassword handles the password reset page
func HandleResetPassword(w http.ResponseWriter, r *http.Request) error {
	service := auth.NewPasswordResetService()

	if r.Method == http.MethodGet {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Redirect(w, r, "/forgot-password", http.StatusFound)
			return nil
		}

		// Verify token is valid
		err := service.ValidateResetToken(r.Context(), token)
		if err != nil {
			errorMsg := "Invalid or expired reset link. Please request a new one."
			return Render(w, r, authviews.ResetPasswordPage("", errorMsg))
		}

		return Render(w, r, authviews.ResetPasswordPage(token, ""))
	}

	if r.Method == http.MethodPost {
		token := r.FormValue("token")
		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")

		if token == "" || newPassword == "" || confirmPassword == "" {
			errorMsg := "All fields are required."
			if r.Header.Get("HX-Request") == "true" {
				return Render(w, r, authviews.ResetPasswordSection(token, errorMsg))
			}
			return Render(w, r, authviews.ResetPasswordPage(token, errorMsg))
		}

		if newPassword != confirmPassword {
			errorMsg := "Passwords do not match."
			if r.Header.Get("HX-Request") == "true" {
				return Render(w, r, authviews.ResetPasswordSection(token, errorMsg))
			}
			return Render(w, r, authviews.ResetPasswordPage(token, errorMsg))
		}

		// Reset password using service
		err := service.ResetPassword(r.Context(), token, newPassword)
		if err != nil {
			errorMsg := err.Error()
			if r.Header.Get("HX-Request") == "true" {
				return Render(w, r, authviews.ResetPasswordSection(token, errorMsg))
			}
			return Render(w, r, authviews.ResetPasswordPage(token, errorMsg))
		}

		// Redirect to login with success message
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/login?message=password_reset_success")
			w.WriteHeader(http.StatusOK)
		} else {
			http.Redirect(w, r, "/login?message=password_reset_success", http.StatusFound)
		}
		return nil
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}

// HandleResetPasswordClear clears any error messages from the reset password form
func HandleResetPasswordClear(w http.ResponseWriter, r *http.Request) error {
	token := r.URL.Query().Get("token")
	return Render(w, r, authviews.ResetPasswordSection(token, ""))
}

// HandleResetPasswordNotificationClear clears only the notification area
func HandleResetPasswordNotificationClear(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, partials.Notification("", true))
}
