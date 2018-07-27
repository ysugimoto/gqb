package gqb_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ysugimoto/gqb"
)

func TestMarshalJSON(t *testing.T) {
	r := gqb.NewResult(map[string]interface{}{
		"string":  "foobarbaz",
		"integer": 1,
	})
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(r)
	assert.NoError(t, err)
}

func TestString(t *testing.T) {
	t.Run("MustString() returns string", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": "foobarbaz",
		})
		assert.Equal(t, "foobarbaz", r.MustString("example"))
	})
	t.Run("String() returns error for no-string type value", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": 1,
		})
		_, err := r.String("example")
		assert.Error(t, err)
	})
	t.Run("String() returns expected string without error", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": "value",
		})
		v, err := r.String("example")
		assert.NoError(t, err)
		assert.Equal(t, "value", v)
	})
}

func TestInt(t *testing.T) {
	t.Run("MustInt() returns int", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": 10,
		})
		assert.Equal(t, 10, r.MustInt("example"))
	})
	t.Run("Int() returns error for no-integer type value", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": "value",
		})
		_, err := r.Int("example")
		assert.Error(t, err)
	})
	t.Run("Int() returns expected int without error", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": 10,
		})
		v, err := r.Int("example")
		assert.NoError(t, err)
		assert.Equal(t, 10, v)
	})
}

func TestInt64(t *testing.T) {
	t.Run("MustInt64() returns int64", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": int64(10),
		})
		assert.Equal(t, int64(10), r.MustInt64("example"))
	})
	t.Run("Int64() returns error for no-64bit-integer type value", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": "value",
		})
		_, err := r.Int64("example")
		assert.Error(t, err)
	})
	t.Run("Int64() returns expected int64 without error", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": int64(10),
		})
		v, err := r.Int64("example")
		assert.NoError(t, err)
		assert.Equal(t, int64(10), v)
	})
}

func TestFloat64(t *testing.T) {
	t.Run("MustFloat64() returns float64", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": float64(10.0),
		})
		assert.Equal(t, 10.0, r.MustFloat64("example"))
	})
	t.Run("Float64() returns error for no-64bit-float type value", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": "value",
		})
		_, err := r.Float64("example")
		assert.Error(t, err)
	})
	t.Run("Float64() returns expected float64 without error", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": float64(10.0),
		})
		v, err := r.Float64("example")
		assert.NoError(t, err)
		assert.Equal(t, float64(10.0), v)
	})
}

func TestBytes(t *testing.T) {
	t.Run("MustBytes() returns string", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": "foobarbaz",
		})
		assert.Equal(t, []byte("foobarbaz"), r.MustBytes("example"))
	})
	t.Run("Bytes() returns error for no-bytes type value", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": 1,
		})
		_, err := r.Bytes("example")
		assert.Error(t, err)
	})
	t.Run("Bytes() returns expected []byte without error", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": "value",
		})
		v, err := r.Bytes("example")
		assert.NoError(t, err)
		assert.Equal(t, []byte("value"), v)
	})
}

func TestDate(t *testing.T) {
	t.Run("MustDate() returns time", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": time.Now().Format("2006-01-02"),
		})
		assert.IsType(t, time.Time{}, r.MustDate("example"))
	})
	t.Run("Date() returns error for no-date-formatted type value", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": 1,
		})
		_, err := r.Date("example")
		assert.Error(t, err)
	})
	t.Run("Date() returns expected time.Time without error", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": time.Now().Format("2006-01-02"),
		})
		v, err := r.Date("example")
		assert.NoError(t, err)
		assert.IsType(t, time.Time{}, v)
	})
}

func TestDatetime(t *testing.T) {
	t.Run("MustDatetime() returns time", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": time.Now().Format("2006-01-02 15:04:05"),
		})
		assert.IsType(t, time.Time{}, r.MustDatetime("example"))
	})
	t.Run("Datetime() returns error for no-date-formatted type value", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": 1,
		})
		_, err := r.Datetime("example")
		assert.Error(t, err)
	})
	t.Run("Datetime() returns expected time.Time without error", func(t *testing.T) {
		r := gqb.NewResult(map[string]interface{}{
			"example": time.Now().Format("2006-01-02 15:04:05"),
		})
		v, err := r.Datetime("example")
		assert.NoError(t, err)
		assert.IsType(t, time.Time{}, v)
	})
}

