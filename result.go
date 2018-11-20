package gqb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const (
	nullString  = "NullString"
	nullFloat64 = "NullFloat64"
	nullInt64   = "NullInt64"
	nullBool    = "NullBool"
	timeStruct  = "Time"
)

// Result is struct for SELECT query result mapper
type Result struct {
	// values stacks all query result column values as interface{}
	values map[string]interface{}
}

// Create Result pointer
func NewResult(values map[string]interface{}) *Result {
	return &Result{
		values: values,
	}
}

// json.Marshaller interface implementation
func (r *Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.values)
}

// Check value corresponds to field existence
func (r *Result) exists(f string) bool {
	_, ok := r.values[f]
	return ok
}

// Check field value is nil
func (r *Result) Nil(f string) bool {
	if v, ok := r.values[f]; !ok {
		return true
	} else {
		return v == nil
	}
}

// Force get field value as string
func (r *Result) MustString(f string) string {
	if s, err := r.String(f); err != nil {
		panic(err)
	} else {
		return s
	}
}

// Get field value as string with caring type conversion
func (r *Result) String(f string) (string, error) {
	v, ok := r.values[f]
	if !ok {
		return "", fmt.Errorf("field %s doesn't exist in result", f)
	}
	switch v.(type) {
	case string:
		return v.(string), nil
	}
	return fmt.Sprint(v), nil
}

// Force get field value as int
func (r *Result) MustInt(f string) int {
	if i, err := r.Int(f); err != nil {
		panic(err)
	} else {
		return i
	}
}

// Get field value as int with caring type conversion
func (r *Result) Int(f string) (int, error) {
	v, ok := r.values[f]
	if !ok {
		return 0, fmt.Errorf("field %s doesn't exist in result", f)
	}
	switch v.(type) {
	case string:
		return strconv.Atoi(v.(string))
	case int:
		return v.(int), nil
	case int8:
		return int(v.(int8)), nil
	case int16:
		return int(v.(int16)), nil
	case int32:
		return int(v.(int32)), nil
	case int64:
		return int(v.(int64)), nil
	case float32:
		return int(v.(float32)), nil
	case float64:
		return int(v.(float64)), nil
	case uint8:
		return int(v.(uint8)), nil
	case uint16:
		return int(v.(uint16)), nil
	case uint32:
		return int(v.(uint32)), nil
	case uint64:
		return int(v.(uint64)), nil
	}

	return 0, fmt.Errorf("field %s couldn't cast to int", f)
}

// Force get field value as int64
func (r *Result) MustInt64(f string) int64 {
	if i64, err := r.Int64(f); err != nil {
		panic(err)
	} else {
		return i64
	}
}

// Get field value as int64 with caring type conversion
func (r *Result) Int64(f string) (int64, error) {
	v, ok := r.values[f]
	if !ok {
		return 0, fmt.Errorf("field %s doesn't exist in result", f)
	}
	switch v.(type) {
	case string:
		return strconv.ParseInt(v.(string), 10, 64)
	case int:
		return int64(v.(int)), nil
	case int8:
		return int64(v.(int8)), nil
	case int16:
		return int64(v.(int16)), nil
	case int32:
		return int64(v.(int32)), nil
	case int64:
		return v.(int64), nil
	case float32:
		return int64(v.(float32)), nil
	case float64:
		return int64(v.(float64)), nil
	case uint8:
		return int64(v.(uint8)), nil
	case uint16:
		return int64(v.(uint16)), nil
	case uint32:
		return int64(v.(uint32)), nil
	case uint64:
		return int64(v.(uint64)), nil
	}

	return 0, fmt.Errorf("field %s couldn't cast to int64", f)
}

// Force get field value as float64
func (r *Result) MustFloat64(f string) float64 {
	if f64, err := r.Float64(f); err != nil {
		panic(err)
	} else {
		return f64
	}
}

// Get field value as float64 with caring type conversion
func (r *Result) Float64(f string) (float64, error) {
	v, ok := r.values[f]
	if !ok {
		return 0, fmt.Errorf("field %s doesn't exist in result", f)
	}
	switch v.(type) {
	case string:
		return strconv.ParseFloat(v.(string), 64)
	case int:
		return float64(v.(int)), nil
	case int8:
		return float64(v.(int8)), nil
	case int16:
		return float64(v.(int16)), nil
	case int32:
		return float64(v.(int32)), nil
	case int64:
		return float64(v.(int64)), nil
	case float32:
		return float64(v.(float32)), nil
	case float64:
		return v.(float64), nil
	case uint8:
		return float64(v.(uint8)), nil
	case uint16:
		return float64(v.(uint16)), nil
	case uint32:
		return float64(v.(uint32)), nil
	case uint64:
		return float64(v.(uint64)), nil
	}

	return 0, fmt.Errorf("field %s couldn't cast to float64", f)
}

// Force get field value as []byte
func (r *Result) MustBytes(f string) []byte {
	if v, err := r.Bytes(f); err != nil {
		panic(err)
	} else {
		return v
	}
}

