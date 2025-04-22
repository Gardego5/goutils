package env

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
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
func Load[T any](optionFuncs ...OptionsFunc) (T, error) {
	options := DEFAULT_OPTIONS
	for _, fn := range optionFuncs {
		fn(&options)
	}

	var cfg T

	typ := reflect.TypeOf(cfg)
	if typ.Kind() != reflect.Struct {
		return cfg, fmt.Errorf("%s is not a struct", typ)
	}

	errs := []error{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag, ok := field.Tag.Lookup(options.TagName)
		if !ok {
			continue
		}

		t, err := parseTag(tag)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to parse tag %s: %w", tag, err))
			continue
		}

		value, ok := os.LookupEnv(t.name)
		if !ok {
			if t.required {
				errs = append(errs, fmt.Errorf("env '%s' required, but not found", t.name))
				continue
			} else {
				value = t.def
			}
		}

		fieldValue := reflect.ValueOf(&cfg).Elem().Field(i)
		val := fieldValue.Addr().Interface()
		if t.json {
			err := json.Unmarshal([]byte(value), val)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to unmarshal JSON %s: %w", value, err))
				continue
			}
		}

		switch val := val.(type) { // var val any := &cfg.Field
		case encoding.TextUnmarshaler:
			val.UnmarshalText([]byte(value))

		default:
			switch fieldValue.Type().Kind() {

			case reflect.Bool:
				valueBool, err := strconv.ParseBool(value)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse bool %s: %w", value, err))
					continue
				}
				fieldValue.SetBool(valueBool)

			case reflect.Int:
				valueInt, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse int %s: %w", value, err))
					continue
				}
				fieldValue.SetInt(valueInt)

			case reflect.Int8:
				valueInt, err := strconv.ParseInt(value, 10, 8)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse int8 %s: %w", value, err))
					continue
				}
				fieldValue.SetInt(valueInt)

			case reflect.Int16:
				valueInt, err := strconv.ParseInt(value, 10, 16)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse int16 %s: %w", value, err))
					continue
				}
				fieldValue.SetInt(valueInt)

			case reflect.Int32:
				valueInt, err := strconv.ParseInt(value, 10, 32)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse int32 %s: %w", value, err))
					continue
				}
				fieldValue.SetInt(valueInt)

			case reflect.Int64:
				valueInt, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse int64 %s: %w", value, err))
					continue
				}
				fieldValue.SetInt(valueInt)

			case reflect.Uint:
				valueUint, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse uint %s: %w", value, err))
					continue
				}
				fieldValue.SetUint(valueUint)

			case reflect.Uint8:
				valueUint, err := strconv.ParseUint(value, 10, 8)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse uint8 %s: %w", value, err))
					continue
				}
				fieldValue.SetUint(valueUint)

			case reflect.Uint16:
				valueUint, err := strconv.ParseUint(value, 10, 16)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse uint16 %s: %w", value, err))
					continue
				}
				fieldValue.SetUint(valueUint)

			case reflect.Uint32:
				valueUint, err := strconv.ParseUint(value, 10, 32)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse uint32 %s: %w", value, err))
					continue
				}
				fieldValue.SetUint(valueUint)

			case reflect.Uint64:
				valueUint, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse uint64 %s: %w", value, err))
					continue
				}
				fieldValue.SetUint(valueUint)

			case reflect.Float32:
				valueFloat, err := strconv.ParseFloat(value, 32)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse float32 %s: %w", value, err))
					continue
				}
				fieldValue.SetFloat(valueFloat)

			case reflect.Float64:
				valueFloat, err := strconv.ParseFloat(value, 64)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to parse float64 %s: %w", value, err))
					continue
				}
				fieldValue.SetFloat(valueFloat)

			case reflect.String:
				fieldValue.SetString(value)

			default:
				errs = append(errs, fmt.Errorf("unsupported type %s", typ))
				continue
			}
		}
	}

	return cfg, errors.Join(errs...)
}

func MustLoad[T any](optionFuncs ...OptionsFunc) T {
	cfg, err := Load[T](optionFuncs...)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	return cfg
}
