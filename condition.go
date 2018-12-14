package gqb

import (
	"fmt"
	"strings"
)

// Condition is common condition struct
type condition struct {
	comparison Comparison
	field      string
	value      interface{}
	combine    CombineType
}

// conditionBuilder::getCombine() interface implementation
func (c condition) Combine() string {
	return string(c.combine)
}

// conditionBuilder::Build() interface implementation
func (c condition) Build(binds []interface{}) (string, []interface{}) {
	var clause string

	switch c.comparison {
	case In:
		q := ""
		values, ok := c.value.([]interface{})
		if ok {
			for _, v := range values {
				q += driverCompat.PlaceHolder(len(binds)+1) + ", "
				binds = bind(binds, v)
			}
			clause = fmt.Sprintf("%s IN (%s)", quote(c.field), strings.Trim(q, ", "))
		}
	case Equal:
		if c.value == nil {
			clause = fmt.Sprintf("%s IS NULL", quote(c.field))
		} else {
			clause = fmt.Sprintf("%s %s %s", quote(c.field), string(c.comparison), driverCompat.PlaceHolder(len(binds)+1))
			binds = bind(binds, c.value)
		}
	case NotEqual:
		if c.value == nil {
			clause = fmt.Sprintf("%s IS NOT NULL", quote(c.field))
		} else {
			clause = fmt.Sprintf("%s %s %s", quote(c.field), string(c.comparison), driverCompat.PlaceHolder(len(binds)+1))
			binds = bind(binds, c.value)
		}
	default:
		clause = fmt.Sprintf("%s %s %s", quote(c.field), string(c.comparison), driverCompat.PlaceHolder(len(binds)+1))
		binds = bind(binds, c.value)
	}
	return clause, binds
}

// Raw clause condition
type rawCondition struct {
	rawClause string
	combine   CombineType
}

// conditionBuilder::getCombine() interface implementation
func (r rawCondition) Combine() string {
	return string(c.combine)
}

// conditionBuilder::Build() interface implementation
func (r rawCondition) Build(binds []interface{}) (string, []interface{}) {
	return r.rawClause, binds
}