// Get field value as []byte with caring type conversion
func (r *Result) Bytes(f string) ([]byte, error) {
	if s, err := r.String(f); err != nil {
		return nil, fmt.Errorf("field %s couldn't cast to []byte", f)
	} else {
		return []byte(s), nil
	}
}

// Force get field value as time.Time with date format
func (r *Result) MustDate(f string) time.Time {
	v := r.values[f]
	if t, ok := v.(time.Time); ok {
		return t
	} else {
		s := v.(string)
		t, _ := time.Parse(dateFormat, s)
		return t
	}
}

// Get field value as time.Time with caring type conversion, time parsing.
// The value must be and date format string
func (r *Result) Date(f string) (time.Time, error) {
	if v, ok := r.values[f]; !ok {
		return time.Time{}, fmt.Errorf("field %s doesn't exist in result", f)
	} else if v == nil {
		return time.Time{}, fmt.Errorf("field %s is nil", f)
	} else if t, ok := v.(time.Time); ok {
		return t, nil
	} else if s, ok := v.(string); ok {
		if t, err := time.Parse(dateFormat, s); err != nil {
			return time.Time{}, err
		} else {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("field %s couldn't cast to time.Time", f)
}

// Force get field value as time.Time with datetime format
func (r *Result) MustDatetime(f string) time.Time {
	v := r.values[f]
	if t, ok := v.(time.Time); ok {
		return t
	} else {
		s := v.(string)
		t, _ := time.Parse(datetimeFormat, s)
		return t
	}
}

// Get field value as time.Time with caring type conversion, time parsing.
// The value must be and dateitme format string
func (r *Result) Datetime(f string) (time.Time, error) {
	if v, ok := r.values[f]; !ok {
		return time.Time{}, fmt.Errorf("field %s doesn't exist in result", f)
	} else if v == nil {
		return time.Time{}, fmt.Errorf("field %s is nil", f)
	} else if t, ok := v.(time.Time); ok {
		return t, nil
	} else if s, ok := v.(string); ok {
		if t, err := time.Parse(datetimeFormat, s); err != nil {
			return time.Time{}, err
		} else {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("field %s couldn't cast to time.Time", f)
}

// Map() assigns query result into supplied struct field values
func (r *Result) Map(dest interface{}) error {
	if dest == nil {
		return fmt.Errorf("destination value must be non-nil")
	}
	v := derefValue(reflect.ValueOf(dest))
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("destination value must be a struct: %d", v.Kind())
	}
	if !v.CanSet() {
		return fmt.Errorf("destination value cannot set")
	}
	rt := v.Type()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if err := r.mapStructField(f, v.Field(i)); err != nil {
			return fmt.Errorf("failed to map value to struct field: %s, %s", f.Name, err.Error())
		}
	}
	return nil
}

// mapStructField() assigns value to struct field
func (r *Result) mapStructField(f reflect.StructField, v reflect.Value) error {
	tag, err := parseTag(string(f.Tag))
	if err != nil {
		return err
	}
	t := f.Type
	var isPtr bool
	if t.Kind() == reflect.Ptr {
		isPtr = true
		t = derefType(t)
	}
	if !v.CanSet() {
		fmt.Printf("%s is cannot set\n", f.Name)
		return nil
	}
	name, ok := tag["db"]
	// tag field doesn't exist or actual result value doesn't exist, no assign
	if !ok || !r.exists(name) {
		return nil
	}
	if err := r.assignBasicTypes(t, v, name, isPtr); err != nil {
		return err
	}
	return nil
}

// assignBasicTypes assigns value for Go's basic types
func (r *Result) assignBasicTypes(t reflect.Type, v reflect.Value, name string, isPtr bool) error {
	switch t.Kind() {
	case reflect.String:
		if s, err := r.String(name); err != nil {
			return err
		} else if isPtr {
			v.Set(reflect.ValueOf(&s))
		} else {
			v.SetString(s)
		}
	case reflect.Bool:
		if i, err := r.Int(name); err != nil {
			return err
		} else if isPtr {
			b := i > 0
			v.Set(reflect.ValueOf(&b))
		} else {
			v.SetBool(i > 0)
		}
	case reflect.Int:
		if i, err := r.Int64(name); err != nil {
			return err
		} else if isPtr {
			ii := int(i)
			v.Set(reflect.ValueOf(&ii))
		} else {
			v.SetInt(i)
		}
	case reflect.Int8:
		if i, err := r.Int64(name); err != nil {
			return err
		} else if isPtr {
			ii := int8(i)
			v.Set(reflect.ValueOf(&ii))
		} else {
			v.SetInt(i)
		}
	case reflect.Int16:
		if i, err := r.Int64(name); err != nil {
			return err
		} else if isPtr {
			ii := int16(i)
			v.Set(reflect.ValueOf(&ii))
		} else {
			v.SetInt(i)
		}
	case reflect.Int32:
		if i, err := r.Int64(name); err != nil {
			return err
		} else if isPtr {
			ii := int32(i)
			v.Set(reflect.ValueOf(&ii))
		} else {
			v.SetInt(i)
		}
	case reflect.Int64:
		if i, err := r.Int64(name); err != nil {
			return err
		} else if isPtr {
			v.Set(reflect.ValueOf(&i))
		} else {
			v.SetInt(i)
		}
	case reflect.Uint:
		if i, err := r.Int64(name); err != nil {
			return err
		} else if isPtr {
			ui := uint(i)
			v.Set(reflect.ValueOf(&ui))
		} else {
			v.SetUint(uint64(i))
		}
	case reflect.Uint8:
		if i, err := r.Int64(name); err != nil {
			return err
		} else if isPtr {
			ui := uint8(i)
			v.Set(reflect.ValueOf(&ui))
		} else {
			v.SetUint(uint64(i))
		}
	case reflect.Uint16:
		if i, err := r.Int64(name); err != nil {
			return err
		} else if isPtr {
			ui := uint16(i)
			v.Set(reflect.ValueOf(&ui))
		} else {
			v.SetUint(uint64(i))
		}
	case reflect.Uint32:
		if i, err := r.Int64(name); err != nil {
			return err
		} else if isPtr {
			ui := uint32(i)
			v.Set(reflect.ValueOf(&ui))
		} else {
			v.SetUint(uint64(i))
		}
	case reflect.Uint64:
		if i, err := r.Int64(name); err != nil {
			return err
		} else if isPtr {
			ui := uint64(i)
			v.Set(reflect.ValueOf(&ui))
		} else {
			v.SetUint(uint64(i))
		}
	case reflect.Float32:
		if i, err := r.Float64(name); err != nil {
			return err
		} else if isPtr {
			f32 := float32(i)
			v.Set(reflect.ValueOf(&f32))
		} else {
			v.SetFloat(i)
		}
	case reflect.Float64:
		if i, err := r.Float64(name); err != nil {
			return err
		} else if isPtr {
			v.Set(reflect.ValueOf(&i))
		} else {
			v.SetFloat(i)
		}
	case reflect.Struct:
		return r.assignStructType(t, v, name, isPtr)
	}
	return nil
}

// assignStructType() assigns value for struct types
func (r *Result) assignStructType(t reflect.Type, v reflect.Value, name string, isPtr bool) error {
	switch t.Name() {
	case nullString:
		i, err := r.String(name)
		nv := sql.NullString{
			String: i,
			Valid:  err == nil,
		}
		if isPtr {
			v.Set(reflect.ValueOf(&nv))
		} else {
			v.Set(reflect.ValueOf(nv))
		}
	case nullFloat64:
		i, err := r.Float64(name)
		nv := sql.NullFloat64{
			Float64: i,
			Valid:   err == nil,
		}
		if isPtr {
			v.Set(reflect.ValueOf(&nv))
		} else {
			v.Set(reflect.ValueOf(nv))
		}
	case nullInt64:
		i, err := r.Int64(name)
		nv := sql.NullInt64{
			Int64: i,
			Valid: err == nil,
		}
		if isPtr {
			v.Set(reflect.ValueOf(&nv))
		} else {
			v.Set(reflect.ValueOf(nv))
		}
	case nullBool:
		i, err := r.Int(name)
		nv := sql.NullBool{
			Bool:  i > 0,
			Valid: err == nil,
		}
		if isPtr {
			v.Set(reflect.ValueOf(&nv))
		} else {
			v.Set(reflect.ValueOf(nv))
		}
	case timeStruct:
		iv := r.values[name]
		if i, ok := iv.(time.Time); ok {
			if isPtr {
				v.Set(reflect.ValueOf(&i))
			} else {
				v.Set(reflect.ValueOf(i))
			}
		} else if s, ok := iv.(string); ok {
			if i, err := time.Parse(datetimeFormat, s); err == nil {
				if isPtr {
					v.Set(reflect.ValueOf(&i))
				} else {
					v.Set(reflect.ValueOf(i))
				}
			} else if i, err := time.Parse(dateFormat, s); err == nil {
				if isPtr {
					v.Set(reflect.ValueOf(&i))
				} else {
					v.Set(reflect.ValueOf(i))
				}
			} else if i, err := time.Parse(timeFormat, s); err == nil {
				if isPtr {
					v.Set(reflect.ValueOf(&i))
				} else {
					v.Set(reflect.ValueOf(i))
				}
			}
		}
	}
	return nil
}

// Short syntax for []*Result
type Results []*Result

// Map() assigns query result into supplied struct field values recursively
func (r Results) Map(dest interface{}) error {
	if dest == nil {
		return fmt.Errorf("destination value must be non-nil")
	}
	v := reflect.ValueOf(dest)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("destination value must be a slice")
	}
	t := v.Type()
	isPtr := t.Kind() == reflect.Ptr
	if isPtr {
		t = t.Elem()
	}
	direct := reflect.Indirect(v)
	for _, result := range r {
		row := reflect.New(t.Elem())
		if err := result.Map(row.Interface()); err != nil {
			return err
		}
		if isPtr {
			direct.Set(reflect.Append(direct, row))
		} else {
			direct.Set(reflect.Append(direct, reflect.Indirect(row)))
		}
	}
	return nil
}
