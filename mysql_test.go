package gqb_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ysugimoto/gqb"
)

func runMysqlTest(t *testing.T) {
	gqb.SetDriver("mysql")

	t.Run("Build error if table not specified", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Where("foo", 1, gqb.Equal).
			Get("")
		assert.Error(t, err)
	})

	t.Run("Table aliasing", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Get(gqb.Alias("example", "E"))
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example` AS `E`", m.query)
	})

	t.Run("Select() only field string", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Select("foo", "bar").
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT `foo`, `bar` FROM `example`", m.query)
		assert.Equal(t, 0, len(m.binds))
	})

	t.Run("Select() contains raw field", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Select("foo", gqb.Raw("COUNT(id) AS cnt")).
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT `foo`, COUNT(id) AS cnt FROM `example`", m.query)
		assert.Equal(t, 0, len(m.binds))
	})

	t.Run("Select() contains alias field", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Select("foo", gqb.Alias("example.bar", "baz")).
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT `foo`, `example`.`bar` AS `baz` FROM `example`", m.query)
		assert.Equal(t, 0, len(m.binds))
	})

	t.Run("Without calling Select() uses asterisk", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).Get("example")

		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example`", m.query)
		assert.Equal(t, 0, len(m.binds))
	})

	t.Run("Group by build query", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			GroupBy("id").
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example` GROUP BY `id`", m.query)
		assert.Equal(t, 0, len(m.binds))
	})

	t.Run("Limit build query", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Limit(10).
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example` LIMIT 10", m.query)
		assert.Equal(t, 0, len(m.binds))
	})

	t.Run("Offset build query", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Offset(10).
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example` OFFSET 10", m.query)
		assert.Equal(t, 0, len(m.binds))
	})

	t.Run("Where() adds WHERE sql", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Where("id", 1, gqb.Equal).
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example` WHERE (`id` = ?)", m.query)
		assert.Equal(t, 1, len(m.binds))
		v, ok := m.binds[0].(int)
		if !ok {
			t.Errorf("first bind parameter should be int")
			return
		}
		assert.Equal(t, 1, v)
	})

	t.Run("Multiple Where() adds WHERE sql", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Where("id", 1, gqb.Equal).
			Where("name", "john", gqb.Equal).
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example` WHERE (`id` = ?) AND (`name` = ?)", m.query)
		assert.Equal(t, 2, len(m.binds))
		v, ok := m.binds[0].(int)
		if !ok {
			t.Errorf("first bind parameter should be int")
			return
		}
		assert.Equal(t, 1, v)
		s, ok := m.binds[1].(string)
		if !ok {
			t.Errorf("second bind parameter should be string")
			return
		}
		assert.Equal(t, "john", s)
	})

	t.Run("OrWhere() adds WHERE with OR expression sql", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Where("id", 1, gqb.Equal).
			OrWhere("name", "john", gqb.Equal).
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example` WHERE (`id` = ?) OR (`name` = ?)", m.query)
		assert.Equal(t, 2, len(m.binds))
		v, ok := m.binds[0].(int)
		if !ok {
			t.Errorf("first bind parameter should be int")
			return
		}
		assert.Equal(t, 1, v)
		s, ok := m.binds[1].(string)
		if !ok {
			t.Errorf("second bind parameter should be string")
			return
		}
		assert.Equal(t, "john", s)
	})

	t.Run("WhereIn() adds WHERE x IN sql", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			WhereIn("id", 1, 2, 3, 4, 5).
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example` WHERE (`id` IN (?, ?, ?, ?, ?))", m.query)
		assert.Equal(t, 5, len(m.binds))
		for i := 0; i < len(m.binds); i++ {
			v, ok := m.binds[i].(int)
			if !ok {
				t.Errorf("bind parameter should be int")
				return
			}
			assert.Equal(t, i+1, v)
		}
	})

	t.Run("Like() adds WHERE x LIKE sql", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Like("name", "joh%").
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example` WHERE (`name` LIKE ?)", m.query)
		assert.Equal(t, 1, len(m.binds))
		s, ok := m.binds[0].(string)
		if !ok {
			t.Errorf("first bind parameter should be string")
			return
		}
		assert.Equal(t, "joh%", s)
	})

	t.Run("WhereGroup() add group where condition within parentheses", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			WhereGroup(func(g *gqb.WhereGroup) {
				g.Where("id", 1, gqb.Equal)
				g.Where("name", "John Smith", gqb.Equal)
			}).
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example` WHERE (`id` = ? AND `name` = ?)", m.query)
		assert.Equal(t, 2, len(m.binds))
		v, ok := m.binds[0].(int)
		if !ok {
			t.Errorf("first bind parameter should be int")
			return
		}
		assert.Equal(t, 1, v)
		s, ok := m.binds[1].(string)
		if !ok {
			t.Errorf("second bind parameter should be string")
			return
		}
		assert.Equal(t, "John Smith", s)
	})

	t.Run("Join() build query", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Join("users", "id", "id", gqb.Equal).
			Where("name", "John Smith", gqb.Equal).
			Get("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT * FROM `example` JOIN `users` ON (`example`.`id` = `users`.`id`) WHERE (`name` = ?)", m.query)
		assert.Equal(t, 1, len(m.binds))
		if v, ok := m.binds[0].(string); !ok {
			t.Errorf("first bind parameter should be string")
		} else {
			assert.Equal(t, "John Smith", v)
		}
	})

	t.Run("All combination select query", func(t *testing.T) {
		m := &mockExecutor{}
		now := time.Now()
		_, err := gqb.New(m).
			Select("id", "name").
			Join("users", "id", "id", gqb.Equal).
			Where("register_at", now, gqb.Lt).
			WhereGroup(func(g *gqb.WhereGroup) {
				g.Where("id", 1, gqb.Equal)
				g.Where("name", "John Smith", gqb.Equal)
			}).
			OrWhereGroup(func(g *gqb.WhereGroup) {
				g.Where("id", 2, gqb.Equal)
				g.Where("name", "Jane Smith", gqb.Equal)
			}).
			OrderBy("register_at", gqb.Desc).
			Limit(10).
			Offset(10).
			Get("example")

		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "SELECT `id`, `name` FROM `example` JOIN `users` ON (`example`.`id` = `users`.`id`) WHERE (`register_at` < ?) AND (`id` = ? AND `name` = ?) OR (`id` = ? AND `name` = ?) ORDER BY `register_at` DESC LIMIT 10 OFFSET 10", m.query)
		assert.Equal(t, 5, len(m.binds))
		if v, ok := m.binds[0].(string); !ok {
			t.Errorf("first bind parameter should be string")
		} else {
			assert.Equal(t, now.Format("2006-01-02 15:04:05"), v)
		}
		if v, ok := m.binds[1].(int); !ok {
			t.Errorf("second bind parameter should be int")
		} else {
			assert.Equal(t, 1, v)
		}
		if v, ok := m.binds[2].(string); !ok {
			t.Errorf("third bind parameter should be string")
		} else {
			assert.Equal(t, "John Smith", v)
		}
		if v, ok := m.binds[3].(int); !ok {
			t.Errorf("fourth bind parameter should be int")
		} else {
			assert.Equal(t, 2, v)
		}
		if v, ok := m.binds[4].(string); !ok {
			t.Errorf("fifth bind parameter should be string")
		} else {
			assert.Equal(t, "Jane Smith", v)
		}
	})

	t.Run("Insert query", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Insert("example", gqb.Data{
				"id":   1,
				"name": "John Smith",
			})
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "INSERT INTO `example` (`id`, `name`) VALUES (?, ?)", m.query)
		assert.Equal(t, 2, len(m.binds))
		if v, ok := m.binds[0].(int); !ok {
			t.Errorf("first bind parameter should be int")
		} else {
			assert.Equal(t, 1, v)
		}
		if v, ok := m.binds[1].(string); !ok {
			t.Errorf("second bind parameter should be string")
		} else {
			assert.Equal(t, "John Smith", v)
		}
	})

	t.Run("Update query", func(t *testing.T) {
		m := &mockExecutor{}
		now := time.Now()
		_, err := gqb.New(m).
			Where("id", 1, gqb.Equal).
			Update("example", gqb.Data{
				"name":       "Jane Smith",
				"updated_at": now,
			})
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "UPDATE `example` SET `name` = ?, `updated_at` = ? WHERE (`id` = ?)", m.query)
		assert.Equal(t, 3, len(m.binds))
		if v, ok := m.binds[0].(string); !ok {
			t.Errorf("first bind parameter should be string")
		} else {
			assert.Equal(t, "Jane Smith", v)
		}
		if v, ok := m.binds[1].(string); !ok {
			t.Errorf("second bind parameter should be string")
		} else {
			assert.Equal(t, now.Format("2006-01-2 15:04:05"), v)
		}
		if v, ok := m.binds[2].(int); !ok {
			t.Errorf("third bind parameter should be int")
		} else {
			assert.Equal(t, 1, v)
		}
	})

	t.Run("Delete query", func(t *testing.T) {
		m := &mockExecutor{}
		_, err := gqb.New(m).
			Where("id", 1, gqb.Equal).
			Delete("example")
		assert.IsType(t, mockError{}, err)
		assert.Equal(t, "DELETE FROM `example` WHERE (`id` = ?)", m.query)
		assert.Equal(t, 1, len(m.binds))
		if v, ok := m.binds[0].(int); !ok {
			t.Errorf("first bind parameter should be int")
		} else {
			assert.Equal(t, 1, v)
		}
	})
}
