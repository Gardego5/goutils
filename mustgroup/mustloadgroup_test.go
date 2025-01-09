package mustgroup_test

import (
	"log/slog"
	"testing"

	"github.com/Gardego5/goutils/env"
	. "github.com/Gardego5/goutils/mustgroup"
)

func TestMustLoadGroup(t *testing.T) {
	type Variables struct {
		LogLevel slog.LevelVar `env:"LOG_LEVEL"`
	}

	t.Run("Successfully loads the environment variable", func(t *testing.T) {
		expectPanic(t, false)

		t.Setenv("LOG_LEVEL", "WARN")

		g := &Group{}
		env := Must(env.Load[Variables]())(g)
		a := Must(somethingWithoutError(true))(g)
		g.Must()

		if env.LogLevel.Level() != slog.LevelWarn.Level() {
			t.Fatalf("expected %s, got %s", slog.LevelWarn, env.LogLevel.Level())
		}
		if a != true {
			t.Fatalf("expected true, got %v", a)
		}
	})

	t.Run("Panics given missing environment variables", func(t *testing.T) {
		expectPanic(t, true)

		g := &Group{}
		_ = Must(env.Load[Variables]())(g)
		_ = Must(somethingWithoutError(true))(g)
		g.Must()

		t.Fatalf("expected a panic")
	})
}
