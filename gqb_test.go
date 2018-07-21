package gqb_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysugimoto/gqb"
)

type sqlResultMock struct{}

func (s sqlResultMock) LastInsertId() (int64, error) {
	return 1, nil
}
func (s sqlResultMock) RowsAffected() (int64, error) {
	return 1, nil
}

type mockExecutor struct{}

func (m mockExecutor) Query(query string, binds ...interface{}) (*sql.Rows, error) {
	return &sql.Rows{}, nil
}
func (m mockExecutor) Exec(query string, binds ...interface{}) (sql.Result, error) {
	return sqlResultMock{}, nil
}

func TestBuildErrorIfTableNotSpecified(t *testing.T) {
	_, _, err := gqb.New(mockExecutor{}).
		Where("foo", 1, gqb.Equal).
		Build(gqb.Select, nil)
	assert.Error(t, err)
}

func TestSelectBuildQuery(t *testing.T) {
	t.Run("Select() only field string", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			Table("example").
			Select("foo", "bar").
			Build(gqb.Select, nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `foo`, `bar` FROM `example`", query)
		assert.Equal(t, len(binds), 0)
	})
	t.Run("Select() contains raw field", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			Table("example").
			Select("foo", gqb.Raw("COUNT(id) AS cnt")).
			Build(gqb.Select, nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `foo`, COUNT(id) AS cnt FROM `example`", query)
		assert.Equal(t, len(binds), 0)
	})
	t.Run("Without calling Select() uses asterisk", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			Table("example").
			Build(gqb.Select, nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `*` FROM `example`", query)
		assert.Equal(t, len(binds), 0)
	})
}

func TestLimitBuildQuery(t *testing.T) {
	query, binds, err := gqb.New(mockExecutor{}).
		Table("example").
		Limit(10).
		Build(gqb.Select, nil)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `*` FROM `example` LIMIT 10", query)
	assert.Equal(t, len(binds), 0)
}

func TestOffsetBuildQuery(t *testing.T) {
	query, binds, err := gqb.New(mockExecutor{}).
		Table("example").
		Offset(10).
		Build(gqb.Select, nil)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `*` FROM `example` OFFSET 10", query)
	assert.Equal(t, len(binds), 0)
}
