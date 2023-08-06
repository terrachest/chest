package env_test

import (
	"testing"

	"github.com/terrachest/server/internal/utils/env"
)

func TestGet(t *testing.T) {
	t.Setenv("foo", "bar")
	got := env.Get("foo")
	if got != "bar" {
		t.Errorf("got = %s; want bar", got)
	}

	t.Setenv("foo", "")
	got = env.Get("foo", "bar")
	if got != "bar" {
		t.Errorf("got = %s; want bar", got)
	}

	got = env.Get("foo", "")
	if got != "" {
		t.Errorf("got = %s; wanted ''", got)
	}
}
