package gqb

import (
	"fmt"
	"strings"
)

type Compat interface {
	Quote(string) string
	RandFunc() string
	PlaceHolder(int) string
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

func (c MysqlCompat) PlaceHolder(index int) string {
	return "?"
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

func (c PostgresCompat) PlaceHolder(index int) string {
	return fmt.Sprintf("$%d", index)
}

type SQLiteCompat struct {
}

func (c SQLiteCompat) Quote(str string) string {
	if str == "" {
		return str
	}
	split := strings.Split(str, ".")
	for i, _ := range split {
		split[i] = `"` + strings.Trim(split[i], `"`) + `"`
	}
	return strings.Join(split, ".")
}

func (c SQLiteCompat) RandFunc() string {
	return "RANDOM()"
}

func (c SQLiteCompat) PlaceHolder(index int) string {
	return "?"
}
