package gqb_test

import (
	"database/sql"
	"testing"
	"time"

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
		Build(gqb.Select, "", nil)
	assert.Error(t, err)
}

func TestSelectBuildQuery(t *testing.T) {
	t.Run("Select() only field string", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			Select("foo", "bar").
			Build(gqb.Select, "example", nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `foo`, `bar` FROM `example`", query)
		assert.Equal(t, 0, len(binds))
	})
	t.Run("Select() contains raw field", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			Select("foo", gqb.Raw("COUNT(id) AS cnt")).
			Build(gqb.Select, "example", nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `foo`, COUNT(id) AS cnt FROM `example`", query)
		assert.Equal(t, 0, len(binds))
	})
	t.Run("Without calling Select() uses asterisk", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			Build(gqb.Select, "example", nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `*` FROM `example`", query)
		assert.Equal(t, 0, len(binds))
	})
}

func TestLimitBuildQuery(t *testing.T) {
	query, binds, err := gqb.New(mockExecutor{}).
		Limit(10).
		Build(gqb.Select, "example", nil)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `*` FROM `example` LIMIT 10", query)
	assert.Equal(t, 0, len(binds))
}

func TestOffsetBuildQuery(t *testing.T) {
	query, binds, err := gqb.New(mockExecutor{}).
		Offset(10).
		Build(gqb.Select, "example", nil)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `*` FROM `example` OFFSET 10", query)
	assert.Equal(t, 0, len(binds))
}

func TestWhereQuery(t *testing.T) {
	t.Run("Where() adds WHERE sql", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			Where("id", 1, gqb.Equal).
			Build(gqb.Select, "example", nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `*` FROM `example` WHERE (`id` = ?)", query)
		assert.Equal(t, 1, len(binds))
		v, ok := binds[0].(int)
		if !ok {
			t.Errorf("first bind parameter should be int")
			return
		}
		assert.Equal(t, 1, v)
	})
	t.Run("Multiple Where() adds WHERE sql", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			Where("id", 1, gqb.Equal).
			Where("name", "john", gqb.Equal).
			Build(gqb.Select, "example", nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `*` FROM `example` WHERE (`id` = ?) AND (`name` = ?)", query)
		assert.Equal(t, 2, len(binds))
		v, ok := binds[0].(int)
		if !ok {
			t.Errorf("first bind parameter should be int")
			return
		}
		assert.Equal(t, 1, v)
		s, ok := binds[1].(string)
		if !ok {
			t.Errorf("second bind parameter should be string")
			return
		}
		assert.Equal(t, "john", s)
	})
	t.Run("OrWhere() adds WHERE with OR expression sql", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			Where("id", 1, gqb.Equal).
			OrWhere("name", "john", gqb.Equal).
			Build(gqb.Select, "example", nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `*` FROM `example` WHERE (`id` = ?) OR (`name` = ?)", query)
		assert.Equal(t, 2, len(binds))
		v, ok := binds[0].(int)
		if !ok {
			t.Errorf("first bind parameter should be int")
			return
		}
		assert.Equal(t, 1, v)
		s, ok := binds[1].(string)
		if !ok {
			t.Errorf("second bind parameter should be string")
			return
		}
		assert.Equal(t, "john", s)
	})
	t.Run("WhereIn() adds WHERE x IN sql", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			WhereIn("id", 1, 2, 3, 4, 5).
			Build(gqb.Select, "example", nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `*` FROM `example` WHERE (`id` IN (?, ?, ?, ?, ?))", query)
		assert.Equal(t, 5, len(binds))
		for i := 0; i < len(binds); i++ {
			v, ok := binds[i].(int)
			if !ok {
				t.Errorf("bind parameter should be int")
				return
			}
			assert.Equal(t, i+1, v)
		}
	})
	t.Run("Like() adds WHERE x LIKE sql", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			Like("name", "joh%").
			Build(gqb.Select, "example", nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `*` FROM `example` WHERE (`name` LIKE ?)", query)
		assert.Equal(t, 1, len(binds))
		s, ok := binds[0].(string)
		if !ok {
			t.Errorf("first bind parameter should be string")
			return
		}
		assert.Equal(t, "joh%", s)
	})
	t.Run("WhereGroup() add group where condition within parentheses", func(t *testing.T) {
		query, binds, err := gqb.New(mockExecutor{}).
			WhereGroup(func(g *gqb.ConditionGroup) {
				g.Where("id", 1, gqb.Equal)
				g.Where("name", "John Smith", gqb.Equal)
			}).
			Build(gqb.Select, "example", nil)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `*` FROM `example` WHERE (`id` = ? AND `name` = ?)", query)
		assert.Equal(t, 2, len(binds))
		v, ok := binds[0].(int)
		if !ok {
			t.Errorf("first bind parameter should be int")
			return
		}
		assert.Equal(t, 1, v)
		s, ok := binds[1].(string)
		if !ok {
			t.Errorf("second bind parameter should be string")
			return
		}
		assert.Equal(t, "John Smith", s)
	})
}

