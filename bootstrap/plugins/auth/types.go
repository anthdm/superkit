package auth

import (
	"auth/app/db"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	UserID   int
	Email    string
	LoggedIn bool
}

func (auth Auth) Check() bool {
	return auth.LoggedIn
}

type User struct {
	ID              int `bun:"id,pk,autoincrement"`
	Email           string
	FirstName       string
	LastName        string
	PasswordHash    string
	EmailVerifiedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func createUserFromFormValues(values SignupFormValues) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(values.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user := User{
		Email:        values.Email,
		FirstName:    values.FirstName,
		LastName:     values.LastName,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	_, err = db.Query.NewInsert().Model(&user).Exec(context.Background())
	return user, err
}

type Session struct {
	ID          int `bun:"id,pk,autoincrement"`
	UserID      int
	Token       string
	IPAddress   string
	UserAgent   string
	ExpiresAt   time.Time
	LastLoginAt time.Time
	CreatedAt   time.Time

	User User `bun:"rel:belongs-to,join:user_id=id"`
}
