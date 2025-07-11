package auth

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/kakarotDevs/brizdoors-goth-template/db"
	"github.com/kakarotDevs/brizdoors-goth-template/models"
)

// PasswordResetService handles password reset business logic
type PasswordResetService struct{}

// NewPasswordResetService creates a new password reset service
func NewPasswordResetService() *PasswordResetService {
	return &PasswordResetService{}
}

// SendPasswordResetEmail sends a password reset email
func (s *PasswordResetService) SendPasswordResetEmail(email, token string) error {
	// For now, we'll just log the reset link
	// In production, you'd integrate with an email service like SendGrid, AWS SES, etc.

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, url.QueryEscape(token))

	fmt.Printf("Password reset email for %s:\n", email)
	fmt.Printf("Reset URL: %s\n", resetURL)
	fmt.Printf("Token: %s\n", token)

	return nil
}

// SendPasswordResetSMS sends a password reset SMS
func (s *PasswordResetService) SendPasswordResetSMS(phone, token string) error {
	// TODO: Implement actual SMS sending using a service like Twilio
	// For now, just log the verification code
	fmt.Printf("SMS verification code for %s: %s\n", phone, token)
	return nil
}

// CreatePasswordReset creates a password reset for the given contact (email or phone)
func (s *PasswordResetService) CreatePasswordReset(ctx context.Context, contact string) error {
	// Determine if contact is email or phone
	isEmail := contains(contact, "@")

	if isEmail {
		// Handle email reset
		user, err := models.GetUserByEmail(ctx, db.DB, contact)
		if err != nil {
			// Don't reveal if user exists or not for security
			return nil
		}

		// Skip if user only has Google OAuth (no password)
		if user.Password == nil {
			return nil
		}

		// Create password reset token
		reset, err := models.CreatePasswordReset(ctx, db.DB, user.ID, contact, "", "email")
		if err != nil {
			return fmt.Errorf("error creating reset link: %w", err)
		}

		// Send reset email
		if err := s.SendPasswordResetEmail(user.Email, reset.Token); err != nil {
			return fmt.Errorf("error sending reset link: %w", err)
		}
	} else {
		// Handle phone reset
		user, err := models.GetUserByPhone(ctx, db.DB, contact)
		if err != nil {
			// Don't reveal if user exists or not for security
			return nil
		}

		// Skip if user only has Google OAuth (no password)
		if user.Password == nil {
			return nil
		}

		// Create password reset token
		reset, err := models.CreatePasswordReset(ctx, db.DB, user.ID, "", contact, "sms")
		if err != nil {
			return fmt.Errorf("error creating verification code: %w", err)
		}

		// Send SMS verification code
		if err := s.SendPasswordResetSMS(contact, reset.Token); err != nil {
			return fmt.Errorf("error sending verification code: %w", err)
		}
	}

	return nil
}

// ValidateResetToken validates a password reset token
func (s *PasswordResetService) ValidateResetToken(ctx context.Context, token string) error {
	_, err := models.GetPasswordResetByToken(ctx, db.DB, token)
	return err
}

// ResetPassword resets a user's password using a valid token
func (s *PasswordResetService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Get and validate reset token
	reset, err := models.GetPasswordResetByToken(ctx, db.DB, token)
	if err != nil {
		return fmt.Errorf("invalid or expired reset link: %w", err)
	}

	// Update user password
	if err := models.SetUserPassword(ctx, db.DB, reset.UserID, newPassword); err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	// Mark token as used
	if err := models.MarkPasswordResetAsUsed(ctx, db.DB, token); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to mark reset token as used: %v\n", err)
	}

	return nil
}

// Helper function to check if string contains @ (email)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 1; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}()))
}
