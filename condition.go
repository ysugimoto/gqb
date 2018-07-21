package gqb

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

type Condition struct {
	comparison Comparison
	field      string
	value      interface{}

	combine Combine
}

func NewCondition(field string, value interface{}, c Comparison) *Condition {
	return &Condition{
		comparison: c,
		field:      field,
		value:      value,
	}
}

type Order struct {
	sort  SortMode
	field string
}

type Join struct {
	on    Condition
	table string
}
