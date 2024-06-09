package validate

import (
	"fmt"
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
		URL string
	}
	foo := Foo{
		URL: "not an url",
	}
	schema := Schema{
		"url": Rules(URL),
	}
	errors, ok := Validate(foo, schema)
	assert.False(t, ok)

	foo.URL = "www.user.com"
	errors, ok = Validate(foo, schema)
	assert.True(t, ok)
	fmt.Println(errors)
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
