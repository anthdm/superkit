package validate

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var createdAt = time.Now()

var testSchema = Schema{
	"createdAt": Rules(Time),
	"startedAt": Rules(TimeBefore(time.Now())),
	"deletedAt": Rules(TimeAfter(createdAt)),
	"email":     Rules(Email),
	"url":       Rules(URL),
	"password": Rules(
		ContainsSpecial,
		ContainsUpper,
		ContainsDigit,
		Min(7),
		Max(50),
	),
	"age":      Rules(GTE(18)),
	"bet":      Rules(GT(0), LTE(10)),
	"username": Rules(Required),
}

func TestValidateRequest(t *testing.T) {
	var (
		email        = "foo@bar.com"
		password     = "superHunter123@"
		firstName    = "Anthony"
		website      = "http://foo.com"
		randomNumber = 123
		randomFloat  = 9.999
	)
	formValues := url.Values{}
	formValues.Set("email", email)
	formValues.Set("password", password)
	formValues.Set("firstName", firstName)
	formValues.Set("url", website)
	formValues.Set("brandom", fmt.Sprint(randomNumber))
	formValues.Set("arandom", fmt.Sprint(randomFloat))
	encodedValues := formValues.Encode()

	req, err := http.NewRequest("POST", "http://foo.com", strings.NewReader(encodedValues))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	type SignupData struct {
		Email                string  `form:"email"`
		Password             string  `form:"password"`
		FirstName            string  `form:"firstName"`
		URL                  string  `form:"url"`
		ARandomRenamedNumber int     `form:"brandom"`
		ARandomRenamedFloat  float64 `form:"arandom"`
	}

	schema := Schema{
		"Email": Rules(Email),
		"Password": Rules(
			Required,
			ContainsDigit,
			ContainsUpper,
			ContainsSpecial,
			Min(7),
		),
		"FirstName":            Rules(Min(3), Max(50)),
		"URL":                  Rules(URL),
		"ARandomRenamedNumber": Rules(GT(100), LT(124)),
		"ARandomRenamedFloat":  Rules(GT(9.0), LT(10.1)),
	}

	var data SignupData
	errors, ok := Request(req, &data, schema)
	assert.True(t, ok)
	assert.Empty(t, errors)

	assert.Equal(t, data.Email, email)
	assert.Equal(t, data.Password, password)
	assert.Equal(t, data.FirstName, firstName)
	assert.Equal(t, data.URL, website)
	assert.Equal(t, data.ARandomRenamedNumber, randomNumber)
	assert.Equal(t, data.ARandomRenamedFloat, randomFloat)
}

func TestTime(t *testing.T) {
	type Foo struct {
		CreatedAt time.Time
	}
	foo := Foo{
		CreatedAt: time.Now(),
	}
	schema := Schema{
		"createdAt": Rules(Time),
	}
	_, ok := Validate(foo, schema)
	assert.True(t, ok)

	foo.CreatedAt = time.Time{}
	_, ok = Validate(foo, schema)
	assert.False(t, ok)
}

func TestURL(t *testing.T) {
	type Foo struct {
		URL string `v:"URL"`
	}
	foo := Foo{
		URL: "not an url",
	}
	schema := Schema{
		"URL": Rules(URL),
	}
	errors, ok := Validate(foo, schema)
	assert.False(t, ok)
	assert.NotEmpty(t, errors)

	validURLS := []string{
		"http://google.com",
		"http://www.google.com",
		"https://www.google.com",
		"https://www.google.com",
		"www.google.com",
		"https://book.com/sales",
		"app.book.com",
		"app.book.com/signup",
	}

	for _, url := range validURLS {
		foo.URL = url
		errors, ok = Validate(foo, schema)
		assert.True(t, ok)
		assert.Empty(t, errors)
	}
}

func TestContainsUpper(t *testing.T) {
	type Foo struct {
		Password string
	}
	foo := Foo{"hunter"}
	schema := Schema{
		"Password": Rules(ContainsUpper),
	}
	errors, ok := Validate(foo, schema)
	assert.False(t, ok)
	assert.NotEmpty(t, errors)

	foo.Password = "Hunter"
	errors, ok = Validate(foo, schema)
	assert.True(t, ok)
	assert.Empty(t, errors)
}

func TestContainsDigit(t *testing.T) {
	type Foo struct {
		Password string
	}
	foo := Foo{"hunter"}
	schema := Schema{
		"Password": Rules(ContainsDigit),
	}
	errors, ok := Validate(foo, schema)
	assert.False(t, ok)
	assert.NotEmpty(t, errors)

	foo.Password = "Hunter1"
	errors, ok = Validate(foo, schema)
	assert.True(t, ok)
	assert.Empty(t, errors)
}

func TestContainsSpecial(t *testing.T) {
	type Foo struct {
		Password string
	}
	foo := Foo{"hunter"}
	schema := Schema{
		"Password": Rules(ContainsSpecial),
	}
	errors, ok := Validate(foo, schema)
	assert.False(t, ok)
	assert.NotEmpty(t, errors)

	foo.Password = "Hunter@"
	errors, ok = Validate(foo, schema)
	assert.True(t, ok)
	assert.Empty(t, errors)
}

func TestRuleIn(t *testing.T) {
	type Foo struct {
		Currency string
	}
	foo := Foo{"eur"}
	schema := Schema{
		"currency": Rules(In([]string{"eur", "usd", "chz"})),
	}
	errors, ok := Validate(foo, schema)
	assert.True(t, ok)
	assert.Empty(t, errors)
	foo = Foo{"foo"}
	errors, ok = Validate(foo, schema)
	assert.False(t, ok)
	assert.Len(t, errors["currency"], 1)
}

func TestValidate(t *testing.T) {
	type User struct {
		Email    string
		Username string
	}
	schema := Schema{
		"email": Rules(Email),
		// Test both lower and uppercase
		"Username": Rules(Min(3), Max(10)),
	}
	user := User{
		Email:    "foo@bar.com",
		Username: "pedropedro",
	}
	errors, ok := Validate(user, schema)
	assert.True(t, ok)
	assert.Empty(t, errors)
}

func TestMergeSchemas(t *testing.T) {
	expected := Schema{
		"Name":      Rules(),
		"Email":     Rules(),
		"FirstName": Rules(),
		"LastName":  Rules(),
	}
	a := Schema{
		"Name":  Rules(),
		"Email": Rules(),
	}
	b := Schema{
		"FirstName": Rules(),
		"LastName":  Rules(),
	}
	c := Merge(a, b)
	assert.Equal(t, expected, c)
}