func TestResultMap(t *testing.T) {
	type Mapper struct {
		Str      string   `db:"str"`
		Id       int64    `db:"id"`
		Id8      int8     `db:"id8"`
		Id16     int16    `db:"id16"`
		Id32     int32    `db:"id32"`
		Uid      uint64   `db:"uid"`
		Uid8     uint8    `db:"uid8"`
		Uid16    uint16   `db:"uid16"`
		Uid32    uint32   `db:"uid32"`
		Rate     float64  `db:"rate"`
		Rate32   float32  `db:"rate32"`
		IsOK     bool     `db:"is_ok"`
		Nullable *string  `db:"nullable"`
		PStr     *string  `db:"p_str"`
		PId      *int64   `db:"p_id"`
		PId8     *int8    `db:"p_id8"`
		PId16    *int16   `db:"p_id16"`
		PId32    *int32   `db:"p_id32"`
		PUid     *uint64  `db:"p_uid"`
		PUid8    *uint8   `db:"p_uid8"`
		PUid16   *uint16  `db:"p_uid16"`
		PUid32   *uint32  `db:"p_uid32"`
		PRate    *float64 `db:"p_rate"`
		PRate32  *float32 `db:"p_rate32"`
		PIsOK    *bool    `db:"p_is_ok"`
	}

	str := "foobar"
	id := 1
	id8 := 2
	id16 := 4
	id32 := 8
	uid := 1
	uid8 := 2
	uid16 := 4
	uid32 := 8
	rate := 0.8
	rate32 := 1.8
	is_ok := 0

	createResult := func() *gqb.Result {
		return gqb.NewResult(map[string]interface{}{
			"str":      str,
			"id":       id,
			"id8":      id8,
			"id16":     id16,
			"id32":     id32,
			"uid":      uid,
			"uid8":     uid8,
			"uid16":    uid16,
			"uid32":    uid32,
			"rate":     rate,
			"rate32":   rate32,
			"is_ok":    is_ok,
			"p_str":    str,
			"p_id":     id,
			"p_id8":    id8,
			"p_id16":   id16,
			"p_id32":   id32,
			"p_uid":    uid,
			"p_uid8":   uid8,
			"p_uid16":  uid16,
			"p_uid32":  uid32,
			"p_rate":   rate,
			"p_rate32": rate32,
			"p_is_ok":  is_ok,
		})
	}

	assertion := func(m Mapper) {
		assert.Equal(t, str, m.Str)
		assert.Equal(t, int64(id), m.Id)
		assert.Equal(t, int8(id8), m.Id8)
		assert.Equal(t, int16(id16), m.Id16)
		assert.Equal(t, int32(id32), m.Id32)
		assert.Equal(t, uint64(uid), m.Uid)
		assert.Equal(t, uint8(uid8), m.Uid8)
		assert.Equal(t, uint16(uid16), m.Uid16)
		assert.Equal(t, uint32(uid32), m.Uid32)
		assert.Equal(t, float64(rate), m.Rate)
		assert.Equal(t, float32(rate32), m.Rate32)
		assert.Equal(t, false, m.IsOK)
		assert.Equal(t, str, *m.PStr)
		assert.Equal(t, int64(id), *m.PId)
		assert.Equal(t, int8(id8), *m.PId8)
		assert.Equal(t, int16(id16), *m.PId16)
		assert.Equal(t, int32(id32), *m.PId32)
		assert.Equal(t, uint64(uid), *m.PUid)
		assert.Equal(t, uint8(uid8), *m.PUid8)
		assert.Equal(t, uint16(uid16), *m.PUid16)
		assert.Equal(t, uint32(uid32), *m.PUid32)
		assert.Equal(t, float64(rate), *m.PRate)
		assert.Equal(t, float32(rate32), *m.PRate32)
		assert.Equal(t, false, *m.PIsOK)
	}

	t.Run("Map to stack struct", func(t *testing.T) {
		r := createResult()
		m := Mapper{}
		assert.NoError(t, r.Map(&m))
		assertion(m)
	})

	t.Run("Map to pointer struct", func(t *testing.T) {
		r := createResult()
		m := &Mapper{}
		assert.NoError(t, r.Map(&m))
		// assertion(*m)
	})

	// t.Run("Map to stack struct for result list", func(t *testing.T) {
	// 	rs := gqb.Results{}
	// 	for i := 0; i < 10; i++ {
	// 		rs = append(rs, createResult())
	// 	}
	// 	ms := []Mapper{}
	// 	assert.NoError(t, rs.Map(&ms))
	// 	assert.Equal(t, 10, len(ms))
	// 	for _, v := range ms {
	// 		assertion(v)
	// 	}
	// })

	// t.Run("Map to pointer struct for result list", func(t *testing.T) {
	// 	rs := gqb.Results{}
	// 	for i := 0; i < 10; i++ {
	// 		rs = append(rs, createResult())
	// 	}
	// 	ms := []*Mapper{}
	// 	assert.NoError(t, rs.Map(&ms))
	// 	assert.Equal(t, 10, len(ms))
	// 	for _, v := range ms {
	// 		assertion(*v)
	// 	}
	// })
}
