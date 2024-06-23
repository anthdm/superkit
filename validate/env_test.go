package validate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyEnv(t *testing.T) {
	assert.Panics(t, func() {
		_ = Env[string]("TEST", Rules(Required))
	})

	assert.Panics(t, func() {
		os.Setenv("TEST", "")
		_ = Env[string]("TEST", Rules(Required))
	})
}

func TestReturnsValue(t *testing.T) {

	os.Setenv("TEST", "value")

	val := Env[string]("TEST", Rules(Required))
	assert.Equal(t, val, "value")
}

func TestDefault(t *testing.T) {
	val := Env[string]("TEST2", Rules(Required), "hello")
	assert.Equal(t, "hello", val)

	os.Setenv("TEST2", "world")
	val = Env[string]("TEST2", Rules(Required), "hello")
	assert.Equal(t, "world", val)
}

func TestInt(t *testing.T) {
	os.Setenv("TEST", "1")
	val := Env[int]("TEST", Rules(LT(2)))
	assert.Equal(t, 1, val)
}

func TestBool(t *testing.T) {
	os.Setenv("TEST", "true")
	val := Env[bool]("TEST", Rules(EQ(true)))
	assert.Equal(t, true, val)
}

func TestFloat(t *testing.T) {
	os.Setenv("TEST", "1.1")
	val := Env[float64]("TEST", Rules(GT(1.0)))
	assert.Equal(t, 1.1, val)
}
