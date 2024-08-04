package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var DEFAULT_OPTIONS = options{
	TagName: "env",
}

// Load the environment variables into a struct.
// To use, call with a generic type and optional options.
// Example:
//
//	type Config struct {
//		Host string `env:"HOST"`
//		Port int `env:"PORT=8080"`
//	}
//	cfg, err := Load(Config{}, TagName("env"))
//	if err != nil {
//		log.Fatal(err)
//	}
func Load[T any](optionFuncs ...OptionsFunc) (*T, error) {
	options := DEFAULT_OPTIONS
	for _, fn := range optionFuncs {
		fn(&options)
	}

	var cfg T

	typ := reflect.TypeOf(cfg)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s is not a struct", typ)
	}

	errs := []error{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(options.TagName)
		if tag == "" {
			continue
		}

		parts := strings.SplitN(tag, "=", 2)
		var name, def string
		switch len(parts) {
		case 0:
			errs = append(errs, fmt.Errorf("empty tag"))
			continue

		case 1:
			name = parts[0]

		case 2:
			name = parts[0]
			def = parts[1]
		}

		value, ok := os.LookupEnv(name)
		if !ok {
			if def != "" {
				value = def
			} else {
				errs = append(errs, fmt.Errorf("env %s not found", name))
				continue
			}
		}

		fieldValue := reflect.ValueOf(&cfg).Elem().Field(i)
		switch fieldValue.Type().Kind() {
		case reflect.String:
			fieldValue.SetString(value)

		case reflect.Int:
			valueInt, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to parse int %s: %w", value, err))
				continue
			}
			fieldValue.SetInt(valueInt)

		default:
			errs = append(errs, fmt.Errorf("unsupported type %s", typ))
			continue
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	} else {
		return &cfg, nil
	}
}
