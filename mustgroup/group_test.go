package mustgroup_test

import (
	"errors"
	"testing"

	. "github.com/Gardego5/goutils/mustgroup"
)

func TestGroup(t *testing.T) {
	t.Run("Group.Error() returns nil given no input", func(t *testing.T) {
		g := &Group{}
		err := g.Error()
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("Group.Error() returns nil given no errors", func(t *testing.T) {
		g := &Group{}
		aValue := Must(somethingWithoutError(true))(g)
		bValue := Must(somethingWithoutError("Hello, World!"))(g)
		err := g.Error()

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
		if aValue != true {
			t.Fatalf("expected true, got %v", aValue)
		}
		if bValue != "Hello, World!" {
			t.Fatalf("expected 'Hello, World!', got %v", bValue)
		}
	})

	t.Run("Group.Error() returns an error given an error", func(t *testing.T) {
		g := &Group{}
		_ = Must(somethingWithError(30, errors.New("an error occurred")))(g)
		err := g.Error()

		if err == nil {
			t.Fatalf("expected an error, got nil")
		}
	})

	t.Run("Group.Must() panics given an error", func(t *testing.T) {
		defer expectPanic(t, true)

		expectedErr := errors.New("an error occurred")

		g := &Group{}
		_ = Must(somethingWithoutError(30))(g)
		_ = Must(somethingWithError("hi", expectedErr))(g)
		g.Must()

		t.Fatalf("expected a panic")
	})

	t.Run("Group.Must() does not panic given no errors", func(t *testing.T) {
		defer expectPanic(t, false)

		g := &Group{}
		a := Must(somethingWithoutError(30))(g)
		b := Must(somethingWithoutError("hi"))(g)
		g.Must()

		if a != 30 {
			t.Fatalf("expected 30, got %v", a)
		}
		if b != "hi" {
			t.Fatalf("expected 'hi', got %v", b)
		}
	})
}
