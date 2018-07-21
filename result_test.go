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
			"example": []byte("foobarbaz"),
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
			"example": []byte("value"),
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
