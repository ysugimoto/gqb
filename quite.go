package gqb

import (
	"fmt"
	"strings"
	"time"
)

var quoteCharacter = `"`

func SetDriver(driverType string) {
	switch driverType {
	case "mysql":
		quoteCharacter = "`"
	default:
		quoteCharacter = `"`
	}
}

// quote() adds back quote prefix/suffix
func quote(str string) string {
	return quoteCharacter + strings.Trim(str, quoteCharacter) + quoteCharacter
}

func toString(v interface{}) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	} else if s, ok := v.(string); ok {
		return formatField(s)
	}
	return ""
}

// formatField() adds back quote by splitting table and column
func formatField(str string) string {
	if str == "" {
		return str
	}
	split := strings.Split(str, ".")
	for i, _ := range split {
		split[i] = quote(split[i])
	}
	return strings.Join(split, ".")
}

// bind() adds some value to bind slice values.
// if value is time.Time struct, stringify with datetime
func bind(b []interface{}, v interface{}) []interface{} {
	if t, ok := v.(time.Time); ok {
		b = append(b, t.Format(timeFormat))
	} else {
		b = append(b, v)
	}
	return b
}
