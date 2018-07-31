package gqb

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"database/sql"
)

// QueryBuilder is struct for stack some conditions, orders, ... with method chain.
type QueryBuilder struct {
	db      Executor
	limit   int64
	offset  int64
	wheres  []ConditionBuilder
	orders  []Order
	selects []interface{}
	joins   []Join
	groupBy []string
}

// Create new Query QueryBuilder
func New(db Executor) *QueryBuilder {
	return &QueryBuilder{
		db: db,
	}
}

// Reset() resets stacks
func (q *QueryBuilder) Reset() {
	q.wheres = []ConditionBuilder{}
	q.selects = []interface{}{}
	q.joins = []Join{}
	q.orders = []Order{}
	q.groupBy = []string{}
	q.limit = 0
	q.offset = 0
}

// Add SELECT fields
func (q *QueryBuilder) Select(fields ...interface{}) *QueryBuilder {
	q.selects = append(q.selects, fields...)
	return q
}

// Add SELECT COUNT fields
func (q *QueryBuilder) SelectCount(field string) *QueryBuilder {
	q.selects = append(q.selects, Raw("COUNT("+quote(field)+")"))
	return q
}

// Add SELECT MAX fields
func (q *QueryBuilder) SelectMax(field string) *QueryBuilder {
	q.selects = append(q.selects, Raw("MAX("+quote(field)+")"))
	return q
}

// Add SELECT MIN fields
func (q *QueryBuilder) SelectMin(field string) *QueryBuilder {
	q.selects = append(q.selects, Raw("MIN("+quote(field)+")"))
	return q
}

// Add SELECT AVG fields
func (q *QueryBuilder) SelectAvg(field string) *QueryBuilder {
	q.selects = append(q.selects, Raw("AVG("+quote(field)+")"))
	return q
}

// Set LIMIT field
func (q *QueryBuilder) Limit(limit int64) *QueryBuilder {
	q.limit = limit
	return q
}

// Set OFFSET field
func (q *QueryBuilder) Offset(offset int64) *QueryBuilder {
	q.offset = offset
	return q
}

// Add JOIN table with condition
func (q *QueryBuilder) Join(table, from, to string, c Comparison) *QueryBuilder {
	q.joins = append(q.joins, Join{
		on: condition{
			comparison: c,
			field:      from,
			value:      to,
		},
		table: table,
	})
	return q
}

// Add WHERE condition group with AND.
// The first argument is generator function which accepts *ConditionGroup as argument.
// After call the generator function, add WHERE stack with called state
func (q *QueryBuilder) WhereGroup(generator func(g *WhereGroup)) *QueryBuilder {
	cg := newWhereGroup(And)
	generator(cg)
	q.wheres = append(q.wheres, cg)
	return q
}

// Add WHERE condition group with OR.
// The first argument is generator function which accepts *ConditionGroup as argument.
// After call the generator function, add WHERE stack with called state
func (q *QueryBuilder) OrWhereGroup(generator func(g *WhereGroup)) *QueryBuilder {
	cg := newWhereGroup(Or)
	generator(cg)
	q.wheres = append(q.wheres, cg)
	return q
}

// Add condition
func (q *QueryBuilder) AddWhere(c ConditionBuilder) *QueryBuilder {
	q.wheres = append(q.wheres, c)
	return q
}

// Add condition with AND combination
func (q *QueryBuilder) Where(field string, value interface{}, comparison Comparison) *QueryBuilder {
	return q.AddWhere(condition{
		comparison: comparison,
		field:      field,
		value:      value,
		combine:    And,
	})
}

// Add condition with OR combination
func (q *QueryBuilder) OrWhere(field string, value interface{}, comparison Comparison) *QueryBuilder {
	return q.AddWhere(condition{
		comparison: comparison,
		field:      field,
		value:      value,
		combine:    Or,
	})
}

// Add IN condition with AND combination
func (q *QueryBuilder) WhereIn(field string, values ...interface{}) *QueryBuilder {
	return q.AddWhere(condition{
		comparison: In,
		field:      field,
		value:      values,
		combine:    And,
	})
}

// Add IN condition with OR combination
func (q *QueryBuilder) OrWhereIn(field string, values ...interface{}) *QueryBuilder {
	return q.AddWhere(condition{
		comparison: In,
		field:      field,
		value:      values,
		combine:    Or,
	})
}

