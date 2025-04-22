package env_test

import (
	"log/slog"
	"testing"

	. "github.com/Gardego5/goutils/env"
)

func TestLoad(t *testing.T) {
	type Variables struct {
		LogLevel slog.LevelVar `env:"LOG_LEVEL"`
	}

	t.Run("Returns an error given no environment variable set", func(t *testing.T) {
		_, err := Load[Variables]()
		if err == nil {
			t.Fatal("expected error for missing LOG_LEVEL")
		}
	})

	t.Run("Returns a struct with the environment variable when present", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "WARN")
		env, err := Load[Variables]()
		if err != nil {
			t.Fatal(err)
		}

		if env.LogLevel.Level() != slog.LevelWarn.Level() {
			t.Fatalf("expected %s, got %s", slog.LevelWarn, env.LogLevel.Level())
		}
	})

	type VariablesWithDefault struct {
		LogLevel slog.Level `env:"LOG_LEVEL=INFO"`
	}

	t.Run("Returns an error given ", func(t *testing.T) {
		env, err := Load[VariablesWithDefault]()
		if err != nil {
			t.Fatal(err)
		}

		if env.LogLevel.Level() != slog.LevelInfo.Level() {
			t.Fatalf("expected default value %s, got %s", slog.LevelInfo, env.LogLevel.Level())
		}
	})

	t.Run("Returns a struct with the environment variable when present", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "WARN")
		env, err := Load[VariablesWithDefault]()
		if err != nil {
			t.Fatal(err)
		}

		if env.LogLevel.Level() != slog.LevelWarn.Level() {
			t.Fatalf("expected %s, got %s", slog.LevelWarn, env.LogLevel.Level())
		}
	})

	t.Run("Numbers types are parsed without error", func(t *testing.T) {
		type Variables struct {
			Bool    bool    `env:"BOOL"`
			Int     int     `env:"INT"`
			Int8    int8    `env:"INT_8"`
			Int16   int16   `env:"INT_16"`
			Int32   int32   `env:"INT_32"`
			Int64   int64   `env:"INT_64"`
			Float32 float32 `env:"FLOAT_32"`
			Float64 float64 `env:"FLOAT_64"`
		}

		t.Setenv("BOOL", "true")
		t.Setenv("INT", "1")
		t.Setenv("INT_8", "127")
		t.Setenv("INT_16", "32767")
		t.Setenv("INT_32", "2147483647")
		t.Setenv("INT_64", "9223372036854775807")
		t.Setenv("FLOAT_32", "3.4028234663852886e+38")
		t.Setenv("FLOAT_64", "1.7976931348623157e+308")

		env, err := Load[Variables]()
		if err != nil {
			t.Fatal(err)
		}
		if !env.Bool {
			t.Fatalf("expected true, got %t", env.Bool)
		}
		if env.Int != 1 {
			t.Fatalf("expected 1, got %d", env.Int)
		}
		if env.Int8 != 127 {
			t.Fatalf("expected 0, got %d", env.Int8)
		}
		if env.Int16 != 32767 {
			t.Fatalf("expected 32767, got %d", env.Int16)
		}
		if env.Int32 != 2147483647 {
			t.Fatalf("expected 2147483647, got %d", env.Int32)
		}
		if env.Int64 != 9223372036854775807 {
			t.Fatalf("expected 9223372036854775807, got %d", env.Int64)
		}
		if env.Float32 != 3.4028234663852886e+38 {
			t.Fatalf("expected 3.4028234663852886e+38, got %f", env.Float32)
		}
		if env.Float64 != 1.7976931348623157e+308 {
			t.Fatalf("expected 1.7976931348623157e+308, got %f", env.Float64)
		}
	})
}
