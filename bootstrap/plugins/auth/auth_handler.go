package auth

import (
	"AABBCCDD/app/db"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/anthdm/superkit/event"
	"github.com/anthdm/superkit/kit"
	"github.com/anthdm/superkit/mail"
	v "github.com/anthdm/superkit/validate"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	userSessionName = "user-session"
)

var authSchema = v.Schema{
	"email":    v.Rules(v.Email),
	"password": v.Rules(v.Required),
}

var signupSchema = v.Schema{
	"email": v.Rules(v.Email),
	"password": v.Rules(
		v.ContainsSpecial,
		v.ContainsUpper,
		v.Min(7),
		v.Max(50),
	),
	"firstName": v.Rules(v.Min(2), v.Max(50)),
	"lastName":  v.Rules(v.Min(2), v.Max(50)),
}

func HandleAuthIndex(kit *kit.Kit) error {
	if kit.Auth().Check() {
		redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/profile")
		return kit.Redirect(http.StatusSeeOther, redirectURL)
	}
	return kit.Render(AuthIndex(AuthIndexPageData{}))
}

func HandleAuthCreate(kit *kit.Kit) error {
	var values LoginFormValues
	errors, ok := v.Request(kit.Request, &values, authSchema)
	if !ok {
		return kit.Render(LoginForm(values, errors))
	}

	var user User
	err := db.Query.NewSelect().
		Model(&user).
		Where("user.email = ?", values.Email).
		Scan(kit.Request.Context())
	if err != nil {
		if err == sql.ErrNoRows {
			errors.Add("credentials", "invalid credentials")
			return kit.Render(LoginForm(values, errors))
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(values.Password))
	if err != nil {
		errors.Add("credentials", "invalid credentials")
		return kit.Render(LoginForm(values, errors))
	}

	skipVerify := kit.Getenv("SUPERKIT_AUTH_SKIP_VERIFY", "false")
	if skipVerify != "true" {
		if user.EmailVerifiedAt.Equal(time.Time{}) {
			errors.Add("verified", "please verify your email")
			return kit.Render(LoginForm(values, errors))
		}
	}

	sessionExpiryStr := kit.Getenv("SUPERKIT_AUTH_SESSION_EXPIRY_IN_HOURS", "48")
	sessionExpiry, err := strconv.Atoi(sessionExpiryStr)
	if err != nil {
		sessionExpiry = 48
	}
	session := Session{
		UserID:      user.ID,
		Token:       uuid.New().String(),
		CreatedAt:   time.Now(),
		LastLoginAt: time.Now(),
		ExpiresAt:   time.Now().Add(time.Hour * time.Duration(sessionExpiry)),
	}
	_, err = db.Query.NewInsert().
		Model(&session).
		Exec(kit.Request.Context())
	if err != nil {
		return err
	}

	// TODO change this with kit.Getenv
	sess := kit.GetSession(userSessionName)
	sess.Values["sessionToken"] = session.Token
	sess.Save(kit.Request, kit.Response)

	redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/profile")

	return kit.Redirect(http.StatusSeeOther, redirectURL)
}

func HandleAuthDelete(kit *kit.Kit) error {
	sess := kit.GetSession(userSessionName)
	defer func() {
		sess.Values = map[any]any{}
		sess.Save(kit.Request, kit.Response)
	}()
	_, err := db.Query.NewDelete().
		Model((*Session)(nil)).
		Where("token = ?", sess.Values["sessionToken"]).
		Exec(kit.Request.Context())
	if err != nil {
		return err
	}
	return kit.Redirect(http.StatusSeeOther, "/")
}

func HandleSignupIndex(kit *kit.Kit) error {
	return kit.Render(SignupIndex(SignupIndexPageData{}))
}

func HandleSignupCreate(kit *kit.Kit) error {
	var values SignupFormValues
	errors, ok := v.Request(kit.Request, &values, signupSchema)
	if !ok {
		return kit.Render(SignupForm(values, errors))
	}
	if values.Password != values.PasswordConfirm {
		errors.Add("passwordConfirm", "passwords do not match")
		return kit.Render(SignupForm(values, errors))
	}
	user, err := createUserFromFormValues(values)
	if err != nil {
		return err
	}
	token, err := createVerificationToken(user.ID)
	if err != nil {
		return err
	}
	event.Emit("user.signup", user)
	mail.NOOPMailer{}.SendEmail(kit.Request.Context(), mail.Contents{
		Title: token,
	})
	return kit.Render(ConfirmEmail(user.Email))
}

func createVerificationToken(userID int) (string, error) {
	expiryStr := kit.Getenv("SUPERKIT_AUTH_EMAIL_VERIFICATION_EXPIRY_IN_HOURS", "1")
	expiry, err := strconv.Atoi(expiryStr)
	if err != nil {
		expiry = 1
	}

	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprint(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiry))),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("SUPERKIT_SECRET")))
}

func HandleEmailVerify(kit *kit.Kit) error {
	tokenStr := kit.Request.URL.Query().Get("token")
	if len(tokenStr) == 0 {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}

	token, err := jwt.ParseWithClaims(
		tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("SUPERKIT_SECRET")), nil
		}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return err
	}
	if !token.Valid {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return kit.Render(EmailVerificationError("Email verification token expired"))
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return kit.Render(EmailVerificationError("Email verification token expired"))
	}

	var user User
	err = db.Query.NewSelect().
		Model(&user).
		Where("id = ?", userID).
		Scan(kit.Request.Context())
	if err != nil {
		return err
	}

	if user.EmailVerifiedAt.After(time.Time{}) {
		return kit.Render(EmailVerificationError("Email already verified"))
	}

	user.EmailVerifiedAt = time.Now()
	_, err = db.Query.NewUpdate().
		Model(&user).
		WherePK().
		Exec(kit.Request.Context())
	if err != nil {
		return err
	}

	return kit.Redirect(http.StatusSeeOther, "/login")
}

func AuthenticateUser(kit *kit.Kit) (kit.Auth, error) {
	auth := Auth{}
	sess := kit.GetSession(userSessionName)
	token, ok := sess.Values["sessionToken"]
	if !ok {
		return auth, nil
	}

	var session Session
	err := db.Query.NewSelect().
		Model(&session).
		Relation("User").
		Where("session.token = ? AND session.expires_at > ?", token, time.Now()).
		Scan(kit.Request.Context())
	if err != nil {
		return auth, nil
	}
	return Auth{

		LoggedIn: true,
		UserID:   session.User.ID,
		Email:    session.User.Email,
	}, nil
}
