package gqb

import (
	"fmt"
	"strings"
)

type Comparison string
type SortMode string
type Combine string

const (
	Equal    Comparison = "="
	NotEqual Comparison = "<>"
	Gt       Comparison = ">"
	Gte      Comparison = "<="
	Lt       Comparison = "<"
	Lte      Comparison = "<="
	In       Comparison = "IN"
	Like     Comparison = "LIKE"

	Desc SortMode = "DESC"
	Asc  SortMode = "ASC"

	And Combine = "AND"
	Or  Combine = "OR"
)

type conditionBuilder interface {
	buildCondition([]interface{}) (string, []interface{})
	getCombine() string
	getComparison() string
}

type Condition struct {
	comparison Comparison
	field      string
	value      interface{}

	combine Combine
}

func (c Condition) getCombine() string {
	return string(c.combine)
}

func (c Condition) getComparison() string {
	return string(c.comparison)
}

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

type Order struct {
	sort  SortMode
	field string
}

type Join struct {
	on    Condition
	table string
}

type ConditionGroup struct {
	conditions []conditionBuilder
	combine    Combine
}

func NewConditionGroup(c Combine) *ConditionGroup {
	return &ConditionGroup{
		conditions: make([]conditionBuilder, 0),
		combine:    c,
	}
}

func (g *ConditionGroup) Where(field string, value interface{}, comparison Comparison) *ConditionGroup {
	g.conditions = append(g.conditions, Condition{
		comparison: comparison,
		field:      field,
		value:      value,
		combine:    And,
	})
	return g
}

func (g *ConditionGroup) OrWhere(field string, value interface{}, comparison Comparison) *ConditionGroup {
	g.conditions = append(g.conditions, Condition{
		comparison: comparison,
		field:      field,
		value:      value,
		combine:    Or,
	})
	return g
}

func (g *ConditionGroup) WhereIn(field string, values ...interface{}) *ConditionGroup {
	g.conditions = append(g.conditions, Condition{
		comparison: In,
		field:      field,
		value:      values,
		combine:    And,
	})
	return g
}

func (g *ConditionGroup) Like(field string, value interface{}) *ConditionGroup {
	g.conditions = append(g.conditions, Condition{
		comparison: Like,
		field:      field,
		value:      value,
		combine:    And,
	})
	return g
}

func (g *ConditionGroup) getCombine() string {
	return string(g.combine)
}
func (g *ConditionGroup) getComparison() string {
	return ""
}

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