// Add NOT IN condition with AND combination
func (q *QueryBuilder) WhereNotIn(field string, values ...interface{}) *QueryBuilder {
	return q.AddWhere(condition{
		comparison: NotIn,
		field:      field,
		value:      values,
		combine:    And,
	})
}

// Add NOT IN condition with OR combination
func (q *QueryBuilder) OrWhereNotIn(field string, values ...interface{}) *QueryBuilder {
	return q.AddWhere(condition{
		comparison: NotIn,
		field:      field,
		value:      values,
		combine:    Or,
	})
}

// Add LIKE condition with AND combination
func (q *QueryBuilder) Like(field string, value interface{}) *QueryBuilder {
	return q.AddWhere(condition{
		comparison: Like,
		field:      field,
		value:      value,
		combine:    And,
	})
}

// Add LIKE condition with OR combination
func (q *QueryBuilder) OrLike(field string, value interface{}) *QueryBuilder {
	return q.AddWhere(condition{
		comparison: Like,
		field:      field,
		value:      value,
		combine:    Or,
	})
}

// Add NOT LIKE condition with AND combination
func (q *QueryBuilder) NotLike(field string, value interface{}) *QueryBuilder {
	return q.AddWhere(condition{
		comparison: NotLike,
		field:      field,
		value:      value,
		combine:    And,
	})
}

// Add NOT LIKE condition with OR combination
func (q *QueryBuilder) OrNotLike(field string, value interface{}) *QueryBuilder {
	return q.AddWhere(condition{
		comparison: NotLike,
		field:      field,
		value:      value,
		combine:    Or,
	})
}

// Add GROUP BY clause
func (q *QueryBuilder) GroupBy(fields ...string) *QueryBuilder {
	q.groupBy = append(q.groupBy, fields...)
	return q
}

// Add ORDER BY cluase
func (q *QueryBuilder) OrderBy(field string, sort SortMode) *QueryBuilder {
	q.orders = append(q.orders, Order{
		field: field,
		sort:  sort,
	})
	return q
}

// Format FROM table
func (q *QueryBuilder) formatTable(table interface{}) (string, error) {
	if v, ok := table.(alias); ok {
		return v.String(), nil
	} else if v, ok := table.(string); ok {
		if v == "" {
			return "", fmt.Errorf("Table name must not be empty")
		}
		return quote(v), nil
	}
	return "", fmt.Errorf("Invalid table specified")
}

// Execute query and get first result
func (q *QueryBuilder) GetOne(table interface{}) (*Result, error) {
	return q.GetOneContext(context.Background(), table)
}

// Execute query and get first result with context
func (q *QueryBuilder) GetOneContext(ctx context.Context, table interface{}) (*Result, error) {
	q.limit = 1
	r, err := q.GetContext(ctx, table)
	if err != nil {
		return nil, err
	} else if len(r) == 0 {
		return nil, sql.ErrNoRows
	}
	return r[0], nil
}

// Execute query and get results
func (q *QueryBuilder) Get(table interface{}) (Results, error) {
	return q.GetContext(context.Background(), table)
}

// Execute query and get results with context
func (q *QueryBuilder) GetContext(ctx context.Context, table interface{}) (Results, error) {
	mainTable, err := q.formatTable(table)
	if err != nil {
		return nil, err
	}
	where, binds := buildWhere(q.wheres, []interface{}{})
	query := strings.TrimSpace(fmt.Sprintf(
		"SELECT %s FROM %s%s%s%s%s%s%s",
		buildSelectFields(q.selects),
		mainTable,
		buildJoin(q.joins, mainTable),
		where,
		buildGroupBy(q.groupBy),
		buildOrderBy(q.orders),
		buildLimit(q.limit),
		buildOffset(q.offset),
	))

	defer q.Reset()
	rows, err := q.db.QueryContext(ctx, query, binds...)
	if err != nil {
		return nil, err
	}
	// gqb close rows pointer automatically so user don't need to care about it.
	// but allocate some more memories to make results
	defer rows.Close()
	return q.scan(rows)
}

