package gqb_test

import (
	"context"
	"database/sql"
	"testing"
)

type sqlResultMock struct{}

func (s sqlResultMock) LastInsertId() (int64, error) {
	return 1, nil
}
func (s sqlResultMock) RowsAffected() (int64, error) {
	return 1, nil
}

type mockError struct{}

func (m mockError) Error() string {
	return "MockError"
}

type mockExecutor struct {
	query string
	binds []interface{}
}

func (m *mockExecutor) QueryContext(ctx context.Context, query string, binds ...interface{}) (*sql.Rows, error) {
	m.query = query
	m.binds = binds
	return nil, mockError{}
}
func (m *mockExecutor) ExecContext(ctx context.Context, query string, binds ...interface{}) (sql.Result, error) {
	m.query = query
	m.binds = binds
	return nil, mockError{}
}

func TestMysql(t *testing.T) {
	runMysqlTest(t)
	runPostgresTest(t)
}

func TestRemains(t *testing.T) {
	t.Errorf("We have to write tests about GROUP BY")
}
