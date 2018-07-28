package gqb

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"database/sql"
)

// Build query mode constants definition
const (
	Select queryMode = iota
	Update
	Insert
	Delete
)

type (
	// Raw type indicates "raw query phrase", it means gqb won't quote any string, columns.
	// You should tableke carefully when use this type, but it's useful for use function like "COUNT(*)".
	Raw string

	// queryMode type is build query type switching.
	queryMode int

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

// Alias() returns formatted alias struct
func Alias(from, to string) alias {
	return alias{
		from: from,
		to:   to,
	}
}

// fmt.Stringer intetface satisfies
func (a alias) String() string {
	return formatField(a.from) + " AS " + formatField(a.to)
}

// fmt.Stringer intetface satisfies
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

// SQL executor interface, this is enough to implement QueryContext() and ExecContext().
// It's useful for running query in transation or not, because Executor accepts both of *sql.DB and *sql.Tx.
type Executor interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

// Builder is struct for stack some conditions, orders, ... with method chain.
type Builder struct {
	db      Executor
	limit   int64
	offset  int64
	wheres  []conditionBuilder
	orders  []Order
	selects []interface{}
	joins   []Join
}

// Create new Query Builder
func New(db Executor) *Builder {
	return &Builder{
		db: db,
	}
}

// Reset() resets stacks
func (q *Builder) Reset() {
	q.wheres = []conditionBuilder{}
	q.selects = []interface{}{}
	q.joins = []Join{}
	q.orders = []Order{}
	q.limit = 0
	q.offset = 0
}

// Add SELECT fields
func (q *Builder) Select(fields ...interface{}) *Builder {
	q.selects = append(q.selects, fields...)
	return q
}

// Set LIMIT field
func (q *Builder) Limit(limit int64) *Builder {
	q.limit = limit
	return q
}

// Set OFFSET field
func (q *Builder) Offset(offset int64) *Builder {
	q.offset = offset
	return q
}

// Add JOIN table with condition
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

// Add WHERE condition group with AND.
// The first argument is generator function which accepts *ConditionGroup as argument.
// After call the generator function, add WHERE stack with called state
func (q *Builder) WhereGroup(generator func(g *ConditionGroup)) *Builder {
	cg := NewConditionGroup(And)
	generator(cg)
	q.wheres = append(q.wheres, cg)
	return q
}

// Add WHERE condition group with OR.
// The first argument is generator function which accepts *ConditionGroup as argument.
// After call the generator function, add WHERE stack with called state
func (q *Builder) OrWhereGroup(generator func(g *ConditionGroup)) *Builder {
	cg := NewConditionGroup(Or)
	generator(cg)
	q.wheres = append(q.wheres, cg)
	return q
}

// Add condition
func (q *Builder) where(field string, value interface{}, comparison Comparison, combine combine) {
	q.wheres = append(q.wheres, Condition{
		comparison: comparison,
		field:      field,
		value:      value,
		combine:    combine,
	})
}

// Add condition with AND combination
func (q *Builder) Where(field string, value interface{}, comparison Comparison) *Builder {
	q.where(field, value, comparison, And)
	return q
}

// Add condition with OR combination
func (q *Builder) OrWhere(field string, value interface{}, comparison Comparison) *Builder {
	q.where(field, value, comparison, Or)
	return q
}

// Add IN condition with AND combination
func (q *Builder) WhereIn(field string, values ...interface{}) *Builder {
	q.where(field, values, In, And)
	return q
}

// Add LIKE condition with AND combination
func (q *Builder) Like(field string, value interface{}) *Builder {
	q.where(field, value, Like, And)
	return q
}

// Add LIKE condition with OR combination
func (q *Builder) OrLike(field string, value interface{}) *Builder {
	q.where(field, value, Like, Or)
	return q
}

// Add ORDER BY condition
func (q *Builder) OrderBy(field string, sort SortMode) *Builder {
	q.orders = append(q.orders, Order{
		field: field,
		sort:  sort,
	})
	return q
}

// Execute query and get first result
func (q *Builder) GetOne(table interface{}) (*Result, error) {
	return q.GetOneContext(context.Background(), table)
}

// Execute query and get first result with context
func (q *Builder) GetOneContext(ctx context.Context, table interface{}) (*Result, error) {
	defLimit := q.limit
	defer func() {
		q.limit = defLimit
	}()
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
func (q *Builder) Get(table interface{}) (Results, error) {
	return q.GetContext(context.Background(), table)
}

// Execute query and get results with context
func (q *Builder) GetContext(ctx context.Context, table interface{}) (Results, error) {
	query, binds, err := q.Build(Select, table, nil)
	if err != nil {
		return nil, err
	}
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

// Execute UPDATE query
func (q *Builder) Update(table interface{}, data Data) (sql.Result, error) {
	return q.UpdateContext(context.Background(), table, data)
}

// Execute UPDATE query with context
func (q *Builder) UpdateContext(ctx context.Context, table interface{}, data Data) (sql.Result, error) {
	query, binds, err := q.Build(Update, table, data)
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, binds...)
}

// Execute INSERT query
func (q *Builder) Insert(table interface{}, data Data) (sql.Result, error) {
	return q.InsertContext(context.Background(), table, data)
}

// Execute INSERT query with context
func (q *Builder) InsertContext(ctx context.Context, table interface{}, data Data) (sql.Result, error) {
	query, binds, err := q.Build(Insert, table, data)
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, binds...)
}

// Execute DELETE query
func (q *Builder) Delete(table interface{}) (sql.Result, error) {
	return q.DeleteContext(context.Background(), table)
}

// Execute DELETE query with context
func (q *Builder) DeleteContext(ctx context.Context, table interface{}) (sql.Result, error) {
	query, binds, err := q.Build(Delete, table, nil)
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, binds...)
}

