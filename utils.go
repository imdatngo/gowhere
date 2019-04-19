package gowhere

import (
	"fmt"
	"reflect"
	"time"
)

// utilsCollection is just a empty struct for a collection of helper methods
type utilsCollection struct{}

// Utils list of predefined functions to make our life easier
var Utils = utilsCollection{}

// ToString converts given value to string
func (u utilsCollection) ToString(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return ""
	case string:
		return v
	case *string:
		return *v
	case time.Time, *time.Time:
		return u.ToDateTime(v)
	case rune:
		return string(v)
	case byte:
		return string(v)
	case []byte:
		return string(v)
	case *bool, *int, *int8, *int16, *int32, *int64, *uint, *uint16, *uint32, *uint64, *uintptr, *float32, *float64, *complex64, *complex128:
		return fmt.Sprintf("%v", reflect.ValueOf(v).Elem())
	default:
		// case bool, int, int8, int16, int32, int64, uint, uint16, uint32, uint64, uintptr, float32, float64, complex64, complex128:
		// win few ns by not using reflect.Value if it's not a ptr
		return fmt.Sprintf("%v", v)
	}
}

// ToSlice converts given value to a slice
func (u utilsCollection) ToSlice(val interface{}) interface{} {
	if val == nil {
		return []interface{}{}
	}

	vt := reflect.TypeOf(val)
	if vt.Kind() == reflect.Array || vt.Kind() == reflect.Slice {
		return val
	}

	return []interface{}{val}
}

// ToSQLVar converts given value to correct type/format for SQL
func (u utilsCollection) ToSQLVar(val interface{}) interface{} {
	switch val.(type) {
	case nil, time.Time, *time.Time, rune, byte, []byte:
		return u.ToString(val)
	default:
		return val
	}
}

// ToDateTime converts given time to datetime format for SQL
func (u utilsCollection) ToDateTime(val interface{}) string {
	switch v := val.(type) {
	case time.Time, *time.Time:
		t, ok := v.(time.Time)
		if !ok {
			tptr, _ := v.(*time.Time)
			t = *tptr
		}
		if t.Hour()+t.Minute()+t.Second()+t.Nanosecond() == 0 {
			return t.Format("2006-01-02")
		}
		return t.Format("2006-01-02 15:04:05")
	default:
		return u.ToString(val)
	}
}

// ToDate converts given time to date, i.e subtract the time to start of day
func (u utilsCollection) ToDate(val interface{}) string {
	switch v := val.(type) {
	case time.Time, *time.Time:
		t, ok := v.(time.Time)
		if !ok {
			tptr, _ := v.(*time.Time)
			t = *tptr
		}
		y, m, d := t.Date()
		date := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
		return date.Format("2006-01-02")
	default:
		return u.ToString(val)
	}
}