// Scan rows to map to result
func (q *QueryBuilder) scan(rows *sql.Rows) (Results, error) {
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
			v := *value
			if b, ok := v.([]byte); ok {
				// Also we need to ensure value is zero value to avoid panic
				if reflect.ValueOf(value).IsValid() {
					values[name] = string(b)
					continue
				}
			}
			// Other types like int, float, decimal will treat as interface directory
			values[name] = v
		}
		results = append(results, NewResult(values))
	}
	return results, err
}

// Execute UPDATE query
func (q *QueryBuilder) Update(table interface{}, data Data) (sql.Result, error) {
	return q.UpdateContext(context.Background(), table, data)
}

// Execute UPDATE query with context
func (q *QueryBuilder) UpdateContext(ctx context.Context, table interface{}, data Data) (sql.Result, error) {
	if data == nil {
		return nil, fmt.Errorf("update data must be non-nil")
	}
	mainTable, err := q.formatTable(table)
	if err != nil {
		return nil, err
	}
	var where, updates string
	binds := []interface{}{}

	for _, k := range data.Keys() {
		updates += quote(k) + " = " + driverCompat.PlaceHolder(len(binds)+1) + ", "
		binds = bind(binds, data[k])
	}
	where, binds = buildWhere(q.wheres, binds)

	query := strings.TrimSpace(fmt.Sprintf(
		"UPDATE %s SET %s%s%s",
		mainTable,
		strings.TrimRight(updates, ", "),
		where,
		buildLimit(q.limit),
	))

	defer q.Reset()
	return q.db.ExecContext(ctx, query, binds...)
}

// Execute INSERT query
func (q *QueryBuilder) Insert(table interface{}, data Data) (sql.Result, error) {
	return q.InsertContext(context.Background(), table, data)
}

// Execute INSERT query with context
func (q *QueryBuilder) InsertContext(ctx context.Context, table interface{}, data Data) (sql.Result, error) {
	if data == nil {
		return nil, fmt.Errorf("insert data must be non-nil")
	}
	mainTable, err := q.formatTable(table)
	if err != nil {
		return nil, err
	}

	var fields, values string
	binds := []interface{}{}

	for _, k := range data.Keys() {
		fields += quote(k) + ", "
		values += driverCompat.PlaceHolder(len(binds)+1) + ", "
		binds = bind(binds, data[k])
	}
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		mainTable,
		strings.TrimRight(fields, ", "),
		strings.TrimRight(values, ", "),
	)
	defer q.Reset()
	return q.db.ExecContext(ctx, query, binds...)
}

// Execute bulk INSERT query
func (q *QueryBuilder) BulkInsert(table interface{}, data []Data) (sql.Result, error) {
	return q.BulkInsertContext(context.Background(), table, data)
}

// Execute bulk INSERT query with context
func (q *QueryBuilder) BulkInsertContext(ctx context.Context, table interface{}, data []Data) (sql.Result, error) {
	if data == nil {
		return nil, fmt.Errorf("insert data must be non-nil")
	}
	mainTable, err := q.formatTable(table)
	if err != nil {
		return nil, err
	}

	var fields string
	valueGroup := []string{}
	binds := []interface{}{}

	for i, d := range data {
		var values string
		for _, k := range d.Keys() {
			if i == 0 {
				fields += quote(k) + ", "
			}
			values += driverCompat.PlaceHolder(len(binds)+1) + ", "
			binds = bind(binds, d[k])
		}
		valueGroup = append(valueGroup, "("+strings.TrimRight(values, ", ")+")")
	}
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		mainTable,
		strings.TrimRight(fields, ", "),
		strings.Join(valueGroup, ", "),
	)
	defer q.Reset()
	return q.db.ExecContext(ctx, query, binds...)
}

// Execute DELETE query
func (q *QueryBuilder) Delete(table interface{}) (sql.Result, error) {
	return q.DeleteContext(context.Background(), table)
}

// Execute DELETE query with context
func (q *QueryBuilder) DeleteContext(ctx context.Context, table interface{}) (sql.Result, error) {
	mainTable, err := q.formatTable(table)
	if err != nil {
		return nil, err
	}
	where, binds := buildWhere(q.wheres, []interface{}{})
	query := strings.TrimSpace(fmt.Sprintf(
		"DELETE FROM %s%s",
		mainTable,
		where,
	))
	defer q.Reset()
	return q.db.ExecContext(ctx, query, binds...)
}
