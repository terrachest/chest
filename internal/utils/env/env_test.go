package env_test

import (
	"privateterraformregistry/internal/utils/env"
	"testing"
)

func TestGet(t *testing.T) {
	t.Setenv("foo", "bar")
	got := env.Get("foo")
	if got != "bar" {
		t.Errorf("Expected bar got %s", got)
	}

	t.Setenv("foo", "")
	got = env.Get("foo", "bar")
	if got != "bar" {
		t.Errorf("Excpected foo got %s", got)
	}

	got = env.Get("foo", "")
	if got != "" {
		t.Errorf("Excpected empty string got %s", got)
	}
}
