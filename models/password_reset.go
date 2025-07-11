package models

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type PasswordReset struct {
	ID        string    `bun:",pk,type:uuid,default:gen_random_uuid()"`
	UserID    string    `bun:",notnull"`
	Token     string    `bun:",unique,notnull"`
	Email     *string   `bun:",nullzero"`
	Phone     *string   `bun:",nullzero"`
	Method    string    `bun:",notnull"` // "email" or "sms"
	ExpiresAt time.Time `bun:",notnull"`
	Used      bool      `bun:",notnull,default:false"`
	CreatedAt time.Time `bun:",null,default:current_timestamp"`
}

var (
	ErrPasswordResetNotFound = errors.New("password reset token not found")
	ErrPasswordResetExpired  = errors.New("password reset token expired")
	ErrPasswordResetUsed     = errors.New("password reset token already used")
)

// GenerateSecureToken creates a cryptographically secure random token
func GenerateSecureToken() (string, error) {
	b := make([]byte, 32) // 256 bits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// CreatePasswordReset creates a new password reset token for a user
func CreatePasswordReset(ctx context.Context, db bun.IDB, userID, email, phone, method string) (*PasswordReset, error) {
	token, err := GenerateSecureToken()
	if err != nil {
		return nil, err
	}

	// Token expires in 1 hour
	expiresAt := time.Now().Add(1 * time.Hour)

	reset := &PasswordReset{
		UserID:    userID,
		Token:     token,
		Email:     &email,
		Phone:     &phone,
		Method:    method,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	// Set the appropriate field based on method
	if method == "email" {
		reset.Phone = nil
	} else if method == "sms" {
		reset.Email = nil
	}

	_, err = db.NewInsert().Model(reset).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return reset, nil
}

// GetPasswordResetByToken retrieves a password reset by token
func GetPasswordResetByToken(ctx context.Context, db bun.IDB, token string) (*PasswordReset, error) {
	reset := new(PasswordReset)
	err := db.NewSelect().Model(reset).Where("token = ?", token).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPasswordResetNotFound
		}
		return nil, err
	}

	// Check if token is expired
	if time.Now().After(reset.ExpiresAt) {
		return nil, ErrPasswordResetExpired
	}

	// Check if token has been used
	if reset.Used {
		return nil, ErrPasswordResetUsed
	}

	return reset, nil
}

// MarkPasswordResetAsUsed marks a password reset token as used
func MarkPasswordResetAsUsed(ctx context.Context, db bun.IDB, token string) error {
	_, err := db.NewUpdate().Model(&PasswordReset{}).
		Set("used = ?", true).
		Where("token = ?", token).
		Exec(ctx)
	return err
}

// CleanupExpiredPasswordResets removes expired password reset tokens
func CleanupExpiredPasswordResets(ctx context.Context, db bun.IDB) error {
	_, err := db.NewDelete().Model(&PasswordReset{}).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)
	return err
}
