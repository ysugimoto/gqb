package gqb

import (
	"fmt"
	"strings"
)

const (
	// Time format, this is used for DATETIME column
	timeFormat = "2006-01-02 15:04:05"

	// Date format, this is used for DATE column
	dateFormat = "2006-01-02"
)

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
		if v := toString(f); v != "" {
			fields += v + ", "
		}
	}
	return strings.TrimRight(fields, ", ")
}

// Create WHERE phrase string.
// gqb uses prepared statement with "?", and add bind parameters slice
// If field is Raw type, field won't escape in order to unexpected quote string is added.
//
// Raw type          -> Raw("COUNT(id)") -> COUNT(id)
// column            -> name             -> `name`
// column with table -> table.name       -> `table`.`name`
func buildWhere(wheres []conditionBuilder, binds []interface{}) (string, []interface{}) {
	if len(wheres) == 0 {
		return "", binds
	}

	first := true
	where := ""
	c := ""

	for _, w := range wheres {
		c = w.getCombine()
		if c != "" {
			c = " " + c + " "
		}
		if first {
			c = ""
			first = false
		}
		var phrase string
		phrase, binds = w.buildCondition(binds)
		where += fmt.Sprintf("%s(%s)", c, phrase)
	}
	return " WHERE " + where, binds
}

// Create ORDER BY phrase string.
func buildOrderBy(orders []Order) string {
	if len(orders) == 0 {
		return ""
	}
	order := []string{}
	for _, o := range orders {
		order = append(order, formatField(o.field)+" "+string(o.sort))
	}
	return " ORDER BY " + strings.Join(order, ", ")
}

// Create JOIN phrase string.
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

// Create LIMIT phrase string.
func buildLimit(limit int64) string {
	if limit == 0 {
		return ""
	}
	return fmt.Sprintf(" LIMIT %d", limit)
}

// Create OFFSET phrase string.
func buildOffset(offset int64) string {
	if offset == 0 {
		return ""
	}
	return fmt.Sprintf(" OFFSET %d", offset)
}
