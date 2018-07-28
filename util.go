package gqb

import (
	"fmt"
	"reflect"

	"github.com/ysugimoto/gqb/compat"
)

var driverCompat compat.Compat = compat.MysqlCompat{}

func SetDriver(driverType string) {
	switch driverType {
	case "mysql":
		driverCompat = compat.MysqlCompat{}
	default:
		driverCompat = compat.PostgresCompat{}
	}
}

// shorthand syntax for compat.Compat.Quote
func quote(str string) string {
	return driverCompat.Quote(str)
}

// parseTag() parses Strcut tag to name-value map
func parseTag(tag string) (map[string]string, error) {
	parsed := make(map[string]string)
	var stack string
	var valueStart bool
	var key string
	for i, b := range []byte(tag) {
		switch b {
		case ':':
			if stack == "" {
				return nil, fmt.Errorf(`syntax error: unexpected ":" is present %s on %d`, tag, i)
			}
			key = stack
			stack = ""
		case '"':
			if !valueStart {
				valueStart = true
			} else {
				parsed[key] = stack
				valueStart = false
			}
			stack = ""
		case ' ':
			continue
		default:
			stack += string(b)
		}
	}
	if stack != "" {
		return nil, fmt.Errorf(`syntax error: invalid sting is remaining: %s`, tag)
	}
	return parsed, nil
}

// derefValue() dereference reflect.Value
func derefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// derefType() dereference reflect.Type
func derefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
