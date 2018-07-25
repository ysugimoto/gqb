package gqb

import (
	"encoding/json"
	"fmt"
	"reflect"
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

func (r *Result) Nil(f string) bool {
	if v, ok := r.values[f]; !ok {
		return true
	} else {
		return v == nil
	}
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
		return "", fmt.Errorf("field %s couldn't cast to string", f)
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
	return []byte(r.MustString(f))
}

func (r *Result) Bytes(f string) ([]byte, error) {
	if s, err := r.String(f); err != nil {
		return nil, fmt.Errorf("field %s couldn't cast to float64", f)
	} else {
		return []byte(s), nil
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

func parseTag(tag string) (map[string]string, error) {
	parsed := make(map[string]string)
	var stack string
	var valueStart bool
	var key string
	for i, b := range []byte(tag) {
		switch b {
		case ':':
			if stack == "" {
				return nil, fmt.Errorf(`syntax error: unexpected ":" is present %s on %d`, tag, i)
			}
			key = stack
			stack = ""
		case '"':
			if !valueStart {
				valueStart = true
			} else {
				parsed[key] = stack
				valueStart = false
			}
			stack = ""
		case ' ':
			continue
		default:
			stack += string(b)
		}
	}
	if stack != "" {
		return nil, fmt.Errorf(`syntax error: invalid sting is remaining: %s`, tag)
	}
	return parsed, nil
}

func (r *Result) Map(dest interface{}) error {
	if dest == nil {
		return fmt.Errorf("destination value must be non-nil")
	}
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("destination value must be a struct")
	}
	rt := v.Type()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if err := r.mapStructField(f, v.Field(i)); err != nil {
			return fmt.Errorf("failed to map value to struct field: %s", f.Name)
		}
	}
	return nil
}

func (r *Result) mapStructField(f reflect.StructField, v reflect.Value) error {
	tag, err := parseTag(string(f.Tag))
	if err == nil {
		return err
	}
	if f.Type.Kind() == reflect.Ptr {
		return r.mapStructField(f, reflect.Indirect(v))
	}
	name, ok := tag["gqb"]
	if !ok {
		return nil
	}
	switch f.Type.Kind() {
	case reflect.String:
		if s, err := r.String(name); err != nil {
			return err
		} else {
			v.SetString(s)
		}
	case reflect.Bool:
		if i, err := r.Int(name); err != nil {
			return err
		} else {
			v.SetBool(i > 0)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := r.Int64(name); err != nil {
			return err
		} else {
			v.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if i, err := r.Int64(name); err != nil {
			return err
		} else {
			v.SetUint(uint64(i))
		}
	case reflect.Float32, reflect.Float64:
		if i, err := r.Float64(name); err != nil {
			return err
		} else {
			v.SetFloat(i)
		}
	}
	return nil
}

// func (r Results) Map(dest interface{}) error {
// 	if dest == nil {
// 		return fmt.Errorf("destination value must be non-nil")
// 	}
// 	v := reflect.ValueOf(dest)
// 	if v.Kind() != reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	if v.Kind() != reflect.Slice {
// 		return fmt.Errorf("destination value must be a slice")
// 	}
// 	for , result := range r {
// 		result.Map(
// 	}
// 	for i := 0; i < v.Len(); i++ {
// 		vv := v.Index(i)
// 		if vv.Kind() != reflect.Ptr {
// 			vv = vv.Elem()
// 		}
// 		rt := vv.Type()
// 		for j := 0; j < rt.NumField(); j++ {
// 			f := rt.Field(j)
// 			if err := mapStructField(f, vv.Field(j)); err != nil {
// 				return fmt.Errorf("failed to map value to struct field: %s", f.Name)
// 			}
// 		}
// 	}
// 	return nil
// }
