package utils

import (
	"net/http"
	"regexp"
	"strings"
)

// FormValidator handles form validation logic
type FormValidator struct{}

// NewFormValidator creates a new form validator instance
func NewFormValidator() *FormValidator {
	return &FormValidator{}
}

// ValidateEmailOrPhone validates if the input is a valid email or phone number
func (fv *FormValidator) ValidateEmailOrPhone(contact string) (bool, string) {
	contact = strings.TrimSpace(contact)

	if contact == "" {
		return false, "Email address or phone number is required."
	}

	// Check if it's an email or phone
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	isEmail := emailRegex.MatchString(contact)
	if isEmail {
		return true, ""
	}

	// For phone numbers, clean and validate as Australian number
	cleanPhone := fv.cleanPhoneNumber(contact)
	if valid, msg := fv.validateAustralianPhone(cleanPhone); !valid {
		return false, msg
	}

	return true, ""
}

// cleanPhoneNumber removes all whitespace and formatting characters
func (fv *FormValidator) cleanPhoneNumber(phone string) string {
	// Remove all whitespace, dashes, parentheses, and plus signs
	re := regexp.MustCompile(`[\s\-\(\)\+]`)
	return re.ReplaceAllString(phone, "")
}

// validateAustralianPhone validates Australian phone number formats
func (fv *FormValidator) validateAustralianPhone(phone string) (bool, string) {
	// Australian mobile numbers: 04XX XXX XXX (10 digits starting with 04)
	// Australian landline: 0X XXXX XXXX (10 digits starting with 02, 03, 07, 08)
	// International format: +61 X XXXX XXXX (11 digits starting with +61)

	// Remove any remaining non-digit characters
	digitsOnly := regexp.MustCompile(`[^0-9]`).ReplaceAllString(phone, "")

	// Check if it's a valid Australian number
	if len(digitsOnly) == 10 {
		// 10-digit Australian number
		if strings.HasPrefix(digitsOnly, "04") {
			// Mobile number
			return true, ""
		} else if strings.HasPrefix(digitsOnly, "02") ||
			strings.HasPrefix(digitsOnly, "03") ||
			strings.HasPrefix(digitsOnly, "07") ||
			strings.HasPrefix(digitsOnly, "08") {
			// Landline number
			return true, ""
		}
	} else if len(digitsOnly) == 11 && strings.HasPrefix(digitsOnly, "61") {
		// International format (+61)
		// Remove the 61 prefix and check the remaining 9 digits
		remaining := digitsOnly[2:]
		if strings.HasPrefix(remaining, "4") {
			// Mobile number
			return true, ""
		} else if strings.HasPrefix(remaining, "2") ||
			strings.HasPrefix(remaining, "3") ||
			strings.HasPrefix(remaining, "7") ||
			strings.HasPrefix(remaining, "8") {
			// Landline number
			return true, ""
		}
	}

	return false, "Please enter a valid Australian phone number (e.g., 0451889326 or +61451889326)."
}

// ValidatePassword validates password requirements
func (fv *FormValidator) ValidatePassword(password string) (bool, string) {
	if password == "" {
		return false, "Password is required."
	}

	if len(password) < 6 {
		return false, "Password must be at least 6 characters long."
	}

	return true, ""
}

// ValidatePasswordMatch validates that passwords match
func (fv *FormValidator) ValidatePasswordMatch(password, confirmPassword string) (bool, string) {
	if password != confirmPassword {
		return false, "Passwords do not match."
	}
	return true, ""
}

// ValidateRequired validates that a field is not empty
func (fv *FormValidator) ValidateRequired(value, fieldName string) (bool, string) {
	if strings.TrimSpace(value) == "" {
		return false, fieldName + " is required."
	}
	return true, ""
}

// FormError represents a form validation error
type FormError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormResponse represents a standardized form response
type FormResponse struct {
	Success bool        `json:"success"`
	Errors  []FormError `json:"errors,omitempty"`
	Message string      `json:"message,omitempty"`
}

// NewFormResponse creates a new form response
func NewFormResponse() *FormResponse {
	return &FormResponse{
		Success: true,
		Errors:  []FormError{},
	}
}

// AddError adds an error to the form response
func (fr *FormResponse) AddError(field, message string) {
	fr.Success = false
	fr.Errors = append(fr.Errors, FormError{
		Field:   field,
		Message: message,
	})
}

// HasErrors checks if the form response has any errors
func (fr *FormResponse) HasErrors() bool {
	return len(fr.Errors) > 0
}

// GetFirstError returns the first error message
func (fr *FormResponse) GetFirstError() string {
	if len(fr.Errors) > 0 {
		return fr.Errors[0].Message
	}
	return ""
}

// IsHTMXRequest checks if the request is an HTMX request
func IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// GetFormValue safely gets a form value with trimming
func GetFormValue(r *http.Request, key string) string {
	return strings.TrimSpace(r.FormValue(key))
}
