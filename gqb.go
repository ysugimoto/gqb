package gqb

import (
	"fmt"
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
	wheres  []Condition
	orders  []Order
	selects []interface{}
	joins   []Join
	table   string
}

func New(db Executor) *Builder {
	return &Builder{
		db: db,
	}
}

func (q *Builder) Reset() {
	q.wheres = []Condition{}
	q.selects = []interface{}{}
	q.joins = []Join{}
	q.orders = []Order{}
	q.limit = 0
	q.offset = 0
}

func (q *Builder) Table(table string) *Builder {
	q.table = table
	return q
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

func (q *Builder) OrderBy(field string, sort SortMode) *Builder {
	q.orders = append(q.orders, Order{
		field: field,
		sort:  sort,
	})
	return q
}

func (q *Builder) GetOne() (*Result, error) {
	defLimit := q.limit
	defer func() {
		q.limit = defLimit
	}()
	q.limit = 1
	r, err := q.Get()
	if err != nil {
		return nil, err
	} else if len(r) == 0 {
		return nil, sql.ErrNoRows
	}
	return r[0], nil
}

func (q *Builder) Get() ([]*Result, error) {
	query, binds, err := q.Build(Select, nil)
	if err != nil {
		return nil, err
	}
	rows, err := q.db.Query(query, binds...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	results := []*Result{}
	for rows.Next() {
		scans := []interface{}{}
		for i := 0; i < len(columns); i++ {
			var s interface{}
			scans = append(scans, &s)
		}
		if err := rows.Scan(scans...); err != nil {
			return nil, err
		}
		v := make(map[string]interface{})
		for i, n := range columns {
			v[n] = *scans[i].(*interface{})
		}
		results = append(results, &Result{
			values: v,
		})
	}
	return results, err
}

func (q *Builder) Update(data Data) error {
	query, binds, err := q.Build(Update, data)
	if err != nil {
		return err
	}
	_, err = q.db.Exec(query, binds...)
	return err
}

func (q *Builder) Insert(data Data) error {
	query, binds, err := q.Build(Insert, data)
	if err != nil {
		return err
	}
	_, err = q.db.Exec(query, binds...)
	return err
}

func (q *Builder) Delete() error {
	query, binds, err := q.Build(Delete, nil)
	if err != nil {
		return err
	}
	_, err = q.db.Exec(query, binds...)
	return err
}

func (q *Builder) Build(mode queryMode, data Data) (string, []interface{}, error) {
	if q.table == "" {
		return "", nil, fmt.Errorf("table not specified")
	}
	switch mode {
	case Select:
		where, binds := buildWhere(q.wheres, []interface{}{})
		return strings.TrimSpace(fmt.Sprintf(
			"SELECT %s FROM %s%s%s%s%s%s",
			buildSelectFields(q.selects),
			quote(q.table),
			buildJoin(q.joins, q.table),
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
			quote(q.table),
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
			quote(q.table),
			strings.TrimRight(fields, ", "),
			strings.TrimRight(values, ", "),
		), binds, nil
	case Delete:
		where, binds := buildWhere(q.wheres, []interface{}{})
		return strings.TrimSpace(fmt.Sprintf(
			"DELETE FROM %s%s",
			quote(q.table),
			where,
		)), binds, nil
	default:
		return "", nil, fmt.Errorf("unexpected query mode specified")
	}
}
