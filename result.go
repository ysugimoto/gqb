package gqb

import (
	"encoding/json"
	"fmt"
	"time"
)

type Results []*Result

type Result struct {
	values map[string]interface{}
}

func NewResult(values map[string]interface{}) *Result {
	return &Result{
		values: values,
	}
}

func (r *Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.values)
}

func (r *Result) MustString(f string) string {
	if s, err := r.String(f); err != nil {
		panic(err)
	} else {
		return s
	}
}

func (r *Result) String(f string) (string, error) {
	if v, ok := r.values[f]; !ok {
		return "", fmt.Errorf("field %s doesn't exist in result", f)
	} else if s, ok := v.(string); !ok {
		if b, err := r.Bytes(f); err != nil {
			return "", fmt.Errorf("field %s couldn't cast to string", f)
		} else {
			return string(b), nil
		}
	} else {
		return s, nil
	}
}

func (r *Result) MustInt(f string) int {
	return r.values[f].(int)
}

func (r *Result) Int(f string) (int, error) {
	if v, ok := r.values[f]; !ok {
		return 0, fmt.Errorf("field %s doesn't exist in result", f)
	} else if i, ok := v.(int); !ok {
		return 0, fmt.Errorf("field %s couldn't cast to int", f)
	} else {
		return i, nil
	}
}

func (r *Result) MustInt64(f string) int64 {
	return r.values[f].(int64)
}

func (r *Result) Int64(f string) (int64, error) {
	if v, ok := r.values[f]; !ok {
		return 0, fmt.Errorf("field %s doesn't exist in result", f)
	} else if i, ok := v.(int64); !ok {
		return 0, fmt.Errorf("field %s couldn't cast to int64", f)
	} else {
		return i, nil
	}
}

func (r *Result) MustFloat64(f string) float64 {
	return r.values[f].(float64)
}

func (r *Result) Float64(f string) (float64, error) {
	if v, ok := r.values[f]; !ok {
		return 0, fmt.Errorf("field %s doesn't exist in result", f)
	} else if i, ok := v.(float64); !ok {
		return 0, fmt.Errorf("field %s couldn't cast to float64", f)
	} else {
		return i, nil
	}
}

func (r *Result) MustBytes(f string) []byte {
	return r.values[f].([]byte)
}

func (r *Result) Bytes(f string) ([]byte, error) {
	if v, ok := r.values[f]; !ok {
		return nil, fmt.Errorf("field %s doesn't exist in result", f)
	} else if b, ok := v.([]byte); !ok {
		return nil, fmt.Errorf("field %s couldn't cast to float64", f)
	} else {
		return b, nil
	}
}

func (r *Result) MustDate(f string) time.Time {
	t, _ := time.Parse(dateFormat, r.MustString(f))
	return t
}

func (r *Result) Date(f string) (time.Time, error) {
	if v, err := r.String(f); err != nil {
		return time.Time{}, err
	} else if t, err := time.Parse(dateFormat, v); err != nil {
		return time.Time{}, fmt.Errorf("field %s couldn't cast to time.Time: %s", f, err.Error())
	} else {
		return t, nil
	}
}

func (r *Result) MustDatetime(f string) time.Time {
	t, _ := time.Parse(timeFormat, r.MustString(f))
	return t
}

func (r *Result) Datetime(f string) (time.Time, error) {
	if v, err := r.String(f); err != nil {
		return time.Time{}, err
	} else if t, err := time.Parse(timeFormat, v); err != nil {
		return time.Time{}, fmt.Errorf("field %s couldn't cast to time.Time: %s", f, err.Error())
	} else {
		return t, nil
	}
}
