package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/utils"
	authviews "github.com/kakarotDevs/brizdoors-goth-template/views/auth"
	"github.com/kakarotDevs/brizdoors-goth-template/views/partials"
)

// FormHandler handles form processing and validation
type FormHandler struct {
	validator *utils.FormValidator
}

// NewFormHandler creates a new form handler
func NewFormHandler() *FormHandler {
	return &FormHandler{
		validator: utils.NewFormValidator(),
	}
}

// HandleLoginForm processes login form submission
func (fh *FormHandler) HandleLoginForm(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return fh.renderLoginError(w, r, "Invalid form data")
	}

	contact := utils.GetFormValue(r, "contact")
	password := utils.GetFormValue(r, "password")

	// Validate inputs
	if valid, msg := fh.validator.ValidateEmailOrPhone(contact); !valid {
		return fh.renderLoginError(w, r, msg)
	}

	if valid, msg := fh.validator.ValidatePassword(password); !valid {
		return fh.renderLoginError(w, r, msg)
	}

	// Call the existing login logic from auth.go
	return HandleLogin(w, r)
}

// HandleForgotPasswordForm processes forgot password form submission
func (fh *FormHandler) HandleForgotPasswordForm(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	contact := utils.GetFormValue(r, "contact")

	// Validate input
	if valid, msg := fh.validator.ValidateEmailOrPhone(contact); !valid {
		return fh.renderForgotPasswordError(w, r, msg)
	}

	// Call the existing forgot password logic from forgot_password.go
	return HandleForgotPassword(w, r)
}

// HandleResetPasswordForm processes reset password form submission
func (fh *FormHandler) HandleResetPasswordForm(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	token := utils.GetFormValue(r, "token")
	newPassword := utils.GetFormValue(r, "new_password")
	confirmPassword := utils.GetFormValue(r, "confirm_password")

	// Validate inputs
	if valid, msg := fh.validator.ValidateRequired(token, "Token"); !valid {
		return fh.renderResetPasswordError(w, r, token, msg)
	}

	if valid, msg := fh.validator.ValidatePassword(newPassword); !valid {
		return fh.renderResetPasswordError(w, r, token, msg)
	}

	if valid, msg := fh.validator.ValidatePasswordMatch(newPassword, confirmPassword); !valid {
		return fh.renderResetPasswordError(w, r, token, msg)
	}

	// Call the existing reset password logic from reset_password.go
	return HandleResetPassword(w, r)
}

// ClearNotification clears only the notification area for a specific form
func (fh *FormHandler) ClearNotification(w http.ResponseWriter, r *http.Request) error {
	// Return empty notification area - client-side JS handles the clearing
	return Render(w, r, partials.Notification("", true))
}

// renderLoginError renders login form with error
func (fh *FormHandler) renderLoginError(w http.ResponseWriter, r *http.Request, errorMsg string) error {
	if utils.IsHTMXRequest(r) {
		return Render(w, r, authviews.LoginSection(errorMsg))
	}
	http.Error(w, errorMsg, http.StatusBadRequest)
	return nil
}

// renderForgotPasswordError renders forgot password form with error
func (fh *FormHandler) renderForgotPasswordError(w http.ResponseWriter, r *http.Request, errorMsg string) error {
	if utils.IsHTMXRequest(r) {
		return Render(w, r, authviews.ForgotPasswordSection(errorMsg))
	}
	http.Error(w, errorMsg, http.StatusBadRequest)
	return nil
}

// renderResetPasswordError renders reset password form with error
func (fh *FormHandler) renderResetPasswordError(w http.ResponseWriter, r *http.Request, token, errorMsg string) error {
	if utils.IsHTMXRequest(r) {
		return Render(w, r, authviews.ResetPasswordSection(token, errorMsg))
	}
	http.Error(w, errorMsg, http.StatusBadRequest)
	return nil
}
