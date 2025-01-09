package mustgroup_test

import "testing"

func somethingWithError[T any](t T, err error) (T, error) { return t, err }
func somethingWithoutError[T any](t T) (T, error)         { return t, nil }

func expectPanic(t *testing.T, shouldPanic bool) {
	err := recover()

	if shouldPanic {
		if err == nil {
			t.Fatalf("expected a panic, got nil")
		}
	} else {
		if err != nil {
			t.Fatalf("expected no panic, got %v", err)
		}
	}
}
