package gqb

import (
	"fmt"
)

// WhereGroup is struct which wraps multiple Condition struct.
// This struct also implement ConditionBuilder inteface, so we can treat as same as Condition struct.
// This is used for grouped condition like "SELECT * FROM example WHERE A = 1 AND (B = 2 AND C = 3)"
// Parentheses inside condition is made by ConditionGroup
type WhereGroup struct {
	conditions  []ConditionBuilder
	combineType CombineType
}

func newWhereGroup(c CombineType) *WhereGroup {
	return &WhereGroup{
		conditions:  make([]ConditionBuilder, 0),
		combineType: c,
	}
}

// ConditionBuilder::Combine() interface implementation
func (w *WhereGroup) Combine() string {
	return string(w.combineType)
}

// ConditionBuilder::Build() interface implementation
func (w *WhereGroup) Build(binds []interface{}) (string, []interface{}) {
	first := true
	where := ""

	for _, cd := range w.conditions {
		c := cd.Combine()
		if c != "" {
			c = " " + c + " "
		}
		if first {
			c = ""
			first = false
		}
		var phrase string
		phrase, binds = cd.Build(binds)
		where += fmt.Sprintf("%s%s", c, phrase)
	}
	return where, binds
}

// Add condition
func (w *WhereGroup) AddWhere(c ConditionBuilder) *WhereGroup {
	w.conditions = append(w.conditions, c)
	return w
}

// Add condition with AND combination
func (w *WhereGroup) Where(field string, value interface{}, comparison Comparison) *WhereGroup {
	return w.AddWhere(condition{
		comparison: comparison,
		field:      field,
		value:      value,
		combine:    And,
	})
}

// Add condition with OR combination
func (w *WhereGroup) OrWhere(field string, value interface{}, comparison Comparison) *WhereGroup {
	return w.AddWhere(condition{
		comparison: comparison,
		field:      field,
		value:      value,
		combine:    Or,
	})
}

// Add IN condition with AND combination
func (w *WhereGroup) WhereIn(field string, values ...interface{}) *WhereGroup {
	return w.AddWhere(condition{
		comparison: In,
		field:      field,
		value:      values,
		combine:    And,
	})
}

// Add IN condition with OR combination
func (w *WhereGroup) OrWhereIn(field string, values ...interface{}) *WhereGroup {
	return w.AddWhere(condition{
		comparison: In,
		field:      field,
		value:      values,
		combine:    Or,
	})
}

// Add LIKE condition with AND combination
func (w *WhereGroup) Like(field string, value interface{}) *WhereGroup {
	return w.AddWhere(condition{
		comparison: Like,
		field:      field,
		value:      value,
		combine:    And,
	})
}

// Add LIKE condition with OR combination
func (w *WhereGroup) OrLike(field string, value interface{}) *WhereGroup {
	return w.AddWhere(condition{
		comparison: Like,
		field:      field,
		value:      value,
		combine:    Or,
	})
}

// Add user specific raw condition with AND combination
func (w *WhereGroup) WhereRaw(raw string) *WhereGroup {
	return w.AddWhere(rawCondition{
		rawClause: raw,
		combine:   And,
	})
}

// Add user specific raw condition with OR combination
func (w *WhereGroup) OrWhereRaw(raw string) *WhereGroup {
	return w.AddWhere(rawCondition{
		rawClause: raw,
		combine:   Or,
	})
}
