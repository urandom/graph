package graph

import (
	"encoding/json"
	"testing"
)

func TestRegisterNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should've been a valid panic")
		}
	}()

	RegisterLinker("test", nil)

	t.Fatalf("Can't register nil constructors")
}

func TestRegisterDup(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should've been a valid panic")
		}
	}()

	RegisterLinker("test", func(opts json.RawMessage) (Linker, error) { return nil, nil })
	RegisterLinker("test", func(opts json.RawMessage) (Linker, error) { return nil, nil })
}

func TestRegister(t *testing.T) {
	if _, ok := operations["test2"]; ok {
		t.Fatalf("test2 shouldn't have been registered yet")
	}

	RegisterLinker("test2", func(opts json.RawMessage) (Linker, error) {
		return nil, nil
	})

	if _, ok := operations["test2"]; !ok {
		t.Fatalf("test2 should be registered")
	}
}
