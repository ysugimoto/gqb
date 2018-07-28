package gqb

import (
	"context"
	"database/sql"
	"sort"
)

type (
	// Raw type indicates "raw query phrase", it means gqb won't quote any string, columns.
	// You should tableke carefully when use this type, but it's useful for use function like "COUNT(*)".
	Raw string

	// Data type is used for INSERT/UPDATE data definition.
	// This is suger syntax for map[string]interface{}, but always fields are sorted by key.
	Data map[string]interface{}

	// alias type is used for SELECT, create alias column name.
	// This will be useful for using JOIN query.
	alias struct {
		from string
		to   string
	}
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
	CombineType string
)

const (
	// Equal compares equavbalence between column and value
	Equal Comparison = "="
	// Alias for Equal
	Eq Comparison = "="

	// NotEqual compares not equavbalence between column and value
	NotEqual Comparison = "<>"
	// Alias for NotEqual
	NotEq Comparison = "<>"

	// Gt compares greater than between column and value
	Gt Comparison = ">"

	// Gte compares greater than equal between column and value
	Gte Comparison = "<="

	// Lt compares less than between column and value
	Lt Comparison = "<"

	// Lte compares less than equal between column and value
	Lte Comparison = "<="

	// In compares within values
	In Comparison = "IN"

	// Like compares value matching phrase
	Like Comparison = "LIKE"

	// Desc indicates decendant
	Desc SortMode = "DESC"

	// Asc indicates ascendant
	Asc SortMode = "ASC"

	// Rand indicates random
	Rand SortMode = "RAND"

	// And concats conditions with AND
	And CombineType = "AND"

	// Or concats conditions with OR
	Or CombineType = "OR"
)

// conditionBuilder is private interface with create WHERE condition string.
type ConditionBuilder interface {
	// buildCondition() builds WHERE condition string and append bind parameters.
	Build([]interface{}) (string, []interface{})

	// getCombine() should return concatenation string AND/OR
	Combine() string
}

// SQL executor interface, this is enough to implement QueryContext() and ExecContext().
// It's useful for running query in transation or not, because Executor accepts both of *sql.DB and *sql.Tx.
type Executor interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

func Alias(from, to string) alias {
	return alias{
		from: from,
		to:   to,
	}
}

// fmt.Stringer intetface implementation
func (a alias) String() string {
	return quote(a.from) + " AS " + quote(a.to)
}

// fmt.Stringer intetface implementation
func (r Raw) String() string {
	return string(r)
}

// Return sorted field name strings
func (d Data) Keys() []string {
	keys := make([]string, len(d))
	i := 0
	for k := range d {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// Order is struct for making ORDER BY phrase
type Order struct {
	sort  SortMode
	field string
}

// Join is struct for making JOIN phrase
type Join struct {
	on    condition
	table string
}
