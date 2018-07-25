package gqb

import (
	"fmt"
	"reflect"
	"strings"

	"database/sql"
)

type Data map[string]interface{}
type Raw string
type queryMode int

const (
	Select queryMode = iota
	Update
	Insert
	Delete
)

type Executor interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

type Builder struct {
	db      Executor
	limit   int64
	offset  int64
	wheres  []conditionBuilder
	orders  []Order
	selects []interface{}
	joins   []Join
}

func New(db Executor) *Builder {
	return &Builder{
		db: db,
	}
}

func (q *Builder) Reset() {
	q.wheres = []conditionBuilder{}
	q.selects = []interface{}{}
	q.joins = []Join{}
	q.orders = []Order{}
	q.limit = 0
	q.offset = 0
}

func (q *Builder) Select(fields ...interface{}) *Builder {
	q.selects = append(q.selects, fields...)
	return q
}

func (q *Builder) Limit(limit int64) *Builder {
	q.limit = limit
	return q
}

func (q *Builder) Offset(offset int64) *Builder {
	q.offset = offset
	return q
}

func (q *Builder) Join(table, from, to string, c Comparison) *Builder {
	q.joins = append(q.joins, Join{
		on: Condition{
			comparison: c,
			field:      from,
			value:      to,
		},
		table: table,
	})
	return q
}

func (q *Builder) WhereGroup(generator func(g *ConditionGroup)) *Builder {
	cg := NewConditionGroup(And)
	generator(cg)
	q.wheres = append(q.wheres, cg)
	return q
}
func (q *Builder) OrWhereGroup(generator func(g *ConditionGroup)) *Builder {
	cg := NewConditionGroup(Or)
	generator(cg)
	q.wheres = append(q.wheres, cg)
	return q
}

func (q *Builder) Where(field string, value interface{}, comparison Comparison) *Builder {
	q.wheres = append(q.wheres, Condition{
		comparison: comparison,
		field:      field,
		value:      value,
		combine:    And,
	})
	return q
}

func (q *Builder) OrWhere(field string, value interface{}, comparison Comparison) *Builder {
	q.wheres = append(q.wheres, Condition{
		comparison: comparison,
		field:      field,
		value:      value,
		combine:    Or,
	})
	return q
}

func (q *Builder) WhereIn(field string, values ...interface{}) *Builder {
	q.wheres = append(q.wheres, Condition{
		comparison: In,
		field:      field,
		value:      values,
		combine:    And,
	})
	return q
}

func (q *Builder) Like(field string, value interface{}) *Builder {
	q.wheres = append(q.wheres, Condition{
		comparison: Like,
		field:      field,
		value:      value,
		combine:    And,
	})
	return q
}

func (q *Builder) OrLike(field string, value interface{}) *Builder {
	q.wheres = append(q.wheres, Condition{
		comparison: Like,
		field:      field,
		value:      value,
		combine:    Or,
	})
	return q
}

func (q *Builder) OrderBy(field string, sort SortMode) *Builder {
	q.orders = append(q.orders, Order{
		field: field,
		sort:  sort,
	})
	return q
}

func (q *Builder) GetOne(table string) (*Result, error) {
	defLimit := q.limit
	defer func() {
		q.limit = defLimit
	}()
	q.limit = 1
	r, err := q.Get(table)
	if err != nil {
		return nil, err
	} else if len(r) == 0 {
		return nil, sql.ErrNoRows
	}
	return r[0], nil
}

func (q *Builder) Get(table string) (Results, error) {
	query, binds, err := q.Build(Select, table, nil)
	if err != nil {
		return nil, err
	}
	rows, err := q.db.Query(query, binds...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return q.scan(rows)
}

func (q *Builder) scan(rows *sql.Rows) (Results, error) {
	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	results := Results{}
	for rows.Next() {
		scans := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			var s interface{}
			scans[i] = &s
		}
		if err := rows.Scan(scans...); err != nil {
			return nil, err
		}
		values := make(map[string]interface{})
		for i, n := range columns {
			name := n.Name()
			value := scans[i].(*interface{})
			// Check nil immediately
			if *value == nil {
				values[name] = *value
				continue
			}
			// We treat charater fields like VARCHAR, TEXT, ...
			// Because Go's sql driver scan as []byte for string type column on interface{},
			// so it's hard to convert to string on marshal JSON.
			t := n.ScanType()
			// Ensure []byte type
			if t.Kind() == reflect.Slice && t.Name() == "RawBytes" {
				// Also we need to ensure value is zero value to avoid panic
				if reflect.ValueOf(value).IsValid() {
					v := *value
					values[name] = string(v.([]byte))
					continue
				}
			}
			// Other types like int, float, decimal will treat as interface directory
			values[name] = *value
		}
		results = append(results, NewResult(values))
	}
	return results, err
}

func (q *Builder) Update(table string, data Data) error {
	query, binds, err := q.Build(Update, table, data)
	if err != nil {
		return err
	}
	_, err = q.db.Exec(query, binds...)
	return err
}

func (q *Builder) Insert(table string, data Data) error {
	query, binds, err := q.Build(Insert, table, data)
	if err != nil {
		return err
	}
	_, err = q.db.Exec(query, binds...)
	return err
}

func (q *Builder) Delete(table string) error {
	query, binds, err := q.Build(Delete, table, nil)
	if err != nil {
		return err
	}
	_, err = q.db.Exec(query, binds...)
	return err
}

func (q *Builder) Build(mode queryMode, table string, data Data) (string, []interface{}, error) {
	if table == "" {
		return "", nil, fmt.Errorf("table not specified")
	}
	switch mode {
	case Select:
		where, binds := buildWhere(q.wheres, []interface{}{})
		return strings.TrimSpace(fmt.Sprintf(
			"SELECT %s FROM %s%s%s%s%s%s",
			buildSelectFields(q.selects),
			quote(table),
			buildJoin(q.joins, table),
			where,
			buildOrderBy(q.orders),
			buildLimit(q.limit),
			buildOffset(q.offset),
		)), binds, nil
	case Update:
		if data == nil {
			return "", nil, fmt.Errorf("update data must be non-nil")
		}
		var where string
		updates, binds := buildUpdateFields(data, []interface{}{})
		where, binds = buildWhere(q.wheres, binds)

		return strings.TrimSpace(fmt.Sprintf(
			"UPDATE %s SET %s%s%s",
			quote(table),
			updates,
			where,
			buildLimit(q.limit),
		)), binds, nil
	case Insert:
		if data == nil {
			return "", nil, fmt.Errorf("insert data must be non-nil")
		}
		binds := []interface{}{}
		fields := ""
		values := ""
		for k, v := range data {
			fields += formatField(k) + ", "
			values += "?, "
			binds = bind(binds, v)
		}
		return fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s)",
			quote(table),
			strings.TrimRight(fields, ", "),
			strings.TrimRight(values, ", "),
		), binds, nil
	case Delete:
		where, binds := buildWhere(q.wheres, []interface{}{})
		return strings.TrimSpace(fmt.Sprintf(
			"DELETE FROM %s%s",
			quote(table),
			where,
		)), binds, nil
	default:
		return "", nil, fmt.Errorf("unexpected query mode specified")
	}
}
