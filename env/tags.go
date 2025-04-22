package env

import (
	"errors"
	"strings"
)

type def struct {
	// The name of the environment variable.
	name string

	// The default value for the field.
	def string

	// Is the field required?
	// Should an error be returned if the field is not set?
	required bool

	// Is the field a JSON object?
	// Should the value be unmarshalled as JSON?
	json bool
}

var ErrUnrecognizedTagOption = errors.New("unrecognized tag option")

func parseTag(tag string) (t def, err error) {
	t.name, t.def, t.required = strings.Cut(tag, "=")
	t.required = !t.required

	parts := strings.Split(t.def, ",")
	if len(parts) > 1 {
		t.def = parts[0]
		for _, part := range parts[1:] {
			switch part {
			case "json":
				t.json = true

			default:
				err = ErrUnrecognizedTagOption
				return
			}
		}
	}

	return
}
