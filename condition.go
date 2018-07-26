package gqb

import (
	"fmt"
	"strings"
)

type (
	// Comparison type indicates how to compare between column and value.
	// This type is used for WHERE.
	Comparison string

	// SortMode type indicates how to sort records.
	// This type is used for ORDER BY
	SortMode string

	// combine type indicates how to concat multiple WHERE conditions.
	// This type is used for WHERE and private type
	combine string
)

const (
	// Equal compares equavbalence between column and value
	Equal Comparison = "="

	// NotEqual compares not equavbalence between column and value
	NotEqual Comparison = "<>"

	// Gt compares greater than between column and value
	Gt Comparison = ">"

	// Gt compares greater than equal between column and value
	Gte Comparison = "<="

	// Lt compares less than between column and value
	Lt Comparison = "<"

	// Lt compares less than equal between column and value
	Lte Comparison = "<="

	// In compares within values
	In Comparison = "IN"

	// Like compares value matching phrase
	Like Comparison = "LIKE"

	// Desc indicates decendant
	Desc SortMode = "DESC"

	// Desc indicates ascendant
	Asc SortMode = "ASC"

	// And concats conditions with AND
	And combine = "AND"

	// Or concats conditions with OR
	Or combine = "OR"
)

// conditionBuilder is private interface with create WHERE condition string.
type conditionBuilder interface {
	// buildCondition() builds WHERE condition string and append bind parameters.
	buildCondition([]interface{}) (string, []interface{})

	// getCombine() should return concatenation string AND/OR
	getCombine() string
}

// Condition is common condition struct
type Condition struct {
	comparison Comparison
	field      string
	value      interface{}
	combine    combine
}

// conditionBuilder::getCombine() interface implementation
func (c Condition) getCombine() string {
	return string(c.combine)
}

// conditionBuilder::buildCondition() interface implementation
func (c Condition) buildCondition(binds []interface{}) (string, []interface{}) {
	var phrase string

	switch c.comparison {
	case In:
		q := ""
		values, ok := c.value.([]interface{})
		if ok {
			for _, v := range values {
				q += "?, "
				binds = bind(binds, v)
			}
			phrase = fmt.Sprintf("%s IN (%s)", formatField(c.field), strings.Trim(q, ", "))
		}
	default:
		phrase = fmt.Sprintf("%s %s ?", formatField(c.field), string(c.comparison))
		binds = bind(binds, c.value)
	}
	return phrase, binds
}

// Order is struct for making ORDER BY phrase
type Order struct {
	sort  SortMode
	field string
}

// Join is struct for making JOIN phrase
type Join struct {
	on    Condition
	table string
}

// ConditionGroup is struct which wraps multiple Condition struct.
// This struct also implement conditionBuilder inteface, so we can treat as same as Condition struct.
// This is used for grouped condition like "SELECT * FROM example WHERE A = 1 AND (B = 2 AND C = 3)"
// Parentheses inside condition is made by ConditionGroup
type ConditionGroup struct {
	conditions []conditionBuilder
	combine    combine
}

// Create ConditionGroup pointer
func NewConditionGroup(c combine) *ConditionGroup {
	return &ConditionGroup{
		conditions: make([]conditionBuilder, 0),
		combine:    c,
	}
}

// Add condition
func (g *ConditionGroup) where(field string, value interface{}, comparison Comparison, combine combine) {
	g.conditions = append(g.conditions, Condition{
		comparison: comparison,
		field:      field,
		value:      value,
		combine:    combine,
	})
}

// Add condition with AND combination
func (g *ConditionGroup) Where(field string, value interface{}, comparison Comparison) *ConditionGroup {
	g.where(field, value, comparison, And)
	return g
}

// Add condition with OR combination
func (g *ConditionGroup) OrWhere(field string, value interface{}, comparison Comparison) *ConditionGroup {
	g.where(field, value, comparison, Or)
	return g
}

// Add condition with IN
func (g *ConditionGroup) WhereIn(field string, values ...interface{}) *ConditionGroup {
	g.where(field, values, In, And)
	return g
}

// Add condition with LIKE
func (g *ConditionGroup) Like(field string, value interface{}) *ConditionGroup {
	g.where(field, value, Like, And)
	return g
}

// Add condition with OR LIKE
func (g *ConditionGroup) OrLike(field string, value interface{}) *ConditionGroup {
	g.where(field, value, Like, Or)
	return g
}

// conditionBuilder::getCombine() interface implementation
func (g *ConditionGroup) getCombine() string {
	return string(g.combine)
}

// conditionBuilder::buildCondition() interface implementation
func (g *ConditionGroup) buildCondition(binds []interface{}) (string, []interface{}) {
	first := true
	where := ""

	for _, cd := range g.conditions {
		c := cd.getCombine()
		if c != "" {
			c = " " + c + " "
		}
		if first {
			c = ""
			first = false
		}
		var phrase string
		phrase, binds = cd.buildCondition(binds)
		where += fmt.Sprintf("%s%s", c, phrase)
	}
	return where, binds
}
