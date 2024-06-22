package validate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyEnv(t *testing.T) {
	assert.Panics(t, func() {
		_ = Env("TEST", Required).Validate()
	})

	assert.Panics(t, func() {
		os.Setenv("TEST", "")

		_ = Env("TEST", Required).Validate()
	})
}

func TestReturnsValue(t *testing.T) {

	os.Setenv("TEST", "value")

	val := Env("TEST", Required).Validate()
	assert.Equal(t, val, "value")
}

func TestDefault(t *testing.T) {
	val := Env("TEST2").Default("hello").Validate()
	assert.Equal(t, "hello", val)

	os.Setenv("TEST2", "world")
	val = Env("TEST2").Default("hello").Validate()
	assert.Equal(t, "world", val)
}
