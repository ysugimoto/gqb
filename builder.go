package gqb

import (
	"fmt"
	"strings"
	"time"
)

const (
	// Time format, this is used for DATETIME column
	timeFormat = "2006-01-02 15:04:05"

	// Date format, this is used for DATE column
	dateFormat = "2006-01-02"
)

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

// Create SELECT column name string.
// If field is Raw type, field won't escape in order to unexpected quote string is added.
//
// Raw type -> Raw("COUNT(id)") -> COUNT(id)
// Others   -> name             -> `name`
func buildSelectFields(selects []interface{}) string {
	if len(selects) == 0 {
		return "*"
	}
	fields := ""
	for _, f := range selects {
		if v, ok := f.(Raw); ok {
			fields += v.String() + ", "
		} else if v, ok := f.(alias); ok {
			fields += v.String() + ", "
		} else if v, ok := f.(string); ok {
			fields += quote(v) + ", "
		}
	}
	return strings.TrimRight(fields, ", ")
}

// Create WHERE clause string.
// gqb uses prepared statement with "?", and add bind parameters slice
// If field is Raw type, field won't escape in order to unexpected quote string is added.
//
// Raw type          -> Raw("COUNT(id)") -> COUNT(id)
// column            -> name             -> `name`
// column with table -> table.name       -> `table`.`name`
func buildWhere(wheres []ConditionBuilder, binds []interface{}) (string, []interface{}) {
	if len(wheres) == 0 {
		return "", binds
	}

	first := true
	where := ""
	c := ""

	for _, w := range wheres {
		c = w.Combine()
		if c != "" {
			c = " " + c + " "
		}
		if first {
			c = ""
			first = false
		}
		var clause string
		clause, binds = w.Build(binds)
		where += fmt.Sprintf("%s(%s)", c, clause)
	}
	return " WHERE " + where, binds
}

// Create ORDER BY clause string.
func buildOrderBy(orders []Order) string {
	if len(orders) == 0 {
		return ""
	}
	order := []string{}
	for _, o := range orders {
		var s string
		if o.sort == Rand {
			s = driverCompat.RandFunc()
		} else {
			s = string(o.sort)
		}
		order = append(order, quote(o.field)+" "+s)
	}
	return " ORDER BY " + strings.Join(order, ", ")
}

// Create JOIN clause string.
func buildJoin(joins []Join, baseTable string) string {
	if len(joins) == 0 {
		return ""
	}

	join := ""
	for _, j := range joins {
		join += fmt.Sprintf(
			" JOIN %s ON (%s.%s %s %s.%s)",
			quote(j.table),
			quote(baseTable),
			quote(j.on.field),
			string(j.on.comparison),
			quote(j.table),
			quote(j.on.value.(string)),
		)
	}
	return join
}

// Create LIMIT clause string.
func buildLimit(limit int64) string {
	if limit == 0 {
		return ""
	}
	return fmt.Sprintf(" LIMIT %d", limit)
}

// Create OFFSET clause string.
func buildOffset(offset int64) string {
	if offset == 0 {
		return ""
	}
	return fmt.Sprintf(" OFFSET %d", offset)
}

// Create GROUP BY clause string.
func buildGroupBy(groupBy []string) string {
	if len(groupBy) == 0 {
		return ""
	}

	var gb string
	for _, g := range groupBy {
		gb += quote(g) + ", "
	}
	return " GROUP BY " + strings.TrimRight(gb, ", ")
}
