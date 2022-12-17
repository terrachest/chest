package modules_test

import (
	"github.com/google/go-cmp/cmp"
	"privateterraformregistry/internal/module"
	"privateterraformregistry/internal/modules"
	"testing"
)

func TestGetModules(t *testing.T) {
	original := modules.Modules{
		Modules: []module.Module{
			{"foo", "bar", "baz", "1.0.0"},
		},
	}

	got := original.GetModules()
	if !cmp.Equal(original.Modules, got) {
		t.Errorf("Modules from GetModules() different to original")
	}

	got = append(got, module.Module{Namespace: "foo", Name: "bar", System: "baz", Version: "1.0.0"})
	if cmp.Equal(original.Modules, got) {
		t.Errorf("Modifying result from GetModules should not modify original")
	}
}

func TestAdd(t *testing.T) {
	ms := modules.Modules{}
	ms.Add(module.Module{
		Namespace: "foo",
		System:    "bar",
		Name:      "baz",
		Version:   "1.0.0",
	})

	if len(ms.GetModules()) != 1 {
		t.Errorf("Wanted length of modules to be 1 got %d", len(ms.GetModules()))
	}

	ms.Add(module.Module{
		Namespace: "foo",
		System:    "bar",
		Name:      "baz",
		Version:   "1.0.0",
	})

	if len(ms.GetModules()) != 1 {
		t.Errorf("Wanted length of modules to be 1 got %d", len(ms.GetModules()))
	}

	ms.Add(module.Module{
		Namespace: "foo",
		System:    "bar",
		Name:      "baz",
		Version:   "1.1.0",
	})

	if len(ms.GetModules()) != 2 {
		t.Errorf("Wanted length of modules to be 2 got %d", len(ms.GetModules()))
	}
}

func TestExists(t *testing.T) {
	ms := modules.Modules{}
	ms.Add(module.Module{
		Namespace: "foo",
		System:    "bar",
		Name:      "baz",
		Version:   "1.0.0",
	})

	moduleExists := ms.Exists(module.Module{
		Namespace: "foo",
		System:    "bar",
		Name:      "baz",
		Version:   "1.0.0",
	})

	if !moduleExists {
		t.Error("expected moduleExists to be true got false")
	}

	moduleExists = ms.Exists(module.Module{
		Namespace: "foo",
		System:    "bar",
		Name:      "bat",
		Version:   "1.0.0",
	})

	if moduleExists {
		t.Error("expected moduleExists to be false got true")
	}
}