func TestJoinQuery(t *testing.T) {
	query, binds, err := gqb.New(mockExecutor{}).
		Join("users", "id", "id", gqb.Equal).
		Where("name", "John Smith", gqb.Equal).
		Build(gqb.Select, "example", nil)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `*` FROM `example` JOIN `users` ON (`example`.`id` = `users`.`id`) WHERE (`name` = ?)", query)
	assert.Equal(t, 1, len(binds))
	if v, ok := binds[0].(string); !ok {
		t.Errorf("first bind parameter should be string")
	} else {
		assert.Equal(t, "John Smith", v)
	}
}

func TestAllCombinationSelectQuery(t *testing.T) {
	now := time.Now()
	query, binds, err := gqb.New(mockExecutor{}).
		Select("id", "name").
		Join("users", "id", "id", gqb.Equal).
		Where("register_at", now, gqb.Lt).
		WhereGroup(func(g *gqb.ConditionGroup) {
			g.Where("id", 1, gqb.Equal)
			g.Where("name", "John Smith", gqb.Equal)
		}).
		OrWhereGroup(func(g *gqb.ConditionGroup) {
			g.Where("id", 2, gqb.Equal)
			g.Where("name", "Jane Smith", gqb.Equal)
		}).
		OrderBy("register_at", gqb.Desc).
		Limit(10).
		Offset(10).
		Build(gqb.Select, "example", nil)

	assert.NoError(t, err)
	assert.Equal(t, "SELECT `id`, `name` FROM `example` JOIN `users` ON (`example`.`id` = `users`.`id`) WHERE (`register_at` < ?) AND (`id` = ? AND `name` = ?) OR (`id` = ? AND `name` = ?) ORDER BY `register_at` DESC LIMIT 10 OFFSET 10", query)
	assert.Equal(t, 5, len(binds))
	if v, ok := binds[0].(string); !ok {
		t.Errorf("first bind parameter should be string")
	} else {
		assert.Equal(t, now.Format("2006-01-02 15:04:05"), v)
	}
	if v, ok := binds[1].(int); !ok {
		t.Errorf("second bind parameter should be int")
	} else {
		assert.Equal(t, 1, v)
	}
	if v, ok := binds[2].(string); !ok {
		t.Errorf("third bind parameter should be string")
	} else {
		assert.Equal(t, "John Smith", v)
	}
	if v, ok := binds[3].(int); !ok {
		t.Errorf("fourth bind parameter should be int")
	} else {
		assert.Equal(t, 2, v)
	}
	if v, ok := binds[4].(string); !ok {
		t.Errorf("fifth bind parameter should be string")
	} else {
		assert.Equal(t, "Jane Smith", v)
	}
}

func TestInsertQuery(t *testing.T) {
	query, binds, err := gqb.New(mockExecutor{}).
		Build(gqb.Insert, "example", gqb.Data{
			"id":   1,
			"name": "John Smith",
		})
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO `example` (`id`, `name`) VALUES (?, ?)", query)
	assert.Equal(t, 2, len(binds))
	if v, ok := binds[0].(int); !ok {
		t.Errorf("first bind parameter should be int")
	} else {
		assert.Equal(t, 1, v)
	}
	if v, ok := binds[1].(string); !ok {
		t.Errorf("second bind parameter should be string")
	} else {
		assert.Equal(t, "John Smith", v)
	}
}

func TestUpdateQuery(t *testing.T) {
	now := time.Now()
	query, binds, err := gqb.New(mockExecutor{}).
		Where("id", 1, gqb.Equal).
		Build(gqb.Update, "example", gqb.Data{
			"name":       "Jane Smith",
			"updated_at": now,
		})
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE `example` SET `name` = ?, `updated_at` = ? WHERE (`id` = ?)", query)
	assert.Equal(t, 3, len(binds))
	if v, ok := binds[0].(string); !ok {
		t.Errorf("first bind parameter should be string")
	} else {
		assert.Equal(t, "Jane Smith", v)
	}
	if v, ok := binds[1].(string); !ok {
		t.Errorf("second bind parameter should be string")
	} else {
		assert.Equal(t, now.Format("2006-01-2 15:04:05"), v)
	}
	if v, ok := binds[2].(int); !ok {
		t.Errorf("third bind parameter should be int")
	} else {
		assert.Equal(t, 1, v)
	}
}
