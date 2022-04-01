package db_models

import "time"

// Represents a user.
type User struct {
	Id                string    `db:"id"`
	CreatedOn         time.Time `db:"created_on"`
	UpdatedOn         time.Time `db:"updated_on"`
	Email             string    `db:"email"`
	Name              string    `db:"name"`
	Bio               string    `db:"bio"`
	ProfilePhotoId    string    `db:"profile_photo_id"`
	PasswordHash      string    `db:"password_hash"`
	SignInAttempts    int       `db:"sign_in_attempts"`
	IsActive          bool      `db:"is_active"`
	AccountResetToken string    `db:"account_reset_token"`
}

func (t User) GetId() string { return t.Id }
