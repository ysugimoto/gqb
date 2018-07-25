package gqb

import (
	"fmt"
	"strings"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"
const dateFormat = "2006-01-02"

func quote(str string) string {
	return "`" + strings.Trim(str, "`") + "`"
}

func formatField(str string) string {
	split := strings.Split(str, ".")
	for i, _ := range split {
		split[i] = quote(split[i])
	}
	return strings.Join(split, ".")
}

func bind(b []interface{}, v interface{}) []interface{} {
	if t, ok := v.(time.Time); ok {
		b = append(b, t.Format(timeFormat))
	} else {
		b = append(b, v)
	}
	return b
}

func buildSelectFields(selects []interface{}) string {
	fields := []string{}
	if len(selects) > 0 {
		for _, f := range selects {
			switch f.(type) {
			case Raw:
				v := f.(Raw)
				fields = append(fields, string(v))
			case string:
				v := f.(string)
				fields = append(fields, formatField(v))
			}
		}
	} else {
		fields = []string{quote("*")}
	}
	return strings.Join(fields, ", ")
}

func buildUpdateFields(data Data, binds []interface{}) (string, []interface{}) {
	fields := []string{}
	for k, v := range data {
		fields = append(fields, formatField(k)+" = ?")
		binds = bind(binds, v)
	}
	return strings.Join(fields, ", "), binds
}

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

func buildLimit(limit int64) string {
	if limit == 0 {
		return ""
	}
	return fmt.Sprintf(" LIMIT %d", limit)
}

func buildOffset(offset int64) string {
	if offset == 0 {
		return ""
	}
	return fmt.Sprintf(" OFFSET %d", offset)
}
