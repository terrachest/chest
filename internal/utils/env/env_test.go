package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetReturnsSetEnvironmentVariable(t *testing.T) {
	t.Setenv("foo", "bar")
	sut := Get("foo")

	assert.Equal(t, sut, "bar")
}

func TestGetDefaultsToDefaultValue(t *testing.T) {
	sut := Get("foo", "bar")

	assert.Equal(t, sut, "bar")
}

func TestGetReturnsEmptyStringWhenNoDefault(t *testing.T) {
	sut := Get("foo")

	assert.Equal(t, sut, "")
}
