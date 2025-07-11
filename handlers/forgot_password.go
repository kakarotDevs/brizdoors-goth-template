package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/kakarotDevs/brizdoors-goth-template/internal/auth"
	authviews "github.com/kakarotDevs/brizdoors-goth-template/views/auth"
	"github.com/kakarotDevs/brizdoors-goth-template/views/partials"
)

// HandleForgotPassword handles the forgot password page
func HandleForgotPassword(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return Render(w, r, authviews.ForgotPasswordPage(""))
	}

	if r.Method == http.MethodPost {
		contact := strings.TrimSpace(r.FormValue("contact"))

		// Validation
		if contact == "" {
			errorMsg := "Email address or phone number is required."
			if r.Header.Get("HX-Request") == "true" {
				return Render(w, r, authviews.ForgotPasswordSection(errorMsg))
			}
			return Render(w, r, authviews.ForgotPasswordPage(errorMsg))
		}

		// Check if it's an email or phone
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

		isEmail := emailRegex.MatchString(contact)
		isPhone := phoneRegex.MatchString(contact)

		if !isEmail && !isPhone {
			errorMsg := "Please enter a valid email address or phone number."
			if r.Header.Get("HX-Request") == "true" {
				return Render(w, r, authviews.ForgotPasswordSection(errorMsg))
			}
			return Render(w, r, authviews.ForgotPasswordPage(errorMsg))
		}

		// Use the password reset service
		service := auth.NewPasswordResetService()
		_ = service.CreatePasswordReset(r.Context(), contact)

		// For security, don't reveal if user exists - just redirect to login
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/login")
			w.WriteHeader(http.StatusOK)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return nil
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}

// HandleForgotPasswordClear clears any error messages from the forgot password form
func HandleForgotPasswordClear(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, authviews.ForgotPasswordSection(""))
}

// HandleForgotPasswordNotificationClear clears only the notification area
func HandleForgotPasswordNotificationClear(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, partials.Notification("", true))
}
