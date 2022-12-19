package env_test

import (
	"privateterraformregistry/internal/utils/env"
	"testing"
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
