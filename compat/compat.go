package compat

import (
	"strings"
)

type Compat interface {
	Quote(string) string
	RandFunc() string
}

type MysqlCompat struct {
}

func (c MysqlCompat) Quote(str string) string {
	if str == "" {
		return str
	}
	split := strings.Split(str, ".")
	for i, _ := range split {
		split[i] = "`" + strings.Trim(split[i], "`") + "`"
	}
	return strings.Join(split, ".")
}

func (c MysqlCompat) RandFunc() string {
	return "RAND()"
}

type PostgresCompat struct {
}

func (c PostgresCompat) Quote(str string) string {
	if str == "" {
		return str
	}
	split := strings.Split(str, ".")
	for i, _ := range split {
		split[i] = `"` + strings.Trim(split[i], `"`) + `"`
	}
	return strings.Join(split, ".")
}

func (c PostgresCompat) RandFunc() string {
	return "RANDOM()"
}
