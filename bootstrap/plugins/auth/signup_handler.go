package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/anthdm/superkit/event"
	"github.com/anthdm/superkit/kit"
	v "github.com/anthdm/superkit/validate"
	"github.com/golang-jwt/jwt/v5"
)

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
	event.Emit("auth.signup", UserWithVerificationToken{
		Token: token,
		User:  user,
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