// Build SQL string and bind paramteres corresponds to query mode
func (q *Builder) Build(mode queryMode, table interface{}, data Data) (string, []interface{}, error) {
	mainTable := toString(table)
	if mainTable == "" {
		return "", nil, fmt.Errorf("table not specified or empty")
	}
	switch mode {
	// Build SELECT query
	case Select:
		where, binds := buildWhere(q.wheres, []interface{}{})
		return strings.TrimSpace(fmt.Sprintf(
			"SELECT %s FROM %s%s%s%s%s%s",
			buildSelectFields(q.selects),
			mainTable,
			buildJoin(q.joins, mainTable),
			where,
			buildOrderBy(q.orders),
			buildLimit(q.limit),
			buildOffset(q.offset),
		)), binds, nil

	// Build UPDATE query
	case Update:
		if data == nil {
			return "", nil, fmt.Errorf("update data must be non-nil")
		}
		where := ""
		binds := []interface{}{}
		updates := ""
		for _, k := range data.Keys() {
			updates += formatField(k) + " = ?, "
			binds = bind(binds, data[k])
		}
		updates = strings.TrimRight(updates, ", ")
		where, binds = buildWhere(q.wheres, binds)

		return strings.TrimSpace(fmt.Sprintf(
			"UPDATE %s SET %s%s%s",
			mainTable,
			updates,
			where,
			buildLimit(q.limit),
		)), binds, nil

	// Build INSERT query
	case Insert:
		if data == nil {
			return "", nil, fmt.Errorf("insert data must be non-nil")
		}
		binds := []interface{}{}
		fields := ""
		values := ""
		for _, k := range data.Keys() {
			fields += formatField(k) + ", "
			values += "?, "
			binds = bind(binds, data[k])
		}
		return fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s)",
			mainTable,
			strings.TrimRight(fields, ", "),
			strings.TrimRight(values, ", "),
		), binds, nil

	// Build DELETE query
	case Delete:
		where, binds := buildWhere(q.wheres, []interface{}{})
		return strings.TrimSpace(fmt.Sprintf(
			"DELETE FROM %s%s",
			mainTable,
			where,
		)), binds, nil

	// Unexpected
	default:
		return "", nil, fmt.Errorf("unexpected query mode specified")
	}
}
