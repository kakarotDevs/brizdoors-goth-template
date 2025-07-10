package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/kakarotDevs/brizdoors-goth-template/utils"
	"github.com/uptrace/bun"
)

type User struct {
	ID        string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	GoogleID  string `bun:"google_id,unique,nullzero"`
	FirstName string `bun:",nullzero"`
	LastName  string `bun:",nullzero"`
	Role      string `bun:",notnull,default:'user'"` // e.g. "user", "admin"
	Email     string `bun:",unique,notnull"`
	Picture   string
	Password  *string   `bun:",nullzero"`
	CreatedAt time.Time `bun:",null,default:current_timestamp"`
	UpdatedAt time.Time `bun:",null,default:current_timestamp"`
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

func GetUserByID(ctx context.Context, db bun.IDB, id string) (*User, error) {
	user := new(User)
	err := db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func GetUserByEmail(ctx context.Context, db bun.IDB, email string) (*User, error) {
	user := new(User)
	err := db.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func GetUserByGoogleID(ctx context.Context, db bun.IDB, googleID string) (*User, error) {
	user := new(User)
	err := db.NewSelect().Model(user).Where("google_id = ?", googleID).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func CreateUser(ctx context.Context, db bun.IDB, user *User) error {
	_, err := db.NewInsert().Model(user).Exec(ctx)
	return err
}

func CreateOrUpdateUser(ctx context.Context, db bun.IDB, user *User) error {
	user.UpdatedAt = time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = user.UpdatedAt
	}
	_, err := db.NewInsert().Model(user).
		On("CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email, picture = EXCLUDED.picture, updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	return err
}

func UpdateUser(ctx context.Context, db bun.IDB, user *User) error {
	user.UpdatedAt = time.Now()
	_, err := db.NewUpdate().Model(user).
		Set("first_name = ?", user.FirstName).
		Set("last_name = ?", user.LastName).
		Set("email = ?", user.Email).
		Set("updated_at = ?", user.UpdatedAt).
		Where("id = ?", user.ID).
		Exec(ctx)
	return err
}

func DeleteUser(ctx context.Context, db bun.IDB, id string) error {
	_, err := db.NewDelete().Model((*User)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// SetUserPassword sets a hashed password for a user
func SetUserPassword(ctx context.Context, db bun.IDB, userID string, newPw string) error {
	hash, err := utils.HashPassword(newPw)
	if err != nil {
		return err
	}
	_, err = db.NewUpdate().Model(&User{}).
		Set("password = ?", hash).
		Set("updated_at = NOW()").
		Where("id = ?", userID).
		Exec(ctx)
	return err
}

// For OAuth users (or login conversion):
func ClearUserPassword(ctx context.Context, db bun.IDB, userID string) error {
	_, err := db.NewUpdate().Model(&User{}).
		Set("password = NULL").
		Set("updated_at = NOW()").
		Where("id = ?", userID).
		Exec(ctx)
	return err
}
