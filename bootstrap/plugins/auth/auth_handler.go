package auth

import (
	"auth/app/db"
	"cmp"
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/anthdm/gothkit/kit"
	v "github.com/anthdm/gothkit/validate"
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
		redirectURL := cmp.Or(os.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN"), "/profile")
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

	// todo: use the kit.Getenv instead of the comp thingy
	skipVerify := cmp.Or(os.Getenv("SUPERKIT_AUTH_SKIP_VERIFY"), "false")
	if skipVerify != "true" {
		if user.EmailVerifiedAt.Equal(time.Time{}) {
			errors.Add("verified", "please verify your email")
			return kit.Render(LoginForm(values, errors))
		}
	}

	sessionExpiryStr := os.Getenv("SUPERKIT_AUTH_SESSION_EXPIRY_IN_HOURS")
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

	redirectURL := cmp.Or(os.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN"), "/profile")

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
	return kit.Render(ConfirmEmail(user.Email))
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
	// TODO: do we really need to check if the user is verified
	// even if we check that already in the login process.
	// if session.User.EmailVerifiedAt.Equal(time.Time{}) {
	// 	return Auth{}, nil
	// }
	return Auth{
		LoggedIn: true,
		UserID:   session.User.ID,
		Email:    session.User.Email,
	}, nil
}
