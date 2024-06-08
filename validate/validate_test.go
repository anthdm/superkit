package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		"email": Rules(Email()),
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
